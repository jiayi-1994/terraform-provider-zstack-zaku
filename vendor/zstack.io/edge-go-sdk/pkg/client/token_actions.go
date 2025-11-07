package client

func (cli *ZeClient) GetToken() (string, error) {
	var resp string
	err := cli.Get("/open-api/token", "", nil, &resp)
	return resp, err
}
