local content

local results = io.popen("df -h / | tail -n1  | awk '{ print $5 }' | sed 's/%//g'")
for r in results:lines()
do
	content = tonumber(r)
end
results:close()

sensor.Value =  string.format("%.1f", content)
sensor.Name = "Root FS usage"
sensor.Unit = "%"
sensor.Class = "None"