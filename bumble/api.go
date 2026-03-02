package bumble

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dvcrn/matrix-bumble/bumble/bumbletypes"
	"github.com/google/uuid"
	"github.com/tmaxmax/go-sse"
)

type BumbleClient struct {
	Domain                string
	Aid                   string
	PinUnauth             string
	DeviceID              string
	UserID                string
	Session               string
	FirstWebVisitID        string
	LastReferredWebVisitID string
	HTTPClient            *http.Client
}

func NewBumbleClient(domain, aid, deviceID, session, userID, firstWebVisitID, lastReferredWebVisitID string) *BumbleClient {
	return &BumbleClient{
		Domain:                domain,
		Aid:                   aid,
		Session:               session,
		DeviceID:              deviceID,
		UserID:                userID,
		FirstWebVisitID:        firstWebVisitID,
		LastReferredWebVisitID: lastReferredWebVisitID,
		HTTPClient:            &http.Client{},
	}
}

func (c *BumbleClient) cookieHeader(session string) string {
	parts := []string{}
	add := func(key, value string) {
		if value == "" {
			return
		}
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}

	add("HDR-X-User-id", c.UserID)
	add("_pin_unauth", c.PinUnauth)
	add("device_id", c.DeviceID)
	add("aid", c.Aid)
	parts = append(parts, "buzz_lang_code=en-us")
	add("session", session)
	parts = append(parts, "session_cookie_name=session", "has_secure_session=1")
	add("first_web_visit_id", c.FirstWebVisitID)
	add("last_referred_web_visit_id", c.LastReferredWebVisitID)
	return strings.Join(parts, "; ")
}

func (c *BumbleClient) createRequest(method, path string, body interface{}) (*http.Request, error) {
	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("https://%s%s", c.Domain, path)
	fmt.Println("Creating request: ", url)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	cookieString := c.cookieHeader(c.Session)

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3.1 Safari/605.1.15")
	req.Header.Set("X-Message-type", strconv.Itoa(body.(*BadooMessage).MessageType))
	pingback, err := genSecretHeader(body, "")
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Pingback", pingback)
	req.Header.Set("x-use-session-cookie", "1")
	req.Header.Set("Cookie", cookieString)

	return req, nil
}

func (c *BumbleClient) executeRequest(req *http.Request) (*bumbletypes.APIResponse, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyCopy := bytes.Buffer{}
	bodyCopy.ReadFrom(resp.Body)

	// fmt.Println("response: ", bodyCopy.String())

	var apiResponse *bumbletypes.APIResponse
	err = json.NewDecoder(&bodyCopy).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	if len(apiResponse.Body) == 0 {
		return nil, errors.New("no response body")
	}

	if apiErr := apiResponse.GetError(); apiErr != nil {
		return apiResponse, apiErr
	}

	return apiResponse, nil
}

func genSecretHeader(e interface{}, r string) (string, error) {
	salt := os.Getenv("BUMBLE_PINGBACK_SALT")
	if salt == "" {
		return "", errors.New("BUMBLE_PINGBACK_SALT is not set")
	}
	return computeMD5(serialize(e) + (r + salt)), nil
}

func computeMD5(input string) string {
	data := []byte(input)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func serialize(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("%v", data)
	}
	return string(bytes)
}

type BadooMessage struct {
	Gpb          string                   `json:"$gpb"`
	MessageID    int                      `json:"message_id"`
	MessageType  int                      `json:"message_type"`
	Version      int                      `json:"version"`
	IsBackground bool                     `json:"is_background"`
	Body         []map[string]interface{} `json:"body"`
}

type ServerGetUserList struct {
	UserFieldFilter UserFieldFilter `json:"user_field_filter"`
	PreferredCount  int             `json:"preferred_count"`
	Offset          int             `json:"offset"`
	FolderId        int             `json:"folder_id"`
}

type ServerOpenChat struct {
	UserFieldFilter UserFieldFilter `json:"user_field_filter"`
	ChatInstanceID  string          `json:"chat_instance_id"`
	MessageCount    int             `json:"message_count"`
}

