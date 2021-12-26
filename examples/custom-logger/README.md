# custom-logger

## 演示功能

自定义实现 logger，步骤：

- 按照 `log.Logger` 的接口定义，实现自己的 logger，比如写入到文件，或者写入到远程
- 使用 `botgo.SetLogger` 将新的 logger 设置到 sdk 中，sdk 就会使用这个 logger 来写日志了