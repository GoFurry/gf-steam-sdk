### 0 internal
| API接口                                        | 参数                                                                              | 描述               |
|----------------------------------------------|---------------------------------------------------------------------------------|------------------|
| internal/api/GetRawBytes                     | `c *client.Client`<br/>`method string`<br/>`url string`<br/>`params url.Values` | HTTP请求原始字节数据     |
| internal/api/GetRawModel                     | `c *client.Client`<br/>`method string`<br/>`url string`<br/>`params url.Values` | HTTP请求原始字节数据返回模型 |
| internal/client/NewClient                    | `cfg *config.SteamConfig`                                                       | 创建工具包的客户端        |
| internal/client/DoRequest                    | `method string`<br/>`baseURL string`<br/>`params url.Values`                    | 通用请求             |
| IFamilyGroupsService/GetSharedLibraryApps/v1 | `access token`                                                                  | 获取家庭组共享的游戏       |


### 1 develop 
api.steampowered.com

| API接口                                             | 封装接口                                 | 强制参数            | 描述                         |
|---------------------------------------------------|--------------------------------------|-----------------|----------------------------|
| IAccountCartService/GetCart/v1                    | sdk.Develop.GetUserCart              | `access token`  | 获取购物车的数据                   |
| IAccountCartService/DeleteCart/v1                 | sdk.Develop.DeleteUserCart           | `access token`  | 清空购物车的数据                   |
| IBillingService/GetRecurringSubscriptionsCount/v1 | sdk.Develop.GetSubscriptionBillCount | `access token`  | 获取 access token 拥有者的订阅账单数量 |
| ICommunityService/GetApps/v1                      | sdk.Develop.GetApps                  | `access token`  | 获取入参对应的商品的简略信息             |
| IFamilyGroupsService/GetChangeLog/v1              | sdk.Develop.GetFamilyChangeLog       | `access token`  | 返回家庭组变更日志                  |
| IFamilyGroupsService/GetFamilyGroup/v1            | sdk.Develop.GetFamilyMembers         | `access token`  | 返回家庭组信息                    |
| IFamilyGroupsService/GetFamilyGroupForUser/v1     | sdk.Develop.GetFamilyGroup           | `access token`  | 返回当前access token用户的家庭组详细信息 |
| IFamilyGroupsService/GetPlaytimeSummary/v1        | sdk.Develop.GetFamilyPlaytime        | `access token`  | 获取家庭组游玩记录信息                |
| IFamilyGroupsService/GetSharedLibraryApps/v1      | sdk.Develop.GetSharedApps            | `access token`  | 获取家庭组共享的游戏                 |

#### 1.1 IAccountCartService
1.1.1 GetCart/v1 <br/>
Get user's cart items <br/>
获取购物车的数据 <br/>
Required: `access token`
```go
cart, err := sdk.Develop.GetUserCart("en", nil)
```
1.1.2 DeleteCart/v1 <br/>
Remove all items from user's cart <br/>
清空购物车的数据 <br/>
Required: `access token`
```go
sdk.Develop.DeleteUserCart(nil)
```
#### 1.2 IBillingService
1.2.1 GetRecurringSubscriptionsCount/v1 <br/>
Get bill count from the access_token's owner <br/>
获取 access token 拥有者的订阅账单数量 <br/>
Required: `access token`
```go
count, err := sdk.Develop.GetSubscriptionBillCount(nil)
```
#### 1.3 ICommunityService
1.3.1 GetApps/v1 <br/>
Get app brief information <br/> 
获取入参对应的商品的简略信息 <br/>
Required: `access token`
```go
apps, err := sdk.Develop.GetApps([]string{"993090", "550"})
```
#### 1.4 IFamilyGroupsService
1.4.1 GetChangeLog/v1 <br/>
Get family change log <br/>
获取家庭组变更日志 <br/>
Required: `access token`
```go
sdk.Develop.GetFamilyChangeLog("1136785")
```
1.4.2 GetFamilyGroup/v1 <br/>
Get family info <br/>
获取家庭组信息 <br/>
Required: `access token`
```go
sdk.Develop.GetFamilyMembers("1136785")
```
1.4.3 GetFamilyGroupForUser/v1 <br/>
Get family group info by user <br/>
获取当前access token用户的家庭组详细信息 <br/>
Required: `access token`
```go
sdk.Develop.GetFamilyGroup("1136785", false)
```
1.4.4 GetPlaytimeSummary/v1 <br/>
Get family playtime <br/>
获取家庭组游玩记录信息 <br/>
Required: `access token`
```go
sdk.Develop.GetFamilyPlaytime("1136785")
```
1.4.5 GetSharedLibraryApps/v1 <br/>
Get family shared apps <br/>
获取家庭组共享的游戏 <br/>
Required: `access token`
```go
sdk.Develop.GetSharedApps("1136785")
```

