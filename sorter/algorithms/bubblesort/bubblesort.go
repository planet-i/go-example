package bubblesort

func BubbleSort(values []int) {
	flag := true
	for i := 0; i < len(values)-1; i++ {
		flag = true
		for j := 0; j < len(values)-i-1; j++ {
			if values[j] > values[j+1] {
				values[j], values[j+1] = values[j+1], values[j] //大石头往后沉
				flag = false
			}
		}
		if flag == true {
			break
		}
	}
}
