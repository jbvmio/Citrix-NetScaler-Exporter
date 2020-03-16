package main

import "github.com/prometheus/client_golang/prometheus"

var netscalerLabels = []string{
	netscalerInstance,
}

var (
	modelID = prometheus.NewDesc(
		"model_id",
		"NetScaler model - reflects the bandwidth available; for example VPX 10 would report as 10.",
		netscalerLabels,
		nil,
	)

	mgmtCPUUsage = prometheus.NewDesc(
		"mgmt_cpu_usage",
		"Current CPU utilisation for management",
		netscalerLabels,
		nil,
	)

	pktCPUUsage = prometheus.NewDesc(
		"pkt_cpu_usage",
		"Current CPU utilisation for packet engines, excluding management",
		netscalerLabels,
		nil,
	)

	memUsage = prometheus.NewDesc(
		"mem_usage",
		"Current memory utilisation",
		netscalerLabels,
		nil,
	)

	flashPartitionUsage = prometheus.NewDesc(
		"flash_partition_usage",
		"Used space in /flash partition of the disk, as a percentage.",
		netscalerLabels,
		nil,
	)

	varPartitionUsage = prometheus.NewDesc(
		"var_partition_usage",
		"Used space in /var partition of the disk, as a percentage. ",
		netscalerLabels,
		nil,
	)

	totRxMB = prometheus.NewDesc(
		"total_received_mb",
		"Total number of Megabytes received by the NetScaler appliance",
		netscalerLabels,
		nil,
	)

	totTxMB = prometheus.NewDesc(
		"total_transmit_mb",
		"Total number of Megabytes transmitted by the NetScaler appliance",
		netscalerLabels,
		nil,
	)

	httpRequests = prometheus.NewDesc(
		"http_requests",
		"Total number of HTTP requests received",
		netscalerLabels,
		nil,
	)

	httpResponses = prometheus.NewDesc(
		"http_responses",
		"Total number of HTTP responses sent",
		netscalerLabels,
		nil,
	)

	tcpCurrentClientConnections = prometheus.NewDesc(
		"tcp_current_client_connections",
		"Client connections, including connections in the Opening, Established, and Closing state.",
		netscalerLabels,
		nil,
	)

	tcpCurrentClientConnectionsEstablished = prometheus.NewDesc(
		"tcp_current_client_connections_established",
		"Current client connections in the Established state, which indicates that data transfer can occur between the NetScaler and the client.",
		netscalerLabels,
		nil,
	)

	tcpCurrentServerConnections = prometheus.NewDesc(
		"tcp_current_server_connections",
		"Server connections, including connections in the Opening, Established, and Closing state.",
		netscalerLabels,
		nil,
	)

	tcpCurrentServerConnectionsEstablished = prometheus.NewDesc(
		"tcp_current_server_connections_established",
		"Current server connections in the Established state, which indicates that data transfer can occur between the NetScaler and the server.",
		netscalerLabels,
		nil,
	)
)
