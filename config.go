package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type ConfigS struct {
	ListenOn        string
	Port            string
	DevMode         bool
	MySQL           string
	SessionSecret   string
	SessionDuration string
}

func LoadConfigFile() ConfigS {
	var dat ConfigS

	path := os.Getenv("CONFIG_FILE")
	if path == "" {
		path = "config.json"
	}

	byt, err := ioutil.ReadFile(path)
	if err != nil {
		Log.FatalF("Failed to open '%s' file: %s", path, err.Error())
	}
	err = json.Unmarshal(byt, &dat)
	if err != nil {
		Log.FatalF("Failed to parse '%s' file: %s", path, err.Error())
	}

	if dat.Port == "" {
		Log.FatalF("Config file '%s' MUST set the port", path)
	}

	if dat.SessionSecret == "" {
		Log.FatalF("Config file '%s' MUST set SessionSecret", path)
	}
	if dat.SessionDuration == "" {
		dat.SessionDuration = "1h"
	}
	SessionDuration, err = time.ParseDuration(dat.SessionDuration)
	panicIfErr(err)

	if !strings.Contains(dat.MySQL, "?parseTime=true") {
		Log.Fatal("Config file '%s' MUST include '?parseTime=true' when connecting to DB", path)
	}

	return dat
}
