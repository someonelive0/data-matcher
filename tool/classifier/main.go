package main

import (
	"bufio"
	"data-matcher/model"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	arg_file = flag.String("file", "", "filename of messages with json in one line")
)

func init() {
	flag.Parse()
}

func main() {
	hostpathmap, err := getUrls(*arg_file)
	if err != nil {
		fmt.Printf("getUrls from [%s] failed: %s\n", *arg_file, err)
		os.Exit(1)
	}

	fmt.Printf("hostpathmap number %d\n", len(hostpathmap))
	for host, paths := range hostpathmap {
		fmt.Println(host, len(paths))
		for i := range paths {
			fmt.Println("  ", paths[i])
		}
	}
}

func getUrls(filename string) (map[string][]string, error) {
	fp, err := os.Open(filename)
	if err != nil {
		log.Printf("open data file [%s] failed: %s\n", filename, err)
		return nil, err
	}
	defer fp.Close()

	count := 0
	m := make(map[string]interface{})
	submap := make(map[string]int)
	// paths := make([]string, 0)
	hostpathmap := make(map[string][]string)

	// 读入消息到内存队列
	fp.Seek(0, 0)
	scanner := bufio.NewScanner(fp)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		err := json.Unmarshal(line, &m)
		if err != nil {
			fmt.Printf("Umarshal failed: %s\n%s\n", err, line)
			continue
		}
		count++

		if event_type, ok := m["event_type"]; ok {

			if event_type.(string) == "http" {
				flowHttp := &model.FlowHttp{}
				if err = json.Unmarshal(line, flowHttp); err != nil {
					fmt.Printf("unmarshal http msg failed %s\n", err)
					continue
				} else {
					// paths = append(paths, flowHttp.Http.Hostname+flowHttp.Http.Url)

					hostname := strings.TrimSpace(flowHttp.Http.Hostname)
					path := strings.TrimSpace(flowHttp.Http.Url)
					if len(hostname) == 0 {
						hostname = "NaN"
					}
					if len(path) == 0 {
						path = "/"
					}

					if hostpaths, found := hostpathmap[hostname]; found {
						hostpathmap[hostname] = append(hostpaths, path)
					} else {
						hostpaths = make([]string, 0)
						hostpathmap[hostname] = append(hostpaths, path)
					}
				}
			}

			submap[event_type.(string)]++
		}
	}

	log.Printf("subjects: %#v", submap)
	if scanner.Err() != nil {
		log.Println("file scanner failed: ", scanner.Err())
		return nil, scanner.Err()
	}

	return hostpathmap, nil
}
