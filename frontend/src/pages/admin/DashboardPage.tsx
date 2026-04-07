import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Users, BookOpen, MessageSquare, CheckCircle, TrendingUp, Eye } from 'lucide-react'
import { Button } from '@/components/ui/button'

const stats = [
  { title: '用户总数', value: '1,234', icon: Users, trend: '+12%' },
  { title: '诗词总数', value: '5,678', icon: BookOpen, trend: '+5%' },
  { title: '评论总数', value: '9,012', icon: MessageSquare, trend: '+18%' },
  { title: '待处理纠错', value: '23', icon: CheckCircle, trend: '-3%' },
]

const recentActivities = [
  { id: 1, action: '用户注册', detail: '新用户 张三 注册了账号', time: '2分钟前' },
  { id: 2, action: '诗词纠错', detail: '用户提交了《静夜思》的纠错申请', time: '15分钟前' },
  { id: 3, action: '评论发布', detail: '用户 李四 评论了《春晓》', time: '30分钟前' },
  { id: 4, action: '诗词添加', detail: '管理员添加了新诗《登鹳雀楼》', time: '1小时前' },
]

export function AdminDashboardPage() {
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
              <stat.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <div className="flex items-center gap-1 mt-1">
                <TrendingUp className={`h-3 w-3 ${stat.trend.startsWith('+') ? 'text-green-500' : 'text-red-500'}`} />
                <span className={`text-xs ${stat.trend.startsWith('+') ? 'text-green-500' : 'text-red-500'}`}>
                  {stat.trend} 较上周
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
            {recentActivities.map((activity) => (
              <div key={activity.id} className="flex items-start justify-between py-3 border-b last:border-0">
                <div>
                  <p className="font-medium">{activity.action}</p>
                  <p className="text-sm text-muted-foreground">{activity.detail}</p>
                </div>
                <span className="text-xs text-muted-foreground">{activity.time}</span>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
