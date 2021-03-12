package main

import (
	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/datachannel"
)

type iceConfigure []string

func newWebRTC(conf *iceConfigure) (*webrtc.RTCPeerConnection, error) {
	return webrtc.New(webrtc.RTCConfiguration{
		IceServers: []webrtc.RTCIceServer{
			{
				URLs: *conf,
			},
		},
	})
}

type rtcWriter struct {
	*webrtc.RTCDataChannel
}

func newWebRTCWriter(dc *webrtc.RTCDataChannel) *rtcWriter {
	return &rtcWriter{dc}
}

func (s *rtcWriter) Write(b []byte) (int, error) {
	err := s.RTCDataChannel.Send(datachannel.PayloadBinary{Data: b})
	return len(b), err
}
