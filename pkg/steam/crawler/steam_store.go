package crawler

import (
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
)

// ============================ Raw HTML 获取原始 HTML ============================

// GetGameStoreRawHTML get game page raw HTML 获取游戏详情页原始 HTML
//   - appID: Game AppID
func (s *CrawlerService) GetGameStoreRawHTML(appID uint64) ([]byte, error) {
	if appID == 0 {
		return nil, errors.NewWithType(errors.ErrTypeParam, "appID is empty", nil)
	}
	return s.GetRawHTML(buildGameStoreURL(appID))
}

// ============================ Save HTML 保存原始 HTML ============================

// SaveGameStoreRawHTML save game page raw HTML 保存游戏详情页原始 HTML
// eg: game_${appID}.html
//   - appID: Game AppID
//   - filename: Custom filename (auto-generate if empty)
func (s *CrawlerService) SaveGameStoreRawHTML(appID uint64, filename string) (string, error) {
	return s.SaveRawHTML(buildGameStoreURL(appID), filename)
}

// ============================ Tool 内部工具方法 ============================

// buildGameStoreURL 构造Steam游戏商店页URL
// 标准化URL拼接规则, 避免手动拼接导致的格式错误
//   - appID: 游戏ID | Game AppID
//
// return:
//   - string: Standardized game store URL
func buildGameStoreURL(appID uint64) string {
	return "https://store.steampowered.com/app/" + util.Uint642String(appID) + "/"
}
