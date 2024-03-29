package matcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"data-matcher/utils"
)

func (p *ManageApi) dumpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s := fmt.Sprintf(`{"run_time": "%s", "pid": %d`, p.Runtime.Format(time.RFC3339), os.Getpid())
	s += fmt.Sprintf(`, "process_mem": %s`, utils.ProcessMem())
	s += fmt.Sprintf(`, "channel": {"flow_channel": {"len": %d, "cap": %d}`, len(p.Flowch), cap(p.Flowch))
	s += fmt.Sprintf(`, "http_channel": {"len": %d, "cap": %d}`, len(p.Httpch), cap(p.Httpch))
	s += fmt.Sprintf(`, "label_http_channel": {"len": %d, "cap": %d}`, len(p.LabelHttpch), cap(p.LabelHttpch))
	s += fmt.Sprintf(`, "label_dns_channel": {"len": %d, "cap": %d} }`, len(p.LabelDnsch), cap(p.LabelDnsch))
	s += fmt.Sprintf(`, "inputer": %s`, p.Inputer.Dump())
	s += fmt.Sprintf(`, "outputer": %s`, p.Outputer.Dump())
	s += fmt.Sprintf(`, "post_worker": %s`, p.PostWorker.Dump())
	b, _ := json.Marshal(p.Workers)
	s += fmt.Sprintf(`, "worker": %s`, b)
	s += `, "rule": {`
	s += fmt.Sprintf(`"value_regex": { "number": %d, "matched": %d }, `, len(p.ValueRegs), 0)
	s += fmt.Sprintf(`"column_dict": { "number": %d, "matched": %d } }`, len(p.ColDicts), 0)
	s += "}"
	fmt.Fprint(w, s)
}

func (p *ManageApi) monitorHostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, utils.HostDump())
}

func (p *ManageApi) monitorHostLoadingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, utils.HostLoading())
}

func (p *ManageApi) statisticHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(p.Stats.Dump())
}

func (p *ManageApi) errorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(p.Myerrors.Dump())
}

func (p *ManageApi) configHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(p.Config.Dump())
}

func (p *ManageApi) configItemHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// vars := mux.Vars(r)
	// //log.Debugf("SyslogApi vars: %v", vars)

	// if len(vars) > 0 {
	// 	if title, ok := vars["title"]; ok {
	// 		if r.Method == "GET" {
	// 			if title == "kafka" {
	// 				toml.NewEncoder(w).Encode(p.Config.Kafkaconfig)
	// 			} else if title == "complete" {
	// 				toml.NewEncoder(w).Encode(p.Config.Completeconfig)
	// 			} else {
	// 				fmt.Fprintf(w, "wrong title: %s", title)
	// 			}
	// 		} else if r.Method == "POST" {
	// 			fmt.Fprintf(w, "not support post")
	// 		}
	// 	} else {
	// 		fmt.Fprintf(w, "not set title")
	// 	}
	// } else {
	// 	toml.NewEncoder(w).Encode(p.Config)
	// }
}

func (p *ManageApi) ruleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	if title, ok := vars["title"]; ok {
		if title == "value_regex" {
			jsonb, _ := json.MarshalIndent(p.ValueRegs, "", " ")
			w.Write(jsonb)
		} else if title == "culumn_dict" {
			jsonb, _ := json.MarshalIndent(p.ColDicts, "", " ")
			w.Write(jsonb)
		} else {
			fmt.Fprintf(w, "wrong title: %s", title)
		}
	} else {
		fmt.Fprintf(w, "not set title")
	}
}

