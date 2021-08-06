package tunnel

/// ProxyMessager Proxy 端消息通信器
type ProxyMessager interface {
	SendOffer(sdp string) error
	HandleAnswer(func(sdp string)) error
}

/// StubMessager Stub 端消息通信器
type StubMessager interface {
	SendAnswer(sdp string) error
	HandleOffer(func(sdp string)) error
}
