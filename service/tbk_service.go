package service

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
	"github.com/shopspring/decimal"
	"strconv"
)

const (
	defaultVersion = "2.0"
	APIHost        = "https://eco.taobao.com/router/rest?"
	timeFormat     = "2006-01-02 15:04:05"
	tklDecryptUrl = "http://www.taokouling.com/index/taobao_tkljm"
)

type TbkService struct {
	AppId     string
	AppSecret string
	AdZoneId  string
}

type ErrorRep struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

type CommonField struct {
	Title          string          `json:"title"`
	PicUrl         string          `json:"pict_url"`
	ZkFinalPrice   string          `json:"zk_final_price"`
	ItemUrl        string          `json:"item_url"`
	CommissionRate string          `json:"commission_rate"`
	Volume         int             `json:"volume"`
	CouponInfo     string          `json:"coupon_info"`
	ShopTitle      string          `json:"shop_title"`
	ShortTitle     string          `json:"short_title"`
	ItemId         int64           `json:"item_id"`
	ClickUrl       string          `json:"click_url"`
	CouponShareUrl string          `json:"coupon_share_url"`
	Commission     float64         `json:"commission"`
	CouponAmount   json.RawMessage `json:"coupon_amount"`
}
type MaterialItem struct {
	CommonField
	ClickUrl string `json:"click_url"`
}

type SearchItem struct {
	CommonField
	Url string `json:"url"`
}

func NewTbkService(app_id, app_secret, ad_zone_id string) *TbkService {
	return &TbkService{
		AppId:     app_id,
		AppSecret: app_secret,
		AdZoneId:  ad_zone_id,
	}
}

func (self *TbkService) Request(method string, params map[string]string) ([]byte, error) {
	values := url.Values{}
	params["method"] = method
	params = self.mergeParams(params)
	params["sign"] = self.getSign(params)
	for k, v := range params {
		values.Set(k, v)
	}
	resp, err := http.Get(APIHost + values.Encode())
	if err != nil {
		return nil, err
	}
	body := resp.Body
	respData, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	if bytes.Contains(respData, []byte("error_response")) {
		errResp := struct {
			ErrorResponse ErrorRep `json:"error_response"`
		}{}
		err = json.Unmarshal(respData, &errResp)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errResp.ErrorResponse.Msg)
	}
	return respData, nil
}

func (self *TbkService) mergeParams(params map[string]string) map[string]string {
	if _, ok := params["v"]; !ok {
		params["v"] = self.getVersion()
	}
	//签名方式
	params["format"] = "json"
	params["sign_method"] = "md5"
	params["app_key"] = self.AppId
	params["adzone_id"] = self.AdZoneId
	params["timestamp"] = self.getTimestamp()
	return params
}

