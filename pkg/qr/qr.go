package qr

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

func Generate(content string) (string, error) {
	var png []byte
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	base64Img := base64.StdEncoding.EncodeToString(png)
	return "data:image/png;base64," + base64Img, nil
}
