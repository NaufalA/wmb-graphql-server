package util

import "encoding/base64"

type PaginationUtil struct {}

func (p *PaginationUtil) EncodeCursor(param string) string {
	return base64.StdEncoding.EncodeToString([]byte(param))
}

func (p *PaginationUtil) DecodeCursor(cursor string) (string, error) {
	param, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return "", err
	}
	return string(param), nil
}
