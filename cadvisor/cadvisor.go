package cadvisor

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"crypto/tls"

	"github.com/google/cadvisor/container"
	cadvisorhttp "github.com/google/cadvisor/http"
	"github.com/google/cadvisor/manager"
	"github.com/google/cadvisor/utils/sysfs"
	"github.com/google/cadvisor/version"
)

var(
	// max number of CPUs that can be used simultaneously. Less than 1 for default (number of cores).
	maxProcs = new(int)
	// print cAdvisor version and exit
	versionFlag = new(bool)
	// HTTP auth file for the web UI
	httpAuthFile = new(string)
	// HTTP auth realm for the web UI
	httpAuthRealm = new(string)
	// HTTP digest file for the web UI
	httpDigestFile = new(string)
	// HTTP digest realm for the web UI
	httpDigestRealm = new(string)
	// Endpoint to expose Prometheus metrics on
	prometheusEndpoint = new(string)
	// Largest interval to allow between container housekeepings
	maxHousekeepingInterval = new(time.Duration)
	// Whether to allow the housekeeping interval to be dynamic
	allowDynamicHousekeeping = new(bool)
	// Enable profiling via web interface host:port/debug/pprof/
	enableProfiling = new(bool)
	// Collector's certificate, exposed to endpoints for certificate based authentication.
	collectorCert = new(string)
	// Key for the collector's certificate
	collectorKey = new(string)
)

var (
	// Metrics to be ignored.
	// Tcp metrics are ignored by default.
	ignoreMetrics metricSetValue = metricSetValue{container.MetricSet{
		container.NetworkTcpUsageMetrics: struct{}{},
		container.NetworkUdpUsageMetrics: struct{}{},
	}}

	// List of metrics that can be ignored.
	ignoreWhitelist = container.MetricSet{
		container.DiskUsageMetrics:       struct{}{},
		container.NetworkUsageMetrics:    struct{}{},
		container.NetworkTcpUsageMetrics: struct{}{},
		container.NetworkUdpUsageMetrics: struct{}{},
	}
)

func init() {
	flag.Var(&ignoreMetrics, "disable_metrics", "comma-separated list of `metrics` to be disabled. Options are 'disk', 'network', 'tcp', 'udp'. Note: tcp and udp are disabled by default due to high CPU usage.")

	// Default logging verbosity to V(2)
	flag.Set("v", "2")
}

func New(argIp *string, argPort *int) {
	*maxProcs = 0
	*versionFlag = false
	*httpAuthFile = ""
	*httpAuthRealm = "localhost"
	*httpDigestFile = ""
	*httpDigestRealm = "localhost"
	*prometheusEndpoint = "/metrics"
	*maxHousekeepingInterval = 60*time.Second
	*allowDynamicHousekeeping = true
	*enableProfiling = false
	*collectorCert = ""
	*collectorKey = ""

	// Just collect container information
	flag.Set("docker_only", "true")

	if *versionFlag {
		fmt.Printf("cAdvisor version %s (%s)\n", version.Info["version"], version.Info["revision"])
		os.Exit(0)
	}

	setMaxProcs()

	memoryStorage, err := NewMemoryStorage()
	if err != nil {
		log.Printf("Failed to initialize storage driver: %s", err)
	}

	sysFs := sysfs.NewRealSysFs()

	collectorHttpClient := createCollectorHttpClient(*collectorCert, *collectorKey)

	containerManager, err := manager.New(memoryStorage, sysFs, *maxHousekeepingInterval, *allowDynamicHousekeeping, ignoreMetrics.MetricSet, &collectorHttpClient)
	if err != nil {
		log.Printf("Failed to create a Container Manager: %s", err)
	}

	mux := http.NewServeMux()

	if *enableProfiling {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	}

	// Register all HTTP handlers.
	err = cadvisorhttp.RegisterHandlers(mux, containerManager, *httpAuthFile, *httpAuthRealm, *httpDigestFile, *httpDigestRealm)
	if err != nil {
		log.Printf("Failed to register HTTP handlers: %v", err)
	}

	cadvisorhttp.RegisterPrometheusHandler(mux, containerManager, *prometheusEndpoint, nil)

	// Start the manager.
	if err := containerManager.Start(); err != nil {
		log.Printf("Failed to start container manager: %v", err)
	}

	// Install signal handler.
	installSignalHandler(containerManager)

	log.Printf("Starting cAdvisor version: %s-%s on port %d", version.Info["version"], version.Info["revision"], *argPort)

	addr := fmt.Sprintf("%s:%d", *argIp, *argPort)
	http.ListenAndServe(addr, mux)
}

func setMaxProcs() {
	// TODO(vmarmol): Consider limiting if we have a CPU mask in effect.
	// Allow as many threads as we have cores unless the user specified a value.
	var numProcs int
	if *maxProcs < 1 {
		numProcs = runtime.NumCPU()
	} else {
		numProcs = *maxProcs
	}
	runtime.GOMAXPROCS(numProcs)

	// Check if the setting was successful.
	actualNumProcs := runtime.GOMAXPROCS(0)
	if actualNumProcs != numProcs {
		log.Printf("Specified max procs of %v but using %v", numProcs, actualNumProcs)
	}
}

func installSignalHandler(containerManager manager.Manager) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	go func() {
		sig := <-c
		if err := containerManager.Stop(); err != nil {
			log.Printf("Failed to stop container manager: %v", err)
		}
		log.Printf("Exiting given signal: %v", sig)
		os.Exit(0)
	}()
}

func createCollectorHttpClient(collectorCert, collectorKey string) http.Client {
	//Enable accessing insecure endpoints. We should be able to access metrics from any endpoint
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	if collectorCert != "" {
		if collectorKey == "" {
			log.Printf("The collector_key value must be specified if the collector_cert value is set.")
		}
		cert, err := tls.LoadX509KeyPair(collectorCert, collectorKey)
		if err != nil {
			log.Printf("Failed to use the collector certificate and key: %s", err)
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
		tlsConfig.BuildNameToCertificate()
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return http.Client{Transport: transport}
}
