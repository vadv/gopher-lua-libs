local tac = require("tac")

local scanner, err = tac.open("./test/test.txt")
if err then error(err) end

local count_line = 3
while true do
    local line = scanner:line()
    if line == nil then break end
    if not(tostring(count_line) == line) then error("must be: "..tostring(count_line)) end
    count_line = count_line - 1
end
scanner:close()

print("done: tac:line()")
