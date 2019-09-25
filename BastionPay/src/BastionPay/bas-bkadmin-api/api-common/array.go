package common

import (
	"reflect"
	"sort"
)

func (t *Tools) Duplicate(a interface{}) (ret []interface{}) {
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}

		ret = append(ret, va.Index(i).Interface())
	}
	return ret
}

func (t *Tools) SliceRemoveDuplicates(slice []string) []string {
	sort.Strings(slice)

	i := 0

	var j int

	for {
		if i >= len(slice)-1 {
			break
		}

		for j = i + 1; j < len(slice) && slice[i] == slice[j]; j++ {

		}

		slice = append(slice[:i+1], slice[j:]...)
		i++
	}

	return slice
}

/**
 * 计算两个数组差
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) ArrayDiff(newRoleId, metaRoleId []string) []string {
	var newIdList []string
	for _, val := range newRoleId {

		var flag bool
		for _, v := range metaRoleId {
			if val == v {
				flag = true
				continue
			}
		}

		if !flag {
			newIdList = append(newIdList, val)
		}
	}

	return newIdList
}

/**
 * 数据去重
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) Unique(slice []string) []string {
	keys := make(map[string]bool)
	var list []string

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}
