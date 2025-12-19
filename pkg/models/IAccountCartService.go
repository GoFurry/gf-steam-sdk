package models

// SteamUserCartResponse Steam 原始响应结构体 IPlayerService/GetOwnedGames
type SteamUserCartResponse struct {
	Response struct {
		Cart struct {
			LineItems []struct {
				LineItemID     string `json:"line_item_id"`
				Type           int    `json:"type"`
				PackageID      int    `json:"packageid"`  // 促销包裹ID
				IsValid        bool   `json:"is_valid"`   // 是否有效
				TimeAdded      int64  `json:"time_added"` // 添加时间
				PriceWhenAdded struct {
					AmountInCents   string `json:"amount_in_cents"`  // 价格
					CurrencyCode    int    `json:"currency_code"`    // 货币码
					FormattedAmount string `json:"formatted_amount"` // 格式化金额
				} `json:"price_when_added"` // 添加时的价格
				Flags struct {
					IsGift    bool `json:"is_gift"`    // 是否为礼物
					IsPrivate bool `json:"is_private"` // 是否私有
				} `json:"flags"` // 标识
			} `json:"line_items"` // 购物车行项目列表
			Subtotal struct {
				AmountInCents   string `json:"amount_in_cents"`  // 小计金额
				CurrencyCode    int    `json:"currency_code"`    // 货币编码
				FormattedAmount string `json:"formatted_amount"` // 格式化小计金额
			} `json:"subtotal"` // 购物车小计
			IsValid bool `json:"is_valid"` // 购物车是否有效
		} `json:"cart"` // 购物车信息
	} `json:"response"` // 响应主体
}

// UserCart IPlayerService/GetOwnedGames 精简模型, 购物车信息
type UserCart struct {
	Items          []CartItems `json:"items"`
	TotalPrice     string      `json:"total_price"`
	FormattedTotal string      `json:"formatted_total"`
}

type CartItems struct {
	LineItemID     string `json:"line_item_id"`
	PackageID      int    `json:"package_id"`
	Price          string `json:"price"`
	FormattedPrice string `json:"formatted_price"`
	AddTime        string `json:"add_time"`
	IsGift         bool   `json:"is_gift"`
	IsPrivate      bool   `json:"is_private"`
}
