package app

import (
	"encoding/json"

	"git.parallelcoin.io/pod/pkg/clog"
	"git.parallelcoin.io/pod/cmd/wallet"
)

func runNode() {
	j, _ := json.MarshalIndent(Config, "", "  ")
	log <- cl.Trc("running with configuration:\n" + string(j))
	walletmain.Main(Config.Wallet)
}
