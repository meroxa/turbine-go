package clearbit

import (
	"net/http"

	"github.com/dghubble/sling"
)

// NameToDomain represents the company returned by a call to Find
type NameToDomain struct {
	Logo   string `json:"logo"`
	Name   string `json:"string"`
	Domain string `json:"domain"`
}

// NameToDomainFindParams wraps the parameters needed to interact with the NameToDomain API
// through the Find method
type NameToDomainFindParams struct {
	Name string `url:"name"`
}

// NameToDomainService gives access to the NameToDomain API.
//
// Our NameToDomain API takes a company name, and returns the domain associated with
// that name.
type NameToDomainService struct {
	baseSling *sling.Sling
	sling     *sling.Sling
}

func newNameToDomainService(sling *sling.Sling, baseURL string) *NameToDomainService {
	return &NameToDomainService{
		baseSling: sling.New(),
		sling:     sling.Base(baseURL).Path("/v1/"),
	}
}

// Find takes a company name and returns the domain associated with that name
func (s *NameToDomainService) Find(params NameToDomainFindParams) (*NameToDomain, *http.Response, error) {
	item := new(NameToDomain)
	ae := new(apiError)
	resp, err := s.sling.New().Get("domains/find").QueryStruct(params).Receive(item, ae)
	return item, resp, relevantError(err, *ae)
}
