package models

import "time"

type Beverage struct {
	Composition map[Item]int // item wise quantities required for the beverage
	PrepTime time.Duration // time taken to prepare the beverage
	Name string // name of the beverage
}
