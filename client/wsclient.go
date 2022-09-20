package client

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type HassClient struct {
	authenticated bool
	address       string
	token         string
	connection    *websocket.Conn
	id            *idGen // go-routine safe id generation for commands

	eventSubscriptions map[uint]func(Event)
}

// Create a new client with the given target address and authentication token. Use `client.connect()` afterwards to establish a connection.
func CreateHassClient(address string, token string) *HassClient {
	return &HassClient{
		address:            address,
		token:              token,
		id:                 &idGen{},
		eventSubscriptions: make(map[uint]func(Event)),
	}
}

// Connects to the homeassistant instance and initiates authentication via token.
// If this method returns without error, you are ready to execute commands to send or subscribe to events.
//
func (client *HassClient) Connect() error {
	if client.authenticated {
		return fmt.Errorf("already connected")
	}

	// TODO: add support for wss (secure websocket connections)
	u := url.URL{Scheme: "ws", Host: client.address, Path: API_ENDPOINT}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	client.connection = c

	if err = client.authenticate(); err != nil {
		return err
	}

	go client.listen()
	return nil
}

func (client *HassClient) Disconnect() {
	log.Printf("Closing connection to %s\n", client.connection.RemoteAddr().String())
	err := client.connection.Close()
	if err != nil {
		log.Printf("Error closing connection %v\n", err)
	}
}

func (client *HassClient) listen() {
	for {
		var msg serverMessage
		err := client.connection.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			return
		}

		if msg.MessageType == string(msg_EVENT) {
			if callback, ok := client.eventSubscriptions[msg.Id]; ok {
				callback(msg.Event)
				continue
			}
		}

		log.Printf("recv: %v\n", msg)
	}
}

// authenticate sends an `auth` message to the server to authenticate via token.
func (client *HassClient) authenticate() error {

	var smsg serverAuthMessage
	// first, listen for servers `auth_required` message
	err := client.connection.ReadJSON(&smsg)
	log.Printf("Got auth_request from homeassistant server %s version %s", client.connection.RemoteAddr().String(), smsg.HassVersion)
	if err != nil {
		return fmt.Errorf("error when receiving server auth request: %w", err)
	}
	if smsg.MessageType != string(msg_AUTH_REQUEST) {
		return fmt.Errorf("expected server to send auth_request, but got: %s", smsg.MessageType)
	}

	// second, respond with `auth` message
	err = client.send(clientAuthMessage{Token: client.token, MessageType: string(msg_AUTH_RESPONSE)})
	log.Printf("Sent client authentication token")
	if err != nil {
		return fmt.Errorf("Error sending authentication token to server: %w", err)
	}

	// last, get servers auth_ok or auth_invalid response
	err = client.connection.ReadJSON(&smsg)
	if err != nil {
		return fmt.Errorf("error when receiving servers auth response: %w", err)
	}

	if smsg.MessageType == string(msg_AUTH_OK) {
		log.Printf("Authentication successful")
		client.authenticated = true
		return nil
	}
	if smsg.MessageType == string(msg_AUTH_INVALID) {
		log.Printf("authentication failed: %s", smsg.Message)
		return fmt.Errorf("authentication failed: %s", smsg.Message)
	}

	// we should not get here
	return fmt.Errorf("received unexpected server resposne: %s", smsg.MessageType)
}

func (client *HassClient) send(message interface{}) error {
	return client.connection.WriteJSON(message)
}

func (client *HassClient) SubscribeEvent(eventType EventType, callback func(Event)) {
	msg := client.createSubscriveEventMessage(eventType)
	client.eventSubscriptions[msg.Id] = callback
	if err := client.send(msg); err != nil {
		log.Printf("error creating subscription: %v", err)
	}
	log.Printf("created subscription for event type %s with id %d", eventType, msg.Id)
}

func (client *HassClient) createSubscriveEventMessage(eventType EventType) *subscribeEventMessage {
	return &subscribeEventMessage{
		EventType: string(eventType),
		typedIdMessage: typedIdMessage{
			Id:          client.id.inc(),
			MessageType: string(msg_SUBSCRIBE_EVENT),
		},
	}
}
