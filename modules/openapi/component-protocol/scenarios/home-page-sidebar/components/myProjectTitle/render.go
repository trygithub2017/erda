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

package myProjectTitle

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
	protocol "github.com/erda-project/erda/modules/openapi/component-protocol"
	"github.com/erda-project/erda/modules/openapi/component-protocol/scenarios/home-page-sidebar/i18n"
)

type MyProjectTitle struct {
	ctxBdl protocol.ContextBundle
	Type   string `json:"type"`
	Props  Props  `json:"props"`
}

type Props struct {
	Visible        bool          `json:"visible"`
	Title          string        `json:"title"`
	Level          int           `json:"level"`
	NoMarginBottom bool          `json:"noMarginBottom"`
	ShowDivider    bool          `json:"showDivider"`
	Size           string        `json:"size"`
	Operations     []interface{} `json:"operations"`
}

func (this *MyProjectTitle) SetCtxBundle(ctx context.Context) error {
	bdl := ctx.Value(protocol.GlobalInnerKeyCtxBundle.String()).(protocol.ContextBundle)
	if bdl.Bdl == nil || bdl.I18nPrinter == nil {
		return fmt.Errorf("invalid context bundle")
	}
	logrus.Infof("inParams:%+v, identity:%+v", bdl.InParams, bdl.Identity)
	this.ctxBdl = bdl
	return nil
}

func (t *MyProjectTitle) Render(ctx context.Context, c *apistructs.Component, scenario apistructs.ComponentProtocolScenario, event apistructs.ComponentEvent, gs *apistructs.GlobalStateData) error {
	if err := t.SetCtxBundle(ctx); err != nil {
		return err
	}
	i18nLocale := t.ctxBdl.Bdl.GetLocale(t.ctxBdl.Locale)
	t.Type = "Title"
	t.Props.Title = i18nLocale.Get(i18n.I18nKeyProjectName)
	t.Props.NoMarginBottom = true
	t.Props.Level = 1
	t.Props.Visible = true
	t.Props.ShowDivider = true
	t.Props.Size = "normal"
	if t.ctxBdl.Identity.OrgID == "" {
		return nil
	}
	orgIntId, err := strconv.Atoi(t.ctxBdl.Identity.OrgID)
	if err != nil {
		return err
	}
	req := &apistructs.PermissionCheckRequest{
		UserID:   t.ctxBdl.Identity.UserID,
		Scope:    apistructs.OrgScope,
		ScopeID:  uint64(orgIntId),
		Resource: apistructs.ProjectResource,
		Action:   apistructs.CreateAction,
	}
	permissionRes, err := t.ctxBdl.Bdl.CheckPermission(req)
	if err != nil {
		return err
	}
	if permissionRes == nil {
		return fmt.Errorf("can not check permission for create project")
	}
	var visible bool
	if permissionRes.Access {
		visible = true
		t.Props.Operations = []interface{}{
			map[string]interface{}{
				"props": map[string]interface{}{
					"text":        i18nLocale.Get(i18n.I18nKeyCreate),
					"visible":     visible,
					"disabled":    false,
					"disabledTip": i18nLocale.Get(i18n.I18nKeyProjectNoPermission),
					"type":        "link",
				},
				"operations": map[string]interface{}{
					"click": map[string]interface{}{
						"command": map[string]interface{}{
							"key":     "goto",
							"target":  "createProject",
							"jumpOut": false,
							"visible": visible,
						},
						"key":    "click",
						"reload": false,
						"show":   false,
					},
				},
			},
		}
	}
	return nil
}

func RenderCreator() protocol.CompRender {
	return &MyProjectTitle{}
}
