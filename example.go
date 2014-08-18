package main

import (
	"crypto/md5"
	"encoding/hex"
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

	res0 := a.Map(func(x inject.Any) inject.Any {
		return string(x.(*Dispatcher).Load())
	})
	res1 := b.Map(func(x inject.Any) inject.Any {
		return string(x.(*Dispatcher).Load())
	})

	// These two values should not match!
	fmt.Println(res0.Map(md5Hash))
	fmt.Println(res1.Map(md5Hash))
}

func md5Hash(v inject.Any) inject.Any {
	hash := md5.Sum([]byte(v.(string)))
	return hex.EncodeToString(hash[:])
}
