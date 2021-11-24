package strutils

import (
	"errors"
	"html"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var emailPattern = regexp.MustCompile("^[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?$")
var domainPattern = regexp.MustCompile(`^(([a-zA-Z0-9-\p{L}]{1,63}\.)?(xn--)?[a-zA-Z0-9\p{L}]+(-[a-zA-Z0-9\p{L}]+)*\.)+[a-zA-Z\p{L}]{2,63}$`)
var urlPattern = regexp.MustCompile(`^((((https?|ftps?|gopher|telnet|nntp)://)|(mailto:|news:))(%[0-9A-Fa-f]{2}|[-()_.!~*';/?:@#&=+$,A-Za-z0-9\p{L}])+)([).!';/?:,][[:blank:]])?$`)

// StringValidator is String processing methods, All operations on this object
type StringValidator struct{}

// NewStringValidator is Creates and returns a String processing methods's pointer.
func NewStringValidator() *StringValidator {
	return &StringValidator{}
}

// IsValidEmail is Validates whether the value is a valid e-mail address.
func (s *StringValidator) IsValidEmail(str string) bool {
	return emailPattern.MatchString(str)
}

// IsValidDomain is Validates whether the value is a valid domain address
func (s *StringValidator) IsValidDomain(str string) bool {
	return domainPattern.MatchString(str)
}

// IsValidURL is Validates whether the value is a valid url
func (s *StringValidator) IsValidURL(str string) bool {
	return urlPattern.MatchString(str)
}

// IsValidMACAddr is Validates whether the value is a valid h/w mac address
func (s *StringValidator) IsValidMACAddr(str string) bool {
	if _, err := net.ParseMAC(str); err == nil {
		return true
	}

	return false
}

/*
This consts for cktypes of IsValidIPAddr
IPAny              // Any IP Address Type
IPv4               // IPv4 (32 chars)
IPv6               // IPv6(39 chars)
IPv4MappedIPv6     // IP4-mapped IPv6 (45 chars) , Ex) ::FFFF:129.144.52.38
IPv4CIDR           // IPv4 + CIDR
IPv6CIDR           // IPv6 + CIDR
IPv4MappedIPv6CIDR //IPv4-mapped IPv6 + CIRD
*/
const (
	none               = 0
	IPAny              = 1
	IPv4               = 32
	IPv6               = 39
	IPv4MappedIPv6     = 45
	IPv4CIDR           = IPv4 + 3
	IPv6CIDR           = IPv6 + 3
	IPv4MappedIPv6CIDR = IPv4MappedIPv6 + 3
)

// IsValidIPAddr is Validates whether the value to be exactly a given validation type (IPv4, IPv6, IPv4MappedIPv6, IPv4CIDR, IPv6CIDR, IPv4MappedIPv6CIDR OR IPAny)
func (s *StringValidator) IsValidIPAddr(str string, cktypes ...int) (bool, error) {

	for _, cktype := range cktypes {

		if cktype != IPAny && cktype != IPv4 && cktype != IPv6 && cktype != IPv4MappedIPv6 && cktype != IPv4CIDR && cktype != IPv6CIDR && cktype != IPv4MappedIPv6CIDR {
			return false, errors.New("Invalid Options")
		}
	}

	l := len(str)
	ret := getIPType(str, l)

	for _, ck := range cktypes {

		if ret != none && (ck == ret || ck == IPAny) {

			switch ret {
			case IPv4, IPv6, IPv4MappedIPv6:
				ip := net.ParseIP(str)

				if ip != nil {
					return true, nil
				}

			case IPv4CIDR, IPv6CIDR, IPv4MappedIPv6CIDR:
				_, _, err := net.ParseCIDR(str)
				if err == nil {
					return true, nil
				}
			}
		}
	}

	return false, errors.New("Invalid IPAddr")
}

// isCIDR is Validates whether the value IP Address with CIRD
func isCIDR(str string, l int) bool {

	if str[l-3] == '/' || str[l-2] == '/' {

		cidrBit := strings.Split(str, "/")
		if 2 == len(cidrBit) {
			bit, err := strconv.Atoi(cidrBit[1])
			//IPv4 : 0~32, IPv6 : 0 ~ 128
			if err == nil && bit >= 0 && bit <= 128 {
				return true
			}
		}
	}

	return false
}

// getIPType is Get a type of IP Address
func getIPType(str string, l int) int {

	if l < 3 { //least 3 chars (::F)
		return none
	}

	hasDot := strings.Index(str[2:], ".")
	hasColon := strings.Index(str[2:], ":")

	switch {
	case hasDot > -1 && hasColon == -1 && l >= 7 && l <= IPv4CIDR:
		if isCIDR(str, l) {
			return IPv4CIDR
		}
		return IPv4
	case hasDot == -1 && hasColon > -1 && l >= 6 && l <= IPv6CIDR:
		if isCIDR(str, l) {
			return IPv6CIDR
		}
		return IPv6

	case hasDot > -1 && hasColon > -1 && l >= 14 && l <= IPv4MappedIPv6:
		if isCIDR(str, l) {
			return IPv4MappedIPv6CIDR
		}
		return IPv4MappedIPv6
	}

	return none
}

const regexDenyFileNameCharList = `[\x00-\x1f|\x21-\x2c|\x3b-\x40|\x5b-\x5e|\x60|\x7b-\x7f]+`
const regexDenyFileName = `|\x2e\x2e\x2f+`

var checkAllowRelativePath = regexp.MustCompile(`(?m)(` + regexDenyFileNameCharList + `)`)
var checkDenyRelativePath = regexp.MustCompile(`(?m)(` + regexDenyFileNameCharList + regexDenyFileName + `)`)

// IsValidFilePath is Validates whether the value is a valid FilePath without relative path
func (s *StringValidator) IsValidFilePath(str string) bool {

	ret := checkDenyRelativePath.MatchString(str)
	return !ret
}

// IsValidFilePathWithRelativePath is Validates whether the value is a valid FilePath (allow with relative path)
func (s *StringValidator) IsValidFilePathWithRelativePath(str string) bool {

	ret := checkAllowRelativePath.MatchString(str)
	return !ret
}

// IsPureTextStrict is Validates whether the value is a pure text, Validation use native
func (s *StringValidator) IsPureTextStrict(str string) (bool, error) {

	l := len(str)

	for i := 0; i < l; i++ {

		c := str[i]

		// deny : control char (00-31 without 9(TAB) and Single 10(LF),13(CR)
		//if c >= 0 && c <= 31 && c != 9 && c != 10 && c != 13 { unsinged value is always >= 0
		if c <= 31 && c != 9 && c != 10 && c != 13 {
			return false, errors.New("Detect Control Character")
		}

		// deny : control char (DEL)
		if c == 127 {
			return false, errors.New("Detect Control Character (DEL)")
		}

		//deny : html tag (< ~ >)
		if c == 60 {

			ds := 0
			for n := i; n < l; n++ {

				// 60 (<) , 47(/) | 33(!) | 63(?)
				if str[n] == 60 && n+1 <= l && (str[n+1] == 47 || str[n+1] == 33 || str[n+1] == 63) {
					ds = 1
					n += 3 //jump to next char
				}

				// 62 (>)
				if (str[n] == 62 && n >= (i+2)) || (ds == 1 && str[n] == 62) {
					return false, errors.New("Detect Tag (<[!|?]~>)")
				}
			}
		}

		//deby : html encoded tag (&xxx;)
		if c == 38 && i+1 <= l && str[i+1] != 35 {

			max := i + 64
			if max > l {
				max = l
			}
			for n := i; n < max; n++ {
				if str[n] == 59 {
					return false, errors.New("Detect HTML Encoded Tag (&XXX;)")
				}
			}
		}
	}

	return true, nil
}

// Requires a string to match a given html tag elements regex pattern
// referrer : http://www.w3schools.com/Tags/
var elementPattern = regexp.MustCompile(`(?im)<(?P<tag>(/*\s*|\?*|\!*)(figcaption|expression|blockquote|plaintext|textarea|progress|optgroup|noscript|noframes|menuitem|frameset|fieldset|!DOCTYPE|datalist|colgroup|behavior|basefont|summary|section|isindex|details|caption|bgsound|article|address|acronym|strong|strike|source|select|script|output|option|object|legend|keygen|ilayer|iframe|header|footer|figure|dialog|center|canvas|button|applet|video|track|title|thead|tfoot|tbody|table|style|small|param|meter|layer|label|input|frame|embed|blink|audio|aside|alert|time|span|samp|ruby|meta|menu|mark|main|link|html|head|form|font|code|cite|body|base|area|abbr|xss|xml|wbr|var|svg|sup|sub|pre|nav|map|kbd|ins|img|div|dir|dfn|del|col|big|bdo|bdi|!--|ul|tt|tr|th|td|rt|rp|ol|li|hr|em|dt|dl|dd|br|u|s|q|p|i|b|a|(h[0-9]+)))([^><]*)([><]*)`)

// Requires a string to match a given urlencoded regex pattern
var urlencodedPattern = regexp.MustCompile(`(?im)(\%[0-9a-fA-F]{1,})`)

// Requires a string to match a given control characters regex pattern (ASCII : 00-08, 11, 12, 14, 15-31)
var controlcharPattern = regexp.MustCompile(`(?im)([\x00-\x08\x0B\x0C\x0E-\x1F\x7F]+)`)

// IsPureTextNormal is Validates whether the value is a pure text, Validation use Regular Expressions
func (s *StringValidator) IsPureTextNormal(str string) (bool, error) {

	decodedStr := html.UnescapeString(str)

	matchedUrlencoded := urlencodedPattern.MatchString(decodedStr)
	if matchedUrlencoded {
		tempBuf, err := url.QueryUnescape(decodedStr)
		if err == nil {
			decodedStr = tempBuf
		}
	}

	matchedElement := elementPattern.MatchString(decodedStr)
	if matchedElement {
		return false, errors.New("Detect HTML Element")
	}

	matchedCc := controlcharPattern.MatchString(decodedStr)
	if matchedCc {
		return false, errors.New("Detect Control Character")
	}

	return true, nil
}