type UserFieldFilter struct {
	Projection    []int `json:"projection"`
	RequestAlbums []struct {
		Count        int `json:"count"`
		Offset       int `json:"offset"`
		AlbumType    int `json:"album_type"`
		PhotoRequest struct {
			ReturnPreviewURL bool `json:"return_preview_url"`
			ReturnLargeURL   bool `json:"return_large_url"`
		} `json:"photo_request,omitempty"`
	} `json:"request_albums,omitempty"`
}

func (c *BumbleClient) FetchUserListData(limit int, offset int) (*bumbletypes.ClientUserList, error) {
	bb := &BadooMessage{
		Gpb:          "badoo.bma.BadooMessage",
		MessageID:    5,
		MessageType:  245,
		Version:      1,
		IsBackground: false,
		Body: []map[string]interface{}{
			{
				"message_type": 245,
				"server_get_user_list": ServerGetUserList{
					UserFieldFilter: UserFieldFilter{
						Projection: []int{200, 210, 340, 230, 640, 580, 300, 860, 280, 590, 591, 250, 700, 762, 592, 880, 582, 930, 585, 583, 305, 330, 763, 1423, 584, 1262, 911, 912},
					},
					PreferredCount: limit,
					Offset:         offset,
					// 0 = matches
					// 4 = favorites
					// 6 = liked me
					// 7 = people I swiped right on
					// 8 = ?
					// 9 = maybe blocklist
					FolderId: 0,
				},
			},
		},
	}

	req, err := c.createRequest("POST", "/mwebapi.phtml?SERVER_GET_USER_LIST", bb)
	if err != nil {
		return nil, err
	}

	resp, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}

	return resp.GetClientUserList(), nil
}

func (c *BumbleClient) OpenWebsocket(ctx context.Context, path string, sequence string, userUid string, deviceID string, cb func(*bumbletypes.EventMessage), shouldReconnect bool) {
	socketUrl := fmt.Sprintf("https://%s%sinit?comet=mode:master,type:5,ua:cr&im=mid:0&seq=%s&auth=uid:%s&t=%s", c.Domain, path, sequence, userUid, uuid.NewString())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, socketUrl, http.NoBody)
	if err != nil {
		fmt.Println(err)
		return
	}

	decodedSession, err := url.QueryUnescape(c.Session)
	if err != nil {
		fmt.Printf("failed to decode session: %v", err)
		return
	}

	cookieString := c.cookieHeader(decodedSession)
	fmt.Println("opening websocket: ", socketUrl)

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3.1 Safari/605.1.15")
	req.Header.Set("X-Message-type", "245")
	req.Header.Set("x-use-session-cookie", "1")
	req.Header.Set("Cookie", cookieString)

	client := sse.Client{
		Backoff: sse.Backoff{
			MaxRetries: 0,
		},
		OnRetry: func(err error, duration time.Duration) { fmt.Println("retrying", err) },
	}

reconnect:
	fmt.Println("establishing websocket connection -- ", req.URL.String())
	conn := client.NewConnection(req)
	conn.SubscribeToAll(func(event sse.Event) {
		fmt.Println("received on SUBSCRIBETOALL")
		// TODO: for debug. remove me.
		func(v interface{}) {
			j, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			buf := bytes.NewBuffer(j)
			fmt.Printf("%v\n", buf.String())
		}(event)
	})

	conn.SubscribeMessages(func(event sse.Event) {
		fmt.Println("received event type", event.Type, event.LastEventID, event.Data)
		fmt.Printf("%v\n", event)
		// TODO: for debug. remove me.
		func(v interface{}) {
			j, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			buf := bytes.NewBuffer(j)
			fmt.Printf("%v\n", buf.String())
		}(event)

		parsed, err := ParseBumbleEventSourceMessage([]byte(event.Data))
		if err != nil {
			fmt.Println("error parsing message", err)
			return
		}

		cb(parsed)
	})

	if err := conn.Connect(); !errors.Is(err, context.Canceled) {
		fmt.Println("error connecting to comet", err)
		if shouldReconnect {
			goto reconnect
		}
	}
}

