package filepath

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds filepath to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//	local filepath = require("filepath")
func Preload(L *lua.LState) {
	L.PreloadModule("filepath", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"abs":            Abs,
	"basename":       Basename,
	"clean":          Clean,
	"dir":            Dir,
	"eval_symlinks":  EvalSymlinks,
	"from_slash":     FromSlash,
	"ext":            Ext,
	"glob":           Glob,
	"is_abs":         IsAbs,
	"join":           Join,
	"list_separator": ListSeparator,
	"match":          Match,
	"rel":            Rel,
	"separator":      Separator,
	"split":          Split,
	"split_list":     SplitList,
	"to_slash":       ToSlash,
	"volume_name":    VolumeName,
}
