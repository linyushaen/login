package model

import (
	"errors"
)

// 根據業務邏輯需要, 自定義一些錯誤.

var (
	ERROR_USER_NOTEXISTS = errors.New("用戶不存在..")
	ERROR_USER_EXISTS    = errors.New("用戶已經存在...")
	ERROR_USER_PWD       = errors.New("密碼不正確")
)
