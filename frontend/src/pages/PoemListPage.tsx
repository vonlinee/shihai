import { useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Search, Filter, ChevronRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { usePoems, useDynasties } from '@/hooks/usePoems'

export function PoemListPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [searchQuery, setSearchQuery] = useState(searchParams.get('keyword') || '')
  const [selectedDynasty, setSelectedDynasty] = useState('')
  const [selectedGenre, setSelectedGenre] = useState('')
  const [page, setPage] = useState(1)

  const { data: poemData, isLoading } = usePoems({
    keyword: searchQuery || undefined,
    dynastyId: selectedDynasty ? Number(selectedDynasty) : undefined,
    genre: selectedGenre || undefined,
    page,
    pageSize: 10,
  })

  const { data: dynasties } = useDynasties()

  const poems = poemData?.list ?? []
  const total = poemData?.total ?? 0

  return (
    <div className="container py-8">
      <div className="flex flex-col md:flex-row gap-6">
        {/* Sidebar Filters */}
        <aside className="w-full md:w-64 space-y-6">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input 
              placeholder="搜索诗词..."
              className="pl-10"
              value={searchQuery}
              onChange={(e) => {
                setSearchQuery(e.target.value)
                setPage(1)
              }}
            />
          </div>
          
          <Card className="ink-border">
            <CardHeader className="pb-3">
              <CardTitle className="font-serif text-lg flex items-center gap-2">
                <Filter className="h-4 w-4" />
                筛选
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <h4 className="text-sm font-medium mb-2">朝代</h4>
                <div className="space-y-1">
                  <button
                    className={`block w-full text-left px-2 py-1 text-sm rounded ${
                      selectedDynasty === '' ? 'text-foreground bg-muted' : 'text-muted-foreground hover:text-foreground hover:bg-muted'
                    }`}
                    onClick={() => { setSelectedDynasty(''); setPage(1) }}
                  >
                    全部
                  </button>
                  {(dynasties ?? []).map((dynasty) => (
                    <button
                      key={dynasty.id}
                      className={`block w-full text-left px-2 py-1 text-sm rounded ${
                        selectedDynasty === String(dynasty.id) ? 'text-foreground bg-muted' : 'text-muted-foreground hover:text-foreground hover:bg-muted'
                      }`}
                      onClick={() => { setSelectedDynasty(String(dynasty.id)); setPage(1) }}
                    >
                      {dynasty.name}
                    </button>
                  ))}
                </div>
              </div>
              <div>
                <h4 className="text-sm font-medium mb-2">体裁</h4>
                <div className="space-y-1">
                  {['全部', '五言绝句', '七言绝句', '五言律诗', '七言律诗', '词'].map((genre) => (
                    <button
                      key={genre}
                      className={`block w-full text-left px-2 py-1 text-sm rounded ${
                        (genre === '全部' && selectedGenre === '') || selectedGenre === genre
                          ? 'text-foreground bg-muted'
                          : 'text-muted-foreground hover:text-foreground hover:bg-muted'
                      }`}
                      onClick={() => {
                        setSelectedGenre(genre === '全部' ? '' : genre)
                        setPage(1)
                      }}
                    >
                      {genre}
                    </button>
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>
        </aside>

        {/* Poem List */}
        <div className="flex-1">
          <div className="flex items-center justify-between mb-6">
            <h1 className="text-2xl font-serif font-bold text-ink">诗词列表</h1>
            <span className="text-sm text-muted-foreground">共 {total} 首</span>
          </div>
          
          {isLoading ? (
            <div className="text-center py-12 text-muted-foreground">加载中...</div>
          ) : poems.length === 0 ? (
            <div className="text-center py-12 text-muted-foreground">暂无诗词数据</div>
          ) : (
            <div className="space-y-4">
              {poems.map((poem) => (
                <Card 
                  key={poem.id} 
                  className="cursor-pointer hover:shadow-lg transition-shadow ink-border"
                  onClick={() => navigate(`/poems/${poem.id}`)}
                >
                  <CardHeader className="pb-3">
                    <div className="flex items-start justify-between">
                      <div>
                        <CardTitle className="font-serif text-xl">{poem.title}</CardTitle>
                        <p className="text-sm text-muted-foreground mt-1">
                          [{poem.dynasty?.name}] {poem.author?.name}
                        </p>
                      </div>
                      <ChevronRight className="h-5 w-5 text-muted-foreground" />
                    </div>
                  </CardHeader>
                  <CardContent>
                    <p className="poetry-text text-foreground/80 whitespace-pre-line mb-4">
                      {poem.content}
                    </p>
                    <div className="flex items-center gap-4 text-sm text-muted-foreground">
                      <span>{poem.genre}</span>
                      <span>·</span>
                      <span>{poem.views} 阅读</span>
                      <span>·</span>
                      <span>{poem.likes} 点赞</span>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}

          {/* Pagination */}
          {total > 10 && (
            <div className="flex justify-center gap-2 mt-8">
              <Button
                variant="outline"
                size="sm"
                disabled={page <= 1}
                onClick={() => setPage(page - 1)}
              >
                上一页
              </Button>
              <span className="flex items-center px-4 text-sm text-muted-foreground">
                第 {page} 页 / 共 {Math.ceil(total / 10)} 页
              </span>
              <Button
                variant="outline"
                size="sm"
                disabled={page >= Math.ceil(total / 10)}
                onClick={() => setPage(page + 1)}
              >
                下一页
              </Button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
