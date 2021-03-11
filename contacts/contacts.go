package contacts

type Configure struct {
	Name     string `yaml:"name"`
	Server   string `yaml:"server"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var configure Configure

func Open(conf Configure) {
	configure = conf
	clientId := conf.Name + "_" + makeRandomString(8)
	openMqtt(conf.Server, clientId, conf.Username, conf.Password)
}

func Close() {
	closeMqtt()
}
