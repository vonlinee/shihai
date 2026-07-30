import { useCallback, useEffect, useMemo, useState } from 'react'
import { createPortal } from 'react-dom'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { Search, Plus, Edit2, Trash2, Eye, X, Wand2 } from 'lucide-react'

import { ChineseVariantToggle } from '@/components/poetry/ChineseVariantToggle'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Combobox, type ComboboxOption } from '@/components/ui/combobox'
import { DataTable, type DataTableColumn } from '@/components/ui/data-table'
import { Input } from '@/components/ui/input'
import { Pagination } from '@/components/ui/Pagination'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useConfirmDialog } from '@/components/ui/ConfirmDialog'
import {
  useAdminBatchDeletePoems,
  useAdminConvertTexts,
  useAdminCreatePoem,
  useAdminDeletePoem,
  useAdminPoemAnnotations,
  useAdminRecognizePingze,
  useAdminUpdatePoem,
} from '@/hooks/useAdmin'
import { useDynasties, useGenreCategories, usePoems, usePoetList } from '@/hooks/usePoems'
import type {
  PoemAnnotationUpsertRequest,
  PoemCreateRequest,
  PoemUpdateRequest,
  TextConversionMode,
} from '@/services/adminService'
import type { GenreCategory, Poem } from '@/types'
import { splitPoemContentInput } from '@/utils/poemContent'
import { toEntityId, toggleAllVisibleIds, toggleSelectedId } from '@/utils/uitools'

import { PoemContentAnnotationEditor } from './annotation/PoemContentAnnotationEditor'
import {
  PoemAnnotationManager,
  toAnnotationDrafts,
  type PoemAnnotationDraft,
} from './annotation/PoemAnnotationManager'

type PoemFormData = {
  title: string
  content: string
  pingze: string
  authorId?: string
  authorName?: string
  dynastyId?: string
  dynastyName?: string
  genreCategory?: string
  genre: string | number | undefined
  translation: string
  appreciation: string
  annotation: string
}
type PoemTextField = 'translation' | 'appreciation' | 'annotation'
type ConvertiblePoemTextField = 'title' | PoemTextField
type PoemConversionTarget = ConvertiblePoemTextField | 'content' | 'poemAnnotations'
type PoemFormTabKey = PoemTextField | 'pingze' | 'poemAnnotations'

const emptyPoemFormData: PoemFormData = {
  title: '',
  content: '',
  pingze: '',
  genre: '',
  translation: '',
  appreciation: '',
  annotation: '',
}

const poemTextFields: { key: PoemTextField; label: string; placeholder: string }[] = [
  { key: 'translation', label: '译文', placeholder: '白话文翻译' },
  { key: 'appreciation', label: '赏析', placeholder: '诗词赏析' },
  { key: 'annotation', label: '注释', placeholder: '字词注释' },
]

const poemFormTabs: { key: PoemFormTabKey; label: string; placeholder?: string }[] = [
  ...poemTextFields,
  { key: 'pingze', label: '平仄' },
  { key: 'poemAnnotations', label: '诗词标注' },
]

const POET_SELECT_PAGE_SIZE = 50

function createEmptyPoemFormData(): PoemFormData {
  return { ...emptyPoemFormData }
}

function formatPoemContentForInput(content: Poem['content']): string {
  return content.join('\n')
}

function formatPoemPingzeForInput(pingze: Poem['pingze'] | null | undefined): string {
  return (pingze ?? []).join('\n')
}

function splitPoemPingzeInput(value: string): string[] {
  return value.split(/\r?\n/).map((line) => line.trim())
}

function hasPoemPingzeInput(value: string): boolean {
  return splitPoemPingzeInput(value).some((line) => line !== '')
}

function isPoemPingzeMark(mark: string): boolean {
  return mark === '平' || mark === '仄' || mark === '?' || !/^\p{Script=Han}$/u.test(mark)
}

function getInvalidPoemPingzeLine(value: string): number | null {
  const lines = splitPoemPingzeInput(value)
  const invalidLineIndex = lines.findIndex((line) => [...line].some((mark) => !isPoemPingzeMark(mark)))
  return invalidLineIndex >= 0 ? invalidLineIndex + 1 : null
}

function toPoemUpdateId(id: string | undefined): string | undefined {
  return id && id.trim() ? id : undefined
}

