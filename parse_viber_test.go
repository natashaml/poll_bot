package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const message = `{
	"event":"message",
	"timestamp":1457764197627,
	"message_token":4912661846655238145,
	"sender":{
		 "id":"01234567890A=",
		 "name":"John McClane",
		 "avatar":"http://avatar.example.com",
		 "country":"UK",
		 "language":"en",
		 "api_version":1
	},
	"message":{
		 "type":"text",
		 "text":"a message to the service",
		 "media":"http://example.com",
		 "location":{
				"lat":50.76891,
				"lon":6.11499
		 },
		 "tracking_data":"tracking data"
	}
}`

func TestParseViberMessage(t *testing.T) {
	m, err := parseCallback([]byte(message))
	require.NoError(t, err)
	require.Equal(t, m.Event, "message")
	require.Equal(t, m.Message.Type, "text")
	require.Equal(t, m.Message.Text, "a message to the service")
	require.Equal(t, m.User.Id, "01234567890A=")
}

const subscribeMessage = `{
	"event":"subscribed",
	"timestamp":1457764197627,
	"user":{
		 "id":"01234567890A=",
		 "name":"John McClane",
		 "avatar":"http://avatar.example.com",
		 "country":"UK",
		 "language":"en",
		 "api_version":1
	},
	"message_token":4912661846655238145
}`

func TestParseViberSubscribe(t *testing.T) {
	m, err := parseCallback([]byte(subscribeMessage))
	require.NoError(t, err)
	require.Equal(t, m.Event, "subscribed")
	require.Equal(t, m.User.Id, "01234567890A=")
}
