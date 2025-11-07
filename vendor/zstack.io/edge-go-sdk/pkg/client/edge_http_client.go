package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/golog"

	"zstack.io/edge-go-sdk/pkg/auth"
	"zstack.io/edge-go-sdk/pkg/errors"
	"zstack.io/edge-go-sdk/pkg/param"
	"zstack.io/edge-go-sdk/pkg/util/httputils"
	"zstack.io/edge-go-sdk/pkg/util/jsonutils"
)

const (
	responseKeyContent  = "content"
	responseKeyTotal    = "totalCount"
	responseKeyResult   = "result"
	responseKeyActionID = "actionId"
)

type ZeHttpClient struct {
	*ZeConfig

	httpClient *http.Client
}

func NewZeHttpClient(config *ZeConfig) *ZeHttpClient {
	rt := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: config.insecureSkipVerify},
	}
	//初始化httpClient
	httpClient := &http.Client{
		Timeout: config.timeout,
		Transport: &auth.ZeAuthProviderTransport{
			Ak:           config.accessKeyId,
			Sk:           config.accessKeySecret,
			ContextPath:  config.contextPath,
			RoundTripper: rt,
		},
	}
	return &ZeHttpClient{
		ZeConfig:   config,
		httpClient: httpClient,
	}
}

func (cli *ZeHttpClient) Page(resource string, params *param.QueryParam, retVal interface{}) (int, error) {
	return cli.PageWithRespKey(resource, responseKeyContent, params, retVal)
}

func (cli *ZeHttpClient) PageWithRespKey(resource, responseKey string, params *param.QueryParam, retVal interface{}) (int, error) {
	params.ReplyWithCount(true)
	err := cli.ListWithRespKey(resource, responseKey, params, retVal)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(params.Get(responseKeyTotal))
}

func (cli *ZeHttpClient) List(resource string, params *param.QueryParam, retVal interface{}) error {
	return cli.ListWithRespKey(resource, responseKeyContent, params, retVal)
}

func (cli *ZeHttpClient) ListWithRespKey(resource, responseKey string, params *param.QueryParam, retVal interface{}) error {
	return cli.ListWithSpec(resource, "", "", responseKey, params, retVal)
}

func (cli *ZeHttpClient) ListWithSpec(resource, resourceId, spec, responseKey string, params *param.QueryParam, retVal interface{}) error {
	var urlStr string
	if params != nil {
		urlStr = cli.getGetURL(resource, resourceId, spec, params.Values)
	} else {
		urlStr = cli.getGetURL(resource, resourceId, spec, nil)
	}

	_, resp, err := cli.httpGet(urlStr)
	if err != nil {
		return err
	}

	if len(responseKey) == 0 {
		return resp.Unmarshal(retVal)
	}

	if params != nil && len(params.Get("replyWithCount")) > 0 {
		total, err := resp.GetString(responseKey, responseKeyTotal)
		if err != nil {
			return err
		}

		params.Set(responseKeyTotal, total)

		return resp.Unmarshal(retVal, responseKey, responseKeyResult)
	}

	return resp.Unmarshal(retVal, responseKey)
}

func (cli *ZeHttpClient) Get(resource, resourceId string, params interface{}, retVal interface{}) error {
	return cli.GetWithRespKey(resource, resourceId, responseKeyContent, params, retVal)
}

func (cli *ZeHttpClient) GetWithRespKey(resource, resourceId, responseKey string, params interface{}, retVal interface{}) error {
	return cli.GetWithSpec(resource, resourceId, "", responseKey, params, retVal)
}

func (cli *ZeHttpClient) GetWithSpec(resource, resourceId, spec, responseKey string, params interface{}, retVal interface{}) error {
	var urlStr string
	if params != nil {
		urlValues, err := param.ConvertStruct2UrlValues(params)
		if err != nil {
			return err
		}
		urlStr = cli.getGetURL(resource, resourceId, spec, urlValues)
	} else {
		urlStr = cli.getGetURL(resource, resourceId, spec, nil)
	}
	_, resp, err := cli.httpGet(urlStr)
	if err != nil {
		return err
	}

	if retVal == nil {
		return nil
	}

	if len(responseKey) == 0 {
		return resp.Unmarshal(retVal)
	}

	return resp.Unmarshal(retVal, responseKey)
}

