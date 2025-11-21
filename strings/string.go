package goutil

// buildLPS 构建最长公共前后缀数组（LPS Array）
func buildLPS(pattern string) []int {
	m := len(pattern)
	lps := make([]int, m)
	length := 0 // 当前最长公共前后缀的长度
	i := 1

	// 从第二个字符开始构建 lps 数组
	for i < m {
		if pattern[i] == pattern[length] {
			length++
			lps[i] = length
			i++
		} else {
			if length != 0 {
				// 回退到之前的最长前缀
				length = lps[length-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}
	return lps
}

// KMP 搜索主函数，返回所有匹配位置的起始索引
func KMP(text, pattern string) []int {
	var result []int
	n := len(text)
	m := len(pattern)

	if m == 0 {
		return result // 空模式，返回空
	}

	// 构建 LPS 数组
	lps := buildLPS(pattern)

	i := 0 // text 的索引
	j := 0 // pattern 的索引

	for i < n {
		if text[i] == pattern[j] {
			i++
			j++
		}

		if j == m {
			// 找到一个匹配
			result = append(result, i-j)
			j = lps[j-1] // 继续查找下一个匹配
		} else if i < n && text[i] != pattern[j] {
			if j != 0 {
				j = lps[j-1]
			} else {
				i++
			}
		}
	}

	return result
}
