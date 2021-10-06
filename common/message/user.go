package message

// 定義一個用戶的結構體

type User struct {
	// 確定字段信息
	// 為了序列化和反序列化成功, 我們必須保證
	// 用戶信息的json字符串的key 和 結構體的字段對應的 tag 名字一致!!!
	UserId     int    `json:"userid"`
	UserPwd    string `json:"userPwd"`
	UserName   string `json:"userName"`
	UserStatus int    `json:"userStatus"` // 用戶狀態
}
