package main

import (
	"bytes"
	"context"
	"os/exec"
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

var oldExecCommandContext = execCommandContext

func makeWeatherRequest(args map[string]interface{}) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	params := reflect.ValueOf(&req).Elem().FieldByName("Params")
	arguments := params.FieldByName("Arguments")
	if arguments.CanSet() {
		arguments.Set(reflect.ValueOf(args))
	} else {
		panic("Cannot set Arguments field via reflection for weather")
	}
	return req
}

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

func TestWeatherHandler_Success(t *testing.T) {
	ctx := context.Background()
	defer func() { execCommandContext = oldExecCommandContext }()
	execCommandContext = func(ctx context.Context, name string, arg ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "echo", "Berlin: +20 C, Sunny")
	}
	req := makeWeatherRequest(map[string]interface{}{"location": "Berlin"})
	res, err := weatherHandler(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Content) == 0 {
		t.Fatalf("expected non-empty content")
	}
	val := getContentText(t, res)
	if !bytes.Contains([]byte(val), []byte("Berlin")) {
		t.Errorf("wanted weather for Berlin, got: %s", val)
	}
}

func TestWeatherHandler_ParamMissingOrNonString(t *testing.T) {
	ctx := context.Background()
	badInputs := []interface{}{nil, 123, 1.2, []string{"Berlin"}, map[string]int{"foo": 1}}
	for _, val := range badInputs {
		req := makeWeatherRequest(map[string]interface{}{"location": val})
		res, err := weatherHandler(ctx, req)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if len(res.Content) == 0 {
			t.Errorf("expected error content, got none for input %v", val)
			continue
		}
		msg := getContentText(t, res)
		if !bytes.Contains([]byte(msg), []byte("string")) {
			t.Errorf("expected error about string type, got: %s", msg)
		}
	}
}

func TestWeatherHandler_ExecFails(t *testing.T) {
	ctx := context.Background()
	defer func() { execCommandContext = oldExecCommandContext }()
	execCommandContext = func(ctx context.Context, name string, arg ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "false") // always fails
	}
	req := makeWeatherRequest(map[string]interface{}{"location": "nowhere"})
	res, err := weatherHandler(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Content) == 0 {
		t.Error("expected error content for exec fail")
		return
	}
	msg := getContentText(t, res)
	if !bytes.Contains([]byte(msg), []byte("exit status")) {
		t.Errorf("expected exec error, got: %s", msg)
	}
}
