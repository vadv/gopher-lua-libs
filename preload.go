package libs

import (
	cert_util "github.com/vadv/gopher-lua-libs/cert_util"
	cmd "github.com/vadv/gopher-lua-libs/cmd"
	db "github.com/vadv/gopher-lua-libs/db"
	filepath "github.com/vadv/gopher-lua-libs/filepath"
	http "github.com/vadv/gopher-lua-libs/http"
	inspect "github.com/vadv/gopher-lua-libs/inspect"
	ioutil "github.com/vadv/gopher-lua-libs/ioutil"
	json "github.com/vadv/gopher-lua-libs/json"
	plugin "github.com/vadv/gopher-lua-libs/plugin"
	regexp "github.com/vadv/gopher-lua-libs/regexp"
	runtime "github.com/vadv/gopher-lua-libs/runtime"
	strings "github.com/vadv/gopher-lua-libs/strings"
	tac "github.com/vadv/gopher-lua-libs/tac"
	tcp "github.com/vadv/gopher-lua-libs/tcp"
	telegram "github.com/vadv/gopher-lua-libs/telegram"
	time "github.com/vadv/gopher-lua-libs/time"
	xmlpath "github.com/vadv/gopher-lua-libs/xmlpath"
	yaml "github.com/vadv/gopher-lua-libs/yaml"

	lua "github.com/yuin/gopher-lua"
)

// Preload(): preload all gopher lua packages
func Preload(L *lua.LState) {
	time.Preload(L)
	strings.Preload(L)
	filepath.Preload(L)
	ioutil.Preload(L)
	http.Preload(L)
	regexp.Preload(L)
	tac.Preload(L)
	inspect.Preload(L)
	yaml.Preload(L)
	plugin.Preload(L)
	cmd.Preload(L)
	json.Preload(L)
	tcp.Preload(L)
	xmlpath.Preload(L)
	db.Preload(L)
	cert_util.Preload(L)
	runtime.Preload(L)
	telegram.Preload(L)
}
