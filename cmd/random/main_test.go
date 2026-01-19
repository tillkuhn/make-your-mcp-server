package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func getContentText(t *testing.T, res *mcp.CallToolResult) string {
	if len(res.Content) == 0 {
		t.Fatalf("Content must not be empty")
	}
	first := res.Content[0]
	textContent, ok := first.(mcp.TextContent)
	if !ok {
		t.Fatalf("Content is not TextContent, got %T", first)
	}
	return textContent.Text
}

func makeRequest(args map[string]interface{}) mcp.CallToolRequest {
	// Dirty hack: create an empty request, then set Arguments field in Params by reflection
	req := mcp.CallToolRequest{}
	params := reflect.ValueOf(&req).Elem().FieldByName("Params")
	arguments := params.FieldByName("Arguments")
	if arguments.CanSet() {
		arguments.Set(reflect.ValueOf(args))
	} else {
		// Fallback: panic and inform developer
		panic("Cannot set Arguments field in Params via reflection")
	}
	return req
}

func TestRandomHandler_Beer(t *testing.T) {
	req := makeRequest(map[string]interface{}{"thing": "beer"})
	resp, err := randomHandler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	value := getContentText(t, resp)
	if value == "" {
		t.Error("expected non-empty beer content")
	}
}

func TestRandomHandler_Job(t *testing.T) {
	req := makeRequest(map[string]interface{}{"thing": "job"})
	resp, err := randomHandler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	value := getContentText(t, resp)
	if value == "" {
		t.Error("expected non-empty job content")
	}
}

func TestRandomHandler_Food(t *testing.T) {
	req := makeRequest(map[string]interface{}{"thing": "food"})
	resp, err := randomHandler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	value := getContentText(t, resp)
	if value == "" {
		t.Error("expected non-empty food content")
	}
}

func TestRandomHandler_Hobby(t *testing.T) {
	req := makeRequest(map[string]interface{}{"thing": "hobby"})
	resp, err := randomHandler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	value := getContentText(t, resp)
	if value == "" {
		t.Error("expected non-empty hobby content")
	}
}
