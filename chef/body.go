package chef

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// luaBody wraps io.Reader and adds methods for calculating hashes and detecting content
type luaBody struct {
	io.Reader
}

// contentType returns the content-type string of luaBody as detected by http.DetectContentType()
func (body *luaBody) contentType() string {
	if json.Unmarshal(body.buffer().Bytes(), &struct{}{}) == nil {
		return "application/json"
	}
	return http.DetectContentType(body.buffer().Bytes())
}

// hash calculates the body content hash
func (body *luaBody) hash() (h string) {
	b := body.buffer()
	// empty buffs should return a empty string
	if b.Len() == 0 {
		h = hashStr("")
	}
	h = hashStr(b.String())
	return
}

// hashStr returns the base64 encoded SHA1 sum of the toHash string
func hashStr(toHash string) string {
	h := sha1.New()
	io.WriteString(h, toHash)
	hashed := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return hashed
}

// buffer creates a  byte.buffer copy from a io.Reader resets read on reader to 0,0
func (body *luaBody) buffer() *bytes.Buffer {
	var b bytes.Buffer
	if body.Reader == nil {
		return &b
	}

	b.ReadFrom(body.Reader)
	_, err := body.Reader.(io.Seeker).Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	return &b
}
