package client

import (
	"zstack.io/edge-go-sdk/pkg/view"
)

func (cli *ZeClient) GetActionResult(actionId string) (map[string]interface{}, error) {
	var resp map[string]interface{}
	return resp, cli.Get("/open-api/v1/result", actionId, nil, &resp)
}

func (cli *ZeClient) ListAuthorizedProject() ([]view.UserProjectSimpleView, error) {
	var resp []view.UserProjectSimpleView
	return resp, cli.List("/open-api/v1/authorized-project", nil, &resp)
}
