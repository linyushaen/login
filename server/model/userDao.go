package model

import (
	"encoding/json"
	"fmt"
	"項目/chatroom/common/message"

	"github.com/garyburd/redigo/redis"
)

// 我們在服務器啟動後, 就初始化一個userDao實例,
// 把它做成全局的變量, 在需要和redis操作時, 就直接使用即可
var (
	MyUserDao *UserDao
)

// 定義一個UserDao 結構體
// 完成對User 結構體的各種操作.

type UserDao struct {
	pool *redis.Pool
}

// 使用工廠模式, 創建一個UserDao實例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {

	userDao = &UserDao{
		pool: pool,
	}
	return
}

// 思考一下在UserDao 應該提供哪些方法給我們
// 1. 根據用戶id 返回 一個User實例+err
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {

	// 通過給定id 去 redis查詢這個用戶
	res, err := redis.String(conn.Do("hget", "users", id))
	if err != nil {
		// 錯誤!
		if err == redis.ErrNil { // 表示在 users 哈希中, 沒有找到對應id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	user = &User{}

	// 這裡我們需要把res 反序列化成User實例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

// 完成登錄的校驗 Login
// 1. Login 完成對用戶的驗證
// 2. 如果用戶的id和pwd都正確, 則返回一個user實例
// 3. 如果用戶的id或pwd有錯誤, 則返回對應的錯誤信息
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {

	// 先從UserDao 的連接池中取出一根連接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}

	// 這時證明這個用戶是獲取到.
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

func (this *UserDao) Register(user *message.User) (err error) {

	// 先從UserDao 的連接池中取出一根連接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}

	// 這時, 說明id在redis還沒有, 則可以完成註冊
	data, err := json.Marshal(user) // 序列化
	if err != nil {
		return
	}

	// 入庫
	_, err = conn.Do("Hset", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存註冊用戶錯誤 err=", err)
		return
	}

	return
}