func (cli *ZeHttpClient) httpGet(urlStr string) (http.Header, jsonutils.JSONObject, error) {
	var respHeader http.Header
	var resp jsonutils.JSONObject
	startTime := time.Now()
	for time.Since(startTime) < 5*time.Minute {
		_, httpRespHeader, httpResp, err := httputils.JSONRequest(cli.httpClient, context.TODO(), httputils.GET, urlStr, nil, cli.debug)
		if err != nil {
			if strings.Contains(err.Error(), "exceeded while awaiting headers") {
				time.Sleep(time.Second * 5)
				continue
			}
			return nil, nil, errors.Wrapf(err, fmt.Sprintf("%s %s", http.MethodGet, urlStr))
		}

		respHeader = httpRespHeader
		resp = httpResp
		break
	}
	return respHeader, resp, nil
}

func (cli *ZeHttpClient) getGetURL(resource, resourceId, spec string, urlValues url.Values) string {
	url := cli.getURL(resource, resourceId, spec)
	if len(urlValues) > 0 {
		url = fmt.Sprintf("%s?%s", url, urlValues.Encode())
	}
	return url
}

// //////////////////////////// Post(创建) ///////////////////////

func (cli *ZeHttpClient) Post(resource string, params interface{}, retVal interface{}) error {
	return cli.PostWithRespKey(resource, "", params, retVal)
}

func (cli *ZeHttpClient) PostWithRespKey(resource, responseKey string, params interface{}, retVal interface{}) error {
	return cli.PostWithSpec(resource, "", "", responseKey, params, retVal)
}

func (cli *ZeHttpClient) PostWithSpec(resource, resourceId, spec, responseKey string, params interface{}, retVal interface{}) error {
	_, err := cli.PostWithAsync(resource, resourceId, spec, responseKey, params, retVal, false)
	return err
}
func (cli *ZeHttpClient) PostWithAsync(resource, resourceId, spec, responseKey string, params interface{}, retVal interface{}, async bool) (string, error) {
	urlStr := cli.getPostURL(resource, resourceId, spec)
	location, _, resp, err := cli.httpPost(urlStr, jsonMarshal(params), async)
	if err != nil {
		return location, err
	}

	if async || retVal == nil {
		return location, nil
	}

	if len(responseKey) == 0 {
		return location, resp.Unmarshal(retVal)
	}

	return location, resp.Unmarshal(retVal, responseKey)
}

func (cli *ZeHttpClient) httpPost(urlStr string, params jsonutils.JSONObject, async bool) (string, http.Header, jsonutils.JSONObject, error) {
	header, respHeader, resp, err := httputils.JSONRequest(cli.httpClient, context.TODO(), httputils.POST, urlStr, params, cli.debug)
	if err != nil {
		return "", nil, nil, errors.Wrapf(err, fmt.Sprintf("%s %s %s", http.MethodPost, urlStr, params.String()))
	}

	var actionID string

	if resp != nil {
		content, err := resp.Get(responseKeyContent)
		if err == nil && content.Contains(responseKeyActionID) {
			actionID, _ = content.GetString(responseKeyActionID)
		}
		if !async {
			resultHeader, result, err := cli.httpWait(header, http.MethodPost, urlStr, params, actionID)
			return actionID, resultHeader, result, err
		}
	}

	return actionID, respHeader, resp, nil

}

func (cli *ZeHttpClient) getPostURL(resource, resourceId, spec string) string {
	return cli.getURL(resource, resourceId, spec)
}

////////////////////////////// Put(更新) ///////////////////////

func (cli *ZeHttpClient) Put(resource, resourceId string, params interface{}, retVal interface{}) error {
	return cli.PutWithRespKey(resource, resourceId, "", params, retVal)
}

