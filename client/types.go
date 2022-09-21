package client

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Context struct {
	Id       string  `json:"id"`
	ParentId *string `json:"parent_id"`
	UserId   string  `json:"user_id"`
}

type Event struct {
	EventType EventType              `json:"event_type"`
	TimeFired time.Time              `json:"time_fired"`
	Origin    string                 `json:"origin"`
	Data      map[string]interface{} `json:"data"`
	Context   Context                `json:"context"`
}

type Service struct {
	Domain      string
	Name        string `json:"name"`
	Id          string
	Description string `json:"description"`
	Fields      map[string]Field
}

type Field struct {
	Description string                 `json:"description"`
	Example     interface{}            `json:"example"`
	Required    bool                   `json:"required"`
	Selecttor   map[string]interface{} `json:"selector"`
}

/* Returns true, if the event is of type `state_changed` In  this case you can safely get the parsed event via
```
if event.IsStateChangedEvent() {
	changeEvent, err := event.GetStateChangedEvent()
}
```
*/
func (evt *Event) IsStateChangedEvent() bool {
	return evt.EventType == EVT_STATE_CHANGED
}

func (evt *Event) GetStateChangedEvent() (*StateChangedEvent, error) {
	if !evt.IsStateChangedEvent() {
		return nil, fmt.Errorf("Event must be of type %s", EVT_STATE_CHANGED)
	}
	e := &StateChangedEvent{}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "json", Result: e})
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(evt.Data); err != nil {
		return nil, fmt.Errorf("Cannot decode StateChangedEvent: %w", err)
	}
	return e, nil
}

type StateChangedEvent struct {
	EntityId string `json:"entity_id"`
	OldState State  `json:"old_state"`
	NewState State  `json:"new_state"`
}

type State struct {
	State       string                 `json:"state"`        // String representation of the current state of the entity. Example off.
	EntityId    string                 `json:"entity_id"`    // Entity ID. Format: <domain>.<object_id>. Example: light.kitchen.
	ObjectId    string                 `json:"object_id"`    // Object ID of entity. Example: kitchen.
	Domain      string                 `json:"domain"`       // Domain of the entity. Example: light.
	Name        string                 `json:"name"`         // Name of the entity. Based on friendly_name attribute with fall back to object ID. Example: Kitchen Ceiling.
	LastChanged string                 `json:"last_changed"` // Time the state changed in the state machine in UTC time. This is not updated when there are only updated attributes. Example: 2017-10-28 08:13:36.715874+00:00.
	LastUpdated string                 `json:"last_updated"` // Time the state was written to the state machine in UTC time. Note that writing the exact same state including attributes will not result in this field being updated. Example: 2017-10-28 08:13:36.715874+00:00.
	Attributes  map[string]interface{} `json:"attributes"`   // A dictionary with extra attributes related to the current state.
	Context     Context                `json:"context"`      // A dictionary with extra attributes related to the context of the state.
}
