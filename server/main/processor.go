package main

import (
	"fmt"
	"io"
	"net"
	"項目/chatroom/common/message"
	processes "項目/chatroom/server/process"
	"項目/chatroom/server/utils"
)

// 先創建一個Processor 的結構體
type Processor struct {
	Conn net.Conn
}

// 功能: 根據客戶端發送消息種類不同, 決定調用哪個函數來處理
func (this *Processor) ServerProcessMes(mes *message.Message) (err error) {

	switch mes.Type {
	case message.LoginMesType:
		// 處理登陸
		// 創建一個UserProcess實例
		up := &processes.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		// 處理註冊
		up := &processes.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	default:
		fmt.Println("消息類型不存在, 無法處理...")
	}
	return
}

func (this *Processor) process2() (err error) {
	// 循環讀客戶端發送的信息
	for {
		// 這裡我們將讀取數據包, 直接封裝成一個函數readPkg(), 返回Message, Err
		// 創建一個Transfer 實例完成讀包任務
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客戶端退出, 服務器端也退出..")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}

		}
		// fmt.Println("mes=", mes)
		err = this.ServerProcessMes(&mes)
		if err != nil {
			return err
		}

	}
}
