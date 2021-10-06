package main

import (
	"fmt"
	"os"
	"項目/chatroom/client/process"
)

// 定義兩個變量, 一個表示用戶id, 一個表示用戶密碼
var userId int
var userPwd string
var userName string

func main() {

	// 接收用戶的選擇
	var key int
	// 判斷是否還繼續顯示菜單
	// var loop = true

	for {
		fmt.Println("----------------歡迎登入多人聊天系統----------------")
		fmt.Println("\t\t\t 1 登陸聊天室")
		fmt.Println("\t\t\t 2 註冊用戶")
		fmt.Println("\t\t\t 3 退出系統")
		fmt.Println("\t\t\t 請選擇(1-3):")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登陸聊天室")
			fmt.Print("請輸入用戶的id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Print("請輸入用戶的密碼:")
			fmt.Scanf("%s\n", &userPwd)
			// loop = false
			// 完成登錄
			// 1. 創建一個UserProcess的實例
			up := &process.UserProcess{}
			up.Login(userId, userPwd)
		case 2:
			fmt.Println("註冊用戶")
			fmt.Println("請輸入用戶id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("請輸入用戶密碼:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("請輸入用戶名字(nickname):")
			fmt.Scanf("%s\n", &userName)
			// 2. 調用UserProcess, 完成註冊的請求
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)
		case 3:
			fmt.Println("退出系統")
			// loop = false
			// 也可以下面
			os.Exit(0)
		default:
			fmt.Println("你的輸入有誤, 請重新輸入")
		}

	}

	// 根據用戶的輸入, 顯示新的提示信息
	// if key == 1 {
	// 	// 說明用戶要登陸
	// 	fmt.Print("請輸入用戶的id:")
	// 	fmt.Scanf("%d\n", &userId)
	// 	fmt.Print("請輸入用戶的密碼:")
	// 	fmt.Scanf("%s\n", &userPwd)

	// 	// 因為

	// 	// 先把登陸的函數, 寫到另外一個文件, 比如login.go
	// 	// 這裡我們會需要重新調用
	// 	test := &process.UserProcess{}
	// 	test.Login(userId, userPwd)
	// 	// if err != nil {
	// 	// 	fmt.Println("登陸失敗")
	// 	// } else {
	// 	// 	fmt.Println("登陸成功")
	// 	// }

	// } else if key == 2 {
	// 	fmt.Println("2")
	// }

}
