package util

///不重复压入[]string
func AppendUniqueString(data *[]string, ID string) {
	exist := false
	for _, v := range *data {
		if v == ID {
			exist = true
			break
		}
	}
	if !exist {
		*data = append(*data, ID)
	}
}


///不重复压入[]string
func AppendUniqueStringEx(data []string, ID string)(res []string) {
	exist := false
	for _, v := range data {
		if v == ID {
			exist = true
			break
		}
	}
	if !exist {
		data = append(data, ID)
	}
	return data
}
///删除存在的ID(string)
func DelExistString(data *[]string, ID string) {
	index := -1
	for i, v := range *data {
		if v == ID {
			index = i
			break
		}
	}
	tmp := *data
	if index != -1 {
		tmp = append(tmp[:index], tmp[index+1:]...)
	}
	*data = tmp
}
