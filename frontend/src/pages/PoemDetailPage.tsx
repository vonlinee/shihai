import { type FormEvent, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  ChevronDown,
  ChevronUp,
  Flag,
  Heart,
  MessageCircle,
  Reply,
  Send,
  Share2,
  ThumbsDown,
  ThumbsUp,
  X,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { toast } from 'sonner'
import { usePoem, useLikePoem } from '@/hooks/usePoems'
import { useComments, useCreateComment, useVoteComment } from '@/hooks/useComments'
import { useAuthStore } from '@/stores/authStore'
import { ReplyChainDialog } from '@/components/comment/ReplyChainDialog'
import { PoemAnnotationSection } from '@/components/poetry/annotation/PoemAnnotationSection'
import type { Comment } from '@/types'

function commentAuthorName(comment: Comment) {
  return comment.user?.name ?? comment.visitorName ?? '匿名用户'
}

type FlatReply = {
  comment: Comment
  parent: Comment
  depth: number
}

const longCommentThreshold = 160

function flattenReplies(root: Comment): FlatReply[] {
  const flatReplies: FlatReply[] = []
  const stack = [...(root.replies ?? [])]
    .reverse()
    .map((reply) => ({ comment: reply, parent: root, depth: 1 }))

  while (stack.length > 0) {
    const current = stack.pop()
    if (!current) continue

    flatReplies.push(current)

    const replies = current.comment.replies ?? []
    for (let index = replies.length - 1; index >= 0; index -= 1) {
      stack.push({
        comment: replies[index],
        parent: current.comment,
        depth: current.depth + 1,
      })
    }
  }

  return flatReplies
}

function findCommentPath(root: Comment, targetId: number): Comment[] {
  const stack = [{ comment: root, path: [root] }]

  while (stack.length > 0) {
    const current = stack.pop()
    if (!current) continue

    if (current.comment.id === targetId) {
      return current.path
    }

    const replies = current.comment.replies ?? []
    for (let index = replies.length - 1; index >= 0; index -= 1) {
      const reply = replies[index]
      stack.push({ comment: reply, path: [...current.path, reply] })
    }
  }

  return []
}

function isFeaturedComment(comment: Comment) {
  return comment.likes >= 10 && comment.likes - comment.dislikes >= 5
}

export function PoemDetailPage() {
  const { id } = useParams()
  const poemId = id
  const { isAuthenticated } = useAuthStore()
  const [commentText, setCommentText] = useState('')
  const [replyTarget, setReplyTarget] = useState<Comment | null>(null)
  const [replyText, setReplyText] = useState('')
  const [replyChain, setReplyChain] = useState<Comment[] | null>(null)
  const [commentPage] = useState(1)

  const { data: poem, isLoading: poemLoading } = usePoem(poemId)
  const { data: commentsData, isLoading: commentsLoading } = useComments(poemId, commentPage)
  const likeMutation = useLikePoem()
  const createCommentMutation = useCreateComment()
  const voteCommentMutation = useVoteComment()

  const handleSubmitComment = () => {
    if (!commentText.trim()) return
    if (!poemId) return
    if (!isAuthenticated) {
      toast.error('请先登录后再发表评论')
      return
    }
    createCommentMutation.mutate(
      {
        poemId,
        content: commentText.trim(),
      },
      {
        onSuccess: () => setCommentText(''),
      },
    )
  }

  const handleStartReply = (comment: Comment) => {
    setReplyTarget(comment)
    setReplyText('')
  }

  const handleSubmitReply = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!replyText.trim()) return
    if (!poemId || !replyTarget) return
    if (!isAuthenticated) {
      toast.error('请先登录后再回复评论')
      return
    }

    createCommentMutation.mutate(
      {
        poemId,
        content: replyText.trim(),
        parentId: replyTarget.id,
      },
      {
        onSuccess: () => {
          setReplyText('')
          setReplyTarget(null)
        },
      },
    )
  }

  if (poemLoading) {
    return (
      <div className="container py-8">
        <div className="max-w-3xl mx-auto text-center py-20 text-muted-foreground">
          加载中...
        </div>
      </div>
    )
  }

  if (!poem) {
    return (
      <div className="container py-8">
        <div className="max-w-3xl mx-auto text-center py-20 text-muted-foreground">
          诗词未找到
        </div>
      </div>
    )
  }

  return (
    <div className="container py-8">
      <div className="max-w-3xl mx-auto">
        <Card className="ink-border">
          <CardHeader className="text-center pb-6">
            <CardTitle className="font-serif text-3xl mb-2">{poem.title}</CardTitle>
            <p className="text-muted-foreground">
              [{poem.dynasty?.name}] {poem.author?.name}
            </p>
          </CardHeader>
          <CardContent className="space-y-8">
            {/* Poem Content */}
            <div className="py-8">
              <PoemAnnotationSection content={poem.content} annotations={poem.annotations ?? []} />
            </div>

            {/* Actions */}
            <div className="flex items-center justify-center gap-4">
              <Button
                variant="outline"
                size="sm"
                onClick={() => poemId && likeMutation.mutate(poemId)}
                disabled={likeMutation.isPending}
              >
                <Heart className="h-4 w-4 mr-2" />
                收藏 ({poem.favorites})
              </Button>
              <Button variant="outline" size="sm" onClick={() => toast.success('已分享')}>
                <Share2 className="h-4 w-4 mr-2" />
                分享
              </Button>
              <Button variant="outline" size="sm" onClick={() => toast.success('已提交纠错申请')}>
                <Flag className="h-4 w-4 mr-2" />
                纠错
              </Button>
            </div>

            {/* Stats */}
            <div className="flex items-center justify-center gap-6 text-sm text-muted-foreground">
              <span>{poem.views} 阅读</span>
              <span>{poem.likes} 点赞</span>
              <span>{poem.favorites} 收藏</span>
            </div>

            {/* Translation */}
            {poem.translation && (
              <div className="bg-muted/50 rounded-lg p-6">
                <h3 className="font-serif font-semibold mb-3">译文</h3>
                <p className="text-muted-foreground leading-relaxed">
                  {poem.translation}
                </p>
              </div>
            )}

            {/* Appreciation */}
            {poem.appreciation && (
              <div>
                <h3 className="font-serif font-semibold mb-3">赏析</h3>
                <p className="text-muted-foreground leading-relaxed">
                  {poem.appreciation}
                </p>
              </div>
            )}

            {/* Annotation */}
            {poem.annotation && (
              <div>
                <h3 className="font-serif font-semibold mb-3">注释</h3>
                <p className="text-muted-foreground leading-relaxed">
                  {poem.annotation}
                </p>
              </div>
            )}

            {/* Comments Section */}
            <div className="border-t pt-6">
              <h3 className="font-serif font-semibold mb-4 flex items-center gap-2">
                <MessageCircle className="h-5 w-5" />
                评论 ({commentsData?.total ?? 0})
              </h3>

              {/* Comment Input */}
              <div className="mb-6 flex gap-2">
                <Input
                  placeholder={isAuthenticated ? '发表你的看法...' : '请先登录后发表评论'}
                  value={commentText}
                  onChange={(e) => setCommentText(e.target.value)}
                  disabled={!isAuthenticated}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') handleSubmitComment()
                  }}
                />
                <Button
                  size="sm"
                  onClick={handleSubmitComment}
                  disabled={!isAuthenticated || createCommentMutation.isPending || !commentText.trim()}
                >
                  <Send className="h-4 w-4" />
                </Button>
              </div>

              {/* Comments List */}
              {commentsLoading ? (
                <p className="text-muted-foreground text-center py-4">加载评论中...</p>
              ) : (commentsData?.list ?? []).length === 0 ? (
                <p className="text-muted-foreground text-center py-8">
                  暂无评论，来发表第一条评论吧
                </p>
              ) : (
                <div className="space-y-4">
                  {(commentsData?.list ?? []).map((comment) => (
                    <CommentThread
                      key={comment.id}
                      rootComment={comment}
                      canReply={isAuthenticated}
                      activeReplyId={replyTarget?.id}
                      replyText={replyText}
                      onReply={handleStartReply}
                      onReplyTextChange={setReplyText}
                      onCancelReply={() => {
                        setReplyTarget(null)
                        setReplyText('')
                      }}
                      onSubmitReply={handleSubmitReply}
                      onVote={(commentId, type) => voteCommentMutation.mutate({ commentId, type })}
                      onShowChain={setReplyChain}
                      submittingReply={createCommentMutation.isPending}
                      voting={voteCommentMutation.isPending}
                    />
                  ))}
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
      {replyChain && (
        <ReplyChainDialog chain={replyChain} onClose={() => setReplyChain(null)} />
      )}
    </div>
  )
}

