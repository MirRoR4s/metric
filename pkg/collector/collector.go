// Package collector provides various collectors that can be used to collect metrics about the application and its environment. These collectors can be registered with a Prometheus registry and will automatically update their metrics in the background.
package collector

import (
	"context"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/MirRoR4s/metric/pkg/metric"
	"github.com/shirou/gopsutil/v4/process"
)

// NewHttpRequestsTotal creates a new collector that tracks the total number of HTTP requests received by the application. It returns a Counter metric that can be registered with a Prometheus registry, and a middleware function that can be used to wrap HTTP handlers to automatically increment the counter for each incoming request.
//
// For example, you can use it like this:
//
//	counter, middleware, err := metric.NewHttpRequestsTotal()
//	if err != nil {
//		log.Fatalf("Failed to create HTTP requests total counter: %v", err)
//	}
//	registry.Register(counter)
//	http.Handle("/metrics", registry.Handler())
//	http.Handle("/hello", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte("Hello, World!"))
//	})))
func NewHttpRequestsTotal() (*metric.Counter, func(http.Handler) http.Handler, error) {
	counter, err := metric.NewCounter("http_requests_total", "Total number of HTTP requests.")
	if err != nil {
		return nil, nil, err
	}
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter.Inc()
			next.ServeHTTP(w, r)
		})
	}
	return counter, middleware, nil
}

// processCollector is a struct that collects various metrics about the current process, such as CPU time, memory usage, and thread count. It periodically updates these metrics in the background and provides a method to write them in Prometheus format.
type processCollector struct {
	pid                    int
	maxVirtualMemory       float64
	cpuSecondsTotal        *metric.Gauge
	openFileDescriptors    *metric.Gauge
	maxOpenFileDescriptors *metric.Gauge
	virtualMemory          *metric.Gauge
	virtualMemoryMax       *metric.Gauge
	residentMemory         *metric.Gauge
	heapMemory             *metric.Gauge
	startTime              *metric.Gauge
	threadCount            *metric.Gauge
	ticker                 *time.Ticker
	cancel                 context.CancelFunc
	samplingPeriod         time.Duration
}

// update periodically updates all memory metrics.
func (p *processCollector) update(ctx context.Context) {
	p.ticker = time.NewTicker(p.samplingPeriod)
	defer p.ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping memory metrics updater")
			return
		case <-p.ticker.C:
			proc, err := process.NewProcess(int32(p.pid))
			if err != nil {
				log.Printf("Error creating process: %v", err)
				continue
			}

			cpuTimes, err := proc.Times()
			if err != nil {
				log.Printf("Error getting CPU times: %v", err)
				continue
			}
			p.cpuSecondsTotal.Set(cpuTimes.User + cpuTimes.System)

			numFDs, err := proc.NumFDs()
			if err != nil {
				log.Printf("Error getting number of file descriptors: %v", err)
				continue
			}
			p.openFileDescriptors.Set(float64(numFDs))
			p.maxOpenFileDescriptors.Set(max(p.maxOpenFileDescriptors.Value(), float64(numFDs)))

			memInfo, err := proc.MemoryInfo()
			if err != nil {
				log.Printf("Error getting memory info: %v", err)
				continue
			}
			p.virtualMemory.Set(float64(memInfo.VMS))
			p.residentMemory.Set(float64(memInfo.RSS))

			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			p.heapMemory.Set(float64(m.HeapAlloc))

			startTime, err := proc.CreateTime()
			if err != nil {
				log.Printf("Error getting process start time: %v", err)
				continue
			}
			p.startTime.Set(float64(startTime) / 1000) // Convert milliseconds to seconds

			threadCount, err := proc.NumThreads()
			if err != nil {
				log.Printf("Error getting thread count: %v", err)
				continue
			}
			p.threadCount.Set(float64(threadCount))
		}
	}
}

