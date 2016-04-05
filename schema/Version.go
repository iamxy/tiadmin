package schema

import ()

type Version struct {
	Version      string `json:"version,omitempty"`
	BuildUTCTime string `json:"buildUTCTime,omitempty"`
}
