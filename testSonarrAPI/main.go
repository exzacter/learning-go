package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// this struct holds sonarr connection details
type Config struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// wrapping http client with sonarr specific config
type SonarrClient struct {
	httpClient *http.Client
	config     Config
}

func NewSonarrClient(cfg Config) *SonarrClient {
	return &SonarrClient{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		config: cfg,
	}
}

// create a function for getting all series within sonarr
func (c *SonarrClient) GetAllSeries(ctx context.Context) (json.RawMessage, error) {
	// api endpoint for all series with a get request
	// %s is the baseurl which is being specified within the Config struct or the environment variable
	url := fmt.Sprintf("%s/api/v3/series", c.config.BaseURL)

	// logging for easier debug
	log.Printf("[CURATARR] Making a request to: %s", url)

	// create request with context for Timeout
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("[CURATARR] Failed to create request to: %v", err)
		return nil, fmt.Errorf("Creating request: %w", err)
	}

	// using x-api-key as a header keeps key out of logs and browser history
	req.Header.Set("X-Api-Key", c.config.APIKey)
	req.Header.Set("Accept", "application/json")

	log.Printf("[CURATARR] Sending GET request")

	// execute the request to the api endpoint and time the request to output if it fails
	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime)
	if err != nil {
		log.Printf("[CURATARR] Sending GET request failed after %v: %v", duration, err)
		return nil, fmt.Errorf("Executing request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("[CURATARR] Response received - Status %d %s - Duration: %v", resp.StatusCode, http.StatusText(resp.StatusCode), duration)

	// check for non succes status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[CURATARR] Error in the response body: %s", string(body))
		return nil, fmt.Errorf("Sonarr returned status %d: %s", resp.StatusCode, string(body))
	}

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[CURATARR] Failed to read response body: %v", err)
		return nil, fmt.Errorf("Reading response: %w", err)
	}

	log.Printf("[CURATARR] Successfully fetched %d bytes", len(body))

	return body, nil
}

// load environment variables from .env
func LoadConfig() Config {

	baseURL := os.Getenv("SONARR_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8989" // default sonarr installation port
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("Sonarr API key not found. API key required to pull tv shows")
	}

	return Config{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Timeout: 30 * time.Second,
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	http.HandleFunc("/", handleIndex)

	http.HandleFunc("/api/fetch-series", handleFetchSeries)

	log.Println("Sever running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html := `<!DOCTYPE html> 
<html>
<head>
    <title>Sonarr API</title>
</head>
<body>
    <h1>Sonarr Series</h1>
    <button onclick="fetchSeries()">Fetch Series</button>
    <pre id="output"></pre>
    <script>
        async function fetchSeries() {
            const output = document.getElementById('output');
            output.textContent = 'Loading...';
            try {
                const response = await fetch('/api/fetch-series');
                const data = await response.json();
            
                const seriesCount = data.length;
                console.log('[STATS] Total series:', seriesCount);
                
                output.textContent = 'Total Series: ' + seriesCount + '\n\n' + JSON.stringify(data, null, 2);
            } catch (error) {
                output.textContent = 'Error: ' + error.message;
            }
        }
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func handleFetchSeries(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /api/v1/sync-series - Syncing your Series' from Sonarr into Curatarr")

	// set cors
	// CORS is ...
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Origin", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "Content-Type")

	// handle preflight?
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// only allow get for fetching handleFetchSeries
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// load config and create sonarr client
	cfg := LoadConfig()
	client := NewSonarrClient(cfg)

	// create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// the actual fetch from sonarr
	seriesJSON, err := client.GetAllSeries(ctx)
	if err != nil {
		log.Printf("Error fetching series: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch series: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(seriesJSON)
}
