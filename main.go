package main

import (
	"fmt"
	"log"
	"os"

	hass "github.com/exceptionfault/hass-golang-ws/client"
)

func main() {
	api_token := os.Getenv("HASS_API_TOKEN")
	host := os.Getenv("HASS_URL")

	client := hass.CreateHassClient(host, api_token).WithEncryption(false)
	err := client.Connect()
	if err != nil {
		panic(err)
	}
	defer client.Disconnect()

	// -----------------------------------------------------------------------------------
	// Get a list of services
	client.GetServices(func(services []*hass.Service, err error) {
		if err != nil {
			panic(err)
		}
		for _, s := range services {
			fmt.Println(s.Domain, s.Id)
			for name, field := range s.Fields {
				fmt.Printf("   %s: %s\n", name, field.Description)
			}
		}
	})

	// -----------------------------------------------------------------------------------
	// Call a service with parameters
	params := &map[string]string{
		"name":    "hass-golang-ws",
		"message": "sent this log",
		"domain":  "text",
	}
	err = client.CallService("logbook", "log", params, func(sr hass.ServerResult) {
		fmt.Println("Result of calling service logbook.log:", sr)
	})
	if err != nil {
		fmt.Println("Error occured during call of service logbook.log:", err)
	}
	fmt.Println("Called service logbook.log")

	// -----------------------------------------------------------------------------------
	// Subscribe to events
	// client.SubscribeEvent(hass.EVT_STATE_CHANGED, func(evt hass.Event) {
	// 	if evt.IsStateChangedEvent() {
	// 		svt, err := evt.GetStateChangedEvent()
	// 		if err != nil {
	// 			log.Printf("Got error: %v", err)
	// 			return
	// 		}
	// 		log.Printf("Entity %s changed from %s to %s", svt.EntityId, svt.OldState.State, svt.NewState.State)
	// 		return
	// 	}
	// 	log.Printf("Got Event %s at %s: %v", evt.EventType, evt.TimeFired, evt.Data)
	// })

	interrupt := make(chan os.Signal, 1)
	sig := <-interrupt
	log.Println(sig)
}
