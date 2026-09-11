import { useCallback, useEffect, useState, type PointerEvent as ReactPointerEvent } from 'react'
import { Outlet, Link, useLocation, Navigate } from 'react-router-dom'
import {
  LayoutDashboard,
  Users,
  BookOpen,
  MessageSquare,
  CheckCircle,
  Layers,
  Pin,
  PanelLeftClose,
  PanelLeftOpen,
} from 'lucide-react'

import { UserMenu } from './UserMenu'
import { Logo } from '@/components/brand/Logo'
import { useAuthStore } from '@/stores/authStore'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { cn } from '@/utils/cn'

const sidebarMinWidth = 220
const sidebarMaxWidth = 360
const sidebarDefaultWidth = 256
const collapsedSidebarWidth = 72

const adminNavItems = [
  { path: '/admin/dashboard', label: '仪表盘', icon: LayoutDashboard },
  { path: '/admin/users', label: '用户与权限', icon: Users },
  { path: '/admin/poems', label: '诗词管理', icon: BookOpen },
  { path: '/admin/work-collections', label: '作品集管理', icon: Layers },
  { path: '/admin/comments', label: '评论管理', icon: MessageSquare },
  { path: '/admin/forum', label: '论坛管理', icon: Pin },
  { path: '/admin/corrections', label: '纠错管理', icon: CheckCircle },
]

function clampSidebarWidth(width: number): number {
  return Math.min(sidebarMaxWidth, Math.max(sidebarMinWidth, width))
}

export function AdminLayout() {
  const location = useLocation()
  const { isAuthenticated } = useAuthStore()
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false)
  const [sidebarWidth, setSidebarWidth] = useState(sidebarDefaultWidth)
  const [isResizingSidebar, setIsResizingSidebar] = useState(false)

  useEffect(() => {
    if (!isResizingSidebar) return

    const previousCursor = document.body.style.cursor
    const previousUserSelect = document.body.style.userSelect
    document.body.style.cursor = 'col-resize'
    document.body.style.userSelect = 'none'

    const handlePointerMove = (event: PointerEvent) => {
      setSidebarWidth(clampSidebarWidth(event.clientX))
    }
    const handlePointerUp = () => {
      setIsResizingSidebar(false)
    }

    window.addEventListener('pointermove', handlePointerMove)
    window.addEventListener('pointerup', handlePointerUp)
    return () => {
      document.body.style.cursor = previousCursor
      document.body.style.userSelect = previousUserSelect
      window.removeEventListener('pointermove', handlePointerMove)
      window.removeEventListener('pointerup', handlePointerUp)
    }
  }, [isResizingSidebar])

  const handleResizeStart = useCallback((event: ReactPointerEvent<HTMLButtonElement>) => {
    event.preventDefault()
    setIsResizingSidebar(true)
  }, [])

  if (!isAuthenticated) {
    return <Navigate to="/" replace />
  }

  const currentTitle = adminNavItems.find((item) => location.pathname === item.path)?.label ?? '管理后台'
  const activeSidebarWidth = isSidebarCollapsed ? collapsedSidebarWidth : sidebarWidth

  return (
    <div className="flex h-screen overflow-hidden bg-background">
      <aside
        className={cn(
          'relative flex h-screen shrink-0 flex-col border-r bg-muted/50 transition-[width] duration-200',
          isResizingSidebar && 'transition-none',
        )}
        style={{ width: activeSidebarWidth }}
      >
        <div
          className={cn(
            'flex shrink-0 items-center gap-3',
            isSidebarCollapsed ? 'flex-col p-4' : 'justify-between p-6',
          )}
        >
          {isSidebarCollapsed ? (
            <div className="group relative h-10 w-10 shrink-0">
              <Link
                to="/"
                title="识海管理"
                aria-label="识海管理首页"
                className="flex h-10 w-10 items-center justify-center rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              >
                <Logo showWordmark={false} compact />
              </Link>
              <button
                type="button"
                aria-label="展开导航栏"
                title="展开导航栏"
                onClick={() => setIsSidebarCollapsed(false)}
                className="pointer-events-none absolute inset-0 flex items-center justify-center rounded-md bg-background/95 text-muted-foreground opacity-0 shadow-sm transition-opacity hover:bg-muted hover:text-foreground group-hover:pointer-events-auto group-hover:opacity-100 group-focus-within:pointer-events-auto group-focus-within:opacity-100"
              >
                <PanelLeftOpen className="h-4 w-4" />
              </button>
            </div>
          ) : (
            <>
              <Link
                to="/"
                title="识海管理"
                aria-label="识海管理首页"
                className="flex min-w-0 items-center rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              >
                <Logo />
              </Link>
              <button
                type="button"
                aria-label="折叠导航栏"
                aria-pressed={false}
                title="折叠导航栏"
                onClick={() => setIsSidebarCollapsed(true)}
                className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
              >
                <PanelLeftClose className="h-4 w-4" />
              </button>
            </>
          )}
        </div>

        <TooltipProvider delayDuration={200}>
          <nav className={cn('min-h-0 flex-1 space-y-1 overflow-y-auto', isSidebarCollapsed ? 'px-3' : 'px-4')}>
            {adminNavItems.map((item) => {
              const isActive = location.pathname === item.path
              const navLinkClassName = cn(
                'flex items-center rounded-lg text-sm font-medium transition-colors',
                isSidebarCollapsed ? 'h-11 justify-center px-0' : 'gap-3 px-4 py-3',
                isActive
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:bg-muted hover:text-foreground',
              )

              if (isSidebarCollapsed) {
                return (
                  <Tooltip key={item.path}>
                    <TooltipTrigger asChild>
                      <Link
                        to={item.path}
                        aria-label={item.label}
                        className={navLinkClassName}
                      >
                        <item.icon className="h-4 w-4 shrink-0" />
                      </Link>
                    </TooltipTrigger>
                    <TooltipContent side="right">{item.label}</TooltipContent>
                  </Tooltip>
                )
              }

              return (
                <Link
                  key={item.path}
                  to={item.path}
                  className={navLinkClassName}
                >
                  <item.icon className="h-4 w-4 shrink-0" />
                  <span className="truncate">{item.label}</span>
                </Link>
              )
            })}
          </nav>
        </TooltipProvider>

        {!isSidebarCollapsed && (
          <button
            type="button"
            aria-label="调整导航栏宽度"
            title="拖动调整宽度"
            onPointerDown={handleResizeStart}
            className="absolute right-0 top-0 z-10 flex h-full w-2 translate-x-1/2 cursor-col-resize touch-none items-center justify-center"
          >
            <span
              className={cn(
                'h-10 w-0.5 rounded-full bg-border opacity-0 transition-opacity',
                isResizingSidebar ? 'opacity-100' : 'hover:opacity-100',
              )}
            />
          </button>
        )}
      </aside>

      <div className="min-w-0 flex h-screen flex-1 flex-col bg-background">
        <header className="z-50 flex h-14 shrink-0 items-center justify-between border-b bg-background/95 px-6 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <h1 className="text-lg font-serif font-bold text-foreground">{currentTitle}</h1>
          <UserMenu />
        </header>

        <main className="min-h-0 min-w-0 flex-1 overflow-auto">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
