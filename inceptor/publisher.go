package main

import (
	"crypto/sha1"
	"encoding/hex"

	"encoding/json"
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

var (
	host = flag.String("host", "a7cf5da.local", "broker host")
)

var wisdom = []string{
	"1: we",
	"2: finally",
	"3: did",
	"4: it",
	"but",
	"also",
	"you",
	"can",
	"give up",
	"keep",
	"it",
	"simple",
}

type payload struct {
	Time time.Time `json:"time"`
	Wisdom string `json:"wisdom"`
	Secret string `json:"secret"`
	Team string `json:"team"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	secret := fmt.Sprintf("%d-%d", rand.Intn(1000), rand.Intn(1000))

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:1883", *host))
	client := mqtt.NewClient(opts)

	log.Printf("Connecting to %v", opts)
	token := client.Connect()
	if !token.WaitTimeout(5 * time.Second) {
		log.Fatalf("Error on connection: %s", token.Error())
	}

	stopSignal := make(chan os.Signal)
	signal.Notify(stopSignal, os.Interrupt)
	go func() {
		<-stopSignal
		log.Printf("Disconnecting...")
		client.Disconnect(uint(time.Second))
		log.Printf("Bye")
		os.Exit(0)
	}()

	go func() {
		client.Subscribe("/test/result", 0, func(client mqtt.Client, message mqtt.Message) {
			var data payload
			if err := json.Unmarshal(message.Payload(), &data); err != nil {
				return
			}
			validSecret := hex.EncodeToString(sha1.New().Sum([]byte(secret + data.Team)))
			if data.Secret == validSecret {
				log.Printf("Received wisdom from %s: %s", data.Wisdom, data.Team)
			} else {
				log.Printf("Bad secret from %s", data.Team)
			}
		})
	}()

	fmt.Println(hex.EncodeToString(sha1.New().Sum([]byte(secret + "teacher"))))

	log.Printf("Start sending data")
	for range time.Tick(2 * time.Second) {
		data, _ := json.Marshal(&payload{
			Time: time.Now(),
			Wisdom: wisdom[rand.Intn(len(wisdom))],
			Secret: secret,
			Team: "teacher",
		})
		res := client.Publish("/test/inception", 0, false, data)
		if !res.WaitTimeout(1 * time.Second) {
			log.Fatalf("Error on publish: %s", res.Error())
		}
	}
}
