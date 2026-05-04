import { Card, CardContent } from '@/components/ui/card'

export function QuizPage() {
  return (
    <div className="container py-8">
      <h1 className="text-2xl font-serif font-bold text-ink mb-6">诗词问答</h1>
      <Card className="ink-border">
        <CardContent className="py-12 text-center">
          <p className="text-muted-foreground">问答功能开发中，敬请期待...</p>
        </CardContent>
      </Card>
    </div>
  )
}
