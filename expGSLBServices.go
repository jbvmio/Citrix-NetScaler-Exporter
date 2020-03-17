package main

import (
	"strconv"

	"github.com/jbvmio/citrix-netscaler-exporter/netscaler"
	"github.com/prometheus/client_golang/prometheus"
)

const gslbServicesSubsystem = "gslb_service"

var gslbServicesLabels = []string{
	netscalerInstance,
	`citrixadc_service_name`,
}

var (
	gslbServicesState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: gslbServicesSubsystem,
			Name:      "state",
			Help:      "Current state of the service. 0 = DOWN, 1 = UP, 2 = OUT OF SERVICE, 3 = UNKNOWN",
		},
		gslbServicesLabels,
	)

	gslbServicesTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: gslbServicesSubsystem,
			Name:      "requests_total",
			Help:      "Total number of requests received on this service",
		},
		gslbServicesLabels,
	)

	gslbServicesTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: gslbServicesSubsystem,
			Name:      "responses_total",
			Help:      "Total number of responses received on this service",
		},
		gslbServicesLabels,
	)

	gslbServicesTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: gslbServicesSubsystem,
			Name:      "request_bytes_total",
			Help:      "Total number of request bytes received on this service",
		},
		gslbServicesLabels,
	)

	gslbServicesTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: gslbServicesSubsystem,
			Name:      "response_bytes_total",
			Help:      "Total number of response bytes received on this service",
		},
		gslbServicesLabels,
	)

	gslbServicesCurrentClientConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: gslbServicesSubsystem,
			Name:      "current_client_connections",
			Help:      "Number of current client connections",
		},
		gslbServicesLabels,
	)

	gslbServicesCurrentServerConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: gslbServicesSubsystem,
			Name:      "current_server_connections",
			Help:      "Number of current connections to the actual servers",
		},
		gslbServicesLabels,
	)

	gslbServicesEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: gslbServicesSubsystem,
			Name:      "established_connections",
			Help:      "Number of server connections in ESTABLISHED state",
		},
		gslbServicesLabels,
	)

	gslbServicesCurrentLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: gslbServicesSubsystem,
			Name:      "current_load",
			Help:      "Load on the service that is calculated from the bound load based monitor",
		},
		gslbServicesLabels,
	)

	gslbServicesVirtualServerServiceHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: gslbServicesSubsystem,
			Name:      "virtual_server_service_hits_total",
			Help:      "Number of times that the service has been provided",
		},
		gslbServicesLabels,
	)
)

func (e *Exporter) collectGSLBServicesState(ns netscaler.NSAPIResponse) {
	e.gslbServicesState.Reset()

	for _, service := range ns.GSLBServiceStats {
		var state float64
		switch service.State {
		case `DOWN`:
			state = 0.0
		case `UP`:
			state = 1.0
		case `OUT OF SERVICE`:
			state = 2.0
		default:
			state = 3.0
		}

		e.gslbServicesState.WithLabelValues(e.nsInstance, service.Name).Set(state)
	}
}

func (e *Exporter) collectGSLBServicesTotalRequests(ns netscaler.NSAPIResponse) {
	e.gslbServicesTotalRequests.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequests, 64)
		e.gslbServicesTotalRequests.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalResponses(ns netscaler.NSAPIResponse) {
	e.gslbServicesTotalResponses.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponses, 64)
		e.gslbServicesTotalResponses.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.gslbServicesTotalRequestBytes.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequestBytes, 64)
		e.gslbServicesTotalRequestBytes.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.gslbServicesTotalResponseBytes.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponseBytes, 64)
		e.gslbServicesTotalResponseBytes.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentClientConns(ns netscaler.NSAPIResponse) {
	e.gslbServicesCurrentClientConns.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentClientConnections, 64)
		e.gslbServicesCurrentClientConns.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentServerConns(ns netscaler.NSAPIResponse) {
	e.gslbServicesCurrentServerConns.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentServerConnections, 64)
		e.gslbServicesCurrentServerConns.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesEstablishedConnections(ns netscaler.NSAPIResponse) {
	e.gslbServicesEstablishedConnections.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.EstablishedConnections, 64)
		e.gslbServicesEstablishedConnections.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

/*
func (e *Exporter) collectGSLBServicesCurrentLoad(ns netscaler.NSAPIResponse) {
	e.gslbServicesCurrentLoad.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentLoad, 64)
		e.gslbServicesCurrentLoad.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesVirtualServerServiceHits(ns netscaler.NSAPIResponse) {
	e.gslbServicesVirtualServerServiceHits.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.ServiceHits, 64)
		e.gslbServicesVirtualServerServiceHits.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}
*/
