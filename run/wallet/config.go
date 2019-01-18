package wallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"git.parallelcoin.io/pod/lib/clog"
	n "git.parallelcoin.io/pod/module/wallet"
	"git.parallelcoin.io/pod/module/wallet/wallet"
	"git.parallelcoin.io/pod/run/logger"
	"github.com/tucnak/climax"
)

var log = clog.NewSubSystem("Wallet", clog.Ninf)

// Config is the default configuration native to ctl
var Config = new(n.Config)

// ConfigAndLog is the combined app and logging configuration data
type ConfigAndLog struct {
	Config *n.Config
	Levels map[string]string
}

// CombinedCfg is the combined app and log levels configuration
var CombinedCfg = ConfigAndLog{
	Config: Config,
	Levels: logger.Levels,
}

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "wallet",
	Brief: "parallelcoin wallet",
	Help:  "check balances, make payments, manage contacts",
	Flags: []climax.Flag{
		{
			Name:     "version",
			Short:    "V",
			Usage:    `--version`,
			Help:     `show version number and quit`,
			Variable: false,
		},
		{
			Name:     "configfile",
			Short:    "C",
			Usage:    "--configfile=/path/to/conf",
			Help:     "path to configuration file",
			Variable: true,
		},
		{
			Name:     "datadir",
			Short:    "D",
			Usage:    "--configfile=/path/to/conf",
			Help:     "path to configuration file",
			Variable: true,
		},
		{
			Name:     "init",
			Usage:    "--init",
			Help:     "resets configuration to defaults",
			Variable: false,
		},
		{
			Name:     "save",
			Usage:    "--save",
			Help:     "saves current configuration",
			Variable: false,
		},
		{
			Name:     "debuglevel",
			Short:    "d",
			Usage:    "--debuglevel=trace",
			Help:     "sets debuglevel, default info, sets the baseline for others not specified",
			Variable: true,
		},
		{
			Name:     "log-database",
			Usage:    "--log-database",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-txscript",
			Usage:    "--log-txscript",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-peer",
			Usage:    "--log-peer",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-netsync",
			Usage:    "--log-netsync",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-rpcclient",
			Usage:    "--log-rpcclient",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "addrmgr",
			Usage:    "--log-addrmgr",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-blockchain-indexers",
			Usage:    "--log-blockchain-indexers",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-blockchain",
			Usage:    "--log-blockchain",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-mining-cpuminer",
			Usage:    "--log-mining-cpuminer",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-mining",
			Usage:    "--log-mining",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-mining-controller",
			Usage:    "--log-mining-controller",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-connmgr",
			Usage:    "--log-connmgr",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-spv",
			Usage:    "--log-spv",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-node-mempool",
			Usage:    "--log-node-mempool",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-node",
			Usage:    "--log-node",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-wallet-wallet",
			Usage:    "--log-wallet-wallet",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-wallet-tx",
			Usage:    "--log-wallet-tx",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-wallet-votingpool",
			Usage:    "--log-wallet-votingpool",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-wallet",
			Usage:    "--log-wallet",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-wallet-chain",
			Usage:    "--log-wallet-chain",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-wallet-rpc-rpcserver",
			Usage:    "--log-wallet-rpc-rpcserver",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-wallet-rpc-legacyrpc",
			Usage:    "--log-wallet-rpc-legacyrpc",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "log-wallet-wtxmgr",
			Usage:    "--log-wallet-wtxmgr",
			Help:     "",
			Variable: true,
		},
		{
			Name:     "addpeers",
			Usage:    "--addpeers=some.peer.com:11047",
			Help:     "adds a peer to the peers database to try to connect to",
			Variable: true,
		},
		{
			Name:     "connectpeers",
			Usage:    "--connectpeers=some.peer.com:11047",
			Help:     "adds a peer to a connect-only whitelist",
			Variable: true,
		},
		{
			Name:     "disablelisten",
			Usage:    "--disablelisten=true",
			Help:     "disables the P2P listener",
			Variable: true,
		},
		{
			Name:     "listeners",
			Short:    "S",
			Usage:    "--listeners=127.0.0.1:11047",
			Help:     "sets an address to listen for P2P connections",
			Variable: true,
		},
		{
			Name:     "maxpeers",
			Usage:    "--maxpeers=100",
			Help:     "sets max number of peers to open connections to at once",
			Variable: true,
		},
		{
			Name:     "disablebanning",
			Usage:    "--disablebanning",
			Help:     "disable banning of misbehaving peers",
			Variable: false,
		},
		{
			Name:     "banduration",
			Usage:    "--banduration=1h",
			Help:     "how long to ban misbehaving peers - valid time units are {s, m, h},  minimum 1 second",
			Variable: true,
		},
		{
			Name:     "banthreshold",
			Usage:    "--banthreshold=100",
			Help:     "maximum allowed ban score before disconnecting and banning misbehaving peers",
			Variable: true,
		},
		{
			Name:     "whitelists",
			Usage:    "--whitelists=127.0.0.1:11047",
			Help:     "add an IP network or IP that will not be banned - eg. 192.168.1.0/24 or ::1",
			Variable: true,
		},
		{
			Name:     "rpcuser",
			Short:    "u",
			Usage:    "--rpcuser=username",
			Help:     "RPC username",
			Variable: true,
		},
		{
			Name:     "rpcpass",
			Short:    "P",
			Usage:    "--rpcpass=password",
			Help:     "RPC password",
			Variable: true,
		},
		{
			Name:     "rpclimituser",
			Short:    "u",
			Usage:    "--rpclimituser=username",
			Help:     "limited user RPC username",
			Variable: true,
		},
		{
			Name:     "rpclimitpass",
			Short:    "P",
			Usage:    "--rpclimitpass=password",
			Help:     "limited user RPC password",
			Variable: true,
		},
		{
			Name:     "rpclisteners",
			Short:    "s",
			Usage:    "--rpclisteners=127.0.0.1:11048",
			Help:     "RPC server to connect to",
			Variable: true,
		},
		{
			Name:     "rpccert",
			Short:    "c",
			Usage:    "--rpccert=/path/to/rpn.cert",
			Help:     "RPC server tls certificate chain for validation",
			Variable: true,
		},
		{
			Name:     "rpckey",
			Short:    "c",
			Usage:    "--rpccert=/path/to/rpn.key",
			Help:     "RPC server tls key for validation",
			Variable: true,
		},
		{
			Name:     "tls",
			Usage:    "--tls=false",
			Help:     "enable TLS",
			Variable: true,
		},
		{
			Name:     "disablednsseed",
			Usage:    "--disablednsseed=false",
			Help:     "disable dns seeding",
			Variable: true,
		},
		{
			Name:     "externalips",
			Usage:    "--externalips=192.168.0.1:11048",
			Help:     "set additional listeners on different address/interfaces",
			Variable: true,
		},
		{
			Name:     "proxy",
			Usage:    "--proxy 127.0.0.1:9050",
			Help:     "connect via SOCKS5 proxy (eg. 127.0.0.1:9050)",
			Variable: true,
		},
		{
			Name:     "proxyuser",
			Usage:    "--proxyuser username",
			Help:     "username for proxy server",
			Variable: true,
		},
		{
			Name:     "proxypass",
			Usage:    "--proxypass password",
			Help:     "password for proxy server",
			Variable: true,
		},
		{
			Name:     "onion",
			Usage:    "--onion 127.0.0.1:9050",
			Help:     "connect via onion proxy (eg. 127.0.0.1:9050)",
			Variable: true,
		},
		{
			Name:     "onionuser",
			Usage:    "--onionuser username",
			Help:     "username for onion proxy server",
			Variable: true,
		},
		{
			Name:     "onionpass",
			Usage:    "--onionpass password",
			Help:     "password for onion proxy server",
			Variable: true,
		},
		{
			Name:     "noonion",
			Usage:    "--noonion=true",
			Help:     "disable onion proxy",
			Variable: true,
		},
		{
			Name:     "torisolation",
			Usage:    "--torisolation=true",
			Help:     "enable tor stream isolation by randomising user credentials for each connection",
			Variable: true,
		},
		{
			Name:     "network",
			Usage:    "--network=mainnet",
			Help:     "connect to specified network: mainnet, testnet, regtestnet or simnet",
			Variable: true,
		},
		{
			Name:     "skipverify",
			Usage:    "--skipverify=false",
			Help:     "do not verify tls certificates (not recommended!)",
			Variable: true,
		},
		{
			Name:     "addcheckpoints",
			Usage:    "--addcheckpoints <height>:<hash>",
			Help:     "add custom checkpoints",
			Variable: true,
		},
		{
			Name:     "disablecheckpoints",
			Usage:    "--disablecheckpoints=true",
			Help:     "disable all checkpoints",
			Variable: true,
		},
		{
			Name:     "dbtype",
			Usage:    "--dbtype=ffldb",
			Help:     "set database backend type",
			Variable: true,
		},
		{
			Name:     "profile",
			Usage:    "--profile=127.0.0.1:3131",
			Help:     "start HTTP profiling server on given address",
			Variable: true,
		},
		{
			Name:     "cpuprofile",
			Usage:    "--cpuprofile=127.0.0.1:3232",
			Help:     "start cpu profiling server on given address",
			Variable: true,
		},
		{
			Name:     "upnp",
			Usage:    "--upnp=true",
			Help:     "enables the use of UPNP to establish inbound port redirections",
			Variable: true,
		},
		{
			Name:     "minrelaytxfee",
			Usage:    "--minrelaytxfee=1",
			Help:     "the minimum transaction fee in DUO/Kb to be considered a nonzero fee",
			Variable: true,
		},
		{
			Name:     "freetxrelaylimit",
			Usage:    "--freetxrelaylimit=100",
			Help:     "limit amount of free transactions relayed in thousand bytes per minute",
			Variable: true,
		},
		{
			Name:     "norelaypriority",
			Usage:    "--norelaypriority=true",
			Help:     "do not require free or low-fee transactions to have high priority for relaying",
			Variable: true,
		},
		{
			Name:     "trickleinterval",
			Usage:    "--trickleinterval=1",
			Help:     "time in seconds between attempts to send new inventory to a connected peer",
			Variable: true,
		},
		{
			Name:     "maxorphantxs",
			Usage:    "--maxorphantxs=100",
			Help:     "set maximum number of orphans transactions to keep in memory",
			Variable: true,
		},
		{
			Name:     "algo",
			Usage:    "--algo=random",
			Help:     "set algorithm to be used by cpu miner",
			Variable: true,
		},
		{
			Name:     "generate",
			Usage:    "--generate=true",
			Help:     "set CPU miner to generate blocks",
			Variable: true,
		},
		{
			Name:     "genthreads",
			Usage:    "--genthreads=-1",
			Help:     "set number of threads to generate using CPU, -1 = all available",
			Variable: true,
		},
		{
			Name:     "miningaddrs",
			Usage:    "--miningaddrs=aoeuaoe0760oeu0",
			Help:     "add an address to the list of addresses to make block payments to from miners",
			Variable: true,
		},
		{
			Name:     "minerlistener",
			Usage:    "--minerlistener=127.0.0.1:11011",
			Help:     "set the port for a miner work dispatch server to listen on",
			Variable: true,
		},
		{
			Name:     "minerpass",
			Usage:    "--minerpass=pa55word",
			Help:     "set the encryption password to prevent leaking or MiTM attacks on miners",
			Variable: true,
		},
		{
			Name:     "blockminsize",
			Usage:    "--blockminsize=80",
			Help:     "mininum block size in bytes to be used when creating a block",
			Variable: true,
		},
		{
			Name:     "blockmaxsize",
			Usage:    "--blockmaxsize=1024000",
			Help:     "maximum block size in bytes to be used when creating a block",
			Variable: true,
		},
		{
			Name:     "blockminweight",
			Usage:    "--blockminweight=500",
			Help:     "mininum block weight to be used when creating a block",
			Variable: true,
		},
		{
			Name:     "blockmaxweight",
			Usage:    "--blockmaxweight=10000",
			Help:     "maximum block weight to be used when creating a block",
			Variable: true,
		},
		{
			Name:     "blockprioritysize",
			Usage:    "--blockprioritysize=256",
			Help:     "size in bytes for high-priority/low-fee transactions when creating a block",
			Variable: true,
		},
		{
			Name:     "uacomment",
			Usage:    "--uacomment=joeblogsminers",
			Help:     "comment to add to the user agent - see BIP 14 for more information.",
			Variable: true,
		},
		{
			Name:     "nopeerbloomfilters",
			Usage:    "--nopeerbloomfilters=false",
			Help:     "disable bloom filtering support",
			Variable: true,
		},
		{
			Name:     "nocfilters",
			Usage:    "--nocfilters=false",
			Help:     "disable committed filtering (CF) support",
			Variable: true,
		},
		{
			Name:     "dropcfindex",
			Usage:    "--dropcfindex",
			Help:     "deletes the index used for committed filtering (CF) support from the database on start up and then exits",
			Variable: false,
		},
		{
			Name:     "sigcachemaxsize",
			Usage:    "--sigcachemaxsize=1000",
			Help:     "the maximum number of entries in the signature verification cache",
			Variable: true,
		},
		{
			Name:     "blocksonly",
			Usage:    "--blocksonly=true",
			Help:     "do not accept transactions from remote peers",
			Variable: true,
		},
		{
			Name:     "txindex",
			Usage:    "--txindex=true",
			Help:     "maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC",
			Variable: true,
		},
		{
			Name:     "droptxindex",
			Usage:    "--droptxindex",
			Help:     "deletes the hash-based transaction index from the database on start up and then exits.",
			Variable: false,
		},
		{
			Name:     "addrindex",
			Usage:    "--addrindex=true",
			Help:     "maintain a full address-based transaction index which makes the searchrawtransactions RPC available",
			Variable: true,
		},
		{
			Name:     "dropaddrindex",
			Usage:    "--dropaddrindex",
			Help:     "deletes the address-based transaction index from the database on start up and then exits",
			Variable: false,
		},
		{
			Name:     "relaynonstd",
			Usage:    "--relaynonstd=true",
			Help:     "relay non-standard transactions regardless of the default settings for the active network",
			Variable: true,
		},
		{
			Name:     "rejectnonstd",
			Usage:    "--rejectnonstd=false",
			Help:     "reject non-standard transactions regardless of the default settings for the active network",
			Variable: true,
		},
	},
	Examples: []climax.Example{
		{
			Usecase:     "--init --rpcuser=user --rpcpass=pa55word --save",
			Description: "resets the configuration file to default, sets rpc username and password and saves the changes to config after parsing",
		},
	},
	Handle: func(ctx climax.Context) int {
		var dl string
		var ok bool
		if dl, ok = ctx.Get("debuglevel"); ok {
			log.Tracef.Print("setting debug level %s", dl)
			log.SetLevel(dl)
			for i := range logger.Levels {
				logger.Levels[i] = dl
			}
		}
		log.Debugf.Print("node version %s", n.Version())
		if ctx.Is("version") {
			fmt.Println("node version", n.Version())
			clog.Shutdown()
		}
		log.Trace.Print("running command")

		var cfgFile string
		if cfgFile, ok = ctx.Get("configfile"); !ok {
			cfgFile = n.DefaultConfigFile
		}
		if ctx.Is("init") {
			log.Debugf.Print("writing default configuration to %s", cfgFile)
			writeDefaultConfig(cfgFile)
			writeLogCfgFile(Config.AppDataDir + "/logconf")
			configNode(&ctx, cfgFile)
		} else {
			log.Infof.Print("loading configuration from %s", cfgFile)
			if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
				log.Warn.Print("configuration file does not exist, creating new one")
				writeDefaultConfig(cfgFile)
				writeLogCfgFile(Config.AppDataDir + "/logconf")
				configNode(&ctx, cfgFile)
			} else {
				log.Debug.Print("reading app configuration from", cfgFile)
				cfgData, err := ioutil.ReadFile(cfgFile)
				if err != nil {
					log.Error.Print(err.Error())
					clog.Shutdown()
				}
				log.Tracef.Print("parsing app configuration\n%s", cfgData)
				err = json.Unmarshal(cfgData, CombinedCfg.Config)
				if err != nil {
					log.Error.Print(err.Error())
					clog.Shutdown()
				}
				logCfgFile := Config.AppDataDir + "/logconf"
				log.Debug.Print("reading logger configuration from", logCfgFile)
				logCfgData, err := ioutil.ReadFile(logCfgFile)
				if err != nil {
					log.Error.Print(err.Error())
					clog.Shutdown()
				}
				log.Tracef.Print("parsing logger configuration\n%s", logCfgData)
				err = json.Unmarshal(logCfgData, &CombinedCfg.Levels)
				if err != nil {
					log.Error.Print(err.Error())
					clog.Shutdown()
				}
				configNode(&ctx, cfgFile)
			}
		}
		runNode()
		clog.Shutdown()
		return 0
	},
}

