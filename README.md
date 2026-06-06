# Maa Universal Pipeline Editor

一个基于 MaaFramework 的通用自动化项目骨架：Go 后端负责连接窗口、加载资源和执行任务，Vue 前端负责编辑 pipeline、管理图片、选择 ROI、查看执行结果。

项目当前附带的 `resource` 只是示例资源包。把配置指向新的资源目录后，可以用于其他游戏、桌面程序或任意可被 MaaFramework 控制和识别的软件。

## 功能

- 可视化编辑 Maa pipeline JSON
- 上传和枚举模板图片
- 从当前窗口截图中框选 ROI
- 枚举桌面窗口并手动连接
- 执行任务、停止任务、查看任务状态
- 查看 OCR / Detect 模型列表
- 支持 Custom Recognition / Custom Action 扩展
- 支持通过配置文件或环境变量切换资源包、库目录、日志目录和目标窗口

## 配置

默认配置在 `config/app.json`：

```json
{
  "LibDir": "./MaaFramework/bin",
  "ResourceDir": "./resource",
  "LogDir": "./log",
  "TargetWindowTitle": "",
  "AutoConnect": false,
  "ServerAddr": ":8080"
}
```

你可以复制一份配置，然后用环境变量指定：

```powershell
$env:MAA_CONFIG = "D:\my-automation\app.json"
```

常用环境变量：

- `MAA_RESOURCE_DIR`: 资源包目录，包含 `pipeline`、`image`、`ocr`、`detect` 等子目录
- `MAA_LIB_DIR`: MaaFramework 原生库目录
- `MAA_LOG_DIR`: 日志和识别标记图目录
- `MAA_SERVER_ADDR`: Web 服务监听地址，例如 `:8081`
- `MAA_TARGET_WINDOW_TITLE`: 自动连接的窗口标题
- `MAA_AUTO_CONNECT`: 是否自动连接目标窗口，`true` 或 `false`
- `MAA_SAVE_DRAW`: 是否保存识别标记图
- `MAA_DEBUG`: 是否启用 MaaFramework debug 模式
- `MAA_STDOUT_LEVEL`: MaaFramework stdout 日志级别

默认 `AutoConnect` 为 `false`，项目启动后不会绑定任何特定软件。你可以在前端窗口面板里手动选择窗口；如果需要自动连接，再设置 `TargetWindowTitle` 和 `AutoConnect`。

## 资源包结构

一个最小资源包可以这样组织：

```text
resource/
  pipeline/
    Startup.json
  image/
    confirm.png
  ocr/
  detect/
```

如果要适配一个新软件，通常只需要：

1. 新建资源目录。
2. 放入该软件的模板图片、模型和 pipeline。
3. 设置 `MAA_RESOURCE_DIR` 或修改 `config/app.json` 的 `ResourceDir`。
4. 启动项目，在窗口面板选择目标软件窗口。
5. 用编辑器调试节点、ROI 和动作。

## Custom 扩展

`custom/recognition.go` 和 `custom/action.go` 已注册到 Maa Resource：

- `ExampleRecognition`
- `ExampleAction`

它们是最小可用示例：识别器命中 ROI 或默认区域，动作点击目标框中心。你可以复制这两个类型，改成自己的识别/动作逻辑，然后在 `service/maa.go` 的 `registerCustomHandlersLocked` 中注册新名称。

## API

常用接口：

- `GET /api/config`
- `POST /api/resources/reload`
- `GET /api/pipelines`
- `GET /api/pipelines/:name`
- `POST /api/pipelines`
- `PUT /api/pipelines/:name`
- `DELETE /api/pipelines/:name`
- `GET /api/windows`
- `POST /api/windows/connect`
- `GET /api/screenshot`
- `POST /api/tasks`
- `POST /api/tasks/run`
- `POST /api/tasks/stop`

## 开发

前端：

```powershell
cd web
npm install
npm test
npm run build
```

后端：

```powershell
go test ./...
go run .
```

当前环境需要先安装 Go 并确保 `go` 在 `PATH` 中。
