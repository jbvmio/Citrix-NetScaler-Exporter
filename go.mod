module github.com/jbvmio/Citrix-NetScaler-Exporter

go 1.12

require (
	github.com/beorn7/perks v0.0.0-20160804104726-4c0e84591b9a // indirect
	github.com/go-kit/kit v0.6.0
	github.com/go-logfmt/logfmt v0.3.0 // indirect
	github.com/go-stack/stack v1.6.0 // indirect
	github.com/golang/protobuf v0.0.0-20170920220647-130e6b02ab05 // indirect
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.0 // indirect
	github.com/pkg/errors v0.8.0
	github.com/prometheus/client_golang v0.8.0
	github.com/prometheus/client_model v0.0.0-20170216185247-6f3806018612 // indirect
	github.com/prometheus/common v0.0.0-20171006141418-1bab55dd05db // indirect
	github.com/prometheus/procfs v0.0.0-20170703101242-e645f4e5aaa8 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
)

replace github.com/jbvmio/citrix-netscaler-exporter/collector => ./collector

replace github.com/jbvmio/citrix-netscaler-exporter/netscaler => ./netscaler
