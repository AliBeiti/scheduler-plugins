package controllerscore

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

func NewControllerScorePlugin(
	ctx context.Context,
	obj runtime.Object,
	handle framework.Handle,
) (framework.Plugin, error) {
	args := &ControllerScoreArgs{}

	if obj != nil {
		rawJSON, err := json.Marshal(obj)
		if err != nil {
			return nil, fmt.Errorf("controllerscore: failed to marshal config obj: %w", err)
		}
		if err := json.Unmarshal(rawJSON, args); err != nil {
			return nil, fmt.Errorf("controllerscore: failed to unmarshal controllerscoreargs: %w", err)
		}
	}
	return newControllerScore(args, handle)
}
