package main

type stubConfigure struct {
	Addr string `yaml:"addr"`
}

func stubServe(conf *stubConfigure) error {
	return nil
}
