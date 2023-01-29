local content

local run = io.popen
local results = run("cat /sys/devices/system/cpu/possible  | sed 's/-//g'")
for r in results:lines()
do
	content = r
end
results:close()

core = tonumber(content) + 1
	
sensor.Value = string.format("%d", core)
sensor.Name = "Cpu cores"
sensor.Unit = ""
sensor.Class = "None"
