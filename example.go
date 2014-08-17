package main

import (
	"fmt"

	"github.com/SimonRichardson/inject/inject"
)

type String struct{}

type Character struct{}

type Hello struct {
	*inject.Module
}

func NewHello() *Hello {
	return &Hello{
		inject.NewModule(func(m inject.IModule) {
			m.Bind(String{}).ToInstance("Hello")
		}),
	}
}

type World struct {
	message string
}

func NewWorld() *World {
	return &World{
		message: inject.As(String{}).GetOrElse("Bad").(string),
	}
}

func (c World) Create() inject.Any {
	return NewWorld()
}

func (c World) String() string {
	return fmt.Sprintf("%s, world!", c.message)
}

func main() {
	hello := inject.Add(NewHello())
	world := hello.GetInstance(World{})

	res := world.Map(func(x inject.Any) inject.Any {
		return x.(*World).String()
	})

	fmt.Println(res)
}
