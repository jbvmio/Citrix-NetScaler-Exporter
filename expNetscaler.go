package main

import "github.com/prometheus/client_golang/prometheus"

const netscalerSubsystem = "netscaler"

var netscalerLabels = []string{
	netscalerInstance,
}

var (
	modelID = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "model_id"),
		"NetScaler model - reflects the bandwidth available; for example VPX 10 would report as 10.",
		netscalerLabels,
		nil,
	)

	mgmtCPUUsage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "mgmt_cpu_usage_pct"),
		"Current CPU utilisation for management as percentage",
		netscalerLabels,
		nil,
	)

	pktCPUUsage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "pkt_cpu_usage_pct"),
		"Current CPU utilisation for packet engines, excluding management as percentage",
		netscalerLabels,
		nil,
	)

	memUsage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "mem_usage_pct"),
		"Current memory utilisation as percentage",
		netscalerLabels,
		nil,
	)

	flashPartitionUsage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "flash_partition_usage_pct"),
		"Used space in /flash partition of the disk, as a percentage.",
		netscalerLabels,
		nil,
	)

	varPartitionUsage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "var_partition_usage_pct"),
		"Used space in /var partition of the disk, as a percentage. ",
		netscalerLabels,
		nil,
	)

	totRxMB = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "total_received_mb"),
		"Total number of Megabytes received by the NetScaler appliance",
		netscalerLabels,
		nil,
	)

	totTxMB = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "total_transmit_mb"),
		"Total number of Megabytes transmitted by the NetScaler appliance",
		netscalerLabels,
		nil,
	)

	httpRequests = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "http_requests_total"),
		"Total number of HTTP requests received",
		netscalerLabels,
		nil,
	)

	httpResponses = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "http_responses_total"),
		"Total number of HTTP responses sent",
		netscalerLabels,
		nil,
	)

	tcpCurrentClientConnections = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "tcp_current_client_connections"),
		"Client connections, including connections in the Opening, Established, and Closing state.",
		netscalerLabels,
		nil,
	)

	tcpCurrentClientConnectionsEstablished = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "tcp_current_client_connections_established"),
		"Current client connections in the Established state, which indicates that data transfer can occur between the NetScaler and the client.",
		netscalerLabels,
		nil,
	)

	tcpCurrentServerConnections = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "tcp_current_server_connections"),
		"Server connections, including connections in the Opening, Established, and Closing state.",
		netscalerLabels,
		nil,
	)

	tcpCurrentServerConnectionsEstablished = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, netscalerSubsystem, "tcp_current_server_connections_established"),
		"Current server connections in the Established state, which indicates that data transfer can occur between the NetScaler and the server.",
		netscalerLabels,
		nil,
	)
)
