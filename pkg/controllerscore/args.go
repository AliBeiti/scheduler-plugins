package controllerscore

type ControllerScoreArgs struct {
	Port int `json:"port,omitempty"`

	Path string `json:"path,omitempty"`

	TimeoutMillis int `json:"timeoutMillis,omitempty"`
}
