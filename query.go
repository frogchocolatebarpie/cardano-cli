package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Tip struct {
	Era          string
	SyncProgress string
	Hash         string
	Epoch        uint64
	Slot         uint64
	Block        uint64
}

type UTxO struct {
	TxHash string
	Index  uint64
	Amount uint64
}

func (c *Cli) UTxOs(addr string) ([]UTxO, error) {
	out, err := c.Exec("query", "utxo", "--address", addr, "--testnet-magic", "1097911063")
	if err != nil {
		return nil, err
	}

	utxos := []UTxO{}
	lines := strings.Split(string(out), "\n")

	if len(lines) < 4 {
		return utxos, nil
	}

	for _, line := range lines[2 : len(lines)-1] {
		args := strings.Fields(line)
		if len(args) < 4 {
			return nil, fmt.Errorf("malformed cli response")
		}

		txHash := args[0]

		index, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}

		amount, err := strconv.Atoi(args[2])
		if err != nil {
			return nil, err
		}

		utxos = append(utxos, UTxO{
			TxHash: txHash,
			Index:  uint64(index),
			Amount: uint64(amount),
		})
	}

	return utxos, nil
}

func (c *Cli) Tip() (*Tip, error) {
	out, err := c.Exec("query", "tip", "--testnet-magic", "1097911063")
	if err != nil {
		return nil, err
	}

	cliTip := &Tip{}
	if err = json.Unmarshal(out, cliTip); err != nil {
		return nil, err
	}

	return cliTip, nil
}

func (c *Cli) ProtocolParameters() ([]byte, error) {

	out, err := c.Exec("query", "protocol-parameters", "--out-file", "protocol.json", "--testnet-magic", "1097911063")

	if err != nil {
		return nil, err
	}

	return out, nil
}
