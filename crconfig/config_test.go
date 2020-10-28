package crconfig_test

import (
	"os"
	"testing"

	"cleverreach.com/crtools/crconfig"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestGets(t *testing.T) {
	test := assert.New(t)

	err := crconfig.Read("notExisting.env")
	test.NotNil(err)

	err = crconfig.Read("testdata.env")
	test.Nil(err)

	test.Equal("just the default", crconfig.Get("NOT_EXISTING", "just the default"))
	test.Equal(true, crconfig.GetBool("ALSO_NOT_EXISTING", true))
	test.EqualValues(3825, crconfig.GetInt("NO_NUMBER", 3825))
	test.EqualValues(8.9, crconfig.GetFloat("NO_FLOAT", 8.9))

	test.Equal("this is the first value", crconfig.Get("FIRST_VAL", "wrong"))
	test.Equal("cool things are hot", crconfig.Get("ANOTHER_THING", "wrong"))
	test.EqualValues(42, crconfig.GetInt("MY_NUMVBER", 23))
	test.EqualValues(3.45, crconfig.GetFloat("TRY_FLOAT", 1.2))
	test.Equal(true, crconfig.GetBool("WORKS_GREAT", false))
}

func TestBind(t *testing.T) {
	type Data struct {
		First    string  `env:"FIRST_VAL,katzenfurz"`
		Second   string  `env:"ANOTHER_THING"`
		Third    int     `env:"MY_NUMVBER"`
		Another  float32 `env:"TRY_FLOAT"`
		Flag     bool    `env:"WORKS_GREAT"`
		Default1 string  `env:"NOT_THERE,foo to the bar"`
		Default2 int     `env:"NOT_THERE,99"`
	}

	test := assert.New(t)

	err := crconfig.Read("testdata.env")
	test.Nil(err)

	d := Data{}
	err = crconfig.Bind(&d)
	test.Nil(err)

	test.Equal("this is the first value", d.First)
	test.Equal("cool things are hot", d.Second)
	test.EqualValues(42, d.Third)
	test.EqualValues(3.45, d.Another)
	test.Equal(true, d.Flag)

	test.Equal("foo to the bar", d.Default1)
	test.EqualValues(99, d.Default2)
}

func TestBindPartially(t *testing.T) {
	type Data1 struct {
		First  string `env:"FIRST_VAL"`
		Second string `env:"ANOTHER_THING"`
		Third  int    `env:"MY_NUMVBER"`
	}
	type Data2 struct {
		Flag     bool   `env:"WORKS_GREAT"`
		Default1 string `env:"NOT_THERE,foo to the bar"`
		Default2 int    `env:"NOT_THERE,99"`
	}

	test := assert.New(t)

	err := crconfig.Read("testdata.env")
	test.Nil(err)

	d1 := Data1{}
	crconfig.Bind(&d1)

	test.Equal("this is the first value", d1.First)
	test.Equal("cool things are hot", d1.Second)
	test.EqualValues(42, d1.Third)

	d2 := Data2{}
	crconfig.Bind(&d2)

	test.Equal(true, d2.Flag)
	test.Equal("foo to the bar", d2.Default1)
	test.EqualValues(99, d2.Default2)
}

func TestGetPrefix(t *testing.T) {
	test := assert.New(t)

	err := crconfig.Read("testdata.env")
	test.Nil(err)

	m := crconfig.GetWithPrefix("TEST_PREFIX_")
	test.Len(m, 3)

	test.Equal("this is the first", m["TEST_PREFIX_UNO"])
	test.Equal("second is here", m["TEST_PREFIX_DUE"])
	test.Equal("and the third", m["TEST_PREFIX_TRES"])
}

func TestSwitches(t *testing.T) {
	test := assert.New(t)

	{ // args with a command
		os.Args = []string{"test_cmd", "-n", "99", "-t", "Hallo Welt"}

		err := crconfig.Read("testdata.env")
		test.Nil(err)

		test.EqualValues(99, crconfig.GetInt("MY_NUMVBER", 23))
		test.EqualValues("Hallo Welt", crconfig.Get("TEST_PREFIX_UNO", "Mega"))
	}
	{ // only switches
		os.Args = []string{"-n", "3825", "-t", "Freudenschwein"}

		err := crconfig.Read("testdata.env")
		test.Nil(err)

		test.EqualValues(3825, crconfig.GetInt("MY_NUMVBER", 23))
		test.EqualValues("Freudenschwein", crconfig.Get("TEST_PREFIX_UNO", "Mega"))
	}
}

func TestVariables(t *testing.T) {
	test := assert.New(t)

	{ // file only
		err := crconfig.Read("testdata.env")
		test.Nil(err)

		test.Equal("Moinsen!", crconfig.Get("BASE_VALUE", ""), "BASE_VALUE")
		test.Equal("Moinsen!", crconfig.Get("COPY_VALUE1", ""), "COPY_VALUE1")
		test.Equal("Moinsen!", crconfig.Get("COPY_VALUE2", ""), "COPY_VALUE2")
		test.Equal("Moinsen!", crconfig.Get("WEIRD_VALUE", ""), "WEIRD_VALUE")
	}

	{ // file and one arg
		os.Args = []string{"test_cmd", "-b", "foose"}

		err := crconfig.Read("testdata.env")
		test.Nil(err)

		test.Equal("foose", crconfig.Get("BASE_VALUE", ""), "BASE_VALUE")
		test.Equal("foose", crconfig.Get("COPY_VALUE1", ""), "COPY_VALUE1")
		test.Equal("foose", crconfig.Get("COPY_VALUE2", ""), "COPY_VALUE2")
		test.Equal("foose", crconfig.Get("WEIRD_VALUE", ""), "WEIRD_VALUE")
	}

	{ // file and two args
		os.Args = []string{"test_cmd", "-b", "foose", "-c", "barse"}

		err := crconfig.Read("testdata.env")
		test.Nil(err)

		test.Equal("foose", crconfig.Get("BASE_VALUE", ""), "BASE_VALUE")
		test.Equal("barse", crconfig.Get("COPY_VALUE1", ""), "COPY_VALUE1")
		test.Equal("barse", crconfig.Get("COPY_VALUE2", ""), "COPY_VALUE2")
		test.Equal("barse", crconfig.Get("WEIRD_VALUE", ""), "WEIRD_VALUE")
	}
}
