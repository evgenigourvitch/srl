package main

import (
	"errors"
	"hash/crc64"
)

const (
	cMaxBodySize                = int64(0x2000)
	cContentTypeApplicationJson = "application/json"
)

var (
	cBlockedResult         = []byte(`{"block":true}`)
	gCrc64Table            = crc64.MakeTable(crc64.ECMA)
	cMethodNotAllowedErr   = errors.New("non-post method")
	cWrongContentTypeErr   = errors.New("wrong content type")
	cWrongPayloadObjectErr = errors.New("wrong payload")
)