func (p *ManageApi) ruleReloadHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("ManageApi receive reload policy command")
	w.WriteHeader(http.StatusOK)
	// test uploaded gather policy file, that is etc/user.xml.upload
	// filename := "etc/" + p.Config.Completeconfig.UserCompleteFile + ".upload"
	// if _, err := os.Stat(filename); err == nil {
	// 	c := NewUserCache(filename)
	// 	if err := c.LoadXML(); err != nil {
	// 		log.Errorf("ManageApi check wrong %s: %s", filename, err)
	// 		// maybe delete the uploaded file
	// 	} else {
	// 		log.Infof("ManageApi existed upload file %s, reload it", filename)
	// 		if err := os.Rename(filename, filename[:len(filename)-7]); err != nil {
	// 			log.Errorf("ManageApi rename %s failed: %s", filename, err)
	// 		}
	// 	}
	// }
	// filename = "etc/" + p.Config.Completeconfig.AccountCompleteFile + ".upload"
	// if _, err := os.Stat(filename); err == nil {
	// 	c := NewAccountCache(filename)
	// 	if err := c.LoadXML(); err != nil {
	// 		log.Errorf("ManageApi check wrong %s: %s", filename, err)
	// 		// maybe delete the uploaded file
	// 	} else {
	// 		log.Infof("ManageApi existed upload file %s, reload it", filename)
	// 		if err := os.Rename(filename, filename[:len(filename)-7]); err != nil {
	// 			log.Errorf("ManageApi rename %s failed: %s", filename, err)
	// 		}
	// 	}
	// }
	// filename = "etc/" + p.Config.Completeconfig.AssetCompleteFile + ".upload"
	// if _, err := os.Stat(filename); err == nil {
	// 	c := NewAssetCache(filename)
	// 	if err := c.LoadXML(); err != nil {
	// 		log.Errorf("ManageApi check wrong %s: %s", filename, err)
	// 		// maybe delete the uploaded file
	// 	} else {
	// 		log.Infof("ManageApi existed upload file %s, reload it", filename)
	// 		if err := os.Rename(filename, filename[:len(filename)-7]); err != nil {
	// 			log.Errorf("ManageApi rename %s failed: %s", filename, err)
	// 		}
	// 	}
	// }
	// filename = "etc/" + p.Config.Completeconfig.AuthCompleteFile + ".upload"
	// if _, err := os.Stat(filename); err == nil {
	// 	c := NewAuthCache(filename)
	// 	if err := c.LoadXML(); err != nil {
	// 		log.Errorf("ManageApi check wrong %s: %s", filename, err)
	// 		// maybe delete the uploaded file
	// 	} else {
	// 		log.Infof("ManageApi existed upload file %s, reload it", filename)
	// 		if err := os.Rename(filename, filename[:len(filename)-7]); err != nil {
	// 			log.Errorf("ManageApi rename %s failed: %s", filename, err)
	// 		}
	// 	}
	// }
	// filename = "etc/" + p.Config.Completeconfig.SysdevCompleteFile + ".upload"
	// if _, err := os.Stat(filename); err == nil {
	// 	c := NewSysdevCache(filename)
	// 	if err := c.LoadXML(); err != nil {
	// 		log.Errorf("ManageApi check wrong %s: %s", filename, err)
	// 		// maybe delete the uploaded file
	// 	} else {
	// 		log.Infof("ManageApi existed upload file %s, reload it", filename)
	// 		if err := os.Rename(filename, filename[:len(filename)-7]); err != nil {
	// 			log.Errorf("ManageApi rename %s failed: %s", filename, err)
	// 		}
	// 	}
	// }
	// filename = "etc/" + p.Config.Completeconfig.HolidayPolicyFile + ".upload"
	// if _, err := os.Stat(filename); err == nil {
	// 	c := NewHolidayPolicy(filename)
	// 	if err := c.LoadJson(); err != nil {
	// 		log.Errorf("ManageApi check wrong %s: %s", filename, err)
	// 		// maybe delete the uploaded file
	// 	} else {
	// 		log.Infof("ManageApi existed upload file %s, reload it", filename)
	// 		if err := os.Rename(filename, filename[:len(filename)-7]); err != nil {
	// 			log.Errorf("ManageApi rename %s failed: %s", filename, err)
	// 		}
	// 	}
	// }

	// s := "reload policy success"
	// if err := p.Usercache.LoadXML(); err != nil {
	// 	s = "Usercache.LoadXML failed: " + err.Error()
	// }
	// if err := p.Accountcache.LoadXML(); err != nil {
	// 	s += "Accountcache.LoadXML failed: " + err.Error()
	// }
	// if err := p.Assetpolicy.LoadXML(); err != nil {
	// 	s += "Assetpolicy.LoadXML failed: " + err.Error()
	// }
	// if err := p.Authcache.LoadXML(); err != nil {
	// 	s += "Authcache.LoadXML failed: " + err.Error()
	// }
	// if err := p.Sysdevcache.LoadXML(); err != nil {
	// 	s += "Sysdevcache.LoadXML failed: " + err.Error()
	// }
	// if err := p.Holidaypolicy.LoadJson(); err != nil {
	// 	s += "Holidaypolicy.LoadJson failed: " + err.Error()
	// }

	// fmt.Fprint(w, s)
}

func (p *ManageApi) DiscoverAppHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	i := 0
	fmt.Fprint(w, "[")
	p.PostWorker.AppMap.Range(func(key, value interface{}) bool {
		app, _ := json.Marshal(key.(string)) // 已经json编码后自带双引号
		num := value.(int)
		if i == 0 {
			fmt.Fprint(w, "\n")
		} else {
			fmt.Fprint(w, ",\n")
		}
		fmt.Fprintf(w, `{"app": %s, "num": %d}`, app, num)

		i++
		return true
	})
	fmt.Fprint(w, "\n]")
}

func (p *ManageApi) DiscoverApiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	i := 0
	fmt.Fprint(w, "[")
	p.PostWorker.ApiMap.Range(func(key, value interface{}) bool {
		api, _ := json.Marshal(key.(string)) // 已经json编码后自带双引号
		num := value.(int)
		if i == 0 {
			fmt.Fprint(w, "\n")
		} else {
			fmt.Fprint(w, ",\n")
		}
		fmt.Fprintf(w, `{"api": %s, "num": %d}`, api, num)

		i++
		return true
	})
	fmt.Fprint(w, "\n]")
}

func (p *ManageApi) DiscoverIpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	i := 0
	fmt.Fprint(w, "[")
	p.PostWorker.IpMap.Range(func(key, value interface{}) bool {
		ip, _ := json.Marshal(key.(string)) // 已经json编码后自带双引号
		num := value.(int)
		if i == 0 {
			fmt.Fprint(w, "\n")
		} else {
			fmt.Fprint(w, ",\n")
		}
		fmt.Fprintf(w, `{"ip": %s, "num": %d}`, ip, num)

		i++
		return true
	})
	fmt.Fprint(w, "\n]")
}
