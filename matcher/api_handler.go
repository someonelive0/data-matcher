package matcher

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"data-matcher/utils"
)

func (p *MatcherApi) dumpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s := fmt.Sprintf(`{"run_time": "%s", "pid": %d, `, p.Runtime.Format(time.RFC3339), os.Getpid())
	s += fmt.Sprintf(`"process_mem": %s, `, utils.ProcessMem())
	// s += fmt.Sprintf(`"input": {"pkts": %d, "tps": %d, "max_tps": %d, "avg_tps": %d, `,
	// 	p.Input.pkts, p.Input.tps, p.Input.max_tps, p.Input.avg_tps)
	// s += fmt.Sprintf(`"stlog_partition": %d, "stlog_offset": %d, "apilog_partition": %d, "apilog_offset": %d, `,
	// 	p.Input.stlog_partition, p.Input.stlog_offset, p.Input.apilog_partition, p.Input.apilog_offset)
	// s += fmt.Sprintf(`"topics": "%s", "invalid_pkts": %d}, `, p.Input.Topics, p.Input.invalid_pkts)
	// s += fmt.Sprintf(`"standard_log_output": {"count": %d, "success": %d, "failed": %d, "partition": %d, "offset": %d, "topics": "%s"}, `,
	// 	p.Stlogoutput.Count, p.Stlogoutput.success, p.Stlogoutput.failed,
	// 	p.Stlogoutput.Topicpartition, p.Stlogoutput.Topicoffset, p.Stlogoutput.Topics)
	// s += fmt.Sprintf(`"standard_log_faildisk": %s, `, p.Stlogoutput.faildisk.Dump())
	// s += fmt.Sprintf(`"api_log_output": {"count": %d, "success": %d, "failed": %d, "partition": %d, "offset": %d, "topics": "%s"}, `,
	// 	p.Apilogoutput.Count, p.Apilogoutput.success, p.Apilogoutput.failed,
	// 	p.Apilogoutput.Topicpartition, p.Apilogoutput.Topicoffset, p.Apilogoutput.Topics)
	// s += fmt.Sprintf(`"api_log_faildisk": %s, `, p.Apilogoutput.faildisk.Dump())
	// s += `"policy": {`
	// s += fmt.Sprintf(`"user_policy": { "number": %d, "matched": %d }, `,
	// 	p.Usercache.Number, p.Usercache.Matched)
	// s += fmt.Sprintf(`"account_policy": { "number": %d, "matched": %d }, `,
	// 	p.Accountcache.Number, p.Accountcache.Matched)
	// s += fmt.Sprintf(`"asset_policy": { "number": %d, "matched": %d }, `,
	// 	p.Assetpolicy.Number, p.Assetpolicy.Matched)
	// s += fmt.Sprintf(`"auth_policy": { "number": %d, "matched": %d }, `,
	// 	p.Authcache.Number, p.Authcache.Matched)
	// s += fmt.Sprintf(`"sysdev_policy": { "number": %d, "matched": %d } }, `,
	// 	p.Sysdevcache.Number, p.Sysdevcache.Matched)
	// s += fmt.Sprintf(`"channel": {"stlog_channel_in": {"len": %d, "cap": %d}, `, len(p.Chan_stlog_0), cap(p.Chan_stlog_0))
	// s += fmt.Sprintf(`"stlog_channel_out": {"len": %d, "cap": %d}, `, len(p.Chan_stlog_1), cap(p.Chan_stlog_1))
	// s += fmt.Sprintf(`"apilog_channel_in": {"len": %d, "cap": %d}, `, len(p.Chan_apilog_0), cap(p.Chan_apilog_0))
	// s += fmt.Sprintf(`"apilog_channel_out": {"len": %d, "cap": %d} } }`, len(p.Chan_apilog_1), cap(p.Chan_apilog_1))
	fmt.Fprint(w, s)
}