---

### 2 store
store.steampowered.com

| API接口                                             | 封装接口                                 | 强制参数            | 描述                         |
|---------------------------------------------------|--------------------------------------|-----------------|----------------------------|
|                     |               |   |                    |

#### 2.1

---

### 3 crawler

| 爬取地址                                                                    | 封装接口                                | 描述                 |
|-------------------------------------------------------------------------|-------------------------------------|--------------------|
| any                                                                     | sdk.Crawler.GetGameStoreRawHTML     | 爬取任意地址的原始 HTML     |
| any                                                                     | sdk.Crawler.SaveGameStoreRawHTML    | 爬取并保存任意地址的原始 HTML  |
| `https://store.steampowered.com/app/{$appid}`                           | sdk.Crawler.GetGameStoreRawHTML     | 获取游戏详情页原始 HTML     |
| `https://store.steampowered.com/app/{$appid}`                           | sdk.Crawler.SaveGameStoreRawHTML    | 保存游戏详情页原始 HTML     |
| `https://store.steampowered.com/`                                       | sdk.Crawler.GetHomePageRawHTML      | 获取 Steam 首页原始 HTML |
| `https://store.steampowered.com/`                                       | sdk.Crawler.SaveHomePageRawHTML     | 保存 Steam 首页原始 HTML |
| `https://store.steampowered.com/{$appid}/reviews/`                      | sdk.Crawler.GetGameReviewRawHTML    | 获取游戏评论页原始 HTML     |
| `https://store.steampowered.com/{$appid}/reviews/`                      | sdk.Crawler.SaveGameReviewRawHTML   | 保存游戏评论页原始 HTML     |
| `https://store.steampowered.com/explore/upcoming`                       | sdk.Crawler.GetUpcomingPageRawHTML  | 获取即将推出推荐页原始 HTML   |
| `https://store.steampowered.com/explore/upcoming`                       | sdk.Crawler.SaveUpcomingPageRawHTML | 保存即将推出推荐页原始 HTML   |
| `https://store.steampowered.com/explore/new`                            | sdk.Crawler.GetNewsRawHTML          | 获取新闻推荐页原始 HTML     |
| `https://store.steampowered.com/explore/new`                            | sdk.Crawler.SaveNewsRawHTML         | 保存新闻推荐页原始 HTML     |
| `https://store.steampowered.com/news/?emclan={$emclan}&emgid={$emgid}`  | sdk.Crawler.GetNewsPageRawHTML      | 获取新闻页原始 HTML       |
| `https://store.steampowered.com/news/?emclan={$emclan}&emgid={$emgid}`  | sdk.Crawler.SaveNewsPageRawHTML     | 保存新闻页原始 HTML       |

