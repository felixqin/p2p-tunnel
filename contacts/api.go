package contacts

type Option struct {
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

var option *Option

func Open(opt *Option) {
	option = opt
	option.clientId = opt.Name + "_" + makeRandomString(8)
	go func() {
		startMqtt(opt.Server, option.clientId, opt.Username, opt.Password)
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
	return findContact(name)
}
