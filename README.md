# tinkoff-broker

[![Lint Status](https://github.com/evsamsonov/tinkoff-broker/actions/workflows/lint.yml/badge.svg)](https://github.com/evsamsonov/tinkoff-broker/actions?workflow=golangci-lint)
[![Test Status](https://github.com/evsamsonov/tinkoff-broker/actions/workflows/test.yml/badge.svg)](https://github.com/evsamsonov/tinkoff-broker/actions?workflow=test)
[![Go Report Card](https://goreportcard.com/badge/github.com/evsamsonov/tinkoff-broker)](https://goreportcard.com/report/github.com/evsamsonov/tinkoff-broker)
[![codecov](https://codecov.io/gh/evsamsonov/tinkoff-broker/branch/master/graph/badge.svg?token=AC751PKE5Y)](https://codecov.io/gh/evsamsonov/tinkoff-broker)

An implementation of [trengin.Broker](https://github.com/evsamsonov/trengin) using
[T-Invest API](https://developer.tbank.ru/invest/intro/intro/)
([invest-api-go-sdk](https://github.com/RussianInvestments/invest-api-go-sdk))
for creating automated trading robots.

## Requirements

- Go 1.22.5+

## Features

- Opens position, changes stop loss and take profit, closes position.
- Tracks open position.
- Supports multiple open positions at the same time.
- Commission in position is approximate.

## How to use

Create a new `Tinkoff` object using constructor `New`. Pass an invest-api-go-sdk
client and a user account identifier.

```go
package main

import (
	"context"
	"log"

	"github.com/evsamsonov/trengin/v2"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"

	tnkbroker "github.com/evsamsonov/tinkoff-broker/v2"
)

func main() {
	ctx := context.Background()

	config := investgo.Config{
		EndPoint:  "invest-public-api.tbank.ru:443", // sandbox: sandbox-invest-public-api.tbank.ru:443
		Token:     "[t-invest-token]",
		AccountId: "[account-id]",
		AppName:   "my-app",
	}
	client, err := investgo.NewClient(ctx, config, zap.NewNop().Sugar())
	if err != nil {
		log.Fatal("Failed to create T-Invest client", zap.Error(err))
	}
	defer func() {
		if err := client.Stop(); err != nil {
			log.Fatal("Failed to stop T-Invest client", zap.Error(err))
		}
	}()

	broker, err := tnkbroker.New(
		client,
		config.AccountId,
		// options...
	)
	if err != nil {
		log.Fatal("Failed to create tinkoff broker")
	}

	tradingEngine := trengin.New(&Strategy{}, broker)
	if err = tradingEngine.Run(ctx); err != nil {
		log.Fatal("Trading engine crashed")
	}
}

type Strategy struct{}

func (s *Strategy) Run(ctx context.Context, actions trengin.Actions) error { panic("implement me") }
```

See more details in [trengin documentation](https://github.com/evsamsonov/trengin).

### Option

You can configure `Tinkoff` using `Option`:

| Methods                           | Returns Option which                                                             |
|-----------------------------------|----------------------------------------------------------------------------------|
| `WithLogger`                      | Sets logger. The default logger is no-op Logger.                                 |
| `WithProtectiveSpread`            | Sets protective spread in percent for executing orders. The default value is 1%. |
| `WithTradeStreamRetryTimeout`     | Defines retry timeout on trade stream error.                                     |
| `WithTradeStreamPingWaitDuration` | Defines duration how long we wait for ping before reconnection.                  |

## Checkup

Use `tinkoff-checkup` to verify that a token, instrument and account can trade
via T-Invest API.

### How to install

```bash
go install github.com/evsamsonov/tinkoff-broker/v2/cmd/tinkoff-checkup@latest
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

Makefile tasks require Docker and Go.

```bash
$ make help
doc                            Run doc server using docker
lint                           Run golang lint using docker
pre-push                       Run golang lint and test
test                           Run tests
generate                       Run go generate
```
