package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cristalhq/hedgedhttp"
	"gopkg.in/yaml.v2"
)

// Config holds the configuration values
type Config struct {
	Timeout        time.Duration `yaml:"timeout"`
	HedgedRequests int           `yaml:"hedged_requests"`
	Mirrors        []string      `yaml:"mirrors"`
}

// LoadConfig reads configuration from a YAML file
func LoadConfig(file string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// Custom client that logs requests and responses
type loggingClient struct {
	client *http.Client
}

func (lc *loggingClient) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()
	log.Printf("Sending request to: %s", req.URL)
	resp, err := lc.client.Do(req)
	if err != nil {
		log.Printf("Request to %s failed: %v", req.URL, err)
		return nil, err
	}
	duration := time.Since(start)
	log.Printf("Received response from %s - Duration: %v", req.URL, duration)
	return resp, nil
}

type MirrorResult struct {
	URL   string
	Time  time.Duration
	Size  int64
	Speed float64 // Bytes per second
	Error error
}

func testMirror(url string, client *loggingClient, wg *sync.WaitGroup, results chan<- MirrorResult) {
	defer wg.Done()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		results <- MirrorResult{URL: url, Error: fmt.Errorf("error creating request: %v", err)}
		return
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		results <- MirrorResult{URL: url, Error: fmt.Errorf("error making request: %v", err)}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		results <- MirrorResult{URL: url, Error: fmt.Errorf("error reading response: %v", err)}
		return
	}
	duration := time.Since(start)
	size := int64(len(body))
	speed := float64(size) / duration.Seconds()

	results <- MirrorResult{
		URL:   url,
		Time:  duration,
		Size:  size,
		Speed: speed,
	}
}

func main() {
	// Load configuration
	configFile := "config.yaml"
	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize logger
	log.SetOutput(os.Stdout)

	// Set up the underlying HTTP client
	baseClient := &http.Client{Timeout: 40 * time.Second}

	// Create the hedgedhttp client with a hedge delay and maximum attempts
	hedgedClient, err := hedgedhttp.NewClient(config.Timeout, config.HedgedRequests, baseClient)
	if err != nil {
		panic(err)
	}

	// Wrap the hedgedhttp client with the logging client
	loggedClient := &loggingClient{client: hedgedClient}

	var wg sync.WaitGroup
	results := make(chan MirrorResult, len(config.Mirrors))

	// Test each mirror concurrently using the loggedClient
	for _, url := range config.Mirrors {
		wg.Add(1)
		go testMirror(url, loggedClient, &wg, results)
	}

	wg.Wait()
	close(results)

	// Summarize the results
	fmt.Println("-------------------------------------------")
	fmt.Println("Mirror Speed Test Results:")
	fmt.Println("-------------------------------------------")
	var fastest, slowest *MirrorResult
	var totalSpeed float64
	for result := range results {
		if result.Error != nil {
			fmt.Printf("URL: %s - Error: %v\n", result.URL, result.Error)
			continue
		}

		fmt.Printf("URL: %s - Time: %v - Size: %d bytes - Speed: %.2f bytes/sec\n",
			result.URL, result.Time, result.Size, result.Speed)

		if fastest == nil || result.Speed > fastest.Speed {
			fastest = &result
		}
		if slowest == nil || result.Speed < slowest.Speed {
			slowest = &result
		}
		totalSpeed += result.Speed
	}

	fmt.Println("-------------------------------------------")
	fmt.Println("Summary:")
	fmt.Println("-------------------------------------------")
	if fastest != nil {
		fmt.Printf("Fastest Mirror: %s - Speed: %.2f bytes/sec\n", fastest.URL, fastest.Speed)
	}
	if slowest != nil {
		fmt.Printf("Slowest Mirror: %s - Speed: %.2f bytes/sec\n", slowest.URL, slowest.Speed)
	}
	fmt.Printf("Average Speed: %.2f bytes/sec\n", totalSpeed/float64(len(config.Mirrors)))
}
