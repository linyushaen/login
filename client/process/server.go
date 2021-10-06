package process

import (
	"fmt"
	"net"
	"os"
	"項目/chatroom/common/message"
	"項目/chatroom/server/utils"
)

// 顯示登錄成功後的界面..
func ShowMenu() {
	fmt.Println("-------恭喜xxx登錄成功---------")
	fmt.Println("-------1. 顯示在線用戶列表-----")
	fmt.Println("-------2. 發送消息------------")
	fmt.Println("-------3. 信息列表------------")
	fmt.Println("-------4. 退出系統------------")
	fmt.Print("請選擇(1-4):")
	var key int
	fmt.Scanln(&key)
	switch key {
	case 1:
		fmt.Println("顯示在線用戶列表")
	case 2:
		fmt.Println("發送消息")
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("你選擇退出了系統...")
		os.Exit(0)
	default:
		fmt.Println("你輸入的選項不正確")
	}
}

// 和服務器端保持通訊
func serverProcessMes(conn net.Conn) {
	// 創建一個transfer實例, 不停的讀取服務器發送的消息
	tf := &utils.Transfer{Conn: conn}
	for {
		fmt.Println("客戶端正在等待讀取服務器發送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err=", err)
			return
		}
		// 如果讀取到消息, 又是下一步處理邏輯
		switch mes.Type {
		case message.NotifyUserStatusMesType: //有人上線了
			// 處理
			// 1. 取出.NotityUserStatusMes
			// 2. 把這個用戶的信息, 狀態保存到客戶map[int]User中
		default:
			fmt.Println("服務器端返回了未知的消息類型")
		}
		// fmt.Printf("mes=%v\n", mes)
	}

}
