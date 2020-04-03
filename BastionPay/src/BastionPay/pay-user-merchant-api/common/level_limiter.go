package common

import (
	"errors"
	"github.com/kataras/iris"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	LevelLimterPrefix = "level_limiter_"
)

///////////////////////////分等级限制///////////////////////////////
type LevelPathLimit struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Name   string `json:"name"` //upper(method)+PATH
	Limit  int    `json:"limit"`
	Time   int    `json:"time"`
	Level  int    `json:"level"`
}

type LevelPathLimitRemain struct {
	LevelPathLimit
	Remian int `json:"remain"`
}

type LevelLimitList [][]*LevelPathLimit

func (this *LevelLimitList) Len() int { // 重写 Len() 方法
	return len(*this)
}

func (this *LevelLimitList) Swap(i, j int) { // 重写 Swap() 方法
	(*this)[i], (*this)[j] = (*this)[j], (*this)[i]
}
func (this *LevelLimitList) Less(i, j int) bool { //
	if len((*this)[j]) == 0 || len((*this)[i]) == 0 {
		return false
	}
	return (*this)[i][0].Level < (*this)[j][0].Level
}

/////////////////////////////////////////////////////////////

var LevelLimiterMap map[string]LevelLimitList
var lvLimitLock sync.Mutex

func DoLevelLimiter(redis *DbRedis, pls []*LevelPathLimit, ctx iris.Context, limitName, account string, level int) (bool, int, int, error) {
	return NewLevelLimiter(redis, pls, ctx, limitName, account, level).Verification()
}

//获取剩余的次数
func GetAllLevelRemainLimits(redis *DbRedis, pls []*LevelPathLimit, account string, level int) (map[string][]*LevelPathLimitRemain, error) {
	limiter := NewLevelLimiter(redis, pls, nil, "", account, level)
	return limiter.GetAllMatchLevelRemainLimits()
}

func NewLevelLimiter(redis *DbRedis, pls []*LevelPathLimit, ctx iris.Context, limitName, account string, level int) *LevelLimiter {
	if len(limitName) == 0 && ctx != nil {
		limitName = strings.ToUpper(ctx.Method()) + ctx.Path()
	}
	if LevelLimiterMap != nil {
		return &LevelLimiter{redis, LevelLimiterMap, limitName, GetRealIp(ctx), account, level}
	}
	LevelLimiterMap2 := make(map[string]LevelLimitList)
	for i := 0; i < len(pls); i++ {
		if len(pls[i].Name) == 0 {
			pls[i].Name = strings.ToUpper(pls[i].Method) + pls[i].Path
		}
		name := pls[i].Name

		llArr, ok := LevelLimiterMap2[name]
		if !ok {
			llArr = make(LevelLimitList, 0)
		}

		llsIndex := -1
		for j := 0; j < len(llArr); j++ {
			if llArr[j][0].Level == pls[i].Level {
				llsIndex = j
				break
			}
		}
		if llsIndex < 0 {
			lls := make([]*LevelPathLimit, 0)
			llArr = append(llArr, lls)
			llsIndex = len(llArr) - 1
		}
		llArr[llsIndex] = append(llArr[llsIndex], pls[i])
		LevelLimiterMap2[name] = llArr
	}
	for _, v := range LevelLimiterMap2 {
		sort.Sort(&v)
	}
	lvLimitLock.Lock()
	defer lvLimitLock.Unlock()
	LevelLimiterMap = LevelLimiterMap2

	return &LevelLimiter{redis, LevelLimiterMap, limitName, GetRealIp(ctx), account, level}
}

type LevelLimiter struct {
	redis   *DbRedis
	plMap   map[string]LevelLimitList //string=method+path，切片 按照level排序
	name    string
	ip      string
	account string
	level   int
}

