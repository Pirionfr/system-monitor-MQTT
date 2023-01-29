local open = io.open
local file = open("/proc/meminfo", "r")
local total = string.match(file:read(), "%d+") --first line
string.match(file:read(), "%d+") --second line
local available = string.match(file:read(), "%d+") --third linbe
file:close()

local usage = 100 - (100*available/total)

sensor.Value = string.format("%.1f", usage)
sensor.Name = "Memory usage"
sensor.Unit = "%"
sensor.Class = "None"
