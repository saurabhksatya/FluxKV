package loadtest

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	serverAddr := startRealServer(t, ctx)
	time.Sleep(200 * time.Millisecond)

	const targetRequests = 100000
	const testDuration = 30 * time.Second

	var (
		totalRequests   int64
		successRequests int64
		failedRequests  int64
		totalLatency    int64
		latencies       []int64
		mu              sync.Mutex
	)

	var wg sync.WaitGroup
	requestsPerWorker := targetRequests / 10
	remainder := targetRequests % 10

	startTime := time.Now()
	deadline := startTime.Add(testDuration)

	for i := 0; i < 10; i++ {
		count := requestsPerWorker
		if i < remainder {
			count++
		}
		wg.Add(1)
		go func(reqCount int) {
			defer wg.Done()
			runWorker(ctx, serverAddr, reqCount, deadline, &totalRequests, &successRequests, &failedRequests, &totalLatency, &latencies, &mu)
		}(count)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		t.Log("Test timeout reached")
	}

	elapsed := time.Since(startTime)
	actualDuration := elapsed
	if elapsed > testDuration {
		actualDuration = testDuration
	}

	throughput := float64(atomic.LoadInt64(&totalRequests)) / actualDuration.Seconds()
	successRate := float64(atomic.LoadInt64(&successRequests)) / float64(atomic.LoadInt64(&totalRequests)) * 100

	t.Logf("=== Load Test Results ===")
	t.Logf("Target requests: %d", targetRequests)
	t.Logf("Test duration: %v", testDuration)
	t.Logf("Actual duration: %v", actualDuration)
	t.Logf("Total requests: %d", atomic.LoadInt64(&totalRequests))
	t.Logf("Successful: %d", atomic.LoadInt64(&successRequests))
	t.Logf("Failed: %d", atomic.LoadInt64(&failedRequests))
	t.Logf("Throughput: %.2f req/s", throughput)
	t.Logf("Success rate: %.2f%%", successRate)

	if len(latencies) > 0 {
		mu.Lock()
		printLatencyStats(latencies)
		mu.Unlock()
	}

	if atomic.LoadInt64(&failedRequests) > 0 {
		t.Logf("Warning: %d requests failed", atomic.LoadInt64(&failedRequests))
	}
}

func runWorker(ctx context.Context, serverAddr string, count int, deadline time.Time,
	totalReqs, successReqs, failedReqs, totalLat *int64, latencies *[]int64, mu *sync.Mutex) {

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		atomic.AddInt64(failedReqs, int64(count))
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	interval := time.Duration(float64(30*time.Second) / float64(count))
	if interval < 100*time.Microsecond {
		interval = 100 * time.Microsecond
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if time.Now().After(deadline) {
				return
			}

			key := fmt.Sprintf("key%d", i)
			value := fmt.Sprintf("value%d", i)

			start := time.Now()

			req := buildRESPSet(key, value)
			_, err := writer.WriteString(req)
			if err != nil {
				atomic.AddInt64(failedReqs, 1)
				continue
			}
			err = writer.Flush()
			if err != nil {
				atomic.AddInt64(failedReqs, 1)
				continue
			}

			resp, err := reader.ReadString('\n')
			if err != nil {
				atomic.AddInt64(failedReqs, 1)
				continue
			}

			latency := time.Since(start).Microseconds()
			atomic.AddInt64(totalLat, latency)

			mu.Lock()
			*latencies = append(*latencies, latency)
			mu.Unlock()

			if strings.HasPrefix(resp, "+OK") {
				atomic.AddInt64(successReqs, 1)
			} else {
				atomic.AddInt64(failedReqs, 1)
			}
			atomic.AddInt64(totalReqs, 1)
		}
	}
}

func startRealServer(t *testing.T, ctx context.Context) string {
	addrStr := "127.0.0.1:8000"

	return addrStr
}

func buildRESPSet(key, value string) string {
	return fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(key), key, len(value), value)
}

func printLatencyStats(latencies []int64) {
	if len(latencies) == 0 {
		return
	}

	sorted := make([]int64, len(latencies))
	copy(sorted, latencies)

	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	p50 := sorted[len(sorted)*50/100]
	p90 := sorted[len(sorted)*90/100]
	p99 := sorted[len(sorted)*99/100]
	p999 := sorted[len(sorted)*999/1000]

	var sum int64
	for _, l := range sorted {
		sum += l
	}
	mean := float64(sum) / float64(len(sorted))

	var variance float64
	for _, l := range sorted {
		diff := float64(l) - mean
		variance += diff * diff
	}
	stddev := math.Sqrt(variance / float64(len(sorted)))

	fmt.Printf("Latency (µs): p50=%d p90=%d p99=%d p99.9=%d mean=%.2f stddev=%.2f max=%d\n",
		p50, p90, p99, p999, mean, stddev, sorted[len(sorted)-1])
}
