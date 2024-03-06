package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func HostDump() string {
	s := fmt.Sprintf(`{"timestamp": "%s", `, time.Now().Format(time.RFC3339))
	s += fmt.Sprintf(`"os": %s, `, HostOS())
	s += fmt.Sprintf(`"mem": %s, `, HostMem())
	s += fmt.Sprintf(`"cpu": %s, `, HostCpu())
	s += fmt.Sprintf(`"disk": %s`, HostDisk())
	s += `}`
	return s
}

func HostLoading() string {
	s := fmt.Sprintf(`{"timestamp": "%s", `, time.Now().Format(time.RFC3339))
	s += fmt.Sprintf(`"measurement": "MiB", `)
	v, _ := mem.VirtualMemory()
	s += fmt.Sprintf(`"mem_total": %v, `, B2MB(v.Total))
	s += fmt.Sprintf(`"mem_available": %v, `, B2MB(v.Available))
	s += fmt.Sprintf(`"mem_used_percent": %.2f, `, v.UsedPercent)

	totalPercent, _ := cpu.Percent(3*time.Second, false)
	perPercents, _ := cpu.Percent(3*time.Second, true)
	// fmt.Printf("total percent:%v per percents:%v", totalPercent, perPercents)
	s += `"cpu_percent": [ `
	for i := range totalPercent {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%.2f", totalPercent[i])
	}

	s += `] ,"cpu_per_percent": [ `
	for i := range perPercents {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%.2f", perPercents[i])
	}

	s += `]}`
	return s
}

func HostOS() string {
	timestamp, _ := host.BootTime()
	t := time.Unix(int64(timestamp), 0)
	s := fmt.Sprintf(`{"boot": "%s", `, t.Local().Format(time.RFC3339))

	kernel_version, _ := host.KernelVersion()
	s += fmt.Sprintf(`"kernel_version": "%s", `, kernel_version)

	platform, family, version, _ := host.PlatformInformation()
	s += fmt.Sprintf(`"platform": "%s", `, platform)
	s += fmt.Sprintf(`"family": "%s", `, family)
	s += fmt.Sprintf(`"version": "%s"}`, version)
	return s
}

func HostMem() string {
	v, _ := mem.VirtualMemory()
	//fmt.Printf("Total: %v, Available: %v, UsedPercent:%f%%\n", v.Total, v.Available, v.UsedPercent)
	return fmt.Sprintf("%v", v)
}

func HostCpu() string {
	physicalCnt, _ := cpu.Counts(false)
	logicalCnt, _ := cpu.Counts(true)
	s := fmt.Sprintf(`{"physical": %d, "logical": %d, `, physicalCnt, logicalCnt)

	// totalPercent, _ := cpu.Percent(3*time.Second, false)
	// perPercents, _ := cpu.Percent(3*time.Second, true)
	// fmt.Printf("total percent:%v per percents:%v", totalPercent, perPercents)

	infos, _ := cpu.Info()
	s += `"info": [`
	for i, info := range infos {
		if i > 0 {
			s += ", "
		}
		b, _ := json.MarshalIndent(info, "", " ")
		s += string(b)
	}
	s += `]}`

	// timescpu, _ := cpu.Times(true)
	// for _, info := range timescpu {
	// 	data, _ := json.MarshalIndent(info, "", " ")
	// 	fmt.Print(string(data))
	// }
	return s
}

func HostDisk() string {
	info, _ := disk.Usage("/")
	b, _ := json.MarshalIndent(info, "", " ")
	return string(b)
}
