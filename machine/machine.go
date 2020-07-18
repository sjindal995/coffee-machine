package machine

import (
	"fmt"
	"github.com/sjindal995/coffee-machine/models"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)
type Machine struct {
	Outlets *semaphore.Weighted //semaphore correspondign to number of outlets
	Beverages map[string]models.Beverage //map containing beverage name to it's model
	Quantity // map of quantity of items
	submissionChan chan models.Beverage // channel to receive beverage requests
}

// map along with mutex containing item quantities present currently in the machine
// lock to avoid race conditions
type Quantity struct {
	QuantityMap map[models.Item]int
	sync.RWMutex
}


func (m *Machine) checkAndUpdateQuantities(beverage models.Beverage) error {
	m.Quantity.Lock()
	defer m.Quantity.Unlock()

	for item, reqQuant := range beverage.Composition {
		availQuant, avail := m.QuantityMap[item]
		if !avail {
			return fmt.Errorf("%v cannot be prepared because %v is not available at %v", beverage.Name, item, time.Now())
		}
		if availQuant < reqQuant {
			return fmt.Errorf("%v cannot be prepared because %v is not sufficient at %v", beverage.Name, item, time.Now())
		}
	}

	for item, reqQuant := range beverage.Composition {
		m.QuantityMap[item] = m.QuantityMap[item] - reqQuant
	}

	return nil
}

func (m *Machine) getBeverage(beverage models.Beverage) error {
	outletAq := m.Outlets.TryAcquire(1)
	if !outletAq {
		return fmt.Errorf("%v cannot be prepared because all outlets are busy at %v", beverage.Name, time.Now())
	}
	defer m.Outlets.Release(1)

	err := m.checkAndUpdateQuantities(beverage)
	if err != nil {
		return err
	}

	time.Sleep(beverage.PrepTime)
	return nil
}

func (m *Machine) runSubmissionChan() {
	for {
		beverage := <-m.submissionChan
		go func() {
			err := m.getBeverage(beverage)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%v is prepared at %v\n", beverage.Name, time.Now())
			}
		}()
	}
}

// Submit beverage request
func (m *Machine) SubmitRequirement(beverageStr string) {
	beverage, ok := m.Beverages[beverageStr]
	if !ok {
		fmt.Printf("%v beverage not available at %v\n", beverageStr, time.Now())
		return
	}
	m.submissionChan <- beverage
}

// Update quantity of an item in the machine
func (m *Machine) UpdateItemQuant(itemStr string, additionalQuant int) {
	item := models.Item(itemStr)
	m.Quantity.Lock()
	defer m.Quantity.Unlock()
	if _, ok := m.QuantityMap[item]; !ok {
		m.QuantityMap[item] = additionalQuant
	} else {
		m.QuantityMap[item] = m.QuantityMap[item] + additionalQuant
	}
}

func InitMachine(config models.MachineConfig) Machine {
	beverages := make(map[string]models.Beverage, len(config.BeveragesComposition))
	for beverageStr, compMap := range config.BeveragesComposition {
		parsedCompMap := make(map[models.Item]int, len(compMap))
		for itemStr, quant := range compMap {
			parsedCompMap[models.Item(itemStr)] = quant
		}
		beverage := models.Beverage{Composition: parsedCompMap, Name: beverageStr, PrepTime: time.Duration(config.PrepTimeInSec) * time.Second}
		beverages[beverageStr] = beverage
	}

	parsedInitialQuantity := make(map[models.Item]int, len(config.InitialQuantity))
	for itemStr, quant := range config.InitialQuantity {
		parsedInitialQuantity[models.Item(itemStr)] = quant
	}

	machine := Machine{
		Outlets: semaphore.NewWeighted(int64(config.Outlets)),
		Beverages: beverages,
		Quantity: Quantity{QuantityMap: parsedInitialQuantity},
		submissionChan: make(chan models.Beverage, config.Outlets),
	}

	go func() {
		machine.runSubmissionChan()
	}()

	return machine
}