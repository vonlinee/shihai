import { lazy, Suspense, type ComponentType, type ReactNode } from 'react'
import { createBrowserRouter, Navigate } from 'react-router-dom'
import { MainLayout } from './components/layout/MainLayout'
import { AdminLayout } from './components/layout/AdminLayout'

function lazyNamed<TModule, TExport extends keyof TModule>(
  loader: () => Promise<TModule>,
  exportName: TExport,
) {
  return lazy(async () => {
    const module = await loader()
    return { default: module[exportName] as ComponentType }
  })
}

function PageLoading() {
  return (
    <div className="flex min-h-[240px] items-center justify-center text-sm text-muted-foreground">
      加载中...
    </div>
  )
}

function withSuspense(element: ReactNode) {
  return <Suspense fallback={<PageLoading />}>{element}</Suspense>
}

const HomePage = lazyNamed(() => import('./pages/HomePage'), 'HomePage')
const PoemListPage = lazyNamed(() => import('./pages/PoemListPage'), 'PoemListPage')
const PoemDetailPage = lazyNamed(() => import('./pages/PoemDetailPage'), 'PoemDetailPage')
const ForumPage = lazyNamed(() => import('./pages/ForumPage'), 'ForumPage')
const QuizPage = lazyNamed(() => import('./pages/QuizPage'), 'QuizPage')
const LoginPage = lazyNamed(() => import('./pages/LoginPage'), 'LoginPage')
const RegisterPage = lazyNamed(() => import('./pages/RegisterPage'), 'RegisterPage')
const UserCenterPage = lazyNamed(() => import('./pages/UserCenterPage'), 'UserCenterPage')
const AdminDashboardPage = lazyNamed(() => import('./pages/admin/DashboardPage'), 'AdminDashboardPage')
const AdminUsersPage = lazyNamed(() => import('./pages/admin/UsersPage'), 'AdminUsersPage')
const AdminPoemsPage = lazyNamed(() => import('./pages/admin/PoemsPage'), 'AdminPoemsPage')
const AdminWorkCollectionsPage = lazyNamed(() => import('./pages/admin/WorkCollectionsPage'), 'AdminWorkCollectionsPage')
const AdminCommentsPage = lazyNamed(() => import('./pages/admin/CommentsPage'), 'AdminCommentsPage')
const AdminForumPage = lazyNamed(() => import('./pages/admin/ForumPage'), 'AdminForumPage')
const AdminCorrectionsPage = lazyNamed(() => import('./pages/admin/CorrectionsPage'), 'AdminCorrectionsPage')

export const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />,
    children: [
      { index: true, element: withSuspense(<HomePage />) },
      { path: 'poems', element: withSuspense(<PoemListPage />) },
      { path: 'poems/:id', element: withSuspense(<PoemDetailPage />) },
      { path: 'forum', element: withSuspense(<ForumPage />) },
      { path: 'forum/:id', element: withSuspense(<ForumPage />) },
      { path: 'quiz', element: withSuspense(<QuizPage />) },
      { path: 'login', element: withSuspense(<LoginPage />) },
      { path: 'register', element: withSuspense(<RegisterPage />) },
      { path: 'user', element: withSuspense(<UserCenterPage />) },
    ],
  },
  {
    path: '/admin',
    element: <AdminLayout />,
    children: [
      { index: true, element: <Navigate to="/admin/dashboard" replace /> },
      { path: 'dashboard', element: withSuspense(<AdminDashboardPage />) },
      { path: 'users', element: withSuspense(<AdminUsersPage />) },
      { path: 'poems', element: withSuspense(<AdminPoemsPage />) },
      { path: 'work-collections', element: withSuspense(<AdminWorkCollectionsPage />) },
      { path: 'comments', element: withSuspense(<AdminCommentsPage />) },
      { path: 'forum', element: withSuspense(<AdminForumPage />) },
      { path: 'corrections', element: withSuspense(<AdminCorrectionsPage />) },
    ],
  },
])