func configNode(ctx *climax.Context, cfgFile string) {
	var err error
	// Apply all configurations specified on commandline
	if ctx.Is("datadir") {
		r, _ := ctx.Get("datadir")
		Config.DataDir = r
	}
	if ctx.Is("addpeers") {
		r, _ := ctx.Get("addpeers")
		Config.AddPeers = strings.Split(r, " ")
	}
	if ctx.Is("connectpeers") {
		r, _ := ctx.Get("connectpeers")
		Config.ConnectPeers = strings.Split(r, " ")
	}
	if ctx.Is("maxpeers") {
		r, _ := ctx.Get("maxpeers")
		Config.MaxPeers, err = strconv.Atoi(r)
		if err != nil {
			log.Error.Print(err.Error())
		}
	}
	if ctx.Is("banduration") {
		r, _ := ctx.Get("banduration")
		error := false
		var bd time.Duration
		switch r[len(r)-1] {
		case 's':
			ts, err := strconv.Atoi(r[:len(r)-1])
			error = err != nil
			bd = time.Duration(ts) * time.Second
		case 'm':
			tm, err := strconv.Atoi(r[:len(r)-1])
			error = err != nil
			bd = time.Duration(tm) * time.Minute
		case 'h':
			th, err := strconv.Atoi(r[:len(r)-1])
			error = err != nil
			bd = time.Duration(th) * time.Hour
		case 'd':
			td, err := strconv.Atoi(r[:len(r)-1])
			error = err != nil
			bd = time.Duration(td) * 24 * time.Hour
		}
		if error {
			log.Errorf.Print("malformed banduration `%s` leaving set at `%s` err: %s", r, Config.BanDuration, err.Error())
		}
		Config.BanDuration = bd
	}
	if ctx.Is("banthreshold") {
		r, _ := ctx.Get("banthreshold")
		bt, err := strconv.Atoi(r)
		if err != nil {
			log.Errorf.Print("malformed banthreshold `%s` leaving set at `%s` err: %s", r, Config.BanThreshold, err.Error())
		} else {
			Config.BanThreshold = uint32(bt)
		}
	}
	if ctx.Is("rpccert") {
		r, _ := ctx.Get("rpccert")
		Config.RPCCert = r
	}
	if ctx.Is("rpckey") {
		r, _ := ctx.Get("rpckey")
		Config.RPCKey = r
	}
	if ctx.Is("proxy") {
		r, _ := ctx.Get("proxy")
		Config.Proxy = r
	}
	if ctx.Is("proxyuser") {
		r, _ := ctx.Get("proxyuser")
		Config.ProxyUser = r
	}
	if ctx.Is("proxypass") {
		r, _ := ctx.Get("proxypass")
		Config.ProxyPass = r
	}
	if ctx.Is("network") {
		r, _ := ctx.Get("network")
		switch r {
		case "testnet":
			Config.TestNet3, Config.SimNet = true, false
		case "regtest":
			Config.TestNet3, Config.SimNet = false, false
		case "simnet":
			Config.TestNet3, Config.SimNet = false, true
		default:
			Config.TestNet3, Config.SimNet = false, false
		}
	}
	if ctx.Is("profile") {
		r, _ := ctx.Get("profile")
		Config.Profile = r
	}
	logger.SetLogging(ctx)
	if ctx.Is("save") {
		log.Infof.Print("saving config file to %s", cfgFile)
		j, err := json.MarshalIndent(Config, "", "  ")
		if err != nil {
			log.Error.Print(err.Error())
		}
		j = append(j, '\n')
		log.Tracef.Print("JSON formatted config file\n%s", j)
		ioutil.WriteFile(cfgFile, j, 0600)
		writeLogCfgFile(Config.DataDir + "/logconf")
	}
}

