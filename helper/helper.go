package helper

import (
	"fmt"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"io/ioutil"
	"unicode"
)

func ReadAdminUI() string {
	data, err := ioutil.ReadFile("./html/admin.html")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

func TokenPage() string {
	data, err := ioutil.ReadFile("./html/get_token.html")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

func ShowBanner() {
	data, err := ioutil.ReadFile("./log/banner.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}

func ReadLogFile(_type string, file string) string {
	if _type == "info" || _type == "error" {
		data, err := ioutil.ReadFile(fmt.Sprintf("./log_files/%s/zendesk_integration_%s_%s.log", _type, file, _type))
		if err != nil {
			return err.Error()
		}
		fmt.Println(string(data))
		return string(data)

	}
	return ""

}


func RemoveUnicodeChar(s string) (string,error) {
	b := make([]byte, len(s))

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	_, _, e := t.Transform(b, []byte(s), true)
	if e != nil {
		return "", e
	}
	return string(b), nil
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}


// func EmptyIfNil(v interface{}) interface{} {
// 	switch reflect.TypeOf(v).Kind() {
// 	case reflect.Array:
// 		s := reflect.ValueOf(v)
// 		if s.Len() > 0 {
// 			return
// 		} else {
// 			return ""
// 		}
// 	}
// }
