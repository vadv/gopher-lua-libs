package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	lua_http "github.com/vadv/gopher-lua-libs/http/client/interface"

	multipartstreamer "github.com/technoweenie/multipartstreamer"
	lua "github.com/yuin/gopher-lua"
)

type luaBot struct {
	client lua_http.LuaHTTPClient
	token  string
	offset int
}

func checkBot(L *lua.LState, n int) *luaBot {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaBot); ok {
		return v
	}
	L.ArgError(n, "telegram_bot_ud excepted")
	return nil
}

func (bot *luaBot) decodeAPIResponse(responseBody io.Reader, resp *APIResponse) (_ []byte, err error) {
	dec := json.NewDecoder(responseBody)
	err = dec.Decode(resp)
	return
}

func (bot *luaBot) makeRequest(endpoint string, params url.Values) (APIResponse, error) {
	method := fmt.Sprintf(APIEndpoint, bot.token, endpoint)
	resp, err := bot.client.PostFormRequest(method, params)
	if err != nil {
		return APIResponse{}, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	_, err = bot.decodeAPIResponse(resp.Body, &apiResp)
	if err != nil {
		return apiResp, err
	}

	if !apiResp.Ok {
		parameters := ResponseParameters{}
		if apiResp.Parameters != nil {
			parameters = *apiResp.Parameters
		}
		return apiResp, Error{apiResp.Description, parameters}
	}

	return apiResp, nil
}

func (bot *luaBot) getUpdates(c *UpdateConfig) ([]Update, error) {
	v := url.Values{}
	if bot.offset != 0 {
		v.Add(`offset`, strconv.Itoa(bot.offset))
	}
	if c != nil {
		if c.Offset > 0 {
			v.Set(`offset`, strconv.Itoa(c.Offset))
		}
		if c.Limit > 0 {
			v.Add(`limit`, strconv.Itoa(c.Limit))
		}
		if c.Timeout > 0 {
			v.Add(`limit`, strconv.Itoa(c.Timeout))
		}
	}
	resp, err := bot.makeRequest(`getUpdates`, v)
	if err != nil {
		return []Update{}, err
	}

	var updates []Update
	json.Unmarshal(resp.Result, &updates)

	for _, u := range updates {
		if u.UpdateID >= bot.offset {
			bot.offset = u.UpdateID + 1
		}
	}

	return updates, nil
}

// Send will send a Chattable item to Telegram.
//
// It requires the Chattable to send.
func (bot *luaBot) send(c Chattable) (Message, error) {
	switch c.(type) {
	case Fileable:
		return bot.sendFile(c.(Fileable))
	default:
		return bot.sendChattable(c)
	}
}

// sendFile determines if the file is using an existing file or uploading
// a new file, then sends it as needed.
func (bot *luaBot) sendFile(config Fileable) (Message, error) {
	if config.useExistingFile() {
		return bot.sendExisting(config.method(), config)
	}

	return bot.uploadAndSend(config.method(), config)
}

// sendChattable sends a Chattable.
func (bot *luaBot) sendChattable(config Chattable) (Message, error) {
	v, err := config.values()
	if err != nil {
		return Message{}, err
	}
	message, err := bot.makeMessageRequest(config.method(), v)
	if err != nil {
		return Message{}, err
	}
	return message, nil
}

// makeMessageRequest makes a request to a method that returns a Message.
func (bot *luaBot) makeMessageRequest(endpoint string, params url.Values) (Message, error) {
	resp, err := bot.makeRequest(endpoint, params)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	return message, nil
}

// sendExisting will send a Message with an existing file to Telegram.
func (bot *luaBot) sendExisting(method string, config Fileable) (Message, error) {
	v, err := config.values()

	if err != nil {
		return Message{}, err
	}

	message, err := bot.makeMessageRequest(method, v)
	if err != nil {
		return Message{}, err
	}

	return message, nil
}

// uploadAndSend will send a Message with a new file to Telegram.
func (bot *luaBot) uploadAndSend(method string, config Fileable) (Message, error) {
	params, err := config.params()
	if err != nil {
		return Message{}, err
	}

	file := config.getFile()

	resp, err := bot.uploadFile(method, params, config.name(), file)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	return message, nil
}

// UploadFile makes a request to the API with a file.
//
// Requires the parameter to hold the file not be in the params.
// File should be a string to a file path, a FileBytes struct,
// a FileReader struct, or a url.URL.
//
// Note that if your FileReader has a size set to -1, it will read
// the file into memory to calculate a size.
func (bot *luaBot) uploadFile(endpoint string, params map[string]string, fieldname string, file interface{}) (APIResponse, error) {
	ms := multipartstreamer.New()

	switch f := file.(type) {
	case string:
		ms.WriteFields(params)

		fileHandle, err := os.Open(f)
		if err != nil {
			return APIResponse{}, err
		}
		defer fileHandle.Close()

		fi, err := os.Stat(f)
		if err != nil {
			return APIResponse{}, err
		}

		ms.WriteReader(fieldname, fileHandle.Name(), fi.Size(), fileHandle)
	case FileBytes:
		ms.WriteFields(params)

		buf := bytes.NewBuffer(f.Bytes)
		ms.WriteReader(fieldname, f.Name, int64(len(f.Bytes)), buf)
	case FileReader:
		ms.WriteFields(params)

		if f.Size != -1 {
			ms.WriteReader(fieldname, f.Name, f.Size, f.Reader)

			break
		}

		data, err := ioutil.ReadAll(f.Reader)
		if err != nil {
			return APIResponse{}, err
		}

		buf := bytes.NewBuffer(data)

		ms.WriteReader(fieldname, f.Name, int64(len(data)), buf)
	case url.URL:
		params[fieldname] = f.String()

		ms.WriteFields(params)
	default:
		return APIResponse{}, errors.New(ErrBadFileType)
	}

	method := fmt.Sprintf(APIEndpoint, bot.token, endpoint)

	req, err := http.NewRequest("POST", method, nil)
	if err != nil {
		return APIResponse{}, err
	}

	ms.SetupRequest(req)

	res, err := bot.client.DoRequest(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return APIResponse{}, err
	}

	var apiResp APIResponse

	err = json.Unmarshal(bytes, &apiResp)
	if err != nil {
		return APIResponse{}, err
	}

	if !apiResp.Ok {
		return APIResponse{}, errors.New(apiResp.Description)
	}

	return apiResp, nil
}
