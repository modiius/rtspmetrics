package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	// -- Media --
	FramesPerSecond prometheus.Gauge

	// -- RTSP --
	ClientConnBytesReceived prometheus.Gauge
	ClientConnBytesSent     prometheus.Gauge
	SessionBytesReceived    prometheus.Gauge

	// -- RTP --
	SessionRTPPacketsReceived prometheus.Counter
	SessionRTPPacketsLost     prometheus.Counter
	SessionRTPPacketsInError  prometheus.Counter
	SessionRTPJitter          prometheus.Gauge

	// -- RTCP --
	SessionRTCPPacketsReceived prometheus.Counter
	SessionRTCPPacketsInError  prometheus.Counter
}

func NewMetrics(name, url string) *Metrics {
	labelSet := prometheus.Labels{"name": name, "url": url}

	return &Metrics{
		FramesPerSecond: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        "frames_per_second",
			Help:        "Frames per second of a video stream.",
			ConstLabels: labelSet,
		}),

		ClientConnBytesReceived: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        "client_conn_bytes_received",
			Help:        "Total bytes received on the client connection.",
			ConstLabels: labelSet,
		}),

		ClientConnBytesSent: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        "client_conn_bytes_sent",
			Help:        "Total bytes sent on the client connection.",
			ConstLabels: labelSet,
		}),

		SessionBytesReceived: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        "session_bytes_received",
			Help:        "Total bytes received during the session.",
			ConstLabels: labelSet,
		}),

		SessionRTPPacketsReceived: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "session_rtp_packets_received_total",
			Help:        "Total RTP packets received during the session.",
			ConstLabels: labelSet,
		}),

		SessionRTPPacketsLost: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "session_rtp_packets_lost_total",
			Help:        "Total RTP packets lost during the session.",
			ConstLabels: labelSet,
		}),

		SessionRTPPacketsInError: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "session_rtp_packets_in_error_total",
			Help:        "Total RTP packets in error during the session.",
			ConstLabels: labelSet,
		}),

		SessionRTPJitter: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        "session_rtp_jitter",
			Help:        "Average RTP jitter during the session.",
			ConstLabels: labelSet,
		}),

		SessionRTCPPacketsReceived: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "session_rtcp_packets_received_total",
			Help:        "Total RTCP packets received during the session.",
			ConstLabels: labelSet,
		}),

		SessionRTCPPacketsInError: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "session_rtcp_packets_in_error_total",
			Help:        "Total RTCP packets in error during the session.",
			ConstLabels: labelSet,
		}),
	}
}
