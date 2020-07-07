package service

import (
	"testing"
)

// TestUserInfo 测试两个获取用户信息的接口返回数据是否一致
func TestUserInfo(t *testing.T) {
	userInfo := NewUserInfo().GetUserInfo([]int{6, 10, 11})
	shortUserInfo := NewUserInfo().GetShortUserInfo([]int{6, 10, 11})

	if len(userInfo) != len(shortUserInfo) {
		t.Error(`返回结果长度不一致`, userInfo, shortUserInfo)
	}

	for userId, row := range shortUserInfo {
		if rowUserInfo, ok := userInfo[userId]; ok {
			if headImg, ok := rowUserInfo.(map[string]interface{})[`head_img`]; ok {
				if row.(map[string]interface{})[`head_img`] != headImg {
					t.Error(`头像地址不一致`, headImg, row.(map[string]interface{})[`head_img`])
				}
			} else {
				t.Error(`缺少头像地址`, userInfo[userId], row)
			}
			if nickname, ok := rowUserInfo.(map[string]interface{})[`nickname`]; ok {
				if row.(map[string]interface{})[`nickname`] != nickname {
					t.Error(`昵称不一致`, nickname, row.(map[string]interface{})[`nickname`])
				}
			} else {
				t.Error(`缺少昵称`, userInfo[userId], row)
			}
		} else {
			t.Error(`user_id 不存在`, userId, userInfo, shortUserInfo)
		}
	}
}