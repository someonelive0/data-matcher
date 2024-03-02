package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
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
	nc, err := utils.NatsConnect(*arg_server, *arg_user, *arg_password)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer nc.Close()
	log.Printf("connect %s success by user %s\n", *arg_server, *arg_user)

	msgch := make(chan *nats.Msg, *arg_chan_size)
	sub, err := utils.QueueSub2Chan(nc, *arg_subject, *arg_queue, msgch)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sub.Unsubscribe()
	log.Printf("sub %s success with queue %s\n", *arg_subject, *arg_queue)

	var wg sync.WaitGroup
	for i := 0; i < *arg_routines; i++ {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			matcher.Worker(name, msgch)
		}(strconv.Itoa(i))
	}

	wg.Wait()
}
