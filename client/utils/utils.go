package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"項目/chatroom/common/message"
)

// 這裡將這些方法關聯到結構體中
type Transfer struct {
	// 分析它應該有哪些字段
	Conn net.Conn
	Buf  [8096]byte // 這是傳輸時, 使用緩沖
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {

	// buf := make([]byte, 8096)
	fmt.Println("讀取客戶端發送的數據...")
	//
	// 如果客戶端關閉了 conn 則, 就不會阻塞
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		// fmt.Println("conn.Read err=", err)
		// err = errors.New("read pkg header error")
		return
	}
	// fmt.Println("buf[0:4]=", this.Buf[0:4])
	// 根據buf[:4] 轉成一個 uint32類型
	var pkgLen uint32 = binary.BigEndian.Uint32(this.Buf[0:4])
	// pkgLen = binary.BigEndian.Uint32(buf[0:4])

	// 根據 pkgLen 讀取消息內容
	n, err := this.Conn.Read(this.Buf[:pkgLen]) // 有覆寫過去, 所以別再想為何長度沒有被Unmarshal到了= =
	if n != int(pkgLen) || err != nil {
		// err = errors.New("read pkg body error")
		return
	}

	// 把pkgLen 反序列化成 -> message.Message

	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	return

}

func (this *Transfer) WritePkg(data []byte) (err error) {

	// 先發送一個長度給對方
	var pkgLen uint32 = uint32(len(data)) // 也可 64 或 16
	// var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	// 發送長度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail ", err)
		return
	}

	// 發送data本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(data) err=", err)
		return
	}
	return
}
