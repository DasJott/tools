package meta_test

import (
	"reflect"
	"testing"

	"cleverreach.com/crtools/meta"
	"github.com/stretchr/testify/assert"
)

type Tester struct {
	FieldOne   string `json:"one" bson:"eins" sql:"uno" klingonisch:"wa',omitempty"`
	FieldTwo   string `json:"two" bson:"zwei" sql:"due" klingonisch:"cha',omitempty"`
	FieldThree string `json:"three" bson:"drei" sql:"tre" klingonisch:"wej,omitempty"`
	FieldFour  string
}

func TestGetFieldName(t *testing.T) {
	test := assert.New(t)

	inst := Tester{}

	v := reflect.TypeOf(inst)
	for i := 0; i < v.NumField(); i++ {
		switch i {
		case 0:
			name, omit := meta.GetJSONName(v.Field(i))
			test.Equal("one", name)
			test.False(omit)
			name, omit = meta.GetBSONName(v.Field(i))
			test.Equal("eins", name)
			test.False(omit)
			name, omit = meta.GetSQLName(v.Field(i))
			test.Equal("uno", name)
			test.False(omit)
			name, omit = meta.GetFieldName("klingonisch", v.Field(i))
			test.Equal("wa'", name)
			test.True(omit)
		case 1:
			name, omit := meta.GetJSONName(v.Field(i))
			test.Equal("two", name)
			test.False(omit)
			name, omit = meta.GetBSONName(v.Field(i))
			test.Equal("zwei", name)
			test.False(omit)
			name, omit = meta.GetSQLName(v.Field(i))
			test.Equal("due", name)
			test.False(omit)
			name, omit = meta.GetFieldName("klingonisch", v.Field(i))
			test.Equal("cha'", name)
			test.True(omit)
		case 2:
			name, omit := meta.GetJSONName(v.Field(i))
			test.Equal("three", name)
			test.False(omit)
			name, omit = meta.GetBSONName(v.Field(i))
			test.Equal("drei", name)
			test.False(omit)
			name, omit = meta.GetSQLName(v.Field(i))
			test.Equal("tre", name)
			test.False(omit)
			name, omit = meta.GetFieldName("klingonisch", v.Field(i))
			test.Equal("wej", name)
			test.True(omit)
		case 3:
			name, omit := meta.GetJSONName(v.Field(i))
			test.Equal("FieldFour", name)
			test.False(omit)
			name, omit = meta.GetBSONName(v.Field(i))
			test.Equal("FieldFour", name)
			test.False(omit)
			name, omit = meta.GetSQLName(v.Field(i))
			test.Equal("FieldFour", name)
			test.False(omit)
			name, omit = meta.GetFieldName("klingonisch", v.Field(i))
			test.Equal("FieldFour", name)
			test.False(omit)
		}
	}
}

func TestGetNames(t *testing.T) {
	test := assert.New(t)

	inst := Tester{}
	names := meta.GetFieldNames(inst, "sql")
	test.Len(names, 4)

	exp := []string{"uno", "due", "tre", "FieldFour"}
	for i := 0; i < len(names); i++ {
		test.Equal(exp[i], names[i])
	}

	// test the omitempty
	names = meta.GetFieldNames(inst, "klingonisch")
	test.Len(names, 1)
	test.Equal(exp[3], names[0])

}

func TestRangeFields(t *testing.T) {
	test := assert.New(t)

	inst := Tester{"alpha", "beta", "gamma", "delta"}
	exp := []string{"alpha", "beta", "gamma", "delta"}

	count, sum := 0, 0
	meta.RangeFields(inst, func(idx int, val interface{}) {
		sum += idx
		count++
		test.EqualValues(exp[idx], val)
	})
	test.Equal(4, count)
	test.Equal(6, sum)

}

func TestGetFieldValueMap(t *testing.T) {
	test := assert.New(t)

	inst := Tester{"alpha", "beta", "gamma", "delta"}
	keys := []string{"uno", "due", "tre", "FieldFour"}
	vals := []string{"alpha", "beta", "gamma", "delta"}

	count := 0

	themap := meta.GetFieldValueMap(inst, "sql")
	for i, key := range keys {
		test.EqualValues(vals[i], themap[key])
		count++
	}

	test.Equal(4, count)
}
