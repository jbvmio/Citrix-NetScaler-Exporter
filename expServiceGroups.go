package main

import (
	"strconv"

	"citrix-netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

const serviceGroupsSubsystem = "servicegroup"

var serviceGroupsLabels = []string{
	netscalerInstance,
	`servicegroup`,
	`member`,
}

var (
	serviceGroupsState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: serviceGroupsSubsystem,
			Name:      "state",
			Help:      "Current state of the server. 0 = DOWN, 1 = UP, 2 = OUT OF SERVICE, 3 = UNKNOWN",
		},
		serviceGroupsLabels,
	)

	serviceGroupsAvgTTFB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: serviceGroupsSubsystem,
			Name:      "average_time_to_first_byte_seconds",
			Help:      "Average TTFB between the NetScaler appliance and the server. TTFB is the time interval between sending the request packet to a service and receiving the first response from the service.",
		},
		serviceGroupsLabels,
	)

	serviceGroupsTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: serviceGroupsSubsystem,
			Name:      "requests_total",
			Help:      "Total number of requests received on this service",
		},
		serviceGroupsLabels,
	)

	serviceGroupsTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: serviceGroupsSubsystem,
			Name:      "responses_total",
			Help:      "Number of responses received on this service.",
		},
		serviceGroupsLabels,
	)

	serviceGroupsTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: serviceGroupsSubsystem,
			Name:      "request_bytes_total",
			Help:      "Total number of request bytes received on this service",
		},
		serviceGroupsLabels,
	)

	serviceGroupsTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: serviceGroupsSubsystem,
			Name:      "response_bytes_total",
			Help:      "Number of response bytes received by this service",
		},
		serviceGroupsLabels,
	)

	serviceGroupsCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: serviceGroupsSubsystem,
			Name:      "current_client_connections",
			Help:      "Number of current client connections.",
		},
		serviceGroupsLabels,
	)

	serviceGroupsSurgeCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: serviceGroupsSubsystem,
			Name:      "surge_queue",
			Help:      "Number of requests in the surge queue.",
		},
		serviceGroupsLabels,
	)

	serviceGroupsCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: serviceGroupsSubsystem,
			Name:      "current_server_connections",
			Help:      "Number of current connections to the actual servers behind the virtual server.",
		},
		serviceGroupsLabels,
	)

	serviceGroupsServerEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: serviceGroupsSubsystem,
			Name:      "server_established_connections",
			Help:      "Number of server connections in ESTABLISHED state.",
		},
		serviceGroupsLabels,
	)

	serviceGroupsCurrentReusePool = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: serviceGroupsSubsystem,
			Name:      "current_reuse_pool",
			Help:      "Number of requests in the idle queue/reuse pool.",
		},
		serviceGroupsLabels,
	)

	serviceGroupsMaxClients = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: serviceGroupsSubsystem,
			Name:      "max_clients",
			Help:      "Maximum open connections allowed on this service.",
		},
		serviceGroupsLabels,
	)
)

func (e *Exporter) collectServiceGroupsState(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsState.Reset()

	var state float64
	switch sg.State {
	case `DOWN`:
		state = 0.0
	case `UP`:
		state = 1.0
	case `OUT OF SERVICE`:
		state = 2.0
	default:
		state = 3.0
	}

	e.serviceGroupsState.WithLabelValues(e.nsInstance, sgName, servername).Set(state)
}

func (e *Exporter) collectServiceGroupsAvgTTFB(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsAvgTTFB.Reset()

	var serviceGroupsAvgTTFBInSeconds float64
	val, _ := strconv.ParseFloat(sg.AvgTimeToFirstByte, 64)
	serviceGroupsAvgTTFBInSeconds = val * 0.001
	e.serviceGroupsAvgTTFB.WithLabelValues(e.nsInstance, sgName, servername).Set(serviceGroupsAvgTTFBInSeconds)
}

func (e *Exporter) collectServiceGroupsTotalRequests(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsTotalRequests.Reset()

	val, _ := strconv.ParseFloat(sg.TotalRequests, 64)
	e.serviceGroupsTotalRequests.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalResponses(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsTotalResponses.Reset()

	val, _ := strconv.ParseFloat(sg.TotalResponses, 64)
	e.serviceGroupsTotalResponses.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalRequestBytes(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsTotalRequestBytes.Reset()

	val, _ := strconv.ParseFloat(sg.TotalRequestBytes, 64)
	e.serviceGroupsTotalRequestBytes.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalResponseBytes(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsTotalResponseBytes.Reset()

	val, _ := strconv.ParseFloat(sg.TotalResponseBytes, 64)
	e.serviceGroupsTotalResponseBytes.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
}

func (e *Exporter) collectServiceGroupsCurrentClientConnections(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsCurrentClientConnections.Reset()

	val, _ := strconv.ParseFloat(sg.CurrentClientConnections, 64)
	e.serviceGroupsCurrentClientConnections.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
}

func (e *Exporter) collectServiceGroupsSurgeCount(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsSurgeCount.Reset()

	val, _ := strconv.ParseFloat(sg.SurgeCount, 64)
	e.serviceGroupsSurgeCount.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
}

func (e *Exporter) collectServiceGroupsCurrentServerConnections(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsCurrentServerConnections.Reset()

	val, _ := strconv.ParseFloat(sg.CurrentServerConnections, 64)
	e.serviceGroupsCurrentServerConnections.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
}

func (e *Exporter) collectServiceGroupsServerEstablishedConnections(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsServerEstablishedConnections.Reset()

	val, _ := strconv.ParseFloat(sg.ServerEstablishedConnections, 64)
	e.serviceGroupsServerEstablishedConnections.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
}

func (e *Exporter) collectServiceGroupsCurrentReusePool(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsCurrentReusePool.Reset()

	val, _ := strconv.ParseFloat(sg.CurrentReusePool, 64)
	e.serviceGroupsCurrentReusePool.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
}

func (e *Exporter) collectServiceGroupsMaxClients(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsMaxClients.Reset()

	val, _ := strconv.ParseFloat(sg.MaxClients, 64)
	e.serviceGroupsMaxClients.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
}
