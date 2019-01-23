package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/cmd/wallet"
	"git.parallelcoin.io/pod/cmd/ctl"
	"github.com/tucnak/climax"
)

var confFile = def.DefaultDataDir + "/conf"

// ConfConfig is the settings that can be set to synchronise across all pod modules
type ConfConfig struct {
	NodeListeners    []string
	NodeRPCListeners []string
	WalletListeners  []string
	NodeUser         string
	NodePass         string
	WalletPass       string
	RPCKey           string
	RPCCert          string
	CAFile           string
	TLS              bool
	SkipVerify       bool
	Proxy            string
	ProxyUser        string
	ProxyPass        string
	Network          string
}

// ConfigCfg is the configuration for this tool
var ConfigCfg ConfConfig

// Apps is the central repository of all the other app configurations
var Apps AppConfigs

// ConfConfigs are the configurations for each app that are applied
type ConfConfigs struct {
	Ctl    CtlConfig
	Node   NodeConfig
	Wallet WalletConfig
	Shell  ShellCfg
}

var f = pu.GenFlag
var t = pu.GenTrig
var s = pu.GenShort
var l = pu.GenLog

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var ConfCommand = climax.Command{
	Name:  "conf",
	Brief: "sets configurations common across modules",
	Help:  "automates synchronising common settings between servers and clients",
	Flags: []climax.Flag{
		t("version", "V", "show version number and quit"),
		t("init", "i", "resets configuration to defaults"),
		t("show", "s", "prints currently configuration"),

		f("nodelistener", "main peer to peer address for apps that connect to the parallelcoin peer to peer network"),
		f("noderpclistener", "address where node listens for RPC"),
		f("walletlistener", "address where wallet listens for RPC"),
		s("user", "u", "username for all the things"),
		s("pass", "P", "password for all the things"),
		s("walletpass", "w", "public password for wallet"),
		f("rpckey", "RPC server certificate key"),
		f("rpccert", "RPC server certificate"),
		f("cafile", "RPC server certificate chain for validation"),
		f("tls", "enable/disable TLS"),
		f("skipverify", "do not verify TLS certificates (not recommended!)"),
		f("proxy", "connect via SOCKS5 proxy"),
		f("proxyuser", "username for proxy"),
		f("proxypass", "password for proxy"),

		f("network", "connect to [mainnet|testnet|regtestnet|simnet]"),
	},
	Examples: []climax.Example{
		{
			Usecase:     "--nodeuser=user --nodepass=pa55word",
			Description: "set the username and password for the node RPC",
		},
	},
	Handle: func(ctx climax.Context) int {
		if ctx.Is("version") {
			fmt.Println("pod conf version", Version())
			os.Exit(0)
		}
		if ctx.Is("init") {
			WriteDefaultConfConfig(confFile)
		} else {
			if _, err := os.Stat(confFile); os.IsNotExist(err) {
				WriteDefaultConfConfig(confFile)
			} else {
				cfgData, err := ioutil.ReadFile(confFile)
				if err != nil {
					WriteDefaultConfConfig(confFile)
				}
				err = json.Unmarshal(cfgData, &ConfigCfg)
				if err != nil {
					WriteDefaultConfConfig(confFile)
				}
			}
		}
		configConf(&ctx, confFile)
		runCtl()
		return 0
	},
}

