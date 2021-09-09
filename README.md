用于放置 git 本地 hooks 脚本.

### 规则

1. commit log 格式：校验是否有携带了 story/bug/task/other 与相关ID
2. 校验是否符合 Angular CommitLog 规范
3. 分支上的提交，自动根据分支名，在 commit log 上补充 story/bug/task/other 等相关信息
4. 自动补充type与scope。type 默认使用 feat，scope 根据变更文件列表计算
5. 提交 commit 的时候，自动识别当前路径下是否存在 golang 文件，如果存在，则执行 golangci-lint run

### 建议使用步骤

#### 全局模式（推荐）

```
git clone http://git.code.oa.com/epc_tools/git_hooks.git .git_hooks
git config --global core.hooksPath `pwd`/.git_hooks
```

注意：请确保本地git版本高于 2.9 才能够自定义 hooksPath