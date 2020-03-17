package main

import (
	"strconv"

	"github.com/jbvmio/citrix-netscaler-exporter/netscaler"
	"github.com/prometheus/client_golang/prometheus"
)

const interfacesSubsystem = "interface"

var interfacesLabels = []string{
	netscalerInstance,
	`interface`,
	`alias`,
}

var (
	interfacesRxBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: interfacesSubsystem,
			Name:      "received_bytes_total",
			Help:      "Number of bytes received by specific interfaces.",
		},
		interfacesLabels,
	)

	interfacesTxBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: interfacesSubsystem,
			Name:      "transmitted_bytes_total",
			Help:      "Number of bytes transmitted by specific interfaces.",
		},
		interfacesLabels,
	)

	interfacesRxPackets = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: interfacesSubsystem,
			Name:      "received_packets_total",
			Help:      "Number of packets received by specific interfaces",
		},
		interfacesLabels,
	)

	interfacesTxPackets = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: interfacesSubsystem,
			Name:      "transmitted_packets_total",
			Help:      "Number of packets transmitted by specific interfaces",
		},
		interfacesLabels,
	)

	interfacesJumboPacketsRx = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: interfacesSubsystem,
			Name:      "jumbo_packets_received_total",
			Help:      "Number of bytes received by specific interfaces",
		},
		interfacesLabels,
	)

	interfacesJumboPacketsTx = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: interfacesSubsystem,
			Name:      "jumbo_packets_transmitted_total",
			Help:      "Number of jumbo packets transmitted by specific interfaces",
		},
		interfacesLabels,
	)

	interfacesErrorPacketsRx = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: interfacesSubsystem,
			Name:      "error_packets_received_total",
			Help:      "Number of error packets received by specific interfaces",
		},
		interfacesLabels,
	)
)

func (e *Exporter) collectInterfacesRxBytes(ns netscaler.NSAPIResponse) {
	e.interfacesRxBytes.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalReceivedBytes, 64)
		e.interfacesRxBytes.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesTxBytes(ns netscaler.NSAPIResponse) {
	e.interfacesTxBytes.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalTransmitBytes, 64)
		e.interfacesTxBytes.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesRxPackets(ns netscaler.NSAPIResponse) {
	e.interfacesRxPackets.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalReceivedPackets, 64)
		e.interfacesRxPackets.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesTxPackets(ns netscaler.NSAPIResponse) {
	e.interfacesTxPackets.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalTransmitPackets, 64)
		e.interfacesTxPackets.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsRx(ns netscaler.NSAPIResponse) {
	e.interfacesJumboPacketsRx.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.JumboPacketsReceived, 64)
		e.interfacesJumboPacketsRx.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsTx(ns netscaler.NSAPIResponse) {
	e.interfacesJumboPacketsTx.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.JumboPacketsTransmitted, 64)
		e.interfacesJumboPacketsTx.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesErrorPacketsRx(ns netscaler.NSAPIResponse) {
	e.interfacesErrorPacketsRx.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.ErrorPacketsReceived, 64)
		e.interfacesErrorPacketsRx.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}
