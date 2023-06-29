package errm

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"
	"testing"
)

var ReadToFile = wrapToPanic(ioutil.ReadFile)
var BytesToStrFunction = wrapToPanic(MyBytesToStr)
var Base64DecodeStringFunction = wrapToPanic(MyBase64DecodeString)
var JSONUnmarshalFunction = wrapToPanic(MyJSONUnmarshal)

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

func TestReadDecodeAndUnMarshal(t *testing.T) {

	r, err := panicToPair(func() info {

		bytes := ReadToFile("testdata/test.base64")
		str := BytesToStrFunction(bytes)

		decodeString := Base64DecodeStringFunction(str)

		return JSONUnmarshalFunction(decodeString)

	})

	assert.Nil(t, err, "not expecting an error")
	assert.Equal(t, info{Test: "monad"}, r)

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

func TestReadFileBase64JSON(t *testing.T) {

	ma := Return("testdata/test.base64")
	m := Bind(ma, ReadFile)
	m = Bind(m, BytesToStr)
	m = Bind(m, Base64DecodeString)
	m = Bind(m, JSONUnmarshal)
	jsonMap, err := m(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(jsonMap)
}

type FuncErr[A, B any] func(A) (B, error)

func compose2[A any, B any, C any](f1 FuncErr[A, B], f2 FuncErr[B, C]) FuncErr[A, C] {
	return func(a A) (C, error) {
		b, err := f1(a)
		if err != nil {
			var result C
			return result, err
		} else {
			return f2(b)
		}

	}
}

func wrapToPanic[A any, B any](f1 FuncErr[A, B]) func(a A) B {
	return func(a A) B {
		b, err := f1(a)
		if err != nil {
			panic(err)
		} else {
			return b
		}
	}

}

func compose3[A any, B any, C any, D any](f1 FuncErr[A, B], f2 FuncErr[B, C], f3 FuncErr[C, D]) FuncErr[A, D] {
	two := compose2(f1, f2)
	return compose2(two, f3)
}

func compose4[A any, B any, C any, D any, E any](
	f1 FuncErr[A, B],
	f2 FuncErr[B, C],
	f3 FuncErr[C, D],
	f4 FuncErr[D, E]) FuncErr[A, E] {
	two := compose3(f1, f2, f3)
	return compose2(two, f4)
}

func compose5[A any, B any, C any, D any, E any, F any](
	f1 FuncErr[A, B],
	f2 FuncErr[B, C],
	f3 FuncErr[C, D],
	f4 FuncErr[D, E],
	f5 FuncErr[E, F]) FuncErr[A, F] {
	two := compose4(f1, f2, f3, f4)
	return compose2(two, f5)
}

func TestBaseForComprehension(t *testing.T) {

	_, err := panicToPair(func() int {

		i1 := panicIfError(strconv.Atoi("100"))
		i2 := panicIfError(strconv.Atoi(intToString(i1) + "_NOT_STRING"))
		i3 := panicIfError(strconv.Atoi("50"))

		return i1 + i2 + i3

	})

	//fmt.Println(r)

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

func MyBytesToStr(v Any) (string, error) {
	vBytes := v.([]byte)

	return string(vBytes), nil

}

func MyBase64DecodeString(v Any) ([]byte, error) {
	vString := v.(string)
	return base64.StdEncoding.DecodeString(vString)

}

type info struct {
	Test string
}

func MyJSONUnmarshal(v Any) (info, error) {
	vBytes := v.([]byte)
	var i info
	err := json.Unmarshal(vBytes, &i)
	return i, err

}

func intToString(i1 int) string {
	return fmt.Sprintf("%d", i1)
}
func panicToPair[T any](f func() T) (r T, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error happened " + toString(r))
		}
	}()

	i := f()

	return i, nil
}

func toString(r any) string {
	return fmt.Sprintf("%v", r)
}

func mytest() (r int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error happened " + toString(r))
		}
	}()

	s1 := aStringFunction("10", false)
	s2 := aStringFunction(s1+"_43234", false)
	i := toIntFunction(s2)

	return i, nil
}
