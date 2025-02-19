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

package oas3_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/erda-project/erda/pkg/swagger/oas3"
)

func TestToYaml(t *testing.T) {
	data, err := ioutil.ReadFile("./testdata/gaia-oas3.json")
	if err != nil {
		t.Error(err)
	}

	valid := json.Valid(data)
	t.Log(valid)
}

//// go test -v -run TestExpandBigSwagger
//func TestExpandBigSwagger(t *testing.T) {
//
//	data, err := ioutil.ReadFile("./testdata/gaia-oas3.json")
//	if err != nil {
//		t.Fatalf("failed to ReadFile, err: %v", err)
//	}
//
//	v3, err := swagger.LoadFromData(data)
//	if err != nil {
//		t.Fatalf("failed to LoadFromData, err: %v", err)
//	}
//
//	var paths []string
//	for path_, _ := range v3.Paths {
//		paths = append(paths, path_)
//	}
//	sort.Strings(paths)
//	for _, path_ := range paths {
//		pathItem := v3.Paths.Find(path_)
//		for _, operation := range map[string]*openapi3.GenerateOperation{
//			http.MethodDelete:  pathItem.Delete,
//			http.MethodGet:     pathItem.Get,
//			http.MethodHead:    pathItem.Head,
//			http.MethodOptions: pathItem.Options,
//			http.MethodPatch:   pathItem.Patch,
//			http.MethodPost:    pathItem.Post,
//			http.MethodPut:     pathItem.Put,
//			http.MethodTrace:   pathItem.Trace,
//			http.MethodConnect: pathItem.Connect,
//		} {
//			if operation == nil {
//				continue
//			}
//			if err := oas3.ExpandOperation(operation, v3); err != nil {
//				t.Fatalf("failed to ExpandOperation: %v", err)
//			}
//			operation.Extensions = nil
//			operation.Parameters = nil
//			// t.Logf("%s %s bodyData", method, path_)
//			// bodyData, _ := json.Marshal(operation.RequestBody)
//			// t.Log(string(bodyData))
//			// t.Logf("%s %s responseData", method, path_)
//			// responsesData, _ := json.Marshal(operation.Responses)
//			// t.Log(string(responsesData))
//		}
//	}
//
//	v3.Components.Schemas = nil
//	indent, err := yaml.Marshal(v3)
//	if err != nil {
//		t.Error(err)
//	}
//	t.Log(string(indent))
//}

//func TestOAS2To3(t *testing.T) {
//	data, err := ioutil.ReadFile("./testdata/swagger_all.json")
//	if err != nil {
//		t.Error(err)
//	}
//	v3, err := swagger.LoadFromData(data)
//	if err != nil {
//		t.Error(err)
//	}
//
//	data, err = json.Marshal(v3)
//	if err != nil {
//		t.Error(err)
//	}
//
//	data, err = oasconv.JSONToYAML(data)
//	if err != nil {
//		t.Error(err)
//	}
//
//	t.Log(string(data))
//}

// go test -v -run TestExpandAllOf
func TestExpandAllOf(t *testing.T) {
	data, err := ioutil.ReadFile("./testdata/allof-oas3.json")
	if err != nil {
		t.Error(err)
	}
	v3, err := oas3.LoadFromData(data)
	if err != nil {
		t.Error(err)
	}

	schema := v3.Paths["/new-resource"].Get.Responses["200"].Value.Content["application/json"].Schema
	if err = oas3.ExpandSchemaRef(schema, v3); err != nil {
		t.Error(err)
	}

	indent, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(indent))
}
