// GENERATED BY textFileToGoConst
// GitHub:     github.com/logrusorgru/textFileToGoConst
// input file: require.lua
// generated:  Fri Nov 11 18:29:59 PST 2022

package tests

const lua_require = `local assertions = require 'assertions'
local orig_require = require
return assertions:new {
    fail_now = true,
    call = orig_require,
}

`