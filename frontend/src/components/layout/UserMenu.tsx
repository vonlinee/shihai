import { useState, useRef, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { User, LogOut, Settings, Shield, BookOpen, ChevronRight } from 'lucide-react'
import { UserAvatar } from '@/components/ui/UserAvatar'
import { useAuthStore } from '@/stores/authStore'
import { ThemeSwitcher } from './ThemeSwitcher'
import { toast } from 'sonner'

export function UserMenu() {
  const [isOpen, setIsOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)
  const { user, isAuthenticated, logout } = useAuthStore()
  const navigate = useNavigate()

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const handleLogout = () => {
    setIsOpen(false)
    logout()
    toast.success('已退出登录')
    navigate('/')
  }

  if (!isAuthenticated || !user) {
    return (
      <div className="flex items-center gap-3">
        <button
          onClick={() => navigate('/login')}
          className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
        >
          登录
        </button>
        <button
          onClick={() => navigate('/register')}
          className="text-sm font-medium bg-primary text-primary-foreground px-4 py-1.5 rounded-md hover:bg-primary/90 transition-colors"
        >
          注册
        </button>
      </div>
    )
  }

  const getRoleLabel = (role: string) => {
    switch (role) {
      case 'admin': return '管理员'
      case 'editor': return '编辑'
      case 'reviewer': return '审核员'
      default: return '用户'
    }
  }

  return (
    <div className="relative" ref={menuRef}>
      {/* Avatar Button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="relative flex items-center gap-2 rounded-full p-0.5 hover:ring-2 hover:ring-primary/30 transition-all focus:outline-none focus:ring-2 focus:ring-primary/50"
        aria-label="用户菜单"
      >
        <UserAvatar
          avatar={user.avatar}
          name={user.name}
          username={user.username}
          className="h-8 w-8 cursor-pointer"
        />
        {/* Online indicator */}
        <span className="absolute bottom-0 right-0 h-2.5 w-2.5 rounded-full bg-green-500 border-2 border-background" />
      </button>

      {/* Dropdown Panel */}
      {isOpen && (
        <div className="absolute right-0 mt-2 w-64 rounded-lg border bg-background shadow-lg z-50 animate-in fade-in-0 zoom-in-95">
          {/* User Info Card */}
          <div className="px-4 py-3 border-b bg-muted/30">
            <div className="flex items-center gap-3">
              <UserAvatar
                avatar={user.avatar}
                name={user.name}
                username={user.username}
                className="h-10 w-10"
              />
              <div className="flex-1 min-w-0">
                <div className="font-medium text-sm truncate">{user.name || user.username}</div>
                <div className="text-xs text-muted-foreground truncate">@{user.username}</div>
              </div>
            </div>
            <div className="mt-2">
              <span className="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium bg-primary/10 text-primary">
                {getRoleLabel(user.role)}
              </span>
            </div>
          </div>

          {/* Menu Items */}
          <div className="py-1">
            <Link
              to="/user"
              onClick={() => setIsOpen(false)}
              className="flex items-center gap-3 px-4 py-2 text-sm hover:bg-muted/50 transition-colors"
            >
              <User className="h-4 w-4 text-muted-foreground" />
              <span className="flex-1">个人中心</span>
              <ChevronRight className="h-3 w-3 text-muted-foreground" />
            </Link>

            <Link
              to="/user"
              onClick={() => setIsOpen(false)}
              className="flex items-center gap-3 px-4 py-2 text-sm hover:bg-muted/50 transition-colors"
            >
              <BookOpen className="h-4 w-4 text-muted-foreground" />
              <span className="flex-1">我的收藏</span>
              <ChevronRight className="h-3 w-3 text-muted-foreground" />
            </Link>

            {user.role === 'admin' && (
              <a
                href="/admin"
                target="_blank"
                rel="noopener noreferrer"
                onClick={() => setIsOpen(false)}
                className="flex items-center gap-3 px-4 py-2 text-sm hover:bg-muted/50 transition-colors"
              >
                <Shield className="h-4 w-4 text-muted-foreground" />
                <span className="flex-1">管理后台</span>
                <ChevronRight className="h-3 w-3 text-muted-foreground" />
              </a>
            )}

            <Link
              to="/user"
              onClick={() => setIsOpen(false)}
              className="flex items-center gap-3 px-4 py-2 text-sm hover:bg-muted/50 transition-colors"
            >
              <Settings className="h-4 w-4 text-muted-foreground" />
              <span className="flex-1">账号设置</span>
              <ChevronRight className="h-3 w-3 text-muted-foreground" />
            </Link>
          </div>

          {/* Theme Switcher */}
          <div className="border-t">
            <ThemeSwitcher />
          </div>

          {/* Divider + Logout */}
          <div className="border-t py-1">
            <button
              onClick={handleLogout}
              className="flex items-center gap-3 px-4 py-2 text-sm w-full text-left hover:bg-muted/50 transition-colors text-red-600"
            >
              <LogOut className="h-4 w-4" />
              <span>退出登录</span>
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
