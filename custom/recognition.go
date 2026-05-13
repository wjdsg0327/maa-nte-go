package custom

import (
	maa "github.com/MaaXYZ/maa-framework-go/v4"
)

// ExampleRecognition 示例自定义识别器
// 实现 CustomRecognitionRunner 接口来添加自定义识别逻辑
type ExampleRecognition struct{}

// Run 执行自定义识别
// ctx: 上下文对象，可用于运行子任务、访问 Tasker 等
// arg: 识别参数，包含当前截图、ROI 区域等信息
// 返回值：识别结果（包含目标矩形框）和是否识别成功
func (r *ExampleRecognition) Run(ctx *maa.Context, arg *maa.CustomRecognitionArg) (*maa.CustomRecognitionResult, bool) {
	// TODO: 在此实现你的识别逻辑
	// 可通过 ctx.RunTask() 运行子任务
	// 可通过 ctx.GetTasker() 获取 Tasker 实例
	//
	// 识别成功时返回目标区域矩形框和 true
	// 未识别到目标时返回 nil, false
	return nil, false
}
