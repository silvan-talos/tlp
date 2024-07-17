# tlp - Telemetry Logging Package

`tlp` is a logging package that unifies logging and telemetry data.

[![Go Reference](https://pkg.go.dev/badge/github.com/silvan-talos/tlp.svg)](https://pkg.go.dev/github.com/silvan-talos/tlp)

[Usage example](example)

## Usage

### Installation

```shell
go get github.com/silvan-talos/tlp
```

### Out-of-the-box integration

Start by directly typing logging commands. Example:

```go
import (
"github.com/silvan-talos/tlp/log"
)

// ... use as needed
log.Error(ctx, "find product", "err", err, "query", query)
log.Info(ctx, "user created successfully", "id", id)
```

### Customize behavior

#### Customize using config file

By default, the logging library checks for a config file named `log-config.yml` in the program running context
directory. Check the config example [here](config/config_example.yml).

For a specific config path, it needs to be customized using code. Example:

```go
package main

import (
    "context"

    "github.com/silvan-talos/tlp/log"
    "github.com/silvan-talos/tlp/config"
)

func main() {
    configPath := "my-location/config.yml"
    var cfg config.Config
    err := config.LoadFromYAML(configPath, &cfg)
    if err != nil {
        // handle err
    }
    logger := log.NewLoggerFromConfig(cfg.Log)
    logger.SetDefault()
    log.Info(context.Background(), "log configured", "configPath", configPath)
}
```

#### Customize using code

The [API documentation](https://pkg.go.dev/github.com/silvan-talos/tlp@v0.0.0-20240717010324-e2378fdccd57/log) provides
a few options to customize logger behavior.

`WithLevel` returns a copy of the original logger with the desired log-level
set, if parsable. `level` format should be either `debug`, `info`, `warn`, `error` or `level(number)`. See log
level [details](#log-levels) to understand the number meaning in level context.

`WithAttrs` creates a copy of the receiver logger and sets an attribute list to be logged for each message.

In order to persist a custom logger and use it from across the packages, you can set it as default using
the `SetDefault` function. See the above code snippet as a demo.

## Extendability

### Log levels

A log level is a number that represents the severity of a log event. The following list shows the default level-number
connection:

- level `DEBUG` = -4
- level `INFO` = 0
- level `WARN` = 4
- level `ERROR` = 8

The naming of custom log levels is `LEVEL(number)`, where `number` has the same meaning as above. `LEVEL(-4)` is
actually `DEBUG`. To log a message to a specific level use `log.Default().Log(ctx, 100, "message")`
or `log.Default().Log(ctx, logging.Level(100), "message")`.

### Log drivers

Any struct that implements
the [`Driver` interface](https://pkg.go.dev/github.com/silvan-talos/tlp@v0.0.0-20240717010324-e2378fdccd57/log#Driver)
can be used as a driver. `tlp` provides you all the logging entry details, alongside the context you may use to extract
other parameters. See [text driver](text/driver.go) and [JSON driver](json/driver.go) as example for driver
implementations.

```go
// set a custom driver

logLevel := logging.Level(-4)
driver := // instantiate driver
logger := NewLogger(driver, logLevel)
logger.SetDefault()

```

### Transaction recorders

Any struct that implements
the [`Recorder` interface](https://pkg.go.dev/github.com/silvan-talos/tlp@v0.0.0-20240717010324-e2378fdccd57/transaction#Recorder)
can be used as a transaction recorder. It offers the possibility to use an actual transaction tracer behind the scenes,
while logging the provided TraceID as usual for correlation.
