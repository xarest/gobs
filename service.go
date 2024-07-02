package gobs

import (
	"context"
	"sync"

	"github.com/traphamxuan/gobs/common"
	"github.com/traphamxuan/gobs/logger"
	"github.com/traphamxuan/gobs/types"
	"github.com/traphamxuan/gobs/utils"
)

type IService interface {
	// Entry point to connect service intances with the others. This method will be called at the beginning of the bootstrap process
	// to build up the dependencies between services. This method will setup `s *Service` lifecycle.
	//
	// Example:
	//
	// func (d *D) Init(ctx context.Context, s *gobs.Service) error {
	// 	s.Deps = []gobs.IService{&B{}, &C{}} // Define dependencies here
	// 	s.OnSetup = func(ctx context.Context, deps []gobs.IService, extraDeps []gobs.CustomService) error {
	// 		// After B & C finish setting up, this callback will be called
	// 		d.B = deps[0].(*B)
	// 		d.C = deps[1].(*C)
	// 		// Other custom setup/configration go here
	// 		return nil
	// 	}
	// 	s.AsyncMode[common.StatusSetup] = true // This line will make OnSetup method be called in concurrent context without blocking others.
	// 	return nil
	// }
	Init(ctx context.Context, s *Service) error
}

type ServiceLifeCycle struct {
	// OnSetup is a callback function that will be called when the main context invokes the bootstrap.Setup(...) method.
	// This method is used to assign the dependencies instances which has setup successfully from gobs to the service instance.
	// The `deps` parameter is a list of dependencies that the service instance depends on. The `extraDeps` parameter is a list of
	// custom dependencies in case service don't share dependencies with the others.
	//
	// Example:
	//
	// sb.OnSetup = func(ctx context.Context, deps []gobs.IService, extraDeps []gobs.CustomService) error {
	// 	o.log = deps[0].(*logger.Logger)
	// 	config := deps[1].(*config.Configuration)
	// 	var dbConfig DatabaseConfig
	// 	if err := config.ParseConfig(&dbConfig); err != nil {
	// 		return err
	// 	}
	// 	o.config = &dbConfig
	// 	return o.Setup(ctx)
	// }
	// func (o *Gorm) Setup(ctx context.Context) error {
	// 	o.log.Debug("Connecting to database")
	// 	dsn := fmt.Sprintf(
	// 		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
	// 		o.config.Username,
	// 		o.config.Password,
	// 		o.config.Host,
	// 		o.config.Port,
	// 		o.config.DBName,
	// 		o.config.SSLMode,
	// 	)
	// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 		Logger: zapgorm2.New(o.log.Desugar()),
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	o.Db = db
	// 	return nil
	// }
	OnSetup func(ctx context.Context, deps []IService, extraDeps []CustomService) error

	// OnStart is a callback function that will be called when the main context invokes the bootstrap.Start(...) method.
	// Service instances passing OnStart method will be marked as `common.StatusStart` status.
	//
	// Example:
	//
	// func (o *Gorm) Init(ctx context.Context, sb *gobs.Service) error {
	// 	sb.OnStart = func(ctx context.Context) error {
	// 		return o.Start(ctx)
	// 	}
	// }
	//
	// func (o *Gorm) Start(c context.Context) error {
	// 	db, err := o.Db.DB()
	// 	if err != nil {
	// 		o.log.Error("Failed to get database connection", zap.Error(err))
	// 		return err
	// 	}
	// 	err = db.Ping()
	// 	if err != nil {
	// 		o.log.Error("Database connection failed", zap.Error(err))
	// 		return err
	// 	}
	// 	o.log.Info("Database connected")
	// 	return nil
	// }
	OnStart func(context.Context) error

	// OnStop is a callback function that will be called when the main context invokes the bootstrap.Stop(...) method.
	// Only services passing OnSetup method (return nil) are able to be invoked this method in stopping phase.
	//
	// Example:
	//
	// func (o *Gorm) Init(ctx context.Context, sb *gobs.Service) error {
	// 	sb.OnStop = func(ctx context.Context) error {
	// 		return o.Stop(ctx)
	// 	}
	// }
	//
	// func (o *Gorm) Stop(c context.Context) error {
	// 	db, err := o.Db.DB()
	// 	if err != nil {
	// 		o.log.Error("Failed to get database connection", zap.Error(err))
	// 		return err
	// 	}
	// 	if err := db.Close(); err != nil {
	// 		o.log.Error("Failed to close database connection", zap.Error(err))
	// 		return err
	// 	}
	// 	o.log.Info("Database connection closed")
	// 	return nil
	// }
	OnStop func(context.Context) error

	// Deps is a list of dependencies that the service instance depends on. This list is just a reference to the type struct.
	// Gobs will automatically look up the existing instances or create new instances based on the type struct.
	// Then set dependencies to the service instance after Init(...) method returns nil.
	Deps []IService

	// ExtraDeps is a list of custom dependencies that the service instance depends on.
	// It provides more information about the instance that the service instance depends on.
	//
	// Example:
	//
	// func (d *D) Init(ctx context.Context, s *gobs.Service) error {
	// 	s.ExtraDeps = []gobs.CustomService{
	// 		{&B{}, "", instanceB}, // D depends on instanceB, type struct B{} with key is default
	// 		{&C{}, "C1", nil}, // D depends on an instance C which has key is C1 in gobs
	// 	}
	// }
	ExtraDeps []CustomService

	// AsyncMode is a map of service status and boolean value. If the value is true, the service instance will be run in parallel goroutine context.
	// Otherwise, the service instance will be run in sequential context.
	//
	// Example:
	//
	// func (d *D) Init(ctx context.Context, s *gobs.Service) error {
	// 	s.AsyncMode = map[common.ServiceStatus]bool {
	// 		common.StatusInit:  false, // Init(...) will be called in sequential context
	// 		common.StatusSetup: true,  // Setup(...) will be called in separate goroutine context
	// 		common.StatusStart: true,  // Start(...) will be called in separate goroutine context
	// 		common.StatusStop:  false, // Stop(...) will be called in sequential context
	// 	}
	// }
	AsyncMode map[common.ServiceStatus]bool
}

