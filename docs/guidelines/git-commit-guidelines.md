# Git 提交规范

本规范用于约束本仓库的 Git 提交信息。Codex、开发者和其他代码协作工具在创建提交时都应遵守。

## 提交格式

提交信息使用 Conventional Commits 风格：

```text
type(scope): 中文说明
```

示例：

```text
fix(poem): 修复了新增诗词时朝代识别失败的问题
feat(collection): 新增作品集后台管理功能
docs(codex): 添加 Codex 项目协作提示词
```

## 字段规则

- `type` 使用英文小写。
- `scope` 使用英文小写，表示影响范围。
- 冒号后必须空一格。
- 提交信息主体必须使用中文，简洁描述“做了什么”。
- 首行尽量不超过 72 个字符。
- 一个提交只表达一个清晰主题，避免把无关改动混在一起。
- 不使用空泛描述，例如 `update`、`fix bug`、`修改代码`。

## 常用 type

- `feat`：新增功能。
  示例：`feat(collection): 新增作品集后台管理功能`
- `fix`：修复缺陷。
  示例：`fix(poem): 修复了新增诗词时朝代识别失败的问题`
- `docs`：文档变更。
  示例：`docs(codex): 添加 Codex 项目协作提示词`
- `refactor`：重构代码，不改变外部行为。
  示例：`refactor(rbac): 优化权限校验服务结构`
- `test`：新增或修改测试。
  示例：`test(collection): 增加作品集多态关联测试`
- `chore`：工程杂项、依赖、脚本等。
  示例：`chore(deps): 更新前端依赖锁定文件`
- `style`：格式、样式调整，不影响逻辑。
  示例：`style(admin): 调整后台表格按钮间距`
- `perf`：性能优化。
  示例：`perf(poem): 优化诗词列表查询性能`
- `build`：构建系统或依赖变更。
  示例：`build(frontend): 调整 Vite 生产构建配置`
- `ci`：持续集成配置变更。
  示例：`ci(test): 添加后端单元测试流水线`
- `revert`：回滚提交。
  示例：`revert(poem): 回滚诗词导入逻辑调整`

## 常用 scope

- `poem`
- `admin`
- `auth`
- `collection`
- `rbac`
- `frontend`
- `backend`
- `docs`
- `database`
- `config`

## 推荐提交流程

1. 查看状态：`git status --short`
2. 检查差异：`git diff`
3. 运行与改动相关的验证命令。
4. 暂存相关文件。
5. 按本规范编写中文提交信息。

示例：

```bash
git commit -m "fix(poem): 修复了新增诗词时朝代识别失败的问题"
```

## 生效范围

- 对 Codex：只要 `AGENTS.md` 引用本规范，并明确要求提交前遵守，Codex 在本仓库执行提交时应按本规范生成提交信息。
- 对开发者：本规范是项目约定，需要人工遵守。
- 对 Git：文档本身不会自动拦截不合规提交；如果需要强制校验，应额外配置 `commit-msg` hook、Husky、commitlint 或服务端提交校验。
