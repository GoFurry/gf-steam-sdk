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
	IFamilyGroupsService = util.STEAM_API_BASE_URL + "IFamilyGroupsService"
)

// ============================ Raw Bytes 原始字节流接口 ============================

// GetFamilyChangeLogRawBytes return family change log. 返回家庭组变更日志.
func (s *DevService) GetFamilyChangeLogRawBytes(familyID string) (respBytes []byte, err error) {
	return api.GetRawBytes(s.buildFamilyChangeLog(familyID))
}

// GetFamilyMembersRawBytes return family info. 返回家庭组信息.
func (s *DevService) GetFamilyMembersRawBytes(familyID string) (respBytes []byte, err error) {
	return api.GetRawBytes(s.buildFamilyMembers(familyID))
}

// GetFamilyGroupRawBytes return family group info by user. 返回当前access token用户的家庭组详细信息.
func (s *DevService) GetFamilyGroupRawBytes(familyID string, included bool) (respBytes []byte, err error) {
	return api.GetRawBytes(s.buildFamilyGroup(familyID, included))
}

// GetFamilyPlaytimeRawBytes return family playtime. 返回家庭组游玩记录信息.
func (s *DevService) GetFamilyPlaytimeRawBytes(familyID string) (respBytes []byte, err error) {
	return api.GetRawBytes(s.buildFamilyPlaytime(familyID))
}

// GetSharedAppsRawBytes return family shared apps. 返回家庭组共享的游戏.
func (s *DevService) GetSharedAppsRawBytes(familyID string) (respBytes []byte, err error) {
	return api.GetRawBytes(s.buildSharedApps(familyID))
}

// ============================ Structed Raw Model 结构化原始模型接口 ============================

// GetFamilyChangeLogRawModel return family change log. 返回家庭组变更日志.
func (s *DevService) GetFamilyChangeLogRawModel(familyID string) (models.FamilyGroupChangeLogResponse, error) {
	return api.GetRawModel[models.FamilyGroupChangeLogResponse](s.buildFamilyChangeLog(familyID))
}

// GetFamilyMembersRawModel return family info. 返回家庭组信息.
func (s *DevService) GetFamilyMembersRawModel(familyID string) (models.FamilyGroupResponse, error) {
	return api.GetRawModel[models.FamilyGroupResponse](s.buildFamilyMembers(familyID))
}

// GetFamilyGroupRawModel return family group info by user. 返回当前access token用户的家庭组详细信息.
func (s *DevService) GetFamilyGroupRawModel(familyID string, included bool) (models.FamilyGroupForUserResponse, error) {
	return api.GetRawModel[models.FamilyGroupForUserResponse](s.buildFamilyGroup(familyID, included))
}

// GetFamilyPlaytimeRawModel return family playtime. 返回家庭组游玩记录信息.
func (s *DevService) GetFamilyPlaytimeRawModel(familyID string) (models.FamilyGroupPlaytimeSummaryResponse, error) {
	return api.GetRawModel[models.FamilyGroupPlaytimeSummaryResponse](s.buildFamilyPlaytime(familyID))
}

// GetSharedAppsRawModel return family shared apps. 返回家庭组共享的游戏.
func (s *DevService) GetSharedAppsRawModel(familyID string) (models.FamilySharedLibraryResponse, error) {
	return api.GetRawModel[models.FamilySharedLibraryResponse](s.buildSharedApps(familyID))
}

// ============================ Brief Model 精简模型接口 ============================

// GetFamilyChangeLogBrief return game brief info. 返回入参对应游戏的简略信息.
func (s *DevService) GetFamilyChangeLogBrief(familyID string) ([]models.FamilyGroupChange, error) {
	rawLog, err := s.GetFamilyChangeLogRawModel(familyID)
	if err != nil {
		return nil, err
	}

	return rawLog.Response.Changes, nil
}

// GetFamilyPlaytimeBrief return family playtime. 返回家庭组游玩记录信息.
func (s *DevService) GetFamilyPlaytimeBrief(familyID string) ([]models.FamilyGroupPlaytimeBrief, error) {
	rawPlaytime, err := s.GetFamilyPlaytimeRawModel(familyID)
	if err != nil {
		return nil, err
	}

	// 转换为精简模型 | Convert to simplified model
	playtimes := make([]models.FamilyGroupPlaytimeBrief, 0, len(rawPlaytime.Response.Entries))
	for _, p := range rawPlaytime.Response.Entries {
		playtime := models.FamilyGroupPlaytimeBrief{
			AppID:         p.AppID,
			SteamID:       p.SteamID,
			FirstPlayed:   util.TimeUnix2String(p.FirstPlayed),
			LatestPlayed:  util.TimeUnix2String(p.LatestPlayed),
			SecondsPlayed: p.SecondsPlayed,
		}
		playtimes = append(playtimes, playtime)
	}

	return playtimes, nil
}

