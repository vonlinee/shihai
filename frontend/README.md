# 识海前端

这是识海古诗词学习平台的前端项目，基于 React、TypeScript、Vite、Tailwind CSS、TanStack Query 和 Radix UI 构建。

## 快速开始

```bash
npm install
npm run dev
```

开发服务器默认运行在 `http://localhost:3000`，`/api` 请求会代理到 `http://localhost:8080`。

## 常用命令

```bash
npm run dev      # 启动开发服务器
npm run build    # 类型检查并构建生产包
npm run lint     # 运行 ESLint
npm run preview  # 本地预览生产构建
```

## 目录结构

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

## 开发文档

修改前端代码前，优先阅读：

- [前端指南](../docs/guidelines/frontend-guidelines.md)
- [项目开发指南索引](../docs/guidelines/README.md)

涉及跨端功能时，还应阅读后端相关指南。
