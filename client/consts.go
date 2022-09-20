package client

// Default Endpoint of Websocket
const API_ENDPOINT string = "/api/websocket"

// Custom type to express available message types
type messageType string

// Custom type to express available event types
type EventType string

// List of all supported message types
const (
	msg_AUTH_REQUEST  messageType = "auth_required"
	msg_AUTH_RESPONSE messageType = "auth"
	msg_AUTH_OK       messageType = "auth_ok"
	msg_AUTH_INVALID  messageType = "auth_invalid"

	msg_RESULT messageType = "result"
	msg_EVENT  messageType = "event"

	msg_SUBSCRIBE_EVENT messageType = "subscribe_events"
)

const (
	EVT_STATE_CHANGED EventType = "state_changed" // Event type: "state_changed"
)
