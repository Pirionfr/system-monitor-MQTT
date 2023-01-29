local content

local run = io.popen
local results = run("uname -r")
for r in results:lines()
do
	content = r
end
results:close()

sensor.Value = content
sensor.Name = "Linux kernel"
sensor.Unit = ""
sensor.Class = "None"
