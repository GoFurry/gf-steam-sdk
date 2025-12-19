package models

// SubscriptionBillCountResponse IBillingService/GetRecurringSubscriptionsCount
type SubscriptionBillCountResponse struct {
	Response struct {
		ActiveSubscriptionsCount   int64 `json:"active_subscriptions_count"`   // 激活的订阅数量
		InactiveSubscriptionsCount int64 `json:"inactive_subscriptions_count"` // 未激活的订阅数量
	} `json:"response"`
}
