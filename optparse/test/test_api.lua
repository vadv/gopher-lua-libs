-- optparse can only work with the arg global; set one up, them mutate it (not reassign it)
arg = {}
local inspect = require("inspect")
local optparse = require("optparse")

function Test_optparse(t)
    local opt = optparse.OptionParser { usage = "%prog [options] [file...]", version = "test 1.23" }
    opt.add_option { "-i", "--input", dest = "input", action = "store" }
    opt.add_option { "-v", "--verbose", dest = "verbose", action = "store_true" }
    table.remove(arg)
    table.insert(arg, "-v")
    table.insert(arg, "-i")
    table.insert(arg, "myinput")
    table.insert(arg, "foo")
    local options, args = opt.parse_args()
    print(string.format("options=%s", inspect(options)))
    print(string.format("args=%s", inspect(args)))
    assert(options.verbose, tostring(options.verbose))
    assert(options.input == "myinput", tostring(options.input))
    assert(args[1] == "foo", tostring(args[1]))
end