package strutils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"io"
	"math"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

// referrer to https://golang.org/pkg/regexp/syntax/
var numericPattern = regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?$`)
var tagElementsPattern = regexp.MustCompile(`(?ims)(?P<tag><(/*\s*|\?*|\!*)(figcaption|expression|blockquote|plaintext|textarea|progress|optgroup|noscript|noframes|menuitem|frameset|fieldset|!DOCTYPE|datalist|colgroup|behavior|basefont|summary|section|isindex|details|caption|bgsound|article|address|acronym|strong|strike|source|select|script|output|option|object|legend|keygen|ilayer|iframe|header|footer|figure|dialog|center|canvas|button|applet|video|track|title|thead|tfoot|tbody|table|style|small|param|meter|layer|label|input|frame|embed|blink|audio|aside|alert|time|span|samp|ruby|meta|menu|mark|main|link|html|head|form|font|code|cite|body|base|area|abbr|xss|xml|wbr|var|svg|sup|sub|pre|nav|map|kbd|ins|img|div|dir|dfn|del|col|big|bdo|bdi|!--|ul|tt|tr|th|td|rt|rp|ol|li|hr|em|dt|dl|dd|br|u|s|q|p|i|b|a|(h[0-9]+))([^><]*)([><]*))`)
var whiteSpacePattern = regexp.MustCompile(`(?im)\s{2,}`)
var entityEncodedPattern = regexp.MustCompile(`(?ims)(&(?:[a-z0-9]{2,8}|#[0-9]{2,3});)`)
var urlEncodedPattern = regexp.MustCompile(`(?ims)(%[a-zA-Z0-9]{2})`)

// for debug
//var detectUnicodeEntities = regexp.MustCompile(`(?ims)u([0-9a-z]{4})`)

// StringProc is String processing methods, All operations on this object
type StringProc struct {
	sync.RWMutex
}

// NewStringProc Creates and returns a String processing methods's pointer.
func NewStringProc() *StringProc {
	return &StringProc{}
}

// AddSlashes is quote string with slashes
func (s *StringProc) AddSlashes(str string) string {

	l := len(str)

	buf := make([]byte, 0, l*2) //prealloca

	for i := 0; i < l; i++ {

		buf = append(buf, str[i])

		switch str[i] {

		case 92: //Dec : /

			if l >= i+1 {
				buf = append(buf, 92)

				if l > i+1 && str[i+1] == 92 {
					i++
				}
			}
		}
	}

	return string(buf)
}

// StripSlashes is Un-quotes a quoted string
func (s *StringProc) StripSlashes(str string) string {

	l := len(str)
	buf := make([]byte, 0, l) //prealloca

	for i := 0; i < l; i++ {

		buf = append(buf, str[i])
		if l > i+1 && str[i+1] == 92 {
			i++
		}
	}

	return string(buf)

}

// Nl2Br is breakstr inserted before looks like space (CRLF , LFCR, SPACE, NL)
func (s *StringProc) Nl2Br(str string) string {

	// BenchmarkNl2Br-8                   	10000000	      3398 ns/op
	// BenchmarkNl2BrUseStringReplace-8   	10000000	      4535 ns/op
	brtag := []byte("<br />")
	l := len(str)
	buf := make([]byte, 0, l) //prealloca

	for i := 0; i < l; i++ {

		switch str[i] {

		case 10, 13: //NL or CR

			buf = append(buf, brtag...)

			if l >= i+1 {
				if l > i+1 && (str[i+1] == 10 || str[i+1] == 13) { //NL+CR or CR+NL
					i++
				}
			}
		default:
			buf = append(buf, str[i])
		}
	}

	return string(buf)
}

// Br2Nl is replaces HTML line breaks to a newline
func (s *StringProc) Br2Nl(str string) string {

	// <br> , <br /> , <br/>
	// <BR> , <BR /> , <BR/>
	nlchar := []byte("\n")

	l := len(str)
	buf := make([]byte, 0, l) //prealloca

	for i := 0; i < l; i++ {

		switch str[i] {

		case 60: //<

			if l >= i+3 {

				/*
					b = 98
					B = 66
					r = 82
					R = 114
					SPACE = 32
					/ = 47
					> = 62
				*/

				if l >= i+3 && ((str[i+1] == 98 || str[i+1] == 66) && (str[i+2] == 82 || str[i+2] == 114) && str[i+3] == 62) { // <br> || <BR>
					buf = append(buf, nlchar...)
					i += 3
					continue
				}

				if l >= i+4 && ((str[i+1] == 98 || str[i+1] == 66) && (str[i+2] == 82 || str[i+2] == 114) && str[i+3] == 47 && str[i+4] == 62) { // <br/> || <BR/>
					buf = append(buf, nlchar...)
					i += 4
					continue
				}

				if l >= i+5 && ((str[i+1] == 98 || str[i+1] == 66) && (str[i+2] == 82 || str[i+2] == 114) && str[i+3] == 32 && str[i+4] == 47 && str[i+5] == 62) { // <br /> || <BR />
					buf = append(buf, nlchar...)
					i += 5
					continue
				}
			}
			fallthrough

		default:
			buf = append(buf, str[i])
		}
	}

	return string(buf)
}

