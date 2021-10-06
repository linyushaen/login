package main

import (
	"fmt"
	"項目/chatroom/server/model"
	"net"
	"time"
)




// 處理和客戶端的通訊
func process(conn net.Conn) {
	
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
