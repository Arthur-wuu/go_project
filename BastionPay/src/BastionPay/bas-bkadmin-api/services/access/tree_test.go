package access

import (
	"github.com/BastionPay/bas-bkadmin-api/models"
)

func (this *AccessList) Tree2(list []*models.Access) []*models.Access {
	data := this.buildData2(list)
	result := this.makeTreeCore2(-1, data)

	return result
}

//func (this *AccessList) buildData2(list []*models.Access) map[int64]map[int64]*models.Access {
//	var data = make(map[int64]map[int64]*models.Access)
//	for _, v := range list {
//		id := v.Id
//		parentId := v.ParentId
//		if _, ok := data[parentId]; !ok {
//			data[parentId] = make(map[int64]*models.Access)
//		}
//
//		data[parentId][id] = v
//	}
//
//	return data
//}

//将
func (this *AccessList) buildData2(list []*models.Access) []*models.Access {
	data := []*models.Access{}

	return []*models.Access{}
}

func (this *AccessList) makeTreeCore2(index int64, data map[int64]map[int64]*models.Access) []*models.Access {
	tmp := make([]*models.Access, 0)
	for id, item := range data[index] {
		if data[id] != nil {
			item.Children = this.makeTreeCore(id, data)
		}

		tmp = append(tmp, item)
	}

	return tmp
}
