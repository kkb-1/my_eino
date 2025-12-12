# my_eino

## 项目介绍

该项目基于字节开源的 eino 框架，魔改更为适用日常开发的工具。

目前已实现：

* react agent :根据模型实际需要，通过特殊 tool（get_tool）暴露隐藏起来的工具，让模型按需获取工具。
* fix eino 原生 react agent 工具调用只检查第一个 message 是否包含 tool_call 短板

todo:

* skill：让 react agent 知道接下来需要获取什么工具。
* 兼容 eino adk.Agent 接口的我的 super agent（继承 react agent 的优点）
* agui 组件搭建：基于 hertz 的 agui 组件
* 基于 agui 与 mcp 的端 tool