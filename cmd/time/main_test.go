package main

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func makeTimeRequest() mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	// set empty Arguments for completeness
	params := reflect.ValueOf(&req).Elem().FieldByName("Params")
	arguments := params.FieldByName("Arguments")
	if arguments.CanSet() {
		arguments.Set(reflect.ValueOf(map[string]interface{}{}))
	} else {
		panic("Cannot set Arguments in Params via reflection for time")
	}
	return req
}

func TestTimeHandler_Basic(t *testing.T) {
	req := makeTimeRequest()
	res, err := timeHandler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil response")
	}
	if len(res.Content) == 0 {
		t.Error("expected content")
	}
	first, ok := res.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Content[0])
	}
	if !strings.Contains(first.Text, "It's now ") {
		t.Errorf("unexpected time response: %s", first.Text)
	}
}
