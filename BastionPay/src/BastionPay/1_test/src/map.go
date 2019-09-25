package src

import (
	"fmt"
)

//map的遍历是无序的，可以用sort包的方法进行排序


func main()  {
	type Role struct {
		Id     int64
		Name   string
		Status string
	}

	type c []Role
	role1:=Role{6,"sa1","1"}
	role2:=Role{7,"sa3","2"}
	role3:=Role{8,"sa4","3"}

	m := make(map[int64]map[int64]Role)

	n := make(map[int64]Role)
	j := make(map[int64]Role)
	k := make(map[int64]Role)

	n[1]=role1
	j[2]=role2
	k[3]=role3

	m[1]=n
	m[2]=j
	m[3]=k
	fmt.Println(m)

//sort.Sort()
}
//
//func (this *Role) build2(r []*Role) map[int64]map[int64]*Role {
//	var m = make(map[int64]map[int64]*Role)
//	type Role struct {
//		Id     int64
//		Name   string
//		Status string
//	}
//	role1:=Role{6,"sa1","1"}
//	role2:=Role{7,"sa3","2"}
//	role3:=Role{8,"sa4","3"}
//
//
//	n := make(map[int64]*Role)
//	j := make(map[int64]*Role)
//	k := make(map[int64]*Role)
//
//	n[1]=&role1
//	j[2]=&role2
//	k[3]=&role3
//
//	m[1]=n
//
//
//	return m
//}



		//
		//
		//for k,v:=range role{
		//	index = append(index,k)
		//	recipientList = append(recipientList, v.Status+v.Name)
		//}
		//
		//fmt.Println(recipientList)
		//fmt.Println(m)




