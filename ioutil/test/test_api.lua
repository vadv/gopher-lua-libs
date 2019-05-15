local ioutil = require("ioutil")

local test_data = "test\n"

local err = ioutil.write_file("./test/file.test", test_data)
if err then error(err) end

local data, err = ioutil.read_file("./test/file.test")
if err then error(err) end
if not(data == test_data) then error("ioutil.read_file()") end
_, err = ioutil.read_file("./test/unknown.test")
if err == nil then error("ioutil.read_file() must be error") end
print("done: ioutil.read_file()")
