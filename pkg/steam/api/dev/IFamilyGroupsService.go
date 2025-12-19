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
	IFamilyGroupsService = util.STEAM_API_BASE_URL + "IFamilyGroupsService"
)

// ============================ Raw Bytes 原始字节流接口 ============================

// GetFamilyChangeLogRawBytes return family change log. 返回家庭组变更日志.
func (s *DevService) GetFamilyChangeLogRawBytes(familyID string) (respBytes []byte, err error) {

	params := url.Values{}
	params.Set("family_groupid", familyID)

	resp, err := s.client.DoRequest("GET", IFamilyGroupsService+"/GetChangeLog/v1/", params)
	if err != nil {
		return respBytes, err
	}

	respBytes, err = sonic.Marshal(resp)
	if err != nil {
		return respBytes, fmt.Errorf("%w: marshal resp failed: %v", errors.ErrAPIResponse, err)
	}

	return respBytes, nil
}

// GetFamilyMembersRawBytes return family info. 返回家庭组信息.
func (s *DevService) GetFamilyMembersRawBytes(familyID string) (respBytes []byte, err error) {

	params := url.Values{}
	params.Set("family_groupid", familyID)

	resp, err := s.client.DoRequest("GET", IFamilyGroupsService+"/GetFamilyGroup/v1/", params)
	if err != nil {
		return respBytes, err
	}

	respBytes, err = sonic.Marshal(resp)
	if err != nil {
		return respBytes, fmt.Errorf("%w: marshal resp failed: %v", errors.ErrAPIResponse, err)
	}

	return respBytes, nil
}

// GetFamilyGroupRawBytes return family group info by user. 返回当前access token用户的家庭组详细信息.
func (s *DevService) GetFamilyGroupRawBytes(familyID string, included bool) (respBytes []byte, err error) {

	params := url.Values{}
	params.Set("family_groupid", familyID)
	if included {
		params.Set("include_family_group_response", "true")
	}

	resp, err := s.client.DoRequest("GET", IFamilyGroupsService+"/GetFamilyGroupForUser/v1/", params)
	if err != nil {
		return respBytes, err
	}

	respBytes, err = sonic.Marshal(resp)
	if err != nil {
		return respBytes, fmt.Errorf("%w: marshal resp failed: %v", errors.ErrAPIResponse, err)
	}

	return respBytes, nil
}

// GetFamilyPlaytimeRawBytes return family playtime. 返回家庭组游玩记录信息.
func (s *DevService) GetFamilyPlaytimeRawBytes(familyID string) (respBytes []byte, err error) {

	params := url.Values{}
	params.Set("family_groupid", familyID)

	resp, err := s.client.DoRequest("POST", IFamilyGroupsService+"/GetPlaytimeSummary/v1/", params)
	if err != nil {
		return respBytes, err
	}

	respBytes, err = sonic.Marshal(resp)
	if err != nil {
		return respBytes, fmt.Errorf("%w: marshal resp failed: %v", errors.ErrAPIResponse, err)
	}

	return respBytes, nil
}

// GetSharedAppsRawBytes return family shared apps. 返回家庭组共享的游戏.
func (s *DevService) GetSharedAppsRawBytes(familyID string) (respBytes []byte, err error) {

	params := url.Values{}
	params.Set("family_groupid", familyID)

	resp, err := s.client.DoRequest("GET", IFamilyGroupsService+"/GetSharedLibraryApps/v1/", params)
	if err != nil {
		return respBytes, err
	}

	respBytes, err = sonic.Marshal(resp)
	if err != nil {
		return respBytes, fmt.Errorf("%w: marshal resp failed: %v", errors.ErrAPIResponse, err)
	}

	return respBytes, nil
}

// ============================ Structed Raw Model 结构化原始模型接口 ============================

// GetFamilyChangeLogRawModel return family change log. 返回家庭组变更日志.
func (s *DevService) GetFamilyChangeLogRawModel(familyID string) (models.FamilyGroupChangeLogResponse, error) {
	bytes, err := s.GetFamilyChangeLogRawBytes(familyID)
	if err != nil {
		return models.FamilyGroupChangeLogResponse{}, err
	}

	var logResp models.FamilyGroupChangeLogResponse
	if err = sonic.Unmarshal(bytes, &logResp); err != nil {
		return models.FamilyGroupChangeLogResponse{}, fmt.Errorf("%w: unmarshal log resp failed: %v", errors.ErrAPIResponse, err)
	}

	return logResp, nil
}

// GetFamilyMembersRawModel return family info. 返回家庭组信息.
func (s *DevService) GetFamilyMembersRawModel(familyID string) (models.FamilyGroupResponse, error) {
	bytes, err := s.GetFamilyMembersRawBytes(familyID)
	if err != nil {
		return models.FamilyGroupResponse{}, err
	}

	var memberResp models.FamilyGroupResponse
	if err = sonic.Unmarshal(bytes, &memberResp); err != nil {
		return models.FamilyGroupResponse{}, fmt.Errorf("%w: unmarshal member resp failed: %v", errors.ErrAPIResponse, err)
	}

	return memberResp, nil
}

// GetFamilyGroupRawModel return family group info by user. 返回当前access token用户的家庭组详细信息.
func (s *DevService) GetFamilyGroupRawModel(familyID string, included bool) (models.FamilyGroupForUserResponse, error) {
	bytes, err := s.GetFamilyGroupRawBytes(familyID, included)
	if err != nil {
		return models.FamilyGroupForUserResponse{}, err
	}

	var memberResp models.FamilyGroupForUserResponse
	if err = sonic.Unmarshal(bytes, &memberResp); err != nil {
		return models.FamilyGroupForUserResponse{}, fmt.Errorf("%w: unmarshal resp failed: %v", errors.ErrAPIResponse, err)
	}

	return memberResp, nil
}

// GetFamilyPlaytimeRawModel return family playtime. 返回家庭组游玩记录信息.
func (s *DevService) GetFamilyPlaytimeRawModel(familyID string) (models.FamilyGroupPlaytimeSummaryResponse, error) {
	bytes, err := s.GetFamilyPlaytimeRawBytes(familyID)
	if err != nil {
		return models.FamilyGroupPlaytimeSummaryResponse{}, err
	}

	var memberResp models.FamilyGroupPlaytimeSummaryResponse
	if err = sonic.Unmarshal(bytes, &memberResp); err != nil {
		return models.FamilyGroupPlaytimeSummaryResponse{}, fmt.Errorf("%w: unmarshal resp failed: %v", errors.ErrAPIResponse, err)
	}

	return memberResp, nil
}

// GetSharedAppsRawModel return family shared apps. 返回家庭组共享的游戏.
func (s *DevService) GetSharedAppsRawModel(familyID string) (models.FamilySharedLibraryResponse, error) {
	bytes, err := s.GetSharedAppsRawBytes(familyID)
	if err != nil {
		return models.FamilySharedLibraryResponse{}, err
	}

	var memberResp models.FamilySharedLibraryResponse
	if err = sonic.Unmarshal(bytes, &memberResp); err != nil {
		return models.FamilySharedLibraryResponse{}, fmt.Errorf("%w: unmarshal resp failed: %v", errors.ErrAPIResponse, err)
	}

	return memberResp, nil
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
