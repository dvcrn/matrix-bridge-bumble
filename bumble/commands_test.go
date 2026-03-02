package bumble

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractCookie(t *testing.T) {
	curlString := "curl 'https://fr1.bumble.com/mwebapi.phtml?SERVER_APP_STARTUP' -H 'Cookie: aid=123; HDR-X-User-id=USER_ID; session=s2%3A195%3AREDACTED; device_id=DEVICE_ID; first_web_visit_id=FIRST; last_referred_web_visit_id=LAST'"
	cookieString := extractCookie(curlString)
	assert.Equal(t, "aid=123; HDR-X-User-id=USER_ID; session=s2%3A195%3AREDACTED; device_id=DEVICE_ID; first_web_visit_id=FIRST; last_referred_web_visit_id=LAST", cookieString)
}

func TestExtractDomain(t *testing.T) {
	curlString := "curl 'https://fr1.bumble.com/mwebapi.phtml?SERVER_APP_STARTUP' -X 'POST'"
	domain := extractDomain(curlString)
	assert.Equal(t, "fr1.bumble.com", domain)
}
