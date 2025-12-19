package models

type SubscriptionBillCountResponse struct {
	Response struct {
		ActiveSubscriptionsCount   int64 `json:"active_subscriptions_count"`
		InactiveSubscriptionsCount int64 `json:"inactive_subscriptions_count"`
	} `json:"response"`
}
