package main

import (
	"math/rand"
	"regexp"
	"fmt"
	"sort"

	//"net/http"

	//"strconv"
	//a "net/url"

	//"strconv"
	"strings"

	//"strconv"

	//"strings"
)

func main(){
	//s := "postgres://user:pass@host.com:5432/paths?k=v#f"
	//urls, _ := a.Parse(s)
	//maps := urls.Query()
	////maps["k"] = "sss"
	//fmt.Println("maps:",maps)
	//maps["k"][0] = "xsacsa"
	//
	//urls.String()
	//
	//
	////var uid *string
	////uidInt, _ := strconv.Atoi(*uid)
	//
	//u, _ := a.Parse(s)
	//query := u.Query()
	//fmt.Println("aaa", query)
	//m, _ := a.ParseQuery(u.RawQuery)
	//
	//fmt.Println(m)
	//fmt.Println(m["k"][0])
	//m.Set("k", "dwqdwq")
	//m.Encode()
	//fmt.Println(m)
	//fmt.Println(s)
	//
	//var ints int
	//fmt.Println(ints)

	//line := "imei		android"
	//data := strings.Split(line, "\t")
	////len(data)
	//fmt.Println(data[0],data[1],data[2], len(data))
str  := "wkreader://app/go/read?bookid=46998"
  arr := strings.Split(str, "&")

  fmt.Println(len(arr))
	fmt.Println(arr)
	//fmt.Println(len(arr))
	s := "157264&chapterid=41722869&force_to_chapter=true&extsourceid=wkr28022"
	s = "123456" + s[6:]
	fmt.Println(s)


	fmt.Println("" == " ")
	ids := []string{"2","4","66"}
	s = GetPkgsByIDs(ids)
	fmt.Println(s)


	bookReg, _ := regexp.Compile(`bookid=([\d]+)`)
	//b :=bookReg.FindStringSubmatch("hap://dsadsa?dsa=sd&dwd&bookid=1238")
	//fmt.Println(b)
	strs := "wklreader://app/go/read?bookid=157264&chapterid=41722869&force_to_chapter=true&extsourceid=wkr28022"
            // wklreader://app/go/read?bookid=33333&chapterid=41722869&force_to_chapter=true&extsourceid=wkr28022
	p := "33333"
	a := bookReg.FindStringIndex(strs)
	// 前a[0] 后a[1]
	pre := strs[:a[0]+7] + p
	end := strs[a[1]:]
	fmt.Println(pre + end)
	fmt.Println(end)






	fmt.Println(a)
	//length := len(b[1])

	num := 1
	i := rand.Intn(num)

	fmt.Println("i" , i)
	sreee :=GetExpidsStr([]uint32{1,3,4,555,654,2})
	fmt.Println("sreee" , sreee)

	sreee2 := GetExpidsStrUniq([]uint32{1,3,3322222,2,2222,2222,3,3,3,4,555,654,2})
	fmt.Println("sreee2" , sreee2)



	keysSlice := []string{"12","23"}
	lists := new(Lists)
	lists.List = make([]Body,len(keysSlice))
	for idx, key := range keysSlice {
		token := "ttt"
		lists.List[idx].Id = key
		lists.List[idx].Token = token
	}

	fmt.Println(*lists)



}

// GetExpidsStr 获取实验 id 以逗号分给
func GetExpidsStr(expids []uint32) (str string) {
	expidsStr := make([]string, 0, len(expids))
	for _, expid := range expids {
		expidsStr = append(expidsStr, fmt.Sprintf("%d", expid))
	}
	return strings.Join(expidsStr, ",")
}

func GetExpidsStrUniq(expids []uint32) (str string) {
	tempMap := make(map[string]int)
	for _, expid := range expids {
		_, ok := tempMap[fmt.Sprintf("%d", expid)]
		if !ok {
			tempMap[fmt.Sprintf("%d", expid)] = 1
		}
	}
	expidsStr := make([]string, 0, len(tempMap))
	for keys, _ := range tempMap {
		expidsStr = append(expidsStr, keys)
		sort.Strings(expidsStr)
	}
	return strings.Join(expidsStr, ",")
}



func  GetPkgsByIDs(ids []string) (pkgNameStr string) {
	pkgNames := make([]string, len(ids))

		pkgNames[0] = ids[0]
	pkgNames[1] = ids[1]
	pkgNames[2] = ids[2]
fmt.Println("pkgNames",pkgNames)
	return strings.Replace(strings.Trim(fmt.Sprint(pkgNames), "[]"), " ", ",", -1)
}


type Lists struct {
	List     []Body    `json:"list,omitempty"`
}

type Body struct {
	Id     string    `json:"id"`
	Token  string    `json:"token"`
}


type WxTokenParams struct {
	Key  string `form:"key" json:"key"`
}