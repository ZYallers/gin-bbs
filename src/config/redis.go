package app

import (
	"os"
	"time"
)

type redisClient struct {
	Host, Port, Pwd string
	Db              int
}

type bbs struct {
	ZSetDiaryCircleOwner      string // 圈子详情页的动态　待审
	ZSetDiaryCircle           string // 圈子详情页的动态　已审
	ZSetDiaryList             string // 首页动态
	ZSetDiaryOwner            string // 个人主页动态
	ZSetUnionDiary            string // 圈子详情页的动态　并集
	HashDiary                 string // 动态详情
	ZSetCircleList            string // 圈子列表
	HashCircle                string // 圈子详情
	HashSpecialPopulation     string // 达人详情
	StringRedisLock           string
	DiaryCommentHash          string // 指定动态的评论列表
	DiaryCommentCounterString string // 指定动态的评论总数
	DiaryLikeZSet             string // 指定动态的点赞列表
	DiaryLikeCounterString    string // 指定动态的点赞总数
	UserLikeDiaryString       string // 用户是否对指定动态点赞过, 格式为 [userId]_[diaryId], 值为 1 的时候代表点赞过, 其余情况都代表未点赞
}

type session struct {
	StringSessToken string
}

type redisKey struct {
	Session session
	Bbs     bbs
}

type redis struct {
	Client               map[string]redisClient
	Key                  redisKey
	NormalTTL, NoDataTTL time.Duration
}

var (
	Redis = redis{
		Client: map[string]redisClient{
			"cache": {
				Host: os.Getenv("redis_host"),
				Port: os.Getenv("redis_port"),
				Pwd:  os.Getenv("redis_password"),
				Db:   0,
			},
			"ssd": {
				Host: os.Getenv("redis_host"),
				Port: os.Getenv("redis_port"),
				Pwd:  os.Getenv("redis_password"),
				Db:   0,
			},
			"session": {
				Host: os.Getenv("redis_session_host"),
				Port: os.Getenv("redis_session_port"),
				Pwd:  os.Getenv("redis_session_password"),
				Db:   0,
			},
		},
		Key: redisKey{
			Session: session{
				StringSessToken: "ci_session:",
			},
			Bbs: bbs{
				ZSetDiaryCircleOwner:      "bbs@diary:circle:owner:zset:",
				ZSetDiaryCircle:           "bbs@diary:circle:zset",
				ZSetDiaryList:             "bbs@diary:list:zset",
				ZSetDiaryOwner:            "bbs@diary:owner:zset",
				ZSetUnionDiary:            "bbs@diary:union:zset:",
				HashDiary:                 "bbs@diary:hash",
				ZSetCircleList:            "bbs@circle:list:zset",
				HashCircle:                "bbs@circle:hash",
				HashSpecialPopulation:     "bbs@special:population:hash",
				StringRedisLock:           "bbs@redis:key:",
				DiaryCommentHash:          "bbs@diarycomment:hash:%d",
				DiaryCommentCounterString: "bbs@diarycommentcounter:string:%d",
				DiaryLikeZSet:             "bbs@diarylike:zset:%d",
				DiaryLikeCounterString:    "bbs@diarylikecounter:string:%d",
				UserLikeDiaryString:       "bbs@userlikediary:string:%d_%d",
			},
		},
		NormalTTL: 432000 * time.Second, // 5d
		NoDataTTL: 10 * time.Second,     // 10s
	}
)
