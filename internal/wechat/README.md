# internal/wechat 目录详细说明

本目录为微信相关功能的核心实现，主要负责微信聊天记录的解密、密钥管理、进程检测、数据结构定义等，支持多平台（Windows/macOS）。

---

## 目录结构

```
internal/wechat/
├── decrypt/      # 聊天记录解密实现（含多平台）
├── key/          # 微信密钥提取与管理（含多平台）
├── model/        # 微信相关数据结构定义
├── process/      # 微信进程检测与管理（含多平台）
├── manager.go    # 微信核心管理器
├── wechat.go     # 微信主流程与接口
```

---

## 主要接口与调用流程

### 1. 微信账号与管理器

#### Account 结构体

表示一个微信账号，包含账号名、平台、版本、数据目录、密钥、进程信息等。

#### Manager 结构体

微信管理器，负责账号、进程的统一管理与调度。

#### 典型调用流程

```go
import (
    "context"
    "github.com/sjzar/chatlog/internal/wechat"
)

// 1. 加载微信进程信息
wechat.Load()

// 2. 获取所有账号
accounts := wechat.GetAccounts()

// 3. 获取指定账号
account, err := wechat.GetAccount("账号名")

// 4. 获取密钥
key, err := account.GetKey(context.Background())

// 5. 解密数据库
err = account.DecryptDatabase(context.Background(), "数据库路径", "输出路径")
```

#### 便捷解密方法

```go
// 通过账号名直接解密数据库
err := wechat.DefaultManager.DecryptDatabase(context.Background(), "账号名", "数据库路径", "输出路径")
```

---

### 2. 解密模块（decrypt/）

#### Decryptor 接口

定义数据库解密的标准接口：

- `Decrypt(ctx, dbfile, key, output)`：解密数据库到输出流
- `Validate(page1, key)`：校验密钥有效性
- `GetPageSize()`、`GetVersion()` 等

#### 解密器工厂

```go
decryptor, err := decrypt.NewDecryptor(platform, version)
```
平台和版本自动适配 Windows/macOS、v3/v4。

#### 通用解密逻辑

`decrypt/common/common.go` 提供跨平台的 AES/HMAC 解密、密钥校验等底层实现。

---

### 3. 密钥提取模块（key/）

#### Extractor 接口

- `Extract(ctx, proc)`：从进程中提取密钥
- `SearchKey(ctx, memory)`：在内存中搜索密钥
- `SetValidate(validator)`：设置密钥校验器

#### 密钥提取器工厂

```go
extractor, err := key.NewExtractor(platform, version)
```
平台和版本自动适配。

---

### 4. 进程检测模块（process/）

#### Detector 接口

- `FindProcesses()`：查找所有微信进程，返回进程信息列表

#### 进程检测器工厂

```go
detector := process.NewDetector(platform)
processes, err := detector.FindProcesses()
```

---

### 5. 数据结构定义（model/）

- `Process`：描述微信进程的结构体，包含 PID、平台、版本、数据目录、账号名等。
- 平台常量：`PlatformWindows`、`PlatformMacOS`
- 状态常量：`StatusOnline`、`StatusOffline`

---

## 详细调用示例

### 获取所有在线微信账号并解密数据库

```go
import (
    "context"
    "github.com/sjzar/chatlog/internal/wechat"
)

func DecryptAllWeChatDBs() {
    // 加载进程信息
    wechat.Load()
    // 获取所有账号
    accounts := wechat.GetAccounts()
    for _, acc := range accounts {
        if acc.Status == "online" {
            // 获取密钥
            key, err := acc.GetKey(context.Background())
            if err != nil {
                // 处理错误
                continue
            }
            // 解密数据库
            err = acc.DecryptDatabase(context.Background(), "数据库路径", "输出路径")
            if err != nil {
                // 处理错误
                continue
            }
            // 处理解密后的数据
        }
    }
}
```

---

## 子模块功能补充说明

- **decrypt/windows/**、**decrypt/darwin/**：分别实现 Windows/macOS 下 v3/v4 版本数据库的解密算法。
- **key/windows/**、**key/darwin/**：分别实现 Windows/macOS 下 v3/v4 版本密钥的提取算法。
- **key/darwin/glance/**：macOS 下的快速密钥扫描工具。
- **process/windows/**、**process/darwin/**：平台相关的进程检测实现。
- **decrypt/common/**：跨平台的 AES/HMAC 解密、密钥校验等底层实现。

---

## 适用场景

- 微信聊天记录的导出、分析、迁移
- 微信数据库的解密与数据恢复
- 多平台（Windows/macOS）微信数据处理工具开发

---

> 如需了解具体实现细节，请查阅各子目录下的源码及注释，或直接调用上述接口进行二次开发。 