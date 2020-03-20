package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var (
	app          = "Citrix-NetScaler-Exporter"
	version      string
	build        string
	username     = flag.String("username", "", "Username with which to connect to the NetScaler API")
	password     = flag.String("password", "", "Password with which to connect to the NetScaler API")
	localMapping = flag.String("mapping", "./mappings.yaml", "Load local mappings file")
	bindPort     = flag.Int("bind_port", 9280, "Port to bind the exporter endpoint to")
	versionFlg   = flag.Bool("version", false, "Display application version")
	debugFlg     = flag.Bool("debug", false, "Enable debug logging?")
	logger       log.Logger
	nsInstance   string
	vipDB        *DB
)

func init() {
	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller, "app", app, "bind_port", *bindPort, "version", "v"+version, "build", build)
}

func main() {
	flag.Parse()

	if *versionFlg {
		fmt.Printf("%s v%s build %s\n", app, version, build)
		os.Exit(0)
	}

	if *username == "" || *password == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	currentMapping = VIPMap{
		mappings: make(map[string]map[string]string),
		lock:     sync.Mutex{},
	}

	currentMapping.loadMappingYaml(*localMapping)

	vipDB = newDB(dbDir)
	vipDB.loadVIPMap(&currentMapping)
	go vipDB.collectAll()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
				<head><title>Citrix NetScaler Exporter</title></head>
				<style>
				label{
				display:inline-block;
				width:75px;
				}
				form label {
				margin: 10px;
				}
				form input {
				margin: 10px;
				}
				</style>
				<body>
				<h1>Citrix NetScaler Exporter</h1>
				<form action="/netscaler">
				<label>Target:</label> <input type="text" name="target" placeholder="https://netscaler.domain.tld"> <br>
				<p>Ignore certificate check?</p>
				<input type="radio" id="yes" name="ignore-cert" value="yes">
				<label for="yes">Yes</label>
				<input type="radio" id="no" name="ignore-cert" value="no" checked>
  				<label for="no">No</label>
				<br>
				<input type="submit" value="Submit">
				</form>
				</body>
				</html>`))
	})

	http.HandleFunc("/netscaler", handler)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/mapping", handleMapping)

	listeningPort := ":" + strconv.Itoa(*bindPort)
	level.Info(logger).Log("msg", "Listening on port "+listeningPort)

	err := http.ListenAndServe(listeningPort, nil)
	vipDB.stopCollect()
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "'target' parameter must be specified", 400)
		return
	}

	ignoreCertCheck := false
	if strings.ToLower(r.URL.Query().Get("ignore-cert")) == "yes" {
		ignoreCertCheck = true
	}

	nsInstance = strings.TrimLeft(target, "https://")
	nsInstance = strings.TrimLeft(nsInstance, "http://")
	nsInstance = strings.Trim(nsInstance, " /")

	if *debugFlg {
		level.Debug(logger).Log("msg", "scraping target", "target", target)
	}

	there, ready := vipDB.exists(target)
	loaded := currentMapping.exists(target)
	switch {
	case !there:
		level.Info(logger).Log("msg", "creating new vip mappings for "+target)
		lbs := lbserver{
			url:    target,
			user:   *username,
			pass:   *password,
			ignore: ignoreCertCheck,
		}
		if loaded {
			vipDB.setLBServer(lbs)
		} else {
			vipDB.setLBServer(lbs)
			err := vipDB.collectVIPMap2(lbs)
			if err != nil {
				level.Error(logger).Log("msg", "error creating new vip mappings: "+err.Error())
				vipDB.removeLBServer(lbs)
				return
			}
		}
	case !ready:
		if !loaded {
			w.WriteHeader(http.StatusOK)
			level.Info(logger).Log("msg", "vip mappings not ready yet for "+target)
			return
		}
	}

	exporter, err := NewExporter(target, *username, *password, ignoreCertCheck, logger, nsInstance)
	if err != nil {
		http.Error(w, "Error creating exporter"+err.Error(), 400)
		level.Error(logger).Log("msg", err)
		return
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	// Delegate http serving to Prometheus client library, which will call Collect.
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func handleMapping(w http.ResponseWriter, r *http.Request) {
	maps, err := currentMapping.getMappingYaml()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(maps)
}
