import { useState } from 'react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useConfirmDialog } from '@/components/ui/ConfirmDialog'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Search, Trash2, CheckCircle, XCircle } from 'lucide-react'
import { useAdminComments } from '@/hooks/useAdmin'
import { useDeleteComment } from '@/hooks/useComments'

export function AdminCommentsPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [page] = useState(1)

  const { data: commentsData, isLoading } = useAdminComments(page, 20)
  const deleteCommentMutation = useDeleteComment()
  const { confirm, ConfirmDialog } = useConfirmDialog()

  const comments = (commentsData?.list ?? []) as Array<{
    id: number;
    content: string;
    user?: { name: string };
    visitorName?: string;
    poemId: number;
    createdAt: string;
    isDeleted: boolean;
  }>

  const handleDelete = async (id: number) => {
    const confirmed = await confirm({
      title: '删除评论',
      description: '确定要删除该评论吗？此操作不可撤销。',
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    deleteCommentMutation.mutate(id)
  }

  return (
    <div className="p-8 space-y-6">
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
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>评论内容</TableHead>
                  <TableHead>用户</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>发布时间</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                  {comments.map((comment) => (
                    <TableRow key={comment.id}>
                      <TableCell>{comment.id}</TableCell>
                      <TableCell className="max-w-xs truncate">{comment.content}</TableCell>
                      <TableCell>{comment.user?.name ?? comment.visitorName ?? '匿名'}</TableCell>
                      <TableCell>
                        <span className={`px-2 py-1 rounded text-xs ${
                          comment.isDeleted
                            ? 'bg-red-500/10 text-red-500'
                            : 'bg-green-500/10 text-green-500'
                        }`}>
                          {comment.isDeleted ? '已删除' : '正常'}
                        </span>
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {new Date(comment.createdAt).toLocaleString()}
                      </TableCell>
                      <TableCell>
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
                      </TableCell>
                    </TableRow>
                  ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
      <ConfirmDialog />
    </div>
  )
}
