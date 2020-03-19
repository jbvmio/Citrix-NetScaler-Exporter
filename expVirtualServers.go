package main

import (
	"strconv"

	"github.com/jbvmio/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

const virtualServersSubsystem = "lb_vserver"

var virtualServersLabels = []string{
	netscalerInstance,
	`citrixadc_lb_name`,
}

var (
	virtualServersWaitingRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: virtualServersSubsystem,
			Name:      "waiting_requests",
			Help:      "Number of requests waiting on a specific virtual server",
		},
		virtualServersLabels,
	)

	virtualServersHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: virtualServersSubsystem,
			Name:      "health",
			Help:      "Percentage of UP services bound to a specific virtual server",
		},
		virtualServersLabels,
	)

	virtualServersInactiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: virtualServersSubsystem,
			Name:      "inactive_services",
			Help:      "Number of inactive services bound to a specific virtual server",
		},
		virtualServersLabels,
	)

	virtualServersActiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: virtualServersSubsystem,
			Name:      "active_services",
			Help:      "Number of active services bound to a specific virtual server",
		},
		virtualServersLabels,
	)

	virtualServersTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: virtualServersSubsystem,
			Name:      "hits_total",
			Help:      "Total virtual server hits",
		},
		virtualServersLabels,
	)

	virtualServersTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: virtualServersSubsystem,
			Name:      "requests_total",
			Help:      "Total virtual server requests",
		},
		virtualServersLabels,
	)

	virtualServersTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: virtualServersSubsystem,
			Name:      "responses_total",
			Help:      "Total virtual server responses",
		},
		virtualServersLabels,
	)

	virtualServersTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: virtualServersSubsystem,
			Name:      "request_bytes_total",
			Help:      "Total virtual server request bytes",
		},
		virtualServersLabels,
	)
	virtualServersTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: virtualServersSubsystem,
			Name:      "response_bytes_total",
			Help:      "Total virtual server response bytes",
		},
		virtualServersLabels,
	)

	virtualServersCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: virtualServersSubsystem,
			Name:      "current_client_connections",
			Help:      "Number of current client connections on a specific virtual server",
		},
		virtualServersLabels,
	)

	virtualServersCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: virtualServersSubsystem,
			Name:      "current_server_connections",
			Help:      "Number of current connections to the actual servers behind the specific virtual server.",
		},
		virtualServersLabels,
	)

	virtualServersState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: virtualServersSubsystem,
			Name:      "state",
			Help:      "Current state of the vserver. 0 = DOWN, 1 = UP, 2 = OUT OF SERVICE, 3 = UNKNOWN",
		},
		virtualServersLabels,
	)
)

func (e *Exporter) collectVirtualServerWaitingRequests(ns netscaler.NSAPIResponse) {
	e.virtualServersWaitingRequests.Reset()

	for _, vs := range ns.VirtualServerStats {
		waitingRequests, _ := strconv.ParseFloat(vs.WaitingRequests, 64)
		e.virtualServersWaitingRequests.WithLabelValues(e.nsInstance, vs.Name).Set(waitingRequests)
	}
}

func (e *Exporter) collectVirtualServerHealth(ns netscaler.NSAPIResponse) {
	e.virtualServersHealth.Reset()

	for _, vs := range ns.VirtualServerStats {
		health, _ := strconv.ParseFloat(vs.Health, 64)
		e.virtualServersHealth.WithLabelValues(e.nsInstance, vs.Name).Set(health)
	}
}

func (e *Exporter) collectVirtualServerInactiveServices(ns netscaler.NSAPIResponse) {
	e.virtualServersInactiveServices.Reset()

	for _, vs := range ns.VirtualServerStats {
		inactiveServices, _ := strconv.ParseFloat(vs.InactiveServices, 64)
		e.virtualServersInactiveServices.WithLabelValues(e.nsInstance, vs.Name).Set(inactiveServices)
	}
}

func (e *Exporter) collectVirtualServerActiveServices(ns netscaler.NSAPIResponse) {
	e.virtualServersActiveServices.Reset()

	for _, vs := range ns.VirtualServerStats {
		activeServices, _ := strconv.ParseFloat(vs.ActiveServices, 64)
		e.virtualServersActiveServices.WithLabelValues(e.nsInstance, vs.Name).Set(activeServices)
	}
}

func (e *Exporter) collectVirtualServerTotalHits(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalHits.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalHits, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.virtualServersTotalHits.WithLabelValues(e.nsInstance, vs.Name).Set(totalHits)
	}
}

func (e *Exporter) collectVirtualServerTotalRequests(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalRequests.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalRequests, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.virtualServersTotalRequests.WithLabelValues(e.nsInstance, vs.Name).Set(totalRequests)
	}
}

func (e *Exporter) collectVirtualServerTotalResponses(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalResponses.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalResponses, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.virtualServersTotalResponses.WithLabelValues(e.nsInstance, vs.Name).Set(totalResponses)
	}
}

func (e *Exporter) collectVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalRequestBytes.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalRequestBytes, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.virtualServersTotalRequestBytes.WithLabelValues(e.nsInstance, vs.Name).Set(totalRequestBytes)
	}
}

func (e *Exporter) collectVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalResponseBytes.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalResponseBytes, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.virtualServersTotalResponseBytes.WithLabelValues(e.nsInstance, vs.Name).Set(totalResponseBytes)
	}
}

func (e *Exporter) collectVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse) {
	e.virtualServersCurrentClientConnections.Reset()

	for _, vs := range ns.VirtualServerStats {
		currentClientConnections, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.virtualServersCurrentClientConnections.WithLabelValues(e.nsInstance, vs.Name).Set(currentClientConnections)
	}
}

func (e *Exporter) collectVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse) {
	e.virtualServersCurrentServerConnections.Reset()

	for _, vs := range ns.VirtualServerStats {
		currentServerConnections, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.virtualServersCurrentServerConnections.WithLabelValues(e.nsInstance, vs.Name).Set(currentServerConnections)
	}
}

func (e *Exporter) collectVirtualServerState(ns netscaler.NSAPIResponse) {
	e.virtualServersState.Reset()

	for _, vs := range ns.VirtualServerStats {
		var state float64
		switch vs.State {
		case `DOWN`:
			state = 0.0
		case `UP`:
			state = 1.0
		case `OUT OF SERVICE`:
			state = 2.0
		default:
			state = 3.0
		}

		e.virtualServersState.WithLabelValues(e.nsInstance, vs.Name).Set(state)
	}
}
