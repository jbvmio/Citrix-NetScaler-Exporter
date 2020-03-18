package main

import (
	"strconv"

	"citrix-netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

const servicesSubsystem = "service"

var servicesLabels = []string{
	netscalerInstance,
	`citrixadc_service_name`,
	`citrixadc_lb_name`,
}

var (
	// TODO - Convert megabytes to bytes
	servicesThroughput = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "throughput_bytes_total",
			Help:      "Number of bytes received or sent by this service",
		},
		servicesLabels,
	)

	servicesAvgTTFB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "average_time_to_first_byte_seconds",
			Help:      "Average TTFB between the NetScaler appliance and the server. TTFB is the time interval between sending the request packet to a service and receiving the first response from the service",
		},
		servicesLabels,
	)

	servicesState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "state",
			Help:      "Current state of the service. 0 = DOWN, 1 = UP, 2 = OUT OF SERVICE, 3 = UNKNOWN",
		},
		servicesLabels,
	)

	servicesTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "requests_total",
			Help:      "Total number of requests received on this service",
		},
		servicesLabels,
	)

	servicesTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "responses_total",
			Help:      "Total number of responses received on this service",
		},
		servicesLabels,
	)

	servicesTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "request_bytes_total",
			Help:      "Total number of request bytes received on this service",
		},
		servicesLabels,
	)

	servicesTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "response_bytes_total",
			Help:      "Total number of response bytes received on this service",
		},
		servicesLabels,
	)

	servicesCurrentClientConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "current_client_connections",
			Help:      "Number of current client connections",
		},
		servicesLabels,
	)

	servicesSurgeCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "surge_queue",
			Help:      "Number of requests in the surge queue",
		},
		servicesLabels,
	)

	servicesCurrentServerConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "current_server_connections",
			Help:      "Number of current connections to the actual servers",
		},
		servicesLabels,
	)

	servicesServerEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "server_established_connections",
			Help:      "Number of server connections in ESTABLISHED state",
		},
		servicesLabels,
	)

	servicesCurrentReusePool = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "current_reuse_pool",
			Help:      "Number of requests in the idle queue/reuse pool.",
		},
		servicesLabels,
	)

	servicesMaxClients = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "max_clients",
			Help:      "Maximum open connections allowed on this service",
		},
		servicesLabels,
	)

	servicesCurrentLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "current_load",
			Help:      "Load on the service that is calculated from the bound load based monitor",
		},
		servicesLabels,
	)

	servicesVirtualServerServiceHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "vserver_service_hits_total",
			Help:      "Number of times that the service has been provided",
		},
		servicesLabels,
	)

	servicesActiveTransactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: servicesSubsystem,
			Name:      "active_transactions",
			Help:      "Number of active transactions handled by this service. (Including those in the surge queue.) Active Transaction means number of transactions currently served by the server including those waiting in the SurgeQ",
		},
		servicesLabels,
	)
)

func (e *Exporter) collectServicesThroughput(ns netscaler.NSAPIResponse) {
	e.servicesThroughput.Reset()

	for _, service := range ns.ServiceStats {
		var throughputInBytes float64
		val, _ := strconv.ParseFloat(service.Throughput, 64)
		// Value is in megabytes. Convert to base unit of bytes
		throughputInBytes = val * 1024 * 1024
		e.servicesThroughput.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(throughputInBytes)
	}
}

func (e *Exporter) collectServicesAvgTTFB(ns netscaler.NSAPIResponse) {
	e.servicesAvgTTFB.Reset()

	for _, service := range ns.ServiceStats {
		var servicesAvgTTFBInSeconds float64
		val, _ := strconv.ParseFloat(service.AvgTimeToFirstByte, 64)
		servicesAvgTTFBInSeconds = val * 0.001
		e.servicesAvgTTFB.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(servicesAvgTTFBInSeconds)
	}
}

func (e *Exporter) collectServicesState(ns netscaler.NSAPIResponse) {
	e.servicesState.Reset()

	for _, service := range ns.ServiceStats {
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
		e.servicesState.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(state)
	}
}

func (e *Exporter) collectServicesTotalRequests(ns netscaler.NSAPIResponse) {
	e.servicesTotalRequests.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequests, 64)
		e.servicesTotalRequests.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}

func (e *Exporter) collectServicesTotalResponses(ns netscaler.NSAPIResponse) {
	e.servicesTotalResponses.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponses, 64)
		e.servicesTotalResponses.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}

func (e *Exporter) collectServicesTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.servicesTotalRequestBytes.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequestBytes, 64)
		e.servicesTotalRequestBytes.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}

func (e *Exporter) collectServicesTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.servicesTotalResponseBytes.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponseBytes, 64)
		e.servicesTotalResponseBytes.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentClientConns(ns netscaler.NSAPIResponse) {
	e.servicesCurrentClientConns.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentClientConnections, 64)
		e.servicesCurrentClientConns.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}

func (e *Exporter) collectServicesSurgeCount(ns netscaler.NSAPIResponse) {
	e.servicesSurgeCount.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.SurgeCount, 64)
		e.servicesSurgeCount.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentServerConns(ns netscaler.NSAPIResponse) {
	e.servicesCurrentServerConns.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentServerConnections, 64)
		e.servicesCurrentServerConns.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}

func (e *Exporter) collectServicesServerEstablishedConnections(ns netscaler.NSAPIResponse) {
	e.servicesServerEstablishedConnections.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServerEstablishedConnections, 64)
		e.servicesServerEstablishedConnections.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentReusePool(ns netscaler.NSAPIResponse) {
	e.servicesCurrentReusePool.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentReusePool, 64)
		e.servicesCurrentReusePool.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}

func (e *Exporter) collectServicesMaxClients(ns netscaler.NSAPIResponse) {
	e.servicesMaxClients.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.MaxClients, 64)
		e.servicesMaxClients.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}

/*
func (e *Exporter) collectServicesCurrentLoad(ns netscaler.NSAPIResponse) {
	e.servicesCurrentLoad.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentLoad, 64)
		e.servicesCurrentLoad.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}

func (e *Exporter) collectServicesVirtualServerServiceHits(ns netscaler.NSAPIResponse) {
	e.servicesVirtualServerServiceHits.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServiceHits, 64)
		e.servicesVirtualServerServiceHits.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}
*/

func (e *Exporter) collectServicesActiveTransactions(ns netscaler.NSAPIResponse) {
	e.servicesActiveTransactions.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ActiveTransactions, 64)
		e.servicesActiveTransactions.WithLabelValues(e.nsInstance, service.Name, getValue(vipDB.db, service.Name)).Set(val)
	}
}
