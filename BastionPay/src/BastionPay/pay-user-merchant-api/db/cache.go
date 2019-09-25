package db


var GCache Cache

type Cache struct{
	//SponsorActivityCache   gcache.Cache
	//SponsorAkCache   gcache.Cache
	//SponsorIdCache   gcache.Cache
	//ActivityCache   gcache.Cache
	//PageCache      gcache.Cache
	//ShareInfoCache gcache.Cache
	//PageShareInfoCache      gcache.Cache
	//RobberCache gcache.Cache
	//
	//
	////走网络 缓存
	//FissionApiActivityYqlCache    gcache.Cache
	//FissionApiActivityListCache    	  gcache.Cache
}

func (this * Cache)Init(){
}

//**************************************************/
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