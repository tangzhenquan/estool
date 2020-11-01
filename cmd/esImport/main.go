package main

import (
	"context"
	"esTool/config"
	elasticDao "esTool/dao/elastic"
	"esTool/importer"
	elasticPkg "esTool/pkg/elastic"
	"esTool/pkg/logger"
	"esTool/pkg/utils"
	"esTool/reader/filereader"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app := &cli.App{
		Name:    "esImport",
		Usage:   "Import log file into elasticsearch",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "conf",
				Usage: "Config file for esImport",
				Value: "esImport.yaml",
			},
		},
		Action: func(c *cli.Context) error {
			if err := config.InitConfig(c.String("conf")); err != nil {
				log.WithError(err).Error("init config error")
				return err

			}
			esImportConfig := config.GetConfig()
			fmt.Println(esImportConfig)
			if err := logger.InitLogger(&esImportConfig.LoggerConfig); err != nil {
				log.WithError(err).Error("init logger error")
				return err
			}

			client, err := elasticPkg.NewElasticClient(ctx, &esImportConfig.ElasticConfig)
			if err != nil {
				log.WithError(err).WithField("config", esImportConfig.ElasticConfig).Error("new elastic error")
				return err
			}
			reader, err := filereader.NewLineReader(esImportConfig.FilePath, esImportConfig.MaxLineSize)
			if err != nil {
				log.WithError(err).WithField("filePath", esImportConfig.FilePath).Error("new line reader error")
				return err
			}
			defer func() {
				if err = reader.Close(); err != nil {
					log.WithError(err).WithField("filePath", esImportConfig.FilePath).Error("reader close error")
				}
			}()

			newImporter := importer.NewImporter(ctx, elasticDao.NewDAO(client), &esImportConfig.ImportConfig, reader)
			stopCh := make(chan struct{})
			utils.SafeExecFunc(func(i ...interface{}) {
				err = newImporter.Start()
				if err != nil {
					log.WithError(err).Error("start error")
				}
				newImporter.PrintStats()
				stopCh <- struct{}{}
			})

			qC := make(chan os.Signal, 1)
			signal.Notify(qC, syscall.SIGINT, syscall.SIGTERM)

			select {
			case s := <-qC:
				log.WithField("info", s.String()).Info("qC")
			case <-stopCh:
				log.Info("stop and exit")
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal(err)
	}

}
