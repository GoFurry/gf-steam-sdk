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
	ICommunityService = util.STEAM_API_BASE_URL + "ICommunityService"
)

// ============================ Raw Bytes 原始字节流接口 ============================

// GetAppsRawBytes return game brief info. 返回入参对应游戏的简略信息.
func (s *DevService) GetAppsRawBytes(appids []string) (respBytes []byte, err error) {

	params := url.Values{}
	for idx, appid := range appids {
		params.Set("appids["+util.Int2String(idx)+"]", appid)
	}

	resp, err := s.client.DoRequest("GET", ICommunityService+"/GetApps/v1/", params)
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

// GetAppsRawModel return game brief info. 返回入参对应游戏的简略信息.
func (s *DevService) GetAppsRawModel(appids []string) (models.GetAppsResponse, error) {
	bytes, err := s.GetAppsRawBytes(appids)
	if err != nil {
		return models.GetAppsResponse{}, err
	}

	var appResp models.GetAppsResponse
	if err = sonic.Unmarshal(bytes, &appResp); err != nil {
		return models.GetAppsResponse{}, fmt.Errorf("%w: unmarshal app resp failed: %v", errors.ErrAPIResponse, err)
	}

	return appResp, nil
}

// ============================ Brief Model 精简模型接口 ============================

// GetAppsBrief return game brief info. 返回入参对应游戏的简略信息.
func (s *DevService) GetAppsBrief(appids []string) ([]models.AppBriefInfo, error) {
	rawApp, err := s.GetAppsRawModel(appids)
	if err != nil {
		return nil, err
	}

	apps := make([]models.AppBriefInfo, 0, len(rawApp.Response.Apps))
	for _, a := range rawApp.Response.Apps {
		app := models.AppBriefInfo{
			ID:               a.AppID,
			Name:             a.Name,
			Type:             a.AppType,
			CommunityVisible: a.CommunityVisibleStats,
			Propagation:      a.Propagation,
			Icon:             fmt.Sprintf(util.STEAM_ICON_URL, a.AppID, a.Icon),
		}
		apps = append(apps, app)
	}

	return apps, nil
}

// ============================ Default Interface 默认接口 ============================

// GetApps return game brief info. 返回入参对应游戏的简略信息.
func (s *DevService) GetApps(appids []string) ([]models.AppBriefInfo, error) {
	return s.GetAppsBrief(appids)
}