func ParseBumbleEventSourceMessage(eventData []byte) (*bumbletypes.EventMessage, error) {
	var res *bumbletypes.EventMessageWrapper
	err := json.Unmarshal(eventData, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Data) == 0 {
		return nil, errors.New("unexpected event content: " + string(eventData))
	}

	if res.Data[0].(string) != "gpb:push" {
		return nil, errors.New("unexpected event type: " + res.Data[0].(string))
	}

	b, err := json.Marshal(res.Data[1])
	if err != nil {
		return nil, err
	}

	var res2 *bumbletypes.EventMessage
	err = json.Unmarshal(b, &res2)
	if err != nil {
		return nil, err
	}

	return res2, nil
}

func (c *BumbleClient) ClientOpenChatTest(userIds []string, limit int) (*bumbletypes.ClientOpenChat, error) {
	messageId := 1

	body := &BadooMessage{
		Gpb:          "badoo.bma.BadooMessage",
		MessageID:    messageId,
		MessageType:  102,
		Version:      1,
		IsBackground: false,
		Body:         []map[string]interface{}{},
	}

	projection := []int{200, 210, 340, 230, 640, 580, 300, 860, 280, 590, 591, 250, 700, 762, 592, 880, 582, 930, 585, 583, 305, 330, 763, 1423, 584, 1262, 911, 912}

	for _, userId := range userIds {
		body.Body = append(body.Body, map[string]interface{}{
			"message_type": 102,
			"server_open_chat": ServerOpenChat{
				UserFieldFilter: UserFieldFilter{
					Projection: projection,
					RequestAlbums: []struct {
						Count        int `json:"count"`
						Offset       int `json:"offset"`
						AlbumType    int `json:"album_type"`
						PhotoRequest struct {
							ReturnPreviewURL bool `json:"return_preview_url"`
							ReturnLargeURL   bool `json:"return_large_url"`
						} `json:"photo_request,omitempty"`
					}([]struct {
						Count        int `json:"count"`
						Offset       int `json:"offset"`
						AlbumType    int `json:"album_type"`
						PhotoRequest struct {
							ReturnPreviewURL bool `json:"return_preview_url"`
							ReturnLargeURL   bool `json:"return_large_url"`
						} `json:"photo_request"`
					}{
						{
							Count:     10,
							Offset:    1,
							AlbumType: 2,
							PhotoRequest: struct {
								ReturnPreviewURL bool `json:"return_preview_url"`
								ReturnLargeURL   bool `json:"return_large_url"`
							}{
								ReturnPreviewURL: true,
								ReturnLargeURL:   true,
							},
						},
					}),
				},
				ChatInstanceID: userId,
				MessageCount:   limit,
			},
		})
	}

	req, err := c.createRequest("POST", "/mwebapi.phtml?SERVER_OPEN_CHAT", body)
	if err != nil {
		return nil, err
	}

	resp, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}

	return resp.GetClientOpenChat(), nil
}

func (c *BumbleClient) ClientOpenChat(userId string, limit int) (*bumbletypes.ClientOpenChat, error) {
	messageId := 1

	body := &BadooMessage{
		Gpb:          "badoo.bma.BadooMessage",
		MessageID:    messageId,
		MessageType:  102,
		Version:      1,
		IsBackground: false,
		Body:         make([]map[string]interface{}, 1),
	}

	projection := []int{200, 210, 340, 230, 640, 580, 300, 860, 280, 590, 591, 250, 700, 762, 592, 880, 582, 930, 585, 583, 305, 330, 763, 1423, 584, 1262, 911, 912}

	body.Body[0] = map[string]interface{}{
		"message_type": 102,
		"server_open_chat": ServerOpenChat{
			UserFieldFilter: UserFieldFilter{
				Projection: projection,
				RequestAlbums: []struct {
					Count        int `json:"count"`
					Offset       int `json:"offset"`
					AlbumType    int `json:"album_type"`
					PhotoRequest struct {
						ReturnPreviewURL bool `json:"return_preview_url"`
						ReturnLargeURL   bool `json:"return_large_url"`
					} `json:"photo_request,omitempty"`
				}([]struct {
					Count        int `json:"count"`
					Offset       int `json:"offset"`
					AlbumType    int `json:"album_type"`
					PhotoRequest struct {
						ReturnPreviewURL bool `json:"return_preview_url"`
						ReturnLargeURL   bool `json:"return_large_url"`
					} `json:"photo_request"`
				}{
					{
						Count:     10,
						Offset:    1,
						AlbumType: 2,
						PhotoRequest: struct {
							ReturnPreviewURL bool `json:"return_preview_url"`
							ReturnLargeURL   bool `json:"return_large_url"`
						}{
							ReturnPreviewURL: true,
							ReturnLargeURL:   true,
						},
					},
				}),
			},
			ChatInstanceID: userId,
			MessageCount:   limit,
		},
	}

	req, err := c.createRequest("POST", "/mwebapi.phtml?SERVER_OPEN_CHAT", body)
	if err != nil {
		return nil, err
	}

	resp, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}

	return resp.GetClientOpenChat(), nil
}

