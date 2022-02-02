package clearbit

import (
	"net/http"

	"github.com/dghubble/sling"
)

const (
	apiVersion = "2018-08-15"
)

type ProspectorResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
	Results  []struct {
		ID   string `json:"id"`
		Name struct {
			FullName   string `json:"fullName"`
			GivenName  string `json:"givenName"`
			FamilyName string `json:"familyName"`
		} `json:"name"`
		Title     string `json:"title"`
		Role      string `json:"role"`
		Seniority string `json:"seniority"`
		Company   struct {
			Name string `json:"name"`
		} `json:"company"`
		Email    string `json:"email"`
		Location string `json:"location"`
		Phone    string `json:"phone"`
		Verified bool   `json:"verified"`
	} `json:"results"`
}

// ProspectorSearchParams wraps the parameters needed to interact with the
// Prospector API
type ProspectorSearchParams struct {
	Domain      string   `url:"domain,omitempty"`
	Role        string   `url:"role,omitempty"`
	Roles       []string `url:"roles[],omitempty"`
	Seniority   string   `url:"seniority,omitempty"`
	Seniorities []string `url:"seniorities[],omitempty"`
	Title       string   `url:"title,omitempty"`
	Titles      []string `url:"titles[],omitempty"`
	City        string   `url:"city,omitempty"`
	Cities      []string `url:"cities[],omitempty"`
	State       string   `url:"state,omitempty"`
	States      []string `url:"states[],omitempty"`
	Country     string   `url:"country,omitempty"`
	Countries   []string `url:"countries[],omitempty"`
	Name        string   `url:"name,omitempty"`
	Page        int      `url:"page,omitempty"`
	PageSize    int      `url:"page_size,omitempty"`
}

// ProspectorService gives access to the Prospector API.
//
// The Prospector API lets you fetch contacts and emails associated with a
// company, employment role, seniority, and job title.
type ProspectorService struct {
	baseSling *sling.Sling
	sling     *sling.Sling
}

func newProspectorService(sling *sling.Sling, baseURL string) *ProspectorService {
	return &ProspectorService{
		baseSling: sling.New(),
		sling:     sling.Base(baseURL).Path("/v1/people/").Set("Api-Version", apiVersion),
	}
}

// Search lets you fetch contacts and emails associated with a company,
// employment role, seniority, and job title.
func (s *ProspectorService) Search(params ProspectorSearchParams) (ProspectorResponse, *http.Response, error) {
	pr := new(ProspectorResponse)
	ae := new(apiError)
	resp, err := s.sling.New().Get("search").QueryStruct(params).Receive(pr, ae)
	return *pr, resp, relevantError(err, *ae)
}