function CommentThread({
  rootComment,
  canReply,
  activeReplyId,
  replyText,
  onReply,
  onReplyTextChange,
  onCancelReply,
  onSubmitReply,
  submittingReply,
  onVote,
  onShowChain,
  voting,
}: {
  rootComment: Comment
  canReply: boolean
  activeReplyId?: number
  replyText: string
  onReply: (comment: Comment) => void
  onReplyTextChange: (value: string) => void
  onCancelReply: () => void
  onSubmitReply: (event: FormEvent<HTMLFormElement>) => void
  submittingReply: boolean
  onVote: (commentId: number, type: 'like' | 'dislike') => void
  onShowChain: (chain: Comment[]) => void
  voting: boolean
}) {
  const flatReplies = flattenReplies(rootComment)
  const isReplyingRoot = activeReplyId === rootComment.id

  return (
    <div className="rounded-lg border bg-background p-4">
      <CommentBody
        comment={rootComment}
        canReply={canReply}
        onReply={onReply}
        onVote={onVote}
        voting={voting}
      />

      {isReplyingRoot && (
        <ReplyForm
          target={rootComment}
          replyText={replyText}
          onReplyTextChange={onReplyTextChange}
          onCancelReply={onCancelReply}
          onSubmitReply={onSubmitReply}
          submittingReply={submittingReply}
        />
      )}

      {flatReplies.length > 0 && (
        <div className="mt-4 border-l-2 border-primary/25 pl-3 sm:pl-5">
          <div className="mb-3 flex items-center gap-2 text-xs text-muted-foreground">
            <Reply className="h-3.5 w-3.5 text-primary" />
            <span>{flatReplies.length} 条回复，按对话顺序平铺展示</span>
          </div>
          <div className="space-y-3">
            {flatReplies.map((flatReply) => {
              const chain = findCommentPath(rootComment, flatReply.comment.id)

              return (
                <div
                  key={flatReply.comment.id}
                  className="rounded-md border border-border/80 border-l-4 border-l-primary/50 bg-muted/30 p-3 sm:p-4"
                >
                  <CommentBody
                    comment={flatReply.comment}
                    canReply={canReply}
                    relation={{
                      parent: flatReply.parent,
                      depth: flatReply.depth,
                    }}
                    onReply={onReply}
                    onVote={onVote}
                    onShowChain={() => onShowChain(chain)}
                    voting={voting}
                  />

                  {activeReplyId === flatReply.comment.id && (
                    <ReplyForm
                      target={flatReply.comment}
                      replyText={replyText}
                      onReplyTextChange={onReplyTextChange}
                      onCancelReply={onCancelReply}
                      onSubmitReply={onSubmitReply}
                      submittingReply={submittingReply}
                    />
                  )}
                </div>
              )
            })}
          </div>
        </div>
      )}
    </div>
  )
}

