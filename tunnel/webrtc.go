package tunnel

import (
	"github.com/pions/webrtc"
)

type IceServers []string

func newWebRTC(conf *IceServers) (*webrtc.RTCPeerConnection, error) {
	return webrtc.New(webrtc.RTCConfiguration{
		IceServers: []webrtc.RTCIceServer{
			{
				URLs: *conf,
			},
		},
	})
}
