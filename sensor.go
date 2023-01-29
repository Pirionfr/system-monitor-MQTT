package main

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
	"log"
	"path/filepath"
)

type Sensor struct {
	ID    string
	Name  string
	Unit  string
	Class string
	Value string
}

const luaSensorTypeName = "sensor"

func (s *Sensor) FillData(v string) {
	L := lua.NewState()
	L.SetGlobal(luaSensorTypeName, luar.New(L, s))
	defer L.Close()
	if err := L.DoFile(fmt.Sprintf("%s.lua", filepath.Join(conf.GetString(ConfigPath), v))); err != nil {
		log.Fatal(err)
	}
}
