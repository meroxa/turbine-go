package clearbit

import (
	"net/http"

	"github.com/dghubble/sling"
)

// Risk represents the risk score returned by a call to Calculate
type Risk struct {
	ID    string `json:"id"`
	Live  bool   `json:"live"`
	Email struct {
		Valid        bool `json:"valid"`
		SocialMatch  bool `json:"socialMatch"`
		CompanyMatch bool `json:"companyMatch"`
		NameMatch    bool `json:"nameMatch"`
		Disposable   bool `json:"disposable"`
		FreeProvider bool `json:"freeProvider"`
		Blacklisted  bool `json:"blacklisted"`
	} `json:"email"`
	Address struct {
		GeoMatch interface{} `json:"geoMatch"`
	} `json:"address"`
	IP struct {
		Proxy       bool        `json:"proxy"`
		GeoMatch    interface{} `json:"geoMatch"`
		Blacklisted bool        `json:"blacklisted"`
		RateLimited interface{} `json:"rateLimited"`
	} `json:"ip"`
	Risk struct {
		Level   string   `json:"level"`
		Score   int      `json:"score"`
		Reasons []string `json:"reasons"`
	} `json:"risk"`
}

// RiskCalculateParams wraps the parameters needed to interact with the Risk API
// through the Calculate method
type RiskCalculateParams struct {
	Email       string `url:"email,omitempty"`
	IP          string `url:"ip,omitempty"`
	CountryCode string `url:"country_code,omitempty"`
	ZipCode     string `url:"zip_code,omitempty"`
	GivenName   string `url:"given_name,omitempty"`
	FamilyName  string `url:"family_name,omitempty"`
	Name        string `url:"name,omitempty"`
}

// RiskService gives access to the Risk API.
//
// Our Risk API takes an email address, an IP address, and additional information
// before returning a risk analysis for the user
type RiskService struct {
	baseSling *sling.Sling
	sling     *sling.Sling
}

func newRiskService(sling *sling.Sling, baseURL string) *RiskService {
	return &RiskService{
		baseSling: sling.New(),
		sling:     sling.Base(baseURL).Path("/v1/"),
	}
}

// Find takes an email address, and an IP address, and returns the risk associated
// with that user
func (s *RiskService) Calculate(params RiskCalculateParams) (*Risk, *http.Response, error) {
	item := new(Risk)
	ae := new(apiError)
	resp, err := s.sling.New().Post("calculate").QueryStruct(params).Receive(item, ae)
	return item, resp, relevantError(err, *ae)
}