#### 3.1 Generic
3.1.1 GetRawHTML <br/>
Crawl any address to get HTML <br/>
爬取任意地址的原始 HTML (通用爬取, 无任何跳过验证的策略) <br/>
```go
htmlBytes, err := sdk.Crawler.GetRawHTML(url)
```
3.1.2 SaveRawHTML <br/>
Crawl any address to save HTML <br/>
爬取并保存任意地址的原始 HTML (通用爬取, 无任何跳过验证的策略) <br/>
```go
savePath, err := sdk.Crawler.SaveRawHTML(url, path)
```
#### 3.2 Game Page
3.2.1 GetGameStoreRawHTML <br/>
Get game page raw HTML <br/>
获取游戏详情页原始 HTML <br/>
```go
htmlBytes, err := sdk.Crawler.GetGameStoreRawHTML(appID)
```
3.2.2 SaveGameStoreRawHTML <br/>
Save game page raw HTML <br/>
保存游戏详情页原始 HTML <br/>
```go
savePath, err := sdk.Crawler.SaveGameStoreRawHTML(appID, path)
```
3.2.3 GetHomePageRawHTML <br/>
Get home page raw HTML  <br/>
获取 Steam 首页原始 HTML <br/>
```go
htmlBytes, err := sdk.Crawler.GetHomePageRawHTML()
```
3.2.4 SaveHomePageRawHTML <br/>
Save home page raw HTML  <br/>
保存 Steam 首页原始 HTML <br/>
```go
savePath, err := sdk.Crawler.SaveHomePageRawHTML(path)
```
3.2.5 GetGameReviewRawHTML <br/>
Get app review page raw HTML<br/>
获取游戏评论页原始 HTML <br/>
```go
htmlBytes, err := sdk.Crawler.GetGameReviewRawHTML(appID)
```
3.2.6 SaveGameReviewRawHTML <br/>
Save home page raw HTML <br/>
保存游戏详情页原始 HTML <br/>
```go
savePath, err := sdk.Crawler.SaveGameReviewRawHTML(appID, path)
```
3.2.7 GetUpcomingPageRawHTML <br/>
Get app upcoming page raw HTML <br/>
获取即将推出推荐页原始 HTML <br/>
```go
htmlBytes, err := sdk.Crawler.GetUpcomingPageRawHTML()
```
3.2.8 SaveUpcomingPageRawHTML <br/>
Save app upcoming page raw HTML <br/>
保存即将推出推荐页原始 HTML <br/>
```go
savePath, err := sdk.Crawler.GetUpcomingPageRawHTML(path)
```
3.2.9 GetNewsRawHTML <br/>
Get app news page raw HTML <br/>
获取新闻推荐页原始 HTML <br/>
```go
htmlBytes, err := sdk.Crawler.GetNewsRawHTML()
```
3.2.10 SaveNewsRawHTML <br/>
Save app news page raw HTML <br/>
保存新闻推荐页原始 HTML <br/>
```go
savePath, err := sdk.Crawler.SaveNewsRawHTML(path)
```
3.2.1 GetNewsPageRawHTML <br/>
Get app news page raw HTML <br/>
获取新闻页原始 HTML <br/>
```go
htmlBytes, err := sdk.Crawler.GetNewsPageRawHTML(emclan, emgid)
```
3.2.2 SaveNewsPageRawHTML <br/>
Save app news page raw HTML <br/>
保存新闻推荐页原始 HTML <br/>
```go
savePath, err := sdk.Crawler.SaveNewsPageRawHTML(emclan, emgid, path)
```


---

### 4 Server

| 封装接口                              | 描述                |
|-----------------------------------|-------------------|
| sdk.Server.QueryServerInfo        | 查询单个服务器的基础信息      |
| sdk.Server.QueryServerPlayers     | 查询单个服务器的玩家信息      |
| sdk.Server.QueryServerRules       | 查询单个服务器的规则信息      |
| sdk.Server.GetServerDetail        | 聚合获取单个服务器的完整信息    |
| sdk.Server.QueryServerInfoList    | 批量查询多个服务器的基础信息    |
| sdk.Server.QueryServerPlayersList | 批量查询多个服务器的玩家信息    |
| sdk.Server.QueryServerRulesList   | 批量查询多个服务器的规则信息    |
| sdk.Server.GetServerDetailList    | 批量查询多个服务器的完整聚合信息  |

