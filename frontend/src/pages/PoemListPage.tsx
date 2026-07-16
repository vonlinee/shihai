import { useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Search, Filter, ChevronDown, ChevronRight } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Pagination } from '@/components/ui/Pagination'
import { usePoems, useDynasties, useGenres } from '@/hooks/usePoems'

export function PoemListPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [searchQuery, setSearchQuery] = useState(searchParams.get('keyword') || '')
  const [selectedDynasty, setSelectedDynasty] = useState('')
  const [selectedGenre, setSelectedGenre] = useState('')
  const [isDynastyExpanded, setIsDynastyExpanded] = useState(true)
  const [isGenreExpanded, setIsGenreExpanded] = useState(true)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)

  const { data: poemData, isLoading } = usePoems({
    keyword: searchQuery || undefined,
    dynasty: selectedDynasty || undefined,
    genre: selectedGenre || undefined,
    page,
    pageSize,
  })

  const { data: dynasties } = useDynasties()
  const { data: genres } = useGenres()

  const poems = poemData?.list ?? []
  const total = poemData?.total ?? 0
  const dynastyOptions = dynasties ?? []
  const genreOptions = genres ?? []

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
                <button
                  type="button"
                  className="mb-2 flex w-full items-center justify-between rounded px-2 py-1 text-left text-sm font-medium hover:bg-muted"
                  onClick={() => setIsDynastyExpanded((current) => !current)}
                >
                  <span>朝代</span>
                  <span className="flex items-center gap-1 text-xs text-muted-foreground">
                    {isDynastyExpanded ? '收起' : '展开'}
                    {isDynastyExpanded ? <ChevronDown className="h-3.5 w-3.5" /> : <ChevronRight className="h-3.5 w-3.5" />}
                  </span>
                </button>
                {isDynastyExpanded && (
                  <div className="space-y-1">
                  <button
                    className={`block w-full text-left px-2 py-1 text-sm rounded ${
                      selectedDynasty === '' ? 'text-foreground bg-muted' : 'text-muted-foreground hover:text-foreground hover:bg-muted'
                    }`}
                    onClick={() => { setSelectedDynasty(''); setPage(1) }}
                  >
                    全部
                  </button>
                  {dynastyOptions.map((dynasty) => (
                    <button
                      key={dynasty.id}
                      className={`block w-full text-left px-2 py-1 text-sm rounded ${
                        selectedDynasty === dynasty.name ? 'text-foreground bg-muted' : 'text-muted-foreground hover:text-foreground hover:bg-muted'
                      }`}
                      onClick={() => { setSelectedDynasty(dynasty.name); setPage(1) }}
                    >
                      {dynasty.name}
                    </button>
                  ))}
                  </div>
                )}
              </div>
              <div>
                <button
                  type="button"
                  className="mb-2 flex w-full items-center justify-between rounded px-2 py-1 text-left text-sm font-medium hover:bg-muted"
                  onClick={() => setIsGenreExpanded((current) => !current)}
                >
                  <span>体裁</span>
                  <span className="flex items-center gap-1 text-xs text-muted-foreground">
                    {isGenreExpanded ? '收起' : '展开'}
                    {isGenreExpanded ? <ChevronDown className="h-3.5 w-3.5" /> : <ChevronRight className="h-3.5 w-3.5" />}
                  </span>
                </button>
                {isGenreExpanded && (
                  <div className="space-y-1">
                  <button
                    className={`block w-full text-left px-2 py-1 text-sm rounded ${
                      selectedGenre === '' ? 'text-foreground bg-muted' : 'text-muted-foreground hover:text-foreground hover:bg-muted'
                    }`}
                    onClick={() => { setSelectedGenre(''); setPage(1) }}
                  >
                    全部
                  </button>
                  {genreOptions.map((genre) => (
                    <button
                      key={genre}
                      className={`block w-full text-left px-2 py-1 text-sm rounded ${
                        selectedGenre === genre
                          ? 'text-foreground bg-muted'
                          : 'text-muted-foreground hover:text-foreground hover:bg-muted'
                      }`}
                      onClick={() => {
                        setSelectedGenre(genre)
                        setPage(1)
                      }}
                    >
                      {genre}
                    </button>
                  ))}
                  </div>
                )}
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
            <div className="space-y-2">
              {poems.map((poem) => (
                <Card 
                  key={poem.id} 
                  className="cursor-pointer hover:shadow-lg transition-shadow ink-border"
                  onClick={() => navigate(`/poems/${poem.id}`)}
                >
                  <CardHeader className="p-3">
                    <div className="flex items-center justify-between gap-3">
                      <div className="min-w-0 space-y-1">
                        <div className="flex min-w-0 flex-col gap-1 sm:flex-row sm:items-baseline sm:gap-3">
                          <CardTitle className="truncate font-serif text-lg">{poem.title}</CardTitle>
                          <p className="shrink-0 text-sm text-muted-foreground">
                            [{poem.dynasty?.name}] {poem.author?.name}
                          </p>
                        </div>
                        <div className="flex items-center gap-3 text-xs text-muted-foreground">
                          <span>{poem.views} 阅读</span>
                          <span>{poem.likes} 点赞</span>
                        </div>
                      </div>
                      <ChevronRight className="h-5 w-5 shrink-0 text-muted-foreground" />
                    </div>
                  </CardHeader>
                </Card>
              ))}
            </div>
          )}

          <Pagination
            className="mt-8"
            page={page}
            pageSize={pageSize}
            total={total}
            onPageChange={setPage}
            onPageSizeChange={(nextPageSize) => {
              setPageSize(nextPageSize)
              setPage(1)
            }}
          />
        </div>
      </div>
    </div>
  )
}
