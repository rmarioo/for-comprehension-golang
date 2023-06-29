package errm

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestSimpleComposition(t *testing.T) {

	r, err := compose4(
		ioutil.ReadFile,
		MyBytesToStr,
		MyBase64DecodeString,
		MyJSONUnmarshal)("testdata/test.base64")

	assert.Nil(t, err, "not expecting an error")
	assert.Equal(t, info{Test: "monad"}, r)
}

type Any interface{}

var ReadToFile = wrapToPanic(ioutil.ReadFile)
var BytesToStrFunction = wrapToPanic(MyBytesToStr)
var Base64DecodeStringFunction = wrapToPanic(MyBase64DecodeString)
var JSONUnmarshalFunction = wrapToPanic(MyJSONUnmarshal)

func TestReadDecodeAndUnMarshal_ForComphrenstion(t *testing.T) {

	r, err := panicToPair(func() info {

		bytes := ReadToFile("testdata/test.base64")
		str := BytesToStrFunction(bytes)

		decodeString := Base64DecodeStringFunction(str)

		return JSONUnmarshalFunction(decodeString)

	})

	assert.Nil(t, err, "not expecting an error")
	assert.Equal(t, info{Test: "monad"}, r)

}

func TestReadDecodeAndUnMarshal_NoForComphrenstion(t *testing.T) {

	bytez, e1 := ioutil.ReadFile("testdata/test.base64")

	if e1 != nil {
		t.Fatal("failed reading file")
	}

	str, e2 := MyBytesToStr(bytez)

	if e2 != nil {
		t.Fatal("converting file to string")
	}

	decodeString, e3 := MyBase64DecodeString(str)

	if e3 != nil {
		t.Fatal("decoding from base 64")
	}

	result, err := MyJSONUnmarshal(decodeString)

	assert.Nil(t, err, "not expecting an error")
	assert.Equal(t, info{Test: "monad"}, result)

}

func TestInvalidBase64ReadDecodeAndUnMarshal(t *testing.T) {

	_, err := panicToPair(func() info {

		bytes := ReadToFile("testdata/invalid.base64")
		str := BytesToStrFunction(bytes)

		decodeString := Base64DecodeStringFunction(str)

		return JSONUnmarshalFunction(decodeString)

	})

	fmt.Println("got error " + toString(err))

	assert.NotNil(t, err, "expecting an error")

}

func TestBaseForComprehension(t *testing.T) {

	_, err := panicToPair(func() int {

		i1 := panicIfError(strconv.Atoi("100"))
		i2 := panicIfError(strconv.Atoi(intToString(i1) + "_NOT_STRING"))
		i3 := panicIfError(strconv.Atoi("50"))

		return i1 + i2 + i3

	})

	assert.NotNil(t, err, "expectin an error ")

}

func TestExecutionStopsAtFirstError(t *testing.T) {

	Atoi := wrapToPanic(strconv.Atoi)
	_, err := panicToPair(func() int {

		i1 := Atoi("100")
		i2 := Atoi(intToString(i1) + "_NOT_STRING")
		i3 := Atoi("50")

		return i1 + i2 + i3

	})

	assert.NotNil(t, err, "expecting an error")

}

func TestExecutionSuccess(t *testing.T) {

	Atoi := wrapToPanic(strconv.Atoi)
	r, err := panicToPair(func() int {

		i1 := Atoi("100")
		i2 := Atoi(intToString(i1))
		i3 := Atoi("50")

		return i1 + i2 + i3

	})

	assert.Nil(t, err, "not expecting error")

	assert.Equal(t, 100+100+50, r)

}

func MyBytesToStr(v []byte) (string, error) {
	vBytes := v

	return string(vBytes), nil

}

func MyBase64DecodeString(v string) ([]byte, error) {
	vString := v
	return base64.StdEncoding.DecodeString(vString)

}

type info struct {
	Test string
}

func MyJSONUnmarshal(v []byte) (info, error) {
	vBytes := v
	var i info
	err := json.Unmarshal(vBytes, &i)
	return i, err

}

func intToString(i1 int) string {
	return fmt.Sprintf("%d", i1)
}
