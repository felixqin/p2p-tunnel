package contacts

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqc mqtt.Client

func openMqtt(broker string, clientId string, uid string, password string) error {
	var handlers = map[string]func(string, []byte) error{
		"to/user/" + uid + "/msg":        handleToUserMessage,
		"to/client/" + clientId + "/msg": handleToUserMessage,
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientId)
	mqc = mqtt.NewClient(opts)

	go func() {
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
					handler(message.Topic(), message.Payload())
				})
				token.Wait()
				err := token.Error()
				if err != nil {
					log.Println("subscribe failed!", topic, err)
				}

				log.Println("subscribe success!", topic)
			}

			return
		}
	}()

	return nil
}

func closeMqtt() error {
	mqc.Disconnect(1000000)
	return nil
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
