# Golang Websocket Client for Homeassistant

This library can be used to commiunicate with [Homeassistant](https://homeassistant.io) via websocket API to send and receive events.

__:warning: This client is in a very early stage of development. The API will change!__

__:warning: This client does not yet support all features of the websocket API nor is it fully tested!__


## Installation

```go
go get github.com/exceptionfault/hass-golang-ws
```


## Usage

Create a client instance and connect:

```go
import (
    hass "github.com/exceptionfault/hass-golang-ws/client"
)

func main() {
    host := "homeassistant.local:8123"
    api_token := "DHJ)DFIK...." // get a long-living token from the UI

    client := hass.CreateHassClient(host, api_token)
    err := client.Connect()
    if err != nil {
        panic(err)
    }
    defer client.Disconnect()
}
```

### Subscribe to Events with Callback

To subscribe to events there are two options which differentiate in the way you get
notified about new events. First method allows you to provide a callback function
which get's called on every received event.

To specify a filter which events you are looking for, you can use predefined constants like `hass.EVT_STATE_CHANGED`.

```go
client.SubscribeEvent(hass.EVT_STATE_CHANGED, 
    // This function get's called on every event
    func(evt hass.Event) {
		log.Printf("Got Event %s at %s: %v", evt.EventType, evt.TimeFired, evt.Data)
    }
)
```

### Get a list of all Services

You can get a list of all Services. The example below prints all services by domain and technical name, including all parameters of the Service.

```go
client.GetServices(func(services []*hass.Service, err error) {
    if err != nil {
        panic(err)
    }
    for _, s := range services {
        fmt.Println(s.Domain, s.Id)
        for name, field := range s.Fields {
            fmt.Printf("    %s: %s\n", name, field.Description)
        }
    }
})

// Example Output:
// light turn_off
//     transition: Duration it takes to get to next state.
//     flash: If the light should flash.
// light toggle
//     transition: Duration it takes to get to next state.
//     rgb_color: Color for the light in RGB-format.
//     color_name: A human readable color name.
//     effect: Light effect.
//     xy_color: Color for the light in XY-format.
//     color_temp: Color temperature for the light in mireds.
//     brightness: Number indicating brightness, where 0 turns ...
```


## Running the sample

The `main.go` contains an example how to connect and do certain actions with the client.
To run the samples make sure to set the following environment variables:
```
export HASS_API_TOKEN=<your-long-living-token>
export HASS_URL=<address-of-homeassistant:8123>
```

Then run `go run main.go`
