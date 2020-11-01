package elastic

import (
	"context"
	elasticapi "github.com/olivere/elastic/v7"
	"reflect"
	"time"
)

type DAO struct {
	client *elasticapi.Client
}

func NewDAO(client *elasticapi.Client) *DAO {
	return &DAO{client: client}
}

func (dao *DAO) getClient() *elasticapi.Client {
	return dao.client
}

func (dao *DAO) CreateIndexIfDoesNotExist(ctx context.Context, indexName, mapping string) error {
	client := dao.getClient()
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	indicesCreateService := client.CreateIndex(indexName)
	if mapping != "" {
		indicesCreateService = indicesCreateService.BodyString(mapping)
	}
	res, err := indicesCreateService.Do(ctx)
	if err != nil {
		return err
	}
	if !res.Acknowledged {
		return CreateIndexError
	}

	return nil
}

type BulkAddRes struct {
	Err          error
	FailedIds    []string
	SuccessCount int
	DupIds       []string
	TooManyIds   []string
	HardError    bool
}

func (dao *DAO) BulkAdd(ctx context.Context, indexName, docType string, data []interface{}) (*elasticapi.BulkResponse, error) {
	client := dao.getClient()
	bulkRequest := client.Bulk()
	for _, doc := range data {
		esRequest := elasticapi.NewBulkIndexRequest().
			Index(indexName).Doc(doc).UseEasyJSON(true).OpType("create")
		bulkRequest = bulkRequest.Add(esRequest)
	}
	//r
	return bulkRequest.Do(ctx)
}

func (dao *DAO) QueryByIds(ctx context.Context, index string, tps, ids []string, ttyp reflect.Type, boost *float64, queryName string) (rows []interface{}, err error) {
	client := dao.getClient()
	idsQuery := elasticapi.NewIdsQuery().Ids(ids...)
	if boost != nil {
		idsQuery.Boost(*boost)
	}
	if queryName != "" {
		idsQuery.QueryName(queryName)
	}

	searchResult, err := client.Search().Index(index).Query(idsQuery).Do(ctx)
	if err != nil {
		return
	}
	for _, item := range searchResult.Each(ttyp) {
		rows = append(rows, item)
	}
	return

}

func (dao *DAO) QueryByTime(ctx context.Context, index string, from, to time.Time, ttyp reflect.Type, size int) (rows []interface{}, err error) {
	client := dao.getClient()
	rangeQuery := elasticapi.NewRangeQuery("time")
	rangeQuery.Gte(from).Lte(to)
	searchResult, err := client.Search().Index(index).Size(size).Query(rangeQuery).Do(ctx)
	if err != nil {
		return
	}
	for _, item := range searchResult.Each(ttyp) {
		rows = append(rows, item)
	}
	return

}
