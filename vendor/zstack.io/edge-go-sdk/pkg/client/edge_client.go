package client

type ZeClient struct {
	*ZeHttpClient
}

func NewZeClient(config *ZeConfig) *ZeClient {
	return &ZeClient{
		ZeHttpClient: NewZeHttpClient(config),
	}
}
