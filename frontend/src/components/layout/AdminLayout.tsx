import { Outlet, Link, useLocation } from 'react-router-dom'
import { LayoutDashboard, Users, BookOpen, MessageSquare, CheckCircle, Shield, Key } from 'lucide-react'

const adminNavItems = [
  { path: '/admin/dashboard', label: '仪表盘', icon: LayoutDashboard },
  { path: '/admin/users', label: '用户管理', icon: Users },
  { path: '/admin/roles', label: '角色管理', icon: Shield },
  { path: '/admin/permissions', label: '权限管理', icon: Key },
  { path: '/admin/poems', label: '诗词管理', icon: BookOpen },
  { path: '/admin/comments', label: '评论管理', icon: MessageSquare },
  { path: '/admin/corrections', label: '纠错管理', icon: CheckCircle },
]

export function AdminLayout() {
  const location = useLocation()

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
      <main className="flex-1 bg-background">
        <Outlet />
      </main>
    </div>
  )
}
