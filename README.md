> ### development branch

[![GoDoc](https://img.shields.io/badge/godoc-documentation-blue.svg)](https://doc.parallelcoin.io/pkg/git.parallelcoin.io/pod/)
# Parallelcoin Pod

Fully integrated all-in-one cli client, full node, wallet server, miner and GUI wallet for Parallelcoin

> ## IMPORTANT
> 
> **Development branch** can be found [here](https://seed1.parallelcoin.io/dev/pod/src/branch/dev), where the current work is located.
> 
> here's our own godoc server documentation for this package: [here](http://89.40.12.55:8008/pkg/git.parallelcoin.io/pod/)
> 
> **This code will not function correctly currently, please be patient while it is fixed on the Development Branch.**

Pod is a multi-application with multiple submodules for different functions. It is self-configuring and configurations can be changed from the commandline as well as editing the json files directly, so the binary itself is the complete distribution for the suite.

It consists of 4 main modules:

1. pod/ctl - command line interface to send queries to a node or wallet and prints the results to the stdout
2. pod/node - full node for parallelcoin network, including optional indexes for address and transaction search, low latency miner controller
3. pod/wallet - wallet server that runs separately from the full node but depends on a full node RPC for much of its functionality, and includes a GUI front end
4. pod/shell - combined full node and wallet server with optional GUI

The shell is currently a simple wallet but will be expanded into a full application framework/shell.

## Building

To make life simpler, there is a builder app in `cmd/` and if you source init.sh `. init.sh` it will set your path to include the `bin/` directory and build and place an executable `bld` in there which builds the main project executable, with the version timestamp set, and puts it also in there so then after that you can test how it works after you make changes.

Otherwise, you can just `go install` in the root directory and `pod` will be placed in your `GOBIN` directory.
