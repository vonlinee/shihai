import { useState } from 'react'
import { BookOpen, CheckCircle2, Heart, UserRound } from 'lucide-react'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useMyCorrections } from '@/hooks/usePoems'
import { useProfile } from '@/hooks/useAuth'
import { useAuthStore } from '@/stores/authStore'
import { cn } from '@/utils/cn'
import type { CorrectionRequest } from '@/types'

type UserCenterSection = 'profile' | 'favorites' | 'corrections'

const menuItems: { key: UserCenterSection; label: string; description: string; icon: typeof UserRound }[] = [
  { key: 'profile', label: '基本信息', description: '账号与角色', icon: UserRound },
  { key: 'favorites', label: '我的收藏', description: '收藏的诗词', icon: Heart },
  { key: 'corrections', label: '我的纠错', description: '提交记录', icon: CheckCircle2 },
]

const correctionStatusLabels: Record<CorrectionRequest['status'], string> = {
  pending: '待处理',
  voting: '处理中',
  processing: '处理中',
  approved: '处理中',
  rejected: '已驳回',
  resolved: '已解决',
  completed: '已解决',
}

const correctionStatusClassNames: Record<CorrectionRequest['status'], string> = {
  pending: 'bg-yellow-500/10 text-yellow-700',
  voting: 'bg-blue-500/10 text-blue-700',
  processing: 'bg-blue-500/10 text-blue-700',
  approved: 'bg-blue-500/10 text-blue-700',
  rejected: 'bg-red-500/10 text-red-700',
  resolved: 'bg-green-500/10 text-green-700',
  completed: 'bg-green-500/10 text-green-700',
}

function formatDate(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString()
}

export function UserCenterPage() {
  const [activeSection, setActiveSection] = useState<UserCenterSection>('profile')
  const { user } = useAuthStore()
  const { data: profile } = useProfile()
  const { data: correctionsData, isLoading: correctionsLoading } = useMyCorrections()

  const displayUser = profile ?? user
  const corrections = correctionsData?.list ?? []
  const activeMenuItem = menuItems.find((item) => item.key === activeSection) ?? menuItems[0]

  return (
    <div className="container py-8">
      <div className="mb-6 flex items-center gap-3">
        <BookOpen className="h-7 w-7 text-cinnabar" />
        <div>
          <h1 className="font-serif text-2xl font-bold text-ink">个人中心</h1>
          <p className="text-sm text-muted-foreground">管理账号资料、收藏与纠错申请</p>
        </div>
      </div>

      <div className="flex items-start gap-6">
        <aside className="w-[260px] shrink-0 rounded-lg border bg-card p-2 shadow-sm">
          <nav className="flex flex-col gap-2">
            {menuItems.map((item) => {
              const Icon = item.icon
              const isActive = activeSection === item.key

              return (
                <button
                  key={item.key}
                  type="button"
                  onClick={() => setActiveSection(item.key)}
                  className={cn(
                    'flex w-full items-center gap-3 rounded-md px-3 py-3 text-left transition-colors',
                    isActive
                      ? 'bg-cinnabar/10 text-cinnabar'
                      : 'text-muted-foreground hover:bg-muted hover:text-foreground',
                  )}
                >
                  <Icon className="h-5 w-5 shrink-0" />
                  <span className="min-w-0">
                    <span className="block text-sm font-medium">{item.label}</span>
                    <span className="block truncate text-xs opacity-80">{item.description}</span>
                  </span>
                </button>
              )
            })}
          </nav>
        </aside>

        <main className="min-w-0 flex-1">
          <Card className="ink-border">
            <CardHeader>
              <CardTitle className="font-serif text-xl">{activeMenuItem.label}</CardTitle>
            </CardHeader>
            <CardContent>
              {activeSection === 'profile' && (
                <div className="grid gap-4 sm:grid-cols-2">
                  <InfoField label="用户名" value={displayUser?.username} />
                  <InfoField label="姓名" value={displayUser?.name} />
                  <InfoField label="角色" value={displayUser?.role} />
                </div>
              )}

              {activeSection === 'favorites' && (
                <p className="text-muted-foreground">暂无收藏</p>
              )}

              {activeSection === 'corrections' && (
                <CorrectionList corrections={corrections} isLoading={correctionsLoading} />
              )}
            </CardContent>
          </Card>
        </main>
      </div>
    </div>
  )
}

function InfoField({ label, value }: { label: string; value?: string }) {
  return (
    <div className="rounded-md border border-border bg-muted/20 p-4">
      <div className="text-sm text-muted-foreground">{label}</div>
      <div className="mt-2 text-base font-medium text-foreground">{value || '-'}</div>
    </div>
  )
}

function CorrectionList({
  corrections,
  isLoading,
}: {
  corrections: CorrectionRequest[]
  isLoading: boolean
}) {
  if (isLoading) {
    return <p className="text-muted-foreground">加载中...</p>
  }

  if (corrections.length === 0) {
    return <p className="text-muted-foreground">暂无纠错申请</p>
  }

  return (
    <div className="space-y-3">
      {corrections.map((correction) => (
        <div key={correction.id} className="rounded-md border border-border p-4">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
            <div className="min-w-0">
              <p className="truncate font-medium text-foreground">
                {correction.poem?.title ?? '未知诗词'}
              </p>
              <p className="mt-1 line-clamp-2 text-sm text-muted-foreground">
                {correction.suggestedText}
              </p>
            </div>
            <span className={`w-fit shrink-0 rounded px-2 py-1 text-xs ${correctionStatusClassNames[correction.status]}`}>
              {correctionStatusLabels[correction.status]}
            </span>
          </div>
          <p className="mt-3 text-xs text-muted-foreground">
            提交于 {formatDate(correction.createdAt)}
          </p>
        </div>
      ))}
    </div>
  )
}
