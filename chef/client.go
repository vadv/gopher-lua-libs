package chef

import (
	"crypto/rsa"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	lua_http "github.com/vadv/gopher-lua-libs/http/client/interface"

	lua "github.com/yuin/gopher-lua"
)

// From https://github.com/go-chef/chef/

const (
	ChefVersion = "11.12.0" // default client version
)

type luaChefClient struct {
	lua_http.LuaHTTPClient
	key  *rsa.PrivateKey
	name string
	url  *url.URL
}

func checkChefClient(L *lua.LState, n int) *luaChefClient {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaChefClient); ok {
		return v
	}
	L.ArgError(n, "chef_client_ud expected")
	return nil
}

func (c *luaChefClient) request(method, path string, body io.Reader) ([]byte, error) {
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	res, err := c.LuaHTTPClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	err = checkResponse(res)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

func (c *luaChefClient) newRequest(method string, requestURL string, body io.Reader) (*http.Request, error) {
	relativeURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, err
	}
	u := c.url.ResolveReference(relativeURL)
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	values := req.URL.Query()
	req.URL.RawQuery = values.Encode()
	myBody := &luaBody{body}
	if body != nil {
		// Detect Content-type
		req.Header.Set("Content-Type", myBody.contentType())
	}
	// Calculate the body hash
	req.Header.Set("X-Ops-Content-Hash", myBody.hash())

	// don't have to check this works, signRequest only emits error when signing hash is not valid, and we baked that in
	c.signRequest(req)
	return req, nil
}

func (c *luaChefClient) signRequest(request *http.Request) error {
	// sanitize the path for the chef-server
	// chef-server doesn't support '//' in the Hash Path.
	var endpoint string
	if request.URL.Path != "" {
		endpoint = path.Clean(request.URL.Path)
		request.URL.Path = endpoint
	} else {
		endpoint = request.URL.Path
	}

	vals := map[string]string{
		"Method":             request.Method,
		"Hashed Path":        hashStr(endpoint),
		"Accept":             "application/json",
		"X-Chef-Version":     ChefVersion,
		"X-Ops-Timestamp":    time.Now().UTC().Format(time.RFC3339),
		"X-Ops-UserId":       c.name,
		"X-Ops-Sign":         "algorithm=sha1;version=1.0",
		"X-Ops-Content-Hash": request.Header.Get("X-Ops-Content-Hash"),
	}

	for _, key := range []string{"Method", "Accept", "X-Chef-Version", "X-Ops-Timestamp", "X-Ops-UserId", "X-Ops-Sign"} {
		request.Header.Set(key, vals[key])
	}

	// To validate the signature it seems to be very particular
	var content string
	for _, key := range []string{"Method", "Hashed Path", "X-Ops-Content-Hash", "X-Ops-Timestamp", "X-Ops-UserId"} {
		content += fmt.Sprintf("%s:%s\n", key, vals[key])
	}
	content = strings.TrimSuffix(content, "\n")
	// generate signed string of headers
	// Since we've gone through additional validation steps above,
	// we shouldn't get an error at this point
	signature, err := generateSignature(c.key, content)
	if err != nil {
		return err
	}

	// TODO: THIS IS CHEF PROTOCOL SPECIFIC
	// Signature is made up of n 60 length chunks
	base64sig := base64BlockEncode(signature, 60)

	// roll over the auth slice and add the apropriate header
	for index, value := range base64sig {
		request.Header.Set(fmt.Sprintf("X-Ops-Authorization-%d", index+1), string(value))
	}

	return nil

}
