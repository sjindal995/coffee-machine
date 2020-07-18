package main

import (
	"github.com/sjindal995/coffee-machine/configuration"
	"github.com/sjindal995/coffee-machine/machine"
	"time"
)

func main() {
	//configPath := os.Args[1]

	config, err := configuration.ParseConfig("configs")
	if err != nil {
		panic(err)
	}

	m := machine.InitMachine(config.MachineConfig)

	m.SubmitRequirement("hot_tea")
	m.SubmitRequirement("hot_tea")
	m.SubmitRequirement("hot_tea")
	time.Sleep(100 * time.Millisecond)
	m.UpdateItemQuant("hot_water", 500)
	m.SubmitRequirement("hot_tea")
	//m.SubmitRequirement("hot_tea")
	time.Sleep(10 * time.Second)
	//m.SubmitRequirement("black_tea")
	//m.SubmitRequirement("hot_tea")
	m.SubmitRequirement("hot_coffee")

	time.Sleep(1000*time.Second)
}
