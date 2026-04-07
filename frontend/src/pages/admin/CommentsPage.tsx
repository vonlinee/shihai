import { useState } from 'react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Search, Trash2, CheckCircle, XCircle } from 'lucide-react'
import { useAdminComments } from '@/hooks/useAdmin'
import { useDeleteComment } from '@/hooks/useComments'

export function AdminCommentsPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [page] = useState(1)

  const { data: commentsData, isLoading } = useAdminComments(page, 20)
  const deleteCommentMutation = useDeleteComment()

  const comments = (commentsData?.list ?? []) as Array<{
    id: number;
    content: string;
    user?: { name: string };
    visitorName?: string;
    poemId: number;
    createdAt: string;
    isDeleted: boolean;
  }>

  const handleDelete = (id: number) => {
    if (!confirm('确定要删除该评论吗？')) return
    deleteCommentMutation.mutate(id)
  }

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-serif font-bold text-ink">评论管理</h1>
      </div>

      <Card className="ink-border">
        <CardHeader className="pb-4">
          <div className="flex items-center gap-4">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="搜索评论内容或用户..."
                className="pl-10"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
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
                    <th className="text-left py-3 px-4 font-medium">评论内容</th>
                    <th className="text-left py-3 px-4 font-medium">用户</th>
                    <th className="text-left py-3 px-4 font-medium">状态</th>
                    <th className="text-left py-3 px-4 font-medium">发布时间</th>
                    <th className="text-left py-3 px-4 font-medium">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {comments.map((comment) => (
                    <tr key={comment.id} className="border-b last:border-0 hover:bg-muted/50">
                      <td className="py-3 px-4">{comment.id}</td>
                      <td className="py-3 px-4 max-w-xs truncate">{comment.content}</td>
                      <td className="py-3 px-4">{comment.user?.name ?? comment.visitorName ?? '匿名'}</td>
                      <td className="py-3 px-4">
                        <span className={`px-2 py-1 rounded text-xs ${
                          comment.isDeleted
                            ? 'bg-red-500/10 text-red-500'
                            : 'bg-green-500/10 text-green-500'
                        }`}>
                          {comment.isDeleted ? '已删除' : '正常'}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-muted-foreground">
                        {new Date(comment.createdAt).toLocaleString()}
                      </td>
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-2">
                          <Button variant="ghost" size="sm">
                            <CheckCircle className="h-4 w-4 text-green-500" />
                          </Button>
                          <Button variant="ghost" size="sm">
                            <XCircle className="h-4 w-4 text-yellow-500" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDelete(comment.id)}
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
        </CardContent>
      </Card>
    </div>
  )
}
