package schema

import ()

type ModelError struct {
	ErrCode int32  `json:"errCode"`
	Reason  string `json:"reason"`
}
