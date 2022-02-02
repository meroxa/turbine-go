package clearbit

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/sling"
)

// Client is a Clearbit client for making Clearbit API requests.
type Client struct {
	Autocomplete *AutocompleteService
	Person       *PersonService
	Company      *CompanyService
	Discovery    *DiscoveryService
	Prospector   *ProspectorService
	Reveal       *RevealService
	Risk         *RiskService
	NameToDomain *NameToDomainService
}

// config represents all the parameters available to configure a Clearbit
// client
type config struct {
	apiKey     string
	httpClient *http.Client
	timeout    time.Duration
	baseURLs   *BaseURLs
}

// Option is an option passed to the NewClient function used to change
// the client configuration
type Option func(*config)

// WithHTTPClient sets the optional http.Client we can use to make requests
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *config) {
		c.httpClient = httpClient
	}
}

// WithAPIKey sets the Clearbit API key.
//
// When this is not provided we'll default to the `CLEARBIT_KEY` environment
// variable.
func WithAPIKey(apiKey string) func(*config) {
	return func(c *config) {
		c.apiKey = apiKey
	}
}

// WithTimeout sets the http timeout
//
// This is just an easier way to set the timeout than directly setting it
// through the withHTTPClient option.
func WithTimeout(d time.Duration) func(*config) {
	return func(c *config) {
		c.timeout = d
	}
}

// WithBaseURL sets the base URL for API requests
//
// This allows for the mocking of the Clearbit service when writing
// tests against the clearbit client
func WithBaseURLs(urls map[string]string) func(*config) {
	return func(c *config) {
		c.baseURLs = NewBaseURLs(urls)
	}
}

type BaseURLs struct {
	Autocomplete string
	Person       string
	Company      string
	Discovery    string
	Prospector   string
	Reveal       string
	Risk         string
	NameToDomain string
}

func NewBaseURLs(overrideURLs map[string]string) *BaseURLs {
	// default base URLs
	baseURLs := &BaseURLs{
		Autocomplete: "https://autocomplete.clearbit.com",
		Person:       "https://person.clearbit.com",
		Company:      "https://company.clearbit.com",
		Discovery:    "https://discovery.clearbit.com",
		Prospector:   "https://prospector.clearbit.com",
		Reveal:       "https://reveal.clearbit.com",
		Risk:         "https://risk.clearbit.com",
		NameToDomain: "https://company.clearbit.com",
	}

	jsonOverrideURLs, err := json.Marshal(overrideURLs)
	if err != nil {
		return baseURLs
	}

	if err := json.Unmarshal(jsonOverrideURLs, &baseURLs); err != nil {
		return baseURLs
	}
	return baseURLs
}

// NewClient returns a new Client.
func NewClient(options ...Option) *Client {
	c := config{
		apiKey:     os.Getenv("CLEARBIT_KEY"),
		httpClient: &http.Client{},
		timeout:    10 * time.Second,
		baseURLs:   NewBaseURLs(map[string]string{}),
	}

	for _, option := range options {
		option(&c)
	}

	c.httpClient.Timeout = c.timeout

	base := sling.New().Client(c.httpClient)
	base.SetBasicAuth(c.apiKey, "")

	return &Client{
		Autocomplete: newAutocompleteService(base.New(), c.baseURLs.Autocomplete),
		Person:       newPersonService(base.New(), c.baseURLs.Person),
		Company:      newCompanyService(base.New(), c.baseURLs.Company),
		Discovery:    newDiscoveryService(base.New(), c.baseURLs.Discovery),
		Prospector:   newProspectorService(base.New(), c.baseURLs.Prospector),
		Reveal:       newRevealService(base.New(), c.baseURLs.Reveal),
		Risk:         newRiskService(base.New(), c.baseURLs.Risk),
		NameToDomain: newNameToDomainService(base.New(), c.baseURLs.NameToDomain),
	}
}
