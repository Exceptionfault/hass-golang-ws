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
	Success bool                   `json:"success"`
	Result  map[string]interface{} `json:"result"`
	Event   Event                  `json:"event"`
}

type subscribeEventMessage struct {
	typedIdMessage
	EventType string `json:"event_type"`
}
