package app

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"git.parallelcoin.io/pod/cmd/ctl"
	cl "git.parallelcoin.io/pod/pkg/util/clog"
	"github.com/tucnak/climax"
)

// CtlFlags is the list of flags and the default values stored in the Usage field
var CtlFlags = GetFlags(CtlCommand)

// DefaultCtlConfig returns an allocated, default CtlCfg
func DefaultCtlConfig(
	datadir string,
) *ctl.Config {
	return &ctl.Config{
		ConfigFile:    filepath.Join(datadir, "ctl/conf.json"),
		DebugLevel:    "off",
		RPCUser:       "user",
		RPCPass:       "pa55word",
		RPCServer:     ctl.DefaultRPCServer,
		RPCCert:       filepath.Join(datadir, "rpc.cert"),
		TLS:           false,
		Proxy:         "",
		ProxyUser:     "",
		ProxyPass:     "",
		TestNet3:      false,
		SimNet:        false,
		TLSSkipVerify: false,
		Wallet:        ctl.DefaultWallet,
	}
}

// WriteCtlConfig writes the current config in the requested location
func WriteCtlConfig(
	cc *ctl.Config,
) {

	j, err := json.MarshalIndent(cc, "", "  ")
	if err != nil {

		log <- cl.Err(err.Error())
	}
	j = append(j, '\n')
	log <- cl.Tracef{"JSON formatted config file\n%s", string(j)}
	EnsureDir(cc.ConfigFile)
	err = ioutil.WriteFile(cc.ConfigFile, j, 0600)
	if err != nil {

		log <- cl.Fatal{
			"unable to write config file %s",
			err.Error(),
		}
		cl.Shutdown()
	}
}

// WriteDefaultCtlConfig writes a default config in the requested location
func WriteDefaultCtlConfig(
	datadir string,
) {

	defCfg := DefaultCtlConfig(datadir)
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {

		log <- cl.Err(err.Error())
	}
	j = append(j, '\n')
	log <- cl.Tracef{"JSON formatted config file\n%s", string(j)}
	EnsureDir(defCfg.ConfigFile)
	err = ioutil.WriteFile(defCfg.ConfigFile, j, 0600)
	if err != nil {

		log <- cl.Fatal{
			"unable to write config file %s",
			err.Error(),
		}
		cl.Shutdown()
	}

	// if we are writing default config we also want to use it
	CtlCfg = defCfg
}

func configCtl(
	ctx *climax.Context,
	cfgFile string,
) {

	var r string
	var ok bool

	// Apply all configurations specified on commandline
	if r, ok = getIfIs(ctx, "debuglevel"); ok {

		CtlCfg.DebugLevel = r
		log <- cl.Trace{
			"set", "debuglevel", "to", r,
		}
	}
	if r, ok = getIfIs(ctx, "rpcuser"); ok {

		CtlCfg.RPCUser = r
		log <- cl.Tracef{
			"set %s to %s", "rpcuser", r,
		}
	}
	if r, ok = getIfIs(ctx, "rpcpass"); ok {

		CtlCfg.RPCPass = r
		log <- cl.Tracef{
			"set %s to %s", "rpcpass", r,
		}
	}
	if r, ok = getIfIs(ctx, "rpcserver"); ok {

		CtlCfg.RPCServer = r
		log <- cl.Tracef{
			"set %s to %s", "rpcserver", r,
		}
	}
	if r, ok = getIfIs(ctx, "rpccert"); ok {

		CtlCfg.RPCCert = r
		log <- cl.Tracef{"set %s to %s", "rpccert", r}
	}
	if r, ok = getIfIs(ctx, "tls"); ok {

		CtlCfg.TLS = r == "true"
		log <- cl.Tracef{"set %s to %s", "tls", r}
	}
	if r, ok = getIfIs(ctx, "proxy"); ok {

		CtlCfg.Proxy = r
		log <- cl.Tracef{"set %s to %s", "proxy", r}
	}
	if r, ok = getIfIs(ctx, "proxyuser"); ok {

		CtlCfg.ProxyUser = r
		log <- cl.Tracef{"set %s to %s", "proxyuser", r}
	}
	if r, ok = getIfIs(ctx, "proxypass"); ok {

		CtlCfg.ProxyPass = r
		log <- cl.Tracef{"set %s to %s", "proxypass", r}
	}
	otn, osn := "false", "false"
	if CtlCfg.TestNet3 {

		otn = "true"
	}
	if CtlCfg.SimNet {

		osn = "true"
	}
	tn, ts := ctx.Get("testnet")
	sn, ss := ctx.Get("simnet")
	if ts {

		CtlCfg.TestNet3 = tn == "true"
	}
	if ss {

		CtlCfg.SimNet = sn == "true"
	}
	if CtlCfg.TestNet3 && CtlCfg.SimNet {

		log <- cl.Error{
			"cannot enable simnet and testnet at the same time. current settings testnet =", otn,
			"simnet =", osn,
		}
	}
	if ctx.Is("skipverify") {

		CtlCfg.TLSSkipVerify = true
		log <- cl.Tracef{
			"set %s to true", "skipverify",
		}
	}
	if ctx.Is("wallet") {

		CtlCfg.RPCServer = CtlCfg.Wallet
		log <- cl.Trc("using configured wallet rpc server")
	}
	if r, ok = getIfIs(ctx, "walletrpc"); ok {

		CtlCfg.Wallet = r
		log <- cl.Tracef{
			"set %s to true", "walletrpc",
		}
	}
	if ctx.Is("save") {

		log <- cl.Info{
			"saving config file to",
			cfgFile,
		}
		j, err := json.MarshalIndent(CtlCfg, "", "  ")
		if err != nil {

			log <- cl.Err(err.Error())
		}
		j = append(j, '\n')
		log <- cl.Trace{
			"JSON formatted config file\n", string(j),
		}
		e := ioutil.WriteFile(cfgFile, j, 0600)
		if e != nil {
			log <- cl.Error{
				"error writing configuration file:", e,
			}
		}
	}
}
