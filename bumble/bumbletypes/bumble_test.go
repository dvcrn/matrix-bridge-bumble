package bumbletypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIResponseGetBodyField(t *testing.T) {
	tests := []struct {
		name     string
		response APIResponse
		field    string
		expected interface{}
	}{
		{
			name: "Empty body",
			response: APIResponse{
				Body: []MessageBody{},
			},
			field:    "ClientUserList",
			expected: nil,
		},
		{
			name: "Field not found",
			response: APIResponse{
				Body: []MessageBody{
					{ClientUserList: &ClientUserList{}},
				},
			},
			field:    "NonExistentField",
			expected: nil,
		},
		{
			name: "Field found but nil",
			response: APIResponse{
				Body: []MessageBody{
					{ClientUserList: nil},
				},
			},
			field:    "ClientUserList",
			expected: nil,
		},
		{
			name: "Field found and not nil",
			response: APIResponse{
				Body: []MessageBody{
					{ClientUserList: &ClientUserList{}, ClientOpenChat: &ClientOpenChat{}},
				},
			},
			field:    "ClientOpenChat",
			expected: &ClientOpenChat{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.getBodyField(tt.field)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAPIResponseGetters(t *testing.T) {
	response := APIResponse{
		Body: []MessageBody{
			{
				ServerErrorMessage:   &ServerErrorMessage{},
				ClientUserList:       &ClientUserList{},
				ClientOpenChat:       &ClientOpenChat{},
				User:                 &UserExtended{},
				ChatMessageReceived:  &ChatMessageReceived{},
				ClientStartup:        &ClientStartup{},
				ClientCommonSettings: &ClientCommonSettings{},
				ClientLoginSuccess:   &ClientLoginSuccess{},
				ClientSessionChanged: &ClientSessionChanged{},
				AppSettings:          &AppSettings{},
				CometConfiguration:   &CometConfiguration{},
			},
		},
	}

	t.Run("GetError", func(t *testing.T) {
		result := response.GetError()
		assert.NotNil(t, result)
		assert.IsType(t, &ServerErrorMessage{}, result)
	})

	t.Run("GetClientUserList", func(t *testing.T) {
		result := response.GetClientUserList()
		assert.NotNil(t, result)
		assert.IsType(t, &ClientUserList{}, result)
	})

	t.Run("GetClientOpenChat", func(t *testing.T) {
		result := response.GetClientOpenChat()
		assert.NotNil(t, result)
		assert.IsType(t, &ClientOpenChat{}, result)
	})

	t.Run("GetUser", func(t *testing.T) {
		result := response.GetUser()
		assert.NotNil(t, result)
		assert.IsType(t, &UserExtended{}, result)
	})

	t.Run("GetChatMessageReceived", func(t *testing.T) {
		result := response.GetChatMessageReceived()
		assert.NotNil(t, result)
		assert.IsType(t, &ChatMessageReceived{}, result)
	})

	t.Run("GetClientStartup", func(t *testing.T) {
		result := response.GetClientStartup()
		assert.NotNil(t, result)
		assert.IsType(t, &ClientStartup{}, result)
	})

	t.Run("GetClientCommonSettings", func(t *testing.T) {
		result := response.GetClientCommonSettings()
		assert.NotNil(t, result)
		assert.IsType(t, &ClientCommonSettings{}, result)
	})

	t.Run("GetClientLoginSuccess", func(t *testing.T) {
		result := response.GetClientLoginSuccess()
		assert.NotNil(t, result)
		assert.IsType(t, &ClientLoginSuccess{}, result)
	})

	t.Run("GetClientSessionChanged", func(t *testing.T) {
		result := response.GetClientSessionChanged()
		assert.NotNil(t, result)
		assert.IsType(t, &ClientSessionChanged{}, result)
	})

	t.Run("GetAppSettings", func(t *testing.T) {
		result := response.GetAppSettings()
		assert.NotNil(t, result)
		assert.IsType(t, &AppSettings{}, result)
	})

	t.Run("GetCometConfiguration", func(t *testing.T) {
		result := response.GetCometConfiguration()
		assert.NotNil(t, result)
		assert.IsType(t, &CometConfiguration{}, result)
	})
}

func TestAPIResponseMultipleBodyEntries(t *testing.T) {
	response := APIResponse{
		Body: []MessageBody{
			{ClientUserList: &ClientUserList{TotalCount: 2}},
			{ClientUserList: &ClientUserList{TotalCount: 3}},
			{ClientOpenChat: &ClientOpenChat{Title: "chat1"}},
		},
	}

	t.Run("GetClientUserList returns first non-nil entry", func(t *testing.T) {
		result := response.GetClientUserList()
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.TotalCount)
	})

	t.Run("GetClientOpenChat returns correct entry", func(t *testing.T) {
		result := response.GetClientOpenChat()
		assert.NotNil(t, result)
		assert.Equal(t, "chat1", result.Title)
	})
}

func TestAPIResponseEmptyBody(t *testing.T) {
	response := APIResponse{
		Body: []MessageBody{},
	}

	t.Run("All getters return nil for empty body", func(t *testing.T) {
		assert.Nil(t, response.GetError())
		assert.Nil(t, response.GetClientUserList())
		assert.Nil(t, response.GetClientOpenChat())
		assert.Nil(t, response.GetUser())
		assert.Nil(t, response.GetChatMessageReceived())
		assert.Nil(t, response.GetClientStartup())
		assert.Nil(t, response.GetClientCommonSettings())
		assert.Nil(t, response.GetClientLoginSuccess())
		assert.Nil(t, response.GetClientSessionChanged())
		assert.Nil(t, response.GetAppSettings())
		assert.Nil(t, response.GetCometConfiguration())
	})
}
