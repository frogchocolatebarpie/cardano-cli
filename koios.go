package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/cardano-community/koios-go-client"
)

type KoiosClient struct {
	inner   koios.Client
	options koios.Option
}

func NewKoiosClient(options koios.Option) *KoiosClient {
	inner, err := koios.New(options)

	if err != nil {
		return nil
	}

	return &KoiosClient{
		options: options,
		inner:   *inner,
	}

}

func (k *KoiosClient) GetTip() *koios.TipResponse {
	tip, err := k.inner.GetTip(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return tip
}

func (k *KoiosClient) IsHealthy() bool {
	tip := k.GetTip()

	if tip == nil {
		log.Fatal(errors.New("No tip response"))
		return false
	}

	return tip.StatusCode == 200
}

func (k *KoiosClient) GetBlock(hash koios.BlockHash) *koios.BlockInfoResponse {
	block, err := k.inner.GetBlockInfo(context.Background(), hash, nil)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return block
}

func (k *KoiosClient) LatestBlock() *koios.Block {
	tip := k.GetTip()
	if tip == nil {
		log.Fatal(errors.New("No tip response"))
		return nil
	}

	hash := koios.BlockHash(tip.Data.Hash)

	block := k.GetBlock(hash)
	if block == nil {
		log.Fatal(errors.New("No block response"))
		return nil
	}

	return block.Data
}

func (k *KoiosClient) BlockHeight() (int, error) {
	tip := k.GetTip()

	if tip == nil {
		return 0, errors.New("No tip response")
	}

	return tip.Data.BlockNo, nil
}

func (k *KoiosClient) NewTxs(fromHeight int, interestedAddrs map[string]bool) (*koios.AddressTxsResponse, error) {
	addrs := []koios.Address{""}
	addressTxs, err := k.inner.GetAddressTxs(context.Background(), addrs, uint64(fromHeight), nil)

	if err != nil {
		return nil, err
	}

	return addressTxs, nil
}

func (k *KoiosClient) SubmitTx() {
	// k.inner.SubmitSignedTx()
}

func main() {
	k := NewKoiosClient(koios.Host(koios.TestnetHost))

	if k == nil {
		return
	}

	isHealthy := k.IsHealthy()
	fmt.Println("IsHealthy: ", isHealthy)

	latestBlock := k.LatestBlock()
	fmt.Println("Lastest Block: ", latestBlock)

	blockHeight, err := k.BlockHeight()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("BlockHeight: ", blockHeight)
}