// WordWrapSimple is Wraps a string to a given number of characters using break characters (TAB, SPACE)
func (s *StringProc) WordWrapSimple(str string, wd int, breakstr string) (string, error) {

	if wd < 1 {
		err := errors.New("wd At least 1 or More")
		return str, err
	}

	strl := len(str)
	breakstrl := len(breakstr)

	buf := make([]byte, 0, (strl+breakstrl)*2)
	bufstr := []byte(str)

	brpos := 0
	for _, v := range bufstr {

		if (v == 9 || v == 32) && brpos >= wd {
			buf = append(buf, []byte(breakstr)...)
			brpos = -1

		} else {
			buf = append(buf, v)
		}
		brpos++
	}

	return string(buf), nil
}

// WordWrapAround is Wraps a string to a given number of characters using break characters (TAB, SPACE)
func (s *StringProc) WordWrapAround(str string, wd int, breakstr string) (string, error) {

	if wd < 1 {
		return "", errors.New("wd At least 1 or More")
	}

	strl := len(str)
	breakstrl := len(breakstr)

	buf := make([]byte, 0, (strl+breakstrl)*2)
	bufstr := []byte(str)

	lastspc := make([]int, 0, strl)
	strpos := 0

	//looking for space or tab
	for _, v := range bufstr {

		if v == 9 || v == 32 {
			lastspc = append(lastspc, strpos)
		}
		strpos++
	}

	inject := make([]int, 0, strl)

	//looking for break point
	beforeBp := 0
	width := wd

	for _, v := range lastspc {

		if beforeBp != v {
			beforeBp = v
		}

		// DEBUG: fmt.Printf("V(%v) (%d <= %d || %d <= %d || %d <= %d) && %d <= %d : ", v, width, beforeBp, width, beforeBp+1, width, beforeBp-1, width, v)
		if (width <= beforeBp || width <= beforeBp+1 || width <= beforeBp-1) && width <= v {
			inject = append(inject, beforeBp)
			width += wd
		} else if width < v && len(lastspc) == 1 {
			inject = append(inject, v)
		}
		//fmt.Println()
	}

	//injection
	breakno := 0
	loopcnt := 0
	injectcnt := len(inject)
	for _, v := range bufstr {

		//fmt.Printf("(%v) %d > %d && %d <= %d\n", v, injectcnt, breakno, inject[breakno], loopcnt)
		if injectcnt > breakno && inject[breakno] == loopcnt {

			buf = append(buf, []byte(breakstr)...)
			if injectcnt > breakno+1 {
				breakno++
			}
		} else {
			buf = append(buf, v)
		}

		loopcnt++
	}

	return string(buf), nil
}

func (s *StringProc) numberToString(obj interface{}) (string, error) {

	var strNum string

	switch obj.(type) {

	case string:
		strNum = obj.(string)
		if !numericPattern.MatchString(strNum) {
			return strNum, fmt.Errorf("Not Support obj.(%v) := %v ", reflect.TypeOf(obj), strNum)
		}
	case int:
		strNum = strconv.FormatInt(int64(obj.(int)), 10)
	case int8:
		strNum = strconv.FormatInt(int64(obj.(int8)), 10)
	case int16:
		strNum = strconv.FormatInt(int64(obj.(int16)), 10)
	case int32:
		strNum = strconv.FormatInt(int64(obj.(int32)), 10)
	case int64:
		strNum = strconv.FormatInt(obj.(int64), 10)
	case uint:
		strNum = strconv.FormatUint(uint64(obj.(uint)), 10)
	case uint8:
		strNum = strconv.FormatUint(uint64(obj.(uint8)), 10)
	case uint16:
		strNum = strconv.FormatUint(uint64(obj.(uint16)), 10)
	case uint32:
		strNum = strconv.FormatUint(uint64(obj.(uint32)), 10)
	case uint64:
		strNum = strconv.FormatUint(obj.(uint64), 10)
	case float32:
		strNum = fmt.Sprintf("%g", obj.(float32))
	case float64:
		strNum = fmt.Sprintf("%g", obj.(float64))
	case complex64:
		strNum = fmt.Sprintf("%g", obj.(complex64))
	case complex128:
		strNum = fmt.Sprintf("%g", obj.(complex128))

	default:
		return strNum, fmt.Errorf("Not Support obj.(%v)", reflect.TypeOf(obj))
	}

	return strNum, nil
}

