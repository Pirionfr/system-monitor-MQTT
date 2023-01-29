package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
)

const DefaultQos = 0

type MQTTClient struct {
	Client mqtt.Client
}

func NewMQTTClient(broker, user, password, clientID string) *MQTTClient {
	c := &MQTTClient{}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", broker))
	opts.SetClientID(clientID)

	if user != "" {
		opts.SetUsername(user)
	}
	if password != "" {
		opts.SetPassword(password)
	}

	c.Client = mqtt.NewClient(opts)
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	return c
}

func (c *MQTTClient) Publish(topic, message string) {
	token := c.Client.Publish(topic, DefaultQos, false, message)
	token.Wait()
}
