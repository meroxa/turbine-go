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
$ meroxa apps init testapp --lang golang
```

The CLI will create a new folder called `testapp` located in the directory where the command was issued. If you want to initialize the app somewhere else, you can append the `--path` flag to the command (`meroxa apps init testapp --lang golang --path ~/anotherdir`). Once you enter the `testapp` directory, the contents will look like this:

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

This will be a full fledged Turbine app that can run. You can even run the tests using the command `meroxa apps run` in the root of the app directory. It just enough to show you what you need to get started.


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

Let's talk about the important parts in this code. Turbine apps have five functions that comprise of the entire DSL. Outside of these functions, you can write whatever code you want to accomplish your tasks:

```go
func (a App) Run(v turbine.Turbine) error
```

`Run` is the main entry point for the application. This is where you can initialize the Turbine framework. This is also the place where, when you deploy your Turbine app to Meroxa, Meroxa will use this as the place to boot up the application.

```go
source, err := v.Resources("source_name")
```

The `Resources` function identifies the upstream or downstream system that you want your code to work with. The `source_name` is the string identifier of the particular system. The string should map to an associated identifier in your `app.json` to configure whats being connected to. For more details, see the `app.json` section.

```go
rr, err := source.Records("collection_name", nil)
```

Once you've got `Resources` set up, you can now stream records from it but you need to identify what records you want. The `Records` function identifies the records or events that you want to stream into your data app.

```go
res, _ := v.Process(rr, Anonymize{})
```

The `Process` function is turbine's way of saying, for the records that are coming in, I want you to process these records against a function. Once your app is deployed on Meroxa, Meroxa will do the work to take each record or event that does get streamed to your app and then run your code against it. This allows Meroxa to scale out your processing relative to the velocity of the records streaming in.

```go
err = dest.Write(res, "collection_name", nil)
```

The `Write` function is optional. It's job is to take any records given to it and stream to the downstream system. In many cases, you might not need to stream data to a another system but this gives you an easy way to do so.


### `app.json`

TODO insert config options

### Testing

Testing should follow standard practices for Go.

## Documentation && Reference

The most comprehensive documentation for Turbine and how to work with Turbine apps is on the Meroxa site: [https://docs.meroxa.com/](https://docs.meroxa.com)

For the Go Reference, check out [https://pkg.go.dev/github.com/meroxa/turbine](https://pkg.go.dev/github.com/meroxa/turbine).

## Contributing

Check out the [/docs/](./docs/) folder for more information on how to contribute.
