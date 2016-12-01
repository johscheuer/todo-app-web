// +build integration

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/uber/tchannel-go/testutils/goroutines"
)

var todoAppServer = "http://localhost:3000"

func getHTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Duration(1 * time.Second),
	}
}

//TestHealthEndpoint checks if the Health endpoint has the right format
func TestHealthEndpoint(t *testing.T) {
	defer validateGoRoutines()

	expectedResponse := map[string]string{
		"redis-master-0": "ok",
		"redis-slave-0":  "ok",
		"self":           "ok",
	}

	resp, err := getHTTPClient().Get(fmt.Sprintf("%s/health", todoAppServer))
	if err != nil || resp.StatusCode != 200 {
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
	defer validateGoRoutines()
	insertItem := "TestCase"

	// Insert Item
	if resp, err := getHTTPClient().Get(fmt.Sprintf("%s/insert/todo/%s", todoAppServer, insertItem)); err != nil || resp.StatusCode != 200 {
		t.Log(err)
		t.FailNow()
	}

	// Read Item
	readResp, err := getHTTPClient().Get(fmt.Sprintf("%s/read/todo", todoAppServer))
	if err != nil || readResp.StatusCode != 200 {
		t.Log(err)
		t.FailNow()
	}

	defer readResp.Body.Close()
	var readResponse []string
	if err = json.NewDecoder(readResp.Body).Decode(&readResponse); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !reflect.DeepEqual([]string{insertItem}, readResponse) {
		t.FailNow()
	}

	// Delete Item
	deleteResp, err := getHTTPClient().Get(fmt.Sprintf("%s/delete/todo/%s", todoAppServer, insertItem))
	if err != nil || deleteResp.StatusCode != 200 {
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

func TestWhoAmI(t *testing.T) {
	defer validateGoRoutines()
	readResp, err := getHTTPClient().Get(fmt.Sprintf("%s/whoami", todoAppServer))
	if err != nil || readResp.StatusCode != 200 {
		t.Log(err)
		t.FailNow()
	}

	//TODO we would ned to set a fix IP address to this container
	defer readResp.Body.Close()
}

func validateGoRoutines() {
	if err := goroutines.IdentifyLeaks(&goroutines.VerifyOpts{
		Excludes: []string{
			"net/http",
		},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Found goroutine leaks on successful test run: %v", err)
	}
}
