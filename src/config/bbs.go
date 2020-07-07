package app

type bbsOption struct {
	HeadImageGirl, HeadImageBoy, VideoSnapshot string
	CommentTemplate                            []string
}

var Bbs = bbsOption{
	HeadImageGirl: "https://hxsapp-oss.hxsapp.com/public/image/default_head_img_new.png",
	HeadImageBoy:  "https://hxsapp-oss.hxsapp.com/bg_default_boy.png",
	VideoSnapshot: "https://hxsapp-user-media-out.oss-cn-hangzhou.aliyuncs.com/snapshots/videos/%s.jpg",
	CommentTemplate: []string{
		`官方认证沙发💯`,
		`疯狂打call☎`,
		`我要给你生猴子🙈`,
		`比心💞`,
		`给大神跪了🙇`,
		`爱了爱了😘`,
		`太硬核了💪`,
		`这波怒赞😎`,
		`这是什么神仙操作😳`,
		`你可真是个宝藏女孩😍`,
		`长得这么好看，一定是个男孩子😜`,
		`这个吼吼吃😋`,
		`就喜欢你可可爱爱的样子👀`,
		`解锁新技能👌`,
		`已阅，嗦粉，撸剧🍜`,
		`仙女本仙💃`,
		`吸一口仙气😏`,
		`加油，奥利给✊`,
		`捕获一只超A的小姐姐💋`,
		`减脂一时爽，一直减一直爽😁`,
		`内容引起极度舒适😂`,
		`这身材 ，我酸了🍋`,
		`呐，给你加鸡腿🍗`,
	}, // 评论模板
}
