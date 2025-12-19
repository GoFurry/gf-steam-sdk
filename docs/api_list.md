### 1 develop 
api.steampowered.com
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

---

### 2 store
store.steampowered.com
#### 2.1

---

### 3 crawler
#### 3.1

---

### 4 Server
#### 4.1

---

### 5 Util
#### 5.1