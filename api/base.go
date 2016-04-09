package api

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/pingcap/tiadmin/schema"
)

type baseController struct {
	beego.Controller
}

func (c *baseController) ServeError(code int32, err string) {
	beego.Error(fmt.Sprintf("code: %d, error: %v", code, err))
	modelError := &schema.ModelError{
		ErrCode: code,
		Reason:  err,
	}
	if json, err := json.Marshal(modelError); err != nil {
		beego.Error("Failed to marshal object, %v, %v", modelError, err)
		c.Abort("500")
	} else {
		c.CustomAbort(500, string(json))
	}
}

func transformMapToEnvironments(envMap map[string]string) []schema.Environment {
	res := []schema.Environment{}
	for k, v := range envMap {
		s := schema.Environment{
			Name:  k,
			Value: v,
		}
		res = append(res, s)
	}
	return res
}

func transformEnvironmentsToMap(envs []schema.Environment) map[string]string {
	res := make(map[string]string)
	for _, v := range envs {
		res[v.Name] = v.Value
	}
	return res
}
