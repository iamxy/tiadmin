package schema

import ()

type Environment struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
