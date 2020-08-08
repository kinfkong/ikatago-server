package nat

import (
	"context"
	"errors"
	"fmt"
	syslog "log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fatedier/beego/utils"
	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/models/auth"
	"github.com/fatedier/frp/models/config"
	"github.com/fatedier/frp/utils/log"
	"github.com/fatedier/frp/utils/version"
	"github.com/fatedier/golib/crypto"
	myconfig "github.com/kinfkong/ikatago-server/config"
	"github.com/spf13/cobra"

	// import the blank assets
	_ "github.com/kinfkong/ikatago-server/nat/assets/frpc/statik"
)

const (
	CfgFileTypeIni = iota
	CfgFileTypeCmd
)

var (
	cfgFile     string
	showVersion bool

	serverAddr      string
	user            string
	protocol        string
	token           string
	logLevel        string
	logFile         string
	logMaxDays      int
	disableLogColor bool

	proxyName         string
	localIp           string
	localPort         int
	remotePort        int
	useEncryption     bool
	useCompression    bool
	customDomains     string
	subDomain         string
	httpUser          string
	httpPwd           string
	locations         string
	hostHeaderRewrite string
	role              string
	sk                string
	multiplexer       string
	serverName        string
	bindAddr          string
	bindPort          int

	kcpDoneCh chan struct{}

	runningService *client.Service
)

// do some stuff
var _ = func() error {
	os.Setenv("KNAT_SERVER_ADDR", "120.53.123.43")
	os.Setenv("KNAT_SERVER_PORT", "7000")
	os.Setenv("KNAT_SERVER_TOKEN", "kinfkong")
	return nil
}()

func InitFRP() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./frpc.ini", "config file of frpc")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of frpc")

	kcpDoneCh = make(chan struct{})
}

var rootCmd = &cobra.Command{
	Use:   "frpc",
	Short: "frpc is the client of frp (https://github.com/fatedier/frp)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}

		// Do not show command usage here.
		err := runClient(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func handleSignal(svr *client.Service) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.Close()
	time.Sleep(250 * time.Millisecond)
	close(kcpDoneCh)
}

func parseClientCommonCfg(fileType int, content string) (cfg config.ClientCommonConf, err error) {
	if fileType == CfgFileTypeIni {
		cfg, err = parseClientCommonCfgFromIni(content)
	} else if fileType == CfgFileTypeCmd {
		cfg, err = parseClientCommonCfgFromCmd()
	}
	if err != nil {
		return
	}

	err = cfg.Check()
	if err != nil {
		return
	}
	return
}

func parseClientCommonCfgFromIni(content string) (config.ClientCommonConf, error) {
	cfg, err := config.UnmarshalClientConfFromIni(content)
	if err != nil {
		return config.ClientCommonConf{}, err
	}
	return cfg, err
}

func parseClientCommonCfgFromCmd() (cfg config.ClientCommonConf, err error) {
	cfg = config.GetDefaultClientConf()

	strs := strings.Split(serverAddr, ":")
	if len(strs) < 2 {
		err = fmt.Errorf("invalid server_addr")
		return
	}
	if strs[0] != "" {
		cfg.ServerAddr = strs[0]
	}
	cfg.ServerPort, err = strconv.Atoi(strs[1])
	if err != nil {
		err = fmt.Errorf("invalid server_addr")
		return
	}

	cfg.User = user
	cfg.Protocol = protocol
	cfg.LogLevel = logLevel
	cfg.LogFile = logFile
	cfg.LogMaxDays = int64(logMaxDays)
	if logFile == "console" {
		cfg.LogWay = "console"
	} else {
		cfg.LogWay = "file"
	}
	cfg.DisableLogColor = disableLogColor

	// Only token authentication is supported in cmd mode
	cfg.AuthClientConfig = auth.GetDefaultAuthClientConf()
	cfg.Token = token

	return
}

func runClient(cfgFilePath string) (err error) {
	var content string
	content, err = config.GetRenderedConfFromFile(cfgFilePath)
	if err != nil {
		return
	}

	cfg, err := parseClientCommonCfg(CfgFileTypeIni, content)
	if err != nil {
		return
	}

	pxyCfgs, visitorCfgs, err := config.LoadAllConfFromIni(cfg.User, content, cfg.Start)
	if err != nil {
		return err
	}

	err = startService(cfg, pxyCfgs, visitorCfgs, cfgFilePath)
	return
}

func startService(cfg config.ClientCommonConf, pxyCfgs map[string]config.ProxyConf, visitorCfgs map[string]config.VisitorConf, cfgFile string) (err error) {
	log.InitLog(cfg.LogWay, cfg.LogFile, cfg.LogLevel,
		cfg.LogMaxDays, cfg.DisableLogColor)

	if cfg.DnsServer != "" {
		s := cfg.DnsServer
		if !strings.Contains(s, ":") {
			s += ":53"
		}
		// Change default dns server for frpc
		net.DefaultResolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("udp", s)
			},
		}
	}
	svr, errRet := client.NewService(cfg, pxyCfgs, visitorCfgs, cfgFile)
	if errRet != nil {
		err = errRet
		return
	}
	runningService = svr
	// Capture the exit signal if we use kcp.
	if cfg.Protocol == "kcp" {
		go handleSignal(svr)
	}

	err = svr.Run()
	if cfg.Protocol == "kcp" {
		<-kcpDoneCh
	}
	return
}

