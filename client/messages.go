package client

/***** Authentication Phase *****/

type serverAuthMessage struct {
	MessageType string `json:"type"`
	HassVersion string `json:"ha_version"`
	Message     string `json:"message"`
}

type clientAuthMessage struct {
	MessageType string `json:"type"`
	Token       string `json:"access_token"`
}

/***** Command Phase *****/

type typedIdMessage struct {
	Id          uint   `json:"id"`
	MessageType string `json:"type"`
}

type serverMessage struct {
	typedIdMessage
	ServerResult
	Event Event `json:"event"`
}

type subscribeEventMessage struct {
	typedIdMessage
	EventType string `json:"event_type"`
}

type entityTarget struct {
	EntityId string `json:"entity_id,omitempty"`
}

type callServiceMessage struct {
	typedIdMessage
	Domain      string             `json:"domain"`
	Service     string             `json:"service"`
	ServiceData *map[string]string `json:"service_data,omitempty"`
	Target      *entityTarget      `json:"target,omitempty"`
}
