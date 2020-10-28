//go:generate go get -u gopkg.in/h2non/gock.v1
//go:generate go get github.com/stretchr/testify/assert

package rest_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"cleverreach.com/crtools/rest"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestGet(t *testing.T) {
	test := assert.New(t)
	defer gock.Off()

	{ // response 200
		gock.New("http://test.com").Get("get").Reply(200).BodyString("krassgeil")
		resp := rest.
			NewClient("http://test.com/get").
			Get()

		test.Equal("krassgeil", resp.String())
		test.Nil(resp.Error)
		test.Equal(200, resp.Status)
	}

	{ // response 404
		gock.New("http://test.com").Get("noget1").Reply(404).BodyString("hau ab")
		resp := rest.
			NewClient("http://test.com/noget1").
			Get()

		test.Equal("hau ab", resp.String())
		test.Nil(resp.Error)
		test.Equal(404, resp.Status)
	}

	{ // token header match
		gock.New("http://test.com").Get("get").MatchHeader("Authorization", "jott 7h4t1sc00l").Reply(201).BodyString("yesss")
		resp := rest.
			NewClient("http://test.com/get").
			Token("jott", "7h4t1sc00l").
			Get()

		test.Equal("yesss", resp.String())
		test.Nil(resp.Error)
		test.Equal(201, resp.Status)
	}

	{ // token header not match
		gock.New("http://test.com").Get("noget").MatchHeader("Authorization", "jott 7h4t1sc00l").Reply(201).BodyString("noooo")
		resp := rest.
			NewClient("http://test.com/noget").
			Token("jott", "7h4t1sn07c00l").
			Get()

		test.Empty(resp.String())
		test.NotNil(resp.Error)
		test.Equal(0, resp.Status)
	}

	{ // some header match
		gock.New("http://test.com").Get("get").MatchHeader("x-testporn", "this is just a test").Reply(200).BodyString("wow, that header matches!")
		resp := rest.
			NewClient("http://test.com/get").
			Header("x-testporn", "this is just a test").
			Get()

		test.Equal("wow, that header matches!", resp.String())
		test.Nil(resp.Error)
		test.Equal(200, resp.Status)
	}

	{ // some header not match
		gock.New("http://test.com").Get("noget").MatchHeader("x-testpron", "this is just a keks").Reply(200).BodyString("no, that header doesn't match!")
		resp := rest.
			NewClient("http://test.com/noget").
			Header("x-testporn", "this is just a test").
			Get()

		test.Empty(resp.String())
		test.NotNil(resp.Error)
		test.Equal(0, resp.Status)
	}

	{ // some param match
		gock.New("http://test.com").Get("get").MatchParam("sdl", "42").Reply(200).BodyString("I love matching params")
		resp := rest.
			NewClient("http://test.com/get").
			Param("sdl", "42").
			Get()

		test.Equal("I love matching params", resp.String())
		test.Nil(resp.Error)
		test.Equal(200, resp.Status)
	}

	{ // some param not match
		gock.New("http://test.com").Get("noget").MatchParam("sdl", "42").Reply(200).BodyString("some things don't match")
		resp := rest.
			NewClient("http://test.com/noget").
			Param("sdl", "this is just a test").
			Get()

		test.Empty(resp.String())
		test.NotNil(resp.Error)
		test.Equal(0, resp.Status)
	}

}
func TestGetObject(t *testing.T) {
	test := assert.New(t)
	defer gock.Off()

	type testy struct {
		Field1 string
		Field2 int
	}

	{ // response of an obejct
		expect := testy{"hallo", 42}
		gock.New("http://test.com/getobject").Reply(200).JSON(&expect)

		res := testy{}
		resp := rest.
			NewClient("http://test.com/getobject").
			GetObject(&res)

		test.Nil(resp.Error)
		test.EqualValues(expect, res)
		test.Equal(200, resp.Status)
	}

	{ // reponse of another object
		gock.New("http://test.com/getobject").Reply(200).JSON(rest.M{"Field1": "FOO", "Field2": 88, "Field3": "gibbets nich"})

		res := testy{}
		resp := rest.
			NewClient("http://test.com/getobject").
			GetObject(&res)

		test.Nil(resp.Error)
		test.Equal("FOO", res.Field1)
		test.Equal(88, res.Field2)
		test.Equal(200, resp.Status)
	}

}

func TestNewInParts(t *testing.T) {
	test := assert.New(t)
	defer gock.Off()

	{
		path := rest.NewClient("/get").URL()
		test.Equal("/get", path)
	}
	{
		path := rest.NewClient("http://test.com", "wurst", "gummi/hans", "/get").URL()
		test.Equal("http://test.com/wurst/gummi/hans/get", path)
	}
	{
		path := rest.NewClient("test.com", "wurst", "gummi/hans", "/get").URL()
		test.Equal("test.com/wurst/gummi/hans/get", path)
	}
	{
		path := rest.NewClient("http://", "test.com", "wurst", "gummi/hans", "/get").URL()
		test.Equal("http://test.com/wurst/gummi/hans/get", path)
	}
	{
		path := rest.NewClient("meine", "wurst", "ist/aus", "gummi").URL()
		test.Equal("/meine/wurst/ist/aus/gummi", path)
	}
	{
		rest.BasePath = "https://"
		path := rest.NewClient("test.com", "/wurst", "gummi/hans", "/get").URL()
		test.Equal("https://test.com/wurst/gummi/hans/get", path)
	}
	{
		rest.BasePath = "https://fummel.com"
		path := rest.NewClient("wurst", "gummi/hans", "/get").URL()
		test.Equal("https://fummel.com/wurst/gummi/hans/get", path)
	}
	rest.BasePath = ""

	{
		gock.New("http://test.com/wurst/gummi/hans").Get("get").Reply(200).BodyString("krassgeil")
		resp := rest.
			NewClient("http://test.com", "wurst", "gummi/hans", "/get").
			Get()

		test.Equal("krassgeil", resp.String())
		test.Nil(resp.Error)
		test.Equal(200, resp.Status)
	}
}