// NumberFmt is format a number with english notation grouped thousands
// TODO : support other country notation
func (s *StringProc) NumberFmt(obj interface{}) (string, error) {

	//check : complex
	switch obj.(type) {
	case complex64, complex128:
		return "", fmt.Errorf("Not Support obj.(%v)", reflect.TypeOf(obj))
	}

	strNum, err := s.numberToString(obj)
	if err != nil {
		return "", err
	}

	bufbyteStr := []byte(strNum)
	bufbyteStrLen := len(bufbyteStr)

	//subffix after dot
	bufbyteTail := make([]byte, bufbyteStrLen-1)

	//init.
	var foundDot, foundPos, dotcnt, bufbyteSize int

	//looking for dot
	for i := bufbyteStrLen - 1; i >= 0; i-- {
		if bufbyteStr[i] == 46 {
			copy(bufbyteTail, bufbyteStr[i:])
			foundDot = i
			foundPos = i
			break
		}
	}

	//make bufbyte size
	if foundDot == 0 { //numeric without dot
		bufbyteSize = int(math.Ceil(float64(bufbyteStrLen) + float64(bufbyteStrLen)/3))
		foundDot = bufbyteStrLen
		foundPos = bufbyteSize - 2

		bufbyteSize--

	} else { //with dot

		var calFoundDot int

		if bufbyteStr[0] == 45 { //if startwith "-"(45)
			calFoundDot = foundDot - 1
		} else {
			calFoundDot = foundDot
		}

		bufbyteSize = int(math.Ceil(float64(calFoundDot) + float64(calFoundDot)/3 + float64(bufbyteStrLen-calFoundDot) - 1))
	}

	//make a buffer byte
	bufbyte := make([]byte, bufbyteSize)

	//skip : need to dot injection
	if 4 > foundDot {
		return strNum, nil
	}

	//injection
	intoPos := foundPos
	for i := foundDot - 1; i >= 0; i-- {
		if dotcnt >= 3 && (bufbyteStr[i] >= 48 && bufbyteStr[i] <= 57 || bufbyteStr[i] == 69 || bufbyteStr[i] == 101 || bufbyteStr[i] == 43) {
			bufbyte[intoPos] = 44
			intoPos--
			dotcnt = 0
		}
		bufbyte[intoPos] = bufbyteStr[i]
		intoPos--
		dotcnt++
	}

	//into dot to tail
	intoPos = foundPos + 1
	if foundDot != bufbyteStrLen {
		for _, v := range bufbyteTail {
			if v == 0 { //NULL
				break
			}

			bufbyte[intoPos] = v
			intoPos++
		}
	}

	return string(bufbyte), nil
}

// padding contol const
const (
	PadLeft  = 0 //left padding
	PadRight = 1 //right padding
	PadBoth  = 2 //both padding
)

// PaddingBoth is Padding method alias with PadBoth Option
func (s *StringProc) PaddingBoth(str string, fill string, mx int) string {
	return s.Padding(str, fill, PadBoth, mx)
}

// PaddingLeft is Padding method alias with PadRight Option
func (s *StringProc) PaddingLeft(str string, fill string, mx int) string {
	return s.Padding(str, fill, PadLeft, mx)
}

// PaddingRight is Padding method alias with PadRight Option
func (s *StringProc) PaddingRight(str string, fill string, mx int) string {
	return s.Padding(str, fill, PadRight, mx)
}

