package main

import (
	"context"
	elasticDao "esTool/dao/elastic"
	elasticPkg "esTool/pkg/elastic"
	"esTool/query"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func main() {
	app := &cli.App{
		Name:    "esQuery",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "indexName",
				Value: "loglog2",
			}, &cli.StringFlag{
				Name:  "docType",
				Value: "log",
			}, &cli.StringFlag{
				Name:  "ids",
				Value: "1,2",
			},&cli.StringFlag{
				Name:  "esUrl",
				Value: "http://192.168.1.194:9200",
			},
		}, Action: func(c *cli.Context) error {
			idsStr := c.String("ids")
			ids := strings.Split(idsStr, ",")
			if len(ids) == 0 {
				return fmt.Errorf("ids invalid")
			}
			elasticConfig := elasticPkg.ConfigT{
				URL: c.String("esUrl"),
			}
			ctx := context.Background()
			client, err := elasticPkg.NewElasticClient(ctx, &elasticConfig)
			if err != nil {
				return err
			}
			queryer := query.NewQueryer(elasticDao.NewDAO(client))
			res , err :=  queryer.QueryByIds(ctx,
				c.String("indexName"),
				ids,
				[]string{c.String("docType")},
			)
			if err != nil{
				return err
			}
			if len(res) > 0{
				for _, item := range res{
					fmt.Println(item.String())
				}
			}else {
				fmt.Println("not found")
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}