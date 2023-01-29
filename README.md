# System Monitor MQTT (SMM)


SMM is a daemon that periodically reads PC sensors values and publishes them on an MQTT broker. Its main feature is the modularity, obtained with the dynamic loading of sensors: a sensor is simply a Lua script that parses sensors' data.

## Building from source

```console
git clone https://github.com/Pirionfr/system-monitor-MQTT
cd system-monitor-MQTT
make build
```

The executable will be placed in `build/smm`.

To install smm, the systemd service, the default sensors and the default config.yml:

```console
sudo make install
```

## Configuration

The configuration is stored in `/etc/smm/config.yml` by default, but it is possible to use another path passing it as the first argument of the executable.

A default `config.yml` is provided in `default/config.yml`.

The only required setting is `mqtt-broker`, that must be the full address of the MQTT broker.

Other options are described in the default `config.yml`.

## Building a sensor

It is possible to add new sensors simply writing their Lua scripts. For this kind of sensors it is needed to provide the absolute path in the `sensors` field of the config.

A sensor script have a prebuilt `sensor` table that must be filled with data that will be published on MQTT. In particular:

```lua
sensor.Value = "..." --The sensor actual value that will directly published
sensor.Name = "Sensor name" --The sensor extended name
sensor.ID = "sensor_id" --The sensor id, cannot contain spaces (default: file name)
sensor.Unit = "..." --The unit of the sensor (like %, °C, ...)
sensor.Class = "..." --The sensor class, important if you are using an MQTT subscriber
```

For this sensor the value will be published on the topic:

`prefix/client_id/sensor.class/sensor.id`

And the debug message on:

`prefix/client_id/sensor.class/sensor.id-debug`

`prefix` is read from the config, the default value is `system-monitor-MQTT`

`client_id` is read from the config, the default value is the current system hostname

If the `ha` flag in the config is set to true, a JSON Home Assistant config for this sensor will be periodically published on the topic:

`homeassistant/sensor.class/client_id_sensor.id/config`

The JSON will be:

```json
{
    "name": "sensor.name",
    "unit_of_measurement": "sensor.unit",
    "device":
    {
        "name": "client_id",
        "model": "client_id",
        "identifiers": "client_id"
    },
    "expire_after": 60,
    "state_topic": "prefix/client_id/sensor.class/sensor.id",
    "unique_id": "client_id_sensor.id"
}
```

The `expire_after` value makes sensors unavailable if the system is unreachable.

For `sensor.class` possible values refere to the Home Assistant [documentation](https://www.home-assistant.io/integrations/sensor/).

## State

When the daemon goes online, it publishes on

`prefix/client_id/state`

The value `Online`, with the relative JSON configuration for Home Assistant, if needed.

When `SIGINT` or `SIGKILL` is received the `Offline` state is published.