function CommentBody({
  comment,
  canReply,
  relation,
  onReply,
  onVote,
  onShowChain,
  voting,
}: {
  comment: Comment
  canReply: boolean
  relation?: {
    parent: Comment
    depth: number
  }
  onReply: (comment: Comment) => void
  onVote: (commentId: number, type: 'like' | 'dislike') => void
  onShowChain?: () => void
  voting: boolean
}) {
  return (
    <div>
      <div className="mb-2 flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
        <div className="min-w-0 space-y-1">
          <div className="flex flex-wrap items-center gap-2">
            <span className="min-w-0 truncate text-sm font-medium">
              {commentAuthorName(comment)}
            </span>
            {isFeaturedComment(comment) && (
              <span className="rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-800">
                精选
              </span>
            )}
          </div>
          {relation && (
            <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
              <span className="inline-flex items-center gap-1 rounded-full bg-primary/10 px-2 py-0.5 text-primary">
                <Reply className="h-3 w-3" />
                回复 @{commentAuthorName(relation.parent)}
              </span>
              {relation.depth > 1 && (
                <span className="rounded-full bg-background px-2 py-0.5">
                  来自第 {relation.depth} 层对话
                </span>
              )}
            </div>
          )}
        </div>
        <span className="shrink-0 text-xs text-muted-foreground">
          {new Date(comment.createdAt).toLocaleString()}
        </span>
      </div>

      <CollapsibleCommentContent content={comment.content} />

      <div className="flex flex-wrap items-center gap-2">
        <button
          type="button"
          className="flex h-7 items-center gap-1 rounded px-1.5 text-xs text-muted-foreground hover:bg-muted hover:text-foreground disabled:cursor-not-allowed disabled:opacity-60"
          onClick={() => onVote(comment.id, 'like')}
          disabled={voting}
        >
          <ThumbsUp className="h-3 w-3" /> {comment.likes}
        </button>
        <button
          type="button"
          className="flex h-7 items-center gap-1 rounded px-1.5 text-xs text-muted-foreground hover:bg-muted hover:text-foreground disabled:cursor-not-allowed disabled:opacity-60"
          onClick={() => onVote(comment.id, 'dislike')}
          disabled={voting}
        >
          <ThumbsDown className="h-3 w-3" /> {comment.dislikes}
        </button>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="h-7 px-2 text-xs text-muted-foreground hover:text-foreground"
          onClick={() => onReply(comment)}
          disabled={!canReply}
          title={canReply ? '回复评论' : '请先登录后回复'}
        >
          <Reply className="mr-1 h-3 w-3" />
          回复{comment.replyCount > 0 ? ` (${comment.replyCount})` : ''}
        </Button>
        {onShowChain && (
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="h-7 px-2 text-xs text-muted-foreground hover:text-foreground"
            onClick={onShowChain}
          >
            查看对话
          </Button>
        )}
      </div>
    </div>
  )
}

