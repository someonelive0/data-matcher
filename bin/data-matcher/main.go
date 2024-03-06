package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"

	"data-matcher/matcher"
	"data-matcher/utils"
)

var (
	arg_debug   = flag.Bool("D", false, "debug")
	arg_version = flag.Bool("v", false, "version")
	arg_config  = flag.String("f", "etc/data-matcher.yaml", "config filename")
	START_TIME  = time.Now()
)

func init() {
	flag.Parse()
	if *arg_version {
		fmt.Print(utils.IDSS_BANNER)
		fmt.Printf("%s\n", utils.SERVICE_VERSION)
		os.Exit(0)
	}
	utils.ShowBannerForApp("data-matcher", utils.SERVICE_VERSION, utils.BUILD_TIME)
	utils.Chdir2PrgPath()
	fmt.Println("pwd:", utils.GetPrgDir())
	utils.InitLog("data-matcher.log", *arg_debug)
	log.Infof("BEGIN... %v, config=%v, debug=%v",
		START_TIME.Format(time.DateTime), *arg_config, *arg_debug)
}

func main() {
	var myerrors = utils.AddLogHook()

	// load config
	var myconfig, err = matcher.LoadConfig(*arg_config)
	if err != nil {
		log.Errorf("loadConfig error %s", err)
		os.Exit(1)
	}
	log.Infof("myconfig: %s", myconfig.Dump())

	runok := true // Exit when run is not ok
	msgch := make(chan *nats.Msg, myconfig.ChannelSize)
	var inputer = matcher.Inputer{}
	err = inputer.Run(msgch, strings.Join(myconfig.NatsConfig.Servers, ","),
		myconfig.NatsConfig.User, myconfig.NatsConfig.Password, myconfig.HttpFlow.Subject, myconfig.NatsConfig.QueueName)
	if err != nil {
		return
	}

	// run workers
	var wg sync.WaitGroup
	for i := 0; i < myconfig.Workers; i++ {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			matcher.Worker(name, msgch)
		}(strconv.Itoa(i))
	}

	// run api interface
	var myapi = matcher.ManageApi{
		RestapiHandler: utils.RestapiHandler{Name: "data-matcher", Runtime: START_TIME},
		Port:           myconfig.ManagePort,
		Config:         myconfig,
		MsgChan:        msgch,
		Myerrors:       myerrors,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := myapi.Run(); err != nil {
			runok = false
		}
	}()

	// Wait for signal and timer
	var signchan = make(chan os.Signal, 5)
	signal.Notify(signchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	var ticker = time.NewTicker(time.Second * 1)

	// GracefullExit, Wait and stop all
	gracefullExit := func() {
		log.Info("GracefullExit")
		ticker.Stop()
		signal.Stop(signchan)
		close(signchan)

		inputer.Stop()
		// stlogOutput.Stop()
		// apilogOutput.Stop()
		// waitChanEmpty(chan_stlog_0, chan_stlog_1)
		myapi.Stop()
		close(msgch)
		wg.Wait() // Wait workers
	}

	// Waitting... signal and timer. block here.
	for {
		select {
		case <-ticker.C:
			if !runok {
				gracefullExit()
			}
		case s, ok := <-signchan:
			if !ok {
				goto END
			}
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Info("Receive SIGNAL: ", s)
				gracefullExit()
			default:
				log.Info("Receive other signal, ignore", s)
			}
		}
	}

END:
	log.Info("END")
}