package db

//import (
//	"github.com/bluele/gcache"
//	"BastionPay/merchant-teammanage-api/config"
//	"errors"
//	"time"
//	"fmt"
//	"github.com/caibirdme/yql"
//)
//
//
//var GCache Cache
//
//type Cache struct{
//	SponsorActivityCache   gcache.Cache
//	SponsorAkCache   gcache.Cache
//	SponsorIdCache   gcache.Cache
//	ActivityCache   gcache.Cache
//	PageCache      gcache.Cache
//	ShareInfoCache gcache.Cache
//	PageShareInfoCache      gcache.Cache
//	RobberCache gcache.Cache
//
//
//	//走网络 缓存
//	FissionApiActivityYqlCache    gcache.Cache
//	FissionApiActivityListCache    	  gcache.Cache
//}
//
//func (this * Cache)Init(){
//}
//
////**************************************************/
//func (this * Cache) GetFissionApiActyList(tp int) (interface{}, error) {
//	if this.FissionApiActivityListCache == nil {
//		return nil, errors.New("not init")
//	}
//	value, err := this.FissionApiActivityListCache.Get(tp)
//	if err != nil {
//		return nil, err
//	}
//
//	return value, nil
//}
//
//func (this * Cache) SetFissionApiActyList(f func(key interface{}) (interface{}, *time.Duration, error)) {
//	this.FissionApiActivityListCache = gcache.New(config.GConfig.Cache.SponsorMaxKey).LRU().LoaderExpireFunc(f).Build()
//}
//
//func (this * Cache) RemoveFissionApiActyList(tp int) {
//	if this.FissionApiActivityListCache == nil {
//		return
//	}
//	this.FissionApiActivityListCache.Remove(tp)
//}
//
////**************************************************/
//func (this * Cache) GetActyYqlByTpAndId(tp int, actyId string) (interface{}, error) {
//	if this.FissionApiActivityYqlCache == nil {
//		return nil, errors.New("not init")
//	}
//	value, err := this.FissionApiActivityYqlCache.Get(fmt.Sprintf("%d_%s", tp, actyId))
//	if err != nil {
//		return nil, err
//	}
//
//	return value, nil
//}
//
//func (this * Cache) SetActyYqlFunc()  {
//	this.FissionApiActivityYqlCache = gcache.New(config.GConfig.Cache.SponsorMaxKey).LRU().Build()
//}
//
//func (this * Cache) SetActyYql(tp int, actyId string, ruler yql.Ruler)  {
//	this.FissionApiActivityYqlCache.SetWithExpire(fmt.Sprintf("%d_%d", tp, actyId), ruler, time.Hour * 24)
//}
//
//func (this * Cache) RemoveActyYql(uuid string) {
//	if this.FissionApiActivityYqlCache == nil {
//		return
//	}
//	this.FissionApiActivityYqlCache.Remove(uuid)
//}
//
////**************************************************/
//
////**************************************************/
//func (this * Cache) GetSponsorByAk(apikey string) (interface{}, error) {
//	if this.SponsorAkCache == nil {
//		return nil, errors.New("not init")
//	}
//	value, err := this.SponsorAkCache.Get(apikey)
//	if err != nil {
//		return nil, err
//	}
//
//	return value, nil
//}
//
//func (this * Cache) SetSponsorAkFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
//	this.SponsorAkCache = gcache.New(config.GConfig.Cache.SponsorMaxKey).LRU().LoaderExpireFunc(f).Build()
//}
//
//func (this * Cache) RemoveSponsorByAk(uuid string) {
//	if this.SponsorAkCache == nil {
//		return
//	}
//	this.SponsorAkCache.Remove(uuid)
//}
//
////**************************************************/
//func (this * Cache) GetSponsorById(id int) (interface{}, error) {
//	if this.SponsorIdCache == nil {
//		return nil, errors.New("not init")
//	}
//	value, err := this.SponsorIdCache.Get(id)
//	if err != nil {
//		return nil, err
//	}
//
//	return value, nil
//}
//
//func (this * Cache) SetSponsorIdFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
//	this.SponsorIdCache = gcache.New(config.GConfig.Cache.SponsorMaxKey).LRU().LoaderExpireFunc(f).Build()
//}
//
//func (this * Cache) RemoveSponsorById(id int) {
//	if this.SponsorIdCache == nil {
//		return
//	}
//	this.SponsorIdCache.Remove(id)
//}
//
////**************************************************/
//func (this * Cache) GetSponsorActivity(spId int) (interface{}, error) {
//	if this.SponsorActivityCache == nil {
//		return nil, errors.New("not init")
//	}
//	value, err := this.SponsorActivityCache.Get(spId)
//	if err != nil {
//		return nil, err
//	}
//
//	return value, nil
//}
//
//func (this * Cache) SetSponsorActivityFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
//	this.SponsorActivityCache = gcache.New(config.GConfig.Cache.SponsorActivityMaxKey).LRU().LoaderExpireFunc(f).Build()
//}
//
//func (this * Cache) RemoveSponsorActivity(spId int) {
//	if this.SponsorActivityCache == nil {
//		return
//	}
//	this.SponsorActivityCache.Remove(spId)
//}
//
////**************************************************/
//func (this * Cache) GetActivity(uuid string) (interface{}, error) {
//	if this.ActivityCache == nil {
//		return nil, errors.New("not init")
//	}
//	value, err := this.ActivityCache.Get(uuid)
//	if err != nil {
//		return nil, err
//	}
//
//	return value, nil
//}
//
//func (this * Cache) SetActivityFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
//	this.ActivityCache = gcache.New(config.GConfig.Cache.ActivityMaxKey).LRU().LoaderExpireFunc(f).Build()
//}
//
//func (this * Cache) RemoveActivity(uuid string) {
//	if this.ActivityCache == nil {
//		return
//	}
//	this.ActivityCache.Remove(uuid)
//}
////**************************************************/
//func (this * Cache) GetPage(uuid string) (interface{}, error) {
//	if this.PageCache == nil {
//		return nil, errors.New("not init")
//	}
//	value, err := this.PageCache.Get(uuid)
//	if err != nil {
//		return nil, err
//	}
//
//	return value, nil
//}
//
//func (this * Cache) RemovePage(uuid string)  {
//	if this.PageCache == nil {
//		return
//	}
//	this.PageCache.Remove(uuid)
//}
//
//func (this * Cache) SetPageFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
//	this.PageCache = gcache.New(config.GConfig.Cache.PageMaxKey).LRU().LoaderExpireFunc(f).Build()
//}
//
////****************************************************/
//
//func (this * Cache) GetShareInfo(uuid string) (interface{}, error) {
//	if this.ShareInfoCache == nil {
//		return nil, errors.New("not init")
//	}
//	value, err := this.ShareInfoCache.Get(uuid)
//	if err != nil {
//		return nil, err
//	}
//
//	return value, nil
//}
//
//func (this * Cache) RemoveShareInfo(uuid string)  {
//	if this.ShareInfoCache == nil {
//		return
//	}
//	this.ShareInfoCache.Remove(uuid)
//}
//
//func (this * Cache) SetShareInfoFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
//	this.ShareInfoCache = gcache.New(config.GConfig.Cache.ShareInfoMaxKey).LRU().LoaderExpireFunc(f).Build()
//}
//
////****************************************************/
//type RobberRedPhoneKey struct{
//	RedUuid string
//	CountryCode string
//	Phone  string
//}
//
//func (this * Cache) GetRobber(key *RobberRedPhoneKey) (interface{}, error) {
//	if this.RobberCache == nil {
//		return nil, errors.New("not init")
//	}
//
//	value, err := this.RobberCache.Get(*key)
//	if err != nil {
//		return nil, err
//	}
//
//	return value, nil
//}
//
//func (this * Cache) RemoveRobber(key *RobberRedPhoneKey)  {
//	if this.RobberCache == nil {
//		return
//	}
//	this.RobberCache.Remove(*key)
//}
//
//func (this * Cache) SetRobberFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
//	this.RobberCache = gcache.New(config.GConfig.Cache.RobberMaxKey).LRU().LoaderExpireFunc(f).Build()
//}
//
////****************************************************/
//
////**************************************************/
//func (this * Cache) GetPageShareInfo(uuid string) (interface{}, error) {
//	if this.PageShareInfoCache == nil {
//		return nil, errors.New("not init")
//	}
//	value, err := this.PageShareInfoCache.Get(uuid)
//	if err != nil {
//		return nil, err
//	}
//
//	return value, nil
//}
//
//func (this * Cache) RemovePageShareInfo(uuid string)  {
//	if this.PageShareInfoCache == nil {
//		return
//	}
//	this.PageShareInfoCache.Remove(uuid)
//}
//
//func (this * Cache) SetPageShareInfoFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
//	this.PageShareInfoCache = gcache.New(config.GConfig.Cache.PageMaxKey).LRU().LoaderExpireFunc(f).Build()
//}
