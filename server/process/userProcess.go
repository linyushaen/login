package processes

import (
	"encoding/json"
	"fmt"
	"net"
	"項目/chatroom/common/message"
	"項目/chatroom/server/model"
	"項目/chatroom/server/utils"
)

type UserProcess struct {
	// 字段?
	Conn net.Conn
	// 增加一個字段, 表示該conn是哪個用戶
	UserId int
}

// 這裡我們編寫通知所有在線用戶的方法
// userId 要通知其它的在線用戶, 我上線
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {

	// 遍歷 onlineUsers, 然後一個一個的發送	NotyifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		// 過濾掉自己
		if id == userId {
			continue
		}
		// 開始通知[單獨的寫一個方法]
		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyMeOnline(userId int) {

	// 組裝我們的NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	// 將notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 將序列化後的notifyUserStatusMes賦值給 mes.Data
	mes.Data = string(data)

	// 對mes再次序列化, 準備發送.
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 發送, 創建我們Transfer實例, 發送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
}

func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {

	// 1. 先從mes 中取出 mes.Data, 並直接反序列化成Registermes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json unmarshal fail err=", err)
		return
	}

	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes
	// 我們需要到redis數據庫去完成註冊.
	// 1.使用model.MyUserDao 到redis去驗證
	err = model.MyUserDao.Register(&registerMes.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "註冊發生未知錯誤..."
		}
	} else {
		registerResMes.Code = 200
	}

	// 3.將 loginResMes 序列化
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}

	// 4. 將data賦值給resMes
	resMes.Data = string(data)

	// 5. 對resMes 進行序列化, 準備發送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	// 6. 發送data, 我們將其封裝到writePkg函數
	// 因為使用分層模式(mvc), 我們先創建一個Transfer 實例, 然後讀取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return err

}

func (this *UserProcess) ServerProcessLogin( /*conn net.Conn, 可以不要了*/ mes *message.Message) (err error) {
	// 核心代碼
	// 1. 先從mes 中取出 mes.Data, 並直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	// 1.先聲明一個 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	// 2.再聲明一個 LoginResMes, 並完成賦值
	var loginResMes message.LoginResMes

	// 我們需要到redis數據庫去完成驗證.
	// 1.使用model.MyUserDao 到redis去驗證
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err != nil {

		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服務器內部錯誤"
		}

	} else {
		loginResMes.Code = 200
		// 這裡, 因為用戶登錄成功, 我們就把該登錄成功的用戶放入到userMgr中
		// 將登錄成的的用戶的userId 賦給 this
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)
		// 通知其它在線用戶, 我上線了
		this.NotifyOthersOnlineUser(loginMes.UserId)
		// 將當前在線用戶的id 放入到loginResMes.UserId
		// 遍歷 userMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println(user, "登錄成功")
	}

	// 如果用戶的id=100, 密碼=123456, 認為合法, 否則不合法
	// if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
	// 	// 合法
	//

	// } else {
	// 	// 不合法
	// 	loginResMes.Code = 500 // 500 狀態碼, 表示該用戶不存在
	// 	loginResMes.Error = "該用戶不存在, 請註冊再使用..."
	// }

	// 3.將 loginResMes 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}

	// 4. 將data賦值給resMes
	resMes.Data = string(data)

	// 5. 對resMes 進行序列化, 準備發送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	// 6. 發送data, 我們將其封裝到writePkg函數
	// 因為使用分層模式(mvc), 我們先創建一個Transfer 實例, 然後讀取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return err
}
