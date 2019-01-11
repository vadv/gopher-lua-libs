local template = require("template")

local mustache, err = template.choose("mustache")
if err then error(err) end

local values = {name="world"}
local result, err = mustache:render("Hello {{name}}!", values)
if err then error(err) end
if not(result == "Hello world!") then error(result) end

local values = {data = {"one", "two"}}
local result, err = mustache:render("{{#data}} {{.}} {{/data}}", values)
if err then error(err) end
if not(result == " one  two ") then error(result) end
