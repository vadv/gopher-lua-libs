package libs

import (
	cloudwatch "github.com/vadv/gopher-lua-libs/aws/cloudwatch"
	"github.com/vadv/gopher-lua-libs/base64"
	cert_util "github.com/vadv/gopher-lua-libs/cert_util"
	chef "github.com/vadv/gopher-lua-libs/chef"
	cmd "github.com/vadv/gopher-lua-libs/cmd"
	crypto "github.com/vadv/gopher-lua-libs/crypto"
	db "github.com/vadv/gopher-lua-libs/db"
	filepath "github.com/vadv/gopher-lua-libs/filepath"
	goos "github.com/vadv/gopher-lua-libs/goos"
	http "github.com/vadv/gopher-lua-libs/http"
	humanize "github.com/vadv/gopher-lua-libs/humanize"
	inspect "github.com/vadv/gopher-lua-libs/inspect"
	ioutil "github.com/vadv/gopher-lua-libs/ioutil"
	json "github.com/vadv/gopher-lua-libs/json"
	log "github.com/vadv/gopher-lua-libs/log"
	pb "github.com/vadv/gopher-lua-libs/pb"
	plugin "github.com/vadv/gopher-lua-libs/plugin"
	pprof "github.com/vadv/gopher-lua-libs/pprof"
	prometheus "github.com/vadv/gopher-lua-libs/prometheus/client"
	regexp "github.com/vadv/gopher-lua-libs/regexp"
	runtime "github.com/vadv/gopher-lua-libs/runtime"
	"github.com/vadv/gopher-lua-libs/shellescape"
	"github.com/vadv/gopher-lua-libs/stats"
	storage "github.com/vadv/gopher-lua-libs/storage"
	strings "github.com/vadv/gopher-lua-libs/strings"
	tac "github.com/vadv/gopher-lua-libs/tac"
	tcp "github.com/vadv/gopher-lua-libs/tcp"
	telegram "github.com/vadv/gopher-lua-libs/telegram"
	template "github.com/vadv/gopher-lua-libs/template"
	time "github.com/vadv/gopher-lua-libs/time"
	xmlpath "github.com/vadv/gopher-lua-libs/xmlpath"
	yaml "github.com/vadv/gopher-lua-libs/yaml"
	zabbix "github.com/vadv/gopher-lua-libs/zabbix"

	lua "github.com/yuin/gopher-lua"
)

// Preload preload all gopher lua packages
func Preload(L *lua.LState) {
	base64.Preload(L)
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
	shellescape.Preload(L)
	telegram.Preload(L)
	zabbix.Preload(L)
	pprof.Preload(L)
	prometheus.Preload(L)
	pb.Preload(L)
	crypto.Preload(L)
	goos.Preload(L)
	storage.Preload(L)
	humanize.Preload(L)
	chef.Preload(L)
	template.Preload(L)
	cloudwatch.Preload(L)
	log.Preload(L)
	stats.Preload(L)
}