func (c *BumbleClient) sendMessage(fromPersonId, toPersonId, msg string) (*bumbletypes.ChatMessageReceived, error) {
	bb := &BadooMessage{
		Gpb:          "badoo.bma.BadooMessage",
		MessageID:    5,
		MessageType:  104,
		Version:      1,
		IsBackground: false,
		Body: []map[string]interface{}{
			{
				"message_type": 104,
				"chat_message": map[string]interface{}{
					"mssg":           msg,
					"message_type":   1,
					"uid":            strconv.FormatInt(time.Now().Unix()*1000, 10),
					"from_person_id": fromPersonId,
					"to_person_id":   toPersonId,
					"read":           false,
				},
			},
		},
	}

	req, err := c.createRequest("POST", "/mwebapi.phtml?SERVER_SEND_CHAT_MESSAGE", bb)
	if err != nil {
		return nil, err
	}

	resp, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}

	return resp.GetChatMessageReceived(), nil
}

type FetchAppStartResponse struct {
	ClientStartup        *bumbletypes.ClientStartup
	ClientCommonSettings *bumbletypes.ClientCommonSettings
	ClientLoginSuccess   *bumbletypes.ClientLoginSuccess
	ClientSessionChanged *bumbletypes.ClientSessionChanged
	User                 *bumbletypes.User
	AppSettings          *bumbletypes.AppSettings
	CometConfiguration   *bumbletypes.CometConfiguration
}

