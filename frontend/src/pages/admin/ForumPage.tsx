import { FormEvent, useMemo, useState } from 'react'
import { Eye, Pin, PinOff, Search, Trash2 } from 'lucide-react'
import { Link } from 'react-router-dom'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Pagination } from '@/components/ui/Pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useConfirmDialog } from '@/components/ui/ConfirmDialog'
import { useAdminDeleteForumPost, useAdminForumPosts, useAdminSetForumPostPinned } from '@/hooks/useAdmin'
import type { ForumPost } from '@/types'
import { cn } from '@/utils/cn'

const pageSize = 20

function formatDateTime(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleString()
}

function authorName(post: ForumPost) {
  return post.user?.name || post.user?.username || '未知用户'
}

export function AdminForumPage() {
  const [page, setPage] = useState(1)
  const [draftKeyword, setDraftKeyword] = useState('')
  const [keyword, setKeyword] = useState('')
  const [includeDeleted, setIncludeDeleted] = useState(false)
  const params = useMemo(
    () => ({
      page,
      pageSize,
      keyword: keyword || undefined,
      includeDeleted,
    }),
    [includeDeleted, keyword, page],
  )
  const { data, isFetching, isLoading } = useAdminForumPosts(params)
  const setPinnedMutation = useAdminSetForumPostPinned()
  const deletePostMutation = useAdminDeleteForumPost()
  const { confirm, ConfirmDialog } = useConfirmDialog()

  const handleSearch = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setKeyword(draftKeyword.trim())
    setPage(1)
  }

  const handleTogglePinned = (post: ForumPost) => {
    setPinnedMutation.mutate({
      id: post.id,
      isPinned: !post.isPinned,
    })
  }

  const handleDelete = async (post: ForumPost) => {
    const confirmed = await confirm({
      title: '删除论坛帖子',
      description: `确定要删除“${post.title}”吗？删除后前台将不再展示。`,
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    deletePostMutation.mutate(post.id)
  }

  const posts = data?.list ?? []

  return (
    <div className="space-y-6 p-8">
      <Card className="ink-border">
        <CardHeader className="pb-4">
          <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
            <form className="flex max-w-md flex-1 gap-2" onSubmit={handleSearch}>
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="搜索标题或正文..."
                  className="pl-10"
                  value={draftKeyword}
                  onChange={(event) => setDraftKeyword(event.target.value)}
                />
              </div>
              <Button type="submit" disabled={isFetching}>
                搜索
              </Button>
            </form>

            <label className="flex items-center gap-2 text-sm text-muted-foreground">
              <input
                type="checkbox"
                className="h-4 w-4 rounded border-input"
                checked={includeDeleted}
                onChange={(event) => {
                  setIncludeDeleted(event.target.checked)
                  setPage(1)
                }}
              />
              包含已删除
            </label>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="py-8 text-center text-sm text-muted-foreground">加载中...</div>
          ) : posts.length === 0 ? (
            <div className="py-8 text-center text-sm text-muted-foreground">暂无论坛帖子</div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>帖子</TableHead>
                  <TableHead>作者</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>回复/浏览</TableHead>
                  <TableHead>更新时间</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {posts.map((post) => (
                  <TableRow key={post.id}>
                    <TableCell className="max-w-lg">
                      <div className="flex min-w-0 items-center gap-2">
                        {post.isPinned && <Pin className="h-4 w-4 shrink-0 text-cinnabar" />}
                        <div className="min-w-0">
                          <div className="truncate font-medium text-foreground">{post.title}</div>
                          <div className="truncate text-xs text-muted-foreground">{post.content}</div>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>{authorName(post)}</TableCell>
                    <TableCell>
                      <span
                        className={cn(
                          'rounded px-2 py-1 text-xs',
                          post.isDeleted ? 'bg-red-500/10 text-red-600' : 'bg-green-500/10 text-green-600',
                        )}
                      >
                        {post.isDeleted ? '已删除' : '正常'}
                      </span>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {post.replyCount} / {post.views}
                    </TableCell>
                    <TableCell className="text-muted-foreground">{formatDateTime(post.updatedAt)}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        <Button asChild variant="ghost" size="icon" title="查看帖子">
                          <Link to={`/forum/${post.id}`}>
                            <Eye className="h-4 w-4" />
                          </Link>
                        </Button>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          title={post.isPinned ? '取消置顶' : '置顶'}
                          disabled={post.isDeleted || setPinnedMutation.isPending}
                          onClick={() => handleTogglePinned(post)}
                        >
                          {post.isPinned ? <PinOff className="h-4 w-4" /> : <Pin className="h-4 w-4" />}
                        </Button>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          title="删除帖子"
                          disabled={post.isDeleted || deletePostMutation.isPending}
                          onClick={() => handleDelete(post)}
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

          <Pagination
            page={page}
            pageSize={pageSize}
            total={data?.total ?? 0}
            onPageChange={setPage}
            disabled={isFetching}
            className="mt-4"
          />
        </CardContent>
      </Card>
      <ConfirmDialog />
    </div>
  )
}
