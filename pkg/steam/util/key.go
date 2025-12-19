package util

import (
	"fmt"

	tool "github.com/GoFurry/gf-steam-sdk/pkg/util"
)

// GetStoreToken 获取 Steam 商店令牌(Webapi_token)
// 指引开发者通过网页登录 Steam 后, 从指定页面的 JSON 响应中提取 Store 令牌,
// 该令牌有效期为 1 天, 是调用 Steam 商店相关接口的关键认证信息
// 操作步骤：
//  1. 确保 Steam 账号已网页登录
//  2. 自动打开令牌获取页面
//  3. 从 JSON 响应中提取 "Webapi_token" 字段值
//
// GetStoreToken gets the Steam Store token (Webapi_token)
// Guides developers to extract the Store token from the JSON response of the specified page after logging in to Steam via web,
// the token is valid for 1 day and is key authentication information for calling Steam Store-related interfaces
// Operation steps:
//  1. Ensure the Steam account is logged in via web
//  2. Automatically open the token acquisition page
//  3. Extract the value of the "Webapi_token" field from the JSON response
func (s *UtilService) GetStoreToken() {
	fmt.Println("[Info] 需要先网页登录Steam, 通过cookie获取令牌(有效期1天)")
	fmt.Println("[Info] Make sure your account is logged in to Steam via web, get token via cookie (valid for 1 day)")
	tool.OpenBrowser("https://store.steampowered.com/pointssummary/ajaxgetasyncconfig")
	fmt.Println("JSON响应中的 Webapi_token 即为商店令牌")
	fmt.Println("Webapi_token in the JSON response is the Store token you need")
	fmt.Println()
	return
}

// GetCommunityToken 获取 Steam 社区令牌(loyalty_webapi_token)
// 指引开发者通过网页登录 Steam 后, 使用 JS 脚本从社区页面提取社区令牌,
// 该令牌有效期为 1 天, 该令牌是调用 Steam 社区相关接口的关键认证信息
// 操作步骤：
//  1. 确保 Steam 账号已网页登录
//  2. 自动打开社区令牌获取页面
//  3. 按 F12 打开控制台, 执行指定 JS 脚本提取令牌
//
// GetCommunityToken gets the Steam Community token (loyalty_webapi_token)
// Guides developers to extract the Community token from the community page using JS scripts after logging in to Steam via web,
// the token is valid for 1 day and is key authentication information for calling Steam Community-related interfaces
// Operation steps:
//  1. Ensure the Steam account is logged in via web
//  2. Automatically open the Community token acquisition page
//  3. Press F12 to open the console and execute the specified JS script to extract the token
func (s *UtilService) GetCommunityToken() {
	fmt.Println("[Info] 需要先网页登录Steam, 通过cookie和JS脚本获取社区令牌")
	fmt.Println("[Info] Make sure your account is logged in to Steam via web, get Community token via cookie and JS script")
	tool.OpenBrowser("https://steamcommunity.com/my/edit/info")
	fmt.Println("使用以下JS脚本获取社区令牌(F12打开控制台后粘贴)：")
	fmt.Println("Use the following JS script to get Community token (paste after pressing F12 to open console):")
	fmt.Println(`  const token = JSON.parse(application_config.dataset.loyalty_webapi_token);`)
	fmt.Println(`  console.log("Steam community token：", token);`)
	fmt.Println()
	return
}

// GetAPIKey 获取 Steam 开发者 API Key
// 指引开发者通过 Steam 开发者页面申请 API Key,
// 该 Key 是调用 Steam Web API 的核心认证信息, 申请前需确保账号已登录且绑定域名
// 操作步骤：
//  1. 确保 Steam 账号已网页登录
//  2. 自动打开 API Key 申请页面
//  3. 按页面指引完成 API Key 申请和绑定
//
// GetAPIKey gets the Steam developer API Key
// Guides developers to apply for an API Key through the Steam developer page,
// the Key is core authentication information for calling the Steam Web API, and the account must be logged in and bound to a domain name before application
// Operation steps:
//  1. Ensure the Steam account is logged in via web
//  2. Automatically open the API Key application page
//  3. Complete API Key application and binding according to page guidance
func (s *UtilService) GetAPIKey() {
	fmt.Println("[Info] 需要先网页登录Steam账号, 前往开发者页面申请API Key")
	fmt.Println("[Info] Make sure your account is logged in to Steam via web, go to the developer page to apply for API Key")
	tool.OpenBrowser("https://steamcommunity.com/dev/apikey")
	fmt.Println()
	return
}
