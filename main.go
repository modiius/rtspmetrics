package main

import (
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/bluenviron/gortsplib/v4/pkg/format/rtph264"
	"github.com/pion/rtp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Sessions []SessionConfig `yaml:"sessions"`
}

type SessionConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func main() {
	// -- Parse Configuration --
	configFilePath := "config.yml"
	rawBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(rawBytes, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// -- Start RTSP Sessions --
	for _, sessionConfig := range config.Sessions {
		go startSession(sessionConfig)
	}

	// -- Start Prometheus Metrics Server --
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Serving metrics on :8080")
	http.ListenAndServe(":8080", nil)
}

func startSession(config SessionConfig) {
	// -- Initialization --
	metrics := NewMetrics(config.Name, config.URL)

	url, err := base.ParseURL(config.URL)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}

	c := gortsplib.Client{}
	err = c.Start(url.Scheme, url.Host)
	if err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}
	defer c.Close()

	// -- Media Discovery and Setup --
	desc, _, err := c.Describe(url)
	if err != nil {
		log.Fatalf("Failed to describe RTSP server: %v", err)
	}

	var forma *format.H264
	medi := desc.FindFormat(&forma)
	if medi == nil {
		log.Fatalf("H264 media not found")
	}

	rtpDec, err := forma.CreateDecoder()
	if err != nil {
		log.Fatalf("Failed to create RTP decoder: %v", err)
	}

	frameDec := &h264Decoder{}
	err = frameDec.initialize()
	if err != nil {
		log.Fatalf("Failed to initialize frame decoder: %v", err)
	}
	defer frameDec.close()

	if forma.SPS != nil {
		frameDec.decode(forma.SPS)
	}
	if forma.PPS != nil {
		frameDec.decode(forma.PPS)
	}

	_, err = c.Setup(desc.BaseURL, medi, 0, 0)
	if err != nil {
		log.Fatalf("Failed to setup media: %v", err)
	}

	// -- Frame Rate Calculation --
	var frameCounter atomic.Uint32
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			fps := frameCounter.Swap(0)
			metrics.FramesPerSecond.Set(float64(fps))
		}
	}()

	// -- Handle RTP Packets --
	c.OnPacketRTP(medi, forma, func(pkt *rtp.Packet) {
		// Decode timestamp
		_, ok := c.PacketPTS2(medi, pkt)
		if !ok {
			log.Printf("Waiting for timestamp")
			return
		}

		// Decode RTP Packet
		au, err := rtpDec.Decode(pkt)
		if err != nil {
			if err != rtph264.ErrNonStartingPacketAndNoPrevious && err != rtph264.ErrMorePacketsNeeded {
				log.Printf("Error decoding RTP packet: %v", err)
			}
			return
		}

		// Process NALUs and update metrics
		for _, nalu := range au {
			img, err := frameDec.decode(nalu)
			if err != nil {
				log.Fatalf("Failed to decode frame: %v", err)
			}
			if img == nil {
				continue
			}

			frameCounter.Add(1)
			stats := c.Stats()

			// RTSP Metrics
			metrics.ClientConnBytesReceived.Set(float64(stats.Conn.BytesReceived))
			metrics.ClientConnBytesSent.Set(float64(stats.Conn.BytesSent))
			metrics.SessionBytesReceived.Set(float64(stats.Session.BytesReceived))

			// RTP Metrics
			metrics.SessionRTPPacketsReceived.Add(float64(stats.Session.RTPPacketsReceived))
			metrics.SessionRTPPacketsLost.Add(float64(stats.Session.RTPPacketsLost))
			metrics.SessionRTPPacketsInError.Add(float64(stats.Session.RTPPacketsInError))
			metrics.SessionRTPJitter.Set(stats.Session.RTPPacketsJitter)

			// RTCP Metrics
			metrics.SessionRTCPPacketsReceived.Add(float64(stats.Session.RTCPPacketsReceived))
			metrics.SessionRTCPPacketsInError.Add(float64(stats.Session.RTCPPacketsInError))
		}
	})

	// -- Start Streaming --
	_, err = c.Play(nil)
	if err != nil {
		log.Fatalf("Failed to play stream: %v", err)
	}

	if err := c.Wait(); err != nil {
		log.Fatal("Failed to keep client connection open.")
	}
}
