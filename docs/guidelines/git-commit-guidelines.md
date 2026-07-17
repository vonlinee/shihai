# Git 提交与分支规范

本规范约束识海仓库的提交信息、提交检查和分支使用。只有用户明确要求提交时，Codex 才创建提交。

## 1. 提交格式

提交首行使用：

```text
type(scope): 中文说明
```

示例：

```text
fix(poem): 修复新增诗词时朝代识别失败的问题
feat(collection): 新增作品集后台管理功能
docs(guidelines): 重构项目开发指南体系
```

## 2. 字段规则

- `type` 使用英文小写。
- `scope` 使用英文小写，表达业务或技术影响范围。
- 冒号后保留一个空格。
- 提交说明使用中文，直接描述完成的变更。
- 首行尽量控制在 72 个字符以内。
- 一个提交只表达一个清晰主题，不混入无关格式化或重构。
- 禁止使用 `update`、`fix bug`、`修改代码` 等无法说明结果的描述。

复杂变更可以增加正文，说明背景、关键取舍和验证方式。破坏性变更必须在正文或 Footer 中使用 `BREAKING CHANGE:` 明确说明影响和迁移方式。

## 3. 常用 type

| type | 用途 | 示例 |
| --- | --- | --- |
| `feat` | 新增功能 | `feat(poem): 新增诗词标注编辑功能` |
| `fix` | 修复缺陷 | `fix(comment): 修复后台评论列表查询失败` |
| `docs` | 文档变更 | `docs(guidelines): 完善前端开发规范` |
| `refactor` | 不改变外部行为的重构 | `refactor(rbac): 整理权限校验服务` |
| `test` | 新增或修改测试 | `test(poem): 增加诗词 ID 解析测试` |
| `perf` | 性能优化 | `perf(poem): 优化诗词列表查询` |
| `style` | 不影响逻辑的样式或格式 | `style(admin): 调整后台表格间距` |
| `build` | 构建系统或依赖 | `build(frontend): 调整 Vite 构建配置` |
| `ci` | 持续集成 | `ci(test): 增加后端测试任务` |
| `chore` | 工程杂项 | `chore(deps): 更新前端依赖锁定文件` |
| `revert` | 回滚已有提交 | `revert(poem): 回滚诗词导入调整` |

## 4. 常用 scope

常用范围包括：

- 业务：`poem`、`comment`、`collection`、`auth`、`rbac`、`admin`。
- 技术：`frontend`、`backend`、`database`、`config`、`deps`、`ci`。
- 文档：`docs`、`guidelines`。

列表不是白名单。新增 scope 时优先使用仓库已有的稳定业务名，避免为单个文件发明过细范围。

## 5. 提交流程

1. 查看工作区：`git status --short`。
2. 检查未暂存差异：`git diff`。
3. 运行与改动相关的验证命令。
4. 只暂存当前主题需要的文件。
5. 再次查看状态：`git status --short`。
6. 检查暂存差异：`git diff --cached`。
7. 检查空白错误：`git diff --cached --check`。
8. 按本规范创建提交。

工作区存在他人或用户改动时，不使用 `git add -A` 混入无关文件；明确列出需要暂存的路径。

## 6. 分支策略

仓库当前分支为：

| 分支 | 用途 |
| --- | --- |
| `master` | 稳定分支 |
| `dev` | 开发集成分支 |

- 普通功能和修复默认从目标集成分支 `dev` 创建分支。
- 分支名称使用清晰前缀，例如 `feat/poem-annotation`、`fix/comment-list`；Codex 创建分支时按运行环境要求使用 `codex/` 前缀。
- 不直接向稳定分支推送未经评审的开发改动。
- 合并目标根据发布流程选择，不能因为文档示例而假设远端一定存在其他分支。
- 若仓库迁移到新的默认分支名，必须同步更新远端默认分支、保护规则、自动化脚本和本文档。

## 7. 生效方式

- 对 Codex：根目录 `AGENTS.md` 引用本规范后，创建提交时必须遵守。
- 对开发者：本规范依靠代码审查和团队协作执行。
- 对 Git：Markdown 文档不会自动阻止不合规提交；需要强制校验时，应另行配置 `commit-msg` hook、commitlint 或服务端规则。
