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

// type rtcWriter struct {
// 	*webrtc.RTCDataChannel
// }

// func newWebRTCWriter(dc *webrtc.RTCDataChannel) *rtcWriter {
// 	return &rtcWriter{dc}
// }

// func (s *rtcWriter) Write(b []byte) (int, error) {
// 	err := s.RTCDataChannel.Send(datachannel.PayloadBinary{Data: b})
// 	return len(b), err
// }
