package access

import (
	"github.com/BastionPay/bas-bkadmin-api/models"
)

func (this *AccessList) Tree(list []*models.Access) []*models.Access {
	data := this.buildData(list)
	result := this.makeTreeCore(-1, data)

	return result
}

func (this *AccessList) buildData(list []*models.Access) map[int64]map[int64]*models.Access {
	var data = make(map[int64]map[int64]*models.Access)
	for _, v := range list {
		id := v.Id
		parentId := v.ParentId
		if _, ok := data[parentId]; !ok {
			data[parentId] = make(map[int64]*models.Access)
		}

		data[parentId][id] = v
	}

	return data
}

func (this *AccessList) makeTreeCore(index int64, data map[int64]map[int64]*models.Access) []*models.Access {
	tmp := make([]*models.Access, 0)
	for id, item := range data[index] {
		if data[id] != nil {
			item.Children = this.makeTreeCore(id, data)
		}

		tmp = append(tmp, item)
	}
	//sort.Sort(AccessListSort(tmp))
	return tmp
}
