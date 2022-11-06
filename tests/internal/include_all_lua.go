package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	packageName = "tests"
	fileName    = "./internal/suite.lua"
	constName   = "lua_suite"
	templateGo  = `// suite.lua for gopher-lua https://pkg.go.dev/github.com/stretchr/testify/suite
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
