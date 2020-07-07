package app

type account struct {
	GetUserInfo       string // 批量获取用户信息, 使用 mget 优化
	GetShortUserInfo  string // 获取用户信息
	GetUserVipVisible string // 用户等级
}

type im struct {
	YunXinIMToken string // 注册云信用户
}

type bonus struct {
	GetUserBean      string
	ConsumeUserBean  string
	RechargeUserBean string
}

type mall struct {
	CouponListById string // Mall优惠券
}

type sdk struct {
	IM      im
	Account account
	Bonus   bonus
	Mall    mall
}

var Sdk = sdk{
	IM: im{
		YunXinIMToken: `http://im.hxsapp.com/im/YunXinIM/token`,
	},
	Account: account{
		GetUserInfo:       `http://account.hxsapp.com/user/UserInfo/getUserInfo`,
		GetShortUserInfo:  `http://account.hxsapp.com/user/userInfo/getShortUserInfo`,
		GetUserVipVisible: `http://account.hxsapp.com/user/UserMember/getUserVipVisible`,
	},
	Bonus: bonus{
		GetUserBean:      `http://bonus.hxsapp.com/bean/UserBean/getUserBean`,
		ConsumeUserBean:  `http://bonus.hxsapp.com/bean/UserBean/consume`,
		RechargeUserBean: `http://bonus.hxsapp.com/bean/UserBean/recharge`,
	},
	Mall: mall{
		CouponListById: `http://mall.hxsapp.com/base/Coupon/getCouponsListById`,
	},
}
