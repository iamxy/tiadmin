package schema

import ()

type StorageMetrics struct {
	Usage    int64 `json:"usage"`
	Capacity int64 `json:"capacity"`
}
