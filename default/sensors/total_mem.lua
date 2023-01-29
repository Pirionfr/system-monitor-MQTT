local open = io.open
local file = open("/proc/meminfo", "r")
local total = string.match(file:read(), "%d+") --first line
file:close()

local memory = total/1000000

sensor.Value = string.format("%.1f", memory)
sensor.Name = "Total memory"
sensor.Unit = "GB"
sensor.Class = "None"