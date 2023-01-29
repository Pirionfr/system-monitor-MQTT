local open = io.open
local file = open("/proc/meminfo", "r")
string.match(file:read(), "%d+") --first line
local free = string.match(file:read(), "%d+") --second line
file:close()

local memory = free/1000000

sensor.Value = string.format("%.1f", memory)
sensor.Name = "Available memory"
sensor.Unit = "GB"
sensor.Class = "None"