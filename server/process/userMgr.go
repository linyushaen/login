package processes

import (
	"fmt"
)

// 因為UserMgr 實例在服務器端有且只有一個
// 因為在很多的地方, 都會使用到, 因此, 我們
// 將其定義為全局變量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 完成對userMgr初始化的工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// 完成對onlineUsers添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {

	this.onlineUsers[up.UserId] = up
}

// 刪除
func (this *UserMgr) DelOnlineUser(userId int) {
	delete(this.onlineUsers, userId)
}

// 返回當前所有在線的用戶
func (this *UserMgr) GetAllOnlineUsers() map[int]*UserProcess {
	return this.onlineUsers
}

// 根據id返回對應的值
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {

	// 如何從map取出一個值, 帶檢測方式
	up, ok := this.onlineUsers[userId]
	if !ok { // 說明, 你要查找的這個用戶, 當前不在線
		err = fmt.Errorf("用戶%d 不存在", userId)
		return
	}
	return
}
