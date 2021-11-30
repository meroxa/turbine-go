package internal

import (
	"github.com/meroxa/meroxa-go/pkg/mock"
	"github.com/meroxa/valve"
)

type LocalClient struct {
	client mock.MockClient
}

func (c LocalClient) GetResource(name string) (*valve.Resource, error) {

}