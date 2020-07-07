// Copyright (c) 2020 HXS R&D Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
//
// @Title account_user_info
// @Description
//
// @Author zhongyongbiao
// @Version 1.0.0
// @Time 2020/5/28 下午5:13
// @Software GoLand
package service

import (
	"github.com/tidwall/gjson"
	"net/http"
	"src/abs"
	app "src/config"
	"src/library/tool"
	"strconv"
	"strings"
	"time"
)

type userInfo struct {
	abs.Service
}

func NewUserInfo() *userInfo {
	return &userInfo{}
}

// GetUserInfo 批量获取用户信息
func (ui *userInfo) GetUserInfo(userId []int) (result map[int]interface{}) {
	result = make(map[int]interface{})
	if len(userId) > 0 {
		var candidate []string
		for _, v := range userId {
			candidate = append(candidate, strconv.Itoa(v)) // 拼装 user_id
		}

		param := make(map[string]interface{})
		param[`user_id`] = strings.Join(candidate, `,`) // 多个 user_id 使用半角逗号分隔
		body := tool.HttpRequestWithSign(app.Sdk.Account.GetUserInfo, param)

		if body != `` {
			remoteData := gjson.Parse(body).Get(`data`).Value() // 只要 data 字段, 丢弃 code, msg 等返回值
			if mapData, ok := remoteData.(map[string]interface{}); ok {
				for k, v := range mapData {
					intUserId, _ := strconv.Atoi(k) // 转换键值类型, 使用 int 类型的 user_id 方便后续处理
					result[intUserId] = v
				}
			}
		}
	}
	return result
}

// GetShortUserInfo 批量获取用户信息, 精简版
// 只返回 head_img/nickname/sex/user_id
func (ui *userInfo) GetShortUserInfo(userId []int) map[int]interface{} {
	result := map[int]interface{}{}
	if len(userId) > 0 {
		var candidate []string
		for _, v := range userId {
			candidate = append(candidate, strconv.Itoa(v)) // 拼装 user_id
		}

		param := make(map[string]interface{})
		param[`ids`] = strings.Join(candidate, `,`) // 多个 user_id 使用半角逗号分隔
		body := tool.HttpRequestWithSign(app.Sdk.Account.GetShortUserInfo, param, http.MethodPost, 3*time.Second)

		if body != `` && gjson.Parse(body).Get(`data.list`).IsObject() {
			remoteData := gjson.Parse(body).Get(`data.list`).Value() // 只要 data 字段, 丢弃 code, msg 等返回值
			if mapData, ok := remoteData.(map[string]interface{}); ok {
				for k, v := range mapData {
					intUserId, _ := strconv.Atoi(k) // 转换键值类型, 使用 int 类型的 user_id 方便后续处理
					result[intUserId] = v
				}
			}
		}
	}
	return result
}

// AppendInfo 批量追加用户信息
// field 需要追加的字段, 例如 head_img, nickname, 参考 http://wiki.sys.hxsapp.net/pages/viewpage.action?pageId=7343723
func (ui *userInfo) AppendInfo(data []map[string]interface{}, field []string) (result []map[string]interface{}) {
	result = make([]map[string]interface{}, len(data))
	if len(data) > 0 {
		var userId []int
		for _, v := range data {
			if !tool.InArray(v[`user_id`].(int), userId) { // 去重处理
				userId = append(userId, v[`user_id`].(int))
			}
		}

		userList := ui.GetShortUserInfo(userId)

		for k, v := range data {
			result[k] = v
			for _, fieldName := range field {
				if existUser, ok := userList[v[`user_id`].(int)]; ok {
					if fieldValue, ok := existUser.(map[string]interface{})[fieldName]; ok {
						result[k][fieldName] = fieldValue
					} else {
						result[k][fieldName] = `` // 如果 account 域没有返回这个字段, 默认返回空字符串
					}
				} else {
					result[k][fieldName] = `` // 如果 account 域没有返回这个字段, 默认返回空字符串
				}
			}
		}
	}
	return result
}
