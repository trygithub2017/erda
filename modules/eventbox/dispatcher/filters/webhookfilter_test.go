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

package filters

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/erda-project/erda/modules/eventbox/constant"
	"github.com/erda-project/erda/modules/eventbox/types"
	"github.com/erda-project/erda/modules/eventbox/webhook"
)

func testWebhookFilter(t *testing.T, f Filter) {
	m := types.Message{
		Sender:  "self",
		Content: "233",
		Labels: map[types.LabelKey]interface{}{
			types.LabelKey(constant.WebhookLabelKey): webhook.EventLabel{
				Event:         "test-event",
				OrgID:         "1",
				ProjectID:     "2",
				ApplicationID: "3",
			},
		},
		Time: 0,
	}
	derr := f.Filter(&m)
	assert.True(t, derr.IsOK())

	http := m.Labels[types.LabelKey("/HTTP")]
	assert.NotNil(t, http, fmt.Sprintf("%+v", m))

	raw, err := json.Marshal(m.Content)
	assert.Nil(t, err)
	em := webhook.EventMessage{}
	assert.Nil(t, json.Unmarshal(raw, &em))
	assert.Equal(t, "test-event", em.Event)
	assert.Equal(t, "", em.Env)
}

func testWebhookFilterWithEnv(t *testing.T, f Filter) {
	m := types.Message{
		Sender:  "self",
		Content: "233",
		Labels: map[types.LabelKey]interface{}{
			types.LabelKey(constant.WebhookLabelKey): webhook.EventLabel{
				Event:         "test-event",
				OrgID:         "1",
				ProjectID:     "2",
				ApplicationID: "3",
				Env:           "test",
			},
		},
		Time: 0,
	}
	derr := f.Filter(&m)
	assert.True(t, derr.IsOK())

	http := m.Labels[types.LabelKey("/HTTP")]
	assert.NotNil(t, http, fmt.Sprintf("%+v", m))

	raw, err := json.Marshal(m.Content)
	assert.Nil(t, err)
	em := webhook.EventMessage{}
	assert.Nil(t, json.Unmarshal(raw, &em))
	assert.Equal(t, "test", em.Env)
}

func testWebhookFilterDINGDINGURL(t *testing.T, f Filter) {
	m := types.Message{
		Sender:  "self",
		Content: "2333",
		Labels: map[types.LabelKey]interface{}{
			types.LabelKey(constant.WebhookLabelKey): webhook.EventLabel{
				Event:         "test-event3",
				OrgID:         "1",
				ProjectID:     "2",
				ApplicationID: "3",
				Env:           "test",
			},
		},
	}
	derr := f.Filter(&m)
	assert.True(t, derr.IsOK())

	dd := m.Labels[types.LabelKey("/DINGDING")]
	assert.NotNil(t, dd, fmt.Sprintf("%+v", m))

}

// func TestWebhookFilter(t *testing.T) {
// 	impl, err := webhook.NewWebHookImpl()
// 	assert.Nil(t, err)
// 	r, err := impl.CreateHook("1", webhook.CreateHookRequest{
// 		Name:   "xxx",
// 		Events: []string{"test-event", "test-event2"},
// 		URL:    "http://test-url",
// 		Active: true,
// 		HookLocation: apistructs.HookLocation{
// 			Org:         "1",
// 			Project:     "2",
// 			Application: "3",
// 			Env:         []string{"test"},
// 		},
// 	})
// 	assert.Nil(t, err)
// 	r2, err := impl.CreateHook("1", webhook.CreateHookRequest{
// 		Name:   "yyy",
// 		Events: []string{"test-event3", "test-event4"},
// 		URL:    "https://oapi.dingtalk.com/robot/send?access_token=xxxx",
// 		Active: true,
// 		HookLocation: apistructs.HookLocation{
// 			Org:         "1",
// 			Project:     "2",
// 			Application: "3",
// 			Env:         []string{"test"},
// 		},
// 	})
// 	assert.Nil(t, err)
// 	defer func() {
// 		assert.Nil(t, impl.DeleteHook("1", string(r)))
// 		assert.Nil(t, impl.DeleteHook("1", string(r2)), string(r2))
// 	}()
// 	f, err := NewWebhookFilter()
// 	assert.Nil(t, err)

// 	t.Run("normal", func(t *testing.T) {
// 		testWebhookFilter(t, f)
// 	})
// 	t.Run("env", func(t *testing.T) {
// 		testWebhookFilterWithEnv(t, f)
// 	})
// 	t.Run("dingding", func(t *testing.T) {
// 		testWebhookFilterDINGDINGURL(t, f)
// 	})
// }
