# 识海前端开发指南

本文档规定识海 React 前端的项目结构、组件、状态、请求、路由、表单和质量要求。业务范围与当前系统边界分别见 [需求设计文档](requirements-design.md) 和 [系统设计文档](system-design.md)。

## 1. 项目概览

前端位于 `frontend/`，使用 React 18、TypeScript 5、Vite 5、React Router、TanStack Query、Zustand、Tailwind CSS、Radix UI 和 lucide-react。

### 1.1 核心目录

```text
frontend/src/
├── components/ # 基础组件、布局和业务组件
├── hooks/      # 数据请求和业务 Hooks
├── pages/      # 页面级组件
├── services/   # HttpClient 和业务接口封装
├── stores/     # Zustand 客户端状态
├── styles/     # 全局样式
├── types/      # 共享 TypeScript 类型
├── utils/      # 通用工具
├── main.tsx    # 应用入口
└── router.tsx  # 路由配置
```

模块内专用的组件、Hook、类型和工具应靠近使用位置；只有被多个业务模块稳定复用的能力才提升到全局目录。新增业务模块不要求机械地同时创建 `components`、`hooks`、`services` 文件。

### 1.2 本地运行与验证

```bash
cd frontend
npm install
npm run dev
npm run build
npm run lint
```

Vite 开发服务器默认使用 `3000` 端口，并将 `/api` 代理到 `VITE_BACKEND_URL`；未设置时使用 `http://localhost:8080`。功能修改至少运行 `npm run build`；涉及 ESLint 规则或大范围代码调整时同时运行 `npm run lint`。

## 2. TypeScript

- 项目开启 `strict`、`noUnusedLocals`、`noUnusedParameters` 和 `noFallthroughCasesInSwitch`。
- 对象类型可以使用 `interface`，联合、映射和组合类型优先使用 `type`；不把这一区分当作机械限制。
- 禁止用 `any` 绕过类型检查；外部未知数据使用 `unknown` 并做收窄。
- 避免双重断言和无依据的非空断言。必须断言时，应能说明运行时保证来自哪里。
- 使用可选链和空值合并处理可空数据，不把合法的 `0`、空字符串或 `false` 误判为缺失。
- 状态常量优先使用 `as const` 对象和联合类型，除非 `enum` 能明显改善与现有代码或外部协议的兼容性。

## 3. 命名与文件

| 类型 | 建议命名 | 示例 |
| --- | --- | --- |
| 页面 | PascalCase + `Page` | `PoemDetailPage.tsx` |
| 业务组件 | PascalCase | `CommentList.tsx` |
| 基础组件 | 遵循现有目录风格 | `button.tsx`、`ConfirmDialog.tsx` |
| Hook | `use` + PascalCase 语义 | `usePoems.ts` |
| Service | camelCase + `Service` | `poemService.ts` |
| Store | camelCase + `Store` | `authStore.ts` |

- 变量和函数使用 camelCase，React 组件使用 PascalCase。
- 布尔值优先使用 `is`、`has`、`can`、`should` 等表达状态的前缀。
- 事件处理函数使用 `handle` 前缀，回调 props 使用 `on` 前缀。
- 文件命名应优先服从所在目录的既有风格，避免仅为统一大小写而重命名大量无关文件。

## 4. 组件

### 4.1 定义方式

组件可以使用函数声明或箭头函数。新代码应在同一文件或同一组件族中保持一致，不使用 Class 组件。

Props 必须有明确类型。公共组件的 Props 类型建议命名为 `{ComponentName}Props`；只在本文件使用的简单 Props 可以内联定义。

页面和业务组件优先使用命名导出，以符合当前路由的 `lazyNamed` 加载方式。基础 UI 组件可以根据 Radix 或现有封装一次导出多个紧密相关的子组件。

### 4.2 职责

- 页面组件负责路由级数据协调和布局。
- 业务组件负责单一业务展示或交互。
- 基础 UI 组件不包含具体业务权限和接口调用。
- Service 负责 HTTP 契约；Hook 负责 Query、Mutation 和页面可复用的状态组合。

组件超过约 200 行时应评估拆分，但行数不是机械上限。高内聚编辑器、数据表封装和生成式基础组件可以保留在单文件中；是否拆分取决于职责数量、状态耦合、可测试性和复用边界。

