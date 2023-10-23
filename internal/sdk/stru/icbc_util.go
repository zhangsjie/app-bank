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
	"fmt"
	"gitlab.yoyiit.com/youyi/app-bank/internal/enum"
	"gitlab.yoyiit.com/youyi/go-core/config"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"go.uber.org/zap"
	"net/url"
	"sort"
	"strings"
	"time"
)

func NewIcbcGlobalRequest() *IcbcGlobalRequest {
	msgId, _ := util.SonyflakeID()
	request := IcbcGlobalRequest{
		AppID:     config.GetString(enum.IcbcAppId, ""),
		MsgID:     msgId,
		SignType:  "RSA2",
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	return &request
}

func GenerateParams(request IcbcGlobalRequest) (*IcbcGlobalRequest, error) {
	// Step 1: 参数排序
	params := make(map[string]string)
	params["app_id"] = request.AppID
	params["msg_id"] = request.MsgID
	params["sign_type"] = request.SignType
	params["timestamp"] = request.Timestamp

	bizContent, err := MarshalBizContent(request.BizContent)
	if err != nil {
		fmt.Println("Error marshaling biz_content:", err)
		return nil, err
	}
	params["biz_content"] = bizContent

	sortedParamKeys := sortParams(params)
	// Step 2: 构造签名文本
	signaturePlain := constructSignaturePlain(sortedParamKeys, params)

	// Step 3: Sign the signature plain text with RSA private key
	privateKey := []byte("your_private_key")
	signature, err := signWithRSA(signaturePlain, privateKey)
	if err != nil {
		fmt.Println("Error signing with RSA:", err)
		return nil, err
	}

	request.Sign = signature
	zap.L().Info(fmt.Sprintf("ICBCGenerateParams==%+v\n", request))
	return &request, nil
}

func sortParams(params map[string]string) []string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func constructSignaturePlain(keys []string, params map[string]string) string {
	var signatureStrings []string
	for _, key := range keys {
		value := params[key]
		encodedValue := url.QueryEscape(value)
		signatureStrings = append(signatureStrings, fmt.Sprintf("%s=%s", key, encodedValue))
	}
	return strings.Join(signatureStrings, "&")
}

func signWithRSA(plainText string, privateKeyBytes []byte) (string, error) {
	hashed := sha256.Sum256([]byte(plainText))

	privateKey, err := parsePrivateKey(privateKeyBytes)
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	return encodedSignature, nil
}

func parsePrivateKey(privateKeyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not an RSA private key")
	}

	return rsaPrivateKey, nil
}

func MarshalBizContent(bizContent interface{}) (string, error) {
	jsonData, err := json.Marshal(bizContent)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