func TestPost(t *testing.T) {
	test := assert.New(t)
	defer gock.Off()

	{ // matching post body
		gock.New("http://test.com").Post("post").BodyString("Hallelujah!").Reply(200).BodyString("Das is'n Ding!")
		buf := bytes.NewBufferString("Hallelujah!")
		resp := rest.
			NewClient("http://test.com/post").
			Post(buf)

		test.Equal("Das is'n Ding!", resp.String())
		test.Nil(resp.Error)
		test.Equal(200, resp.Status)
	}

	{ // not matching post body
		gock.New("http://test.com").Post("post").BodyString("Hallelujah!").Reply(200).BodyString("Das is'n Ding!")
		buf := bytes.NewBufferString("Heureka!")
		resp := rest.
			NewClient("http://test.com/post").
			Post(buf)

		test.NotEqual("Das is'n Ding!", resp.String())
		test.NotNil(resp.Error)
		test.Equal(0, resp.Status)
	}
}

func TestNew(t *testing.T) {
	test := assert.New(t)

	{
		client := rest.NewClient("http://alberto.com/joomla/?option=com_cleverreach&task=hookHandler")

		test.NotNil(client)
		test.Equal("http://alberto.com/joomla?option=com_cleverreach&task=hookHandler", client.URL())
	}

	{
		client := rest.NewClient("http://alberto.com/joomla?option=com_cleverreach&task=hookHandler")

		test.NotNil(client)
		test.Equal("http://alberto.com/joomla?option=com_cleverreach&task=hookHandler", client.URL())
	}
}

func TestParseURL(t *testing.T) {
	test := assert.New(t)
	defer gock.Off()

	{ // parse url and get with parameters
		gock.New("http://test.com").Get("endslash/").MatchParam("eins", "11").MatchParam("zwei", "22").Reply(200).BodyString("What a fat body this is")
		client := rest.NewClientFromURL("http://test.com/endslash/?eins=11&zwei=22")
		test.Equal("http://test.com/endslash/", client.URL())

		resp := client.Get()

		test.Equal("What a fat body this is", resp.String())
		test.Nil(resp.Error)
		test.Equal(200, resp.Status)
	}

	{ // parse url and get with parameters
		gock.New("http://test.com").Get("parsedget").MatchParam("eins", "11").MatchParam("zwei", "22").Reply(200).BodyString("What a fat body this is")
		resp := rest.
			NewClientFromURL("http://test.com/parsedget?eins=11&zwei=22").
			Get()

		test.Equal("What a fat body this is", resp.String())
		test.Nil(resp.Error)
		test.Equal(200, resp.Status)
	}

	{ // parse url and get with different parameters
		gock.New("http://test.com").Get("parsedget").MatchParam("keyvalue", "true").ParamPresent("OnlyParam").Reply(200).BodyString("What a fat body this is")
		resp := rest.
			NewClientFromURL("http://test.com/parsedget?keyvalue=true&OnlyParam").
			Get()

		test.Equal("What a fat body this is", resp.String())
		test.Nil(resp.Error)
		test.Equal(200, resp.Status)
	}

	{ // parse url - negative test
		gock.New("http://test.com").Get("parsedget").MatchParam("keyvalue", "true").ParamPresent("OnlyParam").Reply(200).BodyString("What a fat body this is")
		resp := rest.
			NewClientFromURL("http://test.com/parsedget?keyvalue=true").
			Get()

		test.NotEqual("What a fat body this is", resp.String())
		test.NotNil(resp.Error)
		test.Equal(0, resp.Status)
	}
}

func TestHeader(t *testing.T) {
	test := assert.New(t)
	defer gock.Off()

	{ // parse url and get with parameters
		gock.New("http://test.com").Get("getheader").MatchHeader("papa", "pipi").Reply(200).BodyString("we test header, why body???")
		resp := rest.
			NewClient("http://test.com/getheader").
			Header("papa", "pipi").
			Get()

		test.Equal("we test header, why body???", resp.String())
		test.Nil(resp.Error)
		test.Equal(200, resp.Status)
	}
	{ // parse url and get with parameters
		gock.New("http://test.com").Get("getheader").MatchHeader("papa", "pipi").Reply(200).BodyString("we test header, why body???")
		resp := rest.
			NewClient("http://test.com/getheader").
			Header("popo", "pipi").
			Get()

		test.NotEqual("we test header, why body???", resp.String())
		test.NotNil(resp.Error)
		test.Equal(0, resp.Status)
	}
}

func TestCatchRequest(t *testing.T) {
	test := assert.New(t)

	rest.CatchRequest = func(req *http.Request) *http.Response {
		rec := httptest.NewRecorder()
		rec.WriteString("round like a record baby")
		return rec.Result()
	}
	resp := rest.
		NewClient("http://test.com/get").
		Get()

	test.Equal("round like a record baby", resp.String())
	test.Nil(resp.Error)
	test.Equal(200, resp.Status)

}
