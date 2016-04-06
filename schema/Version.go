package schema

import ()

type Version struct {
	Version      string `json:"version"`
	BuildUTCTime string `json:"buildUTCTime"`
}
