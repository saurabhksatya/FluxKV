package commands

import (
	"fluxKV/internal"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type InfoCommand struct{}

var (
	startTime = time.Now()
	proc, _   = process.NewProcess(int32(os.Getpid()))
)

func (i *InfoCommand) execute(conn net.Conn, cmd []string, db *internal.DataStore) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	var (
		cpuPercent float64
		rss        uint64
	)

	if proc != nil {
		if cpu, err := proc.CPUPercent(); err == nil {
			cpuPercent = cpu
		}

		if memInfo, err := proc.MemoryInfo(); err == nil {
			rss = memInfo.RSS
		}
	}

	info := fmt.Sprintf(
		"# Server\r\n"+
			"redis_version:0.1.0\r\n"+
			"uptime_in_seconds:%d\r\n"+
			"\r\n"+
			"# Memory\r\n"+
			"used_memory:%d\r\n"+
			"used_memory_human:%.2fM\r\n"+
			"process_memory:%d\r\n"+
			"process_memory_human:%.2fM\r\n"+
			"\r\n"+
			"# CPU\r\n"+
			"process_cpu_percent:%.2f\r\n"+
			"runtime_num_cpu:%d\r\n"+
			"runtime_goroutines:%d\r\n"+
			"\r\n"+
			"# Replication\r\n"+
			"role:master\r\n",
		int64(time.Since(startTime).Seconds()),
		mem.Alloc,
		float64(mem.Alloc)/(1024*1024),
		rss,
		float64(rss)/(1024*1024),
		cpuPercent,
		runtime.NumCPU(),
		runtime.NumGoroutine(),
	)

	fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(info), info)
}
