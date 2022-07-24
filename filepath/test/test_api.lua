local filepath = require("filepath")

function Test_join_and_separator(t)
    local path = "1"
    local need_path = path .. filepath.separator() .. "2" .. filepath.separator() .. "3"
    path = filepath.join(path, "2", "3")
    assert(path == need_path, string.format("expected %s; got %s", need_path, path))
end

function Test_glob(t)
    local results = filepath.glob("test" .. filepath.separator() .. "*")
    assert(#results == 1, string.format("expected one glob result; got %d", #results))
    for k, v in pairs(results) do
        if k == 1 then
            assert(v == "test" .. filepath.separator() .. "test_api.lua", v)
        end
    end
end