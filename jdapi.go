package jdsdk

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

const (
	//普通http请求网关
	HttpRouter = "https://router.jd.com/api"

	//设置API读取失败时重试的次数,可以提高API的稳定性,默认为2次
	RestNumeric = 2
)

type ApiReq struct {
	AppKey      string
	AppSecret   string
	AccessToken string
	V           string
	ParamName   string
	CacheLife   int64
	ReqCount    int64
	//GetCache    func(...interface{}) string  //FILE
	//SetCache    func(...interface{}) bool //FILE
	GetCache    func(string) string                         //REDIS
	SetCache    func(interface{}, interface{}, int64) error //REDIS
	WriteErrLog func(ApiLog)
}

type ApiLog struct {
	Id          int64  `json:"id"`
	AppKey      string `json:"app_key"`
	AccessToken string `json:"access_token"`
	V           string `json:"v"`
	Format      string `json:"format"`
	SignMethod  string `json:"sign_method"`
	Timestamp   string `json:"timestamp"`
	Sign        string `json:"sign"`
	Method      string `json:"method"`
	ParamJson   string `json:"param_json"`
	ApiCount    int64  `json:"api_count"`
	ApiErrorInfo
	Result string `json:"result"`
}

type ApiErrorInfo struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

//请求业务参数
type ApiParams map[string]string

func NewClient(appKey, appSecret string) *ApiReq {
	return &ApiReq{
		AppKey:    appKey,
		AppSecret: appSecret,
	}
}

