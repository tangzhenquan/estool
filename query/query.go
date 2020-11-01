package query

import (
	"context"
	"encoding/json"
	elasticDao "esTool/dao/elastic"
	"reflect"
	"time"
)

type Queryer struct {
	dao *elasticDao.DAO
}

func NewQueryer(dao  *elasticDao.DAO) *Queryer {
	return &Queryer{
		dao: dao,
	}
}
type QueryRes map[string]interface{}

func (qr QueryRes)String()string  {
	data , _ := json.Marshal(qr)
	return string(data)
}

func (queryer *Queryer) QueryByIds(ctx context.Context,indexName string, ids, docType []string)([]QueryRes,  error ) {
	subCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	rows , err :=queryer.dao.QueryByIds(subCtx, indexName, docType, ids, reflect.TypeOf(QueryRes{}), nil, "")
	if err != nil{
		return nil, err
	}
	var res []QueryRes
	for _, row := range rows {
		queryItemRes := row.(QueryRes)
		res = append(res, queryItemRes)
	}
	return res, nil
}