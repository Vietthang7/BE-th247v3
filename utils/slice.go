package utils

// Hàm trên trả về những phần tử trong mảng parents mà không có trong mảng subs.
func FindTheOtherElems[T any](parents, subs []T) (theOthers []T) {
	subMap := make(map[any]struct{})
	for _, sub := range subs {
		subMap[sub] = struct{}{}
	}

	for _, id := range parents {
		if _, ok := subMap[id]; !ok {
			theOthers = append(theOthers, id)
		}
	}

	return
}