func (p *MatcherApi) statisticHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"standard_log_v4": `)
	// w.Write(p.Stats_stlog.Dump())
	// fmt.Fprintf(w, `, "api_log_v4": `)
	// w.Write(p.Stats_apilog.Dump())
	fmt.Fprintf(w, `}`)
}

func (p *MatcherApi) errorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(p.Myerrors.Dump())
}

func (p *MatcherApi) configHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// if r.Method == "GET" {
	// 	toml.NewEncoder(w).Encode(p.Config)
	// }
}

func (p *MatcherApi) configItemHandler(w http.ResponseWriter, r *http.Request) {
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

func (p *MatcherApi) policyHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	// if title, ok := vars["title"]; ok {
	// 	if title == "user" {
	// 		jsonb, _ := json.MarshalIndent(p.Usercache, "", " ")
	// 		w.Write(jsonb)
	// 	} else if title == "account" {
	// 		jsonb, _ := json.MarshalIndent(p.Accountcache, "", " ")
	// 		w.Write(jsonb)
	// 	} else if title == "asset" {
	// 		jsonb, _ := json.MarshalIndent(p.Assetpolicy, "", " ")
	// 		w.Write(jsonb)
	// 	} else if title == "auth" {
	// 		jsonb, _ := json.MarshalIndent(p.Authcache, "", " ")
	// 		w.Write(jsonb)
	// 	} else if title == "sysdev" {
	// 		jsonb, _ := json.MarshalIndent(p.Sysdevcache, "", " ")
	// 		w.Write(jsonb)
	// 	} else if title == "holiday" {
	// 		w.Write(p.Holidaypolicy.Dump())
	// 		// w.Write(p.Holidaypolicy.DumpHoliday())
	// 		// w.Write(p.Holidaypolicy.DumpWorktime())
	// 	} else {
	// 		fmt.Fprintf(w, "wrong title: %s", title)
	// 	}
	// } else {
	// 	fmt.Fprintf(w, "not set title")
	// }
}

func (p *MatcherApi) policyReloadHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("MatcherApi receive reload policy command")
	w.WriteHeader(http.StatusOK)
	// test uploaded gather policy file, that is etc/user.xml.upload
	// filename := "etc/" + p.Config.Completeconfig.UserCompleteFile + ".upload"
	// if _, err := os.Stat(filename); err == nil {
	// 	c := NewUserCache(filename)
	// 	if err := c.LoadXML(); err != nil {
	// 		log.Errorf("MatcherApi check wrong %s: %s", filename, err)
	// 		// maybe delete the uploaded file
	// 	} else {
	// 		log.Infof("MatcherApi existed upload file %s, reload it", filename)
	// 		if err := os.Rename(filename, filename[:len(filename)-7]); err != nil {
	// 			log.Errorf("MatcherApi rename %s failed: %s", filename, err)
	// 		}
	// 	}
	// }
	// filename = "etc/" + p.Config.Completeconfig.AccountCompleteFile + ".upload"
	// if _, err := os.Stat(filename); err == nil {
	// 	c := NewAccountCache(filename)
	// 	if err := c.LoadXML(); err != nil {
	// 		log.Errorf("MatcherApi check wrong %s: %s", filename, err)
	// 		// maybe delete the uploaded file
	// 	} else {
	// 		log.Infof("MatcherApi existed upload file %s, reload it", filename)
	// 		if err := os.Rename(filename, filename[:len(filename)-7]); err != nil {
	// 			log.Errorf("MatcherApi rename %s failed: %s", filename, err)
	// 		}
	// 	}
	// }
	// filename = "etc/" + p.Config.Completeconfig.AssetCompleteFile + ".upload"
	// if _, err := os.Stat(filename); err == nil {
	// 	c := NewAssetCache(filename)
	// 	if err := c.LoadXML(); err != nil {
	// 		log.Errorf("MatcherApi check wrong %s: %s", filename, err)
	// 		// maybe delete the uploaded file
	// 	} else {
	// 		log.Infof("MatcherApi existed upload file %s, reload it", filename)
	// 		if err := os.Rename(filename, filename[:len(filename)-7]); err != nil {
	// 			log.Errorf("MatcherApi rename %s failed: %s", filename, err)
	// 		}
	// 	}
	// }
	// filename = "etc/" + p.Config.Completeconfig.AuthCompleteFile + ".upload"
	// if _, err := os.Stat(filename); err == nil {
	// 	c := NewAuthCache(filename)
	// 	if err := c.LoadXML(); err != nil {
	// 		log.Errorf("MatcherApi check wrong %s: %s", filename, err)
	// 		// maybe delete the uploaded file
	// 	} else {
	// 		log.Infof("MatcherApi existed upload file %s, reload it", filename)
	// 		if err := os.Rename(filename, filename[:len(filename)-7]); err != nil {
	// 			log.Errorf("MatcherApi rename %s failed: %s", filename, err)
	// 		}
	// 	}
	// }
	// filename = "etc/" + p.Config.Completeconfig.SysdevCompleteFile + ".upload"
	// if _, err := os.Stat(filename); err == nil {
	// 	c := NewSysdevCache(filename)
	// 	if err := c.LoadXML(); err != nil {
	// 		log.Errorf("MatcherApi check wrong %s: %s", filename, err)
	// 		// maybe delete the uploaded file
	// 	} else {
	// 		log.Infof("MatcherApi existed upload file %s, reload it", filename)
	// 		if err := os.Rename(filename, filename[:len(filename)-7]); err != nil {
	// 			log.Errorf("MatcherApi rename %s failed: %s", filename, err)
	// 		}
	// 	}
	// }
	// filename = "etc/" + p.Config.Completeconfig.HolidayPolicyFile + ".upload"
	// if _, err := os.Stat(filename); err == nil {
	// 	c := NewHolidayPolicy(filename)
	// 	if err := c.LoadJson(); err != nil {
	// 		log.Errorf("MatcherApi check wrong %s: %s", filename, err)
	// 		// maybe delete the uploaded file
	// 	} else {
	// 		log.Infof("MatcherApi existed upload file %s, reload it", filename)
	// 		if err := os.Rename(filename, filename[:len(filename)-7]); err != nil {
	// 			log.Errorf("MatcherApi rename %s failed: %s", filename, err)
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
