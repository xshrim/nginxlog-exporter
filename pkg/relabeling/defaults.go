package relabeling

import "xshrim/nginxlog-exporter/pkg/config"

// DefaultRelabelings are hardcoded relabeling configs that are always there
// and do not need to be explicitly configured
var DefaultRelabelings = []*Relabeling{
	{
		config.RelabelConfig{
			TargetLabel: "method",
			SourceValue: "request",
			Split:       1,
			Whitelist:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"},
		},
	},
	{
		config.RelabelConfig{
			TargetLabel: "status",
			SourceValue: "status",
		},
	},
}
