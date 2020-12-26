package http_tags

import (
	"reflect"
	"net/http"
  "strconv"
  "bytes"
  "mime/multipart"
  "io/ioutil"
	"fmt"
)

var struct_tag = "http"

func SetStructTag(tag string) {
	struct_tag = tag
}

func GetStructTag() string {
	return struct_tag
}


/* never forget to use a pointer as prm*/
func FillInterfaceFromRequest(u interface{}, r *http.Request, ignore map[string]int) {

  s := reflect.ValueOf(u).Elem()
	typeV := s.Type()

	ignored_all := false
	ignored_cnt := int(0)
	if ignored_cnt == len(ignore) {
		ignored_all = true
	}
	for i := 0; i < typeV.NumField(); i++ {

		tag := typeV.Field(i).Tag.Get(struct_tag)
		if tag == "-" {
			continue
		}
		if !ignored_all {
			_, ignored := ignore[tag]
			if ignored {
				ignored_cnt += 1
				if ignored_cnt == len(ignore) {
					ignored_all = true
				}
				continue
			}
		}

    formvalue := r.PostFormValue(tag)
		if len(formvalue) == 0 {
			continue
		}

		field := s.Field(i)
    switch  (field.Type().Kind()) {
      case  reflect.String: {
        field.SetString(formvalue)
      }
      case  reflect.Bool: {
        v, err := strconv.ParseBool(formvalue)
        if err == nil {
          field.SetBool(v)
        } else {
					fmt.Println(err.Error())
				}
      }
      case  reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: {
        v, err:= strconv.ParseInt(formvalue, 10, 64)
        if err == nil {
          field.SetInt(v)
        } else {
					fmt.Println(err.Error())
				}
      }
      case  reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: {
        v, err:= strconv.ParseUint(formvalue, 10, 64)
        if err == nil {
          field.SetUint(v)
        } else {
					fmt.Println(err.Error())
				}
      }
      case  reflect.Float32, reflect.Float64: {
        v, err:= strconv.ParseFloat(formvalue, 64)
        if err == nil {
          field.SetFloat(v)
        } else {
					fmt.Println(err.Error())
				}
      }
      default: continue
    }
  }

}

/* never forget to use a pointer as prm*/
func PutFieldsToRequest(u interface{}, r *http.Request) {
  buf := new(bytes.Buffer)

  writer := multipart.NewWriter(buf)

  s := reflect.ValueOf(u).Elem()
	typeV := s.Type()

	for i := 0; i < typeV.NumField(); i++ {

		field := s.Field(i)
		tag := typeV.Field(i).Tag.Get(struct_tag)
		if tag == "-" {
			continue
		}

    switch  (field.Type().Kind()) {
      case  reflect.String  : {
        writer.WriteField(tag, field.String())
      }
      case  reflect.Bool    : {
        writer.WriteField(tag, strconv.FormatBool(field.Bool()))
      }
      case  reflect.Int,  reflect.Int8,  reflect.Int16,  reflect.Int32,  reflect.Int64: {
          writer.WriteField(tag, strconv.FormatInt(field.Int(), 10))
      }
      case  reflect.Uint8,  reflect.Uint16,  reflect.Uint32,  reflect.Uint64 : {
        writer.WriteField(tag, strconv.FormatUint(field.Uint(), 10))
      }
      case  reflect.Float32,  reflect.Float64 : {
        writer.WriteField(tag, strconv.FormatFloat(field.Float(), 'f', 3, 64))
      }
      default: continue
    }
  }

  writer.Close()

  r.Header.Add("Content-Type", writer.FormDataContentType())
  r.Body = ioutil.NopCloser(buf)
}
