package runtime

import "github.com/meroxa/meroxa-go/pkg/meroxa"

const TurbineLanguageGo = "golang"

type App struct {
	Name     string
	Language string
	GitSHA   string
	Pipeline meroxa.EntityIdentifier
}
