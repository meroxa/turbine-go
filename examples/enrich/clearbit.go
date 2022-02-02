package main

import (
	"github.com/clearbit/clearbit-go/clearbit"
	"log"
	"os"
)

type UserDetails struct {
	FullName        string
	Location        string
	Role            string
	Seniority       string
	Company         string
	GithubUser      string
	GithubFollowers int
}

func EnrichUserEmail(email string) (*UserDetails, error) {
	key := os.Getenv("CLEARBIT_API_KEY")
	client := clearbit.NewClient(clearbit.WithAPIKey(key))
	results, resp, err := client.Person.FindCombined(clearbit.PersonFindParams{
		Email: email,
	})

	if err != nil {
		log.Printf("error looking up email; resp: %%+v", resp.Status)
		return nil, err
	}

	return &UserDetails{
		FullName:        results.Person.Name.FullName,
		Location:        results.Person.Location,
		Role:            results.Person.Employment.Role,
		Seniority:       results.Person.Employment.Seniority,
		Company:         results.Company.Name,
		GithubUser:      results.Person.GitHub.Handle,
		GithubFollowers: results.Person.GitHub.Followers,
	}, nil

}
