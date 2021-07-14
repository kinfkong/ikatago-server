package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	myerrors "github.com/kinfkong/ikatago-server/errors"

	"moul.io/http2curl"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// DoHTTPRequest Sends generic http request
func DoHTTPRequest(method string, url string, headers map[string]string, body []byte) (responseBody string, err error) {
	httpClient := &http.Client{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Close = true
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	command, _ := http2curl.GetCurlCommand(req)
	// log.Printf("DEBUG sending with http: %s\n", command)
	response, err := httpClient.Do(req)
	if err != nil {
		log.Printf("ERROR error requesting with http: %s, error: %v\n", command, err)
		err = errors.New("failed_do_request")
		return
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		log.Printf("ERROR error requesting with http: %s, error: %v\n", command, err)
		err = errors.New("failed_read_body")
		return
	}

	responseBody = string(bodyBytes)

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		log.Printf("ERROR error requesting with http: %s, status code: %v, response: %s\n", command, response.StatusCode, responseBody)
		err = errors.New("invalid_status")
		return
	}

	return
}

// FileExists checks if a file exists and is not a directory before we
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirectoryExists checks if a dir exists and is not a directory before we
func DirectoryExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// GetJSONNumber gets the json number
func GetJSONNumber(val interface{}) (float64, error) {
	switch v := val.(type) {
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	case json.Number:
		return v.Float64()
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case string:
		result, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, myerrors.CreateError(400, "invalid_val_type")
		}
		return result, nil
	case *json.Number:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)
	case *bool:
		if v != nil && *v {
			return 1, nil
		} else {
			return 0, nil
		}
	case *float64:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)
	case *float32:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)

	case *int:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)
	case *int64:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)
	case *int8:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)
	case *int16:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)
	case *int32:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)

	case *uint:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)
	case *uint64:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)
	case *uint8:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)
	case *uint16:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)
	case *uint32:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)
	case *string:
		if v == nil {
			return 0, nil
		}
		return GetJSONNumber(*v)
	default:
		return 0, myerrors.CreateError(400, "invalid_val_type")
	}
}

// GetJSONIntNumber gets the json number
func GetJSONIntNumber(val interface{}) (int, error) {
	f, err := GetJSONNumber(val)
	if err != nil {
		return 0, err
	}
	return (int)(f + 0.5), nil
}

// GetJSONInt64Number gets the json number
func GetJSONInt64Number(val interface{}) (int64, error) {
	f, err := GetJSONNumber(val)
	if err != nil {
		return 0, err
	}
	return (int64)(f + 0.5), nil
}