// Padding is Pad a string to a certain length with another string
// BenchmarkPadding-8                   10000000	       271 ns/op
// BenchmarkPaddingUseStringRepeat-8   	 3000000	       418 ns/op
func (s *StringProc) Padding(str string, fill string, m int, mx int) string {

	byteStr := []byte(str)
	byteStrLen := len(byteStr)
	if byteStrLen >= mx || mx < 1 {
		return str
	}

	var leftsize int
	var rightsize int

	switch m {
	case PadBoth:
		rlsize := float64(mx-byteStrLen) / 2
		leftsize = int(rlsize)
		rightsize = int(rlsize + math.Copysign(0.5, rlsize))

	case PadLeft:
		leftsize = mx - byteStrLen

	case PadRight:
		rightsize = mx - byteStrLen

	}

	buf := make([]byte, 0, mx)

	if m == PadLeft || m == PadBoth {
		for i := 0; i < leftsize; {
			for _, v := range []byte(fill) {
				buf = append(buf, v)
				if i >= leftsize-1 {
					i = leftsize
					break
				} else {
					i++
				}
			}
		}
	}

	buf = append(buf, byteStr...)

	if m == PadRight || m == PadBoth {
		for i := 0; i < rightsize; {
			for _, v := range []byte(fill) {
				buf = append(buf, v)
				if i >= rightsize-1 {
					i = rightsize
					break
				} else {
					i++
				}
			}
		}
	}

	return string(buf)
}

// LowerCaseFirstWords is Lowercase the first character of each word in a string
// INFO : (Support Token Are \t(9)\r(13)\n(10)\f(12)\v(11)\s(32))
func (s *StringProc) LowerCaseFirstWords(str string) string {

	upper := 1
	bufbyteStr := []byte(str)
	retval := make([]byte, len(bufbyteStr))
	for k, v := range bufbyteStr {

		if upper == 1 && v >= 65 && v <= 90 {
			v = v + 32
		}

		upper = 0

		if v >= 9 && v <= 13 || v == 32 {
			upper = 1
		}
		retval[k] = v
	}

	return string(retval)
}

// UpperCaseFirstWords is Uppercase the first character of each word in a string
// INFO : (Support Token Are \t(9)\r(13)\n(10)\f(12)\v(11)\s(32))
func (s *StringProc) UpperCaseFirstWords(str string) string {

	upper := 1
	bufbyteStr := []byte(str)
	retval := make([]byte, len(bufbyteStr))
	for k, v := range bufbyteStr {

		if upper == 1 && v >= 97 && v <= 122 {
			v = v - 32
		}

		upper = 0

		if v >= 9 && v <= 13 || v == 32 {
			upper = 1
		}
		retval[k] = v
	}

	return string(retval)
}

// SwapCaseFirstWords is Switch the first character case of each word in a string
func (s *StringProc) SwapCaseFirstWords(str string) string {

	upper := 1
	bufbyteStr := []byte(str)
	retval := make([]byte, len(bufbyteStr))
	for k, v := range bufbyteStr {

		switch {
		case upper == 1 && v >= 65 && v <= 90:
			v = v + 32

		case upper == 1 && v >= 97 && v <= 122:
			v = v - 32
		}

		upper = 0

		if v >= 9 && v <= 13 || v == 32 {
			upper = 1
		}
		retval[k] = v
	}

	return string(retval)
}

// Unit type control
const (
	_               = uint8(iota)
	LowerCaseSingle // Single Unit character converted to Lower-case
	LowerCaseDouble // Double Unit characters converted to Lower-case

	UpperCaseSingle // Single Unit character converted to Uppper-case
	UpperCaseDouble // Double Unit characters converted to Upper-case

	CamelCaseDouble // Double Unit characters converted to Camel-case
	CamelCaseLong   // Full Unit characters converted to Camel-case
)

var sizeStrLowerCaseSingle = []string{"b", "k", "m", "g", "t", "p", "e", "z", "y"}
var sizeStrLowerCaseDouble = []string{"b", "kb", "mb", "gb", "tb", "pb", "eb", "zb", "yb"}
var sizeStrUpperCaseSingle = []string{"B", "K", "M", "G", "T", "P", "E", "Z", "Y"}
var sizeStrUpperCaseDouble = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
var sizeStrCamelCaseDouble = []string{"B", "Kb", "Mb", "Gb", "Tb", "Eb", "Zb", "Yb"}
var sizeStrCamelCaseLong = []string{"Byte", "KiloByte", "MegaByte", "GigaByte", "TeraByte", "ExaByte", "ZettaByte", "YottaByte"}

