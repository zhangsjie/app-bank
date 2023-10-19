package stru

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

func main() {
	appID := "your_app_id"
	privateKey := []byte("your_private_key")

	// 构造请求参数
	bizContent := AccDetailRequest{
		Fseqno:    "your_fseqno",
		Account:   "your_account",
		Currtype:  1,
		Startdate: "2023-01-01",
		Enddate:   "2023-01-31",
		Serialno:  "",
		Corpno:    "your_corpno",
		Acccompno: "your_acccompno",
		Agreeno:   "your_agreeno",
	}

	request := IcbcGlobalRequest{
		AppID:      appID,
		MsgID:      generateRandomString(16),
		SignType:   "RSA2",
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		BizContent: bizContent,
	}

	// 对请求参数进行签名
	sign, err := signRequest(request, privateKey)
	if err != nil {
		fmt.Println("Failed to sign the request:", err)
		return
	}
	request.Sign = sign

	// 发送请求...
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	randomBytes := make([]byte, length+(length/4))
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err)
	}
	for i, c := range randomBytes {
		result[i] = charset[int(c)%len(charset)]
	}
	return string(result)
}

func signRequest(request IcbcGlobalRequest, privateKey []byte) (string, error) {
	// 将请求参数转换为URL编码的字符串
	rawQuery, err := url.QueryUnescape(string(encodeJson(request.BizContent)))
	if err != nil {
		return "", err
	}

	// 按照参数名排序
	params := make(map[string]string)
	params["app_id"] = request.AppID
	params["msg_id"] = request.MsgID
	params["sign_type"] = request.SignType
	params["timestamp"] = request.Timestamp
	params["biz_content"] = rawQuery

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 拼接参数名和参数值
	var paramString strings.Builder
	for _, key := range keys {
		value := params[key]
		if paramString.Len() > 0 {
			paramString.WriteByte('&')
		}
		paramString.WriteString(key)
		paramString.WriteByte('=')
		paramString.WriteString(value)
	}

	// 使用RSA2算法对参数进行签名
	hash := sha256.New()
	hash.Write([]byte(paramString.String()))
	hashed := hash.Sum(nil)

	privateKeyParsed, err := parsePrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKeyParsed, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	sign := base64.StdEncoding.EncodeToString(signature)
	return sign, nil
}

func encodeJson(data interface{}) []byte {
	encoded, _ := json.Marshal(data)
	return encoded
}

func parsePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing the private key")
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKeyParsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("Invalid private key")
	}

	return privateKeyParsed, nil
}