#### 4.1 Single
4.1.1 QueryServerInfo <br/>
Queries the basic information of a single server <br/>
查询单个服务器的基础信息 <br/>
```go
info, err := sdk.Server.QueryServerInfo(addr)
```
4.1.2 QueryServerPlayers <br/>
Queries the player information of a single server <br/>
查询单个服务器的玩家信息 <br/>
```go
player, err := sdk.Server.QueryServerPlayers(addr)
```
4.1.3 QueryServerRules <br/>
Queries the rule information of a single server <br/>
查询单个服务器的规则信息 <br/>
```go
info, err := sdk.Server.QueryServerRules(addr)
```
4.1.4 QueryServerInfoList <br/>
QueryServerInfoList batch queries the basic information of multiple servers (with rate limit, retry, timeout). <br/>
Uses concurrent query method to ensure the result order is consistent with the input address list, supports exponential backoff retry strategy <br/>
批量查询多个服务器的基础信息(带限流、重试、超时) <br/>
采用并发方式查询, 保证结果顺序与输入地址列表一致, 支持指数退避重试策略 <br/>
```go
infoList, infoErrs, err := sdk.Server.QueryServerInfoList(addrs, TestQPS, TestBurst, TestTimeout, TestRetry)
```
4.1.5 QueryServerPlayersList <br/>
QueryServerPlayersList batch queries the players information of multiple servers (with rate limit, retry, timeout). <br/>
Uses concurrent query method to ensure the result order is consistent with the input address list, supports exponential backoff retry strategy <br/>
批量查询多个服务器的玩家信息(带限流、重试、超时) <br/>
采用并发方式查询, 保证结果顺序与输入地址列表一致, 支持指数退避重试策略 <br/>
```go
infoList, infoErrs, err := sdk.Server.QueryServerPlayersList(addrs, TestQPS, TestBurst, TestTimeout, TestRetry)
```
4.1.6 QueryServerRulesList <br/>
QueryServerRulesList batch queries the rules information of multiple servers (with rate limit, retry, timeout). <br/>
Uses concurrent query method to ensure the result order is consistent with the input address list, supports exponential backoff retry strategy <br/>
批量查询多个服务器的规则信息(带限流、重试、超时) <br/>
采用并发方式查询, 保证结果顺序与输入地址列表一致, 支持指数退避重试策略 <br/>
```go
infoList, infoErrs, err := sdk.Server.QueryServerRulesList(addrs, TestQPS, TestBurst, TestTimeout, TestRetry)
```
#### 4.2 Aggregation
4.2.1 GetServerDetail <br/>
It aggregately gets the complete information of a single server (basic info + players + rules) <br/>
聚合获取单个服务器的完整信息 <br/>
```go
detail, err := sdk.Server.GetServerDetail(addr)
```
4.2.2 GetServerDetailList <br/>
GetServerDetailList batch queries the complete aggregated information of multiple servers (with rate limit, retry, timeout) <br/>
Internally calls the GetServerDetail aggregation interface, uses concurrent query method to ensure the result order is consistent with the input address list <br/>
Supports exponential backoff retry strategy and thread-safe result writing mechanism <br/>
GetServerDetailList 批量查询多个服务器的完整聚合信息(带限流、重试、超时) <br/>
内部调用GetServerDetail聚合接口, 采用并发方式查询, 保证结果顺序与输入地址列表一致 <br/>
支持指数退避重试策略, 线程安全的结果写入机制 <br/>
```go
detailList, detailErrs, err := sdk.Server.GetServerDetailList(addrs, TestQPS, TestBurst, TestTimeout, TestRetry)
```

---

### 5 Util

| 封装接口                       | 描述                        |
|----------------------------|---------------------------|
| sdk.Util.GetStoreToken     | 打开浏览器获取 Steam 商店令牌        |
| sdk.Util.GetCommunityToken | 打开浏览器获取 Steam 社区令牌        |
| sdk.Util.GetAPIKey         | 打开浏览器获取 Steam 开发者 API Key |
| sdk.Util.ParseBBCode       | BBCode 解析为 HTML           |


#### 5.1 Key
5.1.1 GetStoreToken <br/>
Open your browser to get the Steam Store token <br/>
打开浏览器获取 Steam 商店令牌 <br/>
```go
sdk.Util.GetStoreToken()
```
5.1.2 GetCommunityToken <br/>
Open your browser to get the Steam Community token <br/>
打开浏览器获取 Steam 社区令牌 <br/>
```go
sdk.Util.GetCommunityToken()
```
5.1.3 GetAPIKey <br/>
Open your browser to get the Steam developer API Key <br/>
打开浏览器获取 Steam 开发者 API Key <br/>
```go
sdk.Util.GetAPIKey()
```
5.1.4 ParseBBCode <br/>
ParseBBCode recursively parses Steam custom BBCode into HTML string <br/>
将 Steam 自定义 BBCode 递归解析为 HTML 字符串 <br/>
```go
sdk.Util.ParseBBCode(text, limitNumber)
```