function toGenreValue(genre: string | number | undefined): string {
  return typeof genre === 'number' ? String(genre) : genre || ''
}

function findGenreCategoryName(
  genreCategories: GenreCategory[] | undefined,
  genre: string | number | undefined,
): string | undefined {
  const genreValue = toGenreValue(genre)
  if (!genreValue) return undefined
  return genreCategories?.find((category) => category.genres.some((item) => item.name === genreValue))?.name
}

function formatGenreDescription(genre: GenreCategory['genres'][number]): string | undefined {
  const details = [
    genre.lines ? `${genre.lines}句` : '',
    genre.charsPerLine ? `${genre.charsPerLine}字` : '',
  ].filter(Boolean)
  if (details.length > 0) return details.join(' · ')
  return genre.description || undefined
}

function formatTableDateTime(value: string | undefined): string {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  const pad = (numberValue: number) => String(numberValue).padStart(2, '0')
  return [
    `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`,
    `${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`,
  ].join(' ')
}

function toPoemAnnotationPayloads(annotations: PoemAnnotationDraft[]): PoemAnnotationUpsertRequest[] {
  return annotations.map((annotation) => ({
    id: annotation.id,
    targetField: annotation.targetField,
    startLine: annotation.startLine,
    startOffset: annotation.startOffset,
    endLine: annotation.endLine,
    endOffset: annotation.endOffset,
    selectedText: annotation.selectedText,
    title: annotation.title,
    content: annotation.content,
    type: annotation.type,
    displayOrder: annotation.displayOrder,
  }))
}

function formatPoemAnnotationLocation(annotation: PoemAnnotationDraft): string {
  const startLine = annotation.startLine + 1
  const endLine = annotation.endLine + 1
  const startColumn = annotation.startOffset + 1
  const endColumn = Math.max(1, annotation.endOffset)

  if (annotation.startLine === annotation.endLine) {
    return `第 ${startLine} 行 · ${startColumn}-${endColumn} 列`
  }

  return `第 ${startLine}-${endLine} 行 · ${annotation.endLine - annotation.startLine + 1} 行`
}

function PoemAnnotationDraftList({
  annotations,
  onDelete,
}: {
  annotations: PoemAnnotationDraft[]
  onDelete: (index: number) => void
}) {
  if (annotations.length === 0) {
    return (
      <div className="rounded-md border border-dashed py-8 text-center text-sm text-muted-foreground">
        暂无诗词标注，选中正文中的文本后可添加
      </div>
    )
  }

  return (
    <div className="space-y-2">
      {annotations.map((annotation, index) => (
        <div
          key={annotation.id ?? `${annotation.startLine}-${annotation.startOffset}-${index}`}
          className="rounded-md border p-3"
        >
          <div className="flex items-start justify-between gap-3">
            <div className="min-w-0">
              <div className="flex min-w-0 flex-wrap items-center gap-x-2 gap-y-1">
                <span className="font-medium">
                  {index + 1}. {annotation.title || annotation.selectedText}
                </span>
                <span className="rounded bg-muted px-1.5 py-0.5 text-xs text-muted-foreground">
                  {formatPoemAnnotationLocation(annotation)}
                </span>
              </div>
            </div>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="shrink-0"
              onClick={() => onDelete(index)}
            >
              <Trash2 className="h-4 w-4 text-cinnabar" />
            </Button>
          </div>
          <p className="mt-2 text-sm leading-6">{annotation.content || '暂无标注内容'}</p>
        </div>
      ))}
    </div>
  )
}

