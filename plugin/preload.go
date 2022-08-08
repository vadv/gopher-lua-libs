package plugin

import (
	"github.com/vadv/gopher-lua-libs/argparse"
	"github.com/vadv/gopher-lua-libs/aws/cloudwatch"
	"github.com/vadv/gopher-lua-libs/base64"
	"github.com/vadv/gopher-lua-libs/cert_util"
	"github.com/vadv/gopher-lua-libs/chef"
	"github.com/vadv/gopher-lua-libs/cmd"
	"github.com/vadv/gopher-lua-libs/crypto"
	"github.com/vadv/gopher-lua-libs/db"
	"github.com/vadv/gopher-lua-libs/filepath"
	"github.com/vadv/gopher-lua-libs/goos"
	"github.com/vadv/gopher-lua-libs/http"
	"github.com/vadv/gopher-lua-libs/humanize"
	"github.com/vadv/gopher-lua-libs/inspect"
	"github.com/vadv/gopher-lua-libs/ioutil"
	"github.com/vadv/gopher-lua-libs/json"
	"github.com/vadv/gopher-lua-libs/log"
	"github.com/vadv/gopher-lua-libs/pb"
	"github.com/vadv/gopher-lua-libs/pprof"
	prometheus "github.com/vadv/gopher-lua-libs/prometheus/client"
	"github.com/vadv/gopher-lua-libs/regexp"
	"github.com/vadv/gopher-lua-libs/runtime"
	"github.com/vadv/gopher-lua-libs/shellescape"
	"github.com/vadv/gopher-lua-libs/stats"
	"github.com/vadv/gopher-lua-libs/storage"
	"github.com/vadv/gopher-lua-libs/strings"
	"github.com/vadv/gopher-lua-libs/tac"
	"github.com/vadv/gopher-lua-libs/tcp"
	"github.com/vadv/gopher-lua-libs/telegram"
	"github.com/vadv/gopher-lua-libs/template"
	"github.com/vadv/gopher-lua-libs/time"
	"github.com/vadv/gopher-lua-libs/xmlpath"
	"github.com/vadv/gopher-lua-libs/yaml"
	"github.com/vadv/gopher-lua-libs/zabbix"
	lua "github.com/yuin/gopher-lua"
)

// PreloadAll preload all gopher lua packages - note it's needed here to prevent circular deps between plugin and libs
func PreloadAll(L *lua.LState) {
	Preload(L)

	argparse.Preload(L)
	base64.Preload(L)
	cert_util.Preload(L)
	chef.Preload(L)
	cloudwatch.Preload(L)
	cmd.Preload(L)
	crypto.Preload(L)
	db.Preload(L)
	filepath.Preload(L)
	goos.Preload(L)
	http.Preload(L)
	humanize.Preload(L)
	inspect.Preload(L)
	ioutil.Preload(L)
	json.Preload(L)
	log.Preload(L)
	pb.Preload(L)
	pprof.Preload(L)
	prometheus.Preload(L)
	regexp.Preload(L)
	runtime.Preload(L)
	shellescape.Preload(L)
	stats.Preload(L)
	storage.Preload(L)
	strings.Preload(L)
	tac.Preload(L)
	tcp.Preload(L)
	telegram.Preload(L)
	template.Preload(L)
	time.Preload(L)
	xmlpath.Preload(L)
	yaml.Preload(L)
	zabbix.Preload(L)
}
