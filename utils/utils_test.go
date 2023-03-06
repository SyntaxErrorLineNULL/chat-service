package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	// Test case 1
	first := []int{1, 2, 3}
	second := []int{4, 5, 6}
	expected := []int{1, 2, 3, 4, 5, 6}
	expectLen := 6
	result := Merge[int](first, second)

	assert.Equal(t, expected, result)
	assert.Equal(t, expectLen, len(result))

	// Test case 2
	firstTwo := []string{"a", "b", "c"}
	secondTwo := []string{"d", "e", "f"}
	expectedTwo := []string{"a", "b", "c", "d", "e", "f"}
	expectLenTwo := 6
	resultTwo := Merge[string](firstTwo, secondTwo)

	assert.Equal(t, expectedTwo, resultTwo)
	assert.Equal(t, expectLenTwo, len(resultTwo))

	// Test case 3
	firstTree := []float64{1.1, 2.2, 3.3}
	secondTree := []float64{4.4, 5.5, 6.6}
	expectedTree := []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6}
	expectLenTree := 6
	resultTree := Merge[float64](firstTree, secondTree)

	assert.Equal(t, expectedTree, resultTree)
	assert.Equal(t, expectLenTree, len(resultTree))
}