export function PoemsTab() {
  const navigate = useNavigate()
  const [searchQuery, setSearchQuery] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [showPoemDialog, setShowPoemDialog] = useState(false)
  const [editingPoemId, setEditingPoemId] = useState<string | null>(null)
  const [annotatingPoem, setAnnotatingPoem] = useState<Poem | null>(null)
  const [formData, setFormData] = useState<PoemFormData>(() => createEmptyPoemFormData())
  const [annotationDrafts, setAnnotationDrafts] = useState<PoemAnnotationDraft[]>([])
  const [annotationDraftSourceId, setAnnotationDraftSourceId] = useState<string | null>(null)
  const [convertingTarget, setConvertingTarget] = useState<PoemConversionTarget | null>(null)
  const [selectedPoemIds, setSelectedPoemIds] = useState<string[]>([])
  const [poetSearchKeyword, setPoetSearchKeyword] = useState('')
  const [debouncedPoetSearchKeyword, setDebouncedPoetSearchKeyword] = useState('')

  const { data: poemData, isFetching: isPoemsFetching } = usePoems({ keyword: searchQuery || undefined, page, pageSize })
  const editingAnnotationPoemId = showPoemDialog && editingPoemId ? editingPoemId : null
  const { data: editingAnnotations, isLoading: isEditingAnnotationsLoading } = useAdminPoemAnnotations(editingAnnotationPoemId)
  const { data: dynasties } = useDynasties()
  const { data: poetData, isFetching: isPoetsFetching } = usePoetList({
    keyword: debouncedPoetSearchKeyword || undefined,
    dynastyId: formData.dynastyId,
    page: 1,
    pageSize: POET_SELECT_PAGE_SIZE,
  }, { enabled: showPoemDialog && Boolean(formData.dynastyId) })
  const { data: genreCategories } = useGenreCategories()

  const dynastyOptions: ComboboxOption[] = useMemo(
    () => (dynasties ?? []).map((d) => ({ value: String(d.id), label: d.name, description: d.period })),
    [dynasties],
  )
  const authorDynastyOptions: ComboboxOption[] = useMemo(
    () => (dynasties ?? []).map((d) => ({ value: String(d.id), label: d.name })),
    [dynasties],
  )
  const poetOptions: ComboboxOption[] = useMemo(() => {
    const options = (poetData?.list ?? [])
      .filter((p) => !formData.dynastyId || String(p.dynastyId) === String(formData.dynastyId))
      .map((p) => ({ value: String(p.authorId), label: p.name, description: p.dynasty?.name }))
    if (formData.authorId && formData.authorName && !options.some((option) => String(option.value) === String(formData.authorId))) {
      return [{ value: formData.authorId, label: formData.authorName, description: formData.dynastyName }, ...options]
    }
    return options
  }, [formData.authorId, formData.authorName, formData.dynastyId, formData.dynastyName, poetData?.list])
  const genreCategoryOptions: ComboboxOption[] = useMemo(
    () => (genreCategories ?? []).map((category) => ({ value: category.name, label: category.name })),
    [genreCategories],
  )
  const selectedGenreCategory = formData.genreCategory || findGenreCategoryName(genreCategories, formData.genre)
  const selectedGenreCategoryData = useMemo(
    () => (genreCategories ?? []).find((category) => category.name === selectedGenreCategory),
    [genreCategories, selectedGenreCategory],
  )
  const genreOptions: ComboboxOption[] = useMemo(() => {
    const options = (selectedGenreCategoryData?.genres ?? []).map((genre) => ({
      value: genre.name,
      label: genre.name,
      description: formatGenreDescription(genre),
    }))
    const genreValue = toGenreValue(formData.genre)
    if (genreValue && selectedGenreCategory && !options.some((option) => option.value === genreValue)) {
      return [{ value: genreValue, label: genreValue, description: '自定义体裁' }, ...options]
    }
    return options
  }, [formData.genre, selectedGenreCategory, selectedGenreCategoryData?.genres])

  const deletePoemMutation = useAdminDeletePoem()
  const batchDeletePoemsMutation = useAdminBatchDeletePoems()
  const createPoemMutation = useAdminCreatePoem()
  const updatePoemMutation = useAdminUpdatePoem()
  const convertTextsMutation = useAdminConvertTexts()
  const recognizePingzeMutation = useAdminRecognizePingze()
  const { confirm, ConfirmDialog } = useConfirmDialog()
  const poems = poemData?.list ?? []
  const total = poemData?.total ?? 0
  const visiblePoemIds = poems.map((poem) => toEntityId(poem.id))
  const allVisiblePoemsSelected = visiblePoemIds.length > 0 && visiblePoemIds.every((id) => selectedPoemIds.includes(id))
  const poemContentLineCount = Math.max(3, splitPoemContentInput(formData.content).length)
  const invalidPingzeLine = getInvalidPoemPingzeLine(formData.pingze)

  useEffect(() => {
    const timer = window.setTimeout(() => {
      setDebouncedPoetSearchKeyword(poetSearchKeyword.trim())
    }, 250)
    return () => window.clearTimeout(timer)
  }, [poetSearchKeyword])

  useEffect(() => {
    if (!showPoemDialog || !editingPoemId || !editingAnnotations || annotationDraftSourceId === editingPoemId) {
      return
    }
    setAnnotationDrafts(toAnnotationDrafts(editingAnnotations))
    setAnnotationDraftSourceId(editingPoemId)
  }, [annotationDraftSourceId, editingAnnotations, editingPoemId, showPoemDialog])

  useEffect(() => {
    if (!showPoemDialog) return

    const previousBodyOverflow = document.body.style.overflow
    const previousDocumentOverflow = document.documentElement.style.overflow
    document.body.style.overflow = 'hidden'
    document.documentElement.style.overflow = 'hidden'

    return () => {
      document.body.style.overflow = previousBodyOverflow
      document.documentElement.style.overflow = previousDocumentOverflow
    }
  }, [showPoemDialog])

  const handleDelete = useCallback(async (id: string) => {
    const confirmed = await confirm({
      title: '删除诗词',
      description: '确定要删除该诗词吗？此操作不可撤销。',
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    deletePoemMutation.mutate(id, {
      onSuccess: () => setSelectedPoemIds((ids) => ids.filter((selectedId) => selectedId !== id)),
    })
  }, [confirm, deletePoemMutation])

  const handleBatchDelete = async () => {
    if (selectedPoemIds.length === 0) return
    const confirmed = await confirm({
      title: '批量删除诗词',
      description: `确定要删除选中的 ${selectedPoemIds.length} 首诗词吗？此操作不可撤销。`,
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    batchDeletePoemsMutation.mutate(selectedPoemIds, { onSuccess: () => setSelectedPoemIds([]) })
  }

  const closePoemDialog = () => {
    setShowPoemDialog(false)
  }

  const openCreatePoemDialog = () => {
    setEditingPoemId(null)
    setFormData(createEmptyPoemFormData())
    setAnnotationDrafts([])
    setAnnotationDraftSourceId(null)
    setPoetSearchKeyword('')
    setDebouncedPoetSearchKeyword('')
    setShowPoemDialog(true)
  }

  const openEditPoemDialog = useCallback((poem: Poem) => {
    setEditingPoemId(toEntityId(poem.id))
    setFormData({
      title: poem.title,
      content: formatPoemContentForInput(poem.content),
      pingze: formatPoemPingzeForInput(poem.pingze),
      authorId: poem.authorId ? String(poem.authorId) : undefined,
      authorName: poem.author?.name,
      dynastyId: poem.dynastyId ? String(poem.dynastyId) : undefined,
      dynastyName: poem.dynasty?.name,
      genreCategory: poem.genreCategory || findGenreCategoryName(genreCategories, poem.genre),
      genre: poem.genre || '',
      translation: poem.translation || '',
      appreciation: poem.appreciation || '',
      annotation: poem.annotation || '',
    })
    setAnnotationDrafts(toAnnotationDrafts(poem.annotations ?? []))
    setAnnotationDraftSourceId(null)
    setPoetSearchKeyword('')
    setDebouncedPoetSearchKeyword('')
    setShowPoemDialog(true)
  }, [genreCategories])

  const resetPoemDialog = () => {
    setEditingPoemId(null)
    setFormData(createEmptyPoemFormData())
    setAnnotationDrafts([])
    setAnnotationDraftSourceId(null)
    setPoetSearchKeyword('')
    setDebouncedPoetSearchKeyword('')
    setShowPoemDialog(false)
  }

  const buildCreatePayload = (): PoemCreateRequest => {
    const payload: PoemCreateRequest = {
      title: formData.title,
      content: splitPoemContentInput(formData.content),
      genreCategory: selectedGenreCategory || '',
      genre: toGenreValue(formData.genre),
      translation: formData.translation,
      appreciation: formData.appreciation,
      annotation: formData.annotation,
    }
    if (hasPoemPingzeInput(formData.pingze)) payload.pingze = splitPoemPingzeInput(formData.pingze)
    const annotations = toPoemAnnotationPayloads(annotationDrafts)
    if (annotations.length > 0) payload.annotations = annotations
    if (formData.dynastyName) payload.dynastyName = formData.dynastyName
    if (formData.authorName) payload.authorName = formData.authorName
    return payload
  }

  const buildUpdatePayload = (): PoemUpdateRequest => {
    const payload: PoemUpdateRequest = {
      title: formData.title,
      content: splitPoemContentInput(formData.content),
      genreCategory: selectedGenreCategory || '',
      genre: toGenreValue(formData.genre),
      translation: formData.translation,
      appreciation: formData.appreciation,
      annotation: formData.annotation,
      annotations: toPoemAnnotationPayloads(annotationDrafts),
    }
    if (hasPoemPingzeInput(formData.pingze)) payload.pingze = splitPoemPingzeInput(formData.pingze)
    const dynastyId = toPoemUpdateId(formData.dynastyId)
    const authorId = toPoemUpdateId(formData.authorId)
    if (dynastyId) payload.dynastyId = dynastyId
    if (authorId) payload.authorId = authorId
    return payload
  }

  const handlePoemSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (invalidPingzeLine !== null) {
      toast.error(`第 ${invalidPingzeLine} 行平仄只能填写“平”“仄”“?”或非汉字字符`)
      return
    }
    const genreValue = toGenreValue(formData.genre).trim()
    if (genreValue && !selectedGenreCategory) {
      toast.error('请选择体裁大类')
      return
    }
    if (selectedGenreCategory && !genreValue) {
      toast.error('请选择或输入细分类别')
      return
    }
    if (editingPoemId) {
      updatePoemMutation.mutate(
        { id: editingPoemId, data: buildUpdatePayload() },
        { onSuccess: resetPoemDialog },
      )
      return
    }
    createPoemMutation.mutate(buildCreatePayload(), { onSuccess: resetPoemDialog })
  }

  const handleConvertPoemContent = (mode: TextConversionMode) => {
    if (convertTextsMutation.isPending) return
    const contentLines = formData.content.split(/\r?\n/)
    setConvertingTarget('content')
    convertTextsMutation.mutate(
      { mode, texts: contentLines },
      {
        onSuccess: (response) => {
          if (response.texts.length !== contentLines.length) {
            return
          }
          setFormData((current) => ({ ...current, content: response.texts.join('\n') }))
        },
        onSettled: () => setConvertingTarget(null),
      },
    )
  }

  const handleConvertPoemTextField = (field: ConvertiblePoemTextField, mode: TextConversionMode) => {
    if (convertTextsMutation.isPending) return
    const sourceText = formData[field]
    setConvertingTarget(field)
    convertTextsMutation.mutate(
      { mode, texts: [sourceText] },
      {
        onSuccess: (response) => {
          const [convertedText] = response.texts
          if (convertedText === undefined) {
            return
          }
          setFormData((current) => ({ ...current, [field]: convertedText }))
        },
        onSettled: () => setConvertingTarget(null),
      },
    )
  }

  const handleConvertPoemAnnotations = (mode: TextConversionMode) => {
    if (convertTextsMutation.isPending) return
    const textEntries = annotationDrafts.flatMap((annotation, index) => {
      const entries: { index: number; key: 'selectedText' | 'title' | 'content'; text: string }[] = [
        { index, key: 'selectedText', text: annotation.selectedText },
        { index, key: 'content', text: annotation.content },
      ]
      if (annotation.title !== undefined) {
        entries.push({ index, key: 'title', text: annotation.title })
      }
      return entries
    })

    if (textEntries.length === 0) {
      return
    }

    setConvertingTarget('poemAnnotations')
    convertTextsMutation.mutate(
      { mode, texts: textEntries.map((entry) => entry.text) },
      {
        onSuccess: (response) => {
          if (response.texts.length !== textEntries.length) {
            return
          }
          setAnnotationDrafts((drafts) => {
            const nextDrafts = drafts.map((draft) => ({ ...draft }))
            textEntries.forEach((entry, index) => {
              const draft = nextDrafts[entry.index]
              if (!draft) {
                return
              }
              draft[entry.key] = response.texts[index]
            })
            return nextDrafts
          })
        },
        onSettled: () => setConvertingTarget(null),
      },
    )
  }

  const handleRecognizePoemPingze = () => {
    if (recognizePingzeMutation.isPending) return
    const contentLines = splitPoemContentInput(formData.content)
    if (contentLines.length === 0) {
      return
    }
    recognizePingzeMutation.mutate(
      { texts: contentLines },
      {
        onSuccess: (response) => {
          setFormData((current) => ({ ...current, pingze: response.pingze.join('\n') }))
          toast.success('平仄识别完成')
        },
      },
    )
  }

  const poemColumns = useMemo<DataTableColumn<Poem>[]>(
    () => [
      {
        id: 'select',
        header: () => (
          <Checkbox
            aria-label="选择当前页诗词"
            checked={allVisiblePoemsSelected}
            onCheckedChange={() => setSelectedPoemIds((ids) => toggleAllVisibleIds(ids, visiblePoemIds))}
          />
        ),
        cell: ({ row }) => (
          <Checkbox
            aria-label={`选择诗词 ${row.original.title}`}
            checked={selectedPoemIds.includes(toEntityId(row.original.id))}
            onCheckedChange={() => setSelectedPoemIds((ids) => toggleSelectedId(ids, toEntityId(row.original.id)))}
          />
        ),
        enableColumnFilter: false,
        enableSorting: false,
        meta: { align: 'center', draggable: false, fixed: 'left', headerAlign: 'center', minWidth: 40, resizable: false, width: 40 },
      },
      {
        accessorKey: 'title',
        header: '标题',
        cell: ({ row }) => <span className="font-medium">{row.original.title}</span>,
        meta: { filterPlaceholder: '筛选标题', fixed: 'left', width: 220 },
      },
      {
        accessorFn: (poem) => poem.author?.name ?? '',
        id: 'author',
        header: '作者',
        cell: ({ row }) => row.original.author?.name,
        meta: { filterPlaceholder: '筛选作者', width: 160 },
      },
      {
        accessorFn: (poem) => poem.dynasty?.name ?? '',
        id: 'dynasty',
        header: '朝代',
        cell: ({ row }) => row.original.dynasty?.name,
        meta: { filterPlaceholder: '筛选朝代', width: 140 },
      },
      {
        accessorKey: 'genre',
        header: '体裁',
        meta: { filterPlaceholder: '筛选体裁', width: 140 },
      },
      {
        accessorKey: 'views',
        header: '浏览量',
        meta: { align: 'center', width: 120 },
      },
      {
        accessorKey: 'createdAt',
        header: '创建时间',
        cell: ({ row }) => <span className="text-muted-foreground">{formatTableDateTime(row.original.createdAt)}</span>,
        meta: { filterPlaceholder: '筛选创建时间', filterVariant: 'dateRange', width: 150 },
      },
      {
        accessorKey: 'updatedAt',
        header: '最后更新时间',
        cell: ({ row }) => <span className="text-muted-foreground">{formatTableDateTime(row.original.updatedAt)}</span>,
        meta: { filterPlaceholder: '筛选最后更新时间', filterVariant: 'dateRange', width: 150 },
      },
      {
        id: 'actions',
        header: '操作',
        cell: ({ row }) => (
          <div className="flex items-center justify-end gap-1">
            <Button variant="ghost" size="icon" className="h-7 w-7 rounded-sm" onClick={() => navigate(`/poems/${row.original.id}`)}><Eye className="h-4 w-4" /></Button>
            <Button variant="ghost" size="icon" className="h-7 w-7 rounded-sm" onClick={() => openEditPoemDialog(row.original)}><Edit2 className="h-4 w-4" /></Button>
            <Button variant="ghost" size="icon" className="h-7 w-7 rounded-sm" onClick={() => handleDelete(toEntityId(row.original.id))}><Trash2 className="h-4 w-4 text-cinnabar" /></Button>
          </div>
        ),
        enableColumnFilter: false,
        enableSorting: false,
        meta: { align: 'right', draggable: false, fixed: 'right', headerAlign: 'center', width: 132 },
      },
    ],
    [allVisiblePoemsSelected, handleDelete, navigate, openEditPoemDialog, selectedPoemIds, visiblePoemIds],
  )

  return (
    <>
      <div className="flex items-center justify-between">
        <div className="relative max-w-sm flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input placeholder="搜索诗词标题或作者..." className="pl-10" value={searchQuery}
            onChange={(e) => { setSearchQuery(e.target.value); setSelectedPoemIds([]); setPage(1) }} />
        </div>
        <div className="flex items-center gap-2">
          {selectedPoemIds.length > 0 && (
            <Button
              variant="destructive"
              onClick={handleBatchDelete}
              disabled={batchDeletePoemsMutation.isPending}
            >
              <Trash2 className="h-4 w-4 mr-2" />
              删除选中({selectedPoemIds.length})
            </Button>
          )}
          <Button onClick={openCreatePoemDialog}><Plus className="h-4 w-4 mr-2" />添加诗词</Button>
        </div>
      </div>

      <Card className="ink-border">
        <CardContent className="pt-6">
          <DataTable
            columns={poemColumns}
            data={poems}
            emptyText="暂无诗词数据"
            enableColumnDragging
            enableColumnFilters
            getRowId={(poem) => toEntityId(poem.id)}
            loading={isPoemsFetching}
          />
          <Pagination
            className="mt-4"
            page={page}
            pageSize={pageSize}
            total={total}
            disabled={isPoemsFetching}
            onPageChange={(nextPage) => {
              setSelectedPoemIds([])
              setPage(nextPage)
            }}
            onPageSizeChange={(nextPageSize) => {
              setSelectedPoemIds([])
              setPageSize(nextPageSize)
              setPage(1)
            }}
          />
        </CardContent>
      </Card>

      {showPoemDialog && createPortal(
        <div className="fixed left-0 top-0 z-[100] flex h-screen w-screen items-center justify-center bg-black/50">
          <div className="bg-background rounded-lg w-full max-w-5xl max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">{editingPoemId ? '编辑诗词' : '添加诗词'}</h2>
              <button onClick={closePoemDialog} className="p-1 hover:bg-muted rounded"><X className="h-5 w-5" /></button>
            </div>
            <form onSubmit={handlePoemSubmit} className="space-y-4 p-6 pt-4 overflow-y-auto">
              <div className="flex items-center gap-3">
                <label htmlFor="poem-title" className="shrink-0 text-sm font-medium">标题 *</label>
                <Input id="poem-title" value={formData.title} onChange={(e) => setFormData({ ...formData, title: e.target.value })} placeholder="诗词标题" required />
                <ChineseVariantToggle
                  disabled={!formData.title.trim()}
                  isLoading={convertingTarget === 'title'}
                  onChange={(_, mode) => handleConvertPoemTextField('title', mode)}
                />
              </div>
              <div className="grid grid-cols-3 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">体裁</label>
                  <Combobox
                    options={genreOptions}
                    value={formData.genre}
                    onChange={(val) => setFormData({ ...formData, genre: val })}
                    placeholder="选择或输入细分类别"
                    allowCustom
                    groupOptions={genreCategoryOptions}
                    groupValue={selectedGenreCategory}
                    groupPlaceholder="体裁大类"
                    groupSelectedPlaceholder="请选择细分类别"
                    groupRequiredMessage="请先选择体裁大类"
                    onGroupChange={(val) => setFormData({ ...formData, genreCategory: String(val), genre: '' })}
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">作者</label>
                  <Combobox options={poetOptions} value={formData.authorId}
                    onChange={(val, option) => { if (option) setFormData({ ...formData, authorId: String(val), authorName: option.label }); else setFormData({ ...formData, authorId: undefined, authorName: String(val || '') }) }}
                    placeholder={editingPoemId ? '选择作者' : '选择或输入作者'} allowCustom={!editingPoemId}
                    isLoading={isPoetsFetching} virtualListHeight={240}
                    groupOptions={authorDynastyOptions} groupValue={formData.dynastyId}
                    groupPlaceholder="先选择朝代"
                    groupSelectedPlaceholder="请选择朝代"
                    groupRequiredMessage="请先选择朝代"
                    onGroupChange={(val, option) => {
                      setFormData({
                        ...formData,
                        dynastyId: String(val),
                        dynastyName: option.label,
                        authorId: undefined,
                        authorName: '',
                      })
                      setPoetSearchKeyword('')
                      setDebouncedPoetSearchKeyword('')
                    }}
                    onSearchChange={setPoetSearchKeyword} />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">朝代</label>
                  <Combobox options={dynastyOptions} value={formData.dynastyId}
                    onChange={(val, option) => {
                      if (option) {
                        const nextDynastyId = String(val)
                        setFormData({
                          ...formData,
                          dynastyId: nextDynastyId,
                          dynastyName: option.label,
                          ...(nextDynastyId !== formData.dynastyId ? { authorId: undefined, authorName: '' } : {}),
                        })
                      } else {
                        setFormData({ ...formData, dynastyId: undefined, dynastyName: String(val || ''), authorId: undefined, authorName: '' })
                      }
                      setPoetSearchKeyword('')
                      setDebouncedPoetSearchKeyword('')
                    }}
                    placeholder={editingPoemId ? '选择朝代' : '选择或输入朝代'} allowCustom={!editingPoemId} />
                </div>
              </div>
              <PoemContentAnnotationEditor
                contentValue={formData.content}
                annotations={annotationDrafts}
                onContentChange={(content) => setFormData({ ...formData, content })}
                onAnnotationsChange={setAnnotationDrafts}
                isLoadingAnnotations={Boolean(editingPoemId && isEditingAnnotationsLoading && annotationDraftSourceId !== editingPoemId)}
                isConvertingContent={convertingTarget === 'content'}
                onConvertContent={handleConvertPoemContent}
              />
              <Tabs defaultValue="translation" className="space-y-3">
                <TabsList className="grid w-full grid-cols-5">
                  {poemFormTabs.map((field) => (
                    <TabsTrigger key={field.key} value={field.key}>
                      {field.label}
                      {field.key === 'poemAnnotations' && annotationDrafts.length > 0 ? `(${annotationDrafts.length})` : ''}
                    </TabsTrigger>
                  ))}
                </TabsList>
                {poemTextFields.map((field) => (
                  <TabsContent key={field.key} value={field.key} className="space-y-2">
                    <div className="flex items-center justify-between gap-3">
                      <label className="text-sm font-medium">{field.label}</label>
                      <ChineseVariantToggle
                        disabled={!formData[field.key].trim()}
                        isLoading={convertingTarget === field.key}
                        onChange={(_, mode) => handleConvertPoemTextField(field.key, mode)}
                      />
                    </div>
                    <textarea
                      value={formData[field.key]}
                      onChange={(e) => setFormData({ ...formData, [field.key]: e.target.value })}
                      placeholder={field.placeholder}
                      rows={6}
                      className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                    />
                  </TabsContent>
                ))}
                <TabsContent value="pingze" className="space-y-2">
                  <div className="flex items-center justify-between gap-3">
                    <label htmlFor="poem-pingze" className="text-sm font-medium">平仄</label>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      className="gap-1"
                      disabled={!formData.content.trim() || recognizePingzeMutation.isPending}
                      onClick={handleRecognizePoemPingze}
                    >
                      <Wand2 className="h-3.5 w-3.5" />
                      {recognizePingzeMutation.isPending ? '识别中...' : '自动识别'}
                    </Button>
                  </div>
                  <textarea
                    id="poem-pingze"
                    value={formData.pingze}
                    onChange={(e) => setFormData({ ...formData, pingze: e.target.value })}
                    placeholder="平平仄仄平"
                    rows={poemContentLineCount}
                    aria-invalid={invalidPingzeLine !== null}
                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  />
                </TabsContent>
                <TabsContent value="poemAnnotations" className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <label className="text-sm font-medium">诗词标注</label>
                      <ChineseVariantToggle
                        disabled={annotationDrafts.length === 0}
                        isLoading={convertingTarget === 'poemAnnotations'}
                        onChange={(_, mode) => handleConvertPoemAnnotations(mode)}
                      />
                    </div>
                    <span className="text-xs text-muted-foreground">共 {annotationDrafts.length} 条</span>
                  </div>
                  <PoemAnnotationDraftList
                    annotations={annotationDrafts}
                    onDelete={(index) => {
                      setAnnotationDrafts((drafts) => drafts.filter((_, draftIndex) => draftIndex !== index))
                    }}
                  />
                </TabsContent>
              </Tabs>
              <div className="flex justify-end gap-2 pt-2">
                <Button type="button" className="bg-secondary text-secondary-foreground hover:bg-secondary/90" onClick={closePoemDialog}>取消</Button>
                <Button type="submit" disabled={createPoemMutation.isPending || updatePoemMutation.isPending}>
                  {(createPoemMutation.isPending || updatePoemMutation.isPending) ? '提交中...' : editingPoemId ? '保存' : '添加'}
                </Button>
              </div>
            </form>
          </div>
        </div>,
        document.body,
      )}
      <PoemAnnotationManager
        poem={annotatingPoem}
        open={Boolean(annotatingPoem)}
        onClose={() => setAnnotatingPoem(null)}
      />
      <ConfirmDialog />
    </>
  )
}
