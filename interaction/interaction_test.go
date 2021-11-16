package interaction

import (
	"net/http"
	"testing"
)

func TestSignature(t *testing.T) {
	secret := "abcdefg"
	header := http.Header{}
	header.Set(HeaderTimestamp, "1234567890")
	httpBody := "text body"
	sig, err := GenSignature(secret, header, []byte(httpBody))
	if err != nil {
		t.Error(err)
	}
	t.Log(sig)

	header.Set(HeaderSig, sig)
	flag, err := VerifySignature(secret, header, []byte(httpBody))
	if err != nil {
		t.Error(err)
	}
	if !flag {
		t.Error("verify failed, but want ok")
	} else {
		t.Log("verify ok")
	}
}
