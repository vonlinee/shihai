import { createBrowserRouter, Navigate } from 'react-router-dom'
import { MainLayout } from './components/layout/MainLayout'
import { HomePage } from './pages/HomePage'
import { PoemListPage } from './pages/PoemListPage'
import { PoemDetailPage } from './pages/PoemDetailPage'
import { ForumPage } from './pages/ForumPage'
import { QuizPage } from './pages/QuizPage'
import { LoginPage } from './pages/LoginPage'
import { RegisterPage } from './pages/RegisterPage'
import { UserCenterPage } from './pages/UserCenterPage'
import { AdminLayout } from './components/layout/AdminLayout'
import { AdminDashboardPage } from './pages/admin/DashboardPage'
import { AdminUsersPage } from './pages/admin/UsersPage'
import { AdminPoemsPage } from './pages/admin/PoemsPage'
import { AdminCommentsPage } from './pages/admin/CommentsPage'
import { AdminCorrectionsPage } from './pages/admin/CorrectionsPage'
import { RolesPage } from './pages/admin/RolesPage'
import { PermissionsPage } from './pages/admin/PermissionsPage'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />,
    children: [
      { index: true, element: <HomePage /> },
      { path: 'poems', element: <PoemListPage /> },
      { path: 'poems/:id', element: <PoemDetailPage /> },
      { path: 'forum', element: <ForumPage /> },
      { path: 'quiz', element: <QuizPage /> },
      { path: 'login', element: <LoginPage /> },
      { path: 'register', element: <RegisterPage /> },
      { path: 'user', element: <UserCenterPage /> },
    ],
  },
  {
    path: '/admin',
    element: <AdminLayout />,
    children: [
      { index: true, element: <Navigate to="/admin/dashboard" replace /> },
      { path: 'dashboard', element: <AdminDashboardPage /> },
      { path: 'users', element: <AdminUsersPage /> },
      { path: 'roles', element: <RolesPage /> },
      { path: 'permissions', element: <PermissionsPage /> },
      { path: 'poems', element: <AdminPoemsPage /> },
      { path: 'comments', element: <AdminCommentsPage /> },
      { path: 'corrections', element: <AdminCorrectionsPage /> },
    ],
  },
])
