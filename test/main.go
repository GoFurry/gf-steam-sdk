package main

import (
	"fmt"
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
	// 1.1 GetCart v1 required:access_token
	//cart, err := sdk.Develop.GetUserCart("gb", nil)
	//fmt.Println(cart)
	// 1.2 DeleteCart v1 required:access_token
	//fmt.Println(sdk.Develop.DeleteUserCart(nil))

	// 2 IBillingService
	// 2.1 GetRecurringSubscriptionsCount v1 required:access_token
	//count, err := sdk.Develop.GetSubscriptionBillCount(nil)
	//fmt.Println(count)

	// 3 ICommunityService
	// 3.1 GetApps v1
	//apps, err := sdk.Develop.GetApps([]string{"993090", "550"})
	//for _, app := range apps {
	//	fmt.Println(app)
	//}

	// 4 IFamilyGroupsService
	// 4.1 GetChangeLog v1
	//changeLogs, err := sdk.Develop.GetFamilyChangeLog("1136785")
	//cnt = 0
	//for _, clog := range changeLogs {
	//	cnt++
	//	fmt.Println(clog)
	//	if cnt > 10 {
	//		break
	//	}
	//}
	// 4.2 GetFamilyGroup v1
	//family, err := sdk.Develop.GetFamilyMembers("1136785")
	//cnt = 0
	//for _, member := range family.Members {
	//	cnt++
	//	fmt.Println(member)
	//	if cnt > 10 {
	//		break
	//	}
	//}
	// 4.3 GetFamilyGroupForUser v1
	//familyGroup, err := sdk.Develop.GetFamilyGroup("1136785", false)
	//fmt.Println(familyGroup)
	// 4.4 GetPlaytimeSummary v1
	//playtime, err := sdk.Develop.GetFamilyPlaytime("1136785")
	//for _, p := range playtime {
	//	cnt++
	//	fmt.Println(p)
	//	if cnt > 10 {
	//		break
	//	}
	//}
	// 4.5 GetSharedLibraryApps v1
	sharedApps, err := sdk.Develop.GetSharedApps("1136785")
	for _, a := range sharedApps.Apps {
		cnt++
		fmt.Println(a)
		if cnt > 10 {
			break
		}
	}

	// IPlayerService
	// GetOwnedGames v1 required:key/access_token
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
	// GetPlayerSummaries v2 required:key
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
	// GetPlayerAchievements v1 required:key
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
