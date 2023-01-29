local open = io.open

local file = open("/sys/class/hwmon/hwmon2/temp1_input", "r") 
local content = file:read() --first line

file:close()

local temp = tonumber(content)/1000
	
sensor.Value = tostring(temp)
sensor.Name = "CPU Temperature"
sensor.Unit = "°C"
sensor.Class = "temperature"
