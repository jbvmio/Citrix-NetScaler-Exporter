package main

import (
	"strconv"

	"citrix-netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

const csVirtualServersSubsystem = "cs_vserver"

var csVirtualServersLabels = []string{
	netscalerInstance,
	`citrixadc_cs_name`,
}

var (
	csVirtualServersState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "state",
			Help:      "Current state of the server. 0 = DOWN, 1 = UP, 2 = OUT OF SERVICE, 3 = UNKNOWN",
		},
		csVirtualServersLabels,
	)

	csVirtualServersTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "hits_total",
			Help:      "Total virtual server hits",
		},
		csVirtualServersLabels,
	)

	csVirtualServersTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "requests_total",
			Help:      "Total virtual server requests",
		},
		csVirtualServersLabels,
	)

	csVirtualServersTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "responses_total",
			Help:      "Total virtual server responses",
		},
		csVirtualServersLabels,
	)

	csVirtualServersTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "request_bytes_total",
			Help:      "Total virtual server request bytes",
		},
		csVirtualServersLabels,
	)

	csVirtualServersTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "response_bytes_total",
			Help:      "Total virtual server response bytes",
		},
		csVirtualServersLabels,
	)

	csVirtualServersCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "current_client_connections",
			Help:      "Number of current client connections on a specific virtual server",
		},
		csVirtualServersLabels,
	)

	csVirtualServersCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "current_server_connections",
			Help:      "Number of current connections to the actual servers behind the specific virtual server.",
		},
		csVirtualServersLabels,
	)

	csVirtualServersEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "established_connections",
			Help:      "Number of client connections in ESTABLISHED state.",
		},
		csVirtualServersLabels,
	)

	csVirtualServersTotalPacketsReceived = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "packets_received_total",
			Help:      "Total number of packets received",
		},
		csVirtualServersLabels,
	)

	csVirtualServersTotalPacketsSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "packets_sent_total",
			Help:      "Total number of packets sent.",
		},
		csVirtualServersLabels,
	)

	csVirtualServersTotalSpillovers = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "spillovers_total",
			Help:      "Number of times vserver experienced spill over.",
		},
		csVirtualServersLabels,
	)

	csVirtualServersDeferredRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "deferred_requests_total",
			Help:      "Number of deferred request on this vserver",
		},
		csVirtualServersLabels,
	)

	csVirtualServersNumberInvalidRequestResponse = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "invalid_request_response_total",
			Help:      "Number invalid requests/responses on this vserver",
		},
		csVirtualServersLabels,
	)

	csVirtualServersNumberInvalidRequestResponseDropped = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "invalid_request_response_dropped_total",
			Help:      "Number invalid requests/responses dropped on this vserver",
		},
		csVirtualServersLabels,
	)

	csVirtualServersTotalVServerDownBackupHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "vserver_down_backup_hits_total",
			Help:      "Number of times traffic was diverted to backup vserver since primary vserver was DOWN.",
		},
		csVirtualServersLabels,
	)

	csVirtualServersCurrentMultipathSessions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "current_multipath_sessions",
			Help:      "Current Multipath TCP sessions",
		},
		csVirtualServersLabels,
	)

	csVirtualServersCurrentMultipathSubflows = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: csVirtualServersSubsystem,
			Name:      "current_multipath_subflows",
			Help:      "Current Multipath TCP subflows",
		},
		csVirtualServersLabels,
	)
)

func (e *Exporter) collectCSVirtualServerState(ns netscaler.NSAPIResponse) {
	e.csVirtualServersState.Reset()
	for _, vs := range ns.CSVirtualServerStats {
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
		e.csVirtualServersState.WithLabelValues(e.nsInstance, vs.Name).Set(state)
	}
}

func (e *Exporter) collectCSVirtualServerTotalHits(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalHits.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		totalHits, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.csVirtualServersTotalHits.WithLabelValues(e.nsInstance, vs.Name).Set(totalHits)
	}
}

func (e *Exporter) collectCSVirtualServerTotalRequests(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalRequests.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalRequests, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.csVirtualServersTotalRequests.WithLabelValues(e.nsInstance, vs.Name).Set(totalRequests)
	}
}

func (e *Exporter) collectCSVirtualServerTotalResponses(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalResponses.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalResponses, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.csVirtualServersTotalResponses.WithLabelValues(e.nsInstance, vs.Name).Set(totalResponses)
	}
}

func (e *Exporter) collectCSVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalRequestBytes.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalRequestBytes, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.csVirtualServersTotalRequestBytes.WithLabelValues(e.nsInstance, vs.Name).Set(totalRequestBytes)
	}
}

