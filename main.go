package main

func main() {
	cli := &Cli{
		Path: "cardano-cli",
	}

	_, err := cli.ProtocolParameters()

	if err != nil {
		panic(err)
	}

	// TODO: get index from UTXO
	txIns := []string{"7138bded017a3f97a84e74148ef1c2c1c3e6a62a76bb1a4b2d505d64f763938e#1"}
	txOuts := []string{
		"addr_test1vq4p6ky0htqcsq47jsxt256wve4ucwzts54gwuywmqfpdxca25527+0",
		"addr_test1vrnu54c828kyr4xqg8yt5nr62maal637w7flw6hjrk43eys7wfth5+0",
	}

	_, err = cli.BuildRaw(txIns, txOuts, 0, "tx.draft")

	if err != nil {
		panic(err)
	}

	fee, err := cli.CalculateMinFee("tx.draft", 1, 2, 1)
	if err != nil {
		panic(err)
	}

	// TODO: calculate output
	txOuts = []string{
		"addr_test1vq4p6ky0htqcsq47jsxt256wve4ucwzts54gwuywmqfpdxca25527+250000000",
		"addr_test1vrnu54c828kyr4xqg8yt5nr62maal637w7flw6hjrk43eys7wfth5+499654170",
	}
	_, err = cli.BuildRaw(txIns, txOuts, fee, "tx.draft")

	if err != nil {
		panic(err)
	}

	_, err = cli.Sign("tx.draft", "/home/phat/cardano/keys/payment1.skey")

	if err != nil {
		panic(err)
	}

	_, err = cli.Submit("tx.signed")

	if err != nil {
		panic(err)
	}
}
