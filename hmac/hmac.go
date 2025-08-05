package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Hmac struct {
	secret string
}

func New(secret string) Hmac {
	return Hmac{
		secret: secret,
	}
}

func (m *Hmac) Sign(str string) string {
	h := hmac.New(sha256.New, []byte(m.secret))

	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))
}

func (m *Hmac) Verify(str string, signature string) bool {
	expected := m.Sign(str)

	expectedBytes, err1 := hex.DecodeString(expected)
	receivedBytes, err2 := hex.DecodeString(signature)

	if err1 != nil || err2 != nil {
		return false
	}

	return hmac.Equal(expectedBytes, receivedBytes)
}