function CollapsibleCommentContent({ content }: { content: string }) {
  const [isExpanded, setIsExpanded] = useState(false)
  const isLong = content.length > longCommentThreshold || content.split(/\r?\n/).length > 4
  const contentClassName =
    isLong && !isExpanded
      ? 'mb-2 max-h-24 overflow-hidden whitespace-pre-wrap break-words text-sm leading-6 text-foreground/80'
      : 'mb-2 whitespace-pre-wrap break-words text-sm leading-6 text-foreground/80'

  return (
    <div className="mb-2">
      <p className={contentClassName}>{content}</p>
      {isLong && (
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="h-7 px-2 text-xs text-muted-foreground hover:text-foreground"
          onClick={() => setIsExpanded((value) => !value)}
        >
          {isExpanded ? (
            <>
              <ChevronUp className="mr-1 h-3 w-3" />
              收起
            </>
          ) : (
            <>
              <ChevronDown className="mr-1 h-3 w-3" />
              展开
            </>
          )}
        </Button>
      )}
    </div>
  )
}

function ReplyForm({
  target,
  replyText,
  onReplyTextChange,
  onCancelReply,
  onSubmitReply,
  submittingReply,
}: {
  target: Comment
  replyText: string
  onReplyTextChange: (value: string) => void
  onCancelReply: () => void
  onSubmitReply: (event: FormEvent<HTMLFormElement>) => void
  submittingReply: boolean
}) {
  return (
    <form className="mt-3 space-y-2 rounded-md bg-muted/40 p-3" onSubmit={onSubmitReply}>
      <div className="flex items-center justify-between gap-2 text-xs text-muted-foreground">
        <span className="min-w-0 truncate">回复 {commentAuthorName(target)}</span>
        <Button type="button" variant="ghost" size="icon" title="取消回复" onClick={onCancelReply}>
          <X className="h-4 w-4" />
        </Button>
      </div>
      <div className="flex flex-col gap-2 sm:flex-row">
        <Input
          value={replyText}
          onChange={(event) => onReplyTextChange(event.target.value)}
          placeholder="写下你的回复..."
          autoFocus
        />
        <Button type="submit" size="sm" disabled={submittingReply || !replyText.trim()}>
          <Send className="h-4 w-4" />
        </Button>
      </div>
    </form>
  )
}