func (cli *ZeHttpClient) PutWithRespKey(resource, resourceId, responseKey string, params interface{}, retVal interface{}) error {
	return cli.PutWithSpec(resource, resourceId, "actions", responseKey, params, retVal)
}

func (cli *ZeHttpClient) PutWithSpec(resource, resourceId, spec, responseKey string, params interface{}, retVal interface{}) error {
	_, err := cli.PutWithAsync(resource, resourceId, spec, responseKey, params, retVal, false)
	return err
}

func (cli *ZeHttpClient) PutWithAsync(resource, resourceId, spec, responseKey string, params interface{}, retVal interface{}, async bool) (string, error) {
	urlStr := cli.getPutURL(resource, resourceId, spec)
	location, _, resp, err := cli.httpPut(urlStr, jsonMarshal(params), async)
	if err != nil {
		return location, err
	}

	if async || retVal == nil {
		return location, nil
	}

	if len(responseKey) == 0 {
		return location, resp.Unmarshal(retVal)
	}

	return location, resp.Unmarshal(retVal, responseKey)
}

func (cli *ZeHttpClient) httpPut(urlStr string, params jsonutils.JSONObject, async bool) (string, http.Header, jsonutils.JSONObject, error) {
	header, respHeader, resp, err := httputils.JSONRequest(cli.httpClient, context.Background(), httputils.PUT, urlStr, params, cli.debug)
	if err != nil {
		return "", nil, nil, errors.Wrapf(err, fmt.Sprintf("%s %s %s", http.MethodPut, urlStr, params.String()))
	}
	var actionID string

	if resp != nil {
		content, err := resp.Get(responseKeyContent)
		if err == nil && content.Contains(responseKeyActionID) {
			actionID, _ = content.GetString(responseKeyActionID)
		}
		if !async {
			resultHeader, result, err := cli.httpWait(header, http.MethodPut, urlStr, params, actionID)
			return actionID, resultHeader, result, err
		}
	}

	return actionID, respHeader, resp, nil
}

func (cli *ZeHttpClient) getPutURL(resource, resourceId, spec string) string {
	return cli.getURL(resource, resourceId, spec)
}

////////////////////////////// Delete(删除) ///////////////////////

func (cli *ZeHttpClient) Delete(resource, resourceId, deleteMode string) error {
	return cli.DeleteWithSpec(resource, resourceId, "", fmt.Sprintf("deleteMode=%s", deleteMode), nil)
}

func (cli *ZeHttpClient) DeleteWithSpec(resource, resourceId, spec, paramsStr string, retVal interface{}) error {
	_, err := cli.DeleteWithAsync(resource, resourceId, spec, paramsStr, retVal, false)
	return err
}

func (cli *ZeHttpClient) DeleteWithAsync(resource, resourceId, spec, paramsStr string, retVal interface{}, async bool) (string, error) {
	urlStr := cli.getDeleteURL(resource, resourceId, spec, paramsStr)
	location, _, resp, err := cli.httpDelete(urlStr, async)
	if err != nil {
		return location, err
	}

	if async || retVal == nil {
		return location, nil
	}

	return location, resp.Unmarshal(retVal, "")
}

func (cli *ZeHttpClient) httpDelete(urlStr string, async bool) (string, http.Header, jsonutils.JSONObject, error) {
	header, respHeader, resp, err := httputils.JSONRequest(cli.httpClient, context.Background(), httputils.DELETE, urlStr, nil, cli.debug)
	if err != nil {
		return "", nil, nil, errors.Wrapf(err, fmt.Sprintf("%s %s", http.MethodDelete, urlStr))
	}
	var actionID string
	if resp != nil {
		content, err := resp.Get(responseKeyContent)
		if err == nil && content.Contains(responseKeyActionID) {
			actionID, _ = content.GetString(responseKeyActionID)
		}
		if !async {
			resultHeader, result, err := cli.httpWait(header, http.MethodDelete, urlStr, jsonutils.NewDict(), actionID)
			return actionID, resultHeader, result, err
		}
	}

	return actionID, respHeader, resp, nil
}

