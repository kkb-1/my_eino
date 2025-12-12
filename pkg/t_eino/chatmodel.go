package t_eino

import (
	"context"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

var _ model.ToolCallingChatModel = &LearnModel{}

type LearnToolFunc func([]*schema.ToolInfo)

type LearnModel struct {
	tChatModel        model.ToolCallingChatModel
	originalChatModel model.ToolCallingChatModel
}

func NewLearnModel(ctx context.Context, llm model.ToolCallingChatModel, hideToolsCfg, initToolsCfg compose.ToolsNodeConfig) (*LearnModel, compose.ToolsNodeConfig, error) {
	lm := &LearnModel{
		tChatModel:        llm,
		originalChatModel: llm,
	}
	hideInfos, err := genToolInfos(ctx, hideToolsCfg)
	if err != nil {
		return nil, compose.ToolsNodeConfig{}, err
	}
	initInfos, err := genToolInfos(ctx, initToolsCfg)
	if err != nil {
		return nil, compose.ToolsNodeConfig{}, err
	}
	getToolTool, err := ToolGetGetTool(lm, hideInfos, initInfos)
	if err != nil {
		return nil, compose.ToolsNodeConfig{}, err
	}
	tools := make([]tool.BaseTool, len(hideInfos)+len(initInfos)+1, 0)
	tools = append(tools, hideToolsCfg.Tools...)
	tools = append(tools, initToolsCfg.Tools...)
	tools = append(tools, getToolTool)
	totalToolsCfg := compose.ToolsNodeConfig{
		Tools:                tools,
		UnknownToolsHandler:  initToolsCfg.UnknownToolsHandler,
		ExecuteSequentially:  initToolsCfg.ExecuteSequentially,
		ToolArgumentsHandler: initToolsCfg.ToolArgumentsHandler,
		ToolCallMiddlewares:  initToolsCfg.ToolCallMiddlewares,
	}

	return lm, totalToolsCfg, nil
}

func (l *LearnModel) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	return l.tChatModel.Generate(ctx, input, opts...)
}

func (l *LearnModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	return l.tChatModel.Stream(ctx, input, opts...)
}

func (l *LearnModel) WithTools(tools []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	newModel, err := l.originalChatModel.WithTools(tools)
	if err != nil {
		return nil, err
	}
	l.tChatModel = newModel
	return l, err
}
