import { useMemo, useState } from 'react'
import { CheckCircle, Eye, Search, XCircle } from 'lucide-react'

import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { DataTable, type DataTableColumn } from '@/components/ui/data-table'
import { Pagination } from '@/components/ui/Pagination'
import {
  useAdminCorrection,
  useAdminCorrections,
  useAdminUpdateCorrectionStatus,
} from '@/hooks/useAdmin'
import type { CorrectionRequest } from '@/types'

const typeLabels: Record<CorrectionRequest['type'], string> = {
  title: '题目',
  author: '作者',
  dynasty: '朝代',
  content: '原文',
  translation: '译文',
  appreciation: '赏析',
  annotation: '注释',
  other: '其他',
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
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [selectedCorrectionId, setSelectedCorrectionId] = useState<string | null>(null)
  const keyword = searchQuery.trim()
  const queryParams = useMemo(
    () => ({
      page,
      pageSize,
      keyword: keyword || undefined,
    }),
    [keyword, page, pageSize],
  )
  const { data: correctionsData, isLoading } = useAdminCorrections(queryParams)
  const { data: selectedCorrection } = useAdminCorrection(selectedCorrectionId)
  const updateStatusMutation = useAdminUpdateCorrectionStatus()
  const corrections = correctionsData?.list ?? []
  const total = correctionsData?.total ?? 0
  const selectedCorrectionFallback = corrections.find((correction) => correction.id === selectedCorrectionId)
  const correctionDetail = selectedCorrection ?? selectedCorrectionFallback
  const isUpdatingStatus = updateStatusMutation.isPending

  const handleViewCorrection = (correction: CorrectionRequest) => {
    setSelectedCorrectionId((currentId) => (currentId === correction.id ? null : correction.id))
  }

  const handleUpdateStatus = (
    correction: CorrectionRequest,
    status: Extract<CorrectionRequest['status'], 'approved' | 'rejected'>,
  ) => {
    updateStatusMutation.mutate({ id: correction.id, data: { status } })
  }

  const columns = useMemo<DataTableColumn<CorrectionRequest>[]>(
    () => [
      {
        id: 'poem',
        header: '诗词',
        accessorFn: (correction) => correction.poem?.title ?? '未知诗词',
        cell: ({ row }) => (
          <span className="font-medium">{row.original.poem?.title ?? '未知诗词'}</span>
        ),
      },
      {
        accessorKey: 'type',
        header: '纠错类型',
        cell: ({ row }) => typeLabels[row.original.type],
      },
      {
        accessorKey: 'originalText',
        header: '原文',
        cell: ({ row }) => (
          <span className="block max-w-xs truncate text-muted-foreground">
            {row.original.originalText}
          </span>
        ),
      },
      {
        accessorKey: 'suggestedText',
        header: '建议修改',
        cell: ({ row }) => (
          <span className="block max-w-xs truncate">{row.original.suggestedText}</span>
        ),
      },
      {
        id: 'user',
        header: '提交用户',
        accessorFn: (correction) => correction.user?.name || correction.user?.username || '未知用户',
      },
      {
        accessorKey: 'status',
        header: '状态',
        cell: ({ row }) => (
          <span className={`rounded px-2 py-1 text-xs ${statusClassNames[row.original.status]}`}>
            {statusLabels[row.original.status]}
          </span>
        ),
      },
      {
        accessorKey: 'createdAt',
        header: '提交时间',
        enableColumnFilter: false,
        cell: ({ row }) => (
          <span className="text-muted-foreground">{formatDateTime(row.original.createdAt)}</span>
        ),
      },
      {
        id: 'actions',
        header: '操作',
        enableColumnFilter: false,
        enableSorting: false,
        cell: ({ row }) => {
          const correction = row.original
          const canReview = correction.status === 'pending' || correction.status === 'voting'

          return (
          <div className="flex items-center gap-2">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              title="查看详情"
              onClick={() => handleViewCorrection(correction)}
            >
              <Eye className="h-4 w-4" />
            </Button>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              disabled={!canReview || isUpdatingStatus}
              title={canReview ? '通过纠错' : '当前状态不可审核'}
              onClick={() => handleUpdateStatus(correction, 'approved')}
            >
              <CheckCircle className="h-4 w-4 text-green-500" />
            </Button>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              disabled={!canReview || isUpdatingStatus}
              title={canReview ? '驳回纠错' : '当前状态不可审核'}
              onClick={() => handleUpdateStatus(correction, 'rejected')}
            >
              <XCircle className="h-4 w-4 text-cinnabar" />
            </Button>
          </div>
          )
        },
      },
    ],
    [isUpdatingStatus],
  )

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
                onChange={(event) => {
                  setSearchQuery(event.target.value)
                  setPage(1)
                }}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <DataTable
            columns={columns}
            data={corrections}
            emptyText="暂无纠错数据"
            enableColumnDragging
            enableColumnFilters
            getRowId={(correction) => correction.id}
            loading={isLoading}
          />
          {correctionDetail && (
            <div className="mt-4 rounded-md border border-border bg-muted/20 p-4">
              <div className="mb-3 flex items-start justify-between gap-4">
                <div>
                  <h3 className="text-base font-semibold text-foreground">
                    {correctionDetail.poem?.title ?? '未知诗词'} · {typeLabels[correctionDetail.type]}
                  </h3>
                  <p className="mt-1 text-sm text-muted-foreground">
                    {correctionDetail.user?.name || correctionDetail.user?.username || '未知用户'} 提交于{' '}
                    {formatDateTime(correctionDetail.createdAt)}
                  </p>
                </div>
                <span className={`shrink-0 rounded px-2 py-1 text-xs ${statusClassNames[correctionDetail.status]}`}>
                  {statusLabels[correctionDetail.status]}
                </span>
              </div>
              <div className="grid gap-3 text-sm md:grid-cols-2">
                <div>
                  <div className="mb-1 font-medium text-muted-foreground">原文</div>
                  <div className="whitespace-pre-wrap rounded border border-border bg-background p-3">
                    {correctionDetail.originalText || '-'}
                  </div>
                </div>
                <div>
                  <div className="mb-1 font-medium text-muted-foreground">建议修改</div>
                  <div className="whitespace-pre-wrap rounded border border-border bg-background p-3">
                    {correctionDetail.suggestedText || '-'}
                  </div>
                </div>
              </div>
              <div className="mt-3 text-sm">
                <div className="mb-1 font-medium text-muted-foreground">纠错理由</div>
                <div className="whitespace-pre-wrap rounded border border-border bg-background p-3">
                  {correctionDetail.reason || '-'}
                </div>
              </div>
            </div>
          )}
          <Pagination
            className="mt-4"
            page={page}
            pageSize={pageSize}
            total={total}
            onPageChange={setPage}
            onPageSizeChange={(nextPageSize) => {
              setPageSize(nextPageSize)
              setPage(1)
            }}
          />
        </CardContent>
      </Card>
    </div>
  )
}
