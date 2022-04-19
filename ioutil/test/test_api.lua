local ioutil = require("ioutil")

local test_data = "test\n"

function TestWriteRead(t)
    t:Run("write_file", function(t)
        local err = ioutil.write_file("./test/file.test", test_data)
        if err then error(err) end
    end)

    t:Run("read_file", function(t)
        local data, err = ioutil.read_file("./test/file.test")
        if err then error(err) end
        if not(data == test_data) then error("ioutil.read_file()") end
        _, err = ioutil.read_file("./test/unknown.test")
        if err == nil then error("ioutil.read_file() must be error") end
        print("done: ioutil.read_file()")
    end)
end

function TestCopy(t)
    local input_fh, err = io.open("./test/file.test", "r")
    assert(not err, err)
    local output_fh, err = io.open("./test/file2.data", "w")
    assert(not err, err)
    ioutil.copy(output_fh, input_fh)
    input_fh:close()
    output_fh:close()

    local orig_data, err = ioutil.read_file("./test/file.test")
    assert(not err, err)
    local copied_data, err = ioutil.read_file("./test/file2.data")
    assert(not err, err)
    assert(orig_data == copied_data, string.format("'%s' ~= '%s'", orig_data, copied_data))
end