不要因为两处短小代码相似就立即抽象。只有重复逻辑具有稳定语义、会共同变化或能显著降低复杂度时才提取。

### 4.3 列表与条件渲染

- 列表使用稳定业务 ID 作为 `key`；动态列表禁止使用数组下标。
- 对可能为 `0` 的值使用明确条件，例如 `count > 0`，避免渲染出数字 `0`。
- 加载、空数据、失败和权限不足应是独立状态，不能都回退为空白页面。

### 4.4 弹窗和提示

- 禁止使用浏览器原生 `alert`、`confirm` 和 `prompt`。
- 删除、批量删除和权限变更使用项目已有 `ConfirmDialog` 或等效确认组件。
- 非阻塞反馈使用 `sonner`。
- 危险操作的按钮文案说明具体动作，并使用 destructive 语义。
- 弹窗遮罩必须通过正确的 Portal 覆盖浏览器视口，不能受页面容器或侧栏定位影响。

## 5. Hooks 与渲染性能

- 自定义 Hook 以 `use` 开头，并聚焦一个稳定职责。
- `useEffect` 只处理副作用；依赖数组必须准确，订阅、定时器和异步流程需要清理或取消。
- 不在 `useEffect` 中复制可直接计算的派生状态。
- `useMemo` 和 `useCallback` 只在以下场景使用：计算成本明确、下游依赖引用稳定、Effect 依赖需要稳定，或性能测量证明有收益。
- 不要求所有传给子组件的函数都使用 `useCallback`；缓存本身也有复杂度和成本。
- 避免在父组件中创建会导致大范围页面闪烁的全局加载状态。局部文本转换、字段校验等操作应只更新相关控件状态。

## 6. 状态管理

| 状态 | 工具 | 示例 |
| --- | --- | --- |
| 服务端数据 | TanStack Query | 诗词、评论、用户列表 |
| 跨页面客户端状态 | Zustand | 登录信息、主题 |
| 局部 UI 状态 | `useState` / `useReducer` | 弹窗、输入值、选中项 |

- 服务端数据优先由 TanStack Query 管理，不重复复制到 Zustand。
- Query Key 由资源和完整参数组成，并在同一业务域保持一致。
- 写操作使用 Mutation；成功后精确更新缓存或失效相关 Query，避免无关页面整体刷新。
- `staleTime`、重试和失效范围根据数据变化频率设置，不复制固定数值到所有查询。
- Zustand 持久化只保存恢复会话必需的字段，不持久化可从服务端重新获取的大列表。

## 7. 数据请求

### 7.1 Service 与 HttpClient

- 业务 Service 通过 `src/services/api.ts` 导出的 `api` 调用统一 `HttpClient`。
- 页面和组件不直接创建 Axios 实例或散落 `fetch` 调用。
- Service 负责路径、参数和响应类型，不负责页面提示、导航或复杂业务状态。
- HttpClient 统一附加认证信息、解析 `{ code, message, data }` 包络和处理基础网络错误。
- 全局错误 toast 与局部错误 UI 要避免重复提示；需要局部处理的错误应提供明确的抑制或覆盖机制。

### 7.2 分页与参数

- 当前分页参数使用 `page`、`pageSize`。
- 当前分页响应使用 `list`、`total`、`page`、`pageSize`。
- 可选参数值为 `undefined`、`null` 或空字符串时不发送；合法的 `0` 和 `false` 不得被过滤。
- 前端 Service 路径相对 `/api`，例如 `/poems`，不要重复写 `/api`。

### 7.3 Snowflake ID

- 后端 Snowflake ID 在前端业务状态中应归一化为 `string`。
- 路由参数、表格选择、批量操作、`Map`、`Set`、React `key` 和更新删除参数使用字符串 ID。
- 禁止对 Snowflake ID 使用 `Number()`、`parseInt()` 或一元 `+`。
- 兼容历史类型时可以在接口边界短期使用 `string | number`，但进入业务状态前必须转换为字符串。
- 当前部分共享类型仍将 ID 声明为 `number`，属于待迁移的契约债务；修改相关模块时应同步后端 DTO 和测试逐步收敛，不能据此放宽新代码。

## 8. 表单

