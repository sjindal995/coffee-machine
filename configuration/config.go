package configuration

import (
	"github.com/pkg/errors"
	"github.com/sjindal995/coffee-machine/models"
	"github.com/spf13/viper"
)


func ParseConfig(configDirPath string) (*models.Config, error) {
	viper.AddConfigPath(configDirPath)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &models.Config{}
	machineConfigIfc := viper.Get("machine")
	if machineConfigIfc == nil {
		return nil, errors.New("Error while parsing config: machine config not present")
	}

	machineConfig, err := parseMachineConfig(machineConfigIfc)
	if err != nil {
		return nil, err
	}

	config.MachineConfig = *machineConfig
	return config, nil
}

func parseMachineConfig(machineConfigIfc interface{}) (*models.MachineConfig, error) {
	machineConfig := &models.MachineConfig{}
	machineConfig.PrepTimeInSec = 2
	machineConfigMap := machineConfigIfc.(map[string]interface{})

	if outletsConfigMapIfc, ok := machineConfigMap["outlets"]; ok {
		if nOutletsIfc, ok := outletsConfigMapIfc.(map[string]interface{})["count_n"]; ok {
			machineConfig.Outlets = int(nOutletsIfc.(float64))
		} else {
			return nil, errors.New("count_n not present in outlets config")
		}
	} else {
		return nil, errors.New("outlets not present in machine config")
	}

	if intialQuantityMapIfc, ok := machineConfigMap["total_items_quantity"]; ok {
		initialQuantityMap := intialQuantityMapIfc.(map[string]interface{})
		machineConfig.InitialQuantity = make(map[string]int, len(initialQuantityMap))
		for itemStr, quantIfc := range initialQuantityMap {
			machineConfig.InitialQuantity[itemStr] = int(quantIfc.(float64))
		}
	} else {
		return nil, errors.New("total_items_quantity not present in machine config")
	}

	if beveragesIfc, ok := machineConfigMap["beverages"]; ok {
		beveragesMap := beveragesIfc.(map[string]interface{})
		machineConfig.BeveragesComposition = make(map[string]map[string]int, len(beveragesMap))
		for beverageStr, compIfc := range beveragesMap {
			compMap := compIfc.(map[string]interface{})
			machineConfig.BeveragesComposition[beverageStr] = make(map[string]int, len(compMap))
			for itemStr, quantIfc := range compMap {
				machineConfig.BeveragesComposition[beverageStr][itemStr] = int(quantIfc.(float64))
			}
		}
	} else {
		return nil, errors.New("beverages not present in machine config")
	}

	return machineConfig, nil
}