local xml = require("xmlpath")
local inspect = require("inspect")

local data = [[
<channels>
    <channel id="1" xz1="600" />
    <channel id="2"           />
    <channel id="x" xz2="600" />
</channels>
]]
local data_path = "//channel/@id"

local node, err = xml.load(data)
if err then error(err) end

local path, err = xml.compile(data_path)
if err then error(err) end

local iter = path:iter(node)
print(inspect(iter))

for k, v in pairs(iter) do
    print(v:string())
end
