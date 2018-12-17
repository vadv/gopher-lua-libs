package zabbix

import "encoding/json"

type rpcRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Auth    string      `json:"auth,omitempty"`
	Id      int         `json:"id"`
}

type rpcResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Error   zbxError    `json:"error"`
	Result  interface{} `json:"result"`
	Id      int         `json:"id"`
}

type zbxError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (z *zbxError) Error() string {
	return z.Data
}

func (r *rpcResponse) resultToBytes() ([]byte, bool, error) {
	if result, ok := r.Result.(string); ok {
		return []byte(result), true, nil
	}
	result, err := json.Marshal(r.Result)
	return result, false, err
}
