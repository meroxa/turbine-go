package clearbit

import (
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

// Company contains all the company fields gathered from the Company json
// structure. https://dashboard.clearbit.com/docs#enrichment-api-company-api
type Company struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	LegalName     string   `json:"legalName"`
	Domain        string   `json:"domain"`
	DomainAliases []string `json:"domainAliases"`
	Site          struct {
		PhoneNumbers   []string `json:"phoneNumbers"`
		EmailAddresses []string `json:"emailAddresses"`
	} `json:"site"`
	Category struct {
		Sector        string `json:"sector"`
		IndustryGroup string `json:"industryGroup"`
		Industry      string `json:"industry"`
		SubIndustry   string `json:"subIndustry"`
		SicCode       string `json:"sicCode"`
		NaicsCode     string `json:"naicsCode"`
	} `json:"category"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	FoundedYear int      `json:"foundedYear"`
	Location    string   `json:"location"`
	TimeZone    string   `json:"timeZone"`
	UtcOffset   int      `json:"utcOffset"`
	Geo         struct {
		StreetNumber string  `json:"streetNumber"`
		StreetName   string  `json:"streetName"`
		SubPremise   string  `json:"subPremise"`
		City         string  `json:"city"`
		PostalCode   string  `json:"postalCode"`
		State        string  `json:"state"`
		StateCode    string  `json:"stateCode"`
		Country      string  `json:"country"`
		CountryCode  string  `json:"countryCode"`
		Lat          float64 `json:"lat"`
		Lng          float64 `json:"lng"`
	} `json:"geo"`
	Logo     string `json:"logo"`
	Facebook struct {
		Handle string `json:"handle"`
		Likes  int    `json:"likes"`
	} `json:"facebook"`
	LinkedIn struct {
		Handle string `json:"handle"`
	} `json:"linkedin"`
	Twitter struct {
		Handle    string `json:"handle"`
		ID        string `json:"id"`
		Bio       string `json:"bio"`
		Followers int    `json:"followers"`
		Following int    `json:"following"`
		Location  string `json:"location"`
		Site      string `json:"site"`
		Avatar    string `json:"avatar"`
	} `json:"twitter"`
	Crunchbase struct {
		Handle string `json:"handle"`
	} `json:"crunchbase"`
	EmailProvider bool   `json:"emailProvider"`
	Type          string `json:"type"`
	Ticker        string `json:"ticker"`
	Identifiers   struct {
		UsEIN string `json:"usEIN"`
	} `json:"identifiers"`
	Phone   string `json:"phone"`
	Metrics struct {
		AlexaUsRank            int    `json:"alexaUsRank"`
		AlexaGlobalRank        int    `json:"alexaGlobalRank"`
		Employees              int    `json:"employees"`
		EmployeesRange         string `json:"employeesRange"`
		MarketCap              int    `json:"marketCap"`
		Raised                 int    `json:"raised"`
		AnnualRevenue          int    `json:"annualRevenue"`
		EstimatedAnnualRevenue string `json:"estimatedAnnualRevenue"`
		FiscalYearEnd          int    `json:"fiscalYearEnd"`
	} `json:"metrics"`
	IndexedAt time.Time `json:"indexedAt"`
	Tech      []string  `json:"tech"`
	Parent    struct {
		Domain string `json:"domain"`
	} `json:"parent"`
}

// CompanyFindParams wraps the parameters needed to interact with the Company
// API through the Find method
type CompanyFindParams struct {
	Domain string `url:"domain,omitempty"`
}

// CompanyService gives access to the Company API.
// https://dashboard.clearbit.com/docs#enrichment-api-company-api
type CompanyService struct {
	baseSling *sling.Sling
	sling     *sling.Sling
}

func newCompanyService(sling *sling.Sling, baseURL string) *CompanyService {
	return &CompanyService{
		baseSling: sling.New(),
		sling:     sling.Base(baseURL).Path("/v2/companies/"),
	}
}

//Find looks up a company based on its domain
func (s *CompanyService) Find(params CompanyFindParams) (*Company, *http.Response, error) {
	item := new(Company)
	ae := new(apiError)
	resp, err := s.sling.New().Get("find").QueryStruct(params).Receive(item, ae)
	return item, resp, relevantError(err, *ae)
}
