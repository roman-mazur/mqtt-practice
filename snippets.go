package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:1883", *host))
	client := mqtt.NewClient(opts)

	client.Subscribe("/test/inception", 0, func(client mqtt.Client, message mqtt.Message) {
		var data payload
		if err := json.Unmarshal(message.Payload(), &data); err != nil {
			return
		}
	})

	data, _ := json.Marshal(&payload{Wisdom: "hello!", Secret: secret, Team: "teacher"})
	res := client.Publish("/test/result", 0, false, data)
}
