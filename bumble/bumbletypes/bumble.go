package bumbletypes

import (
	"fmt"
	"reflect"
	"strings"
)

type Photo struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	PreviewURL     string `json:"preview_url"`
	LargeURL       string `json:"large_url"`
	LargePhotoSize struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"large_photo_size"`
	PreviewURLExpirationTS int `json:"preview_url_expiration_ts"`
	LargeURLExpirationTS   int `json:"large_url_expiration_ts"`
}

type UserExtended struct {
	User
	ClientSource       int            `json:"client_source"`
	VerificationStatus int            `json:"verification_status"`
	PhotoCount         int            `json:"photo_count"`
	Albums             []Album        `json:"albums"`
	ProfileFields      []ProfileField `json:"profile_fields"`
	ProfileSummary     ProfileSummary `json:"profile_summary"`
	DistanceLong       string         `json:"distance_long"`
	DistanceShort      string         `json:"distance_short"`
	Hometown           Location       `json:"hometown"`
	Residence          Location       `json:"residence"`
}

func (buser *UserExtended) DumpProfile() string {
	profileFields := []string{}
	gender := "male"
	if buser.Gender == 2 {
		gender = "female"
	}
	for _, field := range buser.ProfileFields {
		profileFields = append(profileFields, fmt.Sprintf("%s: %s", field.Name, field.DisplayValue))
	}

	buserInfo := fmt.Sprintf(`Name: %s
Age: %d
Gender: %s
Living in %s: City: %s / Country: %s
From: City: %s / Country: %s

Distance: %s %s

Profile Fields:
%s
	`, buser.Name, buser.Age, gender, buser.Residence.ContextInfo, buser.Residence.City.Name, buser.Residence.Country.Name, buser.Hometown.City.Name, buser.Hometown.Country.Name, buser.DistanceLong, buser.DistanceShort, strings.Join(profileFields, "\n"))

	return buserInfo
}

type Album struct {
	UID                string  `json:"uid"`
	Name               string  `json:"name"`
	OwnerID            string  `json:"owner_id"`
	AccessType         int     `json:"access_type"`
	Accessable         bool    `json:"accessable"`
	Adult              bool    `json:"adult"`
	RequiresModeration bool    `json:"requires_moderation"`
	CountOfPhotos      int     `json:"count_of_photos"`
	IsUploadForbidden  bool    `json:"is_upload_forbidden"`
	Photos             []Photo `json:"photos"`
	AlbumType          int     `json:"album_type"`
	GameMode           int     `json:"game_mode"`
	Caption            string  `json:"caption,omitempty"`
	ExternalProvider   int     `json:"external_provider,omitempty"`
}

type ProfileField struct {
	ID             string `json:"id"`
	Type           int    `json:"type"`
	Name           string `json:"name"`
	DisplayValue   string `json:"display_value"`
	RequiredAction int    `json:"required_action"`
	IconURL        string `json:"icon_url,omitempty"`
	HPElement      int    `json:"hp_element,omitempty"`
}

type ProfileSummary struct {
	GPB string `json:"$gpb"`
}

type Location struct {
	GPB         string  `json:"$gpb"`
	Type        int     `json:"type"`
	Country     Country `json:"country"`
	Region      Region  `json:"region"`
	City        City    `json:"city"`
	ContextInfo string  `json:"context_info"`
}

type Country struct {
	GPB         string `json:"$gpb"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhonePrefix string `json:"phone_prefix"`
	ISOCode     string `json:"iso_code"`
	FlagSymbol  string `json:"flag_symbol"`
	PhoneLength Range  `json:"phone_length"`
}

type Region struct {
	GPB          string `json:"$gpb"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

type City struct {
	GPB         string `json:"$gpb"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ContextInfo string `json:"context_info"`
}

type Range struct {
	GPB      string `json:"$gpb"`
	MinValue int    `json:"min_value"`
	MaxValue int    `json:"max_value"`
}

type User struct {
	UserID           string `json:"user_id"`
	Projection       []int  `json:"projection"`
	AccessLevel      int    `json:"access_level"`
	Name             string `json:"name"`
	Age              int    `json:"age"`
	Gender           int    `json:"gender"`
	IsDeleted        bool   `json:"is_deleted"`
	IsExtendedMatch  bool   `json:"is_extended_match"`
	MatchExtenderID  string `json:"match_extender_id,omitempty"`
	OnlineStatus     int    `json:"online_status"`
	ProfilePhoto     Photo  `json:"profile_photo"`
	IsMatch          bool   `json:"is_match"`
	MatchMode        int    `json:"match_mode"`
	IsCrush          bool   `json:"is_crush"`
	TheirVoteMode    int    `json:"their_vote_mode,omitempty"`
	PreMatchTimeLeft *struct {
		Goal      int  `json:"goal"`
		Progress  int  `json:"progress"`
		StartTS   int  `json:"start_ts"`
		IsDelayed bool `json:"is_delayed"`
	} `json:"pre_match_time_left,omitempty"`
	ReplyTimeLeft *struct {
		Goal     int `json:"goal"`
		Progress int `json:"progress"`
		StartTS  int `json:"start_ts"`
	} `json:"reply_time_left,omitempty"`
	UnreadMessagesCount        int    `json:"unread_messages_count"`
	DisplayMessage             string `json:"display_message,omitempty"`
	IsInappPromoPartner        bool   `json:"is_inapp_promo_partner"`
	IsLocked                   bool   `json:"is_locked"`
	Type                       int    `json:"type"`
	ConnectionStatusIndicator  int    `json:"connection_status_indicator"`
	RematchAction              int    `json:"rematch_action,omitempty"`
	ConnectionExpiredTimestamp int    `json:"connection_expired_timestamp,omitempty"`
}

func (u *User) NormalizeDisplayName() string {
	return fmt.Sprintf("%s 🐝", u.Name)
}