type CustomService struct {
	Service  IService
	Name     string
	Instance interface{}
}

type Service struct {
	ServiceLifeCycle
	*logger.Logger
	following []types.ITask
	followers []types.ITask
	instance  IService
	name      string
	status    common.ServiceStatus
	mutex     map[common.ServiceStatus]*sync.Mutex
}

var _ types.ITask = (*Service)(nil)

func (sb *Service) Name() string {
	return sb.name
}

func NewService(s IService, name string, status common.ServiceStatus, log *logger.Logger) *Service {
	c := &Service{
		ServiceLifeCycle: ServiceLifeCycle{
			AsyncMode: make(map[common.ServiceStatus]bool, common.StatusStop+1),
			OnSetup:   func(_ context.Context, _ []IService, _ []CustomService) error { return nil },
			OnStart:   utils.EmptyFunc,
			OnStop:    utils.EmptyFunc,
		},
		Logger:   log,
		instance: s,
		name:     name,
		status:   status,
		mutex: map[common.ServiceStatus]*sync.Mutex{
			common.StatusUninitialized: {},
			common.StatusInit:          {},
			common.StatusSetup:         {},
			common.StatusStart:         {},
			common.StatusStop:          {},
		},
	}
	c.AddTag("Service/" + name)
	return c
}

func (sb *Service) DependOn(ss common.ServiceStatus) []types.ITask {
	if ss >= common.StatusStop {
		return sb.followers
	}
	return sb.following
}

func (sb *Service) Followers(ss common.ServiceStatus) []types.ITask {
	if ss >= common.StatusStop {
		return sb.following
	}
	return sb.followers
}

func (sb *Service) Run(ctx context.Context, ss common.ServiceStatus) (err error) {
	mutex, ok := sb.mutex[ss]
	if ok && mutex != nil {
		mutex.Lock()
		defer mutex.Unlock()
	}
	switch ss {
	case common.StatusSetup:
		err = sb.OnSetup(ctx, sb.Deps, sb.ExtraDeps)
	case common.StatusStart:
		err = sb.OnStart(ctx)
	case common.StatusStop:
		err = sb.OnStop(ctx)
	default:
		err = nil
	}

	if err != nil {
		return err
	}
	sb.status = ss
	return nil
}

func (sb *Service) IsRunAsync(ss common.ServiceStatus) bool {
	return sb.AsyncMode[ss]
}
