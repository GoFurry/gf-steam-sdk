package dev

import (
	"net/url"

	"github.com/GoFurry/gf-steam-sdk/internal/api"
	"github.com/GoFurry/gf-steam-sdk/internal/client"
	"github.com/GoFurry/gf-steam-sdk/pkg/models"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
)

const (
	IAccountCartService = util.STEAM_API_BASE_URL + "IAccountCartService"
)

// ============================ Raw Bytes 原始字节流接口 ============================

// GetUserCartRawBytes requires access_token, return cart info from the access_token's owner.
func (s *DevService) GetUserCartRawBytes(countryCode string, accessToken *string) (respBytes []byte, err error) {
	return api.GetRawBytes(s.buildUserCart(countryCode, accessToken))
}

// ============================ Structed Raw Model 结构化原始模型接口 ============================

// GetUserCartRawModel requires access_token, return cart info from the access_token's owner.
func (s *DevService) GetUserCartRawModel(countryCode string, accessToken *string) (models.SteamUserCartResponse, error) {
	return api.GetRawModel[models.SteamUserCartResponse](s.buildUserCart(countryCode, accessToken))
}

// ============================ Brief Model 精简模型接口 ============================

// GetUserCartBrief requires access_token, return cart info from the access_token's owner.
func (s *DevService) GetUserCartBrief(countryCode string, accessToken *string) (models.UserCart, error) {
	rawCart, err := s.GetUserCartRawModel(countryCode, accessToken)
	if err != nil {
		return models.UserCart{}, err
	}

	items := make([]models.CartItem, 0, len(rawCart.Response.Cart.LineItems))
	for _, i := range rawCart.Response.Cart.LineItems {
		item := models.CartItem{
			LineItemID:     i.LineItemID,
			PackageID:      i.PackageID,
			Price:          i.PriceWhenAdded.AmountInCents,
			FormattedPrice: i.PriceWhenAdded.FormattedAmount,
			AddTime:        util.TimeUnix2String(i.TimeAdded),
			IsGift:         i.Flags.IsGift,
			IsPrivate:      i.Flags.IsPrivate,
		}
		items = append(items, item)
	}

	return models.UserCart{
		Items:          items,
		TotalPrice:     rawCart.Response.Cart.Subtotal.AmountInCents,
		FormattedTotal: rawCart.Response.Cart.Subtotal.FormattedAmount,
	}, nil
}

// ============================ Default Interface 默认接口 ============================

// GetUserCart requires access_token, return cart info from the access_token's owner.
//   - countryCode changes price
//   - accessToken is required, if globally initialized, use nil
func (s *DevService) GetUserCart(countryCode string, accessToken *string) (models.UserCart, error) {
	return s.GetUserCartBrief(countryCode, accessToken)
}

// DeleteUserCart requires access_token, clear all cart items from the access_token's owner.
//   - accessToken is required, if globally initialized, use nil
func (s *DevService) DeleteUserCart(accessToken *string) error {
	params := url.Values{}
	if accessToken != nil {
		params.Set("access_token", *accessToken)
	}
	_, err := s.client.DoRequest("POST", IAccountCartService+"/DeleteCart/v1/", params)
	return err
}

// ============================ Build 构造入参 ============================

// buildUserCart builds input params.
func (s *DevService) buildUserCart(countryCode string, accessToken *string) (
	c *client.Client,
	method, reqPath string,
	params url.Values,
) {
	params = url.Values{}
	params.Set("user_country", countryCode)
	if accessToken != nil {
		params.Set("access_token", *accessToken)
	}
	return s.client, "GET", IAccountCartService + "/GetCart/v1/", params
}