func (client *ApiReq) Execute(method string, params ApiParams) (*gjson.Result, *ApiErrorInfo) {

	apiErrInfo := &ApiErrorInfo{}

	// system params
	value := url.Values{}
	value.Add("app_key", client.AppKey)
	if client.AccessToken != "" {
		value.Add("access_token", client.AccessToken)
	}
	if client.V != "" {
		value.Add("v", client.V)
	} else {
		value.Add("v", "1.0")
	}
	value.Add("format", "json")
	value.Add("sign_method", "md5")
	value.Add("timestamp", GetNow().Format("2006-01-02 15:04:05"))
	value.Add("method", method)

	var paramData interface{}
	if client.ParamName == "" {
		//jd.union.open.goods.promotiongoodsinfo.query 获取推广商品信息接口
		//这里做个兼容。
		paramTmp := make(map[string]string)
		for k, v := range params {
			//value.Add(k, v)
			paramTmp[k] = v
		}
		paramData = paramTmp
	} else {
		paramTmp := make(map[string]ApiParams)
		paramTmp[client.ParamName] = params
		paramData = paramTmp
	}
	paramJson, jsonErr := json.Marshal(paramData)
	if jsonErr != nil {
		apiErrInfo.Code = 62
		apiErrInfo.Message = ApiErrInfo[apiErrInfo.Code].Error()
		return nil, apiErrInfo
	}
	value.Add("param_json", string(paramJson))

	args := []string{}
	cacheParams := []string{}
	for k, v := range value {
		args = append(args, k+v[0])
		if k != "timestamp" {
			cacheParams = append(cacheParams, v[0])
		}
	}

	//缓存key组合
	sort.Strings(cacheParams)
	cacheId := strings.Join(cacheParams, "_")
	cacheId = method + "." + MD5(cacheId)

	//正常返回结果的json节点
	jsonNode := strings.Replace(method, ".", "_", -1) + "_response"
	//先获取缓存
	var cacheData string
	if client.CacheLife > 0 && client.GetCache != nil {
		cacheData = client.GetCache(cacheId)
	}
	if cacheData == "" {
		// make sign
		sort.Strings(args)
		argsStr := strings.Join(args, "")
		value.Add("sign", MD5(client.AppSecret+argsStr+client.AppSecret))

		//开始请求
		client.ReqCount++ //请求次数+1
		response, err := client.httpSend("POST", HttpRouter, value.Encode())
		fmt.Println("response:", string(response))
		if err != nil {
			//重试N次
			if RestNumeric > 0 && client.ReqCount < RestNumeric {
				//fmt.Println("尝试重新请求", client.ReqCount)
				return client.Execute(method, params)
			}
			apiErrInfo.Message = err.Error()
			return nil, apiErrInfo
		}

		apiError := gjson.GetBytes(response, "error_response")
		if apiError.Exists() {
			//如果存在error_response.说明出错了。
			apiErrInfo.Code = apiError.Get("code").Int()
			apiErrInfo.Message = apiError.Get("zh_desc").String()

			if client.WriteErrLog != nil {
				//把错误日志记录到表里。
				//因为本项目会定时查询优惠券有效性。下架商品太多。所以把下架的错误信息排除记录。
				errLog := ApiLog{}
				errLog.AppKey = value.Get("app_key")
				errLog.V = value.Get("v")
				errLog.Format = value.Get("format")
				errLog.SignMethod = value.Get("sign_method")
				errLog.Timestamp = value.Get("timestamp")
				errLog.AccessToken = value.Get("access_token")
				errLog.Sign = value.Get("sign")
				errLog.Method = method
				errLog.ParamJson = value.Get("param_json")
				errLog.ApiCount = 1
				errLog.Code = apiErrInfo.Code
				errLog.Message = apiErrInfo.Message
				errLog.Result = string(response)

				client.WriteErrLog(errLog)
			}
			if apiErrInfo.Code == 2 {
				//调用超限
			}
			if apiErrInfo.Code == 65 || apiErrInfo.Code == 66 {
				//远程服务调用超时
				time.Sleep(500 * time.Millisecond)
				if RestNumeric > 0 && client.ReqCount < RestNumeric {
					//fmt.Println("尝试重新请求1", client.ReqCount)
					return client.Execute(method, params)
				}
			}
			return &apiError, apiErrInfo
		} else {
			//没有error_response，说明是正常的。
			if client.CacheLife > 0 && client.SetCache != nil {
				client.SetCache(cacheId, string(response), client.CacheLife)
			}
			if responseNode := gjson.GetBytes(response, jsonNode+".result"); responseNode.Exists() {
				var respNode string
				if err := json.Unmarshal([]byte(responseNode.Raw), &respNode); err != nil {
					apiErrInfo.Code = 63 //json格式不合法
					apiErrInfo.Message = ApiErrInfo[apiErrInfo.Code].Error()
					return nil, apiErrInfo
				}
				returnResp := gjson.Parse(respNode)
				client.ReqCount = 0 //还原请求次数
				return &returnResp, nil
			} else {
				apiErrInfo.Code = 69 //获取数据失败
				apiErrInfo.Message = ApiErrInfo[apiErrInfo.Code].Error()
				return nil, apiErrInfo
			}
		}
	} else {
		if err := json.Unmarshal([]byte(cacheData), &cacheData); err != nil {
			apiErrInfo.Code = 69 //获取数据失败
			apiErrInfo.Message = ApiErrInfo[apiErrInfo.Code].Error()
			return nil, apiErrInfo
		}
		if responseNode := gjson.Get(cacheData, jsonNode+".result"); responseNode.Exists() {
			var respNode string
			if err := json.Unmarshal([]byte(responseNode.Raw), &respNode); err != nil {
				apiErrInfo.Code = 63 //json格式不合法
				apiErrInfo.Message = ApiErrInfo[apiErrInfo.Code].Error()
				return nil, apiErrInfo
			}
			returnResp := gjson.Parse(respNode)
			client.ReqCount = 0 //还原请求次数
			return &returnResp, nil
		} else {
			apiErrInfo.Code = 69 //获取数据失败
			apiErrInfo.Message = ApiErrInfo[apiErrInfo.Code].Error()
			return nil, apiErrInfo
		}
	}
}

//检查是否有错误码
func (client *ApiReq) CheckApiErr(resp *gjson.Result) *ApiErrorInfo {
	//这里是xxx_yyy_zzz_response.result节点内的内容检查.
	if resp != nil {
		getCode := resp.Get("code")
		if getCode.Exists() && getCode.Int() != 200 {
			errInfo := ApiErrorInfo{}
			if err := json.Unmarshal([]byte(resp.Raw), &errInfo); err != nil {
				errInfo.Code = 69
				errInfo.Message = ApiErrInfo[errInfo.Code].Error()
			}
			return &errInfo
		}
	}
	return nil
}

func (client *ApiReq) httpSend(method, router, param string) ([]byte, error) {
	var req *http.Request
	var err error
	if method == "POST" {
		req, err = http.NewRequest(method, router, strings.NewReader(param))
	} else {
		req, err = http.NewRequest(method, router+"?"+param, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	httpClient := &http.Client{}
	httpClient.Timeout = 3 * time.Second
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求错误:%d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%X", has)
	return md5str
}
