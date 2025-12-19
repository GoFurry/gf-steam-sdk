package crawler

import (
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
)

// ============================ 获取原始 HTML ============================

// GetGameStoreRawHTML 获取游戏详情页原始 HTML
// 针对 Steam 商店页优化爬取策略, 自动拼接游戏详情页URL, 返回原始HTML字节流
// 参数:
//   - appID: 游戏 ID | Game AppID
//
// 返回值:
//   - []byte: 原始 HTML 字节流 | Raw HTML bytes
//   - error: 参数/爬取错误 | Parameter/crawling error
func (s *CrawlerService) GetGameStoreRawHTML(appID uint64) ([]byte, error) {
	if appID == 0 {
		return nil, errors.NewWithType(errors.ErrTypeParam, "appID is empty", nil)
	}
	return s.GetRawHTML(buildGameStoreURL(appID))
}

// ============================ 保存原始 HTML ============================

// SaveGameStoreRawHTML 获取游戏详情页原始 HTML 并保存到指定路径
// 自动生成标准化文件名(game_${appID}.html), 适配 Steam 游戏页存储规范
// 参数:
//   - appID: 游戏 ID | Game AppID
//   - filename: 自定义文件名 (为空则自动生成: game_${appID}.html) | Custom filename (auto-generate if empty)
//
// 返回值:
//   - string: 完整存储路径 | Full storage path
//   - error: 参数/爬取/存储错误 | Parameter/crawling/storage error
func (s *CrawlerService) SaveGameStoreRawHTML(appID uint64, filename string) (string, error) {
	return s.SaveRawHTML(buildGameStoreURL(appID), "")
}

// ============================ 内部工具方法 ============================

// buildGameStoreURL 构造Steam游戏商店页URL
// 标准化URL拼接规则，避免手动拼接导致的格式错误
// 参数:
//   - appID: 游戏ID | Game AppID
//
// 返回值:
//   - string: 标准化游戏商店页URL | Standardized game store URL
func buildGameStoreURL(appID uint64) string {
	return "https://store.steampowered.com/app/" + util.Uint642String(appID) + "/"
}
