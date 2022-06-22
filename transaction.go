package main

import (
	"strconv"
	"strings"
)

func (c *Cli) BuildRaw(txIns []string, txOuts []string, fee int, outFile string) ([]byte, error) {
	args := []string{"transaction", "build-raw"}

	for _, txIn := range txIns {
		args = append(args, "--tx-in", txIn)
	}

	for _, txOut := range txOuts {
		args = append(args, "--tx-out", txOut)
	}

	args = append(args, "--fee", strconv.Itoa(fee))

	args = append(args, "--out-file", outFile)

	out, err := c.Exec(args...)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Cli) CalculateMinFee(txBodyFile string, txInCount uint64, txOutCount uint64, witnessCount uint64) (int, error) {
	out, err := c.Exec(
		"transaction", "calculate-min-fee",
		"--tx-body-file", txBodyFile,
		"--tx-in-count", strconv.Itoa(int(txInCount)),
		"--tx-out-count", strconv.Itoa(int(txInCount)),
		"--witness-count", strconv.Itoa(int(witnessCount)),
		"--testnet-magic", "1097911063",
		"--protocol-params-file", "protocol.json",
	)

	if err != nil {
		return 0, err
	}

	args := strings.Fields(string(out))

	fee, err := strconv.Atoi(args[0])

	if err != nil {
		return 0, err
	}

	return fee, nil
}

func (c *Cli) Sign(txBodyFile string, signingKeyFile string) ([]byte, error) {
	out, err := c.Exec(
		"transaction", "sign",
		"--tx-body-file", txBodyFile,
		"--signing-key-file", signingKeyFile,
		"--testnet-magic", "1097911063",
		"--out-file", "tx.signed",
	)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Cli) Submit(txFile string) ([]byte, error) {
	out, err := c.Exec(
		"transaction", "submit",
		"--tx-file", txFile,
		"--testnet-magic", "1097911063",
	)

	if err != nil {
		return nil, err
	}

	return out, nil
}
