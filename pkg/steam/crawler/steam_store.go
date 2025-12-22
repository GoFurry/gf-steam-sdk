package crawler

import (
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
)

// ============================ Raw HTML 获取原始 HTML ============================

// GetHomePageRawHTML get home page raw HTML 获取 Steam 首页原始 HTML
func (s *CrawlerService) GetHomePageRawHTML() ([]byte, error) {
	return s.GetRawHTML(buildStoreURL())
}

// GetGameStoreRawHTML get app page raw HTML 获取游戏详情页原始 HTML
//   - appID: Game AppID
func (s *CrawlerService) GetGameStoreRawHTML(appID uint64) ([]byte, error) {
	if appID == 0 {
		return nil, errors.NewWithType(errors.ErrTypeParam, "appID is empty", nil)
	}
	return s.GetRawHTML(buildStoreURL("app/", util.Uint642String(appID)))
}

// GetGameReviewRawHTML get app review page raw HTML 获取游戏评论页原始 HTML
func (s *CrawlerService) GetGameReviewRawHTML(appID uint64) ([]byte, error) {
	return s.GetRawHTML(buildStoreURL("app/", util.Uint642String(appID), "/reviews/"))
}

// GetUpcomingPageRawHTML get app upcoming page raw HTML 获取即将推出推荐页原始 HTML
func (s *CrawlerService) GetUpcomingPageRawHTML() ([]byte, error) {
	return s.GetRawHTML(buildStoreURL("explore/upcoming"))
}

// GetNewsRawHTML get app news page raw HTML 获取新闻推荐页原始 HTML
func (s *CrawlerService) GetNewsRawHTML() ([]byte, error) {
	return s.GetRawHTML(buildStoreURL("explore/new"))
}

// GetNewsPageRawHTML get app news page raw HTML 获取新闻页原始 HTML
func (s *CrawlerService) GetNewsPageRawHTML(emclan, emgid uint64) ([]byte, error) {
	return s.GetRawHTML(buildStoreURL("news/?emclan=", util.Uint642String(emclan), "&emgid=", util.Uint642String(emgid)))
}

// ============================ Save HTML 保存原始 HTML ============================

// SaveHomePageRawHTML save home page raw HTML 保存 Steam 首页原始 HTML
//   - filename: Custom filename (auto-generate if empty)
func (s *CrawlerService) SaveHomePageRawHTML(filename string) (string, error) {
	return s.SaveRawHTML(buildStoreURL(), filename)
}

// SaveGameStoreRawHTML save app page raw HTML 保存游戏详情页原始 HTML
// eg: game_${appID}.html
//   - appID: Game AppID
//   - filename: Custom filename (auto-generate if empty)
func (s *CrawlerService) SaveGameStoreRawHTML(appID uint64, filename string) (string, error) {
	return s.SaveRawHTML(buildStoreURL("app/", util.Uint642String(appID)), filename)
}

// SaveGameReviewRawHTML save home page raw HTML 保存游戏详情页原始 HTML
//   - filename: Custom filename (auto-generate if empty)
func (s *CrawlerService) SaveGameReviewRawHTML(appID uint64, filename string) (string, error) {
	return s.SaveRawHTML(buildStoreURL("app/", util.Uint642String(appID), "/reviews/"), filename)
}

// SaveUpcomingPageRawHTML save app upcoming page raw HTML 保存即将推出推荐页原始 HTML
//   - filename: Custom filename (auto-generate if empty)
func (s *CrawlerService) SaveUpcomingPageRawHTML(filename string) (string, error) {
	return s.SaveRawHTML(buildStoreURL("explore/upcoming"), filename)
}

// SaveNewsRawHTML save app news page raw HTML 保存新闻推荐页原始 HTML
//   - filename: Custom filename (auto-generate if empty)
func (s *CrawlerService) SaveNewsRawHTML(filename string) (string, error) {
	return s.SaveRawHTML(buildStoreURL("explore/new"), filename)
}

// SaveNewsPageRawHTML save app news page raw HTML 保存新闻推荐页原始 HTML
//   - filename: Custom filename (auto-generate if empty)
func (s *CrawlerService) SaveNewsPageRawHTML(emclan, emgid uint64, filename string) (string, error) {
	return s.SaveRawHTML(buildStoreURL("news/?emclan=", util.Uint642String(emclan),
		"&emgid=", util.Uint642String(emgid)),
		filename,
	)
}

// ============================ Tool 内部工具方法 ============================

// buildStoreURL 构造Steam游戏商店页URL
func buildStoreURL(args ...string) (res string) {
	res = util.STEAM_STORE_BASE_URL
	for _, arg := range args {
		res += arg
	}
	return
}
