import { FormEvent, useMemo, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import {
  ChevronLeft,
  Eye,
  MessageCircle,
  Pin,
  Search,
  Send,
  Trash2,
  UserRound,
} from 'lucide-react'
import { toast } from 'sonner'

import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Pagination } from '@/components/ui/Pagination'
import { UserAvatar } from '@/components/ui/UserAvatar'
import { useConfirmDialog } from '@/components/ui/ConfirmDialog'
import {
  useCreateForumPost,
  useCreateForumReply,
  useDeleteForumPost,
  useDeleteForumReply,
  useForumPost,
  useForumPosts,
  useForumReplies,
} from '@/hooks/useForum'
import { useAuthStore } from '@/stores/authStore'
import type { ForumPost, ForumReply } from '@/types'
import { cn } from '@/utils/cn'

const postPageSize = 20
const replyPageSize = 50

function displayName(user?: { name?: string; username?: string }) {
  return user?.name || user?.username || '诗友'
}

function formatTime(value: string) {
  return new Date(value).toLocaleString()
}

function canEditOwn(ownerId: string, currentUserId?: string | number) {
  return String(currentUserId ?? '') === ownerId
}

export function ForumPage() {
  const { id } = useParams()

  if (id) {
    return <ForumDetail postId={id} />
  }

  return <ForumList />
}

function ForumList() {
  const navigate = useNavigate()
  const { isAuthenticated } = useAuthStore()
  const [page, setPage] = useState(1)
  const [keyword, setKeyword] = useState('')
  const [draftKeyword, setDraftKeyword] = useState('')
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const postsQuery = useForumPosts({ page, pageSize: postPageSize, keyword })
  const createPostMutation = useCreateForumPost()

  const pinnedPosts = useMemo(
    () => (postsQuery.data?.list ?? []).filter((post) => post.isPinned),
    [postsQuery.data?.list],
  )

  const handleSearch = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setKeyword(draftKeyword.trim())
    setPage(1)
  }

  const handleCreatePost = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!isAuthenticated) {
      toast.error('请先登录后再发帖')
      return
    }
    const nextTitle = title.trim()
    const nextContent = content.trim()
    if (!nextTitle || !nextContent) return

    createPostMutation.mutate(
      { title: nextTitle, content: nextContent },
      {
        onSuccess: (post) => {
          setTitle('')
          setContent('')
          navigate(`/forum/${post.id}`)
        },
      },
    )
  }

  return (
    <div className="min-h-[calc(100vh-4rem)] bg-muted/20">
      <section className="border-b bg-background">
        <div className="container flex flex-col gap-5 py-6 md:flex-row md:items-end md:justify-between">
          <div>
            <div className="flex items-center gap-3">
              <div className="flex h-12 w-12 items-center justify-center rounded-md bg-cinnabar text-xl font-serif font-bold text-white">
                诗
              </div>
              <div>
                <h1 className="font-serif text-3xl font-bold text-ink">诗海吧</h1>
                <p className="mt-1 text-sm text-muted-foreground">
                  共话诗词、典故、译注与读后感
                </p>
              </div>
            </div>
          </div>
          <form className="flex w-full gap-2 md:w-80" onSubmit={handleSearch}>
            <Input
              value={draftKeyword}
              onChange={(event) => setDraftKeyword(event.target.value)}
              placeholder="搜索帖子"
              className="h-10"
            />
            <Button type="submit" size="icon" title="搜索">
              <Search className="h-4 w-4" />
            </Button>
          </form>
        </div>
      </section>

      <main className="container grid gap-6 py-6 lg:grid-cols-[minmax(0,1fr)_340px]">
        <section className="min-w-0 space-y-3">
          {pinnedPosts.length > 0 && (
            <div className="rounded-md border bg-background">
              {pinnedPosts.map((post) => (
                <PostRow key={post.id} post={post} compact />
              ))}
            </div>
          )}

          <div className="rounded-md border bg-background">
            <div className="flex h-11 items-center justify-between border-b px-4">
              <span className="text-sm font-medium text-foreground">全部帖子</span>
              <span className="text-xs text-muted-foreground">
                {postsQuery.data?.total ?? 0} 主题
              </span>
            </div>
            {postsQuery.isLoading ? (
              <div className="py-14 text-center text-sm text-muted-foreground">加载中...</div>
            ) : (postsQuery.data?.list ?? []).length === 0 ? (
              <div className="py-14 text-center text-sm text-muted-foreground">
                暂无帖子，来开第一帖吧
              </div>
            ) : (
              <div className="divide-y">
                {(postsQuery.data?.list ?? []).map((post) => (
                  <PostRow key={post.id} post={post} />
                ))}
              </div>
            )}
          </div>

          <Pagination
            page={page}
            pageSize={postPageSize}
            total={postsQuery.data?.total ?? 0}
            onPageChange={setPage}
            disabled={postsQuery.isFetching}
            className="border-t-0"
          />
        </section>

        <aside className="space-y-4">
          <Card className="rounded-md">
            <CardContent className="space-y-4 p-4">
              <div className="flex items-center justify-between">
                <h2 className="font-serif text-lg font-semibold text-ink">发表主题</h2>
                {!isAuthenticated && (
                  <Button asChild variant="outline" size="sm">
                    <Link to="/login">登录</Link>
                  </Button>
                )}
              </div>
              <form className="space-y-3" onSubmit={handleCreatePost}>
                <Input
                  value={title}
                  onChange={(event) => setTitle(event.target.value)}
                  maxLength={200}
                  disabled={!isAuthenticated}
                  placeholder={isAuthenticated ? '标题' : '登录后发帖'}
                />
                <textarea
                  value={content}
                  onChange={(event) => setContent(event.target.value)}
                  maxLength={10000}
                  disabled={!isAuthenticated}
                  placeholder="写下你的诗词话题"
                  className="min-h-36 w-full resize-y rounded-md border border-input bg-background px-3 py-2 text-sm leading-6 outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                />
                <Button
                  type="submit"
                  className="w-full"
                  disabled={!isAuthenticated || createPostMutation.isPending || !title.trim() || !content.trim()}
                >
                  <Send className="mr-2 h-4 w-4" />
                  发布
                </Button>
              </form>
            </CardContent>
          </Card>
        </aside>
      </main>
    </div>
  )
}

