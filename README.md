# Gomprehension

**Gomprehension** library simulate a for comprehension and monadic function compositions in golang like the one in other languages like scala or rust

Example you can write this coincise code  
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

instead of this one

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


