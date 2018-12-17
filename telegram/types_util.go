package telegram

import "encoding/json"

func updatesToBytes(updates []Update) []byte {
	result, err := json.Marshal(&updates)
	if err != nil {
		panic(err)
	}
	return result
}

func (m Message) toBytes() []byte {
	result, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return result
}
