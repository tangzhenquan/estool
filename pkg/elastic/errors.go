package elastic

import "errors"

var (
	InvalidConfigError        = errors.New("invalid config error")
	InvalidClientError         = errors.New("invalid client error")
	PingClientError            = errors.New("ping client error")

)