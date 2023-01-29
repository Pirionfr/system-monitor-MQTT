local file = io.open("/proc/uptime", "r")
local line = string.match(file:read(),"%d+")
file:close()


local up_d = tonumber(line)/ (3600 * 24)

sensor.Value = string.format("%.1f", up_d)
sensor.Name = "uptime"
sensor.Unit = "d"
sensor.Class = "duration"