func (self *TbkService) getSign(params map[string]string) string {
	keys := make([]string, 0)
	for k, v := range params {
		if v != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	buffer := bytes.NewBufferString(self.AppSecret)
	for _, k := range keys {
		buffer.WriteString(k + params[k])
	}
	buffer.WriteString(self.AppSecret)
	h := md5.New()
	h.Write(buffer.Bytes())
	sign := hex.EncodeToString(h.Sum(nil))
	return strings.ToUpper(sign)
}

func (self *TbkService) getTimestamp() string {
	return time.Now().Format(timeFormat)
}

func (self *TbkService) getVersion() string {
	return defaultVersion
}

//搜索
func (self *TbkService) Search(params map[string]string) ([]SearchItem, error) {
	resp, err := self.Request("taobao.tbk.dg.material.optional", params)
	if err != nil {
		return nil, err
	}
	respStruct := struct {
		Response struct {
			Total_results int `json:"total_results"`
			ResultList    struct {
				MapData []SearchItem `json:"map_data"`
			} `json:"result_list"`
		} `json:"tbk_dg_material_optional_response"`
	}{}
	err = json.Unmarshal(resp, &respStruct)
	if err != nil {
		return nil, err
	}
	data := respStruct.Response.ResultList.MapData
	for k, v := range data {
		amount := self.parseAmount(v.CouponAmount)
		data[k].ClickUrl = v.Url
		money, _ := decimal.NewFromString(v.ZkFinalPrice)
		rate, _ := decimal.NewFromString(v.CommissionRate)
		data[k].Commission = self.getCommission(money, rate, amount)
	}
	return respStruct.Response.ResultList.MapData, nil
}

//物料精选
func (self *TbkService) Lists(params map[string]string) ([]MaterialItem, error) {
	resp, err := self.Request("taobao.tbk.dg.optimus.material", params)
	if err != nil {
		return nil, err
	}
	respStruct := struct {
		Response struct {
			Total_results int `json:"total_results"`
			ResultList    struct {
				MapData []MaterialItem `json:"map_data"`
			} `json:"result_list"`
		} `json:"tbk_dg_optimus_material_response"`
	}{}
	err = json.Unmarshal(resp, &respStruct)
	if err != nil {
		return nil, err
	}
	data := respStruct.Response.ResultList.MapData
	for k, v := range data {
		amount := self.parseAmount(v.CouponAmount)
		money, _ := decimal.NewFromString(v.ZkFinalPrice)
		rate, _ := decimal.NewFromString(v.CommissionRate)
		rate = rate.Mul(decimal.New(100, 0))
		data[k].Commission = self.getCommission(money, rate, amount)
	}
	return respStruct.Response.ResultList.MapData, nil
}

//解析优惠券金额
func (self *TbkService) parseAmount(raw json.RawMessage) (amount decimal.Decimal) {
	var rawAmount interface{}
	json.Unmarshal(raw, &rawAmount)
	switch rawAmount.(type) {
	case int:
		amount, _ = decimal.NewFromString(strconv.Itoa(rawAmount.(int)))
	case string:
		amount, _ = decimal.NewFromString(rawAmount.(string))
	case float64:
		amount = decimal.NewFromFloat(rawAmount.(float64))
	}
	return amount
}

//生成淘口令
func (self *TbkService) CreateTkl(params map[string]string) (string,error) {
	resp, err := self.Request("taobao.tbk.tpwd.create", params)
	if err != nil {
		return "", err
	}
	respStruct := struct {
		TbkPwdResponse struct{
			Data struct{
				Model string `json:"model"`
			} `json:"data"`
		} `json:"tbk_tpwd_create_response"`
	}{}
	err = json.Unmarshal(resp,&respStruct)
	if err != nil {
		return "",err
	}
	return respStruct.TbkPwdResponse.Data.Model,nil
}
//解密淘口令
func (self *TbkService) DecryptTkl(pwd string) (string,error) {
	client := &http.Client{}
	buffer := bytes.NewBuffer([]byte("text="+pwd))
	reqest, err := http.NewRequest("POST", tklDecryptUrl, buffer)
	if err != nil {
		return "",nil
	}
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqest.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(reqest)
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",nil
	}
	respStruct := make(map[string]interface{})
	err = json.Unmarshal(content,&respStruct)
	if err != nil {
		return  "", nil
	}
	data := respStruct["data"].(map[string]interface{})
	goodsUrl := data["url"].(string)
	if !strings.HasPrefix(goodsUrl,"http") {
		goodsUrl = "https:"+goodsUrl
	}
	u,err := url.Parse(goodsUrl)
	if err != nil {
		return "",nil
	}
	return u.Scheme + "://"+u.Host+u.Path,nil
}
//计算佣金
func (self *TbkService) getCommission(money, rate, amount decimal.Decimal) float64 {
	percent, _ := decimal.NewFromString("1000000") // 比例是 100 * 100 * 90/100
	fee, _ := decimal.NewFromString("90")          // 扣手续费10%
	r := money.Sub(amount).Mul(rate).Mul(fee).DivRound(percent, 2)
	s, _ := r.Float64()
	return s
}