func (u *User) NormalizeUsername() string {
	// return a combination of: short userid, only 10 characters, together with name
	id := u.UserID
	return strings.ToLower(id)
}

type ListSection struct {
	SectionID      string `json:"section_id"`
	Name           string `json:"name"`
	TotalCount     int    `json:"total_count"`
	LastBlock      bool   `json:"last_block"`
	AllowedActions []int  `json:"allowed_actions"`
	SectionType    int    `json:"section_type,omitempty"`
	Users          []User `json:"users"`
	PromoBanners   []any  `json:"promo_banners,omitempty"`
	UserFeature    struct {
		Feature int  `json:"feature"`
		Enabled bool `json:"enabled"`
	} `json:"user_feature,omitempty"`
}

type ClientUserList struct {
	Section       []ListSection `json:"section"`
	TotalSections int           `json:"total_sections"`
	TotalCount    int           `json:"total_count"`
	DelaySec      int           `json:"delay_sec"`
}

type ChatMessageMultimedia struct {
	Format     int                  `json:"format"`
	Visibility MultimediaVisibility `json:"visibility"`
	Photo      Photo                `json:"photo"`
}

type ChatMessage struct {
	UID                     string                   `json:"uid"`
	DateModified            int                      `json:"date_modified"`
	FromPersonID            string                   `json:"from_person_id"`
	ToPersonID              string                   `json:"to_person_id"`
	Mssg                    string                   `json:"mssg"`
	MessageType             int                      `json:"message_type"`
	Read                    bool                     `json:"read"`
	AlbumID                 string                   `json:"album_id"`
	TotalUnread             int                      `json:"total_unread"`
	UnreadFromUser          int                      `json:"unread_from_user"`
	ImageURL                string                   `json:"image_url"`
	FrameURL                string                   `json:"frame_url"`
	CanDelete               bool                     `json:"can_delete"`
	Deleted                 bool                     `json:"deleted"`
	FromPersonInfo          ChatUserInfo             `json:"from_person_info"`
	Sticker                 any                      `json:"sticker"`
	Gift                    any                      `json:"gift"`
	Offensive               bool                     `json:"offensive"`
	DisplayMessage          string                   `json:"display_message"`
	VerificationMethod      VerificationAccessObject `json:"verification_method"`
	DateCreated             int                      `json:"date_created"`
	AccessResponseType      int                      `json:"access_response_type"`
	IsLiked                 bool                     `json:"is_liked"`
	ReplyToUID              string                   `json:"reply_to_uid"`
	FirstResponse           bool                     `json:"first_response"`
	VideoCallMsgInfo        VideoCallMsgInfo         `json:"video_call_msg_info"`
	IsMasked                bool                     `json:"is_masked"`
	AllowReport             bool                     `json:"allow_report"`
	EmojisCount             int                      `json:"emojis_count"`
	HasEmojiCharactersOnly  bool                     `json:"has_emoji_characters_only"`
	UserSubstituteID        string                   `json:"user_substitute_id"`
	AllowReply              bool                     `json:"allow_reply"`
	AllowEditUntilTimestamp int                      `json:"allow_edit_until_timestamp"`
	IsEdited                bool                     `json:"is_edited"`
	AllowForwarding         bool                     `json:"allow_forwarding"`
	ClearChatVersion        int                      `json:"clear_chat_version"`
	Story                   Story                    `json:"story"`
	IsDeclined              bool                     `json:"is_declined"`
	HasLewdPhoto            bool                     `json:"has_lewd_photo"`
	IsReported              bool                     `json:"is_reported"`
	AllowLike               bool                     `json:"allow_like"`
	IsLegacy                bool                     `json:"is_legacy"`
	IsLikelyOffensive       bool                     `json:"is_likely_offensive"`
	GameID                  string                   `json:"game_id"`
	ExperimentalGift        ExperimentalGift         `json:"experimental_gift"`
	Multimedia              ChatMessageMultimedia    `json:"multimedia,omitempty"`
}

type ChatInstance struct {
	GPB                  string `json:"$gpb"`
	UID                  string `json:"uid"`
	DateModified         int    `json:"date_modified"`
	Counter              int    `json:"counter"`
	TheirIconID          string `json:"their_icon_id"`
	MyIconID             string `json:"my_icon_id"`
	OtherAccountDeleted  bool   `json:"other_account_deleted"`
	IsNew                bool   `json:"is_new"`
	FeelsLikeChatting    bool   `json:"feels_like_chatting"`
	MyUnreadMessages     int    `json:"my_unread_messages"`
	TheirUnreadMessages  int    `json:"their_unread_messages"`
	ChatIcebreakerAsk    string `json:"chat_icebreaker_ask"`
	ChatSuggestionPrompt string `json:"chat_suggestion_prompt"`
	IsMatch              bool   `json:"is_match"`
	OpenStickers         bool   `json:"open_stickers"`
}

type ChatSettings struct {
	GPB                           string              `json:"$gpb"`
	ChatInstanceID                string              `json:"chat_instance_id"`
	MultimediaSettings            *MultimediaSettings `json:"multimedia_settings"`
	FeatureOrder                  []int               `json:"feature_order"`
	InputSettings                 *InputSettings      `json:"input_settings"`
	AllowDisablingPrivateDetector bool                `json:"allow_disabling_private_detector"`
	AllowQuestionsGame            bool                `json:"allow_questions_game"`
	IsGoodChat                    bool                `json:"is_good_chat"`
}

type MultimediaSettings struct {
	GPB               string             `json:"$gpb"`
	Feature           ApplicationFeature `json:"feature"`
	MultimediaConfig  MultimediaConfig   `json:"multimedia_config"`
	MultimediaConfigs []MultimediaConfig `json:"multimedia_configs"`
}

type MultimediaConfig struct {
	GPB        string                 `json:"$gpb"`
	Visibility []MultimediaVisibility `json:"visibility"`
	Format     int                    `json:"format"`
	MinLength  int                    `json:"min_length,omitempty"`
	MaxLength  int                    `json:"max_length,omitempty"`
}

