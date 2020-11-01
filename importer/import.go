package importer

import (
	"context"
	"encoding/json"
	elasticDao "esTool/dao/elastic"
	"esTool/pkg/utils"
	"esTool/reader"
	"fmt"
	elasticapi "github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"
)

const (
	failedMaxTries = 3
)

type ConfigT struct {
	Fields                []Field `validate:"required,dive"`
	Sep                   string  `validate:"required"`
	IndexName             string  `validate:"required"`
	DocType               string  `validate:"required"`
	MaxBukBufferSize      int     `validate:"required"`
	MaxReadChanBufferSize int     `validate:"required"`
	Mapping               string
}

type MapStr map[string]interface{}

func (ms MapStr) String() string {
	data, _ := json.Marshal(ms)
	return string(data)
}

type Importer struct {
	dao               *elasticDao.DAO
	reader            reader.Reader
	config            *ConfigT
	stopCh            chan interface{}
	ctx               context.Context
	createdIndexCache map[string]interface{}
	bukBuffer         []interface{}
	stats             Stats
}

func NewImporter(ctx context.Context, dao *elasticDao.DAO, config *ConfigT, reader reader.Reader) *Importer {
	res := &Importer{
		dao:    dao,
		reader: reader,
		config: config,
		ctx:    ctx,
	}
	res.stopCh = make(chan interface{})
	res.bukBuffer = make([]interface{}, 0, config.MaxBukBufferSize+1)
	res.createdIndexCache = make(map[string]interface{})
	return res
}

type ReadRes struct {
	mapStr MapStr
	err    error
}

func (imp *Importer) parseLine(line reader.Line) (error, MapStr) {
	strs := strings.Split(line.Text, imp.config.Sep)
	lenStrs := len(strs)
	if lenStrs == 0 {
		return fmt.Errorf("invalid string"), nil
	}
	res := make(MapStr)
	for index, field := range imp.config.Fields {
		if index < lenStrs {
			value, err := field.Value(strs[index])
			if err != nil {
				return err, nil
			}
			res[field.Name] = value
		} else {
			break
		}
	}
	if len(res) > 0 {
		res["parseTime"] = line.Time
	}
	return nil, res

}

func (imp *Importer) Start() error {
	imp.timerPrintStats()
	resCh := imp.readStart()
	for {
		select {
		case <-imp.ctx.Done():
			log.Debug("Importer$Done")
			return nil
		case res := <-resCh:
			if res == nil {
				return imp.writeLast(imp.ctx)
			}
			if res.err != nil {
				return res.err
			}
			err := imp.write(imp.ctx, res.mapStr)
			if err != nil {
				return err
			}
		}
	}
}

func (imp *Importer) PrintStats() {
	imp.stats.print()
}

func (imp *Importer) timerPrintStats() {
	utils.SafeExecFunc(func(i ...interface{}) {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				imp.PrintStats()
			case <-imp.ctx.Done():
				return
			}
		}
	})

}

func (imp *Importer) indexName() string {
	return imp.config.IndexName
}

func (imp *Importer) createIndex(ctx context.Context, indexName string) error {
	if _, ok := imp.createdIndexCache[indexName]; ok {
		return nil
	}
	subCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := imp.dao.CreateIndexIfDoesNotExist(subCtx, indexName, imp.config.Mapping); err != nil {
		return err
	}
	imp.createdIndexCache[indexName] = struct{}{}
	return nil
}

func (imp *Importer) parseBulkAddRes(res *elasticapi.BulkResponse, writeCount int) error {
	failedBulkBuffer := make([]interface{}, 0, imp.config.MaxBukBufferSize+1)
	var (
		success int
		fail    int
	)
	if res.Errors {
		for i, item := range res.Items {
			for _, result := range item {
				if result.Status < 300 {
					success++
					imp.stats.success++
					continue
				} else if result.Status == 409 {
					imp.stats.duplicates++
					continue
				} else if result.Status < 500 {
					if result.Status == 429 {
						imp.stats.tooMany++
					} else {
						log.Debugf("failed status:%d reason:%s", result.Status, result.Error.Reason)
						imp.stats.nonIndexable++
						continue
					}
				}
				fail++
				imp.stats.fail++
				log.Debugf("Bulk item insert failed (i=%v, status=%v): %s", i, result, result.Error.Reason)
				failedBulkBuffer = append(failedBulkBuffer, imp.bukBuffer[i])
			}
		}
	} else {
		success = writeCount
		imp.stats.success = imp.stats.success + success
	}
	log.Debugf("Bulk item insert success %d", success)
	if fail == len(imp.bukBuffer) {
		return CanNotWriteEsError
	}
	imp.bukBuffer = failedBulkBuffer
	return nil
}
func (imp *Importer) write2Es(ctx context.Context) error {
	indexName := imp.indexName()
	if err := imp.createIndex(ctx, indexName); err != nil {
		return err
	}

	if len(imp.bukBuffer) > 0 {
		subCtx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		res, err := imp.dao.BulkAdd(subCtx, indexName, imp.config.DocType, imp.bukBuffer)
		if err != nil {
			return err
		}
		if err = imp.parseBulkAddRes(res, len(imp.bukBuffer)); err != nil {
			return err
		}
	}
	return nil
}
func (imp *Importer) write(ctx context.Context, mapStr MapStr) error {
	imp.bukBuffer = append(imp.bukBuffer, mapStr)
	if len(imp.bukBuffer) >= imp.config.MaxBukBufferSize {
		if err := imp.write2Es(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (imp *Importer) writeLast(ctx context.Context) error {
	for i := 0; i < failedMaxTries; i++ {
		if err := imp.write2Es(ctx); err != nil {
			return err
		}
		if len(imp.bukBuffer) > 0 {
			log.WithField("failed", imp.bukBuffer).Debug("writeLast failed ids")
			continue
		}
		return nil
	}
	if len(imp.bukBuffer) > 0 {
		log.WithField("failed", imp.bukBuffer).Error("writeLast failed plz check log file ")
	}
	return nil
}

func (imp *Importer) readStart() <-chan *ReadRes {
	ch := make(chan *ReadRes, imp.config.MaxReadChanBufferSize)
	utils.SafeExecFunc(func(i ...interface{}) {
		defer func() {
			close(ch)
			log.Debug("readStart$end")
		}()
		for {
			if imp.ctx.Err() != nil {
				break
			}
			line, err := imp.reader.Next()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.WithError(err).Error("reader read error")
					ch <- &ReadRes{nil, err}
				}
			} else {
				if line.Text == "" {
					continue
				}
				err, mapStr := imp.parseLine(line)
				if err != nil {
					log.WithError(err).WithField("text", line.Text).Error("parse fail")
					continue
				}
				ch <- &ReadRes{mapStr, nil}
			}
		}

	})
	return ch
}
