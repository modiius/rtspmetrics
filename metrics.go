package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// -- Media --

	FramesPerSecond = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "frames_per_second",
		Help: "Frames per second of a video stream.",
	})

	// -- RTSP --

	ClientConnBytesReceived = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "client_conn_bytes_received",
		Help: "Total bytes received on the client connection.",
	})

	ClientConnBytesSent = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "client_conn_bytes_sent",
		Help: "Total bytes sent on the client connection.",
	})

	SessionBytesReceived = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "session_bytes_received",
		Help: "Total bytes received during the session.",
	})

	// -- RTP --

	SessionRTPPacketsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "session_rtp_packets_received_total",
		Help: "Total RTP packets received during the session.",
	})

	SessionRTPPacketsLost = promauto.NewCounter(prometheus.CounterOpts{
		Name: "session_rtp_packets_lost_total",
		Help: "Total RTP packets lost during the session.",
	})

	SessionRTPPacketsInError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "session_rtp_packets_in_error_total",
		Help: "Total RTP packets in error during the session.",
	})

	SessionRTPJitter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "session_rtp_jitter",
		Help: "Average RTP jitter during the session.",
	})

	// -- RTCP --

	SessionRTCPPacketsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "session_rtcp_packets_received_total",
		Help: "Total RTCP packets received during the session.",
	})

	SessionRTCPPacketsInError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "session_rtcp_packets_in_error_total",
		Help: "Total RTCP packets in error during the session.",
	})
)
