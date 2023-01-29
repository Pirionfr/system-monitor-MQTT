local content
local unit
local run = io.popen
local results = run("df -h  / | tail -n1  | awk '{ print $4 }' | sed 's/G//g'") 
for r in results:lines()
do
	content = r:sub(1, -2)
	unit = r:sub(-1, -1)
end
results:close()

sensor.Value = content
sensor.Name = "root FS free"
sensor.Unit = unit .. "iB"
sensor.Class = "data_size"


