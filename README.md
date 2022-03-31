# Turbine

[![nightly tests](https://github.com/meroxa/turbine/actions/workflows/nightly_tests/badge.svg)](https://github.com/meroxa/turbine/actions/workflows/nightly_tests) [![Go Report Card](https://goreportcard.com/badge/github.com/meroxa/turbine)](https://goreportcard.com/report/github.com/stretchr/testify) [![PkgGoDev](https://pkg.go.dev/badge/github.com/meroxa/turbine)](https://pkg.go.dev/github.com/meroxa/turbine)

<p align="center" style="text-align:center;">
  <img alt="turbine logo" src="docs/turbine-outline.svg" width="500" />
</p>

Turbine is a data application framework for building server-side applications that are event-driven, respond to data in real-time, and scale using cloud-native best practices.

The benefits of using Turbine include:

* **Native Developer Tooling:** Turbine doesn't come with any bespoke DSL or patterns. Write software like you normally would!

* **Fits into Existing DevOps Workflows:** Build, test, and deploy. Turbine encourages best practices from the start. Don't test your data app in production ever again.

* **Local Development mirrors Production:** When running locally, you'll immediately see how your app reacts to data. What you get there will be exactly what happens in production but with _scale_ and _speed_.

* **Available in many different programming langauages:** Turbine started out in Go but is available in other languages too:
    * [Go](https://github.com/meroxa/turbine)
    * [Javascript](https://github.com/meroxa/turbine-js)
    * [Python](https://github.com/meroxa/turbine-py)


## Getting Started

To get started, you'll need to [download the Meroxa CLI](https://github.com/meroxa/cli#installation-guide). Once downloaded and installed, you'll need to back to your terminal and initialize a new project:

```bash
$ meroxa apps init testapp --lang=golang
```

The CLI will create a new folder called `testapp` located in the directory where the command was issued. Once you enter the directory, the contents of the directory will look like this:

```bash
$ tree testapp/
testapp
├── README.md
├── app.go
├── app.json
├── app_test.go
└── fixtures
    └── README.md
```

This will be a full fledged Turbine app that can run. You can even run the tests using the command `meroxa apps test` in the root of the app directory. It just enough to show you what you need to get started.


### `app.go`

This is the file were you begin your turbine journey. Any time a turbine app run, this is the entrypoint for the entire application. When the project is first created the file will look like this:

```go
package main

import (
	"github.com/meroxa/turbine"
	"github.com/meroxa/turbine/runner"
)

func main() {
	runner.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

func (a App) Run(v turbine.Turbine) error {
	source, err := v.Resources("source_name")
	if err != nil {
		return err
	}

	rr, err := source.Records("collection_name", nil)
	if err != nil {
		return err
	}

	res, _ := v.Process(rr, Anonymize{})

	dest, err := v.Resources("dest_name")
	if err != nil {
		return err
	}

	err = dest.Write(res, "collection_name", nil)
	if err != nil {
		return err
	}

	return nil
}

type Anonymize struct{}

func (f Anonymize) Process(rr []turbine.Record) ([]turbine.Record, []turbine.RecordWithError) {
	return rr, nil
}
```


### `app.json`

-- TODO walk through the important config options.

### Testing

-- TODO walk through how to do it.

## Documentation && Reference

The most comprehensive documentation for Turbine and how to work with Turbine apps is on the Meroxa site: [https://docs.meroxa.com/](https://docs.meroxa.com)

For the Go Reference, check out [https://pkg.go.dev/github.com/meroxa/turbine](https://pkg.go.dev/github.com/meroxa/turbine).

## Contributing

Check out the [/docs/](./docs/) folder for more information on how to contribute.
