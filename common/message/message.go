package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotyifyUserStatusMes"
)

// 這裡我們定義幾個用戶狀態的常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` // 消息類型
	Data string `json:"data"` // 消息的類型
}

// 定義兩個消息..後面需要再增加
type LoginMes struct {
	UserId   int    `json:"userId"`   // 用戶id
	UserPwd  string `json:"userPwd"`  // 用戶密碼
	UserName string `json:"userName"` // 用戶名
}

type LoginResMes struct {
	Code    int    `json:"code"` // 返回狀態碼 500 表示該用戶為註冊 200 表示登陸成功
	UsersId []int  // 增加字段, 保存用戶Id的切片
	Error   string `json:"error"` // 返回錯誤信息
}

type RegisterMes struct {
	User User `json:"user"` // 類型就是User結構體
}

type RegisterResMes struct {
	Code  int    `json:"code"`  // 返回狀態碼 400 表示該用戶已經占用 200 表示註冊成功
	Error string `json:"error"` // 返回錯誤信息
}

// 為了配合服務器端推送用戶狀態變化的消息
type NotifyUserStatusMes struct {
	UserId int `json:"userid"` // 用戶id
	Status int `json:"status"` // 用戶的狀態
}