func (e *Exporter) collectCSVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalResponseBytes.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalResponseBytes, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.csVirtualServersTotalResponseBytes.WithLabelValues(e.nsInstance, vs.Name).Set(totalResponseBytes)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse) {
	e.csVirtualServersCurrentClientConnections.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		currentClientConnections, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.csVirtualServersCurrentClientConnections.WithLabelValues(e.nsInstance, vs.Name).Set(currentClientConnections)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse) {
	e.csVirtualServersCurrentServerConnections.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		currentServerConnections, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.csVirtualServersCurrentServerConnections.WithLabelValues(e.nsInstance, vs.Name).Set(currentServerConnections)
	}
}

func (e *Exporter) collectCSVirtualServerEstablishedConnections(ns netscaler.NSAPIResponse) {
	e.csVirtualServersEstablishedConnections.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		EstablishedConnections, _ := strconv.ParseFloat(vs.EstablishedConnections, 64)
		e.csVirtualServersEstablishedConnections.WithLabelValues(e.nsInstance, vs.Name).Set(EstablishedConnections)
	}
}

func (e *Exporter) collectCSVirtualServerTotalPacketsReceived(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalPacketsReceived.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalPacketsReceived, _ := strconv.ParseFloat(vs.TotalPacketsReceived, 64)
		e.csVirtualServersTotalPacketsReceived.WithLabelValues(e.nsInstance, vs.Name).Set(totalPacketsReceived)
	}
}

func (e *Exporter) collectCSVirtualServerTotalPacketsSent(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalPacketsSent.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalPacketsSent, _ := strconv.ParseFloat(vs.TotalPacketsSent, 64)
		e.csVirtualServersTotalPacketsSent.WithLabelValues(e.nsInstance, vs.Name).Set(totalPacketsSent)
	}
}

func (e *Exporter) collectCSVirtualServerTotalSpillovers(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalSpillovers.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalSpillovers, _ := strconv.ParseFloat(vs.TotalSpillovers, 64)
		e.csVirtualServersTotalSpillovers.WithLabelValues(e.nsInstance, vs.Name).Set(totalSpillovers)
	}
}

func (e *Exporter) collectCSVirtualServerDeferredRequests(ns netscaler.NSAPIResponse) {
	e.csVirtualServersDeferredRequests.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		deferredRequests, _ := strconv.ParseFloat(vs.DeferredRequests, 64)
		e.csVirtualServersDeferredRequests.WithLabelValues(e.nsInstance, vs.Name).Set(deferredRequests)
	}
}

func (e *Exporter) collectCSVirtualServerNumberInvalidRequestResponse(ns netscaler.NSAPIResponse) {
	e.csVirtualServersNumberInvalidRequestResponse.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		numberInvalidRequestResponse, _ := strconv.ParseFloat(vs.InvalidRequestResponse, 64)
		e.csVirtualServersNumberInvalidRequestResponse.WithLabelValues(e.nsInstance, vs.Name).Set(numberInvalidRequestResponse)
	}
}

func (e *Exporter) collectCSVirtualServerNumberInvalidRequestResponseDropped(ns netscaler.NSAPIResponse) {
	e.csVirtualServersNumberInvalidRequestResponseDropped.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		numberInvalidRequestResponseDropped, _ := strconv.ParseFloat(vs.InvalidRequestResponseDropped, 64)
		e.csVirtualServersNumberInvalidRequestResponseDropped.WithLabelValues(e.nsInstance, vs.Name).Set(numberInvalidRequestResponseDropped)
	}
}

func (e *Exporter) collectCSVirtualServerTotalVServerDownBackupHits(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalVServerDownBackupHits.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalVServerDownBackupHits, _ := strconv.ParseFloat(vs.TotalVServerDownBackupHits, 64)
		e.csVirtualServersTotalVServerDownBackupHits.WithLabelValues(e.nsInstance, vs.Name).Set(totalVServerDownBackupHits)
	}
}

/*
func (e *Exporter) collectCSVirtualServerCurrentMultipathSessions(ns netscaler.NSAPIResponse) {
	e.csVirtualServersCurrentMultipathSessions.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		currentMultipathSessions, _ := strconv.ParseFloat(vs.CurrentMultipathSessions, 64)
		e.csVirtualServersCurrentMultipathSessions.WithLabelValues(e.nsInstance, vs.Name).Set(currentMultipathSessions)
	}
}


func (e *Exporter) collectCSVirtualServerCurrentMultipathSubflows(ns netscaler.NSAPIResponse) {
	e.csVirtualServersCurrentMultipathSubflows.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		currentMultipathSubflows, _ := strconv.ParseFloat(vs.CurrentMultipathSubflows, 64)
		e.csVirtualServersCurrentMultipathSubflows.WithLabelValues(e.nsInstance, vs.Name).Set(currentMultipathSubflows)
	}
}
*/