- 简单、字段较少且无复杂联动的表单可以使用受控状态。
- 字段较多、校验复杂、需要复用 Schema 或动态字段的表单优先使用 `react-hook-form` 和 `zod`。
- 同一表单不要混用多套状态来源。
- 客户端校验用于及时反馈，后端仍必须执行完整校验。
- 提交期间按钮应禁用并展示稳定的进行中状态，防止重复提交。
- 服务端字段错误应映射到对应字段；无法定位的错误显示为表单级错误。
- 当前代码主要使用受控表单，不能把已安装的表单依赖描述为所有表单已经采用。

## 9. 类型契约

- API 请求和响应类型必须与后端 DTO 对齐，JSON 字段使用 camelCase。
- 新增或修改共享 API 请求和响应类型、跨模块共享的领域类型时，每个属性必须添加 JSDoc 字段注释。
- 局部组件 Props、内部临时类型和测试类型只要求为含义不直观的属性添加注释；自动生成类型由生成源维护注释。
- 涉及单位、格式、枚举取值、可选条件、ID 来源或敏感信息的属性必须详细说明，不能只重复属性名；修改属性语义时必须同步更新注释。
- 写请求类型显式声明，不从完整实体类型通过大范围 `Omit` 机械派生。
- 共享类型放在 `src/types/`；只服务单一 Service 或组件的类型放在对应文件附近。
- 接口返回的可空字段必须在类型中表示为可选或 `null`，不能依赖运行时猜测。
- 调整 DTO 时同时搜索 Service、Hook、表格列、表单默认值和测试数据。

```typescript
interface CreatePoemRequest {
  /** 诗词标题。 */
  title: string;
  /** 按正文行拆分的诗句列表。 */
  content: string[];
  /** 作者的 Snowflake ID，以字符串传输。 */
  authorId: string;
  /** 朝代的 Snowflake ID，以字符串传输。 */
  dynastyId: string;
  /** 译文；未提供时省略。 */
  translation?: string;
}
```

## 10. 路由与权限

- 路由定义集中在 `src/router.tsx`，页面不得自行注册路由表。
- 页面可以使用 `useNavigate` 执行登录跳转、返回和业务导航。
- 内部导航使用 React Router，外部链接使用带 `rel="noopener noreferrer"` 的 `<a>`。
- 登录保护应保存来源位置，以便登录成功后返回。
- 页面和布局可以根据权限显示菜单、按钮或无权限状态，但所有真实授权必须由后端接口执行。
- 页面存在不代表业务已经完成；新增菜单前确认对应 Service 和后端路由可用。

## 11. 样式与布局

- 优先使用 Tailwind 工具类和项目现有 CSS 变量；动态计算值可以使用内联样式。
- 使用 `cn()` 合并条件类名，避免手工字符串拼接产生冲突。
- 主题颜色在 `tailwind.config.js` 或全局 CSS 变量中维护，业务组件避免硬编码散落颜色。
- 使用移动端优先的响应式断点，并检查常见桌面和窄屏视口。
- 固定格式控件应有稳定尺寸约束，加载、悬浮和切换状态不能引起布局跳动。
- 模态框、侧栏和固定导航要检查定位上下文、滚动容器和 z-index，避免遮罩不满视口或控件重叠。
- 图标优先使用 `lucide-react`；不熟悉的纯图标按钮提供 tooltip 和可访问名称。

## 12. 代码质量

- `npm run build` 必须通过 TypeScript 检查和 Vite 构建。
- `npm run lint` 不允许警告，因为脚本启用了 `--max-warnings 0`。
- 使用 `@/` 别名处理跨目录导入；同目录和紧邻模块可以使用相对路径。
- 禁止提交临时 `console.log`、断点和调试 UI。确有运行时告警需求时使用项目约定的日志方案。
- Import 顺序保持清晰：React、第三方、项目绝对路径、同目录相对路径；不要只为排序重排未修改文件。
- 当前项目没有配置前端测试脚本。新增重要纯逻辑或复杂交互时，应先明确测试框架和范围，不在文档中声称测试已经自动执行。

## 13. Git

提交格式、检查流程和分支策略统一遵守 [Git 提交规范](git-commit-guidelines.md)，本文档不重复定义。
