package main

import (
	"log"
	"net/http"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/bluenviron/gortsplib/v4/pkg/format/rtph264"
	"github.com/pion/rtp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	c := gortsplib.Client{}

	u, err := base.ParseURL("rtsp://username:password@host:port/path")
	if err != nil {
		panic(err)
	}

	// connect to the server
	err = c.Start(u.Scheme, u.Host)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// find available medias
	desc, _, err := c.Describe(u)
	if err != nil {
		panic(err)
	}

	// find the H264 media and format
	var forma *format.H264
	medi := desc.FindFormat(&forma)
	if medi == nil {
		panic("media not found")
	}

	// setup RTP -> H264 decoder
	rtpDec, err := forma.CreateDecoder()
	if err != nil {
		panic(err)
	}

	// setup H264 -> raw frames decoder
	frameDec := &h264Decoder{}
	err = frameDec.initialize()
	if err != nil {
		panic(err)
	}
	defer frameDec.close()

	// if SPS and PPS are present into the SDP, send them to the decoder
	if forma.SPS != nil {
		frameDec.decode(forma.SPS)
	}
	if forma.PPS != nil {
		frameDec.decode(forma.PPS)
	}

	// setup a single media
	_, err = c.Setup(desc.BaseURL, medi, 0, 0)
	if err != nil {
		panic(err)
	}

	// called when a RTP packet arrives
	c.OnPacketRTP(medi, forma, func(pkt *rtp.Packet) {
		// decode timestamp
		_, ok := c.PacketPTS2(medi, pkt)
		if !ok {
			log.Printf("waiting for timestamp")
			return
		}

		// extract access units from RTP packets
		au, err := rtpDec.Decode(pkt)
		if err != nil {
			if err != rtph264.ErrNonStartingPacketAndNoPrevious && err != rtph264.ErrMorePacketsNeeded {
				log.Printf("ERR: %v", err)
			}
			return
		}

		for _, nalu := range au {
			// convert NALUs into RGBA frames
			img, err := frameDec.decode(nalu)
			if err != nil {
				panic(err)
			}

			// wait for a frame
			if img == nil {
				continue
			}

			// emit metrics after every frame
			updateRTSPMetrics(c.Stats())
		}
	})

	_, err = c.Play(nil)
	if err != nil {
		panic(err)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

func updateRTSPMetrics(stats *gortsplib.ClientStats) {
	// RTSP
	ClientConnBytesReceived.Set(float64(stats.Conn.BytesReceived))
	ClientConnBytesSent.Set(float64(stats.Conn.BytesSent))
	SessionBytesReceived.Set(float64(stats.Session.BytesReceived))

	// RTP
	SessionRTPPacketsReceived.Add(float64(stats.Session.RTPPacketsReceived))
	SessionRTPPacketsLost.Add(float64(stats.Session.RTPPacketsLost))
	SessionRTPPacketsInError.Add(float64(stats.Session.RTPPacketsInError))
	SessionRTPJitter.Set(stats.Session.RTPPacketsJitter)

	// RTCP
	SessionRTCPPacketsReceived.Add(float64(stats.Session.RTCPPacketsReceived))
	SessionRTCPPacketsInError.Add(float64(stats.Session.RTCPPacketsInError))
}
