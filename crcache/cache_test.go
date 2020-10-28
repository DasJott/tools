package cache_test

import (
	"testing"

	cache "cleverreach.com/crtools/crcache"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestCache(t *testing.T) {
	test := assert.New(t)

	c := cache.New()

	val := c.Get("mykey", func() interface{} {
		return "Hello World"
	}, false)
	test.EqualValues("Hello World", val)

	val = c.Get("mykey", func() interface{} {
		return "Furzkissen"
	}, false)
	test.EqualValues("Hello World", val)

	val = c.Get("mykey", func() interface{} {
		return "Furzkissen"
	}, true)
	test.EqualValues("Furzkissen", val)

	c.Delete("mykey")
	val = c.Get("mykey", func() interface{} {
		return "Fallera"
	}, false)
	test.EqualValues("Fallera", val)

	c.Set("mykey", "how cool")
	val = c.Get("mykey", func() interface{} {
		return "Fallera"
	}, false)
	test.EqualValues("how cool", val)
}
