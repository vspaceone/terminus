package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func getConfiguration() (interface{}, error) {
	var config map[string]interface{}

	dat, err := ioutil.ReadFile("config.json")
	if err != nil {
		return config, errors.New("Couldn't read configuration file")
	}
	fmt.Println(string(dat))

	err = json.Unmarshal(dat, &config)
	if err != nil {
		return config, errors.New("Couldn't unmarshal configuration file")
	}

	return config, nil
}
