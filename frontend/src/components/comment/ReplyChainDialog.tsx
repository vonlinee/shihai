import { useState } from 'react'
import { createPortal } from 'react-dom'
import { ChevronDown, ChevronUp, Reply, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import type { Comment } from '@/types'

const longCommentThreshold = 160

function commentAuthorName(comment: Comment) {
  return comment.user?.name ?? comment.visitorName ?? '匿名用户'
}

function isFeaturedComment(comment: Comment) {
  return comment.likes >= 10 && comment.likes - comment.dislikes >= 5
}

function ReplyChainCommentContent({ content }: { content: string }) {
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

export function ReplyChainDialog({ chain, onClose }: { chain: Comment[]; onClose: () => void }) {
  return createPortal(
    <div className="fixed inset-0 z-50">
      <button
        type="button"
        className="absolute inset-0 h-full w-full bg-background/80 backdrop-blur-sm"
        aria-label="关闭回复链"
        onClick={onClose}
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="reply-chain-title"
        className="absolute left-1/2 top-1/2 max-h-[calc(100vh-2rem)] w-[calc(100vw-2rem)] max-w-xl -translate-x-1/2 -translate-y-1/2 overflow-hidden rounded-lg border bg-background shadow-lg"
      >
        <div className="flex items-center justify-between gap-3 border-b px-4 py-3">
          <h4 id="reply-chain-title" className="font-serif text-base font-semibold">
            回复链
          </h4>
          <Button type="button" variant="ghost" size="icon" title="关闭回复链" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </div>
        <div className="max-h-[calc(100vh-8rem)] space-y-3 overflow-y-auto p-4">
          {chain.map((comment, index) => {
            const parent = index > 0 ? chain[index - 1] : null

            return (
              <div
                key={comment.id}
                className={
                  index === 0
                    ? 'rounded-md border bg-background p-3'
                    : 'rounded-md border border-l-4 border-l-primary/50 bg-muted/30 p-3'
                }
              >
                <div className="mb-2 flex flex-col gap-1 sm:flex-row sm:items-start sm:justify-between">
                  <div className="min-w-0">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="truncate text-sm font-medium">
                        {commentAuthorName(comment)}
                      </span>
                      {isFeaturedComment(comment) && (
                        <span className="rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-800">
                          精选
                        </span>
                      )}
                    </div>
                    {parent && (
                      <div className="mt-1 inline-flex items-center gap-1 rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary">
                        <Reply className="h-3 w-3" />
                        回复 @{commentAuthorName(parent)}
                      </div>
                    )}
                  </div>
                  <span className="shrink-0 text-xs text-muted-foreground">
                    {new Date(comment.createdAt).toLocaleString()}
                  </span>
                </div>
                <ReplyChainCommentContent content={comment.content} />
              </div>
            )
          })}
        </div>
      </div>
    </div>,
    document.body,
  )
}