type MultimediaVisibility struct {
	GPB            string `json:"$gpb"`
	VisibilityType int    `json:"visibility_type"`
	Seconds        int    `json:"seconds,omitempty"`
	DisplayValue   string `json:"display_value"`
}

type InputSettings struct {
	GPB           string               `json:"$gpb"`
	InputFeatures []ApplicationFeature `json:"input_features"`
}

type ClientOpenChat struct {
	IsChatAvailable       bool           `json:"is_chat_available"`
	UserOriginatedMessage bool           `json:"user_originated_message"`
	Title                 string         `json:"title"`
	ChatMessages          []*ChatMessage `json:"chat_messages"`
	ChatUser              *User          `json:"chat_user"`
	EncryptedIMWriting    string         `json:"encrypted_im_writing"`
	EncryptedCometURL     string         `json:"encrypted_comet_url"`
	ReadMessagesTimestamp int            `json:"read_messages_timestamp"`
	IsNotInterested       bool           `json:"is_not_interested"`
	ChatSettings          *ChatSettings  `json:"chat_settings"`
	ChatInstance          *ChatInstance  `json:"chat_instance"`
}

type ChatMessageReceived struct {
	UID         string       `json:"uid"`
	Success     bool         `json:"success"`
	ChatMessage *ChatMessage `json:"chat_message"`
}

