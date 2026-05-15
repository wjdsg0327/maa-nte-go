# Custom 使用教程（结合本项目）

这份文档讲的是 MaaFramework 的 `Custom` 识别和 `Custom` 动作在本项目里的用法。本项目是 Go 后端 + Vue 前端，Maa 绑定使用 `maa-framework-go/v4`，现有自定义代码入口在 `custom/recognition.go` 和 `custom/action.go`。

## 1. Custom 是什么

Pipeline 里大部分节点可以用 `TemplateMatch`、`OCR`、`Click` 这类内置能力解决。遇到下面情况时，就适合用 `Custom`：

- 识别逻辑需要自己写代码，例如多步判断、特殊图像处理、读取外部状态。
- 动作逻辑不是简单点击/滑动，例如根据识别结果决定点哪里、动态改变后续流程。
- 需要在一个节点里调用子识别、子任务、子动作。

MaaFramework 里有两种 Custom：

- `recognition: "Custom"`：自定义识别器，返回一个命中的矩形框。
- `action: "Custom"`：自定义动作，返回成功或失败。

## 2. 本项目当前位置

已有文件：

```text
custom/
  recognition.go  # 自定义识别器示例
  action.go       # 自定义动作示例
service/
  maa.go          # Maa 初始化、资源加载、任务执行
resource/
  pipeline/       # Pipeline JSON
```

当前 `custom` 目录里的 `ExampleRecognition` 和 `ExampleAction` 已经实现了接口，但还没有注册到 Maa Resource。也就是说：Pipeline 里写了 `custom_recognition: "ExampleRecognition"` 或 `custom_action: "ExampleAction"` 后，必须先在后端初始化资源时注册它们，否则执行会失败。

## 3. 必须先注册 Custom

在 `service/maa.go` 中，创建并加载 `s.resource` 后，需要调用：

```go
if err := s.resource.RegisterCustomRecognition("ExampleRecognition", &custom.ExampleRecognition{}); err != nil {
    return fmt.Errorf("注册自定义识别失败: %v", err)
}

if err := s.resource.RegisterCustomAction("ExampleAction", &custom.ExampleAction{}); err != nil {
    return fmt.Errorf("注册自定义动作失败: %v", err)
}
```

同时 `service/maa.go` 需要 import 本项目的 `custom` 包：

```go
import (
    // ...
    "study/one/custom"
)
```

建议封装成一个函数，避免以后 Custom 变多时散在初始化逻辑里：

```go
func (s *MaaService) registerCustomHandlersLocked() error {
    if err := s.resource.RegisterCustomRecognition("ExampleRecognition", &custom.ExampleRecognition{}); err != nil {
        return fmt.Errorf("注册 ExampleRecognition 失败: %v", err)
    }
    if err := s.resource.RegisterCustomAction("ExampleAction", &custom.ExampleAction{}); err != nil {
        return fmt.Errorf("注册 ExampleAction 失败: %v", err)
    }
    return nil
}
```

调用位置建议放在 `loadResourcesLocked()` 成功后、`BindResource` 前。因为 `ReloadResources()` 会调用 `resource.Clear()`，清理资源后最好重新注册一次，保证热重载后 Custom 仍可用。

## 4. Pipeline 里怎么填

### 自定义识别

```jsonc
{
  "自定义识别节点": {
    "recognition": "Custom",
    "custom_recognition": "ExampleRecognition",
    "custom_recognition_param": {
      "mode": "check_balance",
      "threshold": 0.8
    },
    "roi": [983, 260, 169, 388],
    "action": "Click"
  }
}
```

字段含义：

- `custom_recognition` 必须和 Go 注册名完全一致。
- `custom_recognition_param` 可以是任意 JSON，会以字符串形式传到 Go。
- `roi` 会传到 `arg.Roi`，格式必须是 `[x,y,w,h]` 四个整数。

### 自定义动作

```jsonc
{
  "自定义动作节点": {
    "recognition": "DirectHit",
    "action": "Custom",
    "custom_action": "ExampleAction",
    "custom_action_param": {
      "clickOffset": [10, 20],
      "nextWhenDone": "领取_补货"
    },
    "target": [100, 200, 80, 40]
  }
}
```

字段含义：

