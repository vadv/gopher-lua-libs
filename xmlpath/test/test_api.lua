local xml = require("xmlpath")
local inspect = require("inspect")

function Test_xmlpath(t)
    local data = [[
<channels>
    <channel id="1" xz1="600" />
    <channel id="2"           />
    <channel id="x" xz2="600" />
</channels>
]]
    local data_path = "//channel/@id"

    local node, err = xml.load(data)
    assert(not err, err)

    local path, err = xml.compile(data_path)
    assert(not err, err)

    local iter = path:iter(node)
    t:Log(inspect(iter))

    for k, v in pairs(iter) do
        t:Logf("%s => %v", k, v:string())
    end
end
