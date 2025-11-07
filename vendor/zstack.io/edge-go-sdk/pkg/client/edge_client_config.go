package client

import (
	"time"
)

const (
	defaultProtocol        = "http"
	defaultEdgePort        = 80
	defaultEdgeContextPath = "/ze"
)

type ZeConfig struct {
	protocol           string
	hostname           string
	port               int
	contextPath        string
	insecureSkipVerify bool

	accessKeyId     string
	accessKeySecret string

	retryInterval int // unit - second
	retryTimes    int

	timeout time.Duration

	debug bool
}

func NewZeConfig(protocol, hostname string, port int, contextPath string) *ZeConfig {
	return &ZeConfig{
		protocol:      protocol,
		hostname:      hostname,
		port:          port,
		contextPath:   contextPath,
		retryInterval: 2,
		retryTimes:    150,
		timeout:       60 * time.Second,
	}
}

func DefaultZeConfig(hostname string) *ZeConfig {
	return NewZeConfig(defaultProtocol, hostname, defaultEdgePort, defaultEdgeContextPath)
}

func (config *ZeConfig) InsecureSkipVerify(insecureSkipVerify bool) *ZeConfig {
	config.insecureSkipVerify = insecureSkipVerify
	return config
}

func (config *ZeConfig) AccessKey(accessKeyId, accessKeySecret string) *ZeConfig {
	config.accessKeyId = accessKeyId
	config.accessKeySecret = accessKeySecret
	return config
}
func (config *ZeConfig) GetAccessKeySecret() string {
	return config.accessKeySecret
}

func (config *ZeConfig) RetryInterval(retryInterval int) *ZeConfig {
	config.retryInterval = retryInterval
	return config
}

func (config *ZeConfig) RetryTimes(retryTimes int) *ZeConfig {
	config.retryTimes = retryTimes
	return config
}

func (config *ZeConfig) Timeout(timeout time.Duration) *ZeConfig {
	config.timeout = timeout
	return config
}

func (config *ZeConfig) Debug(debug bool) *ZeConfig {
	config.debug = debug
	return config
}
