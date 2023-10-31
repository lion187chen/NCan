package main

import (
	"encoding/json"
	"os"
)

type config struct {
	setting
}

type setting struct {
	Server string      `json:"server"`
	User   string      `json:"user"`
	Passwd string      `json:"passwd"`
	Name   string      `json:"name"`
	To     string      `json:"to"`
	From   string      `json:"from"`
	Driver interface{} `json:"driver,omitempty"`
}

func (my *config) init(name string) *config {
	fb, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fb, &my.setting)
	if err != nil {
		panic(err)
	}
	return my
}
