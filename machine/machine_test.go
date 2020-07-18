package machine

import (
	"github.com/sjindal995/coffee-machine/configuration"
	"testing"
	"time"
)

func setupMachine() Machine {
	config, err := configuration.ParseConfig("../configs")
	if err != nil {
		panic(err)
	}

	return InitMachine(config.MachineConfig)

}

func TestFunctional(t *testing.T) {
	m := setupMachine()
	m.SubmitRequirement("hot_tea")
	m.SubmitRequirement("hot_coffee")

	time.Sleep(3 * time.Second)
}


func TestItemInsufficiency(t *testing.T) {
	m := setupMachine()
	m.SubmitRequirement("hot_tea")
	m.SubmitRequirement("hot_coffee")
	m.SubmitRequirement("black_tea")
	time.Sleep(3 * time.Second)
}

func TestItemRefill(t *testing.T) {
	m := setupMachine()
	m.SubmitRequirement("hot_tea")
	m.SubmitRequirement("hot_coffee")
	time.Sleep(100 * time.Millisecond)
	m.UpdateItemQuant("hot_water", 100)
	m.UpdateItemQuant("sugar_syrup", 10)
	m.SubmitRequirement("black_tea")
	time.Sleep(3 * time.Second)
}

func TestOutletUnavailability(t *testing.T) {
	m := setupMachine()
	m.SubmitRequirement("hot_tea")
	m.SubmitRequirement("hot_coffee")
	time.Sleep(100 * time.Millisecond)
	m.UpdateItemQuant("hot_water", 300)
	m.UpdateItemQuant("sugar_syrup", 20)
	m.UpdateItemQuant("tea_leaves_syrup", 30)
	m.UpdateItemQuant("hot_milk", 400)
	m.SubmitRequirement("black_tea")
	m.SubmitRequirement("hot_tea")

	time.Sleep(3 * time.Second)
	m.SubmitRequirement("hot_tea")

	time.Sleep(3 * time.Second)
}

func TestUnavailableItem(t *testing.T) {
	m := setupMachine()
	m.SubmitRequirement("green_tea")
	time.Sleep(1 * time.Second)
}

func TestUnavailableBeverage(t *testing.T) {
	m := setupMachine()
	m.SubmitRequirement("herbal_tea")
	time.Sleep(1 * time.Second)
}