package main

import (
	"context"
	"esTool/pkg/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func genData(ctx context.Context, count int, fileName string) {
	rand.Seed(time.Now().UnixNano())
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for i := 0; i < count; i++ {
		if ctx.Err() != nil {
			return
		}
		add := time.Now().Add(time.Second * time.Duration(rand.Int31n(10000)))
		_, err := file.WriteString(fmt.Sprintf("%s#12.12.12.12#%s#10.2222\n", add.Format("2006-01-02 15:04:05"), strconv.Itoa(int(rand.Int31n(200)))))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app := &cli.App{
		Name:    "tool",
		Usage:   "tools ",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:  "count",
				Value: 10000000,
			}, &cli.StringFlag{
				Name:  "fileName",
				Value: "testData.txt",
			},
		}, Action: func(c *cli.Context) error {
			stopCh := make(chan struct{})
			utils.SafeExecFunc(func(i ...interface{}) {
				genData(ctx, int(c.Int64("count")), c.String("fileName"))
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

		}}
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
