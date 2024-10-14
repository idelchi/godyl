package stringlike

func SetIfEmpty[S ~string](input *S, value S) {
	if *input == "" {
		*input = value
	}
}

// func SetSliceIfNil[S ~string](input *[]S, values ...S) {
// 	if *input == nil {
// 		*input = append([]S(nil), values...)
// 	}
// }

func SetSliceIfNil[S ~[]T, T any](input *S, values ...T) {
	if *input == nil {
		*input = append([]T(nil), values...)
	}
}

func IsEmpty[S ~string](input S) bool {
	return input == ""
}
