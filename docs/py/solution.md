# Chatlog Python 版项目解决方案

## 1. 项目简介与目标

Chatlog 是一个用于解析、解密和分析微信聊天记录的工具，支持多平台（Windows/macOS）、多微信版本，具备命令行、HTTP API、Web UI、MCP 协议等多种交互方式。Python 版目标是：
- 复刻并优化现有 Go 版核心功能
- 保持跨平台兼容性
- 提供良好的可扩展性和可维护性
- 充分利用 Python 生态的丰富库

---

## 2. 总体架构设计

采用分层架构，核心分为：
- **命令行/终端 UI 层**：CLI 工具，支持数据解密、服务启动等
- **Web 层**：HTTP API + Web 前端，支持数据查询、MCP 集成
- **核心业务层**：账号管理、解密、数据库操作、会话与消息处理
- **工具与适配层**：文件监控、压缩/解压、音频编解码、配置管理等

---

## 3. 主要功能模块划分

1. **cli/**         —— 命令行入口与参数解析
2. **core/**        —— 业务核心（账号、解密、数据库、消息等）
3. **web/**         —— HTTP API 服务与 Web UI
4. **wechat/**      —— 微信协议、密钥、解密、平台适配
5. **db/**          —— 微信数据库适配与抽象
6. **utils/**       —— 工具包（压缩、音频、文件监控等）
7. **config/**      —— 配置管理
8. **static/**      —— 前端静态资源
9. **scripts/**     —— 构建、迁移、测试脚本

---

## 4. 技术选型与依赖建议

- **Web 框架**：FastAPI（高性能、类型友好、自动文档）
- **前端**：复用原有 HTML/CSS/JS，或采用 Streamlit/Gradio 进行快速原型
- **数据库**：sqlite3（内置支持，兼容微信数据库）
- **文件监控**：watchdog
- **压缩/解压**：lz4、zstandard（zstd）、zipfile（内置）
- **音频编解码**：pysilk、pydub、ffmpeg-python
- **配置管理**：pydantic/yaml/json
- **命令行**：typer/click
- **日志**：loguru
- **多平台支持**：platform、pywin32（Windows）、pyobjc（macOS）

---

## 5. 目录结构建议

```text
chatlog_py/
  cli/                # 命令行工具
  core/               # 业务核心
  web/                # HTTP API & Web UI
  wechat/             # 微信协议/解密/密钥
  db/                 # 微信数据库适配
  utils/              # 工具包
  config/             # 配置
  static/             # 前端静态资源
  scripts/            # 辅助脚本
  main.py             # 主入口
  requirements.txt    # 依赖
  README.md           # 说明文档
```

---

## 6. 各模块详细说明

### 6.1 cli/
- 负责命令行参数解析、子命令分发（如 key、decrypt、server 等）
- 推荐使用 typer，支持丰富的 CLI 体验

### 6.2 core/
- 账号管理、会话管理、消息抽象、业务流程控制
- 负责 orchestrate 各子模块

### 6.3 web/
- FastAPI 实现 RESTful API，路由与参数校验
- 提供 /api/v1/chatlog、/api/v1/contact、/api/v1/chatroom、/api/v1/session 等接口
- 提供 /sse（MCP 协议）接口
- 静态资源托管（static/）

### 6.4 wechat/
- 微信密钥提取、数据库解密、平台适配（Windows/macOS）
- 兼容多版本微信
- 可用 ctypes/pywin32/pyobjc 调用系统 API

### 6.5 db/
- 微信数据库（sqlite）适配与抽象
- 支持多平台/多版本差异
- 提供统一数据访问接口

### 6.6 utils/
- 文件监控（watchdog）
- 压缩/解压（lz4、zstd、zipfile）
- 音频编解码（pysilk、pydub、ffmpeg-python）
- 通用工具函数

### 6.7 config/
- 配置文件加载与校验（pydantic/yaml/json）
- 支持多环境配置

### 6.8 static/
- 前端 HTML/CSS/JS 资源，建议复用原有 index.htm

### 6.9 scripts/
- 构建、迁移、测试等辅助脚本

---

## 7. 关键设计要点与注意事项

- **解密与密钥提取**：需针对 Windows/macOS 分别实现，优先用 Python 原生/第三方库，必要时用 C/C++ 扩展
- **数据库兼容性**：微信不同版本表结构有差异，需适配
- **多媒体处理**：图片/语音/视频需支持解密与转码
- **API 设计**：RESTful，参数校验，错误处理友好
- **前后端分离**：Web UI 可独立开发，API 只负责数据
- **自动化测试**：建议覆盖核心解密、数据访问、API
- **跨平台**：平台相关代码需抽象隔离，便于维护

---

## 8. 未来可扩展性建议

- 支持更多聊天平台（如 QQ、钉钉等）
- 增加全文索引与搜索（Whoosh/Elasticsearch）
- 数据统计与可视化 Dashboard
- 插件化架构，便于第三方扩展
- 云端部署与多用户支持

---

## 9. 参考依赖（requirements.txt 示例）

```text
fastapi
uvicorn
watchdog
lz4
zstandard
pysilk
pydub
ffmpeg-python
pydantic
typer
loguru
sqlite3  # Python 标准库
platform
pywin32; platform_system=="Windows"
pyobjc; platform_system=="Darwin"
```

---

如需详细模块设计、接口定义或样例代码，可进一步细化。 