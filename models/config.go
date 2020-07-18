package models

type Config struct {
	MachineConfig
}

type MachineConfig struct {
	Outlets int
	InitialQuantity map[string]int
	BeveragesComposition map[string]map[string]int
	PrepTimeInSec int
}
