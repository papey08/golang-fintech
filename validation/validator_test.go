package go_course_validation

import (
	"errors"
	"testing"

	"github.com/papey08/golang-fintech/validation/check"

	"github.com/stretchr/testify/assert"
)

type NestedStruct struct {
	N int    `validate:"max:5"`
	S string `validate:"len:3"`
}

// StructWithNestedStructs is a struct for testing nested validation
type StructWithNestedStructs struct {
	NestedStruct1 NestedStruct
	NestedStruct2 NestedStruct
}

func TestValidate(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "invalid struct: interface",
			args: args{
				v: new(any),
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: map",
			args: args{
				v: map[string]string{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: string",
			args: args{
				v: "some string",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "valid struct with no fields",
			args: args{
				v: struct{}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with untagged fields",
			args: args{
				v: struct {
					f1 string
					f2 string
				}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with unexported fields",
			args: args{
				v: struct {
					foo string `validate:"len:10"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrValidateForUnexportedFields.Error()
			},
		},
		{
			name: "invalid validator syntax",
			args: args{
				v: struct {
					Foo string `validate:"len:abcdef"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrInvalidValidatorSyntax.Error()
			},
		},
		{
			name: "valid struct with tagged fields",
			args: args{
				v: struct {
					Len       string `validate:"len:20"`
					LenZ      string `validate:"len:0"`
					InInt     int    `validate:"in:20,25,30"`
					InNeg     int    `validate:"in:-20,-25,-30"`
					InStr     string `validate:"in:foo,bar"`
					MinInt    int    `validate:"min:10"`
					MinIntNeg int    `validate:"min:-10"`
					MinStr    string `validate:"min:10"`
					MinStrNeg string `validate:"min:-1"`
					MaxInt    int    `validate:"max:20"`
					MaxIntNeg int    `validate:"max:-2"`
					MaxStr    string `validate:"max:20"`
				}{
					Len:       "abcdefghjklmopqrstvu",
					LenZ:      "",
					InInt:     25,
					InNeg:     -25,
					InStr:     "bar",
					MinInt:    15,
					MinIntNeg: -9,
					MinStr:    "abcdefghjkl",
					MinStrNeg: "abc",
					MaxInt:    16,
					MaxIntNeg: -3,
					MaxStr:    "abcdefghjklmopqrst",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong length",
			args: args{
				v: struct {
					Lower    string `validate:"len:24"`
					Higher   string `validate:"len:5"`
					Zero     string `validate:"len:3"`
					BadSpec  string `validate:"len:%12"`
					Negative string `validate:"len:-6"`
				}{
					Lower:    "abcdef",
					Higher:   "abcdef",
					Zero:     "",
					BadSpec:  "abc",
					Negative: "abcd",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong in",
			args: args{
				v: struct {
					InA     string `validate:"in:ab,cd"`
					InB     string `validate:"in:aa,bb,cd,ee"`
					InC     int    `validate:"in:-1,-3,5,7"`
					InD     int    `validate:"in:5-"`
					InEmpty string `validate:"in:"`
				}{
					InA:     "ef",
					InB:     "ab",
					InC:     2,
					InD:     12,
					InEmpty: "",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong min",
			args: args{
				v: struct {
					MinA string `validate:"min:12"`
					MinB int    `validate:"min:-12"`
					MinC int    `validate:"min:5-"`
					MinD int    `validate:"min:"`
					MinE string `validate:"min:"`
				}{
					MinA: "ef",
					MinB: -22,
					MinC: 12,
					MinD: 11,
					MinE: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong max",
			args: args{
				v: struct {
					MaxA string `validate:"max:2"`
					MaxB string `validate:"max:-7"`
					MaxC int    `validate:"max:-12"`
					MaxD int    `validate:"max:5-"`
					MaxE int    `validate:"max:"`
					MaxF string `validate:"max:"`
				}{
					MaxA: "efgh",
					MaxB: "ab",
					MaxC: 22,
					MaxD: 12,
					MaxE: 11,
					MaxF: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 6)
				return true
			},
		},
		{
			name: "wrong field of type []int",
			args: args{
				v: struct {
					MinNums []int `validate:"min:0"`
					MaxNums []int `validate:"max:10"`
				}{
					MinNums: []int{9, 10, 11},
					MaxNums: []int{9, 10, 11},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 1)
				return true
			},
		},
		{
			name: "wrong field of type []string",
			args: args{
				v: struct {
					ShortStrings []string `validate:"len:10"`
					LongStrings  []string `validate:"len:10"`
				}{
					ShortStrings: []string{"abc", "def", "ghi"},
					LongStrings:  []string{"abcdefghij", "klmnopqrst"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 1)
				return true
			},
		},
		{
			name: "wrong fields of type []int and []string",
			args: args{
				v: struct {
					ShortStrings []string `validate:"min:5"`
					LongStrings  []string `validate:"min:5"`
					SmallNums    []int    `validate:"max:10"`
					BigNums      []int    `validate:"max:10"`
					PrimeNums    []int    `validate:"in:2,3,5,7,11"`
					PiDigits     []int    `validate:"in:2,3,5,7,11"`
				}{
					ShortStrings: []string{"abc", "def", "ghijk"},
					LongStrings:  []string{"AntonOcean", "kuai6", "mikhail-chebakov", "TimRazumov"},
					SmallNums:    []int{5, 4, 3, 2, 1},
					BigNums:      []int{2904, 46447, 1210},
					PrimeNums:    []int{5, 5, 7, 2, 3, 2},
					PiDigits:     []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 3)
				return true
			},
		},
		{
			name: "all valid fields",
			args: args{
				v: struct {
					ShortStrings []string `validate:"max:5"`
					LongStrings  []string `validate:"min:5"`
					SmallNums    []int    `validate:"max:10"`
					BigNums      []int    `validate:"min:10"`
					PrimeNums    []int    `validate:"in:2,3,5,7,11"`
					PiDigits     []int    `validate:"in:0,1,2,3,4,5,6,7,8,9"`
				}{
					ShortStrings: []string{"abc", "def", "ghijk"},
					LongStrings:  []string{"AntonOcean", "kuai6", "mikhail-chebakov", "TimRazumov"},
					SmallNums:    []int{5, 4, 3, 2, 1},
					BigNums:      []int{2904, 46447, 1210},
					PrimeNums:    []int{5, 5, 7, 2, 3, 2},
					PiDigits:     []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5},
				},
			},
			wantErr: false,
		},
		{
			name: "correct struct with nested structs",
			args: args{
				v: StructWithNestedStructs{
					NestedStruct1: NestedStruct{
						N: 1,
						S: "abc",
					},
					NestedStruct2: NestedStruct{
						N: 2,
						S: "def",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "struct with 1 invalid nested struct",
			args: args{
				v: StructWithNestedStructs{
					NestedStruct1: NestedStruct{
						N: 5,
						S: "ghi",
					},
					NestedStruct2: NestedStruct{
						N: 6,
						S: "jklmno",
					},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 1)
				return true
			},
		},
		{
			name: "valid Ad struct",
			args: args{
				v: struct {
					ID        int64
					Title     string `validate:"lenInterval:1,99"`
					Text      string `validate:"lenInterval:1,499"`
					AuthorID  int64
					Published bool
				}{
					ID:        10,
					Title:     "Ad with valid title and text",
					Text:      "Text of the valid ad",
					AuthorID:  2,
					Published: false,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid Ad struct",
			args: args{
				v: struct {
					ID        int64
					Title     string `validate:"lenInterval:1,99"`
					Text      string `validate:"lenInterval:1,499"`
					AuthorID  int64
					Published bool
				}{
					ID:        10,
					Title:     "Ad with empty text",
					Text:      "",
					AuthorID:  2,
					Published: false,
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == check.ErrInvalidFieldValue.Error()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
