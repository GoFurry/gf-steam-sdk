package dev

import (
	"fmt"
	"net/url"

	"github.com/GoFurry/gf-steam-sdk/internal/api"
	"github.com/GoFurry/gf-steam-sdk/internal/client"
	"github.com/GoFurry/gf-steam-sdk/pkg/models"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
)

const (
	// ILoyaltyRewardsService 忠诚度奖励服务接口前缀
	ILoyaltyRewardsService = util.STEAM_API_BASE_URL + "ILoyaltyRewardsService"
)

// ============================ Raw Bytes 原始字节流接口 ============================

// GetEquippedProfileItemsRawBytes 返回已装备个人资料道具的原始字节流
func (s *DevService) GetEquippedProfileItemsRawBytes(steamID string, language *string) (respBytes []byte, err error) {
	return api.GetRawBytes(s.buildEquippedProfileItems(steamID, language))
}

// GetReactionsSummaryForUserRawBytes 返回用户互动汇总的原始字节流
func (s *DevService) GetReactionsSummaryForUserRawBytes(steamID string) (respBytes []byte, err error) {
	return api.GetRawBytes(s.buildReactionsSummaryForUser(steamID))
}

// GetLoyaltyRewardsSummaryRawBytes 返回点数汇总的原始字节流
func (s *DevService) GetLoyaltyRewardsSummaryRawBytes(steamID string) (respBytes []byte, err error) {
	return api.GetRawBytes(s.buildLoyaltyRewardsSummary(steamID))
}

// ============================ Structed Raw Model 结构化原始模型接口 ============================

// GetEquippedProfileItemsRawModel 返回已装备个人资料道具的原始结构化模型
func (s *DevService) GetEquippedProfileItemsRawModel(steamID string, language *string) (models.GetEquippedProfileItemsResponse, error) {
	return api.GetRawModel[models.GetEquippedProfileItemsResponse](s.buildEquippedProfileItems(steamID, language))
}

// GetReactionsSummaryForUserRawModel 返回用户互动汇总的原始结构化模型
func (s *DevService) GetReactionsSummaryForUserRawModel(steamID string) (models.GetReactionsSummaryForUserResponse, error) {
	return api.GetRawModel[models.GetReactionsSummaryForUserResponse](s.buildReactionsSummaryForUser(steamID))
}

// GetLoyaltyRewardsSummaryRawModel 返回点数汇总的原始结构化模型
func (s *DevService) GetLoyaltyRewardsSummaryRawModel(steamID string) (models.GetLoyaltyRewardsSummaryResponse, error) {
	return api.GetRawModel[models.GetLoyaltyRewardsSummaryResponse](s.buildLoyaltyRewardsSummary(steamID))
}

// ============================ Brief Model 精简模型接口 ============================

// GetEquippedProfileItemsBrief 返回已装备个人资料道具的精简信息列表
func (s *DevService) GetEquippedProfileItemsBrief(steamID string, language *string) ([]models.ProfileItemBriefInfo, error) {
	rawResp, err := s.GetEquippedProfileItemsRawModel(steamID, language)
	if err != nil {
		return nil, err
	}

	var items []models.ProfileItemBriefInfo
	items = append(items, s.convertToBriefItems(rawResp.Response.ActiveDefinitions, true)...)
	items = append(items, s.convertToBriefItems(rawResp.Response.InactiveDefinitions, false)...)

	return items, nil
}

// GetReactionsSummaryForUserBrief 返回用户互动汇总的精简信息
func (s *DevService) GetReactionsSummaryForUserBrief(steamID string) (models.UserReactionsTotalBrief, error) {
	rawResp, err := s.GetReactionsSummaryForUserRawModel(steamID)
	if err != nil {
		return models.UserReactionsTotalBrief{}, err
	}

	// 构造总览精简模型
	brief := models.UserReactionsTotalBrief{
		TotalGiven:          rawResp.Response.TotalGiven,
		TotalReceived:       rawResp.Response.TotalReceived,
		TotalPointsGiven:    rawResp.Response.TotalPointsGiven,
		TotalPointsReceived: rawResp.Response.TotalPointsReceived,
		Total:               s.convertToReactionBriefItems(rawResp.Response.Total),
		UserReviews:         s.convertToReactionBriefItems(rawResp.Response.UserReviews),
		UGC:                 s.convertToReactionBriefItems(rawResp.Response.UGC),
		Profile:             s.convertToReactionBriefItems(rawResp.Response.Profile),
	}

	return brief, nil
}

// GetLoyaltyRewardsSummaryBrief 返回点数汇总的精简信息
func (s *DevService) GetLoyaltyRewardsSummaryBrief(steamID string) (models.LoyaltyRewardsSummaryBriefInfo, error) {
	rawResp, err := s.GetLoyaltyRewardsSummaryRawModel(steamID)
	if err != nil {
		return models.LoyaltyRewardsSummaryBriefInfo{}, err
	}

	brief := models.LoyaltyRewardsSummaryBriefInfo{
		CurrentPoints: rawResp.Response.Summary.Points,
		EarnedPoints:  rawResp.Response.Summary.PointsEarned,
		SpentPoints:   rawResp.Response.Summary.PointsSpent,
		UpdatedAt:     rawResp.Response.TimestampUpdated,
	}

	return brief, nil
}

