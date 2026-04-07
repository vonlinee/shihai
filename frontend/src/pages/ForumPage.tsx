import { Card, CardContent } from '@/components/ui/card'

export function ForumPage() {
  return (
    <div className="container py-8">
      <h1 className="text-2xl font-serif font-bold text-ink mb-6">诗词论坛</h1>
      <Card className="ink-border">
        <CardContent className="py-12 text-center">
          <p className="text-muted-foreground">论坛功能开发中，敬请期待...</p>
        </CardContent>
      </Card>
    </div>
  )
}