func configConf(ctx *climax.Context, cfgFile string) {
	// First load all of the module configurations and unmarshal into their structs
	confs := []string{
		def.DefaultDataDir + "/ctl/conf",
		def.DefaultDataDir + "/node/conf",
		def.DefaultDataDir + "/wallet/conf",
		def.DefaultDataDir + "/shell/conf",
	}

	// If we can't parse the config files we just reset them to default

	ctlCfg := *ctl.DefaultConfConfig()
	if _, err := os.Stat(confs[0]); os.IsNotExist(err) {
		ctl.WriteDefaultConfConfig(confs[0])
	} else {
		ctlCfgData, err := ioutil.ReadFile(confs[0])
		if err != nil {
			shell.WriteDefaultConfConfig(confs[0])
		} else {
			err = json.Unmarshal(ctlCfgData, &ctlCfg)
			if err != nil {
				shell.WriteDefaultConfConfig(confs[0])
			}
		}
	}

	nodeCfg := *node.DefaultConfConfig()
	if _, err := os.Stat(confs[1]); os.IsNotExist(err) {
		node.WriteDefaultConfConfig(confs[1])
	} else {
		nodeCfgData, err := ioutil.ReadFile(confs[1])
		if err != nil {
			node.WriteDefaultConfConfig(confs[1])
		} else {
			err = json.Unmarshal(nodeCfgData, &nodeCfg)
			if err != nil {
				node.WriteDefaultConfConfig(confs[1])
			}
		}
	}

	walletCfg := *walletrun.DefaultConfConfig()
	if _, err := os.Stat(confs[2]); os.IsNotExist(err) {
		walletrun.WriteDefaultConfConfig(confs[2])
	} else {
		walletCfgData, err := ioutil.ReadFile(confs[2])
		if err != nil {
			walletrun.WriteDefaultConfConfig(confs[2])
		} else {
			err = json.Unmarshal(walletCfgData, &walletCfg)
			if err != nil {
				walletrun.WriteDefaultConfConfig(confs[2])
			}
		}
	}

	shellCfg := *shell.DefaultConfConfig()
	if _, err := os.Stat(confs[3]); os.IsNotExist(err) {
		shell.WriteDefaultConfConfig(confs[3])
	} else {
		shellCfgData, err := ioutil.ReadFile(confs[3])
		if err != nil {
			shell.WriteDefaultConfConfig(confs[3])
		} else {
			err = json.Unmarshal(shellCfgData, &shellCfg)
			if err != nil {
				shell.WriteDefaultConfConfig(confs[3])
			}
		}
	}

	// push all current settings as from the conf configuration to the module configs
	nodeCfg.Node.Listeners = ConfigCfg.NodeListeners
	shellCfg.Node.Listeners = ConfigCfg.NodeListeners
	nodeCfg.Node.RPCListeners = ConfigCfg.NodeRPCListeners
	walletCfg.Wallet.RPCConnect = ConfigCfg.NodeRPCListeners[0]
	shellCfg.Node.RPCListeners = ConfigCfg.NodeRPCListeners
	shellCfg.Wallet.RPCConnect = ConfigCfg.NodeRPCListeners[0]
	ctlCfg.RPCServer = ConfigCfg.NodeRPCListeners[0]

	walletCfg.Wallet.LegacyRPCListeners = ConfigCfg.WalletListeners
	ctlCfg.Wallet = ConfigCfg.NodeRPCListeners[0]
	shellCfg.Wallet.LegacyRPCListeners = ConfigCfg.NodeRPCListeners
	walletCfg.Wallet.LegacyRPCListeners = ConfigCfg.WalletListeners
	ctlCfg.Wallet = ConfigCfg.WalletListeners[0]
	shellCfg.Wallet.LegacyRPCListeners = ConfigCfg.WalletListeners

	nodeCfg.Node.RPCUser = ConfigCfg.NodeUser
	walletCfg.Wallet.PodUsername = ConfigCfg.NodeUser
	walletCfg.Wallet.Username = ConfigCfg.NodeUser
	shellCfg.Node.RPCUser = ConfigCfg.NodeUser
	shellCfg.Wallet.PodUsername = ConfigCfg.NodeUser
	shellCfg.Wallet.Username = ConfigCfg.NodeUser
	ctlCfg.RPCUser = ConfigCfg.NodeUser

	nodeCfg.Node.RPCPass = ConfigCfg.NodePass
	walletCfg.Wallet.PodPassword = ConfigCfg.NodePass
	walletCfg.Wallet.Password = ConfigCfg.NodePass
	shellCfg.Node.RPCPass = ConfigCfg.NodePass
	shellCfg.Wallet.PodPassword = ConfigCfg.NodePass
	shellCfg.Wallet.Password = ConfigCfg.NodePass
	ctlCfg.RPCPass = ConfigCfg.NodePass

	nodeCfg.Node.RPCKey = ConfigCfg.RPCKey
	walletCfg.Wallet.RPCKey = ConfigCfg.RPCKey
	shellCfg.Node.RPCKey = ConfigCfg.RPCKey
	shellCfg.Wallet.RPCKey = ConfigCfg.RPCKey

	nodeCfg.Node.RPCCert = ConfigCfg.RPCCert
	walletCfg.Wallet.RPCCert = ConfigCfg.RPCCert
	shellCfg.Node.RPCCert = ConfigCfg.RPCCert
	shellCfg.Wallet.RPCCert = ConfigCfg.RPCCert

	walletCfg.Wallet.CAFile = ConfigCfg.CAFile
	shellCfg.Wallet.CAFile = ConfigCfg.CAFile

	nodeCfg.Node.TLS = ConfigCfg.TLS
	walletCfg.Wallet.EnableClientTLS = ConfigCfg.TLS
	shellCfg.Node.TLS = ConfigCfg.TLS
	shellCfg.Wallet.EnableClientTLS = ConfigCfg.TLS
	walletCfg.Wallet.EnableServerTLS = ConfigCfg.TLS
	shellCfg.Wallet.EnableServerTLS = ConfigCfg.TLS
	ctlCfg.TLSSkipVerify = ConfigCfg.SkipVerify

	ctlCfg.Proxy = ConfigCfg.Proxy
	nodeCfg.Node.Proxy = ConfigCfg.Proxy
	walletCfg.Wallet.Proxy = ConfigCfg.Proxy
	shellCfg.Node.Proxy = ConfigCfg.Proxy
	shellCfg.Wallet.Proxy = ConfigCfg.Proxy

	ctlCfg.ProxyUser = ConfigCfg.ProxyUser
	nodeCfg.Node.ProxyUser = ConfigCfg.ProxyUser
	walletCfg.Wallet.ProxyUser = ConfigCfg.ProxyUser
	shellCfg.Node.ProxyUser = ConfigCfg.ProxyUser
	shellCfg.Wallet.ProxyUser = ConfigCfg.ProxyUser

	ctlCfg.ProxyPass = ConfigCfg.ProxyPass
	nodeCfg.Node.ProxyPass = ConfigCfg.ProxyPass
	walletCfg.Wallet.ProxyPass = ConfigCfg.ProxyPass
	shellCfg.Node.ProxyPass = ConfigCfg.ProxyPass
	shellCfg.Wallet.ProxyPass = ConfigCfg.ProxyPass

	walletCfg.Wallet.WalletPass = ConfigCfg.WalletPass
	shellCfg.Wallet.WalletPass = ConfigCfg.WalletPass

	var r string
	var ok bool
	var listeners []string
	if r, ok = getIfIs(ctx, "nodelistener"); ok {
		pu.NormalizeAddresses(r, node.DefaultPort, &listeners)
		ConfigCfg.NodeListeners = listeners
		nodeCfg.Node.Listeners = listeners
		shellCfg.Node.Listeners = listeners
	}
	if r, ok = getIfIs(ctx, "noderpclistener"); ok {
		pu.NormalizeAddresses(r, node.DefaultRPCPort, &listeners)
		ConfigCfg.NodeRPCListeners = listeners
		nodeCfg.Node.RPCListeners = listeners
		walletCfg.Wallet.RPCConnect = r
		shellCfg.Node.RPCListeners = listeners
		shellCfg.Wallet.RPCConnect = r
		ctlCfg.RPCServer = r
	}
	if r, ok = getIfIs(ctx, "walletlistener"); ok {
		pu.NormalizeAddresses(r, node.DefaultRPCPort, &listeners)
		ConfigCfg.WalletListeners = listeners
		walletCfg.Wallet.LegacyRPCListeners = listeners
		ctlCfg.Wallet = r
		shellCfg.Wallet.LegacyRPCListeners = listeners
	}
	if r, ok = getIfIs(ctx, "user"); ok {
		ConfigCfg.NodeUser = r
		nodeCfg.Node.RPCUser = r
		walletCfg.Wallet.PodUsername = r
		walletCfg.Wallet.Username = r
		shellCfg.Node.RPCUser = r
		shellCfg.Wallet.PodUsername = r
		shellCfg.Wallet.Username = r
		ctlCfg.RPCUser = r
	}
	if r, ok = getIfIs(ctx, "pass"); ok {
		ConfigCfg.NodePass = r
		nodeCfg.Node.RPCPass = r
		walletCfg.Wallet.PodPassword = r
		walletCfg.Wallet.Password = r
		shellCfg.Node.RPCPass = r
		shellCfg.Wallet.PodPassword = r
		shellCfg.Wallet.Password = r
		ctlCfg.RPCPass = r
	}
	if r, ok = getIfIs(ctx, "walletpass"); ok {
		ConfigCfg.WalletPass = r
		walletCfg.Wallet.WalletPass = ConfigCfg.WalletPass
		shellCfg.Wallet.WalletPass = ConfigCfg.WalletPass
	}

	if r, ok = getIfIs(ctx, "rpckey"); ok {
		r = n.CleanAndExpandPath(r)
		ConfigCfg.RPCKey = r
		nodeCfg.Node.RPCKey = r
		walletCfg.Wallet.RPCKey = r
		shellCfg.Node.RPCKey = r
		shellCfg.Wallet.RPCKey = r
	}
	if r, ok = getIfIs(ctx, "rpccert"); ok {
		r = n.CleanAndExpandPath(r)
		ConfigCfg.RPCCert = r
		nodeCfg.Node.RPCCert = r
		walletCfg.Wallet.RPCCert = r
		shellCfg.Node.RPCCert = r
		shellCfg.Wallet.RPCCert = r
	}
	if r, ok = getIfIs(ctx, "cafile"); ok {
		r = n.CleanAndExpandPath(r)
		ConfigCfg.CAFile = r
		walletCfg.Wallet.CAFile = r
		shellCfg.Wallet.CAFile = r
	}
	if r, ok = getIfIs(ctx, "tls"); ok {
		ConfigCfg.TLS = r == "true"
		nodeCfg.Node.TLS = ConfigCfg.TLS
		walletCfg.Wallet.EnableClientTLS = ConfigCfg.TLS
		shellCfg.Node.TLS = ConfigCfg.TLS
		shellCfg.Wallet.EnableClientTLS = ConfigCfg.TLS
		walletCfg.Wallet.EnableServerTLS = ConfigCfg.TLS
		shellCfg.Wallet.EnableServerTLS = ConfigCfg.TLS
	}
	if r, ok = getIfIs(ctx, "skipverify"); ok {
		ConfigCfg.SkipVerify = r == "true"
		ctlCfg.TLSSkipVerify = r == "true"
	}
	if r, ok = getIfIs(ctx, "proxy"); ok {
		pu.NormalizeAddresses(r, n.DefaultRPCPort, &listeners)
		ConfigCfg.Proxy = r
		ctlCfg.Proxy = ConfigCfg.Proxy
		nodeCfg.Node.Proxy = ConfigCfg.Proxy
		walletCfg.Wallet.Proxy = ConfigCfg.Proxy
		shellCfg.Node.Proxy = ConfigCfg.Proxy
		shellCfg.Wallet.Proxy = ConfigCfg.Proxy
	}
	if r, ok = getIfIs(ctx, "proxyuser"); ok {
		ConfigCfg.ProxyUser = r
		ctlCfg.ProxyUser = ConfigCfg.ProxyUser
		nodeCfg.Node.ProxyUser = ConfigCfg.ProxyUser
		walletCfg.Wallet.ProxyUser = ConfigCfg.ProxyUser
		shellCfg.Node.ProxyUser = ConfigCfg.ProxyUser
		shellCfg.Wallet.ProxyUser = ConfigCfg.ProxyUser
	}
	if r, ok = getIfIs(ctx, "proxypass"); ok {
		ConfigCfg.ProxyPass = r
		ctlCfg.ProxyPass = ConfigCfg.ProxyPass
		nodeCfg.Node.ProxyPass = ConfigCfg.ProxyPass
		walletCfg.Wallet.ProxyPass = ConfigCfg.ProxyPass
		shellCfg.Node.ProxyPass = ConfigCfg.ProxyPass
		shellCfg.Wallet.ProxyPass = ConfigCfg.ProxyPass
	}
	if r, ok = getIfIs(ctx, "network"); ok {
		r = strings.ToLower(r)
		switch r {
		case "mainnet", "testnet", "regtestnet", "simnet":
		default:
			r = "mainnet"
		}
		ConfigCfg.Network = r
		switch r {
		case "mainnet":
			ctlCfg.TestNet3 = false
			ctlCfg.SimNet = false
			nodeCfg.Node.TestNet3 = false
			nodeCfg.Node.SimNet = false
			nodeCfg.Node.RegressionTest = false
			walletCfg.Wallet.SimNet = false
			walletCfg.Wallet.TestNet3 = false
			shellCfg.Node.TestNet3 = false
			shellCfg.Node.RegressionTest = false
			shellCfg.Node.SimNet = false
			shellCfg.Wallet.TestNet3 = false
			shellCfg.Wallet.SimNet = false
		case "testnet":
			ctlCfg.TestNet3 = true
			ctlCfg.SimNet = false
			nodeCfg.Node.TestNet3 = true
			nodeCfg.Node.SimNet = false
			nodeCfg.Node.RegressionTest = false
			walletCfg.Wallet.SimNet = false
			walletCfg.Wallet.TestNet3 = true
			shellCfg.Node.TestNet3 = true
			shellCfg.Node.RegressionTest = false
			shellCfg.Node.SimNet = false
			shellCfg.Wallet.TestNet3 = true
			shellCfg.Wallet.SimNet = false
		case "regtestnet":
			ctlCfg.TestNet3 = false
			ctlCfg.SimNet = false
			nodeCfg.Node.TestNet3 = false
			nodeCfg.Node.SimNet = false
			nodeCfg.Node.RegressionTest = true
			walletCfg.Wallet.SimNet = false
			walletCfg.Wallet.TestNet3 = false
			shellCfg.Node.TestNet3 = false
			shellCfg.Node.RegressionTest = true
			shellCfg.Node.SimNet = false
			shellCfg.Wallet.TestNet3 = false
			shellCfg.Wallet.SimNet = false
		case "simnet":
			ctlCfg.TestNet3 = false
			ctlCfg.SimNet = true
			nodeCfg.Node.TestNet3 = false
			nodeCfg.Node.SimNet = true
			nodeCfg.Node.RegressionTest = false
			walletCfg.Wallet.SimNet = true
			walletCfg.Wallet.TestNet3 = false
			shellCfg.Node.TestNet3 = false
			shellCfg.Node.RegressionTest = false
			shellCfg.Node.SimNet = true
			shellCfg.Wallet.TestNet3 = false
			shellCfg.Wallet.SimNet = true
		}
	}
	WriteConfConfig(cfgFile, ConfigCfg)
	// Now write the configs for all the others reading them and overwriting the changed values
	ctl.WriteConfConfig(confs[0], &ctlCfg)
	node.WriteConfConfig(confs[1], &nodeCfg)
	walletrun.WriteConfConfig(confs[2], &walletCfg)
	shell.WriteConfConfig(confs[3], &shellCfg)
	if ctx.Is("show") {
		j, err := json.MarshalIndent(ConfigCfg, "", "  ")
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(string(j))
	}
}