func (c *BumbleClient) FetchAppStart() (*FetchAppStartResponse, error) {
	bb := &BadooMessage{
		Gpb:          "badoo.bma.BadooMessage",
		MessageID:    1,
		MessageType:  2,
		Version:      1,
		IsBackground: false,
		Body: []map[string]interface{}{
			{
				"message_type": 2,
				"server_app_startup": map[string]interface{}{
					"app_build":                      "MoxieWebapp",
					"app_name":                       "moxie",
					"app_version":                    "1.0.0",
					"can_send_sms":                   false,
					"user_agent":                     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Safari/605.1.15",
					"screen_width":                   1680,
					"screen_height":                  1050,
					"language":                       0,
					"is_cold_start":                  true,
					"external_provider_redirect_url": "https://bumble.com/static/external-auth-result.html?",
					"locale":                         "en-au",
					"system_locale":                  "en-US",
					"app_platform_type":              5,
					"app_product_type":               400,
					"device_info": map[string]interface{}{
						"webcam_available": true,
						"form_factor":      3,
					},
					"build_configuration":         2,
					"build_fingerprint":           "31214",
					"supported_features":          []int{141, 145, 11, 15, 1, 2, 13, 46, 4, 248, 6, 18, 155, 70, 160, 140, 130, 189, 187, 220, 223, 180, 197, 161, 232, 29, 227, 237, 239, 254, 190, 290, 291, 296, 250, 264, 294, 295, 310, 100, 148, 262},
					"supported_minor_features":    []int{472, 317, 2, 216, 244, 232, 19, 130, 225, 246, 31, 125, 183, 114, 254, 8, 9, 83, 41, 427, 115, 288, 420, 477, 93, 226, 413, 267, 39, 290, 398, 453, 180, 281, 40, 455, 280, 499, 471, 397, 411, 352, 447, 146, 469, 118, 63, 391, 523, 293, 431, 620, 574, 405, 547, 451, 571, 319, 297, 558, 394, 593, 628, 603, 602, 537, 305, 561, 324, 554, 505, 696, 576, 707, 726, 624, 797, 821, 829, 651, 148, 920, 916, 877, 936, 657, 935, 742, 941, 309, 329, 307, 553},
					"supported_notifications":     []int{83, 73, 3, 72, 49, 46, 109, 81, 44, 96, 89},
					"supported_payment_providers": []int{26, 100, 35, 100001, 191, 194, 83, 192, 195, 180},
					"supported_promo_blocks": []map[string]interface{}{
						{"context": 92, "position": 13, "types": []int{71}},
						{"context": 45, "position": 21, "types": []int{148}},
						{"context": 89, "position": 5, "types": []int{160, 358}},
						{"context": 8, "position": 13, "types": []int{111, 112, 113}},
						{"context": 53, "position": 18, "types": []int{136, 93, 12}},
						{"context": 45, "position": 18, "types": []int{327}},
						{"context": 45, "position": 15, "types": []int{410, 93, 134, 135, 136, 137, 327, 308, 309, 334, 187, 61, 422, 423}},
						{"context": 10, "position": 1, "types": []int{265, 266, 286}},
						{"context": 148, "position": 21, "types": []int{179, 180, 283}},
						{"context": 26, "position": 13, "types": []int{354}},
						{"context": 26, "position": 4, "types": []int{355, 356}},
						{"context": 26, "position": 1, "types": []int{354}},
						{"context": 26, "position": 18, "types": []int{357}},
						{"context": 130, "position": 13, "types": []int{268, 267}},
						{"context": 113, "position": 1, "types": []int{228}},
						{"context": 3, "position": 1, "types": []int{80, 423}},
						{"context": 3, "position": 4, "types": []int{80, 228, 423}},
						{"context": 119, "position": 1, "types": []int{80, 282, 81, 90, 422, 140}},
						{"context": 43, "position": 1, "types": []int{96, 307}},
						{"context": 43, "position": 18, "types": []int{369}},
						{"context": 119, "position": 18, "types": []int{369}},
						{"context": 10, "position": 18, "types": []int{358, 174}},
						{"context": 10, "position": 8, "types": []int{358}},
						{"context": 26, "position": 16, "types": []int{286, 371}},
						{"context": 10, "position": 6, "types": []int{286, 373, 372}},
						{"context": 246, "position": 13, "types": []int{404}},
					},
					"supported_user_substitutes": []map[string]interface{}{
						{"context": 1, "types": []int{3}},
					},
					"supported_onboarding_types": []int{9},
					"user_field_filter_client_login_success": map[string]interface{}{
						"projection": []int{210, 220, 230, 200, 91, 890, 340, 10, 11, 231, 71, 93, 100},
					},
					"a_b_testing_settings": map[string]interface{}{
						"tests": []map[string]interface{}{
							{"test_id": "bumble__gifs_with_old_input"},
							{"test_id": "venmo_new_desktop_flow"},
						},
					},
					"dev_features": []string{"bumble_bizz", "bumble_snooze", "bumble_questions", "bumble__pledge", "bumble__request_photo_verification", "bumble_moves_making_impact_", "bumble__photo_verification_filters", "bumble_gift_cards", "bumble__antighosting_xp_dead_chat_followup", "bumble_private_detector", "bumble_distance_expansion", "bumble_live_in_the_hive", "bumble__shared_experiences__recommend_to_a_friend__v1", "arkose.integration"},
					"device_id":    c.DeviceID,
					"supported_screens": []map[string]interface{}{
						{"type": 23, "version": 4},
						{"type": 26, "version": 0},
						{"type": 13, "version": 0},
						{"type": 14, "version": 0},
						{"type": 15, "version": 0},
						{"type": 16, "version": 0},
						{"type": 18, "version": 0},
						{"type": 19, "version": 0},
						{"type": 20, "version": 0},
						{"type": 21, "version": 0},
						{"type": 25, "version": 1},
						{"type": 27, "version": 0},
						{"type": 28, "version": 0},
						{"type": 57, "version": 0},
						{"type": 29, "version": 1},
						{"type": 69, "version": 0},
						{"type": 297, "version": 0},
						{"type": 298, "version": 0},
						{"type": 92, "version": 0},
						{"type": 378, "version": 0},
						{"type": 379, "version": 0},
						{"type": 382, "version": 0},
						{"type": 386, "version": 1},
						{"type": 387, "version": 0},
						{"type": 488, "version": 0},
						{"type": 401, "version": 1},
						{"type": 412, "version": 0},
						{"type": 503, "version": 0},
						{"type": 505, "version": 0},
						{"type": 507, "version": 0},
						{"type": 522, "version": 0},
						{"type": 397, "version": 0},
						{"type": 399, "version": 0},
						{"type": 536, "version": 0},
						{"type": 534, "version": 0},
					},
					"supported_landings": []map[string]interface{}{
						{"source": 25, "params": []int{20, 3}, "search_settings_types": []int{3}},
						{"source": 129, "params": []int{3}},
					},
					"app_domain": "com.bumble",
				},
			},
		},
	}

	req, err := c.createRequest("POST", "/mwebapi.phtml?SERVER_APP_STARTUP", bb)
	if err != nil {
		return nil, err
	}

	resp, err := c.executeRequest(req)
	if err != nil {
		// check if it's an API error, or actual error
		if resp.GetError() != nil {
			// if it's an API error, check for error code 79 == Message: Client does not support all mandatory screens (0030-0002-0003), ID: 0030-0002-0003, ETA: 562, Type: 79
			// then ignore the error and proceed
			if resp.GetError().ErrorCode != "79" {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	res := FetchAppStartResponse{
		ClientStartup:        resp.GetClientStartup(),
		ClientCommonSettings: resp.GetClientCommonSettings(),
		ClientLoginSuccess:   resp.GetClientLoginSuccess(),
		ClientSessionChanged: resp.GetClientSessionChanged(),
		AppSettings:          resp.GetAppSettings(),
		CometConfiguration:   resp.GetCometConfiguration(),
	}

	user := resp.GetUser()
	if user != nil {
		res.User = &user.User
	}

	return &res, nil
}

func (c *BumbleClient) ServerUpdateSession(hotpanelSessionID string) error {
	bb := &BadooMessage{
		Gpb:          "badoo.bma.BadooMessage",
		MessageID:    13,
		MessageType:  199,
		Version:      1,
		IsBackground: false,
		Body: []map[string]interface{}{
			{
				"message_type": 199,
				"server_update_session": map[string]string{
					"hotpanel_session_id": hotpanelSessionID,
				},
			},
		},
	}

	req, err := c.createRequest("POST", "/mwebapi.phtml?SERVER_UPDATE_SESSION", bb)
	if err != nil {
		return err
	}

	_, err = c.executeRequest(req)
	if err != nil {
		return err
	}

	return nil
}

// func (c *BumbleClient) ServerAppStartup() (*bumbletypes.APIResponse, error) {
// 	bb := &BadooMessage{
// 		Gpb:          "badoo.bma.BadooMessage",
// 		MessageID:    1,
// 		MessageType:  2,
// 		Version:      1,
// 		IsBackground: false,
// 		Body: []map[string]interface{}{
// 			{
// 				"message_type": 2,
// 				"server_app_startup": map[string]interface{}{
// 					"app_build":                      "MoxieWebapp",
// 					"app_name":                       "moxie",
// 					"app_version":                    "1.0.0",
// 					"can_send_sms":                   false,
// 					"user_agent":                     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
// 					"screen_width":                   1680,
// 					"screen_height":                  1050,
// 					"language":                       0,
// 					"is_cold_start":                  true,
// 					"external_provider_redirect_url": "https://fr1.bumble.com/static/external-auth-result.html?",
// 					"locale":                         "en-au",
// 					"system_locale":                  "en-US",
// 					"app_platform_type":              5,
// 					"app_product_type":               400,
// 					"device_info": map[string]interface{}{
// 						"webcam_available": true,
// 						"form_factor":      3,
// 					},
// 					"build_configuration":         2,
// 					"build_fingerprint":           "31139",
// 					"supported_features":          []int{141, 145, 11, 15, 1, 2, 13, 46, 4, 248, 6, 18, 155, 70, 160, 140, 130, 189, 187, 220, 223, 180, 197, 161, 232, 29, 227, 237, 239, 254, 190, 290, 291, 296, 250, 264, 294, 295, 310, 100, 148, 262},
// 					"supported_minor_features":    []int{472, 317, 2, 216, 244, 232, 19, 130, 225, 246, 31, 125, 183, 114, 254, 8, 9, 83, 41, 427, 115, 288, 420, 477, 93, 226, 413, 267, 39, 290, 398, 453, 180, 281, 40, 455, 280, 499, 471, 397, 411, 352, 447, 146, 469, 118, 63, 391, 523, 293, 431, 620, 574, 405, 547, 451, 571, 319, 297, 558, 394, 593, 628, 603, 602, 537, 305, 561, 324, 554, 505, 696, 576, 707, 726, 624, 797, 821, 829, 651, 148, 920, 916, 877, 936, 657, 935, 742, 941, 309, 329, 307, 553},
// 					"supported_notifications":     []int{83, 73, 3, 72, 49, 46, 109, 81, 44, 96, 89},
// 					"supported_payment_providers": []int{26, 100, 35, 100001, 191, 194, 83, 192, 195},
// 					"supported_promo_blocks": []map[string]interface{}{
// 						{"context": 92, "position": 13, "types": []int{71}},
// 						{"context": 45, "position": 21, "types": []int{148}},
// 						{"context": 89, "position": 5, "types": []int{160, 358}},
// 						{"context": 8, "position": 13, "types": []int{111, 112, 113}},
// 						{"context": 53, "position": 18, "types": []int{136, 93, 12}},
// 						{"context": 45, "position": 18, "types": []int{327}},
// 						{"context": 45, "position": 15, "types": []int{410, 93, 134, 135, 136, 137, 327, 308, 309, 334, 187, 61, 422, 423}},
// 						{"context": 10, "position": 1, "types": []int{265, 266, 286}},
// 						{"context": 148, "position": 21, "types": []int{179, 180, 283}},
// 						{"context": 26, "position": 13, "types": []int{354}},
// 						{"context": 26, "position": 4, "types": []int{355, 356}},
// 						{"context": 26, "position": 1, "types": []int{354}},
// 						{"context": 26, "position": 18, "types": []int{357}},
// 						{"context": 130, "position": 13, "types": []int{268, 267}},
// 						{"context": 113, "position": 1, "types": []int{228}},
// 						{"context": 3, "position": 1, "types": []int{80, 423}},
// 						{"context": 3, "position": 4, "types": []int{80, 228, 423}},
// 						{"context": 119, "position": 1, "types": []int{80, 282, 81, 90, 422, 140}},
// 						{"context": 43, "position": 1, "types": []int{96, 307}},
// 						{"context": 43, "position": 18, "types": []int{369}},
// 						{"context": 119, "position": 18, "types": []int{369}},
// 						{"context": 10, "position": 18, "types": []int{358, 174}},
// 						{"context": 10, "position": 8, "types": []int{358}},
// 						{"context": 26, "position": 16, "types": []int{286, 371}},
// 						{"context": 10, "position": 6, "types": []int{286, 373, 372}},
// 						{"context": 246, "position": 13, "types": []int{404}},
// 					},
// 					"supported_user_substitutes": []map[string]interface{}{
// 						{"context": 1, "types": []int{3}},
// 					},
// 					"supported_onboarding_types": []int{9},
// 					"user_field_filter_client_login_success": map[string]interface{}{
// 						"projection": []int{210, 220, 230, 200, 91, 890, 340, 10, 11, 231, 71, 93, 100},
// 					},
// 					"a_b_testing_settings": map[string]interface{}{
// 						"tests": []map[string]string{
// 							{"test_id": "bumble__gifs_with_old_input"},
// 							{"test_id": "venmo_new_desktop_flow"},
// 						},
// 					},
// 					"dev_features": []string{
// 						"bumble_bizz", "bumble_snooze", "bumble_questions", "bumble__pledge",
// 						"bumble__request_photo_verification", "bumble_moves_making_impact_",
// 						"bumble__photo_verification_filters", "bumble_gift_cards",
// 						"bumble__antighosting_xp_dead_chat_followup", "bumble_private_detector",
// 						"bumble_distance_expansion", "bumble_live_in_the_hive",
// 						"bumble__shared_experiences__recommend_to_a_friend__v1",
// 					},
// 					"device_id": c.DeviceID,
// 					"supported_screens": []map[string]int{
// 						{"type": 23, "version": 4}, {"type": 26, "version": 0}, {"type": 13, "version": 0},
// 						{"type": 14, "version": 0}, {"type": 15, "version": 0}, {"type": 16, "version": 0},
// 						{"type": 18, "version": 0}, {"type": 19, "version": 0}, {"type": 20, "version": 0},
// 						{"type": 21, "version": 0}, {"type": 25, "version": 1}, {"type": 27, "version": 0},
// 						{"type": 28, "version": 0}, {"type": 57, "version": 0}, {"type": 29, "version": 1},
// 						{"type": 69, "version": 0}, {"type": 297, "version": 0}, {"type": 298, "version": 0},
// 						{"type": 92, "version": 0}, {"type": 378, "version": 0}, {"type": 379, "version": 0},
// 						{"type": 382, "version": 0}, {"type": 386, "version": 1}, {"type": 387, "version": 0},
// 						{"type": 488, "version": 0}, {"type": 401, "version": 1}, {"type": 412, "version": 0},
// 						{"type": 503, "version": 0}, {"type": 505, "version": 0}, {"type": 507, "version": 0},
// 						{"type": 522, "version": 0}, {"type": 397, "version": 0}, {"type": 399, "version": 0},
// 						{"type": 536, "version": 0}, {"type": 534, "version": 0},
// 					},
// 					"supported_landings": []map[string]interface{}{
// 						{"source": 25, "params": []int{20, 3}, "search_settings_types": []int{3}},
// 						{"source": 129, "params": []int{3}},
// 					},
// 					"app_domain": "com.bumble.fr1",
// 				},
// 			},
// 		},
// 	}

// 	req, err := c.createRequest("POST", "/mwebapi.phtml?SERVER_APP_STARTUP", bb)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp, err := c.executeRequest(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp.getapp, nil
// }

func (c *BumbleClient) GetUser(userID string) (*bumbletypes.UserExtended, error) {
	body := &BadooMessage{
		// "$gpb":          "badoo.bma.BadooMessage",
		// "message_id":    19,
		// "message_type":  403,
		// "version":       1,
		// "is_background": false,
		Gpb:          "badoo.bma.BadooMessage",
		MessageID:    19,
		MessageType:  403,
		Version:      1,
		IsBackground: false,
		Body: []map[string]interface{}{
			{
				"message_type": 403,
				"server_get_user": map[string]interface{}{
					"user_id": userID,
					"user_field_filter": map[string]interface{}{
						"game_mode":  0,
						"projection": []int{200, 340, 230, 310, 370, 762, 890, 493, 530, 540, 291, 490, 1160, 1161, 210, 380},
						"request_music_services": map[string]interface{}{
							"top_artists_limit":  10,
							"supported_services": []int{29},
						},
						"request_albums": []map[string]interface{}{
							{
								"person_id":  userID,
								"album_type": 2,
								"offset":     1,
							},
							{
								"person_id":         userID,
								"album_type":        12,
								"external_provider": 12,
							},
						},
					},
					"client_source": 10,
				},
			},
		},
	}

	req, err := c.createRequest("POST", "/mwebapi.phtml?SERVER_GET_USER", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Message-type", "403")
	req.Header.Set("x-use-session-cookie", "1")

	resp, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}

	return resp.GetUser(), nil
}
