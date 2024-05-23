package test

import (
	"context"
	"fmt"
	"github.com/bcds/go-hpc-common/types"
	"github.com/bcds/gosdk/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"testing"
)

func TestEthereumWebsocket(t *testing.T) {
	client, err := ethclient.Dial("wss://ws-sepolia.reservoir.tools")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0x8bc970f1db3dc6377ef5230d2b893abadb833f93")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog) // pointer to event log
		}
	}
}
