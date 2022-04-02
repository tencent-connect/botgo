// Package signature 用于处理平台和机器人开发者之间的互动请求中的签名验证
package signature

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/tencent-connect/botgo/log"
)

const (
	// HeaderSig 请求签名
	HeaderSig = "X-Signature-Ed25519"
	// HeaderTimestamp 跟请求签名对应的时间戳，用于验证签名
	HeaderTimestamp = "X-Signature-Timestamp"
)

type ed25519Key struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

// Verify 验证签名，需要传入 http 头，httpBody
// 请在方法外部从 http request 上读取了 body 之后再交给签名验证方法进行验证，避免重复读取
func Verify(secret string, header http.Header, httpBody []byte) (bool, error) {
	// 生成密钥
	key, err := genKey(secret)
	if err != nil {
		log.Errorf("genPublicKey error, %v", err)
		return false, err
	}
	sigBuffer, err := decodeSigBuffer(header.Get(HeaderSig))
	if err != nil {
		log.Errorf("decodeSigBuffer error, %v", err)
		return false, err
	}
	content, err := genOriginalContent(header.Get(HeaderTimestamp), httpBody)
	if err != nil {
		log.Errorf("get original content error: %v", err)
	}
	return ed25519.Verify(key.PublicKey, content, sigBuffer), nil
}

// Generate 生成签名，sdk 中的改方法，主要用于与验证签名方法配合进行验证
func Generate(secret string, header http.Header, httpBody []byte) (string, error) {
	key, err := genKey(secret)
	if err != nil {
		log.Errorf("genPrivateKey error, %v", err)
		return "", err
	}
	content, err := genOriginalContent(header.Get(HeaderTimestamp), httpBody)
	if err != nil {
		log.Errorf("get original content error: %v", err)
	}
	return hex.EncodeToString(ed25519.Sign(key.PrivateKey, content)), nil
}

func genOriginalContent(timestamp string, body []byte) ([]byte, error) {
	if timestamp == "" {
		return nil, errors.New("timestamp is nil")
	}
	// 按照 timestamp+Body 顺序组成签名体
	var msg bytes.Buffer
	msg.WriteString(timestamp)
	msg.Write(body)
	return msg.Bytes(), nil
}

// genKey 根据 seed 生成公钥，私钥，私钥用于请求方加密，公钥用于服务方验证
func genKey(secret string) (*ed25519Key, error) {
	seed, err := getSeed(secret)
	if err != nil {
		return nil, err
	}
	publicKey, privateKey, err := ed25519.GenerateKey(strings.NewReader(seed))
	if err != nil {
		return nil, err
	}
	return &ed25519Key{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

func decodeSigBuffer(signature string) ([]byte, error) {
	if signature == "" {
		return nil, errors.New("not found signature")
	}
	sigBuf, err := hex.DecodeString(signature)
	if err != nil {
		return nil, fmt.Errorf("hex decode signature failed: %v", err)
	}
	if len(sigBuf) != ed25519.SignatureSize || sigBuf[63]&224 != 0 {
		return nil, errors.New("signature decode result is not a valid buf")
	}
	return sigBuf, nil
}

// getSeed 使用 secret 生成算法 seed
func getSeed(secret string) (string, error) {
	if secret == "" {
		return "", errors.New("secret invalid")
	}
	seed := secret
	for len(seed) < ed25519.SeedSize {
		seed = strings.Repeat(seed, 2)
	}
	return seed[:ed25519.SeedSize], nil
}