// WriteConfConfig creates and writes the config file in the requested location
func WriteConfConfig(cfgFile string, cfg ConfConfig) {
	j, err := json.MarshalIndent(ConfigCfg, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	j = append(j, '\n')
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		panic(err.Error())
	}
}

// WriteDefaultConfConfig creates and writes a default config file in the requested location
func WriteDefaultConfConfig(cfgFile string) {
	defCfg := DefaultConfConfig()
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	j = append(j, '\n')
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		panic(err.Error())
	}
	// if we are writing default config we also want to use it
	ConfigCfg = *defCfg
}

func DefaultConfConfig() *ConfConfig {
	u := pu.GenKey()
	p := pu.GenKey()
	return &ConfConfig{
		NodeListeners:    []string{"127.0.0.1:11047"},
		NodeRPCListeners: []string{"127.0.0.1:11048"},
		WalletListeners:  []string{"127.0.0.1:11046"},
		NodeUser:         u,
		NodePass:         p,
		WalletPass:       wallet.InsecurePubPassphrase,
		RPCKey:           walletmain.DefaultRPCKeyFile,
		RPCCert:          walletmain.DefaultRPCCertFile,
		CAFile:           walletmain.DefaultCAFile,
		TLS:              false,
		SkipVerify:       false,
		Proxy:            "",
		ProxyUser:        "",
		ProxyPass:        "",
		Network:          "mainnet",
	}
}
