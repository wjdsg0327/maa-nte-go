package custom

import (
	"encoding/json"

	maa "github.com/MaaXYZ/maa-framework-go/v4"
)

// ExampleRecognition 示例自定义识别器
// 实现 CustomRecognitionRunner 接口来添加自定义识别逻辑
type ExampleRecognition struct{}

type ExampleRecognitionParam struct {
	Label string `json:"label"`
}

// Run 执行自定义识别
// ctx: 上下文对象，可用于运行子任务、访问 Tasker 等
// arg: 识别参数，包含当前截图、ROI 区域等信息
// 返回值：识别结果（包含目标矩形框）和是否识别成功
func (r *ExampleRecognition) Run(ctx *maa.Context, arg *maa.CustomRecognitionArg) (*maa.CustomRecognitionResult, bool) {
	var param ExampleRecognitionParam
	_ = json.Unmarshal([]byte(arg.CustomRecognitionParam), &param)

	box := arg.Roi
	if box.Width() <= 0 || box.Height() <= 0 {
		box = maa.Rect{100, 100, 200, 80}
	}

	detail := map[string]interface{}{
		"label": param.Label,
		"task":  arg.CurrentTaskName,
	}
	detailBytes, _ := json.Marshal(detail)

	return &maa.CustomRecognitionResult{
		Box:    box,
		Detail: string(detailBytes),
	}, true
}
