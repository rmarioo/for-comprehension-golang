# Gomprehension

**Gomprehension** library simulate a [for comprehension](https://edward-huang.com/scala/functional-programming/2021/11/30/why-do-functional-programmers-prefer-for-comprehension-over-imperative-code-block/) and monadic function compositions in golang like the one in other languages like scala or rust

Example you can write this coincise code  and leverage for comprehension.
This is needed in complex cases when function need to be composed and a function need to access to more than on previous results
```golang
var ReadToFile = wrapToPanic(ioutil.ReadFile)
var BytesToStrFunction = wrapToPanic(MyBytesToStr)
var Base64DecodeStringFunction = wrapToPanic(MyBase64DecodeString)
var JSONUnmarshalFunction = wrapToPanic(MyJSONUnmarshal)


func TestReadDecodeAndUnMarshal(t *testing.T) {

r, err := panicToPair(func () info {

    bytes := ReadToFile("testdata/test.base64")

    str := BytesToStrFunction(bytes)

    decodeString := Base64DecodeStringFunction(str)

    return JSONUnmarshalFunction(decodeString)

})

assert.Nil(t, err, "not expecting an error")
assert.Equal(t, info{Test: "monad"}, r)

}
```

instead of this usual verbose way to check each time the error.
NOTE:  In following code error handling code is more that business code

```golang
func TestReadDecodeAndUnMarshal_NoForComphrenstion(t *testing.T) {

bytez, e1 := ioutil.ReadFile("testdata/test.base64")

if e1 != nil {
    t.Fatal("failed reading file");
}

str, e2 := MyBytesToStr(bytez)

if e2 != nil {
    t.Fatal("converting file to string");
}

decodeString, e3 := MyBase64DecodeString(str)

if e3 != nil {
    t.Fatal("decoding from base 64");
}

result, err := MyJSONUnmarshal(decodeString)

assert.Nil(t, err, "not expecting an error")
assert.Equal(t, info{Test: "monad"}, result)

}
```

if there is no need to for comphrension a simple composition can be used 

```golang
func TestSimpleComposition(t *testing.T) {

	r, err := compose4(
		ioutil.ReadFile,
		MyBytesToStr,
		MyBase64DecodeString,
		MyJSONUnmarshal)("testdata/test.base64")

	assert.Nil(t, err, "not expecting an error")
	assert.Equal(t, info{Test: "monad"}, r)
}
```


