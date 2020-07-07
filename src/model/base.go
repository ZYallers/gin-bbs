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
// @Title base
// @Description
//
// @Author zhongyongbiao
// @Version 1.0.0
// @Time 2020/5/28 下午5:40
// @Software GoLand
package model

import "src/abs"

const (
	StateDelete  = 0 // 删除
	StateOnline  = 1 // 已上线
	StateOffline = 2 // 下线
)

type baseModel struct {
	abs.Model
}

// DefaultState 默认写入新数据的状态（为了将来全部内容进入审核库）
func (b *baseModel) DefaultState() int {
	return StateOnline
}
