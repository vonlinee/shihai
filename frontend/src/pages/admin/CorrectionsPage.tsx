import { useMemo, useState } from 'react'
import { CheckCircle, Eye, Search, XCircle } from 'lucide-react'

import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useAdminCorrections } from '@/hooks/useAdmin'
import type { CorrectionRequest } from '@/types'

const typeLabels: Record<CorrectionRequest['type'], string> = {
  content: '原文',
  translation: '译文',
  appreciation: '赏析',
  annotation: '注释',
}

const statusLabels: Record<CorrectionRequest['status'], string> = {
  pending: '待审核',
  voting: '投票中',
  approved: '已通过',
  rejected: '已驳回',
  completed: '已完成',
}

const statusClassNames: Record<CorrectionRequest['status'], string> = {
  pending: 'bg-yellow-500/10 text-yellow-600',
  voting: 'bg-blue-500/10 text-blue-600',
  approved: 'bg-green-500/10 text-green-600',
  rejected: 'bg-red-500/10 text-red-600',
  completed: 'bg-muted text-muted-foreground',
}

function formatDateTime(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleString()
}

export function AdminCorrectionsPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const keyword = searchQuery.trim()
  const queryParams = useMemo(
    () => ({
      page: 1,
      pageSize: 20,
      keyword: keyword || undefined,
    }),
    [keyword],
  )
  const { data: correctionsData, isLoading } = useAdminCorrections(queryParams)
  const corrections = correctionsData?.list ?? []

  return (
    <div className="space-y-6 p-8">
      <Card className="ink-border">
        <CardHeader className="pb-4">
          <div className="flex items-center gap-4">
            <div className="relative max-w-sm flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="搜索诗词或用户..."
                className="pl-10"
                value={searchQuery}
                onChange={(event) => setSearchQuery(event.target.value)}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="py-8 text-center text-sm text-muted-foreground">加载中...</div>
          ) : corrections.length === 0 ? (
            <div className="py-8 text-center text-sm text-muted-foreground">暂无纠错数据</div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>诗词</TableHead>
                  <TableHead>纠错类型</TableHead>
                  <TableHead>原文</TableHead>
                  <TableHead>建议修改</TableHead>
                  <TableHead>提交用户</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>提交时间</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {corrections.map((correction) => (
                  <TableRow key={correction.id}>
                    <TableCell>{correction.id}</TableCell>
                    <TableCell className="font-medium">{correction.poem?.title ?? '未知诗词'}</TableCell>
                    <TableCell>{typeLabels[correction.type]}</TableCell>
                    <TableCell className="max-w-xs truncate text-muted-foreground">{correction.originalText}</TableCell>
                    <TableCell className="max-w-xs truncate">{correction.suggestedText}</TableCell>
                    <TableCell>{correction.user?.name || correction.user?.username || '未知用户'}</TableCell>
                    <TableCell>
                      <span className={`rounded px-2 py-1 text-xs ${statusClassNames[correction.status]}`}>
                        {statusLabels[correction.status]}
                      </span>
                    </TableCell>
                    <TableCell className="text-muted-foreground">{formatDateTime(correction.createdAt)}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Button variant="ghost" size="sm" disabled title="详情接口待接入">
                          <Eye className="h-4 w-4" />
                        </Button>
                        <Button variant="ghost" size="sm" disabled title="审核接口待接入">
                          <CheckCircle className="h-4 w-4 text-green-500" />
                        </Button>
                        <Button variant="ghost" size="sm" disabled title="驳回接口待接入">
                          <XCircle className="h-4 w-4 text-cinnabar" />
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
    </div>
  )
}
