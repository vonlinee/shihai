import { Outlet, Link, useLocation, Navigate } from 'react-router-dom'
import { LayoutDashboard, Users, BookOpen, MessageSquare, CheckCircle } from 'lucide-react'
import { UserMenu } from './UserMenu'
import { useAuthStore } from '@/stores/authStore'

const adminNavItems = [
  { path: '/admin/dashboard', label: '仪表盘', icon: LayoutDashboard },
  { path: '/admin/users', label: '用户与权限', icon: Users },
  { path: '/admin/poems', label: '诗词管理', icon: BookOpen },
  { path: '/admin/comments', label: '评论管理', icon: MessageSquare },
  { path: '/admin/corrections', label: '纠错管理', icon: CheckCircle },
]

export function AdminLayout() {
  const location = useLocation()
  const { isAuthenticated } = useAuthStore()

  if (!isAuthenticated) {
    return <Navigate to="/" replace />
  }

  const currentTitle = adminNavItems.find((item) => location.pathname === item.path)?.label ?? '管理后台'

  return (
    <div className="min-h-screen flex">
      {/* Sidebar */}
      <aside className="w-64 bg-muted/50 border-r">
        <div className="p-6">
          <Link to="/" className="flex items-center gap-2 font-serif text-xl font-bold text-ink">
            <span className="text-cinnabar">识</span>海管理
          </Link>
        </div>
        <nav className="px-4 space-y-1">
          {adminNavItems.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className={`flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium transition-colors ${
                location.pathname === item.path
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:text-foreground hover:bg-muted'
              }`}
            >
              <item.icon className="h-4 w-4" />
              {item.label}
            </Link>
          ))}
        </nav>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col bg-background min-h-0">
        {/* Sticky Top Bar */}
        <header className="sticky top-0 z-50 h-14 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 flex items-center justify-between px-6 shrink-0">
          <h1 className="text-lg font-serif font-bold text-ink">{currentTitle}</h1>
          <UserMenu />
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-auto">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
