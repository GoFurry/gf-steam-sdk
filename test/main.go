package main

import (
	"log"
	"os"

	steamConfig "github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam"
)

var cnt = 0

func main() {
	cfg := steamConfig.NewDefaultConfig().
		WithAccessToken(readAccessToken()).
		WithAPIKey(readKey()).
		WithProxyURL("http://127.0.0.1:7897")
	sdk, err := steam.NewSteamSDK(cfg)
	if err != nil {
		log.Fatalf("[Main] 创建 Steam SDK 失败: %v", err)
	}
	defer sdk.Close()

	// 1 IAccountCartService
	// 1.1 GetCart v1 need:access_token
	//cart, err := sdk.Develop.GetUserCart("gb", nil)
	//fmt.Println(cart)
	// 1.2 DeleteCart v1 need:access_token
	//fmt.Println(sdk.Develop.DeleteUserCart(nil))

	// 2 IBillingService
	// 2.1 GetRecurringSubscriptionsCount v1 need:access_token
	//count, err := sdk.Develop.GetSubscriptionBillCount(nil)
	//fmt.Println(count)

	// IPlayerService
	// GetOwnedGames v1 need:key/access_token
	//ownedGames, err := sdk.Develop.GetOwnedGames("76561198370695025", true)
	//cnt = 0
	//for _, ownedGame := range ownedGames {
	//	cnt++
	//	fmt.Println(ownedGame)
	//	if cnt > 10 {
	//		break
	//	}
	//}

	// ISteamUser
	// GetPlayerSummaries v2 need:key
	//summaries, err := sdk.Develop.GetPlayerSummaries("76561198370695025")
	//cnt = 0
	//for _, summarie := range summaries {
	//	cnt++
	//	fmt.Println(summarie)
	//	if cnt > 10 {
	//		break
	//	}
	//}

	// ISteamUserStats
	// GetPlayerAchievements v1 need:key
	//achievements, err := sdk.Develop.GetPlayerAchievements("76561198370695025", 550, "zh")
	//cnt = 0
	//for _, achievement := range achievements {
	//	cnt++
	//	fmt.Println(achievement)
	//	if cnt > 10 {
	//		break
	//	}
	//}
}

func readKey() string {
	file, _ := os.ReadFile("./key.txt")
	return string(file)
}

func readAccessToken() string {
	file, _ := os.ReadFile("./token.txt")
	return string(file)
}
