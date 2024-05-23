package test

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"log"
	"testing"
)

func TestEthereumWebsocket(t *testing.T) {
	// 使用WS的方式连接
	client, err := ethclient.Dial("wss://ws-sepolia.reservoir.tools")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// address监控的合约地址
	address := common.HexToAddress("0x8bc970f1db3dc6377ef5230d2b893abadb833f93")
	//topicHash := crypto.Keccak256Hash([]byte("LotteryEvent"))

	// 过滤处理
	query := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		//Topics:    [][]common.Hash({topicHash}),
	}

	logs := make(chan types.Log)
	//
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	// 消息的接收以及处理
	for {
		select {
		// select 阻塞监控多个channel,任意一个channel有消息则解除阻塞，并且执行case内
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog) // pointer to event log
		}
	}
}
