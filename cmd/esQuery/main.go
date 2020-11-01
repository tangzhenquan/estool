package main

import (
	"context"
	elasticDao "esTool/dao/elastic"
	elasticPkg "esTool/pkg/elastic"
	"esTool/query"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

func main() {
	app := &cli.App{
		Name:    "esQuery",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "indexName",
				Value: "loglog102",
			}, &cli.StringFlag{
				Name:  "from",
				Value: "2009-10-10 10:10:10",
			}, &cli.IntFlag{
				Name:  "size",
				Value: 5,
			}, &cli.StringFlag{
				Name:  "esUrl",
				Value: "http://192.168.2.194:9200",
			}, &cli.StringFlag{
				Name:  "esUser",
				Value: "",
			}, &cli.StringFlag{
				Name:  "esPasswd",
				Value: "",
			},
		}, Action: func(c *cli.Context) error {
			from, err := time.Parse("2006-01-02 15:04:05", c.String("from"))
			if err != nil {
				return fmt.Errorf("cant't parse from arg")
			}
			to := time.Now()
			elasticConfig := elasticPkg.ConfigT{
				URL:    c.String("esUrl"),
				User:   c.String("esUser"),
				Passwd: c.String("esPasswd"),
			}
			ctx := context.Background()
			client, err := elasticPkg.NewElasticClient(ctx, &elasticConfig)
			if err != nil {
				return err
			}
			queryer := query.NewQueryer(elasticDao.NewDAO(client))
			res, err := queryer.QueryByTime(ctx,
				c.String("indexName"),
				from,
				to,
				c.Int("size"),
			)
			if err != nil {
				return err
			}
			if len(res) > 0 {
				for _, item := range res {
					fmt.Println(item.String())
				}
			} else {
				fmt.Println("not found")
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
