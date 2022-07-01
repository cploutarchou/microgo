package url_signer

import (
	"fmt"
	"strings"
	"time"
)

type Signer struct {
	Secret []byte
}

func (s *Signer) GenerateToken(data string) string {
	var urlToSing string
	crypt := New(s.Secret, Timestamp)
	if strings.Contains(data, "?") {
		urlToSing = fmt.Sprintf("%s&hash=", data)
	} else {
		urlToSing = fmt.Sprintf("%s&hash=", data)
	}
	tokenBytes := crypt.Sign([]byte(urlToSing))
	token := string(tokenBytes)
	return token
}

func (s *Signer) Verify(token string) bool {
	crypt := New(s.Secret, Timestamp)
	_, err := crypt.Unsigned([]byte(token))
	if err != nil {
		return false
	}
	return true
}

func (s *Signer) Expired(token string, minUntilExpire int64) bool {
	crypt := New(s.Secret, Timestamp)
	ts := crypt.Parse([]byte(token))
	return time.Since(ts.Timestamp) > time.Duration(minUntilExpire)*time.Minute

}
