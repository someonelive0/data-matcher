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
	"data-matcher/model"
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
		fmt.Printf("%s\n", utils.APP_VERSION)
		os.Exit(0)
	}
	utils.ShowBannerForApp("data-matcher", utils.APP_VERSION, utils.BUILD_TIME)
	utils.Chdir2PrgPath()
	pwd, _ := utils.GetPrgDir()
	fmt.Println("pwd:", pwd)
	if err := utils.InitLog("data-matcher.log", *arg_debug); err != nil {
		fmt.Printf("init log failed: %s\n", err)
		os.Exit(1)
	}
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
	runok := true                                                       // Exit when run is not ok
	var flowch = make(chan *nats.Msg, myconfig.ChannelSize)             // input channel
	var httpch = make(chan *model.FlowHttp, myconfig.ChannelSize)       // to post worker
	var label_httpch = make(chan *model.FlowHttp, myconfig.ChannelSize) // to outputer
	var label_dnsch = make(chan *model.FlowDns, myconfig.ChannelSize)   // to outputer
	var stats = matcher.NewMyStatistic(START_TIME)

	var inputer = matcher.Inputer{ // http flow inputer, 如有多个flow要输入，则建立多个inputer
		Flowch:     flowch,
		NatsConfig: &myconfig.NatsConfig,
		Flow:       &myconfig.Flow,
		Stats:      stats,
	}
	if err = inputer.Run(); err != nil {
		return
	}

	// run outputer
	var wg sync.WaitGroup
	var outputer = matcher.Outputer{
		LabelHttpch: label_httpch,
		LabelDnsch:  label_dnsch,
		NatsConfig:  &myconfig.NatsConfig,
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
	var appmap sync.Map
	var apimap sync.Map
	var ipmap sync.Map
	var wgWokers sync.WaitGroup // 单独等待workers的结束
	var workers []*matcher.Worker = make([]*matcher.Worker, 0)
	for i := 0; i < myconfig.Workers; i++ {
		worker := &matcher.Worker{
			Name:        strconv.Itoa(i),
			Flowch:      flowch,
			Httpch:      httpch,
			LabelHttpch: label_httpch,
			LabelDnsch:  label_dnsch,
			ValueRegs:   regs,
			ColDicts:    dicts,
			Appmap:      &appmap,
			Apimap:      &apimap,
			Ipmap:       &ipmap,
		}
		if err = worker.Init(); err != nil {
			log.Errorf("worker init failed %s", err)
			return
		}
		workers = append(workers, worker)
	}

	for i := 0; i < myconfig.Workers; i++ {
		wgWokers.Add(1)
		go func(i int) {
			defer wgWokers.Done()
			workers[i].Run()
		}(i)
	}

	// run post-worker
	var post_worker = matcher.PostWorker{
		Httpch:     httpch,
		NatsConfig: &myconfig.NatsConfig,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = post_worker.Run(); err != nil {
			runok = false
		}
	}()

	// run manage api
	var myapi = matcher.ManageApi{
		RestapiHandler: utils.RestapiHandler{Name: "data-matcher", Runtime: START_TIME},
		Port:           myconfig.ManagePort,
		Config:         myconfig,
		Stats:          stats,
		Flowch:         flowch,
		Httpch:         httpch,
		LabelHttpch:    label_httpch,
		LabelDnsch:     label_dnsch,
		Inputer:        &inputer,
		Outputer:       &outputer,
		Workers:        workers,
		PostWorker:     &post_worker,
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
		close(flowch)
		wgWokers.Wait()
		close(httpch)
		close(label_httpch)
		close(label_dnsch)
		post_worker.Stop()
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
