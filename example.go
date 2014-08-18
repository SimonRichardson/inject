package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/SimonRichardson/inject/inject"
)

type Scheme struct{}
type Host struct{}
type Port struct{}
type Path struct{}

type GoogleModule struct {
	*inject.Module
}

func NewGoogleModule() *GoogleModule {
	return &GoogleModule{
		inject.NewModule(func(m inject.IModule) {
			m.Bind(Scheme{}).ToInstance("https")
			m.Bind(Host{}).ToInstance("google.co.uk")
			m.Bind(Port{}).ToInstance(443)
			m.Bind(Path{}).ToInstance("/robots.txt")
		}),
	}
}

type GithubModule struct {
	*inject.Module
}

func NewGithubModule() *GithubModule {
	return &GithubModule{
		inject.NewModule(func(m inject.IModule) {
			m.Bind(Scheme{}).ToInstance("https")
			m.Bind(Host{}).ToInstance("github.com")
			m.Bind(Port{}).ToInstance(443)
			m.Bind(Path{}).ToInstance("/robots.txt")
		}),
	}
}

type Dispatcher struct {
	scheme string
	host   string
	port   int
	path   string
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		scheme: inject.As(Scheme{}).GetOrElse(inject.Constant("http")).(string),
		host:   inject.As(Host{}).GetOrElse(inject.Constant("x.x.x.x")).(string),
		port:   inject.As(Port{}).GetOrElse(inject.Constant(80)).(int),
		path:   inject.As(Path{}).GetOrElse(inject.Constant("error.txt")).(string),
	}
}

func (c Dispatcher) Get() inject.Any {
	return NewDispatcher()
}

func (c Dispatcher) Load() []byte {
	// Handle errors.
	resp, _ := http.Get(fmt.Sprintf("%s://%s:%v%s", c.scheme, c.host, c.port, c.path))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func main() {
	x := inject.Add(NewGoogleModule())
	y := inject.Add(NewGithubModule())

	a := x.GetInstance(Dispatcher{})
	b := y.GetInstance(Dispatcher{})

	a.Map(func(x inject.Any) inject.Any {
		fmt.Println(x)
		return string(x.(*Dispatcher).Load())
	})
	b.Map(func(x inject.Any) inject.Any {
		fmt.Println(x)
		return string(x.(*Dispatcher).Load())
	})

	//fmt.Println(res0)
	//fmt.Println(res1)
}
