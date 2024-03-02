package utils

import (
	"fmt"
	"runtime"
)

/*
runtime.MemStats这个结构体包含的字段比较多，但是大多都很有用，去掉那些注释来看各个属性，会发现各个属性都是很有价值的：
1、Alloc uint64 //golang语言框架堆空间分配的字节数
2、TotalAlloc uint64 //从服务开始运行至今分配器为分配的堆空间总 和，只有增加，释放的时候不减少
3、Sys uint64 //服务现在系统使用的内存
4、Lookups uint64 //被runtime监视的指针数
5、Mallocs uint64 //服务malloc的次数
6、Frees uint64 //服务回收的heap objects的字节数
7、HeapAlloc uint64 //服务分配的堆内存字节数
8、HeapSys uint64 //系统分配的作为运行栈的内存
9、HeapIdle uint64 //申请但是未分配的堆内存或者回收了的堆内存（空闲）字节数
10、HeapInuse uint64 //正在使用的堆内存字节数
10、HeapReleased uint64 //返回给OS的堆内存，类似C/C++中的free。
11、HeapObjects uint64 //堆内存块申请的量
12、StackInuse uint64 //正在使用的栈字节数
13、StackSys uint64 //系统分配的作为运行栈的内存
14、MSpanInuse uint64 //用于测试用的结构体使用的字节数
15、MSpanSys uint64 //系统为测试用的结构体分配的字节数
16、MCacheInuse uint64 //mcache结构体申请的字节数(不会被视为垃圾回收)
17、MCacheSys uint64 //操作系统申请的堆空间用于mcache的字节数
18、BuckHashSys uint64 //用于剖析桶散列表的堆空间
19、GCSys uint64 //垃圾回收标记元信息使用的内存
20、OtherSys uint64 //golang系统架构占用的额外空间
21、NextGC uint64 //垃圾回收器检视的内存大小
22、LastGC uint64 // 垃圾回收器最后一次执行时间。
23、PauseTotalNs uint64 // 垃圾回收或者其他信息收集导致服务暂停的次数。
24、PauseNs [256]uint64 //一个循环队列，记录最近垃圾回收系统中断的时间
25、PauseEnd [256]uint64 //一个循环队列，记录最近垃圾回收系统中断的时间开始点。
26、NumForcedGC uint32 //服务调用runtime.GC()强制使用垃圾回收的次数。
27、GCCPUFraction float64 //垃圾回收占用服务CPU工作的时间总和。如果有100个goroutine，垃圾回收的时间为1S,那么久占用了100S。
28、BySize //内存分配器使用情况

如果我们要观察应用层代码使用的内存大小，可以观察Alloc字段。
如果我们要观察程序从系统申请的内存以及归还给系统的情况，可以观察HeapIdle和HeapReleased字段。
*/
func ProcessMem() string {

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	gonum := runtime.NumGoroutine()

	s := fmt.Sprintf(`{ "measurement": "MiB", "alloc": %d, "heap_idle": %d, "heap_released": %d, `,
		B2MB(ms.Alloc), B2MB(ms.HeapIdle), B2MB(ms.HeapReleased))
	s += fmt.Sprintf(`"heap_inuse": %d, `, B2MB(ms.HeapInuse))
	s += fmt.Sprintf(`"sys_heap_stack": "%d,%d,%d", `, B2MB(ms.Sys), B2MB(ms.HeapSys), B2MB(ms.StackSys))
	s += fmt.Sprintf(`"num_gc": %d, `, ms.NumGC)

	s += fmt.Sprintf(`"num_goroutine": %d }`, gonum)
	return s
}

// windows has no syscall.Rusage
func ProcessCpu() string {
	// var rusage syscall.Rusage
	// if err := syscall.Getrusage(os.Getpid(), &rusage); err != nil {
	// 	return ""
	// }
	return ""
}

func B2MB(b uint64) uint64 {
	B := b / 1024 / 1024
	if B == 0 {
		B = 1
	}
	return B
}

//===========================================
/*
func runInWindows(cmd string) (string, error) {
	result, err := exec.Command("cmd", "/c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

func RunCommand(cmd string) (string, error) {
	if runtime.GOOS == "windows" {
		return runInWindows(cmd)
	} else {
		return runInLinux(cmd)
	}
}

func runInLinux(cmd string) (string, error) {
	fmt.Println("Running Linux cmd:" + cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

//根据进程名判断进程是否运行
func CheckProRunning(program string) (bool, error) {
	a := `ps ux | awk '/` + program + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	if err != nil {
		return false, err
	}
	return pid != "", nil
}

//根据进程名称获取进程ID
func GetPid(program string) (string, error) {
	a := `ps ux | awk '/` + program + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	return pid, err
}
*/
