import { useNavigate } from 'react-router-dom'
import { Search, BookOpen, MessageSquare, HelpCircle, ChevronRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useRandomPoems } from '@/hooks/usePoems'
import { useAnnouncements } from '@/hooks/useAnnouncements'

export function HomePage() {
  const navigate = useNavigate()
  const { data: featuredPoems, isLoading: poemsLoading } = useRandomPoems(3)
  const { data: announcementsData, isLoading: announcementsLoading } = useAnnouncements(1, 5)

  return (
    <div className="space-y-12">
      {/* Hero Section */}
      <section className="relative py-20 lg:py-32 overflow-hidden">
        <div className="container relative z-10">
          <div className="max-w-3xl mx-auto text-center space-y-6">
            <h1 className="text-4xl md:text-6xl font-serif font-bold text-ink">
              品味<span className="text-cinnabar">诗词</span>之美
            </h1>
            <p className="text-lg text-muted-foreground max-w-xl mx-auto">
              识海古诗词学习平台，汇聚千年文化精华，让您在诗词的海洋中自由遨游
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center max-w-md mx-auto">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input 
                  placeholder="搜索诗词、诗人..." 
                  className="pl-10"
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      navigate(`/poems?keyword=${e.currentTarget.value}`)
                    }
                  }}
                />
              </div>
              <Button onClick={() => navigate('/poems')}>
                开始探索
              </Button>
            </div>
          </div>
        </div>
        {/* Decorative background */}
        <div className="absolute inset-0 -z-10 opacity-5">
          <div className="absolute top-10 left-10 text-9xl font-serif text-ink">诗</div>
          <div className="absolute bottom-10 right-10 text-9xl font-serif text-ink">词</div>
        </div>
      </section>

      {/* Featured Poems */}
      <section className="container">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-serif font-bold text-ink">精选诗词</h2>
          <Button variant="ghost" onClick={() => navigate('/poems')}>
            查看更多 <ChevronRight className="h-4 w-4 ml-1" />
          </Button>
        </div>
        {poemsLoading ? (
          <div className="text-center py-12 text-muted-foreground">加载中...</div>
        ) : (
          <div className="grid md:grid-cols-3 gap-6">
            {(featuredPoems ?? []).map((poem) => (
              <Card 
                key={poem.id} 
                className="cursor-pointer hover:shadow-lg transition-shadow ink-border"
                onClick={() => navigate(`/poems/${poem.id}`)}
              >
                <CardHeader className="pb-3">
                  <CardTitle className="font-serif text-xl">{poem.title}</CardTitle>
                  <p className="text-sm text-muted-foreground">
                    [{poem.dynasty?.name}] {poem.author?.name}
                  </p>
                </CardHeader>
                <CardContent>
                  <p className="poetry-text text-foreground/80 whitespace-pre-line">
                    {poem.content}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </section>

      {/* Features */}
      <section className="container py-12">
        <div className="grid md:grid-cols-3 gap-8">
          <div className="text-center space-y-4">
            <div className="w-16 h-16 mx-auto rounded-full bg-primary/10 flex items-center justify-center">
              <BookOpen className="h-8 w-8 text-primary" />
            </div>
            <h3 className="text-lg font-serif font-semibold">海量诗词</h3>
            <p className="text-sm text-muted-foreground">
              收录唐诗宋词等经典作品，涵盖各个朝代的名家名作
            </p>
          </div>
          <div className="text-center space-y-4">
            <div className="w-16 h-16 mx-auto rounded-full bg-primary/10 flex items-center justify-center">
              <MessageSquare className="h-8 w-8 text-primary" />
            </div>
            <h3 className="text-lg font-serif font-semibold">互动交流</h3>
            <p className="text-sm text-muted-foreground">
              评论、论坛功能让您与其他诗词爱好者交流心得
            </p>
          </div>
          <div className="text-center space-y-4">
            <div className="w-16 h-16 mx-auto rounded-full bg-primary/10 flex items-center justify-center">
              <HelpCircle className="h-8 w-8 text-primary" />
            </div>
            <h3 className="text-lg font-serif font-semibold">趣味问答</h3>
            <p className="text-sm text-muted-foreground">
              通过诗词问答检验学习成果，寓教于乐
            </p>
          </div>
        </div>
      </section>

      {/* Announcements */}
      <section className="container pb-12">
        <Card className="ink-border">
          <CardHeader>
            <CardTitle className="font-serif text-xl">公告通知</CardTitle>
          </CardHeader>
          <CardContent>
            {announcementsLoading ? (
              <div className="text-center py-4 text-muted-foreground">加载中...</div>
            ) : (
              <div className="space-y-3">
                {(announcementsData?.list ?? []).map((item) => (
                  <div 
                    key={item.id} 
                    className="flex items-center justify-between py-2 border-b last:border-0"
                  >
                    <span className="text-sm">{item.title}</span>
                    <span className="text-xs text-muted-foreground">
                      {new Date(item.createdAt).toLocaleDateString()}
                    </span>
                  </div>
                ))}
                {(announcementsData?.list ?? []).length === 0 && (
                  <p className="text-center text-muted-foreground py-4">暂无公告</p>
                )}
              </div>
            )}
          </CardContent>
        </Card>
      </section>
    </div>
  )
}
