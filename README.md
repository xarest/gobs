# godi
Golang dependencies injection framework to manage life-cycle and scope of instances of an application at runtime

## Code convention
All components have their own dependencies and life-cycle
- Initialization
- Setup/Configuration
- Start/Run
- Stop

If you want to add this life-cycle setting for a instances, please implement `gobs.IService`
```go
type Product struct {
	log *logger.Logger
	orm *orm.Orm
	s3  *s3.S3
}

var _ gobs.IService = (*Product)(nil)

func (p *Product) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Dependencies = []gobs.BlockIdentifier{
		{S: &logger.Logger{}},
		{S: &orm.Orm{}},
		{S: &s3.S3{}},
	}
	onSetup := func(ctx context.Context, dependencies []gobs.BlockIdentifier) error {
		p.log = dependencies[0].S.(*logger.Logger)
		p.orm = dependencies[1].S.(*orm.Orm)
		p.s3 = dependencies[2].S.(*s3.S3)
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}
```
then put this instance to the main thread at init step. All other components required this instance will find this instance with the same manner.
```go
sm := gobs.NewBootstrap()
sm.AddDefault(&service.Product{}, "service.Product")
```
