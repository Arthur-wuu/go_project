package kbspirit

import "BastionPay/bas-tv-proxy/api"

// 搜索环境
type Env struct {
	// Markets  map[string]bool
	Input            string
	Markets          []string
	Markets_FullName []string
	Language         int64
	Count            int64
	Delist           bool
	Kuozhan          bool
	CategoryID       []string
	VipLevel         int
	SuffixFlag       bool
	ChineseFlag      bool
}

type KuoZhan struct {
	LeiXing   string //leixing
	ZiLeiXing string //zileixing
}

// NewSearchEnv 新建搜索环境
func NewSearchEnv() *Env {
	return &Env{
		// Markets:  make(map[string]bool),
		Markets:    make([]string, 0),
		Language:   0,
		Count:      0,
		Delist:     true,
		Kuozhan:    false,
		CategoryID: make([]string, 0),
	}
}

type Searcher interface {
	Init() error
	Search(key string, env *Env) []*api.JPBShuChu
	UnInit()
	Update(market string, arrs []*api.JPBShuJu)
}
