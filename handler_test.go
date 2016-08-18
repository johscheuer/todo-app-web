// +build integration

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var todoAppServer = "http://localhost:3000"

func getHTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
}

//TestHealthEndpoint checks if the Health endpoint has the right format
func TestHealthEndpoint(t *testing.T) {
	expectedResponse := map[string]string{
		"redis-master": "ok",
		"redis-slave":  "ok",
		"self":         "ok",
	}

	resp, err := getHTTPClient().Get(fmt.Sprintf("%s/health", todoAppServer))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var Response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(expectedResponse, Response) {
		t.Fail()
	}
}

func TestInsertReadAndDeleteItem(t *testing.T) {
	insertItem := "TestCase"

	// Insert Item
	if _, err := getHTTPClient().Get(fmt.Sprintf("%s/insert/todo/%s", todoAppServer, insertItem)); err != nil {
		t.Log(err)
		t.FailNow()
	}

	// Read Item
	readResp, err := getHTTPClient().Get(fmt.Sprintf("%s/read/todo", todoAppServer))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer readResp.Body.Close()
	var readResponse []string
	if err := json.NewDecoder(readResp.Body).Decode(&readResponse); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !reflect.DeepEqual([]string{insertItem}, readResponse) {
		t.FailNow()
	}

	// Delete Item
	deleteResp, err := getHTTPClient().Get(fmt.Sprintf("%s/delete/todo/%s", todoAppServer, insertItem))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer deleteResp.Body.Close()
	var deleteResponse []string
	if err := json.NewDecoder(deleteResp.Body).Decode(&deleteResponse); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !reflect.DeepEqual([]string{}, deleteResponse) {
		t.FailNow()
	}
}

//TODO func checkResponse

/* TODO Tests for:
- whoAmIHandler
*/