// GetSharedAppsBrief return family shared apps. 返回家庭组共享的游戏.
func (s *DevService) GetSharedAppsBrief(familyID string) (models.FamilySharedLibraryAppBrief, error) {
	rawShared, err := s.GetSharedAppsRawModel(familyID)
	if err != nil {
		return models.FamilySharedLibraryAppBrief{}, err
	}

	// 转换为精简模型 | Convert to simplified model
	apps := make([]models.SharedAppBrief, 0, len(rawShared.Response.Apps))
	for _, a := range rawShared.Response.Apps {
		app := models.SharedAppBrief{
			AppID:          a.AppID,
			OwnerSteamIDs:  a.OwnerSteamIDs,
			Name:           a.Name,
			Cover:          fmt.Sprintf(util.STEAM_CAPSULE_URL, a.AppID),
			Icon:           fmt.Sprintf(util.STEAM_ICON_URL, a.AppID, a.ImgIconHash),
			ExcludeReason:  a.ExcludeReason,
			RtTimeAcquired: util.TimeUnix2String(a.RtTimeAcquired),
			RtLastPlayed:   util.TimeUnix2String(a.RtLastPlayed),
			RtPlaytime:     a.RtPlaytime,
			AppType:        a.AppType,
		}
		apps = append(apps, app)
	}

	res := models.FamilySharedLibraryAppBrief{
		OwnerSteamID: rawShared.Response.OwnerSteamID,
		Apps:         apps,
	}

	return res, nil
}

// ============================ Default Interface 默认接口 ============================

// GetFamilyChangeLog return game brief info. 返回入参对应游戏的简略信息.
func (s *DevService) GetFamilyChangeLog(familyID string) ([]models.FamilyGroupChange, error) {
	return s.GetFamilyChangeLogBrief(familyID)
}

// GetFamilyMembers return family info. 返回家庭组信息.
func (s *DevService) GetFamilyMembers(familyID string) (models.FamilyGroup, error) {
	member, err := s.GetFamilyMembersRawModel(familyID)
	if err != nil {
		return models.FamilyGroup{}, err
	}
	return member.Response, nil
}

// GetFamilyGroup return family group info by user. 返回当前access token用户的家庭组详细信息.
func (s *DevService) GetFamilyGroup(familyID string, included bool) (models.FamilyGroupForUserResponse, error) {
	return s.GetFamilyGroupRawModel(familyID, included)
}

// GetFamilyPlaytime return family playtime. 返回家庭组游玩记录信息.
func (s *DevService) GetFamilyPlaytime(familyID string) ([]models.FamilyGroupPlaytimeBrief, error) {
	return s.GetFamilyPlaytimeBrief(familyID)
}

// GetSharedApps return family shared apps. 返回家庭组共享的游戏.
func (s *DevService) GetSharedApps(familyID string) (models.FamilySharedLibraryAppBrief, error) {
	return s.GetSharedAppsBrief(familyID)
}

// ============================ Build 构造入参 ============================

// buildFamilyChangeLog builds input params.
func (s *DevService) buildFamilyChangeLog(familyID string) (
	c *client.Client,
	method, reqPath string,
	params url.Values,
) {
	params = url.Values{}
	params.Set("family_groupid", familyID)
	return s.client, "GET", IFamilyGroupsService + "/GetChangeLog/v1/", params
}

// buildFamilyMembers builds input params.
func (s *DevService) buildFamilyMembers(familyID string) (
	c *client.Client,
	method, reqPath string,
	params url.Values,
) {
	params = url.Values{}
	params.Set("family_groupid", familyID)
	return s.client, "GET", IFamilyGroupsService + "/GetFamilyGroup/v1/", params
}

// buildFamilyGroup builds input params.
func (s *DevService) buildFamilyGroup(familyID string, included bool) (
	c *client.Client,
	method, reqPath string,
	params url.Values,
) {
	params = url.Values{}
	params.Set("family_groupid", familyID)
	if included {
		params.Set("include_family_group_response", "true")
	}
	return s.client, "GET", IFamilyGroupsService + "/GetFamilyGroupForUser/v1/", params
}

// buildFamilyPlaytime builds input params.
func (s *DevService) buildFamilyPlaytime(familyID string) (
	c *client.Client,
	method, reqPath string,
	params url.Values,
) {
	params = url.Values{}
	params.Set("family_groupid", familyID)
	return s.client, "POST", IFamilyGroupsService + "/GetPlaytimeSummary/v1/", params
}

// buildSharedApps builds input params.
func (s *DevService) buildSharedApps(familyID string) (
	c *client.Client,
	method, reqPath string,
	params url.Values,
) {
	params = url.Values{}
	params.Set("family_groupid", familyID)
	return s.client, "GET", IFamilyGroupsService + "/GetSharedLibraryApps/v1/", params
}
