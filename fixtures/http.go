package fixtures

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rabi-food-core/config"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

var url = fmt.Sprintf("http://localhost:%s", config.TestPort)

type PostInput struct {
	Body     any
	URI      string
	Response any
	Token    string
}

func Post(t *testing.T, input PostInput) (statusCode int) {
	b, err := json.Marshal(input.Body)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, url+input.URI, bytes.NewBuffer(b))
	require.Nil(t, err)

	req.Header.Set("Content-Type", "application/json")
	if input.Token != "" {
		req.Header.Set("Authorization", "Bearer "+input.Token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	statusCode = resp.StatusCode

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	if resp.Header.Get("Content-Type") != "application/json" {
		if v, ok := input.Response.(*string); ok {
			*v = string(responseBody)
		}

		return
	}

	err = json.Unmarshal(responseBody, input.Response)
	require.Nil(t, err)

	return
}

type GetInput struct {
	URI      string
	Response any
	Token    string
	Query    any
}

func Get(t *testing.T, input GetInput) (statusCode int) {
	req, err := http.NewRequest(http.MethodGet, url+input.URI, nil)
	require.Nil(t, err)

	if input.Token != "" {
		req.Header.Set("Authorization", "Bearer "+input.Token)
	}

	if input.Query != nil {
		q := req.URL.Query()
		v := reflect.ValueOf(input.Query)
		typeOfS := v.Type()

		for i := 0; i < v.NumField(); i++ {
			q.Add(typeOfS.Field(i).Name, fmt.Sprintf("%v", v.Field(i).Interface()))
		}

		req.URL.RawQuery = q.Encode()
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	statusCode = resp.StatusCode

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	if v, ok := input.Response.(*string); ok {
		*v = string(responseBody)
		return
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		return
	}

	err = json.Unmarshal(responseBody, input.Response)
	require.Nil(t, err)

	return
}

type PatchInput struct {
	URI      string
	Body     any
	Response any
	Token    string
	Query    any
}

func Patch(t *testing.T, input PatchInput) (statusCode int) {
	b, err := json.Marshal(input.Body)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPatch, url+input.URI, bytes.NewBuffer(b))
	require.Nil(t, err)

	req.Header.Set("Content-Type", "application/json")
	if input.Token != "" {
		req.Header.Set("Authorization", "Bearer "+input.Token)
	}

	if input.Query != nil {
		q := req.URL.Query()
		v := reflect.ValueOf(input.Query)
		typeOfS := v.Type()

		for i := 0; i < v.NumField(); i++ {
			q.Add(typeOfS.Field(i).Name, fmt.Sprintf("%v", v.Field(i).Interface()))
		}

		req.URL.RawQuery = q.Encode()
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	statusCode = resp.StatusCode

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	if v, ok := input.Response.(*string); ok {
		*v = string(responseBody)
		return
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		return
	}

	err = json.Unmarshal(responseBody, input.Response)
	require.Nil(t, err)

	return
}

type DeleteInput struct {
	URI   string
	Body  any
	Token string
}

func Delete(t *testing.T, input DeleteInput) (string, int) {
	b, err := json.Marshal(input.Body)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodDelete, url+input.URI, bytes.NewBuffer(b))
	require.Nil(t, err)

	req.Header.Set("Content-Type", "application/json")
	if input.Token != "" {
		req.Header.Set("Authorization", "Bearer "+input.Token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	return string(responseBody), resp.StatusCode
}