//HumanByteSize is Byte Size convert to Easy Readable Size String
func (s *StringProc) HumanByteSize(obj interface{}, decimals int, unit uint8) (string, error) {

	if unit < LowerCaseSingle || unit > CamelCaseLong {
		return "", fmt.Errorf("Not allow unit parameter : %v", unit)
	}

	strNum, err := s.numberToString(obj)
	if err != nil {
		return "", err
	}

	var bufStrFloat64 float64

	switch obj.(type) {
	case string:
		bufStrFloat64, err = strconv.ParseFloat(strNum, 64)
		if err != nil {
			return "", fmt.Errorf("Not Support %v (obj.(%v))", obj, reflect.TypeOf(obj))
		}

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32:

		float64Type := reflect.TypeOf(float64(0))
		tmpVal := reflect.Indirect(reflect.ValueOf(obj))

		/*
			//impossible?
			if tmpVal.Type().ConvertibleTo(float64Type) == false {
				bufStrFloat64, err = strconv.ParseFloat(reflect.ValueOf(obj).String(), 64)
				if err != nil {
					return "", fmt.Errorf("Not Support %v (obj.(%v))", obj, reflect.TypeOf(obj))
				}

			} else {
				bufStrFloat64 = tmpVal.Convert(float64Type).Float()
			}
		*/

		bufStrFloat64 = tmpVal.Convert(float64Type).Float()

	case float64:
		bufStrFloat64 = obj.(float64)

	default:
		return "", fmt.Errorf("Not Support obj.(%v)", reflect.TypeOf(obj))
	}

	var sizeStr []string

	switch unit {
	case LowerCaseSingle:
		sizeStr = sizeStrLowerCaseSingle
	case LowerCaseDouble:
		sizeStr = sizeStrLowerCaseDouble
	case UpperCaseSingle:
		sizeStr = sizeStrUpperCaseSingle
	case UpperCaseDouble:
		sizeStr = sizeStrUpperCaseDouble
	case CamelCaseDouble:
		sizeStr = sizeStrCamelCaseDouble
	case CamelCaseLong:
		sizeStr = sizeStrCamelCaseLong
	}

	strNumLen := len(strNum)

	factor := int(math.Floor(float64(strNumLen)-1) / 3)

	decimalsFmt := `%.` + strconv.Itoa(decimals) + `f%s`
	humanSize := bufStrFloat64 / math.Pow(1024, float64(factor))

	var unitStr string
	if len(sizeStr) > factor {
		unitStr = sizeStr[factor]
	} else {
		unitStr = "NaN"
	}

	return fmt.Sprintf(decimalsFmt, humanSize, unitStr), nil
}

//HumanFileSize is File Size convert to Easy Readable Size String
func (s *StringProc) HumanFileSize(filepath string, decimals int, unit uint8) (string, error) {

	fd, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	defer s.closeFd(fd)

	stat, err := fd.Stat() // impossible?. maybe it can be broken fd after file open. anyway can't make a test case..
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	if stat.IsDir() {
		return "", fmt.Errorf("%v isn't file", filepath)
	}

	return s.HumanByteSize(stat.Size(), decimals, unit)
}

// compare with map
var recursiveDepth = 0
var recursiveDepthKeypList = struct {
	sync.RWMutex
	ar []string
}{ar: make([]string, 32)}

func (s *StringProc) compareMap(compObj1 reflect.Value, compObj2 reflect.Value) (bool, error) {

	recursiveDepth++
	var valueCompareErr bool

	for _, k := range compObj1.MapKeys() {

		recursiveDepthKeypList.Lock()
		recursiveDepthKeypList.ar = append(recursiveDepthKeypList.ar, k.String())
		recursiveDepthKeypList.Unlock()

		//check : Type
		if compObj1.MapIndex(k).Kind() != compObj2.MapIndex(k).Kind() {
			return false, fmt.Errorf("Different Type : (obj1[%v] type is `%v`) != (obj2[%v] type is `%v`)", k, compObj1.MapIndex(k).Kind(), k, compObj2.MapIndex(k).Kind())
		}

		switch compObj1.MapIndex(k).Kind() {

		//String
		case reflect.String:
			if compObj1.MapIndex(k).String() != compObj2.MapIndex(k).String() {
				valueCompareErr = true
			}

		//Integer
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if compObj1.MapIndex(k).Int() != compObj2.MapIndex(k).Int() {
				valueCompareErr = true
			}

		//Un-signed Integer
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if compObj1.MapIndex(k).Uint() != compObj2.MapIndex(k).Uint() {
				valueCompareErr = true
			}

		//Float
		case reflect.Float32, reflect.Float64:
			if compObj1.MapIndex(k).Float() != compObj2.MapIndex(k).Float() {
				valueCompareErr = true
			}

		//Boolean
		case reflect.Bool:
			if compObj1.MapIndex(k).Bool() != compObj2.MapIndex(k).Bool() {
				valueCompareErr = true
			}

		//Complex
		case reflect.Complex64, reflect.Complex128:
			if compObj1.MapIndex(k).Complex() != compObj2.MapIndex(k).Complex() {
				valueCompareErr = true
			}

		//Map : recursive loop
		case reflect.Map:
			retval, err := s.compareMap(compObj1.MapIndex(k), compObj2.MapIndex(k))
			if !retval {
				return retval, err
			}

		default:
			return false, fmt.Errorf("Not Support Compare : (obj1[%v] := %v) != (obj2[%v] := %v)", k, compObj1.MapIndex(k), k, compObj2.MapIndex(k))
		}

		if valueCompareErr {
			if recursiveDepth == 1 {
				return false, fmt.Errorf("Different Value : (obj1[%v] := %v) != (obj2[%v] := %v)", k, compObj1.MapIndex(k), k, compObj2.MapIndex(k))
			}

			recursiveDepthKeypList.Lock()
			depthStr := strings.Join(recursiveDepthKeypList.ar, "][")
			recursiveDepthKeypList.Unlock()
			return false, fmt.Errorf("Different Value : (obj1[%v] := %v) != (obj2[%v] := %v)", depthStr, compObj1.MapIndex(k).Interface(), depthStr, compObj2.MapIndex(k))

		}
	}

	return true, nil
}

