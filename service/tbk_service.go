package service

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	defaultVersion = "2.0"
	APIHost        = "https://eco.taobao.com/router/rest?"
	timeFormat     = "2006-01-02 15:04:05"
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
	Title          string  `json:"title"`
	PicUrl         string  `json:"pict_url"`
	ZkFinalPrice   string  `json:"zk_final_price"`
	ItemUrl        string  `json:"item_url"`
	CommissionRate string  `json:"commission_rate"`
	Volume         int     `json:"volume"`
	CouponInfo     string  `json:"coupon_info"`
	ShopTitle      string  `json:"shop_title"`
	ShortTitle     string  `json:"short_title"`
	ItemId         int64   `json:"item_id"`
	ClickUrl       string  `json:"click_url"`
	CouponShareUrl string  `json:"coupon_share_url"`
	Commission     float64 `json:"commission"`
}
type MaterialItem struct {
	CommonField
	ClickUrl     string  `json:"click_url"`
	CouponAmount float64 `json:"coupon_amount"`
}

type SearchItem struct {
	CommonField
	CouponAmount string `json:"coupon_amount,number"`
	Url          string `json:"url"`
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
	log.Println(string(buffer.Bytes()))
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
	return respStruct.Response.ResultList.MapData, nil
}

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
	return respStruct.Response.ResultList.MapData, nil
}
