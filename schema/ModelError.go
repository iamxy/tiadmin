package schema

import ()

type ModelError struct {
	ErrCode int32  `json:"errCode,omitempty"`
	Reason  string `json:"reason,omitempty"`
}
