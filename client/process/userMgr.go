package process

import (
	_ "fmt"
	"項目/chatroom/common/message"
)

// 客戶端要維護的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)

// 編寫一個方法, 處理返回的NotityUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	// 適當優化
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok { // 原來沒有
		user = &message.User{
			UserId:     notifyUserStatusMes.UserId,
			UserStatus: notifyUserStatusMes.Status,
		}
	}

	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user
}
