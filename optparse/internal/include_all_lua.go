package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	packageName = "optparse"
	fileName    = "./internal/optparse.lua"
	constName   = "lua_optparse"
	templateGo  = `// optparse.lua for gopher-lua
package %s

const %s = "%s"
`
)

func main() {

	out, err := os.Create("lua_const.go")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer out.Close()

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	base46Code := base64.StdEncoding.EncodeToString(data)
	content := fmt.Sprintf(templateGo, packageName, constName, base46Code)
	out.WriteString(content)
}