func (cli *ZeHttpClient) getDeleteURL(resource, resourceId, spec, paramsStr string) string {
	url := cli.getURL(resource, resourceId, spec)
	if len(paramsStr) > 0 {
		if strings.Contains(url, "?") {
			url = fmt.Sprintf("%s&%s", url, paramsStr)
		} else {
			url = fmt.Sprintf("%s?%s", url, paramsStr)
		}
	}
	return url
}

////////////////////////////// 公共方法 ///////////////////////

func (cli *ZeHttpClient) do(method, url string, body io.Reader) (*http.Response, jsonutils.JSONObject, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	resp, err := cli.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	rbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	rbody = bytes.TrimSpace(rbody)

	var jrbody jsonutils.JSONObject = nil
	if len(rbody) > 0 && (rbody[0] == '{' || rbody[0] == '[') {
		var err error
		jrbody, err = jsonutils.Parse(rbody)
		if err != nil && cli.debug {
			// ignore the error
			fmt.Fprintf(os.Stderr, "parsing json failed: %s", err)
		}
	}
	return resp, jrbody, nil
}

func (cli *ZeHttpClient) getURL(resource, resourceId, spec string) string {
	url := cli.getRequestURL(resource)
	if len(resourceId) > 0 {
		url = fmt.Sprintf("%s/%s", url, resourceId)
		if len(spec) > 0 {
			url = fmt.Sprintf("%s/%s", url, spec)
		}
	}
	return url
}

func (cli *ZeHttpClient) getRequestURL(resource string) string {
	return fmt.Sprintf("%s://%s:%d%s%s", cli.protocol, cli.hostname, cli.port, cli.contextPath, resource)
}

func jsonMarshal(params interface{}) jsonutils.JSONObject {
	return jsonutils.Marshal(params)
}

func (cli *ZeHttpClient) httpWait(header http.Header, action string, requestURL string, params jsonutils.JSONObject, actionID string) (http.Header, jsonutils.JSONObject, error) {
	return retryCallback(func() (http.Header, jsonutils.JSONObject, error) {
		location := cli.getGetURL("/open-api/v1/result", actionID, "", nil)
		resp, err := httputils.Request(cli.httpClient, context.TODO(), httputils.GET, location, header, nil, cli.debug)
		if err != nil {
			return nil, nil, errors.Wrap(err, fmt.Sprintf("wait location %s", location))
		}

		if resp.StatusCode != 200 {
			if resp.StatusCode == 202 {
				httputils.CloseResponse(resp)
				return nil, nil, errors.NewJobRunningError(fmt.Sprintf("StatusCode: %d, Job Still Running", resp.StatusCode))
			}
			_, result, err := parseJSONResponseForHttpWait(resp, cli.debug)
			return nil, nil, fmt.Errorf("StatusCode: %d, Reponse: %v, Error: %v", resp.StatusCode, result, err)
		}

		return parseJSONResponseForHttpWait(resp, cli.debug)
	}, action, requestURL, params.String(), cli.retryInterval, cli.retryTimes)
}

func parseJSONResponseForHttpWait(resp *http.Response, debug bool) (http.Header, jsonutils.JSONObject, error) {
	resultHeader, result, err := httputils.ParseJSONResponse("", resp, debug, nil)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			golog.Errorf("httpWait not found %s", err.Error())
			return nil, nil, err //errors.ErrNotFound
		}

		return nil, nil, err
	}
	return resultHeader, result, nil
}

func retryCallback(fn func() (http.Header, jsonutils.JSONObject, error), action, requestURL, params string, interval, retryTimes int) (http.Header, jsonutils.JSONObject, error) {
	for {
		header, body, err := fn()
		if err == nil {
			return header, body, nil
		}

		if retryTimes == 0 {
			return header, body, err
		}
		if !errors.IsJobRunningError(err) {
			return header, body, err
		}

		golog.Debugf("Wait for job %s %s %s complete , lastest result ： %s", action, requestURL, params, err.Error())
		time.Sleep(time.Duration(interval) * time.Second)
		retryTimes--
	}
}
