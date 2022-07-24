local tac = require("tac")

function Test_tac(t)
    local scanner, err = tac.open("./test/test.txt")
    assert(not err, err)

    local count_line = 3
    while true do
        local line = scanner:line()
        if line == nil then
            break
        end
        assert(tostring(count_line) == line, "must be: " .. tostring(count_line))
        count_line = count_line - 1
    end
    scanner:close()
end