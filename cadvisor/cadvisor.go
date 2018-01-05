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
	"strings"
	"syscall"
	"time"

	"github.com/google/cadvisor/cache/memory"
	"github.com/google/cadvisor/container"
	cadvisorhttp "github.com/google/cadvisor/http"
	"github.com/google/cadvisor/manager"
	"github.com/google/cadvisor/storage"
	"github.com/google/cadvisor/utils/sysfs"
	"github.com/google/cadvisor/version"

	"crypto/tls"
)

//var argIp = flag.String("listen_ip", "", "IP to listen on, defaults to all IPs")
//var argPort = flag.Int("port", 8080, "port to listen")

// var maxProcs = flag.Int("max_procs", 0, "max number of CPUs that can be used simultaneously. Less than 1 for default (number of cores).")
// var versionFlag = flag.Bool("version", false, "print cAdvisor version and exit")
// var httpAuthFile = flag.String("http_auth_file", "", "HTTP auth file for the web UI")
// var httpAuthRealm = flag.String("http_auth_realm", "localhost", "HTTP auth realm for the web UI")
// var httpDigestFile = flag.String("http_digest_file", "", "HTTP digest file for the web UI")
// var httpDigestRealm = flag.String("http_digest_realm", "localhost", "HTTP digest file for the web UI")
// var prometheusEndpoint = flag.String("prometheus_endpoint", "/metrics", "Endpoint to expose Prometheus metrics on")
// var maxHousekeepingInterval = flag.Duration("max_housekeeping_interval", 60*time.Second, "Largest interval to allow between container housekeepings")
// var allowDynamicHousekeeping = flag.Bool("allow_dynamic_housekeeping", true, "Whether to allow the housekeeping interval to be dynamic")
// var enableProfiling = flag.Bool("profiling", false, "Enable profiling via web interface host:port/debug/pprof/")
// var collectorCert = flag.String("collector_cert", "", "Collector's certificate, exposed to endpoints for certificate based authentication.")
// var collectorKey = flag.String("collector_key", "", "Key for the collector's certificate")

var(
	maxProcs = new(int)
	versionFlag = new(bool)
	httpAuthFile = new(string)
	httpAuthRealm = new(string)
	httpDigestFile = new(string)
	httpDigestRealm = new(string)
	prometheusEndpoint = new(string)
	maxHousekeepingInterval = new(time.Duration)
	allowDynamicHousekeeping = new(bool)
	enableProfiling = new(bool)
	collectorCert = new(string)
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

var (
	storageDriver   = flag.String("storage_driver", "", fmt.Sprintf("Storage `driver` to use. Data is always cached shortly in memory, this controls where data is pushed besides the local cache. Empty means none. Options are: <empty>, %s", strings.Join(storage.ListDrivers(), ", ")))
	storageDuration = flag.Duration("storage_duration", 2*time.Minute, "How long to keep data stored (Default: 2min).")
)

type metricSetValue struct {
	container.MetricSet
}

func (ml *metricSetValue) String() string {
	var values []string
	for metric := range ml.MetricSet {
		values = append(values, string(metric))
	}
	return strings.Join(values, ",")
}

func (ml *metricSetValue) Set(value string) error {
	ml.MetricSet = container.MetricSet{}
	if value == "" {
		return nil
	}
	for _, metric := range strings.Split(value, ",") {
		if ignoreWhitelist.Has(container.MetricKind(metric)) {
			(*ml).Add(container.MetricKind(metric))
		} else {
			return fmt.Errorf("unsupported metric %q specified in disable_metrics", metric)
		}
	}
	return nil
}

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

// NewMemoryStorage creates a memory storage with an optional backend storage option.
func NewMemoryStorage() (*memory.InMemoryCache, error) {
	backendStorage, err := storage.New(*storageDriver)
	if err != nil {
		return nil, err
	}
	if *storageDriver != "" {
		log.Printf("Using backend storage type %q", *storageDriver)
	}
	log.Printf("Caching stats in memory for %v", *storageDuration)
	return memory.New(*storageDuration, backendStorage), nil
}
