package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var conf *viper.Viper

const (
	UpdatePeriod = "update-period"
	MQTTBroker   = "mqtt-broker"
	MQTTUser     = "mqtt-user"
	MQTTPassword = "mqtt-user"
	Sensors      = "sensors"
	HA           = "ha"
	Prefix       = "prefix"
	ClientID     = "client-id"
	ConfigPath   = "conf-path"
)

func init() {
	conf = viper.New()

	conf.SetDefault(UpdatePeriod, 10)
	conf.SetDefault(MQTTBroker, "localhost:1883")
	conf.SetDefault(Sensors, []string{})
	conf.SetDefault(HA, false)
	conf.SetDefault(Prefix, "system-monitor-MQTT")
	hostname, _ := os.Hostname()
	conf.SetDefault(ClientID, hostname)
	conf.SetDefault(ConfigPath, "/etc/smm/sensors")

	conf.SetConfigName("config")    // name of config file (without extension)
	conf.SetConfigType("yml")       // REQUIRED if the config file does not have the extension in the name
	conf.AddConfigPath("/etc/smm/") // path to look for the config file in
	conf.AddConfigPath("default/")  // optionally look for config in the working directory
	err := conf.ReadInConfig()      // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		log.Fatalf("fatal error config file: %v", err)
	}

	conf.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	conf.WatchConfig()
}

func main() {
	d := NewDaemon(conf)
	d.run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		for {
			<-sigs
			d.NotifyState(false)
			done <- true
		}
	}()

	<-done

}
