package controllerscore

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	v1 "k8s.io/pkg/scheduler/framework/v1alhpa1"
)

func NewControllerScorePlugin(
	obj runtime.Unknown,
	handle framework.Handle,
	_ v1.PluginConfig,
) (framework.Plugin, error) {
	args := &ControllerScoreArgs{}
	if len(obj.Raw) > 0 {
		if err := json.Unmarshal(obj.Raw, args); err != nil {
			return nil, fmt.Errorf("ControllerScore: failed to decode args: %w", err)
		}
	}
	return newControllerScore(args, handle)
}
