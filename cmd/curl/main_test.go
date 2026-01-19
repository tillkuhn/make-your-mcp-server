package main

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"reflect"
	"testing"
)

func makeCurlRequest(args map[string]interface{}) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	params := reflect.ValueOf(&req).Elem().FieldByName("Params")
	arguments := params.FieldByName("Arguments")
	if arguments.CanSet() {
		arguments.Set(reflect.ValueOf(args))
	} else {
		panic("Cannot set Arguments field in Params via reflection for curl")
	}
	return req
}

func TestCurlHandler_Http(t *testing.T) {
	if _, err := exec.LookPath("curl"); err != nil {
		t.Skip("curl binary not available")
	}
	// Start a test http server for reliability
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer ts.Close()
	req := makeCurlRequest(map[string]interface{}{"url": ts.URL})
	resp, err := curlHandler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if len(resp.Content) == 0 {
		t.Error("expected content")
	}
	first, ok := resp.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", resp.Content[0])
	}
	if first.Text != "ok" {
		t.Errorf("expected http reply 'ok', got '%s'", first.Text)
	}
}
