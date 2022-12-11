package adepttech

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"golang.org/x/oauth2"
)

func EncodeStructAsParams(obj interface{}) url.Values {
	values := make(url.Values)
	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

	for i := 0; i < objValue.NumField(); i++ {
		k := objType.Field(i).Name
		jsonTag := objValue.Type().Field(i).Tag.Get("json")
		if len(jsonTag) > 0 {
			switch jsonTag {
			case "-":
			default:
				parts := strings.Split(jsonTag, ",")
				k = parts[0]
			}
		}
		raw := objValue.Field(i).Interface()
		switch v := raw.(type) {
		case string:
			values.Add(k, v)
		case int, int64, int32, uint, float32, float64:
			values.Add(k, fmt.Sprintf("%v", v))
		default:
			b, err := json.Marshal(v)
			if err != nil {
				panic(err) // possible?
			}
			values.Add(k, string(b))

		}
	}
	return values
}

func MarshalToken(token *oauth2.Token) *MarshalledToken {
	return &MarshalledToken{token: token}
}

func UnmarshalToken(b []byte) (*MarshalledToken, error) {
	token := new(oauth2.Token)
	err := json.Unmarshal([]byte(b), token)
	if err != nil {
		return nil, err
	}
	return &MarshalledToken{token: token}, nil
}

type MarshalledToken struct {
	token *oauth2.Token
}

func (mt *MarshalledToken) Token() (*oauth2.Token, error) {
	return mt.token, nil
}

func (mt *MarshalledToken) Bytes() ([]byte, error) {
	return json.Marshal(mt.token)
}
