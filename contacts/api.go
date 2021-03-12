package contacts

import "fmt"

type Configure struct {
	Name     string `yaml:"name"`
	Server   string `yaml:"server"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	clientId string
}

type Contact struct {
	Name     string
	ClientId string
	Owner    string
}

type Offer struct {
	Sdp string `json:"sdp"`
}

type Answer struct {
	Sdp string `json:"sdp"`
}

var configure Configure

func Open(conf *Configure) {
	configure = *conf
	configure.clientId = conf.Name + "_" + makeRandomString(8)
	go func() {
		startMqtt(conf.Server, configure.clientId, conf.Username, conf.Password)
	}()
}

func Close() {
	stopMqtt()
}

func HandleOfferFunc(handler func(fromClient string, offer *Offer)) {
	offerHandler = handler
}

func HandleAnswerFunc(handler func(fromClient string, answer *Answer)) {
	answerHandler = handler
}

func SendOffer(clientId string, offer *Offer) error {
	return sendMessageToClient(clientId, "sendOffer", offer)
}

func SendAnswer(clientId string, answer *Answer) error {
	return sendMessageToClient(clientId, "sendAnswer", answer)
}

func FindContact(name string) (*Contact, error) {
	for _, contact := range contacts {
		if contact.Name == name {
			return contact, nil
		}
	}

	return nil, fmt.Errorf("contact not found")
}
