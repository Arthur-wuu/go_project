package src

import "fmt"

func quickSort2(theArray []int, start int, end int)[]int {
	if (start<end){
		m, n := start, end
		base := theArray[m]
		for {
			if (m < n){              // 0 < 6
				for{
					if((m < n)&&(theArray[n]>=base)){
						n--
					}else{
						theArray[m] = theArray[n]
						break
					}
				}
				for{
					if((m < n)&&(theArray[m]<=base)){    // 10  10   1
						m++
					}else{
						theArray[n] = theArray[m]
						break
					}
				}
			}else{
				break
			}
			theArray[m] = base
			quickSort(theArray, start, m-1)
			quickSort(theArray, n+1, end)
		}
	}
	return theArray
}
func main() {
	var theArray = []int{10,9,7,88,6,3}
	fmt.Print("排序前")
	fmt.Println(theArray)
	fmt.Print("排序后")
	arrayResult := quickSort2(theArray,0,len(theArray)-1)
	fmt.Println(arrayResult)
}


