# Go Bootstrap (GoBs)
Golang dependencies injection framework to manage life-cycle and scope of instances of an application at runtime.

![](gobs-run-13-instances-async.gif)

Documentation website at [gobs.xarest.com](https://gobs.xarest.com)

## Code convention
All components have their own dependencies and life-cycle
- Initialization
- Setup/Configuration
- Start/Run
- Stop

If you want to add this life-cycle setting for a instances, please implement `gobs.IServiceInit`
```go
type D struct {
	B *B
	C *C
}

var _ gobs.IServiceInit = (*D)(nil)

func (d *SD) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(B), new(C)},
		common.StatusSetup: true,
	}, nil
}

var _ gobs.IServiceSetup = (*D)(nil)

func (d *D) Setup(ctx context.Context, deps gobs.Dependencies) error {
	if err := deps.Assign(&s.s2, &s.s3); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	return nil
}

```
then put this instance to the main thread at init step. All other components required this instance will find this instance with the same manner.
```go
ctx := context.BackGround()
sm := gobs.NewBootstrap()
bs.AddDefault(&D{})
bs.StartBootstrap(ctx)
```
With dependencies:
- B -> A
- C -> A, B
- D -> B, C

Output log will be
```
Service test_test.D is added with status Uninitialized
Service test_test.B is added with status Uninitialized
Service test_test.C is added with status Uninitialized
Service test_test.A is added with status Uninitialized
EXECUTE Init WITH 4 SERVICES
Service test_test.A Init successfully
Service test_test.B Init successfully
Service test_test.C Init successfully
Service test_test.D Init successfully
EXECUTE Setup WITH 4 SERVICES
Service test_test.A Setup successfully
Service test_test.B Setup successfully
Service test_test.C Setup successfully
Service test_test.D Setup successfully
EXECUTE Start WITH 4 SERVICES
Service test_test.A Start successfully
Service test_test.B Start successfully
Service test_test.C Start successfully
Service test_test.D Start successfully
EXECUTE Stop WITH 4 SERVICES
Service test_test.D Stop successfully
Service test_test.C Stop successfully
Service test_test.B Stop successfully
Service test_test.A Stop successfully
```