func writeLogCfgFile(logCfgFile string) {
	log.Info.Print("writing log configuration file", logCfgFile)
	j, err := json.MarshalIndent(logger.Levels, "", "  ")
	if err != nil {
		log.Error.Print(err.Error())
	}
	j = append(j, '\n')
	log.Tracef.Print("JSON formatted logging config file\n%s", j)
	ioutil.WriteFile(logCfgFile, j, 0600)

}
func writeDefaultConfig(cfgFile string) {
	defCfg := defaultConfig()
	defCfg.ConfigFile = cfgFile
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		log.Error.Print(err.Error())
	}
	j = append(j, '\n')
	log.Tracef.Print("JSON formatted config file\n%s", j)
	ioutil.WriteFile(cfgFile, j, 0600)
	// if we are writing default config we also want to use it
	Config = defCfg
}

func defaultConfig() *n.Config {
	return &n.Config{
		ConfigFile:             n.DefaultConfigFile,
		DataDir:                n.DefaultDataDir,
		AppDataDir:             n.DefaultAppDataDir,
		LogDir:                 n.DefaultLogDir,
		RPCKey:                 n.DefaultRPCKeyFile,
		RPCCert:                n.DefaultRPCCertFile,
		WalletPass:             wallet.InsecurePubPassphrase,
		CAFile:                 "",
		LegacyRPCMaxClients:    n.DefaultRPCMaxClients,
		LegacyRPCMaxWebsockets: n.DefaultRPCMaxWebsockets,
		AddPeers:               []string{},
		ConnectPeers:           []string{},
	}
}