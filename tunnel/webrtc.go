package tunnel

import (
	"github.com/pions/webrtc"
)

type IceOptions []string

func newWebRTC(conf *IceOptions) (*webrtc.RTCPeerConnection, error) {
	return webrtc.New(webrtc.RTCConfiguration{
		IceServers: []webrtc.RTCIceServer{
			{
				URLs: *conf,
			},
		},
	})
}
