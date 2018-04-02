package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type SearchBody struct {
	Took		int				`json:"took"`
	TimedOut	bool			`json:"timed_out"`
	Shards		interface{}		`json:"_shards"`
	Hits		HitBody			`json:"hits"`
}

type HitBody struct {
	Total		int				`json:"total"`
	MaxScore	float32			`json:"max_score"`
	Hits		[]HitContent	`json:"hits"`
}

type HitContent struct {
	Index		string			`json:"_index"`
	Type		string			`json:"_type"`
	Id			string			`json:"_id"`
	Score		float32			`json:"_score"`
	Source		Source			`json:"_source"`
	Sort		[]uint64		`json:"sort"`
}

type Source struct {
	Log			string			`json:"log"`
	Stream		string			`json:"stream"`
	Docker		interface{}		`json:"docker"`
	Kubernetes	Kubernetes		`json:"kubernetes"`
	Timestamp	string			`json:"@timestamp"`
	Tag			string			`json:"tag"`
}

type Kubernetes struct {
	ContainerName	string		`json:"container_name"`
	NamespaceName	string		`json:"namespace_name"`
	PodName			string		`json:"pod_name"`
	PodId			string		`json:"pod_id"`
	Labels			interface{}	`json:"labels"`
	Host			string		`json:"host"`
	MasterURL		string		`json:"master_url"`
	NamespaceId		string		`json:"namespace_id"`
}

func postQuery(url string, fabricNamespace string, podName string) (*string, error) {
	efkQueryString := "(\"ProcessProposal -> DEBU \" AND (Entry OR Exit)) OR NewCCCC OR generateDockerfile"

	var mustQuery [] interface{}
	// query_string for log
	queryString := make(map[string]interface{})
	queryString["query"] = efkQueryString
	queryStringBody := make(map[string]interface{})
	queryStringBody["query_string"] = queryString
	mustQuery = append(mustQuery, queryStringBody)
	// query k8s namespace
	namespace := make(map[string]interface{})
	namespace["query"] = fabricNamespace
	matchNamespace := make(map[string]interface{})
	matchNamespace["kubernetes.namespace_name"] = namespace
	matchNamespaceBody := make(map[string]interface{})
	matchNamespaceBody["match_phrase"] = matchNamespace
	mustQuery = append(mustQuery, matchNamespaceBody)
	// query k8s pod name
	containerName := make(map[string]interface{})
	containerName["query"] = podName
	matchContainerName := make(map[string]interface{})
	matchContainerName["kubernetes.container_name"] = containerName
	matchContainerNameBody := make(map[string]interface{})
	matchContainerNameBody["match_phrase"] = matchContainerName
	mustQuery = append(mustQuery, matchContainerNameBody)
	// encapsulate to full json object
	mustBody := make(map[string]interface{})
	mustBody["must"] = mustQuery
	boolBody := make(map[string]interface{})
	boolBody["bool"] = mustBody
	queryBody := make(map[string]interface{})
	queryBody["query"] = boolBody

	// Marshal request to json, and send to server
	byteData, err := json.Marshal(queryBody)
	if err != nil {
		log.Printf("Error cannot marshal post body: %s", err)
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(byteData))
	if err != nil {
		log.Printf("Error construct request body: %s", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error post request to server: %s", err)
		return nil, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error read respons from server: %s", err)
		return nil, err
	}

	resString := string(resBody)
	return &resString, nil
}

// Make a Regex to delete all non-printable characters
func removeNonPrintable(originalString string) (*string, error) {
	reg, err := regexp.Compile("[\x00\x08\x0B\x0C\x0E-\x1F]")
	if err != nil {
		log.Printf("Error cannot compile regex: %s", err)
		return nil, err
	}
	printableString := reg.ReplaceAllString(originalString, "")
	return &printableString, nil
}

func extractLogs(searchData *string) (*[]HitContent, error) {
	// Delete all non-printable characters
	pureString, err := removeNonPrintable(*searchData)
	if err != nil {
		log.Printf("Error remove non printable characters: %s", err)
		return nil, err
	}

	// Here needs reply struct judgements

	// Extract string to a new json
	var searchBody SearchBody
	json.Unmarshal([]byte(*pureString), &searchBody)
	return &searchBody.Hits.Hits, nil
}

func main() {
	addr := flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	interval := flag.Uint64("interval", 60, "The interval (seconds) to grab data from fabric net.")
	elaSearchAddr := flag.String("elaSearchAddr", "127.0.0.1", "The address of elasticsearch.")
	elaSearchPort := flag.String("elaSearchPort", "3000", "The port of elasticsearch.")
	fabricNamespace := flag.String("fabricNamespace", "fabric-net", "The namespaces of fabric net.")
	containerName := flag.String("containerName", "peer0-org1", "The container name to grab data")

	go func() {
		for {
			elaSearchURL := "http://" + *elaSearchAddr + ":" + *elaSearchPort
			log.Printf("Info get logs from elasticsearch server: %s", elaSearchURL)
			res, err := postQuery(elaSearchURL, *fabricNamespace, *containerName)
			if err != nil {
				log.Printf("Error cannot query logs from ElasticSearch: %s", err.Error())
				return
			}
			logsArray, err := extractLogs(res)
			log.Printf("Info start to analysis logs")
			analysisLogs(logsArray)
			log.Printf("Info sleep for %d seconds", *interval)
			time.Sleep(time.Duration(*interval) * time.Second)
		}
	}()

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
