# tinkoff-broker

[![Lint Status](https://github.com/evsamsonov/tinkoff-broker/actions/workflows/lint.yml/badge.svg)](https://github.com/evsamsonov/tinkoff-broker/actions?workflow=golangci-lint)
[![Test Status](https://github.com/evsamsonov/tinkoff-broker/actions/workflows/test.yml/badge.svg)](https://github.com/evsamsonov/tinkoff-broker/actions?workflow=test)
[![Go Report Card](https://goreportcard.com/badge/github.com/evsamsonov/tinkoff-broker)](https://goreportcard.com/report/github.com/evsamsonov/tinkoff-broker)
[![codecov](https://codecov.io/gh/evsamsonov/tinkoff-broker/branch/master/graph/badge.svg?token=AC751PKE5Y)](https://codecov.io/gh/evsamsonov/tinkoff-broker)

An implementation of [trengin.Broker](http://github.com/evsamsonov/trengin) using [Tinkoff Invest API](https://tinkoff.github.io/investAPI/) 
for creating automated trading robots. 

## Features
- Opens position, changes stop loss and take profit, closes position.
- Tracks open position.
- Doesn't support multiple open positions at the same time.
- Commission in position is approximate.

## How to use

Create a new `Tinkoff` object using constructor `New`. Pass [full-access token](https://tinkoff.github.io/investAPI/token/),
user account identifier, [FIGI](https://tinkoff.github.io/investAPI/faq_identification/) of a trading instrument.

```go
package main

import (
	"context"
	"log"

	tnkbroker "github.com/evsamsonov/tinkoff-broker"
	"github.com/evsamsonov/trengin/v2"
)

func main() {
	tinkoffBroker, err := tnkbroker.New(
		"tinkoff-token",
		"123",
		"BBG004730N88",
		// options...
	)
	if err != nil {
		log.Fatal("Failed to create tinkoff broker")
	}
	
	tradingEngine := trengin.New(&Strategy{}, tinkoffBroker)
	if err = tradingEngine.Run(context.Background()); err != nil {
		log.Fatal("Trading engine crashed")
	}
}

type Strategy struct{}
func (s *Strategy) Run(ctx context.Context) error { panic("implement me") }
func (s *Strategy) Actions() trengin.Actions { panic("implement me") }
```

See more details in [trengin documentation](http://github.com/evsamsonov/trengin).

### Option

You can configure `Tinkoff` to use `Option`

| Methods                           | Returns Option which                                                             |
|-----------------------------------|----------------------------------------------------------------------------------|
| `WithLogger`                      | Sets logger. The default logger is no-op Logger.                                 |
| `WithAppName`                     | Sets [x-app-name](https://tinkoff.github.io/investAPI/grpc/#appname).            |
| `WithProtectiveSpread`            | Sets protective spread in percent for executing orders. The default value is 1%. |
| `WithTradeStreamRetryTimeout`     | Defines retry timeout on trade stream error.                                     |
| `WithTradeStreamPingWaitDuration` | Defines duration how long we wait for ping before reconnection.                  |

## Checkup

Use `tinkoff-checkup` for checking the ability to trade with a specific token, instrument and account. 

### How to install

```bash
go install github.com/evsamsonov/tinkoff-broker/cmd/tinkoff-checkup@latest
```

### How to use 

```bash
tinkoff-checkup [ACCOUNT_ID] [INSTRUMENT_FIGI] [-v]
 ```

| Flag | Description         |
|------|---------------------|
| `-v` | Print logger output |

## Development

### Makefile 

Makefile tasks are required docker and golang.

```bash
$ make help    
doc                            Run doc server using docker
lint                           Run golang lint using docker
pre-push                       Run golang lint and test
test                           Run tests
```
