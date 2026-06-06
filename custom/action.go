package custom

import (
	"encoding/json"

	maa "github.com/MaaXYZ/maa-framework-go/v4"
)

// ExampleAction 示例自定义动作
// 实现 CustomActionRunner 接口来添加自定义动作逻辑
type ExampleAction struct{}

type ExampleActionParam struct {
	Offset []int `json:"offset"`
}

// Run 执行自定义动作
// ctx: 上下文对象，可用于运行子任务、执行点击/滑动等操作
// arg: 动作参数，包含识别到的目标区域、识别详情等信息
// 返回值：true 表示动作成功，false 表示失败
func (a *ExampleAction) Run(ctx *maa.Context, arg *maa.CustomActionArg) bool {
	box := arg.Box
	if box.Width() <= 0 || box.Height() <= 0 {
		box = maa.Rect{100, 100, 1, 1}
	}

	x := int32(box.X() + box.Width()/2)
	y := int32(box.Y() + box.Height()/2)

	var param ExampleActionParam
	_ = json.Unmarshal([]byte(arg.CustomActionParam), &param)
	if len(param.Offset) == 2 {
		x += int32(param.Offset[0])
		y += int32(param.Offset[1])
	}

	job := ctx.GetTasker().GetController().PostClick(x, y)
	if job == nil {
		return false
	}
	return job.Wait().Status().Success()
}
