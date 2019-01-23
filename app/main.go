package app

import (
	"git.parallelcoin.io/pod/cmd/ctl"
<<<<<<< HEAD
	"git.parallelcoin.io/pod/cmd/conf"
	"git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/cmd/shell"
=======
	"git.parallelcoin.io/pod/cmd/node"
>>>>>>> master
	"github.com/tucnak/climax"
)

var interrupt <-chan struct{}

// PodApp is the climax main app controller for pod
var PodApp = climax.Application{
	Name:     "pod",
	Brief:    "multi-application launcher for Parallelcoin Pod",
	Version:  version(),
	Commands: []climax.Command{},
	Topics:   []climax.Topic{},
	Groups:   []climax.Group{},
	Default:  nil,
}

// Main is the real pod main
func Main() int {
	PodApp.AddCommand(CtlCommand)
	PodApp.AddCommand(NodeCommand)
	PodApp.AddCommand(WalletCommand)
	PodApp.AddCommand(ShellCommand)
	PodApp.AddCommand(ConfCommand)
	return PodApp.Run()
}
