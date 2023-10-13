package main

import (
	"fmt"
	"net/http"
	"testing"
)

const url = "http://127.0.0.1:3000"
const folder = ".\\example"
var core = newCore(folder, 3000)

func TestSearchFolder(t *testing.T) {
    core.searchFolder(folder, "")
}

func TestStartServer(t *testing.T) {
    core.disableStartMessage = true
    core.logs = false
    core.startServer()
}

func TestNormalRoute(t *testing.T) {
    err := testRequest(url)
    if (err != nil) {
        t.Error(err.Error())
    }
}

func TestGetFile(t *testing.T) {
    const endpoint = url + "/page2/text.txt"
    err := testRequest(endpoint)
    if (err != nil) {
        t.Error(err.Error())
    }
}

func TestPrivateFile(t *testing.T) {
    const endpoint = url + "/secret.txt"
    err := testRequest(endpoint)
    if (err == nil) {
        t.Error("Endpoint " + endpoint + " returned wrong status code \nExpected: 404\n Got: 200")
    }
}

func TestDynamicRoute(t *testing.T) {
    const endpoint = url + "/user/example"
    err := testRequest(endpoint)
    if (err != nil) {
        t.Error(err.Error())
    }
}

func testRequest(endpoint string) error {
    req, err := http.Get(endpoint)
    if (err != nil) {
        return err
    }
    if (req.StatusCode != 200) {
        return fmt.Errorf("Endpoint `%v` returned wrong status code \nExpected: 200\n Got: %v", endpoint, req.StatusCode)
    }
    return nil
}