//是否限制，剩余次数(-1表示无限制)，时间，err
func (this *LevelLimiter) Verification() (bool, int, int, error) {
	plArr, ok := this.plMap[this.name]
	if !ok || len(plArr) == 0 {
		return true, -1, 0, nil
	}
	var vaildPls []*LevelPathLimit
	for i := len(plArr) - 1; i >= 0; i-- { //level不多，遍历足够
		if len(plArr[i]) == 0 {
			continue
		}
		if this.level >= plArr[i][0].Level {
			vaildPls = plArr[i]
			break
		}
	}
	if vaildPls == nil || len(vaildPls) == 0 {
		return true, -1, 0, nil
	}

	minRemian := -1
	minTime := int(0)
	for i := 0; i < len(vaildPls); i++ {
		count, err := this.countDown(vaildPls[i])
		if err != nil {
			return false, 0, 0, err
		}
		remian := vaildPls[i].Limit - count
		if remian < 0 {
			return false, vaildPls[i].Limit, vaildPls[i].Time, nil
		}
		if minRemian < 0 || remian < minRemian {
			minRemian = remian
			minTime = vaildPls[i].Time
		}
	}
	return true, minRemian, minTime, nil
}

func (this *LevelLimiter) countDown(pl *LevelPathLimit) (int, error) {
	key := this.genLimitKey(pl)

	result, err := this.redis.Do("INCR", key)
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, nil
	}

	res := int(result.(int64))
	if res == 1 {
		_, err = this.redis.Do("EXPIRE", key, pl.Time)
		if err != nil {
			return 0, err
		}
	}

	return res, err
}

func (this *LevelLimiter) genLimitKey(pl *LevelPathLimit) string {
	name := this.name
	if len(name) == 0 {
		name = pl.Name
	}
	if len(this.account) > 0 {
		return LevelLimterPrefix + this.account + "_" + name + "_" + strconv.Itoa(pl.Level) + "_" + strconv.Itoa(pl.Time)
	} else {
		return LevelLimterPrefix + this.ip + "_" + name + "_" + strconv.Itoa(pl.Level) + "_" + strconv.Itoa(pl.Time)
	}
}

func (this *LevelLimiter) GetAllMatchLevels() [][]*LevelPathLimit {
	res := make(LevelLimitList, 0)
	for _, llArr := range LevelLimiterMap {
		index := -1
		for i := len(llArr) - 1; i >= 0; i-- { //level不多，遍历足够
			if len(llArr[i]) == 0 {
				continue
			}
			if this.level >= llArr[i][0].Level {
				index = i
				break
			}
		}
		if index < 0 {
			continue
		}
		res = append(res, llArr[index])
	}
	return res
}

func (this *LevelLimiter) GetAllMatchLevelRemainLimits() (map[string][]*LevelPathLimitRemain, error) {
	llArr := this.GetAllMatchLevels()
	//redis 操作
	keys := make([]interface{}, 0)
	for j := 0; j < len(llArr); j++ {
		for i := 0; i < len(llArr[j]); i++ {
			key := this.genLimitKey(llArr[j][i])
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return nil, nil
	}
	reply, err := this.redis.Do("mget", keys...)
	if err != nil {
		return nil, err
	}
	replyArr, ok := reply.([]interface{})
	if !ok {
		return nil, errors.New("types err")
	}

	index := 0
	res := make(map[string][]*LevelPathLimitRemain, 0)
	for j := 0; j < len(llArr); j++ {
		tt := make([]*LevelPathLimitRemain, 0)
		for i := 0; i < len(llArr[j]); i++ {
			count := -1
			switch replyArr[index].(type) {
			case nil:
				count = llArr[j][i].Limit
			case []byte, string:
				repStr, _ := replyArr[index].([]byte)
				count, _ = strconv.Atoi(string(repStr))
				count = llArr[j][i].Limit - count
				if count < 0 {
					count = 0
				}
			}
			index++

			lvRemain := &LevelPathLimitRemain{
				LevelPathLimit{
					llArr[j][i].Method,
					llArr[j][i].Path,
					llArr[j][i].Name,
					llArr[j][i].Limit,
					llArr[j][i].Time,
					llArr[j][i].Level},
				count,
			}
			tt = append(tt, lvRemain)
		}
		if len(tt) > 0 {
			res[tt[0].Name] = tt
		}
	}
	return res, nil
}
