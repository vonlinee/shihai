import { Link, useNavigate } from 'react-router-dom'
import { BookOpen, MessageSquare, HelpCircle, User, LogOut } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/authStore'
import { toast } from 'sonner'

export function Header() {
  const { user, isAuthenticated, logout } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    toast.success('已退出登录')
    navigate('/')
  }

  const navItems = [
    { path: '/poems', label: '诗词', icon: BookOpen },
    { path: '/forum', label: '论坛', icon: MessageSquare },
    { path: '/quiz', label: '问答', icon: HelpCircle },
  ]

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-16 items-center justify-between">
        {/* Logo */}
        <Link to="/" className="flex items-center gap-2 font-serif text-2xl font-bold text-ink">
          <span className="text-cinnabar">识</span>海
        </Link>

        {/* Navigation */}
        <nav className="hidden md:flex items-center gap-6">
          {navItems.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className="flex items-center gap-2 text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
            >
              <item.icon className="h-4 w-4" />
              {item.label}
            </Link>
          ))}
        </nav>

        {/* User Actions */}
        <div className="flex items-center gap-4">
          {isAuthenticated ? (
            <>
              {user?.role === 'admin' && (
                <Button variant="ghost" size="sm" onClick={() => navigate('/admin')}>
                  管理后台
                </Button>
              )}
              <Button variant="ghost" size="sm" onClick={() => navigate('/user')}>
                <User className="h-4 w-4 mr-2" />
                {user?.name}
              </Button>
              <Button variant="ghost" size="sm" onClick={handleLogout}>
                <LogOut className="h-4 w-4" />
              </Button>
            </>
          ) : (
            <>
              <Button variant="ghost" size="sm" onClick={() => navigate('/login')}>
                登录
              </Button>
              <Button size="sm" onClick={() => navigate('/register')}>
                注册
              </Button>
            </>
          )}
        </div>
      </div>
    </header>
  )
}
