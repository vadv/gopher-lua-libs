local regexp = require("regexp")

function Test_compiled_regex_match(t)
    local reg, err = regexp.compile("(gopher){2}")
    if err then
        error(err)
    end
    assert(not reg:match("gopher"), "must not be matched")
    assert(reg:match("gophergopher"), "must be matched")
end

function Test_find_all_string_submatch(t)
    t:Run("1", function(t)
        local reg, err = regexp.compile("string: (.*)$")
        assert(not err, err)
        local result = reg:find_all_string_submatch("my string: 'hello world'")
        assert(result[1][2] == "'hello world'", "not found: " .. tostring(result[1][2]))
    end)

    t:Run("2", function(t)
        local reg, err = regexp.compile("string: '(.*)\\s+(.*)'$")
        assert(not err, err)
        local result = reg:find_all_string_submatch("my string: 'hello world'")
        assert(result[1][2] == "hello", "not found: " .. tostring(result[1][2]))
        assert(result[1][3] == "world", "not found: " .. tostring(result[1][3]))
    end)
end

function Test_match(t)
    local found, err = regexp.match("(gopher){2}", "gophergopher")
    assert(not err, err)
    assert(found, "must be matched")
end

function Test_find_all_string_submatch(t)
    local result, err = regexp.find_all_string_submatch("string: '(.*)\\s+(.*)'$", "my string: 'hello world'")
    assert(not err, err)
    assert(result[1][2] == "hello", "not found: " .. tostring(result[1][2]))
    assert(result[1][3] == "world", "not found: " .. tostring(result[1][3]))
end