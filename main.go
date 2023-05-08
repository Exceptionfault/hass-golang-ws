package main

import (
	"fmt"
	"log"
	"os"

	hass "github.com/bemble/hass-golang-ws/client"
)

func main() {
	api_token := os.Getenv("HASS_API_TOKEN")
	host := os.Getenv("HASS_URL")
	scheme := os.Getenv("HASS_SCHEME")

	client := hass.CreateHassClient(host, api_token, scheme)
	err := client.Connect()
	if err != nil {
		panic(err)
	}
	defer client.Disconnect()

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
