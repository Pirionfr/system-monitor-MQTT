package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

const HAExpireAfter = 60

type Daemon struct {
	MQTT         *MQTTClient
	ClientID     string
	Prefix       string
	UpdatePeriod int
	HaEnable     bool
	Sensors      []string
}

func NewDaemon(config *viper.Viper) *Daemon {
	d := &Daemon{
		MQTT:         NewMQTTClient(config.GetString(MQTTBroker), config.GetString(MQTTUser), config.GetString(MQTTPassword), config.GetString(ClientID)),
		UpdatePeriod: config.GetInt(UpdatePeriod),
		ClientID:     config.GetString(ClientID),
		Prefix:       fmt.Sprintf("%s/%s", config.GetString(Prefix), config.GetString(ClientID)),
		HaEnable:     config.GetBool(HA),
		Sensors:      config.GetStringSlice(Sensors),
	}
	return d
}

func (d *Daemon) run() {
	haRegister := false
	CycleCount := 0

	go func() {
		t := time.NewTicker(time.Duration(d.UpdatePeriod) * time.Second)
		for range t.C {
			d.NotifyState(true)

			if CycleCount >= d.UpdatePeriod/HAExpireAfter && d.HaEnable {
				d.HARegister()
				haRegister = true
				CycleCount = 0
			}

			for _, v := range d.Sensors {
				sensor := &Sensor{
					ID: v,
				}
				sensor.FillData(v)
				d.MQTT.Publish(fmt.Sprintf("%s/%s/%s", d.Prefix, sensor.Class, sensor.ID), sensor.Value)
			}

			if haRegister {
				haRegister = false
			}
			CycleCount += CycleCount
		}
	}()
}

type State struct {
	Name        string `json:"name"`
	Class       string `json:"class"`
	Device      Device `json:"device"`
	ExpireAfter int32  `json:"expire_after"`
	StateTopic  string `json:"state_topic"`
	UniqueId    string `json:"unique_id"`
}

type Device struct {
	Name        string `json:"name"`
	Identifiers string `json:"identifiers"`
}

func (d *Daemon) HARegisterState() {
	topic := fmt.Sprintf("homeassistant/binary_sensor/%s_state/config", d.ClientID)

	state := State{
		Name:  fmt.Sprintf("%s State", d.ClientID),
		Class: "connectivity",
		Device: Device{
			Name:        d.ClientID,
			Identifiers: d.ClientID,
		},
		ExpireAfter: HAExpireAfter,
		StateTopic:  fmt.Sprintf("%s/state", d.Prefix),
		UniqueId:    fmt.Sprintf("%s_state", d.ClientID),
	}

	bMsg, _ := json.Marshal(state)

	d.MQTT.Publish(topic, string(bMsg))
}

type SensorMsg struct {
	Name              string `json:"name"`
	Class             string `json:"class"`
	UnitOfMeasurement string `json:"unit_of_measurement,omitempty"`
	Device            Device `json:"device"`
	ExpireAfter       int32  `json:"expire_after"`
	StateTopic        string `json:"state_topic"`
	UniqueId          string `json:"unique_id"`
}

func (d *Daemon) HARegisterSensor(sensor *Sensor) {
	topic := fmt.Sprintf("homeassistant/sensor/%s_%s/config", d.ClientID, sensor.ID)

	state := SensorMsg{
		Name:  fmt.Sprintf("%s %s", d.ClientID, sensor.Name),
		Class: sensor.Class,
		Device: Device{
			Name:        d.ClientID,
			Identifiers: d.ClientID,
		},
		ExpireAfter: HAExpireAfter,
		StateTopic:  fmt.Sprintf("%s/%s/%s", d.Prefix, sensor.Class, sensor.ID),
		UniqueId:    fmt.Sprintf("%s_%s", d.ClientID, sensor.ID),
	}
	if sensor.Unit != "" {
		state.UnitOfMeasurement = sensor.Unit
	}

	bMsg, _ := json.Marshal(state)
	d.MQTT.Publish(topic, string(bMsg))
}

func (d *Daemon) HARegister() {
	d.HARegisterState()
	for _, v := range d.Sensors {
		sensor := &Sensor{
			ID: v,
		}
		sensor.FillData(v)
		d.HARegisterSensor(sensor)
	}
}

func (d *Daemon) NotifyState(state bool) {
	topic := fmt.Sprintf("%s/state", d.Prefix)
	message := "OFF"
	if state {
		message = "ON"
	}
	d.MQTT.Publish(topic, message)
}
