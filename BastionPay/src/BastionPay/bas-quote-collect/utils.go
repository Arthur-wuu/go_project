package main

//func ToApiMoneyInfo(m *collect.MoneyInfo) *apiquote.MoneyInfo {
//	m2 := new(apiquote.MoneyInfo)
//	m2.Symbol = m.Symbol
//	m2.Price = m.Price
//	m2.Volume_24h  = m.Volume_24h
//	m2.Market_cap = m.Market_cap
//	m2.Percent_change_1h  = m.Percent_change_1h
//	m2.Percent_change_24h = m.Percent_change_24h
//	m2.Percent_change_7d  = m.Percent_change_7d
//	m2.Last_updated      = m.Last_updated
//	return  m2
//}
//
//func ToApiCodeInfo(c * collect.CodeInfo) *apiquote.CodeInfo  {
//	c2 := new(apiquote.CodeInfo)
//	c2.Id  = c.Id
//	c2.Name = c.Name
//	c2.Symbol  = c.Symbol
//	c2.Website_slug = c.Website_slug
//	c2.Timestamp   = c.Timestamp
//	return c2
//}
//
//func ToLocalCodeInfo(c *apiquote.CodeInfo )  *collect.CodeInfo{
//	c2 := new(collect.CodeInfo)
//	c2.Id  = c.Id
//	c2.Name = c.Name
//	c2.Symbol  = c.Symbol
//	c2.Website_slug = c.Website_slug
//	c2.Timestamp   = c.Timestamp
//	return c2
//}
//
//func ToApiKXian(c *quote.KXian)*apiquote.KXian {
//	c2 := new(apiquote.KXian)
//	c2.Timestamp = c.Timestamp
//	c2.LowPrice = c.LowPrice
//	c2.HighPrice = c.HighPrice
//	c2.LastPrice = c.LastPrice
//	c2.ClosePrice = c.ClosePrice
//	c2.OpenPrice = c.OpenPrice
//	return c2
//}
//
//func ToApiKXians(c []*quote.KXian)[]*apiquote.KXian {
//	if c == nil || len(c) == 0 {
//		return nil
//	}
//	arr := make([]*apiquote.KXian, 0)
//	for i:=0; i < len(c); i++ {
//		c2 := new(apiquote.KXian)
//		c2.Timestamp = c[i].Timestamp
//		c2.LowPrice = c[i].LowPrice
//		c2.HighPrice = c[i].HighPrice
//		c2.LastPrice = c[i].LastPrice
//		c2.ClosePrice = c[i].ClosePrice
//		c2.OpenPrice = c[i].OpenPrice
//		arr = append(arr, c2)
//	}
//	return arr
//}
