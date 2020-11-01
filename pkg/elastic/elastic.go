package elastic

import (
	"context"
	"fmt"
	elasticapi "github.com/olivere/elastic/v7"
)

type ConfigT struct {
	// 访问地址
	URL string `validate:"required"`

	// 用户名
	User string

	// 密码
	Passwd string
}

func (ct *ConfigT) Invalid() bool {
	return len(ct.URL) == 0
}

func NewElasticClient(ctx context.Context, config *ConfigT) (*elasticapi.Client, error) {
	if config.Invalid() {
		return nil, InvalidConfigError
	}
	clientOpts := []elasticapi.ClientOptionFunc{elasticapi.SetURL(config.URL), elasticapi.SetSniff(false)}
	if len(config.User) > 0 && len(config.Passwd) > 0 {
		clientOpts = append(clientOpts, elasticapi.SetBasicAuth(config.User, config.Passwd))
	}

	client, err := elasticapi.NewClient(clientOpts...)
	if err != nil {
		return nil, err
	}
	err = ping(ctx, client, config.URL)
	if err != nil {
		return nil, err
	}

	return client, nil

}

func ping(ctx context.Context, client *elasticapi.Client, url string) error {

	if client != nil {
		info, code, err := client.Ping(url).Do(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Elasticsearch returned with code %d and version %s \n", code, info.Version.Number)
		return nil
	}
	return InvalidClientError
}
