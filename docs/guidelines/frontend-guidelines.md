# 识海（shihai）前端开发规范文档

> 本文档规定了识海古诗词学习平台前端项目的开发规范，适用于所有参与前端开发的成员。技术栈：React 18 + TypeScript 5 + Vite 5 + TailwindCSS 3 + Zustand + TanStack Query。

---

## 项目概览

识海前端基于 React、TypeScript、Vite、Tailwind CSS、TanStack Query、Radix UI 和 lucide-react 构建。

### 技术栈

- React
- TypeScript
- Vite
- React Router
- TanStack Query
- Tailwind CSS
- Radix UI
- lucide-react

### 核心目录

```text
frontend/
├── src/components/ # 通用组件和布局组件
├── src/hooks/      # React hooks 和数据请求 hooks
├── src/pages/      # 页面级组件
├── src/services/   # API 客户端和服务封装
├── src/stores/     # 前端状态管理
├── src/styles/     # 全局样式
├── src/types/      # TypeScript 类型定义
├── src/utils/      # 通用工具
├── src/main.tsx    # 应用入口
└── src/router.tsx  # 路由配置
```

### 开发约定

- 页面放在 `src/pages/`，可复用 UI 放在 `src/components/`。
- API 调用集中在 `src/services/`，页面通过 hooks 或服务层调用接口。
- 路由变更同步更新 `src/router.tsx`。
- 类型优先维护在 `src/types/` 或相关服务文件中，避免重复定义。
- 后台管理页面应保持信息密度、可扫描性和稳定交互，不做营销式落地页。

### 本地运行

```bash
cd frontend
npm install
npm run dev
```

Vite 开发服务器默认端口为 `3000`，`/api` 会代理到 `http://localhost:8080`。

### 常用验证

```bash
cd frontend
npm run build
```

如果修改涉及代码质量规则，可额外运行：

```bash
cd frontend
npm run lint
```

### 修改前阅读

- 新增或调整前端页面、组件、路由、API 调用：阅读本文档。
- 修改后台管理功能：额外阅读 [系统设计文档](system-design.md) 中的后台和权限相关内容。
- 修改跨端功能：同时阅读 [后端指南](backend-guidelines.md)、[Go 开发规范](go-development-guidelines.md) 和后端相关规范。

---

## 目录