func (p *processCollector) WritePrometheus() string {
	return p.cpuSecondsTotal.WritePrometheus() + "\n" +
		p.openFileDescriptors.WritePrometheus() + "\n" +
		p.maxOpenFileDescriptors.WritePrometheus() + "\n" +
		p.virtualMemory.WritePrometheus() + "\n" +
		p.virtualMemoryMax.WritePrometheus() + "\n" +
		p.residentMemory.WritePrometheus() + "\n" +
		p.heapMemory.WritePrometheus() + "\n" +
		p.startTime.WritePrometheus() + "\n" +
		p.threadCount.WritePrometheus()
}

type ProcessCollectorOption func(*processCollector)

// WithPID sets the PID of the process to collect metrics for. If not set, it defaults to the current process's PID.
func WithPID(pid int) ProcessCollectorOption {
	return func(p *processCollector) {
		p.pid = pid
	}
}

// WithMaxVirtualMemory sets the maximum virtual memory size in bytes that the process is allowed to use. This value is set on the virtualMemoryMax gauge and can be used to track how much of the allowed virtual memory is being used.
func WithMaxVirtualMemory(maxVirtualMemory float64) ProcessCollectorOption {
	return func(p *processCollector) {
		p.maxVirtualMemory = maxVirtualMemory
	}
}

// WithSamplingPeriod sets the sampling period for updating the process metrics. The default sampling period is 500 milliseconds, but you can adjust it based on your needs. A shorter sampling period will provide more up-to-date metrics but may increase CPU usage, while a longer sampling period will reduce CPU usage but may result in less timely metrics.
func WithSamplingPeriod(samplingPeriod time.Duration) ProcessCollectorOption {
	return func(p *processCollector) {
		p.samplingPeriod = samplingPeriod
	}
}

// NewProcess creates a new process metrics collector that periodically updates various metrics related to the current process, such as CPU time, memory usage, and thread count.
//
// See https://prometheus.io/docs/instrumenting/writing_clientlibs/#process-metrics for more details on the specific metrics collected.
func NewProcess(ctx context.Context, opts ...ProcessCollectorOption) (*processCollector, error) {
	p := &processCollector{
		pid:              os.Getpid(),
		maxVirtualMemory: float64(1 << 30),
		samplingPeriod:   500 * time.Millisecond, // Default sampling period
	}
	for _, opt := range opts {
		opt(p)
	}
	cpuSecondsTotal, err := metric.NewGauge("process_cpu_seconds_total", "Total user and system CPU time spent in seconds.")
	if err != nil {
		return nil, err
	}
	p.cpuSecondsTotal = cpuSecondsTotal

	openFileDescriptors, err := metric.NewGauge("process_open_fds", "Number of open file descriptors.")
	if err != nil {
		return nil, err
	}
	p.openFileDescriptors = openFileDescriptors

	maxOpenFileDescriptors, err := metric.NewGauge("process_max_fds", "Maximum number of open file descriptors.")
	if err != nil {
		return nil, err
	}
	p.maxOpenFileDescriptors = maxOpenFileDescriptors

	virtualMemoryBytes, err := metric.NewGauge("process_virtual_memory_bytes", "Virtual memory size in bytes.")
	if err != nil {
		return nil, err
	}
	p.virtualMemory = virtualMemoryBytes

	virtualMemoryMax, err := metric.NewGauge("process_virtual_memory_max_bytes", "Maximum virtual memory size in bytes.")
	if err != nil {
		return nil, err
	}
	p.virtualMemoryMax = virtualMemoryMax
	p.virtualMemoryMax.Set(p.maxVirtualMemory)

	residentMemoryBytes, err := metric.NewGauge("process_resident_memory_bytes", "Resident memory size in bytes.")
	if err != nil {
		return nil, err
	}
	p.residentMemory = residentMemoryBytes

	heapMemoryBytes, err := metric.NewGauge("process_heap_memory_bytes", "Heap memory size in bytes.")
	if err != nil {
		return nil, err
	}
	p.heapMemory = heapMemoryBytes

	startTime, err := metric.NewGauge("process_start_time_seconds", "Start time of the process since unix epoch in seconds.")
	if err != nil {
		return nil, err
	}
	p.startTime = startTime

	threadCount, err := metric.NewGauge("process_thread_count", "Number of threads currently used by the process.")
	if err != nil {
		return nil, err
	}
	p.threadCount = threadCount

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	go p.update(ctx)
	return p, nil
}
