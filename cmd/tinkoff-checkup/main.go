/*
Tinkoff-checkup checks all methods of Tinkoff Broker.
It opens position, changes conditional orders, closes position.
This can be useful for development, checking the ability
to trade with a specific token, instrument and account.

How to install:

	go install github.com/evsamsonov/tinkoff-broker/cmd/tinkoff-checkup@latest

Usage:

	tinkoff-checkup [ACCOUNT_ID] [INSTRUMENT_FIGI] [flags]

The flags are:

	-v
	    Print logger output
*/
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/evsamsonov/trengin/v2"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"

	tnkbroker "github.com/evsamsonov/tinkoff-broker"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(
			"This command checks all methods of Tinkoff Broker.\n" +
				"It opens position, changes conditional orders, closes position.",
		)
		fmt.Println("\nUsage: tinkoff-checkup [ACCOUNT_ID] [INSTRUMENT_FIGI] [-v]")
		return
	}
	accountID := os.Args[1]
	instrumentFIGI := os.Args[2]

	verbose := flag.Bool("v", false, "")
	if err := flag.CommandLine.Parse(os.Args[3:]); err != nil {
		log.Fatalf("Failed to parse args: %s", err)
	}

	checkupParams := NewCheckupParams(accountID, instrumentFIGI)
	if err := checkupParams.AskUser(); err != nil {
		log.Fatalf("Failed to get checkup params: %s", err)
	}

	checkuper, err := NewTinkoffCheckuper(*verbose)
	if err != nil {
		log.Fatalf("Failed to create tinkoff checkuper: %s", err)
	}
	if err := checkuper.CheckUp(checkupParams); err != nil {
		log.Fatalf("Failed to check up: %s", err)
	}
	fmt.Println("Check up is successful! 🍺")
}

type CheckUpArgs struct {
	accountID        string
	instrumentFIGI   string
	tinkoffToken     string
	stopLossOffset   float64
	takeProfitOffset float64
	positionType     trengin.PositionType
}

func NewCheckupParams(accountID, instrumentFIGI string) CheckUpArgs {
	return CheckUpArgs{
		accountID:      accountID,
		instrumentFIGI: instrumentFIGI,
	}
}

func (c *CheckUpArgs) AskUser() error {
	fmt.Printf("Paste Tinkoff token: ")
	tokenBytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return fmt.Errorf("read token: %w", err)
	}
	fmt.Println()
	c.tinkoffToken = string(tokenBytes)

	var positionType string
	fmt.Print("Enter position direction [long, short]: ")
	if _, err = fmt.Scanln(&positionType); err != nil {
		return fmt.Errorf("read stop loss indent: %w", err)
	}
	if positionType != "long" && positionType != "short" {
		return fmt.Errorf("read position direction: %w", err)
	}
	c.positionType = trengin.Long
	if positionType == "short" {
		c.positionType = trengin.Short
	}

	var StopLossOffset, TakeProfitOffset float64
	fmt.Print("Enter stop loss indent [0 - skip]: ")
	if _, err = fmt.Scanln(&StopLossOffset); err != nil {
		return fmt.Errorf("read stop loss indent: %w", err)
	}
	c.stopLossOffset = StopLossOffset

	fmt.Print("Enter take profit indent [0 - skip]: ")
	if _, err = fmt.Scanln(&TakeProfitOffset); err != nil {
		return fmt.Errorf("read take profit indent: %w", err)
	}
	c.takeProfitOffset = TakeProfitOffset
	return nil
}

type TinkoffCheckuper struct {
	logger *zap.Logger
}

func NewTinkoffCheckuper(verbose bool) (*TinkoffCheckuper, error) {
	logger := zap.NewNop()
	if verbose {
		var err error
		logger, err = zap.NewDevelopment(zap.IncreaseLevel(zap.DebugLevel))
		if err != nil {
			return nil, fmt.Errorf("create logger: %w", err)
		}
	}
	return &TinkoffCheckuper{
		logger: logger,
	}, nil
}

func (t *TinkoffCheckuper) CheckUp(params CheckUpArgs) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	tinkoffConfig := investgo.Config{
		EndPoint:  "invest-public-api.tinkoff.ru:443",
		Token:     params.tinkoffToken,
		AppName:   "evsamsonov.gdealer",
		AccountId: params.accountID,
	}
	tinkoffClient, err := investgo.NewClient(ctx, tinkoffConfig, t.logger.Sugar())
	if err != nil {
		return fmt.Errorf("create tinkoff client: %w", err)
	}

	tinkoffBroker, err := tnkbroker.New(
		tinkoffClient,
		params.accountID,
		tnkbroker.WithLogger(t.logger),
	)
	if err != nil {
		return fmt.Errorf("create tinkoff broker: %w", err)
	}

	g.Go(func() error {
		defer cancel()
		if err := tinkoffBroker.Run(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return fmt.Errorf("tinkoff broker: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		defer cancel()
		t.WaitAnyKey("Press any key for open position...")

		openPositionAction := trengin.OpenPositionAction{
			Type:             params.positionType,
			FIGI:             params.instrumentFIGI,
			Quantity:         1,
			StopLossOffset:   params.stopLossOffset,
			TakeProfitOffset: params.takeProfitOffset,
		}
		position, positionClosed, err := tinkoffBroker.OpenPosition(ctx, openPositionAction)
		if err != nil {
			return fmt.Errorf("open position: %w", err)
		}
		fmt.Printf(
			"Position was opened. Open price: %f, stop loss: %f, take profit: %f, commission: %f\n",
			position.OpenPrice,
			position.StopLoss,
			position.TakeProfit,
			position.Commission,
		)

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			case pos := <-positionClosed:
				fmt.Printf(
					"Position was closed. Conditional orders was removed. "+
						"Close price: %f, profit: %f, commission: %f\n",
					pos.ClosePrice,
					pos.Profit(),
					position.Commission,
				)
			}
			return nil
		})
		t.WaitAnyKey("Press any key for reduce by half conditional orders...")

		changeConditionalOrderAction := trengin.ChangeConditionalOrderAction{
			PositionID: position.ID,
			StopLoss:   position.OpenPrice - params.stopLossOffset/2*position.Type.Multiplier(),
			TakeProfit: position.OpenPrice + params.takeProfitOffset/2*position.Type.Multiplier(),
		}
		position, err = tinkoffBroker.ChangeConditionalOrder(ctx, changeConditionalOrderAction)
		if err != nil {
			return fmt.Errorf("change condition order: %w", err)
		}
		fmt.Printf(
			"Conditional orders was changed. New stop loss: %f, new take profit: %f\n",
			position.StopLoss,
			position.TakeProfit,
		)
		t.WaitAnyKey("Press any key for close position...")

		closePositionAction := trengin.ClosePositionAction{PositionID: position.ID}
		position, err = tinkoffBroker.ClosePosition(ctx, closePositionAction)
		if err != nil {
			return fmt.Errorf("close position: %w", err)
		}
		return nil
	})

	return g.Wait()
}
func (t *TinkoffCheckuper) WaitAnyKey(msg string) {
	fmt.Print(msg)
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}
