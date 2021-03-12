package contacts

import (
	"log"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqc mqtt.Client

func startMqtt(broker string, clientId string, uid string, password string) {
	var handlers = map[string]func(string, []byte){
		"to/user/" + uid + "/msg":        handleToUserMessage,
		"to/client/" + clientId + "/msg": handleToUserMessage,
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientId)
	opts.SetConnectionLostHandler(func(mqtt.Client, error) {
		handleServerConnectionLost()
	})

	mqc = mqtt.NewClient(opts)
	for {
		log.Println("mqtt client connect to", broker, "...")
		token := mqc.Connect()
		token.Wait()
		err := token.Error()
		if err != nil {
			log.Println("connect error", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("mqtt client connect success")

		// subscribe topics
		for topic, handler := range handlers {
			token := mqc.Subscribe(topic, 2, func(client mqtt.Client, message mqtt.Message) {
				log.Println("message arrived!", message.Topic())
				// log.Printf("message payload: %s\n", message.Payload())
				handler(parseIdFormTopic(message.Topic()), message.Payload())
			})
			token.Wait()
			err := token.Error()
			if err != nil {
				log.Println("subscribe failed!", topic, err)
			}

			log.Println("subscribe success!", topic)
		}

		handleServerConnected()
		break
	}
}

func stopMqtt() {
	mqc.Disconnect(1000000)
}

func parseIdFormTopic(topic string) string {
	list := strings.Split(topic, "/")
	if len(list) <= 2 {
		return ""
	}

	return list[2]
}

func mqttPublish(topic string, payload []byte) error {
	log.Println("publish topic", topic)
	token := mqc.Publish(topic, 0, false, payload)
	token.Wait()
	err := token.Error()
	if err != nil {
		return err
	}

	return nil
}

func mqttSendToUser(uid string, payload []byte) error {
	topic := "to/user/" + uid + "/msg"
	return mqttPublish(topic, payload)
}

func mqttSendToClient(clientId string, payload []byte) error {
	topic := "to/client/" + clientId + "/msg"
	return mqttPublish(topic, payload)
}
