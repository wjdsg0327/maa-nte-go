package custom

import (
	maa "github.com/MaaXYZ/maa-framework-go/v4"
)

// ExampleAction 示例自定义动作
// 实现 CustomActionRunner 接口来添加自定义动作逻辑
type ExampleAction struct{}

// Run 执行自定义动作
// ctx: 上下文对象，可用于运行子任务、执行点击/滑动等操作
// arg: 动作参数，包含识别到的目标区域、识别详情等信息
// 返回值：true 表示动作成功，false 表示失败
func (a *ExampleAction) Run(ctx *maa.Context, arg *maa.CustomActionArg) bool {
	// TODO: 在此实现你的动作逻辑
	// 可通过 ctx.RunTask() 运行子任务
	// 可通过 ctx.RunAction() 执行内置动作（点击、滑动等）
	//
	// 动作成功返回 true，失败返回 false
	return false
}
