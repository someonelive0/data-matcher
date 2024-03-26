package engine

import (
	"fmt"
	"testing"

	"github.com/dexyk/stringosim"
)

// From https://www.golangcodes.com/forum.php?mod=viewthread&tid=74
func TestJaccard(t *testing.T) {
	// Jaccard
	// Jaccard距离可以通过设置n-gram的大小来计算，这将用于比较。如果省略了该大小，则将使用默认值1。
	// 输出分值越小表示两者越接近
	srcTxt := "我是中国人"
	targetTxt := "我美国人"
	fmt.Println("stringosim.Jaccard", srcTxt, targetTxt)
	fmt.Println(stringosim.Jaccard([]rune(srcTxt), []rune(targetTxt), []int{1}))
	// 0.4666666666666667
	fmt.Println(stringosim.Jaccard([]rune(srcTxt), []rune(targetTxt), []int{2}))
	// 0.6111111111111112
	fmt.Println(stringosim.Jaccard([]rune(srcTxt), []rune(targetTxt), []int{3}))
	// 0.7222222222222222

	srcTxt = "我是中国人"
	targetTxt = "我中国人"
	fmt.Println("\nstringosim.Jaccard", srcTxt, targetTxt)
	fmt.Println(stringosim.Jaccard([]rune(srcTxt), []rune(targetTxt), []int{1}))
	// 0.16666666666666663
	fmt.Println(stringosim.Jaccard([]rune(srcTxt), []rune(targetTxt), []int{2}))
	// 0.33333333333333337
	fmt.Println(stringosim.Jaccard([]rune(srcTxt), []rune(targetTxt), []int{3}))
	// 0.4666666666666667

	fmt.Println("\nstringosim Jaro and Jaro-Winkler", srcTxt, targetTxt)
	// Jaro and Jaro-Winkler
	// 值越大，则越相似
	fmt.Println(stringosim.Jaro([]rune(srcTxt), []rune(targetTxt)))
	// 0.9333333333333332
	fmt.Println(stringosim.JaroWinkler([]rune(srcTxt), []rune(targetTxt)))
	// 0.94

	// QGram
	// 越小越相似
	fmt.Println("\nstringosim QGram", srcTxt, targetTxt)
	fmt.Println(stringosim.QGram([]rune(srcTxt), []rune(targetTxt)))
	// 5

	// Cosine
	// 越小越相似
	fmt.Println("\nstringosim Cosine", srcTxt, targetTxt)
	fmt.Println(stringosim.Cosine([]rune(srcTxt), []rune(targetTxt)))
	// 0.19417703597461977

	// Levenshtein
	// 越小越相似
	fmt.Println("\nstringosim Levenshtein", srcTxt, targetTxt)
	fmt.Println(stringosim.Levenshtein([]rune(srcTxt), []rune(targetTxt)))
	// 1
}
