package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

func Generate2FA(username string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "QuantaVerse", // 你的应用名称
		AccountName: username,      // 用户名
	})
	if err != nil {
		return "", "", fmt.Errorf("生成密钥失败: %w", err)
	}

	url := key.URL() // 生成 OTP URL

	pngBytes, err := qrcode.Encode(url, qrcode.Medium, 256) //生成二维码的bytes
	if err != nil {
		return "", "", fmt.Errorf("qrcode encode error: %w", err)
	}
	pngBase64 := base64.StdEncoding.EncodeToString(pngBytes)
	pngDataURL := "data:image/png;base64," + pngBase64
	return key.Secret(), pngDataURL, nil
}
