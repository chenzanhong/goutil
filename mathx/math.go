package goutil

// 快速幂：计算 (base^exp) % mod
func PowMod(base, exp, mod int) int {
	result := 1
	for exp > 0 {
		if exp%2 == 1 {
			result = (result * base) % mod
		}
		base = (base * base) % mod
		exp /= 2
	}
	return result
}

// 计算组合数 C(n, k) % mod，使用费马小定理求逆元
func Combination(n, k, mod int) int {
	if k > n || k < 0 {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}

	// C(n, k) = n! / (k! * (n-k)!)
	// 使用模逆元：a^(-1) ≡ a^(mod-2) % mod
	numerator := 1   // n!
	denominator := 1 // k! * (n-k)!

	for i := 1; i <= k; i++ {
		numerator = (numerator * (n - i + 1)) % mod
		denominator = (denominator * i) % mod
	}

	// C(n,k) = numerator / denominator % mod
	inv := PowMod(denominator, mod-2, mod)
	return (numerator * inv) % mod
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MergeSort 对 []int 进行归并排序，返回新的已排序切片
func MergeSort(arr []int) []int {
	// 基础情况：如果数组长度小于等于1，直接返回
	if len(arr) <= 1 {
		return arr
	}

	// 分割数组
	mid := len(arr) / 2
	left := MergeSort(arr[:mid])  // 递归排序左半部分
	right := MergeSort(arr[mid:]) // 递归排序右半部分

	// 合并两个已排序的数组
	return merge(left, right)
}

// merge 合并两个已排序的切片并返回新的已排序切片
func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	// 比较并合并
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// 添加剩余元素
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}
