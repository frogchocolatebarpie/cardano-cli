package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/echovl/cardano-go"
)

type Tip struct {
	Era          string
	SyncProgress string
	Hash         string
	Epoch        uint64
	Slot         uint64
	Block        uint64
}

type cliTx struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	CborHex     string `json:"cborHex"`
}

type Cli struct{}

func (c *Cli) runCommand(args ...string) ([]byte, error) {
	out := &bytes.Buffer{}

	args = append(args, "--testnet-magic", "1097911063")

	cmd := exec.Command("cardano-cli", args...)
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (c *Cli) UTxOs(addr cardano.Address) ([]cardano.UTxO, error) {
	out, err := c.runCommand("query", "utxo", "--address", addr.Bech32())
	if err != nil {
		return nil, err
	}

	utxos := []cardano.UTxO{}
	lines := strings.Split(string(out), "\n")

	if len(lines) < 3 {
		return utxos, nil
	}

	for _, line := range lines[2 : len(lines)-1] {
		amount := cardano.NewValue(0)
		args := strings.Fields(line)
		if len(args) < 4 {
			return nil, fmt.Errorf("malformed cli response")
		}
		txHash, err := cardano.NewHash32(args[0])
		if err != nil {
			return nil, err
		}
		index, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}
		lovelace, err := strconv.Atoi(args[2])
		if err != nil {
			return nil, err
		}
		amount.Coin = cardano.Coin(lovelace)

		assets := strings.Split(line, "+")
		for _, asset := range assets[1 : len(assets)-1] {
			args := strings.Fields(asset)
			quantity := args[0]
			unit := strings.ReplaceAll(args[1], ".", "")
			unitBytes, err := hex.DecodeString(unit)
			if err != nil {
				return nil, err
			}
			policyID := cardano.NewPolicyIDFromHash(unitBytes[:28])
			assetName := string(unitBytes[28:])
			assetValue, err := strconv.ParseUint(quantity, 10, 64)
			if err != nil {
				return nil, err
			}
			currentAssets := amount.MultiAsset.Get(policyID)
			if currentAssets != nil {
				currentAssets.Set(
					cardano.NewAssetName(assetName),
					cardano.BigNum(assetValue),
				)
			} else {
				amount.MultiAsset.Set(
					policyID,
					cardano.NewAssets().
						Set(
							cardano.NewAssetName(string(assetName)),
							cardano.BigNum(assetValue),
						),
				)
			}
		}

		utxos = append(utxos, cardano.UTxO{
			Spender: addr,
			TxHash:  txHash,
			Index:   uint64(index),
			Amount:  amount,
		})
	}

	return utxos, nil
}

func (c *Cli) Tip() (*cardano.NodeTip, error) {
	out, err := c.runCommand("query", "tip")
	if err != nil {
		return nil, err
	}

	cliTip := &Tip{}
	if err = json.Unmarshal(out, cliTip); err != nil {
		return nil, err
	}

	return &cardano.NodeTip{
		Epoch: cliTip.Epoch,
		Block: cliTip.Block,
		Slot:  cliTip.Slot,
	}, nil
}

func main() {
	cli := &Cli{}
	address, err := cardano.NewAddress("addr_test1vp9uhllavnhwc6m6422szvrtq3eerhleer4eyu00rmx8u6c42z3v8")
	if err != nil {
		panic(err)
	}

	out, err := cli.UTxOs(address)

	if err != nil {
		panic(err)
	}

	fmt.Println(out[0].TxHash)
}