type ServerErrorMessage struct {
	GPB          string `json:"$gpb"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	ErrorID      string `json:"error_id"`
	ErrorETA     int    `json:"error_eta"`
	Type         int    `json:"type"`
}

func (e *ServerErrorMessage) Error() string {
	return fmt.Sprintf("Error Code: %s, Message: %s, ID: %s, ETA: %d, Type: %d", e.ErrorCode, e.ErrorMessage, e.ErrorID, e.ErrorETA, e.Type)
}

type MessageBody struct {
	ClientUserList *ClientUserList `json:"client_user_list,omitempty"`
	ClientOpenChat *ClientOpenChat `json:"client_open_chat,omitempty"`
	User           *UserExtended   `json:"user,omitempty"`
	// UserExtended         *UserExtended         `json:"user,omitempty"`
	ChatMessageReceived  *ChatMessageReceived  `json:"chat_message_received,omitempty"`
	ClientStartup        *ClientStartup        `json:"client_startup,omitempty"`
	ClientCommonSettings *ClientCommonSettings `json:"client_common_settings,omitempty"`
	ClientLoginSuccess   *ClientLoginSuccess   `json:"client_login_success,omitempty"`
	ClientSessionChanged *ClientSessionChanged `json:"client_session_changed,omitempty"`
	AppSettings          *AppSettings          `json:"app_settings,omitempty"`
	CometConfiguration   *CometConfiguration   `json:"comet_configuration,omitempty"`
	ServerErrorMessage   *ServerErrorMessage   `json:"server_error_message,omitempty"`
	MessageType          int                   `json:"message_type"`
}

type APIResponse struct {
	GPB            string        `json:"$gpb"`
	MessageType    int           `json:"message_type"`
	Version        int           `json:"version"`
	MessageID      int           `json:"message_id"`
	ObjectType     int           `json:"object_type"`
	Body           []MessageBody `json:"body"`
	ResponsesCount int           `json:"responses_count"`
	IsBackground   bool          `json:"is_background"`
	Vhost          string        `json:"vhost"`
}

func (a *APIResponse) getBodyField(fieldName string) interface{} {
	if len(a.Body) == 0 {
		return nil
	}

	for _, body := range a.Body {
		value := reflect.ValueOf(body).FieldByName(fieldName)
		if value.IsValid() && !value.IsNil() {
			return value.Interface()
		}
	}

	return nil
}

func (a *APIResponse) GetError() *ServerErrorMessage {
	result := a.getBodyField("ServerErrorMessage")
	if result == nil {
		return nil
	}
	return result.(*ServerErrorMessage)
}

func (a *APIResponse) GetClientUserList() *ClientUserList {
	result := a.getBodyField("ClientUserList")
	if result == nil {
		return nil
	}
	return result.(*ClientUserList)
}

func (a *APIResponse) GetClientOpenChat() *ClientOpenChat {
	result := a.getBodyField("ClientOpenChat")
	if result == nil {
		return nil
	}
	return result.(*ClientOpenChat)
}

func (a *APIResponse) GetUser() *UserExtended {
	result := a.getBodyField("User")
	if result == nil {
		return nil
	}
	return result.(*UserExtended)
}

func (a *APIResponse) GetChatMessageReceived() *ChatMessageReceived {
	result := a.getBodyField("ChatMessageReceived")
	if result == nil {
		return nil
	}
	return result.(*ChatMessageReceived)
}

func (a *APIResponse) GetClientStartup() *ClientStartup {
	result := a.getBodyField("ClientStartup")
	if result == nil {
		return nil
	}
	return result.(*ClientStartup)
}

func (a *APIResponse) GetClientCommonSettings() *ClientCommonSettings {
	result := a.getBodyField("ClientCommonSettings")
	if result == nil {
		return nil
	}
	return result.(*ClientCommonSettings)
}

func (a *APIResponse) GetClientLoginSuccess() *ClientLoginSuccess {
	result := a.getBodyField("ClientLoginSuccess")
	if result == nil {
		return nil
	}
	return result.(*ClientLoginSuccess)
}

func (a *APIResponse) GetClientSessionChanged() *ClientSessionChanged {
	result := a.getBodyField("ClientSessionChanged")
	if result == nil {
		return nil
	}
	return result.(*ClientSessionChanged)
}

func (a *APIResponse) GetAppSettings() *AppSettings {
	result := a.getBodyField("AppSettings")
	if result == nil {
		return nil
	}
	return result.(*AppSettings)
}

func (a *APIResponse) GetCometConfiguration() *CometConfiguration {
	result := a.getBodyField("CometConfiguration")
	if result == nil {
		return nil
	}
	return result.(*CometConfiguration)
}

type EventMessage struct {
	GPB            string             `json:"$gpb"`
	MessageType    int                `json:"message_type"`
	Version        int                `json:"version"`
	MessageID      int                `json:"message_id"`
	ObjectType     int                `json:"object_type,omitempty"`
	Body           []EventMessageBody `json:"body"`
	ResponsesCount int                `json:"responses_count"`
}

type EventMessageWrapper struct {
	Cmd  string        `json:"cmd,omitempty"`
	Nd   int           `json:"nd,omitempty"`
	Seq  string        `json:"seq,omitempty"`
	Data []interface{} `json:"data,omitempty"`
}

type EventMessageBody struct {
	GPB           string         `json:"$gpb"`
	ChatIsWriting *ChatIsWriting `json:"chat_is_writing,omitempty"`
	ChatMessage   *ChatMessage   `json:"chat_message,omitempty"`
	PersonNotice  *PersonNotice  `json:"person_notice,omitempty"`
	MessageType   int            `json:"message_type"`
	ObjectType    int            `json:"object_type,omitempty"`
}

type ChatIsWriting struct {
	GPB          string `json:"$gpb"`
	WhoIsWriting string `json:"who_is_writing"`
	WhoIsWaiting string `json:"who_is_waiting"`
}

type ChatUserInfo struct {
	GPB            string `json:"$gpb"`
	Name           string `json:"name"`
	Age            int    `json:"age"`
	Gender         int    `json:"gender"`
	LargeImageID   string `json:"large_image_id"`
	InterestsCount int    `json:"interests_count"`
	Wish           string `json:"wish"`
	NumberOfPhotos int    `json:"number_of_photos"`
	ProfileRating  int    `json:"profile_rating"`
	LastSeenOnline string `json:"last_seen_online"`
}

type Sticker struct {
	GPB string `json:"$gpb"`
}

type PurchasedGift struct {
	GPB string `json:"$gpb"`
}

type VerificationAccessObject struct {
	GPB string `json:"$gpb"`
}

type VideoCallMsgInfo struct {
	GPB string `json:"$gpb"`
}

type Story struct {
	GPB string `json:"$gpb"`
}

type ExperimentalGift struct {
	GPB string `json:"$gpb"`
}

type PersonNotice struct {
	GPB                 string               `json:"$gpb"`
	Type                int                  `json:"type"`
	Folder              int                  `json:"folder"`
	DisplayValue        string               `json:"display_value"`
	IntValue            int                  `json:"int_value"`
	FilterDisplayValues []FilterDisplayValue `json:"filter_display_values"`
}

type FilterDisplayValue struct {
	GPB      string `json:"$gpb"`
	Filter   int    `json:"filter"`
	IntValue int    `json:"int_value"`
}

type ProfileFieldType string

const (
	ProfileFieldTypeLocation         ProfileFieldType = "location"
	ProfileFieldTypeAboutMe          ProfileFieldType = "aboutme_text"
	ProfileFieldTypeHeight           ProfileFieldType = "lifestyle_height"
	ProfileFieldTypeExercise         ProfileFieldType = "lifestyle_exercise"
	ProfileFieldTypeEducation        ProfileFieldType = "lifestyle_education"
	ProfileFieldTypeDrinking         ProfileFieldType = "lifestyle_drinking"
	ProfileFieldTypeSmoking          ProfileFieldType = "lifestyle_smoking"
	ProfileFieldTypeGender           ProfileFieldType = "lifestyle_gender"
	ProfileFieldTypeDatingIntentions ProfileFieldType = "lifestyle_dating_intentions"
	ProfileFieldTypeFamilyPlans      ProfileFieldType = "lifestyle_family_plans"
	ProfileFieldTypeZodiac           ProfileFieldType = "lifestyle_zodiak"
	ProfileFieldTypePolitics         ProfileFieldType = "lifestyle_politics"
)

func getProfileField(profileFields []ProfileField, fieldType ProfileFieldType) string {
	for _, f := range profileFields {
		if f.ID == string(fieldType) {
			return f.DisplayValue
		}
	}
	return ""
}

// App Startup

type ClientUserVerifiedGet struct {
	DisplayMessage     string                         `json:"display_message"`
	Methods            []UserVerificationMethodStatus `json:"methods"`
	VerificationStatus int                            `json:"verification_status"`
}

type UserVerificationMethodStatus struct {
	Type                        int                            `json:"type"`
	Name                        string                         `json:"name"`
	DisplayMessage              string                         `json:"display_message"`
	VerificationData            string                         `json:"verification_data,omitempty"`
	IsConnected                 bool                           `json:"is_connected"`
	IsConfirmed                 bool                           `json:"is_confirmed"`
	AllowUnlink                 bool                           `json:"allow_unlink"`
	PhoneNumberVerificationType int                            `json:"phone_number_verification_type,omitempty"`
	PhoneNumber                 string                         `json:"phone_number,omitempty"`
	PhonePrefix                 string                         `json:"phone_prefix,omitempty"`
	AccessRestrictions          VerificationAccessRestrictions `json:"access_restrictions"`
	AllowedSwitchTypes          []int                          `json:"allowed_switch_types,omitempty"`
	VerificationPromo           PromoBlock                     `json:"verification_promo,omitempty"`
	ExternalProviderData        ExternalProvider               `json:"external_provider_data,omitempty"`
	ExternalAccountURL          string                         `json:"external_account_url,omitempty"`
}

type VerificationAccessRestrictions struct {
	AvailableOptions []int  `json:"available_options,omitempty"`
	ChosenOption     int    `json:"chosen_option"`
	Disclaimer       string `json:"disclaimer,omitempty"`
	PeopleWithAccess int    `json:"people_with_access"`
	FeedConnected    bool   `json:"feed_connected,omitempty"`
}

type PromoBlock struct {
	Mssg               string           `json:"mssg"`
	Header             string           `json:"header"`
	PromoBlockType     int              `json:"promo_block_type"`
	PromoBlockPosition int              `json:"promo_block_position"`
	Buttons            []CallToAction   `json:"buttons"`
	ExtraTexts         []PromoBlockText `json:"extra_texts,omitempty"`
}

type CallToAction struct {
	Text   string `json:"text,omitempty"`
	Action int    `json:"action"`
	Type   int    `json:"type"`
}

type PromoBlockText struct {
	Type int    `json:"type"`
	Text string `json:"text"`
}

type ExternalProvider struct {
	ID              string                   `json:"id"`
	DisplayName     string                   `json:"display_name"`
	LogoURL         string                   `json:"logo_url"`
	Type            int                      `json:"type"`
	AuthData        ExternalProviderAuthData `json:"auth_data"`
	ReadPermissions []string                 `json:"read_permissions"`
}

type ExternalProviderAuthData struct {
	Type     int    `json:"type"`
	OauthURL string `json:"oauth_url"`
	AppID    string `json:"app_id"`
}

type ExtendedGender struct {
	// Additional properties can be added if needed
}

type ClientStartup struct {
	Host                 []string             `json:"host"`
	LocationServices     []int                `json:"location_services"`
	LanguageID           int                  `json:"language_id"`
	StartupType          int                  `json:"startup_type"`
	AppFeature           []ApplicationFeature `json:"app_feature"`
	AllowReview          bool                 `json:"allow_review"`
	FullSiteURL          string               `json:"full_site_url"`
	ForbidRegisterViaSMS bool                 `json:"forbid_register_via_sms"`
	PartnerID            int                  `json:"partner_id"`
	IsPushEnabled        bool                 `json:"is_push_enabled"`
	StartupSettings      StartupSettings      `json:"startup_settings"`
	PlatformID           int                  `json:"platform_id"`
}

type ApplicationFeature struct {
	Feature               int                   `json:"feature"`
	Enabled               bool                  `json:"enabled"`
	RequiredAction        int                   `json:"required_action,omitempty"`
	DisplayMessage        string                `json:"display_message,omitempty"`
	DisplayTitle          string                `json:"display_title,omitempty"`
	DisplayAction         string                `json:"display_action,omitempty"`
	GoalProgress          GoalProgress          `json:"goal_progress,omitempty"`
	AllowWebrtcCallConfig AllowWebrtcCallConfig `json:"allow_webrtc_call_config,omitempty"`
	ProductType           int                   `json:"product_type,omitempty"`
	PaymentAmount         int                   `json:"payment_amount,omitempty"`
}

type GoalProgress struct {
	Goal     int `json:"goal"`
	Progress int `json:"progress"`
}

type AllowWebrtcCallConfig struct {
	DialingTimeout        int  `json:"dialing_timeout"`
	BusyToneLength        int  `json:"busy_tone_length"`
	AllowTurnOffCamera    bool `json:"allow_turn_off_camera"`
	AllowMuteMic          bool `json:"allow_mute_mic"`
	MaxNonFatalRetries    int  `json:"max_non_fatal_retries"`
	StatePollingPeriodSec int  `json:"state_polling_period_sec"`
}

type StartupSettings struct {
	Providers []ExternalStatsProvider `json:"providers"`
}

type ExternalStatsProvider struct {
	Type            int    `json:"type"`
	URL             string `json:"url"`
	MaximumPoolSize int    `json:"maximumPoolSize"`
	MaximumTimeout  int    `json:"maximumTimeout"`
}

type ClientCommonSettings struct {
	PartnerID                               int                       `json:"partner_id"`
	LanguageID                              int                       `json:"language_id"`
	ApplicationFeatures                     []ApplicationFeature      `json:"application_features"`
	AllowReview                             bool                      `json:"allow_review"`
	ABTestingSettings                       ABTestingSettings         `json:"a_b_testing_settings"`
	ExternalEndpoints                       []ExternalEndpoint        `json:"external_endpoints"`
	DevFeatures                             []DevFeature              `json:"dev_features"`
	UserCountry                             Country                   `json:"user_country"`
	WebSpecificOptions                      WebSpecificOptions        `json:"web_specific_options"`
	PhoneVerificationForcePin               bool                      `json:"phone_verification_force_pin"`
	ServerGetSocialFriendsConnectionsPeriod int                       `json:"server_get_social_friends_connections_period"`
	AllowFacebookLikeOnReview               bool                      `json:"allow_facebook_like_on_review"`
	DefaultUnitType                         int                       `json:"default_unit_type"`
	LanguageISOCode                         string                    `json:"language_iso_code"`
	MaxRequestedUsersCount                  int                       `json:"max_requested_users_count"`
	RequiredChecks                          []int                     `json:"required_checks"`
	ReportLongProfileViewMilliseconds       []int                     `json:"report_long_profile_view_milliseconds"`
	EncountersQueueSettings                 EncountersQueueSettings   `json:"encounters_queue_settings"`
	ShowClientWhatsNew                      bool                      `json:"show_client_whats_new"`
	ShowServerWhatsNew                      bool                      `json:"show_server_whats_new"`
	FreezeConnectionsEnabled                bool                      `json:"freeze_connections_enabled"`
	CrushWhatsNewAllowed                    bool                      `json:"crush_whats_new_allowed"`
	AllowFullscreenPhotosInOtherProfiles    bool                      `json:"allow_fullscreen_photos_in_other_profiles"`
	SigninProviders                         []int                     `json:"signin_providers"`
	PhotoUploadSettings                     PhotoUploadSettings       `json:"photo_upload_settings"`
	WebPushInitParams                       WebPushInitParams         `json:"web_push_init_params"`
	SdkIntegrations                         []SdkIntegration          `json:"sdk_integrations"`
	VideoUploadSettings                     VideoUploadSettings       `json:"video_upload_settings"`
	OwnProfileSettings                      OwnProfileSettings        `json:"own_profile_settings"`
	ChatSettings                            GlobalChatSettings        `json:"chat_settings"`
	TooltipConfigs                          []TooltipConfig           `json:"tooltip_configs"`
	LivestreamSettings                      GlobalLivestreamSettings  `json:"livestream_settings"`
	PushOfferID                             string                    `json:"push_offer_id"`
	EmailDomains                            []string                  `json:"email_domains"`
	DateFormat                              DateFormat                `json:"date_format"`
	GroupChatsSettings                      ConversationsSettings     `json:"group_chats_settings"`
	IsBumbleEventModerator                  bool                      `json:"is_bumble_event_moderator"`
	LiveVideoSettings                       LiveVideoSettings         `json:"live_video_settings"`
	HideYourTurnBadgeAfterSec               int                       `json:"hide_your_turn_badge_after_sec"`
	MaxNumberOfLoginMethods                 int                       `json:"max_number_of_login_methods"`
	SnapCameraSettings                      []SnapCameraFeatureConfig `json:"snap_camera_settings"`
	LocationSettings                        LocationSettings          `json:"location_settings"`
	AutoplayVideoSettings                   VideoAutoplaySettings     `json:"autoplay_video_settings"`
	MatchScreenSettings                     MatchScreenSettings       `json:"match_screen_settings"`
}

type ABTestingSettings struct {
	Tests []ABTest `json:"tests"`
}

type ABTest struct {
	TestID      string `json:"test_id"`
	VariationID string `json:"variation_id"`
	SettingsID  string `json:"settings_id"`
}

type ExternalEndpoint struct {
	Type int    `json:"type"`
	URL  string `json:"url"`
}

type DevFeature struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

type WebSpecificOptions struct {
	OpenSearchFilter bool `json:"open_search_filter"`
	HighlightLexemes bool `json:"highlight_lexemes"`
}

type EncountersQueueSettings struct {
	Tooltips                          []Tooltip `json:"tooltips"`
	ShowVotingButtonsForNumberOfVotes int       `json:"show_voting_buttons_for_number_of_votes"`
	NumOfCompletedVotesToReport       int       `json:"num_of_completed_votes_to_report"`
	CacheExpirationTimeSec            int       `json:"cache_expiration_time_sec"`
}

type Tooltip struct {
	Type                    int            `json:"type"`
	Frequency               int            `json:"frequency,omitempty"`
	Text                    string         `json:"text,omitempty"`
	Buttons                 []CallToAction `json:"buttons,omitempty"`
	StatsOnly               bool           `json:"stats_only,omitempty"`
	PotentialMatchFrequency int            `json:"potential_match_frequency,omitempty"`
}

type PhotoUploadSettings struct {
	MinPhotoSize PhotoSize `json:"min_photo_size"`
	MaxSizeFast  PhotoSize `json:"max_size_fast"`
	MaxSizeSlow  PhotoSize `json:"max_size_slow"`
}

type PhotoSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type WebPushInitParams struct {
	ApplicationServerKey string         `json:"application_server_key"`
	WebServiceURL        string         `json:"web_service_url"`
	WebSitePushID        string         `json:"web_site_push_id"`
	Params               []GenericParam `json:"params"`
}

type GenericParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SdkIntegration struct {
	Type   int    `json:"type"`
	AppKey string `json:"app_key,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

type VideoUploadSettings struct {
	MinDuration             int         `json:"min_duration"`
	MaxDuration             int         `json:"max_duration"`
	MaxSizeBytes            int         `json:"max_size_bytes"`
	Format                  VideoFormat `json:"format"`
	MaxRecordingDurationSec int         `json:"max_recording_duration_sec"`
	AudioFormat             AudioFormat `json:"audio_format"`
}

type VideoFormat struct {
	Encoding              int       `json:"encoding"`
	MaxBitRateKbps        int       `json:"max_bit_rate_kbps"`
	MaxResolution         PhotoSize `json:"max_resolution"`
	MaxPortraitResolution PhotoSize `json:"max_portrait_resolution"`
}

type AudioFormat struct {
	Type         int  `json:"type"`
	SampleRateHz int  `json:"sample_rate_hz"`
	BitRateKbps  int  `json:"bit_rate_kbps"`
	AudioStereo  bool `json:"audio_stereo"`
	VbrEnable    bool `json:"vbr_enable"`
}

type OwnProfileSettings struct {
	LayoutElements []OwnProfileLayoutElement `json:"layout_elements"`
}

type OwnProfileLayoutElement struct {
	Type      int `json:"type"`
	Reference int `json:"reference"`
}

type GlobalChatSettings struct {
	MaxCharactersInMessage      int                    `json:"max_characters_in_message"`
	DelayBeforeWarnBadBlockerMs int                    `json:"delay_before_warn_bad_blocker_ms"`
	DelayBeforeShowGoodOpenerMs int                    `json:"delay_before_show_good_opener_ms"`
	AudioRecordingSettings      AudioRecordingSettings `json:"audio_recording_settings"`
}

type AudioRecordingSettings struct {
	StartRecordingDelayMs int         `json:"start_recording_delay_ms"`
	MaxRecordingLengthMs  int         `json:"max_recording_length_ms"`
	AudioFormat           AudioFormat `json:"audio_format"`
	WaveformLength        int         `json:"waveform_length"`
}

type TooltipConfig struct {
	Context  int       `json:"context"`
	Tooltips []Tooltip `json:"tooltips,omitempty"`
}

type GlobalLivestreamSettings struct {
	MaxAllowedMessageAge int   `json:"max_allowed_message_age"`
	ConnectionTimeoutSec int   `json:"connection_timeout_sec"`
	MaxMessageSize       int   `json:"max_message_size"`
	SuppressEvents       []int `json:"suppress_events"`
}

type DateFormat struct {
	Day   DateElement `json:"day"`
	Month DateElement `json:"month"`
	Year  DateElement `json:"year"`
}

type DateElement struct {
	Position    int    `json:"position"`
	Placeholder string `json:"placeholder"`
	Title       string `json:"title"`
	Length      int    `json:"length"`
	A11yText    string `json:"a11y_text"`
}

type ConversationsSettings struct {
	ConversationsPollingPeriodSec int `json:"conversations_polling_period_sec"`
	MaxNumOfParticipants          int `json:"max_num_of_participants"`
	MaxGroupNameLength            int `json:"max_group_name_length"`
}

type LiveVideoSettings struct {
	Playback         LiveVideoPlaybackSettings `json:"playback"`
	StatsIntervalSec int                       `json:"stats_interval_sec"`
	MaxStatsDelayMs  int                       `json:"max_stats_delay_ms"`
}

type LiveVideoPlaybackSettings struct {
	IsDefaultMute       bool `json:"is_default_mute"`
	ForwardTimeSec      int  `json:"forward_time_sec"`
	BackwardTimeSec     int  `json:"backward_time_sec"`
	CanChangeResolution bool `json:"can_change_resolution"`
	Autoplay            bool `json:"autoplay"`
}

type SnapCameraFeatureConfig struct {
	Context  int      `json:"context"`
	GroupIDs []string `json:"group_ids"`
}

type LocationSettings struct {
	SignificantDistanceMeters int `json:"significant_distance_meters"`
	MinUpdateTimeIntervalSec  int `json:"min_update_time_interval_sec"`
	MinUpdateDistanceMeters   int `json:"min_update_distance_meters"`
}

type VideoAutoplaySettings struct {
	ProfileVideoOptions []int `json:"profile_video_options"`
}

type MatchScreenSettings struct {
	BffMatchExpirationTimeSec int `json:"bff_match_expiration_time_sec"`
}

type ClientLoginSuccess struct {
	SessionID                    string   `json:"session_id"`
	DeprecatedNewPeople          int      `json:"deprecated_new_people"`
	DeprecatedNewMessages        int      `json:"deprecated_new_messages"`
	IsFirstLogin                 bool     `json:"is_first_login"`
	Host                         []string `json:"host"`
	AllowReview                  bool     `json:"allow_review"`
	UserInfo                     User     `json:"user_info"`
	EncryptedUserID              string   `json:"encrypted_user_id"`
	ExternalProviderTokenRefresh int      `json:"external_provider_token_refresh"`
}

type ClientSessionChanged struct {
	NewSessionID          string `json:"new_session_id"`
	IsRegistrationSession bool   `json:"is_registration_session"`
}

type AppSettings struct {
	InterfaceLanguage                   int                   `json:"interface_language"`
	InterfaceSound                      bool                  `json:"interface_sound"`
	InterfaceMetric                     bool                  `json:"interface_metric"`
	NotifyMessages                      bool                  `json:"notify_messages"`
	NotifyWantYou                       bool                  `json:"notify_want_you"`
	NotifyMutual                        bool                  `json:"notify_mutual"`
	PrivacyShowDistance                 bool                  `json:"privacy_show_distance"`
	NotifyAlerts                        bool                  `json:"notify_alerts"`
	PrivacyShowOnlineStatus             bool                  `json:"privacy_show_online_status"`
	EmailMessages                       bool                  `json:"email_messages"`
	EmailAlerts                         bool                  `json:"email_alerts"`
	EmailNews                           bool                  `json:"email_news"`
	EmailMatches                        bool                  `json:"email_matches"`
	NotifyPhotoratings                  bool                  `json:"notify_photoratings"`
	EmailPhotoratings                   bool                  `json:"email_photoratings"`
	VerificationHideDetails             bool                  `json:"verification_hide_details"`
	VerificationOnlyVerifiedMessages    bool                  `json:"verification_only_verified_messages"`
	PrivacyShowInPublicSearch           bool                  `json:"privacy_show_in_public_search"`
	HideAccount                         bool                  `json:"hide_account"`
	LetOtherUsersShareMyProfile         bool                  `json:"let_other_users_share_my_profile"`
	PrivacyShowOnlyToPeopleILikeOrVisit bool                  `json:"privacy_show_only_to_people_i_like_or_visit"`
	SkipNonModeratedAlert               bool                  `json:"skip_non_moderated_alert"`
	EmailWantYou                        bool                  `json:"email_want_you"`
	ContactNotInterestAction            int                   `json:"contact_not_interest_action"`
	PrivacyAllowSearchByEmail           bool                  `json:"privacy_allow_search_by_email"`
	Menu                                AppSettingsMenu       `json:"menu"`
	PrivacyEnableTargetedAds            bool                  `json:"privacy_enable_targeted_ads"`
	EnabledGameModes                    []EnabledGameMode     `json:"enabled_game_modes"`
	PrivacyAdsState                     int                   `json:"privacy_ads_state"`
	NotificationsMenu                   AppSettingsMenu       `json:"notifications_menu"`
	InvisibleMode                       AppSetting            `json:"invisible_mode"`
	PrivacyShowGender                   bool                  `json:"privacy_show_gender"`
	BillingEmail                        string                `json:"billing_email"`
	AutoplayVideoSettings               VideoAutoplaySettings `json:"autoplay_video_settings"`
	Context                             int                   `json:"context"`
}

type AppSettingsMenu struct {
	Sections []AppSettingsMenuSection `json:"sections"`
	Title    string                   `json:"title,omitempty"`
}

type AppSettingsMenuSection struct {
	Items []AppSettingsMenuItem `json:"items"`
	Name  string                `json:"name,omitempty"`
	Text  string                `json:"text,omitempty"`
}

type AppSettingsMenuItem struct {
	Type                int                 `json:"type"`
	Name                string              `json:"name"`
	NotificationSetting NotificationSetting `json:"notification_setting,omitempty"`
}

type NotificationSetting struct {
	Type          int    `json:"type"`
	SendEmail     bool   `json:"send_email,omitempty"`
	SendCloudPush bool   `json:"send_cloud_push,omitempty"`
	SendInapp     bool   `json:"send_inapp,omitempty"`
	EmailApproved bool   `json:"email_approved,omitempty"`
	Description   string `json:"description,omitempty"`
	Category      string `json:"category,omitempty"`
	StatsID       int    `json:"stats_id,omitempty"`
	SendSMS       bool   `json:"send_sms,omitempty"`
}

type EnabledGameMode struct {
	GameMode              int  `json:"game_mode"`
	Available             bool `json:"available"`
	CanToggleAvailability bool `json:"can_toggle_availability"`
}

type AppSetting struct {
	State           int              `json:"state"`
	TogglingOptions []TogglingOption `json:"toggling_options"`
	TogglingReasons []TogglingReason `json:"toggling_reasons"`
}

type TogglingOption struct {
	Header          string `json:"header"`
	TogglePeriodSec int    `json:"toggle_period_sec"`
	// PreferredAlternative TogglingOption `json:"preferred_alternative,omitempty"`
	Message    string `json:"message,omitempty"`
	AcceptText string `json:"accept_text,omitempty"`
	RefuseText string `json:"refuse_text,omitempty"`
}

type TogglingReason struct {
	Text string `json:"text"`
	ID   string `json:"id"`
}

type CometConfiguration struct {
	Path     string `json:"path"`
	Sequence string `json:"sequence"`
}

type CookieContainer struct {
	SessionCookieName string `json:"session_cookie_name"`
	DeviceID          string `json:"device_id"`
	BuzzLangCode      string `json:"buzz_lang_code"`
	FirstWebVisit     string `json:"first_web_visit"`
	DnsDisplayed      string `json:"dnsDisplayed"`
	CcpaApplies       string `json:"ccpaApplies"`
	SignedLspa        string `json:"signedLspa"`
	SpSu              string `json:"_sp_su"`
	GclAu             string `json:"_gcl_au"`
	Ga                string `json:"_ga"`
	CcpaUUID          string `json:"ccpaUUID"`
	HDRXUserID        string `json:"HDR-X-User-id"`
	Session           string `json:"session"`
	Aid               string `json:"aid,omitempty"`
}

type BadooMessage struct {
	GPB            string        `json:"$gpb"`
	MessageType    int           `json:"message_type"`
	Version        int           `json:"version"`
	MessageID      int           `json:"message_id"`
	ObjectType     int           `json:"object_type"`
	Body           []MessageBody `json:"body"`
	ResponsesCount int           `json:"responses_count"`
	IsBackground   bool          `json:"is_background"`
	Vhost          string        `json:"vhost"`
}

type ServerAppStartup struct {
	AppBuild                          string                             `json:"app_build"`
	AppName                           string                             `json:"app_name"`
	AppVersion                        string                             `json:"app_version"`
	CanSendSMS                        bool                               `json:"can_send_sms"`
	UserAgent                         string                             `json:"user_agent"`
	ScreenWidth                       int                                `json:"screen_width"`
	ScreenHeight                      int                                `json:"screen_height"`
	Language                          int                                `json:"language"`
	IsColdStart                       bool                               `json:"is_cold_start"`
	ExternalProviderRedirectURL       string                             `json:"external_provider_redirect_url"`
	Locale                            string                             `json:"locale"`
	SystemLocale                      string                             `json:"system_locale"`
	AppPlatformType                   int                                `json:"app_platform_type"`
	AppProductType                    int                                `json:"app_product_type"`
	DeviceInfo                        *DeviceInfo                        `json:"device_info"`
	BuildConfiguration                int                                `json:"build_configuration"`
	BuildFingerprint                  string                             `json:"build_fingerprint"`
	SupportedFeatures                 []int                              `json:"supported_features"`
	SupportedMinorFeatures            []int                              `json:"supported_minor_features"`
	SupportedNotifications            []int                              `json:"supported_notifications"`
	SupportedPaymentProviders         []int                              `json:"supported_payment_providers"`
	SupportedPromoBlocks              []*SupportedPromoBlock             `json:"supported_promo_blocks"`
	SupportedUserSubstitutes          []*SupportedUserSubstitute         `json:"supported_user_substitutes"`
	SupportedOnboardingTypes          []int                              `json:"supported_onboarding_types"`
	UserFieldFilterClientLoginSuccess *UserFieldFilterClientLoginSuccess `json:"user_field_filter_client_login_success"`
	ABTestingSettings                 *ABTestingSettings                 `json:"a_b_testing_settings"`
	DevFeatures                       []string                           `json:"dev_features"`
	DeviceID                          string                             `json:"device_id"`
	SupportedScreens                  []*SupportedScreen                 `json:"supported_screens"`
	SupportedLandings                 []*SupportedLanding                `json:"supported_landings"`
	AppDomain                         string                             `json:"app_domain"`
}

type DeviceInfo struct {
	WebcamAvailable bool `json:"webcam_available"`
	FormFactor      int  `json:"form_factor"`
}

type SupportedPromoBlock struct {
	Context  int   `json:"context"`
	Position int   `json:"position"`
	Types    []int `json:"types"`
}

type SupportedUserSubstitute struct {
	Context int   `json:"context"`
	Types   []int `json:"types"`
}

type UserFieldFilterClientLoginSuccess struct {
	Projection []int `json:"projection"`
}

type SupportedScreen struct {
	Type    int `json:"type"`
	Version int `json:"version"`
}

type SupportedLanding struct {
	Source              int   `json:"source"`
	Params              []int `json:"params"`
	SearchSettingsTypes []int `json:"search_settings_types"`
}