// AnyCompare is compares two same basic type (without prt) dataset (slice,map,single data).
// TODO : support interface, struct ...
// NOTE : Not safe , Not Test Complete. Require more test data based on the complex dataset.
func (s *StringProc) AnyCompare(obj1 interface{}, obj2 interface{}) (bool, error) {

	compObjVal1 := reflect.ValueOf(obj1)
	compObjVal2 := reflect.ValueOf(obj2)

	compObjType1 := reflect.TypeOf(obj1)
	compObjType2 := reflect.TypeOf(obj2)

	if !compObjVal1.IsValid() || !compObjVal2.IsValid() {
		return false, fmt.Errorf("Invalid, obj1(%v) != obj2(%v)", obj1, obj2)
	}

	if compObjType1.Kind() != compObjType2.Kind() {
		return false, fmt.Errorf("Not Compare type, obj1.(%v) != obj2.(%v)", compObjType1.Kind(), compObjType2.Kind())
	}

	recursiveDepthKeypList.Lock()
	recursiveDepthKeypList.ar = make([]string, 0)
	recursiveDepthKeypList.Unlock()

	switch obj1.(type) {

	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, bool:
		if compObjType1.Comparable() && compObjType2.Comparable() {
			return obj1 == obj2, nil
		}

	default:

		switch {

		case compObjVal1.Kind() == reflect.Slice:

			if compObjVal1.Len() != compObjVal2.Len() {
				return false, fmt.Errorf("Different Size : obj1(%d) != obj2(%d)", compObjVal1.Len(), compObjVal2.Len())
			}

			for i := 0; i < compObjVal1.Len(); i++ {
				if compObjVal1.Index(i).Interface() != compObjVal2.Index(i).Interface() {
					return false, fmt.Errorf("Different Value : (obj1[%d] := %v) != (obj2[%d] := %v)", i, compObjVal1.Index(i).Interface(), i, compObjVal2.Index(i).Interface())
				}
			}

		case compObjVal1.Kind() == reflect.Map:
			if compObjVal1.Len() != compObjVal2.Len() {
				return false, fmt.Errorf("Different Size : obj1(%d) != obj2(%d)", compObjVal1.Len(), compObjVal2.Len())
			}

			recursiveDepth = 0
			retval, err := s.compareMap(compObjVal1, compObjVal2)
			if !retval {
				return retval, err
			}

		default:
			return reflect.DeepEqual(obj1, obj2), nil
			//return false, fmt.Errorf("Not Support Compare : (obj1[%v]) , (obj2[%v])", compObjVal1.Kind(), compObjVal2.Kind())

		}
	}
	return true, nil
}

func (s *StringProc) isHex(c byte) bool {

	if (c >= 48 && c <= 57) || (c >= 65 && c <= 70) || (c >= 97 && c <= 102) { //0~9, a~f, A~F
		return true
	}

	return false
}

func (s *StringProc) unHex(c byte) byte { //from golang. unhex

	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}

	return 0
}

