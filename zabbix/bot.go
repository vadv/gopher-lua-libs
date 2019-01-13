package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	lua_http "github.com/vadv/gopher-lua-libs/http/client/interface"
	lua "github.com/yuin/gopher-lua"
)

type luaBot struct {
	client   lua_http.LuaHTTPClient
	baseURL  string
	apiURL   string
	loginURL string
	chatURL  string
	user     string
	password string
	debug    bool
	// auth
	id   int
	auth string
}

func checkBot(L *lua.LState, n int) *luaBot {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaBot); ok {
		return v
	}
	L.ArgError(n, "zabbix_bot_ud excepted")
	return nil
}

func (b *luaBot) updateURLs() error {
	if b.baseURL == `` {
		return fmt.Errorf("base_url must be set")
	}

	delimiter := `/`
	if strings.HasSuffix(b.baseURL, `/`) {
		delimiter = ``
	}
	if b.apiURL == `` {
		b.apiURL = b.baseURL + delimiter + `api_jsonrpc.php`
	}
	if b.loginURL == `` {
		b.loginURL = b.baseURL + delimiter + `index.php?login=1`
	}
	if b.chatURL == `` {
		b.chatURL = b.baseURL + delimiter + `chart.php`
	}

	return nil
}

func (b *luaBot) sendRequest(method string, data interface{}) (rpcResponse, error) {
	if b.debug {
		log.Printf("[DEBUG] send request method: `%s` data: %#v\n", method, data)
	}
	id := b.id
	b.id = id + 1
	request := rpcRequest{Jsonrpc: "2.0", Method: method, Params: data, Auth: b.auth, Id: id}
	if method == `APIInfo.version` || method == `user.login` {
		request.Auth = ``
	}
	encoded, err := json.Marshal(request)
	if err != nil {
		return rpcResponse{Error: zbxError{Code: -1, Data: err.Error()}}, err
	}
	if b.debug {
		log.Printf("[DEBUG] send request body: %s\n", encoded)
	}
	httpRequest, err := http.NewRequest(`POST`, b.apiURL, bytes.NewBuffer(encoded))
	if err != nil {
		return rpcResponse{Error: zbxError{Code: -2, Data: err.Error()}}, err
	}
	httpRequest.Header.Set(`Content-Type`, `application/json-rpc`)
	httpResponse, err := b.client.DoRequest(httpRequest)
	if err != nil {
		return rpcResponse{Error: zbxError{Code: -3, Data: err.Error()}}, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode >= 300 {
		return rpcResponse{Error: zbxError{Code: httpResponse.StatusCode, Data: "Bad response code"}}, fmt.Errorf("response code: %d", httpResponse.StatusCode)
	}
	response := rpcResponse{}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, httpResponse.Body)
	if err != nil {
		return rpcResponse{Error: zbxError{Code: -4, Data: err.Error()}}, err
	}
	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		return rpcResponse{Error: zbxError{Code: -4, Data: err.Error()}}, err
	}
	if response.Error.Code != 0 || response.Error.Data != `` {
		return response, fmt.Errorf("error %d: %s", response.Error.Code, response.Error.Data)
	}
	return response, nil
}

func (b *luaBot) login() error {
	params := make(map[string]string, 0)
	params["user"] = b.user
	params["password"] = b.password
	response, err := b.sendRequest("user.login", params)
	if err != nil {
		return err
	}
	if response.Error.Code != 0 {
		return &response.Error
	}
	b.auth = response.Result.(string)
	return nil
}

func (b *luaBot) logout() error {
	response, err := b.sendRequest("user.logout", make(map[string]string, 0))
	if err != nil {
		return err
	}
	if response.Error.Code != 0 {
		return &response.Error
	}
	return nil
}

func (b *luaBot) apiVersion() (string, error) {
	response, err := b.sendRequest("APIInfo.version", make(map[string]string, 0))
	if err != nil {
		return "", err
	}
	if response.Error.Code != 0 {
		return "", &response.Error
	}
	version := response.Result.(string)
	return version, nil
}

func (b *luaBot) saveGraph(filename string, itemId, period, width, height int64) error {

	dst, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()

	values := url.Values{
		"name":      {b.user},
		"password":  {b.password},
		"autologin": {"1"},
		"enter":     {"Sign in"},
	}
	response, err := b.client.PostFormRequest(b.loginURL, values)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("response status from %s: %d", b.loginURL, response.StatusCode)
	}
	curTime := strconv.FormatInt(time.Now().Unix(), 10)
	chartURL := fmt.Sprintf(`%s?period=%d&itemids%%5B0%%5D=%d&width=%d&height=%d&curtime=%s`, b.chatURL, period, itemId, width, height, curTime)
	request, err := http.NewRequest(`GET`, chartURL, nil)
	if err != nil {
		return err
	}
	response, err = b.client.DoRequest(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return fmt.Errorf(`chart url %s response code: %d`, chartURL, response.StatusCode)
	}
	_, err = io.Copy(dst, response.Body)
	return err
}
