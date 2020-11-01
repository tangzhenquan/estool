package elastic

import (
	"context"
	elasticPkg "esTool/pkg/elastic"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

const (
	url = "http://192.168.1.194:9200"
	indexName = "loglog2"
	docType = "log"
	user = ""
	passwd = ""
)

type ResItem map[string]interface{}

func TestDAO(t *testing.T)  {
	config := elasticPkg.ConfigT{
		URL:  url,
		User: user,
		Passwd: passwd,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var (
		err error
	)
	client, err := elasticPkg.NewElasticClient(ctx, &config)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	dao := NewDAO(client)
	testDAO_CreateIndexIfDoesNotExist(ctx, t, dao)
	testDAO_QueryByIds(ctx, t, dao)
	testDAO_BulkAdd(ctx, t, dao)
}

func testDAO_CreateIndexIfDoesNotExist(ctx context.Context, t *testing.T,dao *DAO) {

	if err := dao.CreateIndexIfDoesNotExist(ctx,indexName, "" ); err != nil{
		t.Error(err)
		t.FailNow()
	}
	return
}


func testDAO_QueryByIds(ctx context.Context, t *testing.T,dao *DAO) {

	//rows, err := dao.IdsQuery(10001, EVENT_GATHER_INDEX, seeyou.EV_DEVICE_ACTIVATE, []string{"867515022483027"}, reflect.TypeOf({}), nil, "")
	rows, err := dao.QueryByIds(ctx, indexName, []string{docType}, []string{"tjxxgnUBwrsCcw-Z-XcT","uDxxgnUBwrsCcw-Z-XcT"}, reflect.TypeOf(ResItem{}), nil, "")
	if err != nil {
		t.Error(err)
	}
	for _, row := range rows {
		res := row.(ResItem)
		t.Log(res)
	}
	return
}



func testDAO_BulkAdd(ctx context.Context, t *testing.T,dao *DAO) {
	var data []interface{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		item := make(map[string]interface{})
		item["time"] = time.Now().Add(time.Second*time.Duration(rand.Int31n(10000)))
		item["ip"] = "12.12.12.12"
		item["count"] = 100
		item["price"] = 10.2222
		data = append(data, item)
	}

	res, err := dao.BulkAdd(ctx, indexName, docType, data)
	if err != nil {
		t.Error(err)
	}
	if res.Errors {
		t.Error("has error")
	}
}