- `custom_action` 必须和 Go 注册名完全一致。
- `custom_action_param` 会以字符串形式传到 Go。
- `target` 会转换成 `arg.Box`，可以填 `true`、节点名、`[x,y]`、`[x,y,w,h]`。

## 5. 前端编辑器里怎么填

在 Pipeline 编辑页面：

1. 选择节点。
2. `RECOGNITION` 里把识别方式选成 `Custom`。
3. `识别器名称` 填注册名，例如 `ExampleRecognition`。
4. `识别器参数` 填合法 JSON，例如：

```json
{"mode":"check_balance","threshold":0.8}
```

如果需要 ROI，可以用 ROI 区域旁边的“框选”按钮从当前截图中拖拽获取。

动作同理：

1. `ACTION` 里把动作类型选成 `Custom`。
2. `动作名称` 填注册名，例如 `ExampleAction`。
3. `动作参数` 填合法 JSON。
4. `目标位置` 填 `true`、节点名或坐标。

## 6. 写一个可用的 Custom Recognition

`custom/recognition.go` 当前示例返回 `nil, false`，表示永远识别失败。下面是一个最小可用版本：如果传入 ROI 有宽高，就直接命中 ROI；否则命中一个默认区域。

```go
package custom

import (
    "encoding/json"

    maa "github.com/MaaXYZ/maa-framework-go/v4"
)

type ExampleRecognition struct{}

type ExampleRecognitionParam struct {
    Label string `json:"label"`
}

func (r *ExampleRecognition) Run(ctx *maa.Context, arg *maa.CustomRecognitionArg) (*maa.CustomRecognitionResult, bool) {
    var param ExampleRecognitionParam
    _ = json.Unmarshal([]byte(arg.CustomRecognitionParam), &param)

    box := arg.Roi
    if box.Width() <= 0 || box.Height() <= 0 {
        box = maa.Rect{100, 100, 200, 80}
    }

    detail := map[string]any{
        "label": param.Label,
        "task":  arg.CurrentTaskName,
    }
    detailBytes, _ := json.Marshal(detail)

    return &maa.CustomRecognitionResult{
        Box:    box,
        Detail: string(detailBytes),
    }, true
}
```

如果你要在自定义识别里复用内置 OCR，可以这样做：

```go
detail, err := ctx.RunRecognition("临时OCR", arg.Img, map[string]any{
    "临时OCR": map[string]any{
        "recognition": "OCR",
        "expected":    "确认",
        "roi":         []int{100, 100, 400, 200},
    },
})
if err == nil && detail != nil && detail.Hit {
    return &maa.CustomRecognitionResult{
        Box:    detail.Box,
        Detail: detail.DetailJson,
    }, true
}
return nil, false
```

## 7. 写一个可用的 Custom Action

`custom/action.go` 当前示例返回 `false`，表示动作失败。下面是一个最小可用版本：点击 `target`/识别框中心点。

```go
package custom

import (
    maa "github.com/MaaXYZ/maa-framework-go/v4"
)

type ExampleAction struct{}

func (a *ExampleAction) Run(ctx *maa.Context, arg *maa.CustomActionArg) bool {
    box := arg.Box
    if box.Width() <= 0 || box.Height() <= 0 {
        box = maa.Rect{100, 100, 1, 1}
    }

    x := int32(box.X() + box.Width()/2)
    y := int32(box.Y() + box.Height()/2)

    job := ctx.GetTasker().GetController().PostClick(x, y)
    return job.Wait().Success()
}
```

如果想根据参数控制偏移：

```go
type ExampleActionParam struct {
    Offset []int `json:"offset"`
}
```

在 `Run` 里解析 `arg.CustomActionParam`：

```go
var param ExampleActionParam
_ = json.Unmarshal([]byte(arg.CustomActionParam), &param)

if len(param.Offset) == 2 {
    x += int32(param.Offset[0])
    y += int32(param.Offset[1])
}
```

Pipeline 参数：

```json
{
  "offset": [10, 20]
}
```

## 8. 动态改变流程

Custom 最大的价值之一是可以动态改后续节点。

例如识别到某种状态后，让当前节点跳到 `领取_补货`：

```go
err := ctx.OverrideNext(arg.CurrentTaskName, []maa.NextItem{
    {Name: "领取_补货"},
})
if err != nil {
    return false
}
return true
```

也可以在 Custom 里直接运行一个子任务：