// FRP represents the frp
type FRP struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	configFile string
}

// InitWithConfig inits from the config
func (frp *FRP) InitWithConfig(configObject map[string]interface{}) error {
	frpConfigFile, ok := configObject["config_file"]
	if !ok {
		syslog.Printf("ERROR minssing config_file\n")
		return errors.New("missing_config_File")
	}

	// check file exists
	if !utils.FileExists(frpConfigFile.(string)) {
		syslog.Printf("ERROR config file not found")
		return errors.New("file_not_found")
	}
	frp.configFile = frpConfigFile.(string)

	crypto.DefaultSalt = "frp"
	rand.Seed(time.Now().UnixNano())
	InitFRP()

	return nil
}

// RunAsync runs the frp
func (frp *FRP) RunAsync() error {

	go func() {
		err := runClient(frp.configFile)
		if err != nil {
			syslog.Fatal(err)
		}
	}()

	return frp.waitUntilReady(60)
}

func (frp *FRP) waitUntilReady(timeout int) error {
	ready := false
	endTime := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		if time.Now().After(endTime) {
			break
		}
		if runningService != nil {
			ctl := runningService.GetController()
			if ctl == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			pm := runningService.GetController().GetProxyManager()
			if pm == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			pss := pm.GetAllProxyStatus()
			for _, ps := range pss {
				// check if this is the target proxy
				// if ps.Cfg.GetBaseInfo()
				basicInfo := ps.Cfg.GetBaseInfo()
				if basicInfo == nil {
					continue
				}
				serverListenPort, _ := myconfig.GetServerListenPort()
				if basicInfo.LocalPort == serverListenPort {
					// found, check status
					if ps.Status == "running" {
						// everything is done
						// get the remote port
						host := strings.TrimSpace(strings.Split(ps.RemoteAddr, ":")[0])
						portString := strings.Split(ps.RemoteAddr, ":")[1]
						port, err := strconv.Atoi(portString)

						if err != nil {
							continue
						}
						if len(host) == 0 {
							host = runningService.GetClientCommonConf().ServerAddr
						}
						frp.Port = port
						frp.Host = host
						ready = true
					}
				}
			}
		}
		if ready {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if ready {
		return nil
	}
	syslog.Printf("ERROR cannot connect to the frp server\n")
	return errors.New("timeout")
}

// GetInfo gets the ssh info
func (frp *FRP) GetInfo() (Info, error) {
	return Info{
		Host: frp.Host,
		Port: frp.Port,
	}, nil
}
