package savestusignin

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// 解析Base64编码的公钥
func parsePublicKey(publicKeyBase64 string) (*rsa.PublicKey, error) {
	der, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("Base64解码失败: %v", err)
	}

	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		block, _ := pem.Decode(der)
		if block == nil {
			return nil, fmt.Errorf("解析DER公钥失败: %v", err)
		}
		pub, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("解析PEM公钥失败: %v", err)
		}
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("公钥不是RSA类型")
	}

	return rsaPub, nil
}

// 加密密码
func encryptPassword(password string, publicKey *rsa.PublicKey) (string, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(password)) // 使用 rand.Reader
	if err != nil {
		return "", fmt.Errorf("加密失败: %v", err)
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}
