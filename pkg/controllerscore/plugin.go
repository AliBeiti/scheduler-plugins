package controllerscore

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"k8s.io/klog/v2"

	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const PluginName = "ControllerScore"

const maxNodeScore int64 = 100

const (
	defaultPort      = 6000
	defaultPath      = "/getScore"
	defaultTimeoutMs = 200
)

type ControllerScore struct {
	handle     framework.Handle
	port       int
	path       string
	timeout    time.Duration
	httpClient *http.Client
}

var _ framework.FilterPlugin = &ControllerScore{}
var _ framework.ScorePlugin = &ControllerScore{}

func (cs *ControllerScore) Name() string {
	return PluginName
}

func (cs *ControllerScore) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

func getInternalIP(node *v1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == v1.NodeInternalIP {
			return addr.Address
		}
	}
	return ""
}

func (cs *ControllerScore) fetchControllerScore(nodeIP string) (float64, error) {
	url := fmt.Sprintf("http://%s:%d%s", nodeIP, cs.port, cs.path)
	klog.V(3).Infof("ControllerScore: GET %s", url)
	ctx, cancel := context.WithTimeout(context.Background(), cs.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		klog.Errorf("ControllerScore: building request for %s: %v", url, err)
		return 0, err
	}

	resp, err := cs.httpClient.Do(req)
	if err != nil {
		klog.Errorf("ControllerScore: HTTP get %s failed: %v", url, err)
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("controller returned HTTP %d", resp.StatusCode)
	}

	var payload struct {
		Score float64 `json:"score"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		klog.Errorf("ControllerScore: decoding response from %s: %v", url, err)
		return 0, err
	}
	klog.V(3).Infof("ControllerScore: node %s returned raw score %f", nodeIP, payload.Score)
	return payload.Score, nil

}

func (cs *ControllerScore) Filter(
	ctx context.Context,
	state *framework.CycleState,
	pod *v1.Pod,
	nodeInfo *framework.NodeInfo,
) *framework.Status {
	node := nodeInfo.Node()
	klog.V(4).Infof("ControllerScore: filter called for pod %s on node %s", pod.Name, node.Name)
	ip := getInternalIP(node)
	if ip == "" {
		return framework.NewStatus(framework.Unschedulable, " no InternalIP")
	}

	rawScore, err := cs.fetchControllerScore(ip)
	if err != nil {
		rawScore = 0.0
	}
	if rawScore >= 1 {
		return framework.NewStatus(framework.Unschedulable, "node score == 1.0")
	}
	return framework.NewStatus(framework.Success)
}

func (cs *ControllerScore) Score(
	ctx context.Context,
	state *framework.CycleState,
	pod *v1.Pod,
	nodeName string,
) (int64, *framework.Status) {
	nodeInfo, err := cs.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Success)
	}
	node := nodeInfo.Node()
	klog.V(4).Infof("ControllerScore: score called for pod %s on node %s", pod.Name, node.Name)
	ip := getInternalIP(node)
	if ip == "" {
		return 0, framework.NewStatus(framework.Error, "no InternalIP")
	}

	rawScore, err := cs.fetchControllerScore(ip)
	if err != nil {
		rawScore = 0.0
	}
	if rawScore < 0.0 {
		rawScore = 0.0
	}
	if rawScore > 1.0 {
		rawScore = 1.0
	}

	schedulerSocre := int64((1.0 - rawScore) * float64(maxNodeScore))
	return schedulerSocre, framework.NewStatus(framework.Success)
}

func newControllerScore(obj interface{}, handle framework.Handle) (framework.Plugin, error) {
	args := obj.(*ControllerScoreArgs)

	port := args.Port
	if port == 0 {
		port = defaultPort
	}
	path := args.Path
	if path == "" {
		path = defaultPath
	}
	timeoutMs := args.TimeoutMillis
	if timeoutMs <= 0 {
		timeoutMs = defaultTimeoutMs
	}

	client := &http.Client{
		Timeout: time.Duration(timeoutMs) * time.Millisecond,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: time.Duration(timeoutMs) * time.Millisecond,
			}).DialContext,
		},
	}

	return &ControllerScore{
		handle:     handle,
		port:       port,
		path:       path,
		timeout:    time.Duration(timeoutMs) * time.Millisecond,
		httpClient: client,
	}, nil

}
