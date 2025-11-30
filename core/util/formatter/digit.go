package formatter

import (
	"fmt"
	"math"
)

func FormatDigits(f float64, totalDigits int) string {
	// 計算需要幾位小數 = 總位數 - 整數部分的位數
	intPart := math.Abs(f) // 先取絕對值
	intDigits := 0

	if intPart >= 1 {
		intDigits = int(math.Floor(math.Log10(intPart))) + 1
	} else if intPart >= 0 {
		intDigits = 1
	}

	// 整數部分超過總位數則不顯示小數
	decimalPlaces := max(totalDigits-intDigits, 0)

	// fmt 格式化補0
	return fmt.Sprintf("%.*f", decimalPlaces, f)
}
