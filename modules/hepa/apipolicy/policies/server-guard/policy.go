// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package serverguard

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/modules/hepa/apipolicy"
	"github.com/erda-project/erda/modules/hepa/k8s"
)

type Policy struct {
	apipolicy.BasePolicy
}

func (policy Policy) NeedSerialUpdate() bool {
	return true
}

func (policy Policy) CreateDefaultConfig(ctx map[string]interface{}) apipolicy.PolicyDto {
	dto := &PolicyDto{
		ExtraLatency:   50,
		RefuseCode:     429,
		RefuseResponse: "System is busy, please try it later.",
	}
	dto.Switch = false
	return dto
}

func (policy Policy) MergeDiceConfig(conf map[string]interface{}) (apipolicy.PolicyDto, error) {
	dto := &PolicyDto{
		ExtraLatency:   50,
		RefuseCode:     429,
		RefuseResponse: "System is busy, please try it later.",
	}
	if len(conf) == 0 {
		return dto, nil
	}
	dto.Switch = true
	if value, ok := conf["qps"]; ok {
		if vv, ok := value.(float64); ok {
			if vv != 0 {
				dto.MaxTps = int64(vv)
			} else {
				dto.Switch = false
			}
		}
	}
	if value, ok := conf["max_delay"]; ok {
		if vv, ok := value.(float64); ok {
			if vv != 0 {
				dto.ExtraLatency = int64(vv)
			}
		}
	}
	if value, ok := conf["deny_status"]; ok {
		if vv, ok := value.(float64); ok {
			if vv != 0 {
				dto.RefuseCode = int64(vv)
			}
		}
	}
	if value, ok := conf["deny_content"]; ok {
		if vv, ok := value.(string); ok {
			if vv != "" {
				dto.RefuseResponse = vv
			}
		}
	}
	dto.AdjustDto()
	return dto, nil
}

func (policy Policy) UnmarshalConfig(config []byte) (apipolicy.PolicyDto, error, string) {
	policyDto := &PolicyDto{}
	err := json.Unmarshal(config, policyDto)
	if err != nil {
		return nil, errors.Wrapf(err, "json parse config failed, config:%s", config), "Invalid config"
	}
	ok, msg := policyDto.IsValidDto()
	if !ok {
		return nil, errors.Errorf("invalid policy dto, msg:%s", msg), msg
	}
	return policyDto, nil, ""
}

func (policy Policy) ParseConfig(dto apipolicy.PolicyDto, ctx map[string]interface{}) (apipolicy.PolicyConfig, error) {
	res := apipolicy.PolicyConfig{}
	policyDto, ok := dto.(*PolicyDto)
	if !ok {
		return res, errors.Errorf("invalid config:%+v", dto)
	}
	if !policyDto.Switch {
		emptyStr := ""
		// use empty str trigger regions update
		res.IngressAnnotation = &apipolicy.IngressAnnotation{
			LocationSnippet: &emptyStr,
		}
		res.IngressController = &apipolicy.IngressController{
			HttpSnippet:   &emptyStr,
			ServerSnippet: &emptyStr,
		}
		res.AnnotationReset = true
		return res, nil
	}
	value, ok := ctx[apipolicy.CTX_IDENTIFY]
	if !ok {
		return res, errors.Errorf("get identify failed:%+v", ctx)
	}
	id, ok := value.(string)
	if !ok {
		return res, errors.Errorf("convert failed:%+v", value)
	}
	value, ok = ctx[apipolicy.CTX_K8S_CLIENT]
	if !ok {
		return res, errors.Errorf("get k8s client failed:%+v", ctx)
	}
	adapter, ok := value.(k8s.K8SAdapter)
	if !ok {
		return res, errors.Errorf("convert failed:%+v", value)
	}
	count, err := adapter.CountIngressController()
	if err != nil {
		count = 1
		logrus.Errorf("get ingress controller count failed, err:%+v", err)
	}
	tps := int64(math.Ceil(float64(policyDto.MaxTps) / float64(count)))
	burst := int64(math.Ceil(float64(policyDto.ExtraLatency) * float64(tps) / 1000))
	limitReqZone := fmt.Sprintf("limit_req_zone 1 zone=server-guard-%s:1m rate=%dr/s;\n", id, tps)
	var limitReq string
	if burst != 0 {
		limitReq = fmt.Sprintf("limit_req zone=server-guard-%s burst=%d;\n", id, burst)
	} else {
		limitReq = fmt.Sprintf("limit_req zone=server-guard-%s nodelay;\n", id)
	}
	limitReqStatus := fmt.Sprintf("limit_req_status %d;\n", LIMIT_INNER_STATUS)
	errorPage := fmt.Sprintf("error_page %d = @LIMIT-%s;\n", LIMIT_INNER_STATUS, id)
	var namedLocation string
	if policyDto.RefuseResponse == "" {
		namedLocation = fmt.Sprintf(`
location @LIMIT-%s {
    return %d;
}
`, id, policyDto.RefuseCode)
	} else if policyDto.RefuseCode >= 300 && policyDto.RefuseCode < 400 {
		namedLocation = fmt.Sprintf(`
location @LIMIT-%s {
    return %d "%s";
}
`, id, policyDto.RefuseCode, policyDto.RefuseResponse)
	} else {
		namedLocation = fmt.Sprintf(`
location @LIMIT-%s {
    more_set_headers 'Content-Type: text/plain; charset=utf-8';
    return %d "%s";
}
`, id, policyDto.RefuseCode, policyDto.RefuseResponse)
	}
	locationSnippet := limitReq + limitReqStatus + errorPage
	res.IngressAnnotation = &apipolicy.IngressAnnotation{
		LocationSnippet: &locationSnippet,
	}
	httpSnippet := limitReqZone
	serverSnippet := namedLocation
	res.IngressController = &apipolicy.IngressController{
		HttpSnippet:   &httpSnippet,
		ServerSnippet: &serverSnippet,
	}
	return res, nil
}

func init() {
	err := apipolicy.RegisterPolicyEngine("safety-server-guard", &Policy{})
	if err != nil {
		panic(err)
	}
}
