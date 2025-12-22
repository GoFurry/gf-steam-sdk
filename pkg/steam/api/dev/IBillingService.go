package dev

import (
	"net/url"

	"github.com/GoFurry/gf-steam-sdk/internal/api"
	"github.com/GoFurry/gf-steam-sdk/internal/client"
	"github.com/GoFurry/gf-steam-sdk/pkg/models"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
)

const (
	IBillingService = util.STEAM_API_BASE_URL + "IBillingService"
)

// ============================ Raw Bytes 原始字节流接口 ============================

// GetSubscriptionBillCountRawBytes requires access_token, return bill count from the access_token's owner.
// If globally init access_token, you can use access_token = nil.
// 返回 access_token 拥有者的订阅账单数量, 若 access_token 全局初始化则可填 nil.
func (s *DevService) GetSubscriptionBillCountRawBytes(accessToken *string) (respBytes []byte, err error) {
	return api.GetRawBytes(s.buildSubscriptionBill(accessToken))
}

// ============================ Default Interface 默认接口 ============================

// GetSubscriptionBillCount requires access_token, return bill count from the access_token's owner.
// If globally init access_token, you can use access_token = nil.
// 返回 access_token 拥有者的订阅账单数量, 若 access_token 全局初始化则可填 nil.
func (s *DevService) GetSubscriptionBillCount(accessToken *string) (models.SubscriptionBillCountResponse, error) {
	return api.GetRawModel[models.SubscriptionBillCountResponse](s.buildSubscriptionBill(accessToken))
}

// ============================ Build 构造入参 ============================

// buildSubscriptionBill builds input params.
func (s *DevService) buildSubscriptionBill(accessToken *string) (
	c *client.Client,
	method, reqPath string,
	params url.Values,
) {
	params = url.Values{}
	if accessToken != nil {
		params.Set("access_token", *accessToken)
	}
	return s.client, "GET", IBillingService + "/GetRecurringSubscriptionsCount/v1/", params
}
