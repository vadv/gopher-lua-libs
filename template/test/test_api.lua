local template = require("template")

function Test_template(t)
    local mustache, err = template.choose("mustache")
    assert(not err, err)

    tests = {
        {
            name = "hello world",
            values = { name = "world" },
            template = "Hello {{name}}!",
            expected = "Hello world!",
        },
        {
            name = "one two",
            values = { data = { "one", "two" } },
            template = "{{#data}} {{.}} {{/data}}",
            expected = " one  two ",
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            local result, err = mustache:render(tt.template, tt.values)
            assert(not err, err)
            assert(result == tt.expected, string.format([[expected "%s"; got "%s"]], tt.expected, result))
        end)
    end
end
