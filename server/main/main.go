package main

import (
	// "encoding/binary"
	// "encoding/json"
	"fmt"
	"項目/chatroom/server/model"

	// "io"
	"net"
	// "項目/chatroom/common/message"
	"time"
)

// func readPkg(conn net.Conn) (mes message.Message, err error) {

// 	buf := make([]byte, 8096)
// 	fmt.Println("讀取客戶端發送的數據...")
// 	//
// 	// 如果客戶端關閉了 conn 則, 就不會阻塞
// 	_, err = conn.Read(buf[:4])
// 	if err != nil {
// 		// fmt.Println("conn.Read err=", err)
// 		// err = errors.New("read pkg header error")
// 		return
// 	}
// 	fmt.Println("buf[0:4]=", buf[0:4])
// 	// 根據buf[:4] 轉成一個 uint32類型
// 	var pkgLen uint32
// 	pkgLen = binary.BigEndian.Uint32(buf[0:4])

// 	// 根據 pkgLen 讀取消息內容
// 	n, err := conn.Read(buf[:pkgLen]) // 有覆寫過去, 所以別再想為何長度沒有被Unmarshal到了= =
// 	if n != int(pkgLen) || err != nil {
// 		// err = errors.New("read pkg body error")
// 		return
// 	}

// 	// 把pkgLen 反序列化成 -> message.Message

// 	err = json.Unmarshal(buf[:pkgLen], &mes)
// 	if err != nil {
// 		fmt.Println("json.Unmarshal err=", err)
// 		return
// 	}

// 	return

// }

// func writePkg(conn net.Conn, data []byte) (err error) {

// 	// 先發送一個長度給對方
// 	var pkgLen uint32 = uint32(len(data)) // 也可 64 或 16
// 	var buf [4]byte
// 	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
// 	// 發送長度
// 	n, err := conn.Write(buf[:4])
// 	if n != 4 || err != nil {
// 		fmt.Println("conn.Write(bytes) fail ", err)
// 		return
// 	}

// 	// 發送data本身
// 	n, err = conn.Write(data)
// 	if n != int(pkgLen) || err != nil {
// 		fmt.Println("conn.Write(data) err=", err)
// 		return
// 	}
// 	return
// }

// 編寫一個函數serverProcessLogin函數, 專門處理登錄請求
// func ServerProcessLogin(conn net.Conn, mes *message.Message) (err error) {
// 	// 核心代碼
// 	// 1. 先從mes 中取出 mes.Data, 並直接反序列化成LoginMes
// 	var loginMes message.LoginMes
// 	err = json.Unmarshal([]byte(mes.Data), &loginMes)
// 	if err != nil {
// 		fmt.Println("json.Unmarshal fail err=", err)
// 		return
// 	}

// 	// 1.先聲明一個 resMes
// 	var resMes message.Message
// 	resMes.Type = message.LoginResMesType

// 	// 2.再聲明一個 LoginResMes, 並完成賦值
// 	var loginResMes message.LoginResMes

// 	// 如果用戶的id=100, 密碼=123456, 認為合法, 否則不合法
// 	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
// 		// 合法
// 		loginResMes.Code = 200

// 	} else {
// 		// 不合法
// 		loginResMes.Code = 500 // 500 狀態碼, 表示該用戶不存在
// 		loginResMes.Error = "該用戶不存在, 請註冊再使用..."
// 	}

// 	// 3.將 loginResMes 序列化
// 	data, err := json.Marshal(loginResMes)
// 	if err != nil {
// 		fmt.Println("json.Marshal fail err=", err)
// 		return
// 	}

// 	// 4. 將data賦值給resMes
// 	resMes.Data = string(data)

// 	// 5. 對resMes 進行序列化, 準備發送
// 	data, err = json.Marshal(resMes)
// 	if err != nil {
// 		fmt.Println("json.Marshal fail err=", err)
// 		return
// 	}
// 	// 6. 發送data, 我們將其封裝到writePkg函數
// 	err = writePkg(conn, data)
// 	return err
// }

// 編寫一個ServerProcessMes函數
// // 功能: 根據客戶端發送消息種類不同, 決定調用哪個函數來處理
// func ServerProcessMes(conn net.Conn, mes *message.Message) (err error) {

// 	switch mes.Type {
// 	case message.LoginMesType:
// 		// 處理登陸
// 		err = ServerProcessLogin(conn, mes)
// 	case message.LoginResMesType:
// 		// 處理註冊
// 	default:
// 		fmt.Println("消息類型不存在, 無法處理...")
// 	}
// 	return
// }

// 處理和客戶端的通訊
func process(conn net.Conn) {
	// 這裡需要延時關閉conn
	defer conn.Close()

	// 這裡調用總控, 創建一個總控
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客戶端和服務器通訊協程錯誤err=", err)
		return
	}

}

func init() {
	// 當服務器啟動時, 我們就去初始化我們的redis的連接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
}

// 這裡我們編寫一個函數, 完成對UserDao的初始化任務
func initUserDao() {
	// 這裡的pool 本身就是一個全局的變量
	// 這裡需要注意一個初始化順序問題
	// initPool, 再 initUserDao
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {

	// 提示信息
	fmt.Println("服務器[新的結構]在8889端口監聽....")
	listen, err := net.Listen("tcp", "127.0.0.1:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	// 一但監聽成功, 就等待客戶端來鏈接服務器
	for {
		fmt.Println("等待客戶端來鏈接服務器.....")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}

		// 一但鏈接成功, 則啟動一個協程和客戶端保持通訊
		go process(conn)
	}

}
