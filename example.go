package main

import (
	"fmt"

	"github.com/SimonRichardson/inject/inject"
)

type Hello struct {
	*inject.Module
}

func NewHello() *Hello {
	return &Hello{
		inject.NewModule(func(m inject.IModule) {
			m.Bind("string").ToInstance("Hello")
		}),
	}
}

type World struct {
	message string
}

func NewWorld() *World {
	return &World{
		message: inject.As("string").GetOrElse("Bad").(string),
	}
}

func (c World) Create() *World {
	return NewWorld()
}

func (c World) String() string {
	return fmt.Sprintf("%s, world!", c.message)
}

func main() {
	hello := inject.Add(NewHello())
	world := hello.GetInstance(NewWorld())
	res := world.Map(func(x inject.Any) inject.Any {
		return x.(*World).String()
	})
	fmt.Println(res)
}
