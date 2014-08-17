inject
======

Injection Points

## Example

Create the module with the various configuration points.

```
type Address struct {}
type Host struct {}

type NetworkModule struct {
    *inject.Module
}

func NewNetworkModule() *NetworkModule {
    return &NetworkModule{
        inject.NewModule(func(m inject.IModule) {
            m.Bind(Address{}).ToInstance("127.0.0.1")
            m.Bind(Host{}).ToInstance(8080)
        }),
    }
}
```

Now use the injection points from the module.

```
type Dispatcher struct {
    host string
    port int
}

func NewDispatcher() *Dispatcher {
    return &Dispatcher{
        host: inject.As(Address{}).GetOrElse(inject.Constant("x.x.x.x")).(string),
        port: inject.As(Host{}).GetOrElse(inject.Constant(80)).(int),
    }
}

func (c Dispatcher) Create() inject.Any {
    return NewDispatcher()
}

func (c Dispatcher) Load() []byte {
    // Handle errors.
    resp, _ := http.Get(fmt.Sprintf("%s:%v", c.host, c.port))
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}
```

Run the example:

```
func main() {
    module := inject.Add(NewNetworkModule())
    dispatcher := module.GetInstance(Dispatcher{})

    result := dispatcher.Map(func(x inject.Any) inject.Any {
        return x.(*Dispatcher).Load().(string)
    })

    fmt.Println(result)
}
```