// DecodeUnicodeEntities Decodes Unicode Entities
func (s *StringProc) DecodeUnicodeEntities(val string) (string, error) {

	var tmpret []byte

	l := len(val)
	for i := 0; i < l; i++ {

		if val[i] == 37 && val[i+1] == 117 && l >= i+6 { // % + u

			var tmpval []byte
			tmpval = append(tmpval, val[i+2], val[i+3], val[i+4], val[i+5])

			runeval, err := strconv.ParseInt(string(tmpval), 16, 64)
			if err != nil {
				return "", err
			}

			tmprune := []byte(string(rune(runeval)))
			tmpret = append(tmpret, tmprune...)
			i += 5 //jump %uXXXX

		} else if val[i] == 37 { //control character or other
			tmpret = append(tmpret, s.unHex(val[i+1])<<4|s.unHex(val[i+2]))
			i += 2
		} else {
			tmpret = append(tmpret, val[i])
		}
	}

	return string(tmpret), nil
}

// DecodeURLEncoded Decodes URL-encoded string (including unicode entities)
// NOTE : golang.url.unescape not support unicode entities (%uXXXX)
func (s *StringProc) DecodeURLEncoded(val string) (string, error) {

	var tmpret []byte

	l := len(val)
	for i := 0; i < l; i++ {

		if l <= i+1 { // panic: runtime error: index out of range
			tmpret = append(tmpret, val[i])
			break
		}

		// 37 = %, 117 = u (UnicodeEntity)
		if val[i] == 37 && val[i+1] != 117 && l >= i+3 && s.isHex(val[i+1]) && s.isHex(val[i+2]) {

			tmpret = append(tmpret, s.unHex(val[i+1])<<4|s.unHex(val[i+2]))
			i += 2
			continue
		}

		if val[i] == 37 && val[i+1] == 117 && l >= i+6 { // % + u

			var tmpval []byte
			tmpval = append(tmpval, val[i+2], val[i+3], val[i+4], val[i+5])

			runeval, err := strconv.ParseInt(string(tmpval), 16, 64)
			if err != nil {
				return "", err
			}

			tmprune := []byte(string(rune(runeval)))
			tmpret = append(tmpret, tmprune...)
			i += 5
			continue
		}

		tmpret = append(tmpret, val[i])
	}

	return string(tmpret), nil
}

// StripTags is remove all tag in string
func (s *StringProc) StripTags(str string) (string, error) {

	var retval bool
	notproccnt := 0

	//looking for html entities code in str
ENTITY_DECODE:

	if notproccnt > 3 {
		goto END
	}

	retval = entityEncodedPattern.MatchString(str)
	if retval {
		str = html.UnescapeString(str)
	}

	//looking for html entities code in str
	retval = urlEncodedPattern.MatchString(str)
	if retval {
		tmpstr, err := url.QueryUnescape(str)
		if err == nil {
			if tmpstr == str {
				notproccnt++
			}

			str = tmpstr
			goto ENTITY_DECODE
		} else {

			//url.QueryUnescape not support UnicodeEntities
			tmpstr, err := s.DecodeURLEncoded(str)
			if err == nil {
				if tmpstr == str {
					notproccnt++
				}
				str = tmpstr
				goto ENTITY_DECODE
			} else {
				return str, err
			}
		}
	}
END:

	//remove tag elements
	cleanedStr := tagElementsPattern.ReplaceAllString(str, "")

	//remove multiple whitespace
	cleanedStr = whiteSpacePattern.ReplaceAllString(cleanedStr, "\n")

	return cleanedStr, nil
}

// ConvertToStr is Convert basic data type to string
func (s *StringProc) ConvertToStr(obj interface{}) (string, error) {

	switch obj.(type) {
	case bool:
		if obj.(bool) {
			return "true", nil
		}
		return "false", nil

	default:
		return s.numberToString(obj)
	}
}

// ConvertToArByte returns Convert basic data type to []byte
func (s *StringProc) ConvertToArByte(obj interface{}) ([]byte, error) {

	switch obj.(type) {

	case bool:
		if obj.(bool) {
			return []byte("true"), nil
		}
		return []byte("false"), nil

		/*
			case byte:
				return []byte{obj.(byte)}, nil
		*/

	case []uint8:
		return reflect.ValueOf(obj).Bytes(), nil

	case string:
		return []byte(obj.(string)), nil

	case int:
		return []byte(strconv.FormatInt(int64(obj.(int)), 10)), nil
	case int8:
		return []byte(strconv.FormatInt(int64(obj.(int8)), 10)), nil
	case int16:
		return []byte(strconv.FormatInt(int64(obj.(int16)), 10)), nil
	case int32:
		return []byte(strconv.FormatInt(int64(obj.(int32)), 10)), nil
	case int64:
		return []byte(strconv.FormatInt(obj.(int64), 10)), nil
	case uint:
		return []byte(strconv.FormatUint(uint64(obj.(uint)), 10)), nil
	case uint8: //same byte
		return []byte(strconv.FormatUint(uint64(obj.(uint8)), 10)), nil
	case uint16:
		return []byte(strconv.FormatUint(uint64(obj.(uint16)), 10)), nil
	case uint32:
		return []byte(strconv.FormatUint(uint64(obj.(uint32)), 10)), nil
	case uint64:
		return []byte(strconv.FormatUint(obj.(uint64), 10)), nil
	case float32:
		return []byte(fmt.Sprintf("%g", obj.(float32))), nil
	case float64:
		return []byte(fmt.Sprintf("%g", obj.(float64))), nil
	case complex64:
		return []byte(fmt.Sprintf("%g", obj.(complex64))), nil
	case complex128:
		return []byte(fmt.Sprintf("%g", obj.(complex128))), nil

	default:
		return nil, fmt.Errorf("not support type(%s)", reflect.TypeOf(obj).String())
	}
}

