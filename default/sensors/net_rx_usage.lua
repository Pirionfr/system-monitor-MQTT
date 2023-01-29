local content

local run = io.popen
local results = run("cat /proc/net/dev | grep 'enp6s0:' | awk '{ print $2 }'")
for r in results:lines()
do
	content = r
end
results:close()

content = content/1073741824

sensor.Value = string.format("%.3f", content)
sensor.Name = "Network RX usage"
sensor.Unit = "GiB"
sensor.Class = "None"
