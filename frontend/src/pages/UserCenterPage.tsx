import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useAuthStore } from '@/stores/authStore'
import { useProfile } from '@/hooks/useAuth'

export function UserCenterPage() {
  const { user } = useAuthStore()
  const { data: profile } = useProfile()

  const displayUser = profile ?? user

  return (
    <div className="container py-8">
      <h1 className="text-2xl font-serif font-bold text-ink mb-6">个人中心</h1>
      <div className="grid md:grid-cols-3 gap-6">
        <Card className="ink-border">
          <CardHeader>
            <CardTitle className="font-serif text-lg">基本信息</CardTitle>
          </CardHeader>
          <CardContent>
            <p>用户名: {displayUser?.username}</p>
            <p>姓名: {displayUser?.name}</p>
            <p>角色: {displayUser?.role}</p>
          </CardContent>
        </Card>
        <Card className="ink-border">
          <CardHeader>
            <CardTitle className="font-serif text-lg">我的收藏</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground">暂无收藏</p>
          </CardContent>
        </Card>
        <Card className="ink-border">
          <CardHeader>
            <CardTitle className="font-serif text-lg">我的纠错</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground">暂无纠错申请</p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
