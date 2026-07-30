import { useCallback, useMemo, useState } from 'react'
import { Edit2, Plus, Trash2 } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { DataTable, type DataTableColumn } from '@/components/ui/data-table'
import { useConfirmDialog } from '@/components/ui/ConfirmDialog'
import {
  useAdminBatchDeleteCiTunes,
  useAdminCiTunes,
  useAdminCreateCiTune,
  useAdminDeleteCiTune,
  useAdminPoemTypes,
  useAdminUpdateCiTune,
} from '@/hooks/useAdmin'
import type { CiTune } from '@/types'
import { toEntityId, toggleAllVisibleIds, toggleSelectedId } from '@/utils/uitools'

import { CiTuneFormDialog, type CiTuneFormPayload } from './CiTuneFormDialog'

export function CiTuneManagementTab() {
  const { data: ciTunes, isLoading } = useAdminCiTunes()
  const { data: poemTypes } = useAdminPoemTypes()
  const createMutation = useAdminCreateCiTune()
  const updateMutation = useAdminUpdateCiTune()
  const deleteMutation = useAdminDeleteCiTune()
  const batchDeleteMutation = useAdminBatchDeleteCiTunes()
  const { confirm, ConfirmDialog } = useConfirmDialog()
  const [showDialog, setShowDialog] = useState(false)
  const [editingCiTune, setEditingCiTune] = useState<CiTune | null>(null)
  const [selectedCiTuneIds, setSelectedCiTuneIds] = useState<string[]>([])

  const ciTuneList = ciTunes ?? []
  const visibleCiTuneIds = ciTuneList.map((ciTune) => toEntityId(ciTune.id))
  const allVisibleCiTunesSelected = visibleCiTuneIds.length > 0 && visibleCiTuneIds.every((id) => selectedCiTuneIds.includes(id))
  const poemTypeOptions = useMemo(
    () => (poemTypes ?? [])
      .filter((poemType) => poemType.category === '词')
      .map((poemType) => ({
        value: toEntityId(poemType.id),
        label: poemType.name,
        description: poemType.description,
      })),
    [poemTypes],
  )

  const openCreate = () => {
    setEditingCiTune(null)
    setShowDialog(true)
  }

  const openEdit = useCallback((ciTune: CiTune) => {
    setEditingCiTune(ciTune)
    setShowDialog(true)
  }, [])

  const closeDialog = () => {
    setShowDialog(false)
    setEditingCiTune(null)
  }

  const handleSubmit = (payload: CiTuneFormPayload) => {
    if (editingCiTune) {
      updateMutation.mutate(
        { id: toEntityId(editingCiTune.id), data: payload },
        { onSuccess: closeDialog },
      )
      return
    }
    createMutation.mutate(payload, { onSuccess: closeDialog })
  }

  const handleDelete = useCallback(
    async (id: string) => {
      const confirmed = await confirm({
        title: '删除词牌',
        description: '确定要删除该词牌吗？已有词作品上的词牌关联会被清空。',
        confirmText: '删除',
        destructive: true,
      })
      if (!confirmed) return
      deleteMutation.mutate(id, {
        onSuccess: () =>
          setSelectedCiTuneIds((ids) => ids.filter((selectedId) => selectedId !== id)),
      })
    },
    [confirm, deleteMutation],
  )

  const handleBatchDelete = async () => {
    if (selectedCiTuneIds.length === 0) return
    const confirmed = await confirm({
      title: '批量删除词牌',
      description: `确定要删除选中的 ${selectedCiTuneIds.length} 个词牌吗？已有词作品上的词牌关联会被清空。`,
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    batchDeleteMutation.mutate(selectedCiTuneIds, {
      onSuccess: () => setSelectedCiTuneIds([]),
    })
  }

  const ciTuneColumns = useMemo<DataTableColumn<CiTune>[]>(
    () => [
      {
        id: 'select',
        header: () => (
          <Checkbox
            aria-label="选择全部词牌"
            checked={allVisibleCiTunesSelected}
            onCheckedChange={() =>
              setSelectedCiTuneIds((ids) => toggleAllVisibleIds(ids, visibleCiTuneIds))
            }
          />
        ),
        cell: ({ row }) => {
          const ciTuneId = toEntityId(row.original.id)
          return (
            <Checkbox
              aria-label={`选择词牌 ${row.original.name}`}
              checked={selectedCiTuneIds.includes(ciTuneId)}
              onCheckedChange={() =>
                setSelectedCiTuneIds((ids) => toggleSelectedId(ids, ciTuneId))
              }
            />
          )
        },
        enableColumnFilter: false,
        enableSorting: false,
        meta: { align: 'center', draggable: false, fixed: 'left', headerAlign: 'center', minWidth: 40, resizable: false, width: 40 },
      },
      {
        accessorKey: 'name',
        header: '词牌名',
        cell: ({ row }) => <span className="font-medium">{row.original.name}</span>,
        meta: { filterPlaceholder: '筛选词牌', fixed: 'left', minWidth: 180 },
      },
      {
        accessorFn: (ciTune) => ciTune.poemType?.name ?? '',
        id: 'poemType',
        header: '所属词体裁',
        cell: ({ row }) => (
          <span className="text-muted-foreground">{row.original.poemType?.name ?? '未归类'}</span>
        ),
        meta: { filterPlaceholder: '筛选词体裁', width: 160 },
      },
      {
        accessorFn: (ciTune) => ciTune.aliases.join('，'),
        id: 'aliases',
        header: '别名',
        cell: ({ row }) => (
          <span className="block max-w-sm truncate text-muted-foreground">
            {row.original.aliases.length > 0 ? row.original.aliases.join('，') : '-'}
          </span>
        ),
        meta: { filterPlaceholder: '筛选别名', minWidth: 220 },
      },
      {
        accessorFn: (ciTune) => ciTune.description ?? '',
        id: 'description',
        header: '说明',
        cell: ({ row }) => (
          <span className="block max-w-md truncate text-muted-foreground">
            {row.original.description || '-'}
          </span>
        ),
        meta: { filterPlaceholder: '筛选说明', minWidth: 280 },
      },
      {
        id: 'actions',
        header: '操作',
        cell: ({ row }) => (
          <div className="flex items-center justify-end gap-1">
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 rounded-sm"
              onClick={() => openEdit(row.original)}
            >
              <Edit2 className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 rounded-sm"
              onClick={() => handleDelete(toEntityId(row.original.id))}
            >
              <Trash2 className="h-4 w-4 text-cinnabar" />
            </Button>
          </div>
        ),
        enableColumnFilter: false,
        enableSorting: false,
        meta: { align: 'right', draggable: false, fixed: 'right', headerAlign: 'center', width: 120 },
      },
    ],
    [
      allVisibleCiTunesSelected,
      handleDelete,
      openEdit,
      selectedCiTuneIds,
      visibleCiTuneIds,
    ],
  )

  return (
    <>
      <div className="mb-4 flex items-center justify-end gap-2">
        {selectedCiTuneIds.length > 0 && (
          <Button
            variant="destructive"
            onClick={handleBatchDelete}
            disabled={batchDeleteMutation.isPending}
          >
            <Trash2 className="h-4 w-4 mr-2" />
            删除选中({selectedCiTuneIds.length})
          </Button>
        )}
        <Button onClick={openCreate}>
          <Plus className="h-4 w-4 mr-2" />
          添加词牌
        </Button>
      </div>

      <Card className="ink-border">
        <CardContent className="pt-6">
          <DataTable
            columns={ciTuneColumns}
            data={ciTuneList}
            emptyText="暂无词牌数据"
            enableColumnDragging
            enableColumnFilters
            getRowId={(ciTune) => toEntityId(ciTune.id)}
            loading={isLoading}
          />
        </CardContent>
      </Card>

      {showDialog && (
        <CiTuneFormDialog
          editingCiTune={editingCiTune}
          isSubmitting={createMutation.isPending || updateMutation.isPending}
          poemTypeOptions={poemTypeOptions}
          onClose={closeDialog}
          onSubmit={handleSubmit}
        />
      )}
      <ConfirmDialog />
    </>
  )
}
