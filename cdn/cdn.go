package cdn

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	Method  = "POST"
	BaseUrl = "https://cdn.api.qcloud.com/v2/index.php"
)

func Signature(secretKey string, params map[string]interface{}) (string, map[string]interface{}) {
	/*add common params*/
	timestamp := time.Now().Unix()
	rd := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000)
	params["Timestamp"] = timestamp
	params["Nonce"] = rd
	/**sort all the params to make signPlainText**/
	sigUrl := Method + strings.TrimLeft(BaseUrl, "https://") + "?"
	sigParam := ""
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	isfirst := true
	for _, key := range keys {
		if !isfirst {
			sigUrl = sigUrl + "&"
			sigParam = sigParam + "&"
		}
		isfirst = false
		if strings.Contains(key, "_") {
			strings.Replace(key, ".", "_", -1)
		}
		value := typeSwitcher(params[key])
		sigUrl = sigUrl + key + "=" + value
		sigParam = sigParam + key + "=" + value
	}
	//fmt.Println("signPlainText: ", sigUrl)
	unencode_sign, _sign := sign(sigUrl, secretKey)
	sigParam = "Signature=" + _sign + "&" + sigParam
	params["Signature"] = unencode_sign
	return sigParam, params
}

func SendRequest(params string) string {
	requesturl := BaseUrl
	var response string
	if Method == "GET" {
		params_str := "?" + params
		requesturl = requesturl + params_str
		response = httpGet(requesturl)
	} else if Method == "POST" {
		res, err := httpPost(requesturl, params)
		if err != nil {
			println(err.Error())
			return err.Error()
		}

		response = string(res)
	} else {
		fmt.Println("unsuppported http method")
		return "unsuppported http method"
	}
	return response
}

func typeSwitcher(t interface{}) string {
	switch v := t.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	case int64:
		return strconv.Itoa(int(v))
	default:
		return ""
	}
}

func sign(signPlainText string, secretKey string) (string, string) {
	key := []byte(secretKey)
	hash := hmac.New(sha1.New, key)
	hash.Write([]byte(signPlainText))
	sig := base64.StdEncoding.EncodeToString([]byte(string(hash.Sum(nil))))
	encd_sig := url.QueryEscape(sig)
	return sig, encd_sig
}

func httpGet(url string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(3) * time.Second}
	fmt.Println(url)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer resp.Body.Close()
	body, erro := ioutil.ReadAll(resp.Body)
	if erro != nil {
		fmt.Println("http wrong error")
		return erro.Error()
	}
	return string(body)
}

func httpPost(requesturl string, params string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", requesturl, strings.NewReader(params))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