// ReverseStr is Reverse a String , According to value type between ascii or rune
// TODO : improve performance (use goroutin)
func (s *StringProc) ReverseStr(str string) string {
	/*
	   data : "0123456789" * 100
	   BenchmarkReverseStr-8              	   50000	     34127 ns/op	    5120 B/op	       2 allocs/op
	   BenchmarkReverseNormalStr-8        	 1000000	      1187 ns/op	    2048 B/op	       2 allocs/op
	   BenchmarkReverseReverseUnicode-8   	  100000	     29343 ns/op	    5120 B/op	       2 allocs/op
	*/

	if len(str) != utf8.RuneCountInString(str) {
		return s.ReverseUnicode(str)
	}

	return s.ReverseNormalStr(str)
}

// ReverseNormalStr is Reverse a None-unicode String
func (s *StringProc) ReverseNormalStr(str string) string {

	bufbyteStr := []byte(str)
	bufbyteStrLen := len(bufbyteStr)
	swapSize := int(math.Ceil(float64(bufbyteStrLen) / 2))

	headNo := 0
	tailNo := bufbyteStrLen - 1
	for i := 0; i < swapSize; i++ {
		bufbyteStr[tailNo], bufbyteStr[headNo] = bufbyteStr[headNo], bufbyteStr[tailNo]
		headNo++
		tailNo--
	}

	return string(bufbyteStr[:])
}

// ReverseUnicode is Reverse a unicode String
func (s *StringProc) ReverseUnicode(str string) string {

	bufRuneStr := []rune(str)
	bufRuneStrl := len(bufRuneStr)
	swapSize := int(math.Ceil(float64(bufRuneStrl) / 2))

	headNo := 0
	tailNo := bufRuneStrl - 1
	for i := 0; i < swapSize; i++ {
		bufRuneStr[tailNo], bufRuneStr[headNo] = bufRuneStr[headNo], bufRuneStr[tailNo]
		headNo++
		tailNo--
	}

	return string(bufRuneStr[:])
}

// FileMD5Hash is MD5 checksum of the file
func (s *StringProc) FileMD5Hash(filepath string) (string, error) {

	fd, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	defer s.closeFd(fd)

	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, fd); err != nil {
		return "", err
	}

	return hex.EncodeToString(md5Hash.Sum(nil)), nil
}

// MD5Hash is MD5 checksum of the string
func (s *StringProc) MD5Hash(str string) (string, error) {

	md5Hash := md5.New()
	if _, err := io.WriteString(md5Hash, str); err != nil {
		return "", err
	}

	return hex.EncodeToString(md5Hash.Sum(nil)), nil
}

func (s *StringProc) closeFd(fd *os.File) {

	err := fd.Close()
	if err != nil {
		fmt.Printf("Error : %+v\n", err)
	}

}

// RegExpNamedGroups is Captures the text matched by regex into the group name
// NOTE : Not Support the Multiple Groups with The Same Name
func (s *StringProc) RegExpNamedGroups(regex *regexp.Regexp, val string) (map[string]string, error) {

	ok := false
	err := errors.New("not all success patterns were matched")

	retval := map[string]string{}
	extractSubExpNames := regex.SubexpNames()

	ret := regex.FindStringSubmatch(val)
	if len(ret) > 0 {
		for no, val := range ret {
			if no != 0 && val != "" {
				if extractSubExpNames[no] != "" {
					retval[extractSubExpNames[no]] = val
					ok = true
				}
			}
		}
	}

	if ok {
		err = nil
	}

	return retval, err
}
