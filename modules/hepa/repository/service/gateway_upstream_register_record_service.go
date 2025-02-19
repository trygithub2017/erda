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

package service

import (
	"github.com/pkg/errors"
	"github.com/xormplus/xorm"

	"github.com/erda-project/erda/modules/hepa/common"
	. "github.com/erda-project/erda/modules/hepa/common/vars"
	"github.com/erda-project/erda/modules/hepa/repository/orm"
)

type GatewayUpstreamRegisterRecordServiceImpl struct {
	engine *orm.OrmEngine
}

func NewGatewayUpstreamRegisterRecordServiceImpl() (*GatewayUpstreamRegisterRecordServiceImpl, error) {
	engine, error := orm.GetSingleton()
	if error != nil {
		return nil, errors.Wrap(error, "new GatewayUpstreamRegisterRecordServiceImpl failed")
	}
	return &GatewayUpstreamRegisterRecordServiceImpl{engine}, nil
}

func (impl *GatewayUpstreamRegisterRecordServiceImpl) Insert(session *xorm.Session, item *orm.GatewayUpstreamRegisterRecord) error {
	if session == nil || item == nil || len(item.UpstreamId) == 0 || item.RegisterId == "" {
		return errors.New(ERR_INVALID_ARG)
	}
	_, err := orm.Insert(session, item)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (impl *GatewayUpstreamRegisterRecordServiceImpl) count(upstreamId string) (int64, error) {
	if len(upstreamId) == 0 {
		return 0, errors.New(ERR_INVALID_ARG)
	}
	total, err := orm.Count(impl.engine, &orm.GatewayUpstreamRegisterRecord{}, "upstream_id = ?", upstreamId)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return total, nil
}

func (impl *GatewayUpstreamRegisterRecordServiceImpl) GetPage(upstreamId string, page *common.Page) (*common.PageQuery, error) {
	total, err := impl.count(upstreamId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	page.SetTotalNum(total)
	if total == 0 {
		return &common.PageQuery{Result: []orm.GatewayUpstreamRegisterRecord{}, Page: page}, nil
	}
	var result []orm.GatewayUpstreamRegisterRecord
	err = orm.SelectPage(impl.engine.Desc("create_time"), &result, page, "upstream_id = ?", upstreamId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &common.PageQuery{Result: result, Page: page}, nil
}

func (impl *GatewayUpstreamRegisterRecordServiceImpl) Get(upstreamId string, registerId string) (*orm.GatewayUpstreamRegisterRecord, error) {
	if len(upstreamId) == 0 || registerId == "" {
		return nil, errors.New(ERR_INVALID_ARG)
	}
	record := &orm.GatewayUpstreamRegisterRecord{}
	succ, err := orm.Get(impl.engine, record, "upstream_id = ? and register_id = ?", upstreamId, registerId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !succ {
		return nil, nil
	}
	return record, nil
}
