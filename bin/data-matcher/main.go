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

	"data-matcher/engine"
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

	// load rules
	rulesConf, err := engine.NewRulesConfig(myconfig.RulesFile)
	if err != nil {
		log.Errorf("load rules config [%s] error %s", myconfig.RulesFile, err)
		os.Exit(1)
	}
	regs := rulesConf.GetValueReg()
	log.Infof("load value regex in rules config number %d", len(regs))
	dicts := rulesConf.GetColDict()
	log.Infof("load column dicts in rules config number %d", len(dicts))

	// run inputer, receive nats msg to channel
	runok := true // Exit when run is not ok
	var msgch = make(chan *nats.Msg, myconfig.ChannelSize)
	var outch = make(chan *nats.Msg, myconfig.ChannelSize)
	var stats = matcher.NewMyStatistic(START_TIME)
	var inputer = matcher.Inputer{ // http flow inputer, 如有多个flow要输入，则建立多个inputer
		Msgch:      msgch,
		NatsConfig: &myconfig.NatsConfig,
		HttpFlow:   &myconfig.HttpFlow,
		Stats:      stats,
	}
	if err = inputer.Run(); err != nil {
		return
	}

	// run outputer
	var wg sync.WaitGroup
	var outputer = matcher.Outputer{
		Outch:       outch,
		NatsConfig:  &myconfig.NatsConfig,
		HttpOutFlow: &myconfig.HttpFlow,
		Stats:       stats,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = outputer.Run(); err != nil {
			runok = false
		}
	}()

	// run workers
	var workers []*matcher.Worker = make([]*matcher.Worker, 0)
	for i := 0; i < myconfig.Workers; i++ {
		worker := &matcher.Worker{
			Name:      strconv.Itoa(i),
			Msgch:     msgch,
			Outch:     outch,
			ValueRegs: regs,
			ColDicts:  dicts,
		}
		if err = worker.Init(); err != nil {
			log.Errorf("worker init failed %s", err)
			return
		}
		workers = append(workers, worker)
	}

	for i := 0; i < myconfig.Workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			workers[i].Run()
		}(i)
	}

	// run manage api
	var myapi = matcher.ManageApi{
		RestapiHandler: utils.RestapiHandler{Name: "data-matcher", Runtime: START_TIME},
		Port:           myconfig.ManagePort,
		Config:         myconfig,
		Stats:          stats,
		Msgch:          msgch,
		Outch:          outch,
		Inputer:        &inputer,
		Outputer:       &outputer,
		Workers:        workers,
		ValueRegs:      regs,
		ColDicts:       dicts,
		Myerrors:       myerrors,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = myapi.Run(); err != nil {
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
		close(msgch)
		close(outch)
		outputer.Stop()
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
