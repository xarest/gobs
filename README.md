# Go Bootstrap (GoBs)
Golang dependencies injection framework to manage life-cycle and scope of instances of an application at runtime

## Code convention
All components have their own dependencies and life-cycle
- Initialization
- Setup/Configuration
- Start/Run
- Stop

If you want to add this life-cycle setting for a instances, please implement `gobs.IService`
```go
type D struct {
	B *B
	C *C
}

var _ gobs.IService = (*D)(nil)

func (d *D) Init(ctx context.Context, co *gobs.Component) error {
	co.Deps = []gobs.IService{&B{}, &C{}} // Define dependencies here
	onSetup := func(ctx context.Context, deps []gobs.IService, extraDeps []gobs.CustomService) error {
		// After B & C finish setting up, this callback will be called
		d.B = deps[0].(*B)
		d.C = deps[1].(*C)
		// Other custom setup/configration go here
		return nil
	}
	co.OnSetup = &onSetup
	return nil
}
```
then put this instance to the main thread at init step. All other components required this instance will find this instance with the same manner.
```go
sm := gobs.NewBootstrap()
bs.AddDefault(&D{})
bs.Setup(context.BackGround())
```
With dependencies:
- B -> A
- C -> A, B
- D -> B, C

Output log will be
```
Service D is added
Service B is added
Service C is added
Service A is added
Service A is notifying 0 followers
Service A setup successfully
Service B is waiting for A
Service B is done waiting for A
Service B is notifying 0 followers
Service B setup successfully
Service D is waiting for B
Service D is done waiting for B
Service D is waiting for C
Service C is waiting for A
Service C is done waiting for A
Service C is waiting for B
Service C is done waiting for B
Service C is notifying 1 followers
Service C setup successfully
Service D is done waiting for C
Service D is notifying 0 followers
Service D setup successfully
```
