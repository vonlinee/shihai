import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { Heart, Share2, Flag, MessageCircle, ThumbsUp, ThumbsDown, Send } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { toast } from 'sonner'
import { usePoem, useLikePoem } from '@/hooks/usePoems'
import { useComments, useCreateComment, useVoteComment } from '@/hooks/useComments'
import { useAuthStore } from '@/stores/authStore'

export function PoemDetailPage() {
  const { id } = useParams()
  const poemId = Number(id)
  const { isAuthenticated } = useAuthStore()
  const [commentText, setCommentText] = useState('')
  const [commentPage] = useState(1)

  const { data: poem, isLoading: poemLoading } = usePoem(poemId)
  const { data: commentsData, isLoading: commentsLoading } = useComments(poemId, commentPage)
  const likeMutation = useLikePoem()
  const createCommentMutation = useCreateComment()
  const voteCommentMutation = useVoteComment()

  const handleSubmitComment = () => {
    if (!commentText.trim()) return
    if (!isAuthenticated) {
      toast.error('请先登录后再发表评论')
      return
    }
    createCommentMutation.mutate(
      { poemId, content: commentText },
      { onSuccess: () => setCommentText('') },
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
            <div className="text-center py-8">
              <p className="poetry-text text-xl leading-loose whitespace-pre-line">
                {poem.content}
              </p>
            </div>

            {/* Actions */}
            <div className="flex items-center justify-center gap-4">
              <Button
                variant="outline"
                size="sm"
                onClick={() => likeMutation.mutate(poemId)}
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
              <div className="flex gap-2 mb-6">
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
                    <div key={comment.id} className="border rounded-lg p-4">
                      <div className="flex items-center justify-between mb-2">
                        <span className="font-medium text-sm">
                          {comment.user?.name ?? comment.visitorName ?? '匿名用户'}
                        </span>
                        <span className="text-xs text-muted-foreground">
                          {new Date(comment.createdAt).toLocaleString()}
                        </span>
                      </div>
                      <p className="text-sm text-foreground/80 mb-2">{comment.content}</p>
                      <div className="flex items-center gap-3">
                        <button
                          className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground"
                          onClick={() => voteCommentMutation.mutate({ commentId: comment.id, type: 'like' })}
                        >
                          <ThumbsUp className="h-3 w-3" /> {comment.likes}
                        </button>
                        <button
                          className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground"
                          onClick={() => voteCommentMutation.mutate({ commentId: comment.id, type: 'dislike' })}
                        >
                          <ThumbsDown className="h-3 w-3" /> {comment.dislikes}
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
