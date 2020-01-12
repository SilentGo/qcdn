package modules

import (
	"encoding/json"
	"fmt"
	"qcdn/cdn"
	"strconv"
	"strings"
)

type Client struct {
	Detail    int
	Offset    int
	Limit     int
	Hosts     string
	SecretID  string
	SecretKey string
}

type DescribeCdnHostsRsp struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	CodeDesc string `json:"codeDesc"`
	Data     struct {
		Hosts []struct {
			ID          int    `json:"id"`
			AppID       int    `json:"app_id"`
			OwnerUin    int64  `json:"owner_uin"`
			ProjectID   int    `json:"project_id"`
			Host        string `json:"host"`
			HostType    string `json:"host_type"`
			ServiceType string `json:"service_type"`
			Origin      string `json:"origin"`
			Cache       []struct {
				Type int    `json:"type"`
				Rule string `json:"rule"`
				Time int    `json:"time"`
				Unit string `json:"unit"`
			} `json:"cache"`
			Status         int    `json:"status"`
			Disabled       int    `json:"disabled"`
			Message        string `json:"message"`
			CreateTime     string `json:"create_time"`
			UpdateTime     string `json:"update_time"`
			Deleted        string `json:"deleted"`
			FwdHostType    string `json:"fwd_host_type"`
			FwdHost        string `json:"fwd_host"`
			MiddleResource int    `json:"middle_resource"`
			Refer          struct {
				Type     int           `json:"type"`
				List     []interface{} `json:"list"`
				NullFlag int           `json:"null_flag"`
			} `json:"refer"`
			Readonly    int           `json:"readonly"`
			TegStatus   string        `json:"teg_status"`
			Cname       string        `json:"cname"`
			CacheMode   string        `json:"cache_mode"`
			FurlCache   string        `json:"furl_cache"`
			HTTP2       int           `json:"http2"`
			SslType     int           `json:"ssl_type"`
			PidConfig   []interface{} `json:"pid_config"`
			HTTPSConfig struct {
				HTTP2             int    `json:"http2"`
				Spdy              string `json:"spdy"`
				SslStapling       string `json:"ssl_stapling"`
				SslStaplingVerify string `json:"ssl_stapling_verify"`
				SslVerifyClient   string `json:"ssl_verify_client"`
			} `json:"https_config"`
			SslDeployTime interface{} `json:"ssl_deploy_time"`
			SslExpireTime interface{} `json:"ssl_expire_time"`
			SslCertName   string      `json:"ssl_cert_name"`
			SslCertID     string      `json:"ssl_cert_id"`
			Seo           string      `json:"seo"`
			TestURL       string      `json:"test_url"`
			HostID        int         `json:"host_id"`
			Capping       struct {
				Bandwidth int64  `json:"bandwidth"`
				Unit      string `json:"unit"`
				Hy        string `json:"hy"`
				Active    string `json:"active"`
			} `json:"capping"`
		} `json:"hosts"`
		Total int `json:"total"`
	} `json:"data"`
}

func NewClient(secretid, secretkey string, hosts string, detail, offset, limit int) *Client {
	return &Client{
		Detail:    detail,
		Offset:    offset,
		Limit:     limit,
		Hosts:     hosts,
		SecretKey: secretkey,
		SecretID:  secretid,
	}
}

func (c *Client) DescribeCdnHosts() *DescribeCdnHostsRsp {
	rsp := new(DescribeCdnHostsRsp)
	params := make(map[string]interface{})
	params["Action"] = "DescribeCdnHosts"
	params["detail"] = c.Detail
	params["offset"] = c.Offset
	params["limit"] = c.Limit
	params["SecretId"] = c.SecretID

	request, _ := cdn.Signature(c.SecretKey, params)
	response := cdn.SendRequest(request)

	err := json.Unmarshal([]byte(response), &rsp)
	if err != nil {
		fmt.Println(err)
	}

	return rsp
}

func (c *Client) GetHostInfoByHost() *DescribeCdnHostsRsp {
	rsp := new(DescribeCdnHostsRsp)
	hosts := strings.Split(c.Hosts, ";")
	params := make(map[string]interface{})
	params["Action"] = "GetHostInfoByHost"
	params["SecretId"] = c.SecretID
	for index, host := range hosts {
		str := "hosts." + strconv.Itoa(index)
		host = strings.TrimSpace(host)
		params[str] = host
	}

	request, _ := cdn.Signature(c.SecretKey, params)
	respone := cdn.SendRequest(request)

	err := json.Unmarshal([]byte(respone), &rsp)
	if err != nil {
		fmt.Print(err)
	}
	return rsp
}