1. [项目结构规范](#一项目结构规范)
2. [TypeScript 规范](#二typescript-规范)
3. [命名规范](#三命名规范)
4. [组件规范](#四组件规范)
5. [Hooks 规范](#五hooks-规范)
6. [状态管理规范](#六状态管理规范)
7. [数据请求规范](#七数据请求规范)
8. [样式规范](#八样式规范)
9. [路由规范](#九路由规范)
10. [表单规范](#十表单规范)
11. [类型定义规范](#十一类型定义规范)
12. [代码质量规范](#十二代码质量规范)
13. [Git 提交与分支管理](#十三git-提交与分支管理)

---

## 一、项目结构规范

```
src/
├── components/          # 可复用组件
│   ├── ui/              # 基础 UI 组件（shadcn/ui 生成，一般不手动修改）
│   ├── layout/          # 布局组件（Header、Sidebar、Footer 等）
│   ├── poetry/          # 诗词业务组件
│   ├── comment/         # 评论业务组件
│   └── admin/           # 管理后台业务组件
├── pages/               # 页面级组件（与路由一一对应）
│   └── admin/           # 管理后台页面
├── hooks/               # 自定义 React Hooks
├── services/            # API 请求封装
│   ├── httpClient.ts    # HTTP 客户端抽象层
│   ├── api.ts           # API 基础实例
│   ├── authService.ts   # 认证相关请求
│   ├── poemService.ts   # 诗词相关请求
│   └── index.ts         # 统一导出
├── stores/              # Zustand 全局状态
├── types/               # TypeScript 类型定义
│   └── index.ts         # 统一类型导出
├── utils/               # 工具函数
├── styles/              # 全局样式
├── router.tsx           # 路由配置
└── main.tsx             # 应用入口
```

**规则：**

- `pages/` 下只放页面组件，每个路由对应一个文件，**不放可复用组件**。
- `components/` 按业务域划分子目录，`ui/` 目录为基础组件，**不在其中添加业务逻辑**。
- `hooks/` 只放 Custom Hooks，文件名以 `use` 开头。
- `services/` 只负责网络请求，不包含业务逻辑判断。
- `stores/` 每个状态域对应一个文件，以 `Store` 命名（如 `authStore.ts`）。
- 新增业务模块时，在 `components/`、`hooks/`、`services/` 下同步创建对应文件，保持各层完整。

---

## 二、TypeScript 规范

项目开启 `strict` 模式，以下规则强制执行。

### 2.1 类型声明

- 优先使用 `interface` 声明对象类型，使用 `type` 声明联合类型、交叉类型等。
- 禁止使用 `any`，必要时使用 `unknown` 并做类型收窄。
- 函数必须有明确的参数类型和返回值类型，React 组件 props 除外（可由 TypeScript 推断）。

```typescript
// 正确
interface User {
  id: number;
  username: string;
}

type Status = 'pending' | 'approved' | 'rejected';

function getUserName(user: User): string {
  return user.username;
}

// 错误
function getUserName(user: any) {
  return user.username;
}
```

### 2.2 类型断言

- 尽量避免类型断言（`as`），优先通过类型守卫收窄。
- 禁止使用双重断言（`as unknown as XXX`）绕过类型检查，除非有充分理由并添加注释说明。

```typescript
// 正确：类型守卫
if (typeof value === 'string') {
  console.log(value.toUpperCase());
}

// 尽量避免
const user = response as User;
```

### 2.3 可选链与非空断言

- 优先使用可选链（`?.`）和空值合并（`??`）处理可能为空的值。
- 禁止滥用非空断言（`!`），只在确定不为 null/undefined 时使用，并添加注释。

```typescript
// 正确
const name = user?.profile?.name ?? '匿名';

// 需要注释说明
const token = localStorage.getItem('token')!; // 此处由登录拦截保证必然存在
```

### 2.4 枚举

- 使用 `const` 对象替代 `enum`，避免枚举的编译产物膨胀问题。

```typescript
// 推荐
const CorrectionStatus = {
  Pending: 'pending',
  Voting: 'voting',
  Approved: 'approved',
  Rejected: 'rejected',
} as const;
type CorrectionStatus = typeof CorrectionStatus[keyof typeof CorrectionStatus];

// 避免
enum CorrectionStatus {
  Pending = 'pending',
  ...
}
```

---

## 三、命名规范

### 3.1 文件命名

| 类型 | 命名格式 | 示例 |
|------|----------|------|
| 页面组件 | PascalCase + `Page` 后缀 | `PoemDetailPage.tsx` |
| 业务组件 | PascalCase | `PoemCard.tsx`, `CommentList.tsx` |
| 基础 UI 组件 | PascalCase | `Button.tsx`, `Input.tsx` |
| Custom Hook | camelCase + `use` 前缀 | `usePoems.ts`, `useAuth.ts` |
| Service | camelCase + `Service` 后缀 | `poemService.ts` |
| Store | camelCase + `Store` 后缀 | `authStore.ts` |
| 类型文件 | camelCase | `index.ts` |
| 工具函数 | camelCase | `cn.ts`, `format.ts` |

### 3.2 变量与函数

- 变量、函数使用 camelCase。
- 常量使用 UPPER_SNAKE_CASE（模块级配置常量）或 PascalCase（对象常量）。
- React 组件使用 PascalCase。
- 布尔值变量以 `is`、`has`、`can`、`should` 开头。

```typescript
// 变量
const poemList: Poem[] = [];
const isLoading = false;
const hasPermission = true;

// 常量
const API_BASE_URL = '/api/v1';
const MAX_COMMENT_LENGTH = 500;

// 组件
function PoemCard({ poem }: { poem: Poem }) { ... }
```

### 3.3 事件处理函数

事件处理函数以 `handle` 开头，描述所处理的事件：

```typescript
const handleSubmit = (e: React.FormEvent) => { ... };
const handleCommentLike = (commentId: number) => { ... };
const handlePageChange = (page: number) => { ... };
```

---

## 四、组件规范

### 4.1 组件定义

- 统一使用函数组件 + 箭头函数形式，禁止使用 Class 组件。
- Props 类型用 `interface` 定义，命名为 `{ComponentName}Props`。
- 组件文件只导出一个主组件（允许辅助子组件在同文件内定义但不导出）。

```tsx
interface PoemCardProps {
  poem: Poem;
  onFavorite?: (id: number) => void;
  className?: string;
}

const PoemCard = ({ poem, onFavorite, className }: PoemCardProps) => {
  return (
    <div className={cn('rounded-lg border p-4', className)}>
      <h3 className="font-serif text-lg">{poem.title}</h3>
      {/* ... */}
    </div>
  );
};

export default PoemCard;
```

### 4.2 组件职责

- 页面组件（`pages/`）：负责数据获取、状态协调、布局，调用业务组件。
- 业务组件（`components/{domain}/`）：关注单一业务展示或交互，通过 props 接收数据。
- 基础组件（`components/ui/`）：无业务逻辑，只负责 UI 渲染，完全由 props 驱动。

```
页面组件（数据获取 + 布局）
    ↓
业务组件（业务展示 + 交互）
    ↓
基础 UI 组件（纯 UI 渲染）
```

### 4.3 组件拆分原则

- 单个组件不超过 **200 行**，超出时拆分子组件。
- 当同一逻辑在 2 个以上组件中重复，提取为 Custom Hook 或通用组件。
- 组件只做一件事，避免将多个不相关功能塞入同一组件。

### 4.4 Key 规范

列表渲染必须提供稳定的 `key`，禁止使用数组下标作为 key（动态增删时）：

```tsx
// 正确
poems.map((poem) => <PoemCard key={poem.id} poem={poem} />)

// 错误
poems.map((poem, index) => <PoemCard key={index} poem={poem} />)
```

### 4.5 条件渲染

对于可能为 `0` 的数字，避免直接用 `&&` 短路渲染：

```tsx
// 正确
{count > 0 && <Badge count={count} />}

// 错误（count 为 0 时会渲染数字 0）
{count && <Badge count={count} />}
```

### 4.6 弹窗与提示

- 禁止使用浏览器原生 `alert`、`confirm`、`prompt` 以及 `window.alert`、`window.confirm`、`window.prompt`。
- 删除、批量删除、权限变更等确认类操作，统一使用项目已有弹窗组件，例如 `src/components/ui/ConfirmDialog.tsx`。
- 成功、失败、信息类非阻塞提示统一使用项目已有 toast 方案，例如 `sonner`。
- 危险操作的确认按钮文案必须明确，样式使用 destructive 语义，避免只写“确定”。
- 新增页面或交互提交前，应检查是否引入了原生弹窗调用。

---

## 五、Hooks 规范

### 5.1 自定义 Hook

- 必须以 `use` 开头（ESLint 强制）。
- 一个 Hook 只聚焦一个职责（如数据获取、表单状态、权限检查）。
- 对外暴露的返回值结构要稳定，不要随意增删字段。

```typescript
// hooks/usePoems.ts
export function usePoems(params: PoemSearchParams) {
  const { data, isLoading, error } = useQuery({
    queryKey: ['poems', params],
    queryFn: () => poemService.getList(params),
  });

  return {
    poems: data?.items ?? [],
    total: data?.total ?? 0,
    isLoading,
    error,
  };
}
```

### 5.2 useEffect 规范

- 必须指定正确的依赖数组，禁止为了消除 warning 而乱填依赖。
- 有副作用清理需求（订阅、定时器）时必须返回清理函数。
- 避免在 `useEffect` 中直接更新 state 并再次触发 effect（会导致无限循环）。

```typescript
// 正确：带清理的 effect
useEffect(() => {
  const timer = setInterval(() => fetchData(), 30000);
  return () => clearInterval(timer);
}, [fetchData]);
```

### 5.3 useMemo 和 useCallback

- 仅在有明确性能瓶颈时使用，不要过度优化。
- 传给子组件的函数 prop 用 `useCallback` 包裹，避免不必要的子组件重渲染。
- 计算量较大的派生值用 `useMemo` 缓存。

---

## 六、状态管理规范

项目使用 Zustand 管理全局状态，TanStack Query 管理服务端状态。

### 6.1 状态分类

| 状态类型 | 工具 | 示例 |
|----------|------|------|
| 服务端数据（列表、详情） | TanStack Query | 诗词列表、评论列表 |
| 全局客户端状态 | Zustand | 登录用户信息、主题 |
| 局部 UI 状态 | useState / useReducer | 弹窗开关、表单值 |

**原则：服务端数据优先交给 TanStack Query 管理，不要把接口数据手动存入 Zustand。**

### 6.2 Zustand Store 规范

- 每个 Store 文件对应一个业务域，命名 `{domain}Store.ts`。
- 使用 `interface` 定义 State 类型，包含数据字段和操作方法。
- 需要持久化的状态使用 `persist` 中间件，通过 `partialize` 只持久化必要字段。

```typescript
// stores/authStore.ts
interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (user: User, token: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      login: (user, token) => set({ user, token, isAuthenticated: true }),
      logout: () => set({ user: null, token: null, isAuthenticated: false }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ user: state.user, token: state.token }),
    }
  )
);
```

### 6.3 TanStack Query 规范

- `queryKey` 保持一致性，按 `['资源', 参数对象]` 格式组织。
- 增删改操作使用 `useMutation`，成功后通过 `invalidateQueries` 刷新相关数据。
- 设置合理的 `staleTime`，避免频繁的无意义请求。

```typescript
// 查询
const { data } = useQuery({
  queryKey: ['poems', { page, pageSize, keyword }],
  queryFn: () => poemService.getList({ page, pageSize, keyword }),
  staleTime: 60 * 1000, // 1 分钟内不重新请求
});

// 变更
const { mutate: createComment } = useMutation({
  mutationFn: commentService.create,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['comments', poemId] });
    toast.success('评论发布成功');
  },
});
```

---

## 七、数据请求规范

### 7.1 Service 层规范

- 每个业务域对应一个 Service 文件，只负责 HTTP 请求和参数/响应映射。
- 不在 Service 中做业务逻辑判断（如权限判断、数据聚合）。
- 通过 `httpClient` 抽象层发起请求，不直接使用 `axios` 或 `fetch`。

```typescript
// services/poemService.ts
import { httpClient } from './api';
import type { Poem, PoemSearchParams, PaginatedResponse } from '@/types';

export const poemService = {
  getList: (params: PoemSearchParams) =>
    httpClient.get<PaginatedResponse<Poem>>('/v1/poems', params),

  getById: (id: number) =>
    httpClient.get<Poem>(`/v1/poems/${id}`),

  create: (data: CreatePoemRequest) =>
    httpClient.post<Poem>('/v1/poems', data),

  update: (id: number, data: UpdatePoemRequest) =>
    httpClient.put<Poem>(`/v1/poems/${id}`, data),

  delete: (id: number) =>
    httpClient.delete<void>(`/v1/poems/${id}`),
};
```

### 7.2 错误处理

- HTTP 客户端（`httpClient.ts`）统一处理 401 跳转登录、toast 错误提示。
- 组件层通过 TanStack Query 的 `error` 状态渲染错误 UI，不需要重复 try/catch。
- 对于需要自定义错误处理的场景，在 `useMutation.onError` 中处理。

### 7.3 请求参数

- 列表接口分页参数统一使用 `page`（从 1 开始）和 `pageSize`。
- 空字符串、undefined、null 的可选参数不传给后端（在 httpClient 中过滤）。

### 7.4 雪花 ID 处理

- 后端响应中的雪花 ID 按字符串处理，前端类型中的 ID 优先声明为 `string`；兼容历史数据时可临时使用 `string | number`，但进入业务状态前应归一化为字符串。
- 表格选择、批量操作、`Map`/`Set` key、React `key`、路由参数和删除/更新接口参数都必须使用字符串 ID。
- 禁止对雪花 ID 使用 `Number()`、`parseInt()`、一元 `+` 等方式强转为数字，除非能确认该字段不是雪花 ID。
- 向后端提交 ID 时优先传字符串，避免超过 `Number.MAX_SAFE_INTEGER` 后出现精度丢失、误选多行或删除错误数据。

---

## 八、样式规范

### 8.1 TailwindCSS 使用规范

- 优先使用 Tailwind 工具类，禁止内联 `style` 属性（动态值除外）。
- 使用 `cn()` 工具函数合并类名，支持条件类名：

```tsx
import { cn } from '@/utils/cn';

<div className={cn(
  'rounded-lg border p-4',
  isActive && 'border-primary bg-primary/10',
  className  // 允许外部传入的 className
)} />
```

- 自定义主题颜色通过 `tailwind.config.js` 的 `theme.extend` 统一定义，禁止在组件中硬编码颜色值（如 `#C41E3A`）。

### 8.2 响应式设计

- 使用 Tailwind 响应式前缀，遵循移动端优先原则：`sm:` / `md:` / `lg:` / `xl:`。

```tsx
<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
```

### 8.3 颜色使用

项目定义了古典中国风主题色，优先使用语义化颜色变量：

| 颜色变量 | 用途 |
|----------|------|
| `ink` | 主文字、标题、强调色 |
| `paper` | 背景、卡片底色 |
| `cinnabar` | 操作按钮、高亮、警示 |
| `primary` / `secondary` | 通用主色调 |

```tsx
// 正确
<h1 className="text-ink font-serif">静夜思</h1>

// 错误
<h1 style={{ color: '#8B4513' }}>静夜思</h1>
```

### 8.4 字体规范

- 诗词标题、诗词正文使用 `font-serif`（宋体/楷体风格）。
- 普通 UI 文本使用 `font-sans`（默认）。

### 8.5 动画

- 优先使用 `tailwindcss-animate` 提供的动画类（`animate-fade-in`、`animate-slide-up`）。
- 自定义动画在 `tailwind.config.js` 的 `keyframes` 中定义，不写内联 style。

---

## 九、路由规范

### 9.1 路由配置

- 所有路由在 `router.tsx` 中统一配置，不在页面组件内部跳转路由。
- 路由路径使用 kebab-case，与资源名一致。

```typescript
// 路径规范示例
/poems              // 诗词列表
/poems/:id          // 诗词详情
/user-center        // 个人中心
/admin/users        // 管理后台用户管理
```

### 9.2 路由守卫

- 需要登录的页面通过路由守卫拦截，未登录时重定向到 `/login`，并携带 `from` 参数（登录后可跳回）。
- 需要特定角色的页面（如管理后台）通过 `useRbac` hook 在页面组件中做权限判断，无权时展示无权限页面。

```tsx
// 路由守卫示例
const ProtectedRoute = ({ children }: { children: ReactNode }) => {
  const { isAuthenticated } = useAuthStore();
  const location = useLocation();
  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }
  return <>{children}</>;
};
```

### 9.3 导航

- 使用 `useNavigate` 进行编程式导航，不使用 `window.location.href`（认证流程的 401 重定向除外）。
- 外部链接使用 `<a target="_blank" rel="noopener noreferrer">`。

---

## 十、表单规范

### 10.1 表单库

- 所有表单使用 `react-hook-form` + `zod` 的组合实现：react-hook-form 管理表单状态，zod 做校验 Schema。
- 通过 `@hookform/resolvers/zod` 将 zod schema 接入 react-hook-form。

```tsx
const loginSchema = z.object({
  username: z.string().min(1, '请输入用户名'),
  password: z.string().min(6, '密码至少 6 位'),
});

type LoginFormData = z.infer<typeof loginSchema>;

const { register, handleSubmit, formState: { errors } } = useForm<LoginFormData>({
  resolver: zodResolver(loginSchema),
});
```

### 10.2 校验规则

- 校验规则统一在 zod schema 中定义，不在 JSX 中直接写 `validate` 函数。
- 错误提示信息使用中文，语义清晰。
- 必填字段使用 `.min(1, '提示')` 而非 `.nonempty()`，保持风格一致。

### 10.3 提交状态

- 提交按钮在请求进行中必须 disabled 并显示 loading 状态，防止重复提交。

```tsx
<Button type="submit" disabled={isSubmitting}>
  {isSubmitting ? '提交中...' : '提交'}
</Button>
```

---

## 十一、类型定义规范

### 11.1 类型文件组织

- 全局共用类型统一在 `src/types/index.ts` 中定义和导出。
- 仅单个组件/模块使用的局部类型，在该文件内定义，不放到全局类型文件。

### 11.2 API 类型

- 接口请求/响应类型必须与后端 DTO 保持一致，字段名使用 camelCase（后端 JSON 已统一为 camelCase）。
- 分页响应统一使用 `PaginatedResponse<T>` 泛型。

```typescript
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}
```

### 11.3 类型复用

- 优先使用 `Pick`、`Omit`、`Partial` 等工具类型基于现有类型派生新类型，避免重复定义。

```typescript
// 创建请求类型从实体类型派生
type CreatePoemRequest = Omit<Poem, 'id' | 'createdAt' | 'updatedAt' | 'views' | 'likes'>;
type UpdatePoemRequest = Partial<CreatePoemRequest>;
```

---

## 十二、代码质量规范

### 12.1 ESLint

- 项目配置了 ESLint（TypeScript + React Hooks + React Refresh 规则），**提交前代码必须通过 lint**。
- `lint` 命令设置了 `--max-warnings 0`，即不允许任何警告存在于提交代码中。

```bash
npm run lint        # 检查
```

### 12.2 TypeScript 严格模式

`tsconfig.json` 开启了以下严格选项，必须遵守：

- `strict: true` — 所有严格类型检查
- `noUnusedLocals: true` — 禁止未使用的局部变量
- `noUnusedParameters: true` — 禁止未使用的函数参数
- `noFallthroughCasesInSwitch: true` — switch 禁止 fallthrough

### 12.3 路径别名

- 使用 `@/` 路径别名替代相对路径导入（tsconfig 已配置 `@/* → ./src/*`）。
- 跨目录引用（超过两层 `../`）必须使用 `@/` 别名。

```typescript
// 正确
import { useAuthStore } from '@/stores/authStore';
import type { Poem } from '@/types';

// 错误
import { useAuthStore } from '../../stores/authStore';
```

### 12.4 Import 顺序

按以下顺序分组，组间空行：

```typescript
// 1. React 相关
import { useState, useEffect } from 'react';

// 2. 第三方库
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';

// 3. 项目内绝对路径（@/）
import { useAuthStore } from '@/stores/authStore';
import { poemService } from '@/services';
import type { Poem } from '@/types';

// 4. 相对路径（同目录组件等）
import PoemCard from './PoemCard';
```

### 12.5 Console 和调试代码

- 禁止在提交代码中保留 `console.log`、`console.warn` 调试输出。
- 临时调试完成后必须清理，提交前执行 lint 检查。

### 12.6 组件导出

- 页面组件使用 `export default`。
- 业务组件和工具函数优先使用命名导出（`export const`），便于 tree-shaking 和 IDE 自动导入。

---

## 十三、Git 提交与分支管理

### 提交信息

提交信息的格式、类型和示例以 [Git 提交规范](git-commit-guidelines.md) 为准。

前端相关提交的 `scope` 可以使用业务模块或前端模块名，例如 `frontend`、`admin`、`auth`、`poem-list`、`poem-detail`、`http-client`。

### 分支管理

| 分支 | 用途 |
|------|------|
| `main` | 生产稳定版本，只接受 PR 合并 |
| `develop` | 前端开发集成分支 |
| `feat/xxx` | 功能开发分支，从 `develop` 拉取 |
| `fix/xxx` | Bug 修复分支 |

- 禁止直接向 `main` 推送代码。
- 功能开发完成后通过 PR 合并到 `develop`，经测试后再合并到 `main`。

---

*最后更新：2026-06-03*
