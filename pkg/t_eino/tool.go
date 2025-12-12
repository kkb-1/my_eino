package t_eino

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/mark3labs/mcp-go/client"
	mcp2 "github.com/mark3labs/mcp-go/mcp"
)

const (
	todoToolStateKey = "tool:todo"
	getToolStateKey  = "tool:get_tool"
)

type HideTool struct {
	hide  map[string]*schema.ToolInfo
	alive []*schema.ToolInfo
}

type getToolArguments struct {
	Name string `json:"name"`
}

type TODO struct {
	Content string `json:"content"`
	Status  string `json:"status" jsonschema:"enum=pending,enum=in_progress,enum=completed"`
}

type writeTodosArguments struct {
	Todos []TODO `json:"todos"`
}

// 按需显露 tool
func ToolGetGetTool(lm *LearnModel, hide, alive []*schema.ToolInfo) (tool.BaseTool, error) {
	var hideMap map[string]*schema.ToolInfo
	for _, info := range hide {
		hideMap[info.Name] = info
	}
	t, err := utils.InferTool("get_tool", WriteTodosToolDescription, func(ctx context.Context, input getToolArguments) (output string, err error) {
		if err := compose.ProcessState(ctx, func(ctx context.Context, s *state) error {
			s.lock.Lock()
			defer s.lock.Unlock()
			ht, ok := s.toolState[getToolStateKey].(*HideTool)
			if !ok {
				ht = &HideTool{
					hide:  hideMap,
					alive: alive,
				}
				s.toolState[getToolStateKey] = ht
			}

			info, ok := ht.hide[input.Name]
			if !ok {
				return errors.New("this tool is not exist")
			}
			ht.alive = append(ht.alive, info)
			_, err = lm.WithTools(ht.alive)
			return err
		}); err != nil {
			return "", err
		}

		return fmt.Sprintf("get tool %s success", input.Name), nil
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}

func ToolGetWriteTodos() (tool.BaseTool, error) {
	t, err := utils.InferTool("write_todos", WriteTodosToolDescription, func(ctx context.Context, input writeTodosArguments) (output string, err error) {
		if err := compose.ProcessState(ctx, func(ctx context.Context, s *state) error {
			s.lock.Lock()
			defer s.lock.Unlock()
			s.toolState[todoToolStateKey] = input.Todos
			return nil
		}); err != nil {
			return "", err
		}

		todos, err := sonic.MarshalString(input.Todos)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Updated todo list to %s", todos), nil
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}

func ToolsGetSandbox(ctx context.Context) ([]tool.BaseTool, error) {
	cli, err := client.NewStreamableHttpClient("http://localhost:8080/mcp")
	if err != nil {
		return nil, err
	}
	// sse client  needs to manually start asynchronous communication
	// while stdio does not require it.
	err = cli.Start(ctx)
	if err != nil {
		return nil, err
	}
	// Initialize
	initRequest := mcp2.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp2.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp2.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}

	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		return nil, err
	}
	return mcp.GetTools(ctx, &mcp.Config{Cli: cli})
}
