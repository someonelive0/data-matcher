package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"

	"data-matcher/matcher"
	"data-matcher/utils"
)

var (
	param_debug   = flag.Bool("D", false, "debug")
	param_version = flag.Bool("v", false, "version")
	param_config  = flag.String("f", "etc/data-matcher.yaml", "config filename")
	arg_server    = flag.String("server", "nats://localhost:4222", "nats server hostname or ip")
	arg_user      = flag.String("user", "", "nats user")
	arg_password  = flag.String("password", "", "")
	arg_subject   = flag.String("subject", "foo", "subject of nats to receive message")
	arg_queue     = flag.String("queue", "foo", "queue name for the subject")
	arg_chan_size = flag.Int("chan-size", 1000000, "channel size to cache message")
	arg_routines  = flag.Int("routines", 8, "routines to process message")
	START_TIME    = time.Now()
)

func init() {
	flag.Parse()
	if *param_version {
		fmt.Print(utils.IDSS_BANNER)
		fmt.Printf("%s\n", utils.SERVICE_VERSION)
		os.Exit(0)
	}
	utils.ShowBannerForApp("data-matcher", utils.SERVICE_VERSION, utils.BUILD_TIME)
	utils.Chdir2PrgPath()
	fmt.Println("pwd:", utils.GetPrgDir())
	utils.InitLog("data-matcher.log", *param_debug)
	log.Infof("BEGIN... %v, config=%v, debug=%v",
		START_TIME.Format("2006-01-02 15:04:05"), *param_config, *param_debug)
}

func main() {
	var myerrors = utils.AddLogHook()

	// load config
	// var myconfig, err = matcher.LoadConfig(*param_config)
	// if err != nil {
	// 	log.Errorf("loadConfig error %s", err)
	// 	os.Exit(1)
	// }
	// log.Infof("MyConfig: %s", myconfig.Dump())

	runok := true // Exit when run is not ok
	msgch := make(chan *nats.Msg, *arg_chan_size)
	var inputer = matcher.Inputer{}
	err := inputer.Run(msgch, *arg_server, *arg_user, *arg_password, *arg_subject, *arg_queue)
	if err != nil {
		return
	}

	// run workers
	var wg sync.WaitGroup
	for i := 0; i < *arg_routines; i++ {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			matcher.Worker(name, msgch)
		}(strconv.Itoa(i))
	}

	wg.Wait()

	// run api interface
	var myapi = matcher.ManageApi{
		RestapiHandler: utils.RestapiHandler{Name: "data-matcher", Runtime: START_TIME},
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
		// msgInput.Stop()
		// stlogOutput.Stop()
		// apilogOutput.Stop()
		// waitChanEmpty(chan_stlog_0, chan_stlog_1)
		myapi.Stop()
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