function PostRow({ post, compact = false }: { post: ForumPost; compact?: boolean }) {
  return (
    <Link
      to={`/forum/${post.id}`}
      className={cn(
        'grid gap-3 px-4 py-4 transition-colors hover:bg-muted/40 md:grid-cols-[88px_minmax(0,1fr)_120px]',
        compact && 'bg-primary/5',
      )}
    >
      <div className="flex items-center gap-4 md:justify-center">
        <StatBlock value={post.replyCount} label="回复" />
        <StatBlock value={post.views} label="浏览" muted />
      </div>
      <div className="min-w-0">
        <div className="flex min-w-0 items-center gap-2">
          {post.isPinned && (
            <span className="inline-flex h-5 shrink-0 items-center gap-1 rounded-sm bg-cinnabar px-1.5 text-xs text-white">
              <Pin className="h-3 w-3" />
              置顶
            </span>
          )}
          <h2 className="truncate text-base font-semibold text-foreground">{post.title}</h2>
        </div>
        <p className="mt-1 line-clamp-2 text-sm leading-6 text-muted-foreground">{post.content}</p>
      </div>
      <div className="flex items-center gap-2 text-xs text-muted-foreground md:justify-end">
        <UserAvatar
          avatar={post.user?.avatar}
          name={post.user?.name}
          username={post.user?.username}
          className="h-7 w-7"
        />
        <div className="min-w-0">
          <div className="truncate text-foreground">{displayName(post.user)}</div>
          <div className="truncate">{formatTime(post.updatedAt)}</div>
        </div>
      </div>
    </Link>
  )
}

function StatBlock({ value, label, muted = false }: { value: number; label: string; muted?: boolean }) {
  return (
    <div className={cn('w-12 text-center text-xs', muted ? 'text-muted-foreground' : 'text-cinnabar')}>
      <div className="font-semibold">{value}</div>
      <div>{label}</div>
    </div>
  )
}

