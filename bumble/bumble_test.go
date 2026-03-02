package bumble

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEventResponse(t *testing.T) {
	expectedMessage := "(Well I mean you’re still here so)"

	testData := []byte(`{
  "cmd": "ev",
  "nd": 1,
  "seq": "0",
  "data": [
    "gpb:push",
    {
      "$gpb": "badoo.bma.BadooMessage",
      "message_type": 6004,
      "version": 0,
      "message_id": 0,
      "responses_count": 1,
      "body": [
        {
          "$gpb": "badoo.bma.MessageBody",
          "message_type": 6004,
          "chat_message": {
            "$gpb": "badoo.bma.ChatMessage",
            "uid": "1",
            "date_modified": 0,
            "from_person_id": "from",
            "to_person_id": "to",
            "mssg": "` + expectedMessage + `"
          }
        }
      ]
    }
  ]
}`)

	res, err := ParseBumbleEventSourceMessage(testData)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, res.Body[0].ChatMessage.Mssg)
}

func TestGenSecretHeader(t *testing.T) {
	t.Setenv("BUMBLE_PINGBACK_SALT", "test-salt")

	tdString := `{"$gpb":"badoo.bma.BadooMessage","body":[{"message_type":245,"server_get_user_list":{"user_field_filter":{"projection":[200,210,340,230,640,580,300,860,280,590,591,250,700,762,592,880,582,930,585,583,305,330,763,1423,584,1262,911,912]},"preferred_count":30,"folder_id":0}}],"message_id":4,"message_type":245,"version":1,"is_background":false}`

	var asGolang *map[string]interface{}
	err := json.Unmarshal([]byte(tdString), &asGolang)
	assert.NoError(t, err)

	secretHeader, err := genSecretHeader(*asGolang, "")
	assert.NoError(t, err)
	assert.Equal(t, computeMD5(serialize(*asGolang)+"test-salt"), secretHeader)
}
