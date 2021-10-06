package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"項目/chatroom/common/message"
	"項目/chatroom/server/utils"
)

type UserProcess struct {
	// 暫時不需要字段
}

func (this *UserProcess) Register(userId int,
	userPwd string, userName string) (err error) {
	// 1. 鏈接到服務器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 延時關閉
	defer conn.Close()

	// 2. 準備通過conn發送消息給服務器
	var mes message.Message
	mes.Type = message.RegisterMesType
	// 3. 創建一個LoginMes 結構體
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 4. 將registerMes 序列化
	data, err := json.Marshal(registerMes) // data返回為切片
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 5. 把data賦給 mes.Data字段
	mes.Data = string(data)

	// 6. 將 mes進行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 創建一個Transfer 實例
	tf := &utils.Transfer{
		Conn: conn,
	}

	// 發送data給服務器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("註冊發送信息錯誤 err=", err)
	}

	mes, err = tf.ReadPkg() // mes 就是 RegisterResMes

	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}

	// 將mes的Data部分反序列化成 RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("註冊成功, 你重新登錄一把")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}

	return
}

// 給關聯一個用戶登錄的方法
// 寫一個函數, 完成登錄
func (this *UserProcess) Login(userId int, userPwd string) (err error) {

	// 下一步就要開始定協議..
	// fmt.Printf("userId=%d userPwd=%s\n", userId, userPwd)

	// return nil

	// 1. 鏈接到服務器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 延時關閉
	defer conn.Close()

	// 2. 準備通過conn發送消息給服務器
	var mes message.Message
	mes.Type = message.LoginMesType
	// 3. 創建一個LoginMes 結構體
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	// 4. 將loginMes 序列化
	data, err := json.Marshal(loginMes) // data返回為切片
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 5. 把data賦給 mes.Data字段
	mes.Data = string(data)

	// 6. 將 mes進行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 7. 到這個時候 data就是我們要發送的消息
	// 7.1 先把 data 的長度發送給服務器
	// 先獲取到 data的長度->轉成一個表示長度的byte切片
	var pkgLen uint32 = uint32(len(data)) // 也可 64 或 16
	// pkgLen = uint32(len(data))
	var buf [4]byte
	// buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	// 發送長度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail ", err)
		return
	}

	// fmt.Printf("客戶端, 發送消息的長度=%d 內容=%s", len(data), string(data))

	// 發送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail=", err)
		return
	}

	// 休眠20
	// time.Sleep(20 * time.Second)
	// fmt.Println("休眠了20..")
	// 這裡還需要處理服務器端返回的消息.
	// 創建一個Transfer 實例 (這裡也可以做一個優化用成全局變量)
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}
	// 將mes的Data部分反序列化成 LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		// fmt.Println("登錄成功")

		// 可以顯示當前在線用戶列表, 遍歷loginResMes.UsersId
		fmt.Println("當前在線用戶列表如下:")
		for _, v := range loginResMes.UsersId {

			// 如果我們要求不顯示自己在線, 下面我們增加一個代碼
			if v == userId {
				continue
			}

			fmt.Println("用戶id:\t", v)
			// 完成 客戶端的 onlineUsers 完成初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")

		// 這裡我們還需要在客戶端啟動一個協程
		// 該協程保持和服務器端的通訊.如果服務器有數據推送給客戶端
		// 則接收並顯示在客戶端的終端
		go serverProcessMes(conn)

		// 1. 顯示我們登錄成功的菜單[循環]...
		for {
			ShowMenu()
		}

	} else {
		fmt.Println(loginResMes.Error)
	}

	return
}
