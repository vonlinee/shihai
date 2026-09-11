import { Link } from 'react-router-dom'
import { BookOpen, MessageSquare, HelpCircle } from 'lucide-react'
import { Logo } from '@/components/brand/Logo'
import { UserMenu } from './UserMenu'

export function Header() {
  const navItems = [
    { path: '/poems', label: '诗词', icon: BookOpen },
    { path: '/forum', label: '论坛', icon: MessageSquare },
    { path: '/quiz', label: '问答', icon: HelpCircle },
  ]

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-16 items-center px-6">
        <Link
          to="/"
          aria-label="识海首页"
          className="flex shrink-0 items-center rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <Logo />
        </Link>

        {/* Navigation */}
        <nav className="hidden md:flex items-center gap-6 ml-8">
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

        {/* Spacer to push avatar to far right */}
        <div className="flex-1" />

        {/* User Menu (GitHub-style avatar dropdown) */}
        <UserMenu />
      </div>
    </header>
  )
}
