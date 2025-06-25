package statsserver

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    types "k8s.io/cri-api/pkg/apis/runtime/v1"

    "github.com/cri-o/cri-o/internal/config/cgmgr"
    "github.com/cri-o/cri-o/internal/lib/sandbox"
    "github.com/cri-o/cri-o/internal/log"
)

// fetchThreadCount fetches the thread count for a container using its PID.
func fetchThreadCount(containerPid int) (int, error) {
    taskDir := filepath.Join("/proc", fmt.Sprintf("%d", containerPid), "task")
    entries, err := os.ReadDir(taskDir)
    if err != nil {
        return 0, err
    }
    return len(entries), nil
}

// generateSandboxProcessThreadMetrics generates process-related metrics for a single container in the sandbox.
func generateSandboxProcessThreadMetrics(sb *sandbox.Sandbox, pids *cgmgr.PidsStats) []*types.Metric {
    
    threadCount, err := fetchThreadCount(int(pids.Current))

    if err != nil {
        log.Warnf(context.Background(), "Failed to fetch thread count for container PID %d: %v", int(pids.Current), err)
        threadCount = 0
    }

    processMetrics := []*containerMetric{
        {
            desc: containerProcessThreadsTotal,
            valueFunc: func() metricValues {
                return metricValues{{
                    value:      uint64(threadCount),
                    metricType: types.MetricType_GAUGE,
                }}
            },
        },
    }

    return computeSandboxMetrics(sb, processMetrics, "process")
}