function ForumDetail({ postId }: { postId: string }) {
  const navigate = useNavigate()
  const { user, isAuthenticated } = useAuthStore()
  const [replyText, setReplyText] = useState('')
  const [replyTo, setReplyTo] = useState<ForumReply | null>(null)
  const [replyPage, setReplyPage] = useState(1)
  const postQuery = useForumPost(postId)
  const repliesQuery = useForumReplies(postId, replyPage, replyPageSize)
  const createReplyMutation = useCreateForumReply(postId)
  const deletePostMutation = useDeleteForumPost()
  const deleteReplyMutation = useDeleteForumReply(postId)
  const { confirm, ConfirmDialog } = useConfirmDialog()
  const canManagePost = Boolean(postQuery.data && canEditOwn(postQuery.data.userId, user?.id))

  const handleCreateReply = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!isAuthenticated) {
      toast.error('请先登录后再回复')
      return
    }
    const content = replyText.trim()
    if (!content) return

    createReplyMutation.mutate(
      { content, parentId: replyTo?.id },
      {
        onSuccess: () => {
          setReplyText('')
          setReplyTo(null)
        },
      },
    )
  }

  const handleDeletePost = async () => {
    const confirmed = await confirm({
      title: '删除帖子',
      description: '删除后帖子将不再展示。',
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    deletePostMutation.mutate(postId, {
      onSuccess: () => navigate('/forum'),
    })
  }

  const handleDeleteReply = async (reply: ForumReply) => {
    const confirmed = await confirm({
      title: '删除回复',
      description: '删除后该回复将不再展示。',
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    deleteReplyMutation.mutate(reply.id)
  }

  if (postQuery.isLoading) {
    return (
      <div className="container py-8">
        <div className="py-20 text-center text-sm text-muted-foreground">加载中...</div>
      </div>
    )
  }

  if (!postQuery.data) {
    return (
      <div className="container py-8">
        <Button type="button" variant="ghost" onClick={() => navigate('/forum')}>
          <ChevronLeft className="mr-2 h-4 w-4" />
          返回论坛
        </Button>
        <div className="py-20 text-center text-sm text-muted-foreground">帖子不存在或已删除</div>
      </div>
    )
  }

  const post = postQuery.data

  return (
    <div className="min-h-[calc(100vh-4rem)] bg-muted/20">
      <ConfirmDialog />
      <main className="container py-6">
        <Button type="button" variant="ghost" onClick={() => navigate('/forum')} className="mb-4">
          <ChevronLeft className="mr-2 h-4 w-4" />
          返回论坛
        </Button>

        <article className="rounded-md border bg-background">
          <header className="border-b px-5 py-5">
            <div className="flex flex-wrap items-center gap-2">
              {post.isPinned && (
                <span className="inline-flex h-5 items-center gap-1 rounded-sm bg-cinnabar px-1.5 text-xs text-white">
                  <Pin className="h-3 w-3" />
                  置顶
                </span>
              )}
              <h1 className="min-w-0 flex-1 font-serif text-2xl font-bold leading-9 text-ink">
                {post.title}
              </h1>
              {canManagePost && (
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  title="删除帖子"
                  onClick={handleDeletePost}
                  disabled={deletePostMutation.isPending}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              )}
            </div>
            <div className="mt-3 flex flex-wrap items-center gap-4 text-xs text-muted-foreground">
              <span className="inline-flex items-center gap-1">
                <UserRound className="h-3.5 w-3.5" />
                {displayName(post.user)}
              </span>
              <span className="inline-flex items-center gap-1">
                <Eye className="h-3.5 w-3.5" />
                {post.views} 浏览
              </span>
              <span className="inline-flex items-center gap-1">
                <MessageCircle className="h-3.5 w-3.5" />
                {post.replyCount} 回复
              </span>
              <span>{formatTime(post.createdAt)}</span>
            </div>
          </header>
          <div className="whitespace-pre-wrap px-5 py-6 text-sm leading-8 text-foreground">
            {post.content}
          </div>
        </article>

        <section className="mt-5 rounded-md border bg-background">
          <div className="flex h-11 items-center justify-between border-b px-5">
            <h2 className="text-sm font-medium text-foreground">回帖</h2>
            <span className="text-xs text-muted-foreground">{repliesQuery.data?.total ?? 0} 楼</span>
          </div>
          {repliesQuery.isLoading ? (
            <div className="py-14 text-center text-sm text-muted-foreground">加载回帖中...</div>
          ) : (repliesQuery.data?.list ?? []).length === 0 ? (
            <div className="py-14 text-center text-sm text-muted-foreground">暂无回帖</div>
          ) : (
            <div className="divide-y">
              {(repliesQuery.data?.list ?? []).map((reply, index) => (
                <ReplyRow
                  key={reply.id}
                  reply={reply}
                  floor={(replyPage - 1) * replyPageSize + index + 1}
                  currentUserId={user?.id}
                  onReply={() => setReplyTo(reply)}
                  onDelete={() => handleDeleteReply(reply)}
                  deleting={deleteReplyMutation.isPending}
                />
              ))}
            </div>
          )}
        </section>

        <Pagination
          page={replyPage}
          pageSize={replyPageSize}
          total={repliesQuery.data?.total ?? 0}
          onPageChange={setReplyPage}
          disabled={repliesQuery.isFetching}
          className="mt-4 border-t-0"
        />

        <section className="mt-5 rounded-md border bg-background p-5">
          <form className="space-y-3" onSubmit={handleCreateReply}>
            <div className="flex items-center justify-between">
              <h2 className="font-serif text-lg font-semibold text-ink">
                {replyTo ? `回复 ${displayName(replyTo.user)}` : '发表回复'}
              </h2>
              {replyTo && (
                <Button type="button" variant="ghost" size="sm" onClick={() => setReplyTo(null)}>
                  取消
                </Button>
              )}
            </div>
            {replyTo && (
              <div className="rounded-md bg-muted px-3 py-2 text-xs leading-5 text-muted-foreground">
                {replyTo.content}
              </div>
            )}
            <textarea
              value={replyText}
              onChange={(event) => setReplyText(event.target.value)}
              maxLength={5000}
              disabled={!isAuthenticated}
              placeholder={isAuthenticated ? '写下你的回复' : '请先登录后回复'}
              className="min-h-32 w-full resize-y rounded-md border border-input bg-background px-3 py-2 text-sm leading-6 outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
            />
            <div className="flex justify-end gap-2">
              {!isAuthenticated && (
                <Button asChild variant="outline">
                  <Link to="/login">登录</Link>
                </Button>
              )}
              <Button
                type="submit"
                disabled={!isAuthenticated || createReplyMutation.isPending || !replyText.trim()}
              >
                <Send className="mr-2 h-4 w-4" />
                回复
              </Button>
            </div>
          </form>
        </section>
      </main>
    </div>
  )
}

function ReplyRow({
  reply,
  floor,
  currentUserId,
  onReply,
  onDelete,
  deleting,
}: {
  reply: ForumReply
  floor: number
  currentUserId?: string | number
  onReply: () => void
  onDelete: () => void
  deleting: boolean
}) {
  const canDelete = canEditOwn(reply.userId, currentUserId)

  return (
    <div className="grid gap-4 px-5 py-5 md:grid-cols-[160px_minmax(0,1fr)]">
      <aside className="flex items-center gap-3 md:block">
        <UserAvatar
          avatar={reply.user?.avatar}
          name={reply.user?.name}
          username={reply.user?.username}
          className="h-10 w-10 md:h-12 md:w-12"
        />
        <div className="min-w-0 md:mt-2">
          <div className="truncate text-sm font-medium text-foreground">{displayName(reply.user)}</div>
          <div className="text-xs text-muted-foreground">第 {floor} 楼</div>
        </div>
      </aside>
      <div className="min-w-0">
        {reply.parent && (
          <div className="mb-3 rounded-md bg-muted px-3 py-2 text-xs leading-5 text-muted-foreground">
            回复 {displayName(reply.parent.user)}：{reply.parent.content}
          </div>
        )}
        <p className="whitespace-pre-wrap text-sm leading-7 text-foreground">{reply.content}</p>
        <div className="mt-4 flex items-center justify-between gap-3 text-xs text-muted-foreground">
          <span>{formatTime(reply.createdAt)}</span>
          <div className="flex items-center gap-1">
            <Button type="button" variant="ghost" size="sm" onClick={onReply}>
              回复
            </Button>
            {canDelete && (
              <Button
                type="button"
                variant="ghost"
                size="icon"
                title="删除回复"
                onClick={onDelete}
                disabled={deleting}
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
