import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Users, BookOpen, MessageSquare, CheckCircle, TrendingUp, Eye } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAdminDashboard } from '@/hooks/useAdmin'
import type { DashboardStat } from '@/services/adminService'

const statIcons: Record<DashboardStat['key'], typeof Users> = {
  users: Users,
  poems: BookOpen,
  comments: MessageSquare,
  pendingCorrections: CheckCircle,
}

const defaultStats: DashboardStat[] = [
  { key: 'users', title: '用户总数', value: 0, trendPercent: 0 },
  { key: 'poems', title: '诗词总数', value: 0, trendPercent: 0 },
  { key: 'comments', title: '评论总数', value: 0, trendPercent: 0 },
  { key: 'pendingCorrections', title: '待处理纠错', value: 0, trendPercent: 0 },
]

function formatNumber(value: number): string {
  return new Intl.NumberFormat('zh-CN').format(value)
}

function formatTrend(value: number): string {
  const prefix = value > 0 ? '+' : ''
  return `${prefix}${value}%`
}

function formatRelativeTime(value: string): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'

  const diffSeconds = Math.max(0, Math.floor((Date.now() - date.getTime()) / 1000))
  if (diffSeconds < 60) return '刚刚'
  const diffMinutes = Math.floor(diffSeconds / 60)
  if (diffMinutes < 60) return `${diffMinutes}分钟前`
  const diffHours = Math.floor(diffMinutes / 60)
  if (diffHours < 24) return `${diffHours}小时前`
  const diffDays = Math.floor(diffHours / 24)
  if (diffDays < 30) return `${diffDays}天前`
  return date.toLocaleDateString('zh-CN')
}

export function AdminDashboardPage() {
  const { data: dashboard, isLoading } = useAdminDashboard()
  const stats = dashboard?.stats ?? defaultStats
  const recentActivities = dashboard?.recentActivities ?? []

  return (
    <div className="p-8 space-y-8">
      <div className="flex justify-end">
        <Button variant="outline" size="sm">
          <Eye className="h-4 w-4 mr-2" />
          查看站点
        </Button>
      </div>

      {/* Stats Grid */}
      <div className="grid md:grid-cols-4 gap-6">
        {stats.map((stat) => (
          <Card key={stat.title} className="ink-border">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {stat.title}
              </CardTitle>
              {(() => {
                const Icon = statIcons[stat.key]
                return <Icon className="h-4 w-4 text-muted-foreground" />
              })()}
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{isLoading ? '-' : formatNumber(stat.value)}</div>
              <div className="flex items-center gap-1 mt-1">
                <TrendingUp className={`h-3 w-3 ${stat.trendPercent >= 0 ? 'text-green-500' : 'text-red-500'}`} />
                <span className={`text-xs ${stat.trendPercent >= 0 ? 'text-green-500' : 'text-red-500'}`}>
                  {formatTrend(stat.trendPercent)} 较上周
                </span>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Recent Activity */}
      <Card className="ink-border">
        <CardHeader>
          <CardTitle className="font-serif text-lg">最近活动</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {!isLoading && recentActivities.length === 0 && (
              <div className="py-8 text-center text-sm text-muted-foreground">暂无最近活动</div>
            )}
            {recentActivities.map((activity) => (
              <div key={activity.id} className="flex items-start justify-between py-3 border-b last:border-0">
                <div>
                  <p className="font-medium">{activity.action}</p>
                  <p className="text-sm text-muted-foreground">{activity.detail}</p>
                </div>
                <span className="text-xs text-muted-foreground">{formatRelativeTime(activity.createdAt)}</span>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
