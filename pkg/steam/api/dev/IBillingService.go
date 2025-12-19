package dev

import (
	"fmt"
	"net/url"

	"github.com/GoFurry/gf-steam-sdk/pkg/models"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
	"github.com/bytedance/sonic"
)

const (
	IBillingService = util.STEAM_API_BASE_URL + "IBillingService"
)

// ============================ Raw Bytes 原始字节流接口 ============================

// GetSubscriptionBillCountRawBytes requires access_token, return bill count from the access_token's owner.
// If globally init access_token, you can use access_token = nil.
// 返回 access_token 拥有者的订阅账单数量, 若 access_token 全局初始化则可填 nil.
func (s *DevService) GetSubscriptionBillCountRawBytes(accessToken *string) (respBytes []byte, err error) {

	params := url.Values{}
	if accessToken != nil {
		params.Set("access_token", *accessToken)
	}

	resp, err := s.client.DoRequest("GET", IBillingService+"/GetRecurringSubscriptionsCount/v1/", params)
	if err != nil {
		return respBytes, err
	}

	respBytes, err = sonic.Marshal(resp)
	if err != nil {
		return respBytes, fmt.Errorf("%w: marshal resp failed: %v", errors.ErrAPIResponse, err)
	}

	return respBytes, nil
}

// ============================ Default Interface 默认接口 ============================

// GetSubscriptionBillCount requires access_token, return bill count from the access_token's owner.
// If globally init access_token, you can use access_token = nil.
// 返回 access_token 拥有者的订阅账单数量, 若 access_token 全局初始化则可填 nil.
func (s *DevService) GetSubscriptionBillCount(accessToken *string) (models.SubscriptionBillCountResponse, error) {
	bytes, err := s.GetSubscriptionBillCountRawBytes(accessToken)
	if err != nil {
		return models.SubscriptionBillCountResponse{}, err
	}

	var cartResp models.SubscriptionBillCountResponse
	if err = sonic.Unmarshal(bytes, &cartResp); err != nil {
		return models.SubscriptionBillCountResponse{}, fmt.Errorf("%w: unmarshal bill count resp failed: %v", errors.ErrAPIResponse, err)
	}

	return cartResp, nil
}
