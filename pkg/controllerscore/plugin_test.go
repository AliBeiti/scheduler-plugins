package controllerscore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	schedruntime "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
)

// Test that Filter rejects score == 1.0
func TestFilterRejectsScoreOne(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]float64{"score": 1.0})
	}))
	defer ts.Close()

	plugin := &ControllerScore{
		handle:     nil,
		port:       0, // ignored by override
		path:       "",
		timeout:    50 * time.Millisecond,
		httpClient: ts.Client(),
	}
	// Override fetchControllerScore to hit ts.URL
	plugin.fetchControllerScore = func(nodeIP string) (float64, error) {
		resp, err := plugin.httpClient.Get(ts.URL)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		var p struct{ Score float64 }
		json.NewDecoder(resp.Body).Decode(&p)
		return p.Score, nil
	}

	node := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "node1"},
		Status: v1.NodeStatus{
			Addresses: []v1.NodeAddress{
				{Type: v1.NodeInternalIP, Address: "127.0.0.1"},
			},
		},
	}
	nodeInfo := &framework.NodeInfo{Node: node}
	status := plugin.Filter(context.Background(), nil, &v1.Pod{}, nodeInfo)
	if !status.IsUnschedulable() {
		t.Fatalf("expected Unschedulable, got %v", status)
	}
}

// Test that Score(0.25) → 75
func TestScoreConvertsCorrectly(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]float64{"score": 0.25})
	}))
	defer ts.Close()

	plugin := &ControllerScore{
		port:       0,
		path:       "",
		timeout:    50 * time.Millisecond,
		httpClient: ts.Client(),
	}
	plugin.fetchControllerScore = func(nodeIP string) (float64, error) {
		resp, err := plugin.httpClient.Get(ts.URL)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		var p struct{ Score float64 }
		json.NewDecoder(resp.Body).Decode(&p)
		return p.Score, nil
	}

	// Fake a scheduler.Framework with one node “node1”
	node := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "node1"},
		Status: v1.NodeStatus{
			Addresses: []v1.NodeAddress{{Type: v1.NodeInternalIP, Address: "127.0.0.1"}},
		},
	}
	registry := map[string]framework.PluginFactory{}
	cfg := v1alpha1.Plugins{}
	pluginConfig := []v1alpha1.PluginConfig{}
	f, _ := schedruntime.NewFramework(registry, cfg, pluginConfig, schedruntime.WithClientSet(nil),
		schedruntime.WithSnapshotSharedLister(fakeLister{node}))
	plugin.handle = f

	score, status := plugin.Score(context.Background(), nil, &v1.Pod{}, "node1")
	if !status.IsSuccess() {
		t.Fatalf("expected Success, got %v", status)
	}
	if score != 75 {
		t.Errorf("expected 75, got %d", score)
	}
}

// fakeLister returns our one node when asked.
type fakeLister struct{ node *v1.Node }

func (f fakeLister) NodeInfos() framework.NodeInfoLister { return fakeNodeInfoLister{f.node} }
func (f fakeLister) PodInfos() framework.PodInfoLister   { panic("not used") }

type fakeNodeInfoLister struct{ node *v1.Node }

func (f fakeNodeInfoLister) List() ([]*framework.NodeInfo, error) {
	return []*framework.NodeInfo{{Node: f.node}}, nil
}
func (f fakeNodeInfoLister) Get(name string) (*framework.NodeInfo, error) {
	if name == f.node.Name {
		return &framework.NodeInfo{Node: f.node}, nil
	}
	return nil, fmt.Errorf("not found")
}
