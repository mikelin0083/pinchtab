package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
)

func handleRecordStart(c *Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		file, err := r.RequireString("file")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		ext := filepath.Ext(file)
		var format string
		switch ext {
		case ".gif":
			format = "gif"
		case ".webm":
			format = "webm"
		case ".mp4":
			format = "mp4"
		default:
			return mcp.NewToolResultError(fmt.Sprintf("unsupported format %q — use .gif, .webm, or .mp4", ext)), nil
		}

		payload := map[string]any{"format": format}
		if fps, ok := optInt(r, "fps"); ok {
			payload["fps"] = fps
		}
		if quality, ok := optInt(r, "quality"); ok {
			payload["quality"] = quality
		}
		if scale, ok := optFloat(r, "scale"); ok {
			payload["scale"] = scale
		}
		if tabID := optString(r, "tabId"); tabID != "" {
			payload["tabId"] = tabID
		}

		body, code, err := c.Post(ctx, "/record/start", payload)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return resultFromBytes(body, code)
	}
}

func handleRecordStop(c *Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		file := optString(r, "file")
		if file == "" {
			return mcp.NewToolResultError("file parameter is required"), nil
		}

		body, code, err := c.Post(ctx, "/record/stop", map[string]any{})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if code >= 400 {
			return resultFromBytes(body, code)
		}

		if err := os.WriteFile(file, body, 0600); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("write file: %v", err)), nil
		}

		return jsonResult(map[string]any{
			"status": "saved",
			"file":   file,
			"bytes":  len(body),
		})
	}
}

func handleRecordStatus(c *Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		body, code, err := c.Get(ctx, "/record/status", nil)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return resultFromBytes(body, code)
	}
}
