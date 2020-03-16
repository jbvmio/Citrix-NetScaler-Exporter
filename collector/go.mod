module github.com/jbvmio/citrix-netscaler-exporter/collector

go 1.12

replace github.com/jbvmio/citrix-netscaler-exporter => ../

require (
	github.com/go-kit/kit v0.10.0
	github.com/prometheus/client_golang v1.5.1
)
