import { useState } from 'react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Search, Plus, Edit2, Trash2, Eye } from 'lucide-react'
import { usePoems } from '@/hooks/usePoems'
import { useAdminDeletePoem } from '@/hooks/useAdmin'
import { useNavigate } from 'react-router-dom'

export function AdminPoemsPage() {
  const navigate = useNavigate()
  const [searchQuery, setSearchQuery] = useState('')
  const [page, setPage] = useState(1)

  const { data: poemData, isLoading } = usePoems({
    keyword: searchQuery || undefined,
    page,
    pageSize: 10,
  })

  const deletePoemMutation = useAdminDeletePoem()

  const poems = poemData?.list ?? []
  const total = poemData?.total ?? 0

  const handleDelete = (id: number) => {
    if (!confirm('确定要删除该诗词吗？')) return
    deletePoemMutation.mutate(id)
  }

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-serif font-bold text-ink">诗词管理</h1>
        <Button>
          <Plus className="h-4 w-4 mr-2" />
          添加诗词
        </Button>
      </div>

      <Card className="ink-border">
        <CardHeader className="pb-4">
          <div className="flex items-center gap-4">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="搜索诗词标题或作者..."
                className="pl-10"
                value={searchQuery}
                onChange={(e) => {
                  setSearchQuery(e.target.value)
                  setPage(1)
                }}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-center py-8 text-muted-foreground">加载中...</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="text-left py-3 px-4 font-medium">ID</th>
                    <th className="text-left py-3 px-4 font-medium">标题</th>
                    <th className="text-left py-3 px-4 font-medium">作者</th>
                    <th className="text-left py-3 px-4 font-medium">朝代</th>
                    <th className="text-left py-3 px-4 font-medium">体裁</th>
                    <th className="text-left py-3 px-4 font-medium">浏览量</th>
                    <th className="text-left py-3 px-4 font-medium">添加时间</th>
                    <th className="text-left py-3 px-4 font-medium">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {poems.map((poem) => (
                    <tr key={poem.id} className="border-b last:border-0 hover:bg-muted/50">
                      <td className="py-3 px-4">{poem.id}</td>
                      <td className="py-3 px-4 font-medium">{poem.title}</td>
                      <td className="py-3 px-4">{poem.author?.name}</td>
                      <td className="py-3 px-4">{poem.dynasty?.name}</td>
                      <td className="py-3 px-4">{poem.genre}</td>
                      <td className="py-3 px-4">{poem.views}</td>
                      <td className="py-3 px-4 text-muted-foreground">
                        {new Date(poem.createdAt).toLocaleDateString()}
                      </td>
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-2">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => navigate(`/poems/${poem.id}`)}
                          >
                            <Eye className="h-4 w-4" />
                          </Button>
                          <Button variant="ghost" size="sm">
                            <Edit2 className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDelete(poem.id)}
                          >
                            <Trash2 className="h-4 w-4 text-cinnabar" />
                          </Button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          {/* Pagination */}
          {total > 10 && (
            <div className="flex justify-center gap-2 mt-4">
              <Button
                variant="outline"
                size="sm"
                disabled={page <= 1}
                onClick={() => setPage(page - 1)}
              >
                上一页
              </Button>
              <span className="flex items-center px-4 text-sm text-muted-foreground">
                第 {page} 页
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
        </CardContent>
      </Card>
    </div>
  )
}