// ============================ Default Interface 默认接口 ============================

// GetEquippedProfileItems 返回已装备个人资料道具的精简信息
func (s *DevService) GetEquippedProfileItems(steamID string, language *string) ([]models.ProfileItemBriefInfo, error) {
	return s.GetEquippedProfileItemsBrief(steamID, language)
}

// GetReactionsSummaryForUser 返回用户互动汇总的精简信息
func (s *DevService) GetReactionsSummaryForUser(steamID string) (models.UserReactionsTotalBrief, error) {
	return s.GetReactionsSummaryForUserBrief(steamID)
}

// GetLoyaltyRewardsSummary 返回点数汇总的精简信息
func (s *DevService) GetLoyaltyRewardsSummary(steamID string) (models.LoyaltyRewardsSummaryBriefInfo, error) {
	return s.GetLoyaltyRewardsSummaryBrief(steamID)
}

// ============================ 工具方法 ============================

// convertToBriefItems 转换原始道具定义为精简模型
func (s *DevService) convertToBriefItems(rawItems []models.ProfileItemDefinition, isActive bool) []models.ProfileItemBriefInfo {
	briefItems := make([]models.ProfileItemBriefInfo, 0, len(rawItems))
	for _, item := range rawItems {
		brief := models.ProfileItemBriefInfo{
			ID:            item.DefID,
			AppID:         item.AppID,
			Name:          item.CommunityItemData.ItemName,
			Title:         item.CommunityItemData.ItemTitle,
			Description:   item.CommunityItemData.ItemDescription,
			PointCost:     item.PointCost,
			IsActive:      isActive,
			IsAnimated:    item.CommunityItemData.Animated,
			LargeImageURL: util.STEAM_COMMUNITY_ASSETS_IMAGES_URL + util.Int642String(item.AppID) + "/" + item.CommunityItemData.ItemImageLarge,
		}
		// 补充可选字段
		if item.CommunityItemData.ItemImageSmall != "" {
			brief.SmallImageURL = util.STEAM_COMMUNITY_ASSETS_IMAGES_URL + util.Int642String(item.AppID) + "/" + item.CommunityItemData.ItemImageSmall
		}
		if item.CommunityItemData.ItemMovieWebm != "" {
			brief.WebmMovieURL = util.STEAM_COMMUNITY_ASSETS_IMAGES_URL + util.Int642String(item.AppID) + "/" + item.CommunityItemData.ItemMovieWebm
		}
		if item.CommunityItemData.ItemMovieMp4 != "" {
			brief.Mp4MovieURL = util.STEAM_COMMUNITY_ASSETS_IMAGES_URL + util.Int642String(item.AppID) + "/" + item.CommunityItemData.ItemMovieMp4
		}
		briefItems = append(briefItems, brief)
	}
	return briefItems
}

// convertToReactionBriefItems 转换原始互动汇总项为精简模型
func (s *DevService) convertToReactionBriefItems(rawItems []models.ReactionSummaryItem) []models.ReactionSummaryBriefInfo {
	briefItems := make([]models.ReactionSummaryBriefInfo, 0, len(rawItems))
	for _, item := range rawItems {
		brief := models.ReactionSummaryBriefInfo{
			ReactionID:     item.ReactionID,
			Given:          item.Given,
			Received:       item.Received,
			PointsGiven:    item.PointsGiven,
			PointsReceived: item.PointsReceived,
			IconURL:        fmt.Sprintf("%s%d.png", util.STEAM_LOYALTY_REACTION_ICON_BASE_URL, item.ReactionID),
		}
		briefItems = append(briefItems, brief)
	}
	return briefItems
}

// ============================ Build 构造入参 ============================

// buildEquippedProfileItems 构造GetEquippedProfileItems接口的请求参数
func (s *DevService) buildEquippedProfileItems(steamID string, language *string) (
	c *client.Client,
	method, reqPath string,
	params url.Values,
) {
	params = url.Values{}
	params.Set("steamid", steamID)
	if language != nil {
		params.Set("language", *language)
	}
	return s.client, "GET", ILoyaltyRewardsService + "/GetEquippedProfileItems/v1/", params
}

// buildReactionsSummaryForUser 构造GetReactionsSummaryForUser接口的请求参数
func (s *DevService) buildReactionsSummaryForUser(steamID string) (
	c *client.Client,
	method, reqPath string,
	params url.Values,
) {
	params = url.Values{}
	params.Set("steamid", steamID)
	return s.client, "GET", ILoyaltyRewardsService + "/GetReactionsSummaryForUser/v1/", params
}

// buildLoyaltyRewardsSummary 构造GetSummary接口的请求参数
func (s *DevService) buildLoyaltyRewardsSummary(steamID string) (
	c *client.Client,
	method, reqPath string,
	params url.Values,
) {
	params = url.Values{}
	params.Set("steamid", steamID)
	return s.client, "GET", ILoyaltyRewardsService + "/GetSummary/v1/", params
}
