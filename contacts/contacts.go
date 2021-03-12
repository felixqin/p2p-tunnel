package contacts

type Configure struct {
	Name     string `yaml:"name"`
	Server   string `yaml:"server"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	clientId string
}

var configure Configure

func Open(conf Configure) {
	configure = conf
	configure.clientId = conf.Name + "_" + makeRandomString(8)
	go func() {
		startMqtt(conf.Server, configure.clientId, conf.Username, conf.Password)
	}()
}

func Close() {
	stopMqtt()
}
