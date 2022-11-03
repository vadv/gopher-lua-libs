# gopher-lua-libs
[![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs)


Package contains is a libs for [gopher-lua](https://github.com/yuin/gopher-lua).

## License

Development version, available on github, released under BSD 3-clause.

## Installation

```
go get github.com/vadv/gopher-lua-libs
```

## Index

* [argparse](/argparse) argparse CLI parsing <https://github.com/luarocks/argparse>
* [base64](/base64) [encoding/base64](https://pkg.go.dev/encoding/base64) api
* [cloudwatch](/aws/cloudwatch) aws cloudwatch log access
* [cert_util](/cert_util) monitoring ssl certs
* [chef](/chef) chef client api
* [cmd](/cmd) cmd port
* [crypto](/crypto) calculate md5, sha256 hash for string
* [db](/db) access to databases
* [filepath](/filepath) path.filepath port
* [goos](/goos) os port
* [http](/http) http.client && http.server
* [humanize](/humanize) humanize [github.com/dustin/go-humanize](https://github.com/dustin/go-humanize) port
* [inspect](/inspect) pretty print [github.com/kikito/inspect.lua](https://github.com/kikito/inspect.lua)
* [ioutil](/ioutil) io/ioutil port
* [json](/json) json implementation
* [log](/log) log port
* [plugin](/plugin) run lua code in lua code
* [pprof](/pprof) pprof http-server for golang from lua
* [prometheus](/prometheus/client) prometheus exporter
* [regexp](/regexp) regexp port
* [runtime](/runtime) runtime port
* [pb](/pb) [https://github.com/cheggaaa/pb](https://github.com/cheggaaa/pb) port (v3)
* [shellescape](/shellescape) shellescape <https://github.com/alessio/shellescape> port
* [stats](/stats) stats [https://github.com/montanaflynn/stats](https://github.com/montanaflynn/stats) port
* [storage](/storage) package for store persist data and share values between lua states
* [strings](/strings) strings port (utf supported)
* [tac](/tac) tac line-by-line scanner (from end of file to up)
* [tcp](/tcp) raw tcp client lib
* [telegram](/telegram) telegram bot
* [template](/template) template engines
* [time](/time) time port
* [xmlpath](/xmlpath) [gopkg.in/xmlpath.v2](https://gopkg.in/xmlpath.v2) port
* [yaml](/yaml) [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2) port
* [zabbix](/zabbix) zabbix bot


## Usage

For the quick overview you can use standalone interpreter with listed libs. Examples and documentation for modules can be found in their directories.
```
go install github.com/vadv/gopher-lua-libs/cmd/glua-libs@latest

glua-libs example.lua
```

This example shows basic usage of this libs in your code

```golang
package main


import (
        "log"
        "flag"
        "os"

        lua "github.com/yuin/gopher-lua"
        libs "github.com/vadv/gopher-lua-libs"

)
var (
        exec = flag.String("execute", "", "execute lua script")
)


func main() {
        flag.Parse()
        state := lua.NewState()
        defer state.Close()
        libs.Preload(state)
        if *exec != `` {
                if err := state.DoFile(*exec); err != nil {
                        log.Printf("[ERROR] Error executing file: ", err)
                }
        } else {
                log.Println("Target file was not given!")
                os.Exit(0)
        }
}


```