```go
detail, err := ctx.RunTask("领取_点击补货")
if err != nil {
    return false
}
return detail != nil && detail.Status.Success()
```

## 9. 结合你的 `领取` Pipeline 的例子

现在 `resource/pipeline/领取.json` 里有一段是：

```jsonc
{
  "领取_送货上门确定": {
    "recognition": "OCR",
    "expected": "确认",
    "roi": [983, 260, 169, 388],
    "action": "Click"
  }
}
```

如果这个确认按钮逻辑变复杂，比如有时要点确认，有时要先关闭弹窗，可以改成：

```jsonc
{
  "领取_送货上门确定": {
    "recognition": "Custom",
    "custom_recognition": "ExampleRecognition",
    "custom_recognition_param": {
      "label": "送货上门确认"
    },
    "roi": [983, 260, 169, 388],
    "action": "Custom",
    "custom_action": "ExampleAction",
    "custom_action_param": {
      "offset": [0, 0]
    }
  }
}
```

这时流程是：

1. Maa 截图。
2. 把截图、ROI、参数传给 `ExampleRecognition.Run()`。
3. `Run()` 返回一个框和 true，表示识别成功。
4. Maa 把识别框作为 `arg.Box` 传给 `ExampleAction.Run()`。
5. `ExampleAction.Run()` 自己决定点击、跳转、或返回失败。

## 10. 调试方法

推荐用前端的“执行当前节点”调试。现在本项目的执行结果面板会显示：

- task 状态
- 节点列表
- recognition 的 `algorithm`、`hit`、`box`、`detail`
- action 的 `action`、`success`、`box`、`detail`
- 原始 JSON

Custom Recognition 的 `Detail` 建议返回 JSON 字符串，方便结果面板阅读：

```go
detailBytes, _ := json.Marshal(map[string]any{
    "reason": "matched roi",
    "score":  0.92,
})
return &maa.CustomRecognitionResult{
    Box:    box,
    Detail: string(detailBytes),
}, true
```

后端日志位置：

```text
log/
debug/
```

如果 Custom 没有被调用，优先检查：

1. 注册名是否和 Pipeline 完全一致。
2. 是否在 `ReloadResources()` 后重新注册。
3. Pipeline 是否保存成功，并且资源 reload 没有失败。
4. `custom_recognition_param` / `custom_action_param` 是否是合法 JSON。
5. Recognition 返回 `nil, false` 会让节点识别失败；Action 返回 `false` 会让节点动作失败。

## 11. 常见坑

### 名字必须完全一致

Go 注册：

```go
s.resource.RegisterCustomAction("ExampleAction", &custom.ExampleAction{})
```

Pipeline 必须写：

```json
"custom_action": "ExampleAction"
```

大小写不一致也不行。

### 参数是字符串，需要自己 JSON 解析

`arg.CustomRecognitionParam` 和 `arg.CustomActionParam` 在 Go 里是 `string`，不是自动解析好的结构体。

### ROI 必须是四个整数

正确：

```json
"roi": [983, 260, 169, 388]
```

错误：

```json
"roi": [983, 260.169, 388]
```

### ReloadResources 后要保证注册还在

本项目保存 Pipeline 后会调用 `service.Service.ReloadResources()`。如果 `resource.Clear()` 清掉了自定义注册，需要在 reload 后重新注册。最稳的做法是把注册函数放到 `loadResourcesLocked()` 末尾，每次加载资源都执行。

## 12. 推荐落地顺序

1. 在 `service/maa.go` 接入 `registerCustomHandlersLocked()`。
2. 把 `ExampleRecognition` 改成至少能返回 ROI 命中。
3. 把 `ExampleAction` 改成点击 `arg.Box` 中心。
4. 在前端新建一个测试节点，填 `Custom` 识别和 `Custom` 动作。
5. 使用“执行当前节点”看结果面板里的 detail。
6. 再把 Custom 用到真实的 `领取` 流程里。

## 13. 参考文件

- `custom/recognition.go`
- `custom/action.go`
- `service/maa.go`
- `resource/pipeline/领取.json`
- `docs/zh_cn/3.1-任务流水线协议.md`
- `maa-framework-go/custom_recognition.go`
- `maa-framework-go/custom_action.go`
- `maa-framework-go/context.go`
- `maa-framework-go/examples/custom-recognition/main.go`
- `maa-framework-go/examples/custom-action/main.go`
