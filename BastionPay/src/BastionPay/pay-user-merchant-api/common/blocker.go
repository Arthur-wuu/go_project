package common

import (
	//"go.uber.org/zap"
	"errors"
	"github.com/kataras/iris"
	"strconv"
	"sync"
)

const (
	TermBlockLimitPrefix = "TermBlockLimit_"
	TermBlockLockPrefix  = "TermBlockLock_"
)

// 在t1时间段内达到n次限制将锁死账户t2时间；这和limiter有差异
type TermBlock struct {
	Name     string
	Limit    int64
	Time     int
	Locktime int
}

type TermBlockResponse struct {
	Remain_count int64   //剩余次数
	Lock_time    float64 //小时
	OnBlock      bool    //是否被锁
	OpenFlag     bool    //该功能是否开启
}

func (this *TermBlockResponse) update(count int64, tb *TermBlock) {
	remain := tb.Limit - count
	if this.Remain_count == 0 || remain < this.Remain_count {
		this.Remain_count = remain
		this.Lock_time = float64(tb.Locktime) / 3600
	}
	if remain <= 0 {
		this.OnBlock = true
	}
}

var TermBlockMap map[string][]*TermBlock
var termBlockLock sync.Mutex

func NewTermBlocker(redis *DbRedis, pls []*TermBlock, blockName, account string, ctx iris.Context) *TermBlocker {
	if TermBlockMap != nil {
		return &TermBlocker{redis, TermBlockMap, blockName, GetRealIp(ctx), account}
	}

	TermBlockMapTemp := make(map[string][]*TermBlock, 0)
	for i := 0; i < len(pls); i++ {
		if pls[i] == nil || pls[i].Time == 0 || pls[i].Limit == 0 || pls[i].Locktime == 0 {
			continue
		}
		mm, ok := TermBlockMapTemp[pls[i].Name]
		if !ok {
			mm = make([]*TermBlock, 0)
		}
		mm = append(mm, pls[i])
		TermBlockMapTemp[pls[i].Name] = mm
	}

	termBlockLock.Lock()
	defer termBlockLock.Unlock()
	TermBlockMap = TermBlockMapTemp

	return &TermBlocker{redis, TermBlockMap, blockName, GetRealIp(ctx), account}
}

//短周期锁账号
type TermBlocker struct {
	redis     *DbRedis
	pls       map[string][]*TermBlock
	blockName string
	ip        string
	account   string
}

func (this *TermBlocker) IsBlock() (bool, error) {
	//是否 被lock
	_, ok := this.pls[this.blockName]
	if !ok {
		return false, nil
	}
	result, err := this.redis.Do("GET", this.genLockKey())
	if err != nil {
		return false, err
	}
	if result == nil {
		return false, nil
	}
	return true, nil
}

func (this *TermBlocker) Done(clearFlag bool) (*TermBlockResponse, error) {
	res := new(TermBlockResponse)
	lts, ok := this.pls[this.blockName]
	if !ok {
		return res, nil
	}
	if clearFlag {
		return res, this.clearAll(lts)
	}

	res.OpenFlag = true
	for i := 0; i < len(lts); i++ { //这里应该不能批处理了
		count, err := this.incrLimit(this.genLimitKey(lts[i]), lts[i].Time)
		if err != nil {
			return res, err
		}
		res.update(count, lts[i])
		if res.OnBlock {
			return res, this.setBlock(lts[i].Locktime)
		}
	}
	return res, nil
}

func (this *TermBlocker) clearAll(tbs []*TermBlock) error {
	keyArr := make([]interface{}, 0)
	for i := 0; i < len(tbs); i++ {
		if tbs[i] == nil {
			continue
		}
		keyArr = append(keyArr, this.genLimitKey(tbs[i]))
	}
	keyArr = append(keyArr, this.genLockKey())
	_, err := this.redis.Do("del", keyArr...)
	return err
}

func (this *TermBlocker) incrLimit(key string, timeout int) (int64, error) {
	result, err := this.redis.Do("INCR", key)
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, nil
	}

	res, ok := result.(int64)
	if !ok {
		return 0, errors.New("redis type err")
	}
	if res == 1 {
		_, err = this.redis.Do("EXPIRE", key, timeout)
		if err != nil {
			return 0, err
		}
	}
	return res, nil
}

func (this *TermBlocker) setBlock(timeout int) error {
	_, err := this.redis.Do("setex", this.genLockKey(), timeout, "1")
	return err
}

func (this *TermBlocker) genLimitKey(tb *TermBlock) string {
	return TermBlockLimitPrefix + this.account + "_" + this.blockName + "_" + strconv.Itoa(tb.Time)
}

func (this *TermBlocker) genLockKey() string {
	return TermBlockLockPrefix + this.account + "_" + this.blockName
}
