import { useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Combobox, type ComboboxOption } from '@/components/ui/combobox'
import { Pagination } from '@/components/ui/Pagination'
import { useConfirmDialog } from '@/components/ui/ConfirmDialog'
import { Search, Plus, Edit2, Trash2, Eye, X, BookOpen, Crown, User } from 'lucide-react'
import { usePoems, useDynasties, usePoets, usePoetList, useGenres } from '@/hooks/usePoems'
import {
  useAdminDeletePoem, useAdminCreatePoem, useAdminUpdatePoem,
  useAdminBatchDeletePoems,
  useAdminCreateDynasty, useAdminUpdateDynasty, useAdminDeleteDynasty,
  useAdminBatchDeleteDynasties,
  useAdminCreatePoet, useAdminUpdatePoet, useAdminDeletePoet,
  useAdminBatchDeletePoets,
} from '@/hooks/useAdmin'
import { useNavigate } from 'react-router-dom'
import { splitPoemContentInput } from '@/utils/poemContent'
import type { Poem } from '@/types'
import type { PoemCreateRequest, PoemUpdateRequest } from '@/services/adminService'

type TabKey = 'poems' | 'dynasties' | 'poets'
type PoemFormData = {
  title: string
  content: string
  authorId?: string
  authorName?: string
  dynastyId?: string
  dynastyName?: string
  genre: string | number | undefined
  translation: string
  appreciation: string
  annotation: string
}
type PoemTextField = 'translation' | 'appreciation' | 'annotation'

const tabs: { key: TabKey; label: string; icon: typeof BookOpen }[] = [
  { key: 'poems', label: '诗词', icon: BookOpen },
  { key: 'dynasties', label: '朝代', icon: Crown },
  { key: 'poets', label: '诗人', icon: User },
]

const emptyPoemFormData: PoemFormData = {
  title: '',
  content: '',
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

function createEmptyPoemFormData(): PoemFormData {
  return { ...emptyPoemFormData }
}

function formatPoemContentForInput(content: Poem['content']): string {
  return content.join('\n')
}

function toPoemUpdateId(id: string | undefined): string | undefined {
  return id && id.trim() ? id : undefined
}

function toEntityId(id: string | number): string {
  return String(id)
}

function toggleSelectedId(selectedIds: string[], id: string): string[] {
  return selectedIds.includes(id)
    ? selectedIds.filter((selectedId) => selectedId !== id)
    : [...selectedIds, id]
}

function toggleAllVisibleIds(selectedIds: string[], visibleIds: string[]): string[] {
  const allVisibleSelected = visibleIds.length > 0 && visibleIds.every((id) => selectedIds.includes(id))
  if (allVisibleSelected) {
    return selectedIds.filter((id) => !visibleIds.includes(id))
  }
  return Array.from(new Set([...selectedIds, ...visibleIds]))
}

// ─── Main Component ──────────────────────────────────────────────────────────

export function AdminPoemsPage() {
  const [activeTab, setActiveTab] = useState<TabKey>('poems')

  return (
    <div className="p-8 space-y-6">
      {/* Tab Bar */}
      <div className="flex border-b gap-0">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`flex items-center gap-2 px-5 py-2.5 text-sm font-medium border-b-2 transition-colors -mb-px ${
              activeTab === tab.key
                ? 'border-primary text-primary'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <tab.icon className="h-4 w-4" />
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'poems' && <PoemsTab />}
      {activeTab === 'dynasties' && <DynastiesTab />}
      {activeTab === 'poets' && <PoetsTab />}
    </div>
  )
}

// ─── Poems Tab ───────────────────────────────────────────────────────────────

function PoemsTab() {
  const navigate = useNavigate()
  const [searchQuery, setSearchQuery] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [showPoemDialog, setShowPoemDialog] = useState(false)
  const [editingPoemId, setEditingPoemId] = useState<string | null>(null)
  const [formData, setFormData] = useState<PoemFormData>(() => createEmptyPoemFormData())
  const [selectedPoemIds, setSelectedPoemIds] = useState<string[]>([])

  const { data: poemData, isLoading } = usePoems({ keyword: searchQuery || undefined, page, pageSize })
  const { data: dynasties } = useDynasties()
  const { data: poets } = usePoets()
  const { data: genres } = useGenres()

  const dynastyOptions: ComboboxOption[] = (dynasties ?? []).map((d) => ({ value: String(d.id), label: d.name, description: d.period }))
  const poetOptions: ComboboxOption[] = (poets ?? []).map((p) => ({ value: String(p.id), label: p.name, description: p.dynasty?.name }))
  const genreOptions: ComboboxOption[] = (genres ?? []).map((g) => ({ value: g, label: String(g) }))

  const deletePoemMutation = useAdminDeletePoem()
  const batchDeletePoemsMutation = useAdminBatchDeletePoems()
  const createPoemMutation = useAdminCreatePoem()
  const updatePoemMutation = useAdminUpdatePoem()
  const { confirm, ConfirmDialog } = useConfirmDialog()
  const poems = poemData?.list ?? []
  const total = poemData?.total ?? 0
  const visiblePoemIds = poems.map((poem) => toEntityId(poem.id))
  const allVisiblePoemsSelected = visiblePoemIds.length > 0 && visiblePoemIds.every((id) => selectedPoemIds.includes(id))

  const handleDelete = async (id: string) => {
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
  }

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
    setShowPoemDialog(true)
  }

  const openEditPoemDialog = (poem: Poem) => {
    setEditingPoemId(toEntityId(poem.id))
    setFormData({
      title: poem.title,
      content: formatPoemContentForInput(poem.content),
      authorId: poem.authorId ? String(poem.authorId) : undefined,
      authorName: poem.author?.name,
      dynastyId: poem.dynastyId ? String(poem.dynastyId) : undefined,
      dynastyName: poem.dynasty?.name,
      genre: poem.genre || '',
      translation: poem.translation || '',
      appreciation: poem.appreciation || '',
      annotation: poem.annotation || '',
    })
    setShowPoemDialog(true)
  }

  const resetPoemDialog = () => {
    setEditingPoemId(null)
    setFormData(createEmptyPoemFormData())
    setShowPoemDialog(false)
  }

  const buildCreatePayload = (): PoemCreateRequest => {
    const payload: PoemCreateRequest = {
      title: formData.title,
      content: splitPoemContentInput(formData.content),
      genre: typeof formData.genre === 'number' ? String(formData.genre) : formData.genre || '',
      translation: formData.translation,
      appreciation: formData.appreciation,
      annotation: formData.annotation,
    }
    if (formData.dynastyName) payload.dynastyName = formData.dynastyName
    if (formData.authorName) payload.authorName = formData.authorName
    return payload
  }

  const buildUpdatePayload = (): PoemUpdateRequest => {
    const payload: PoemUpdateRequest = {
      title: formData.title,
      content: splitPoemContentInput(formData.content),
      genre: typeof formData.genre === 'number' ? String(formData.genre) : formData.genre || '',
      translation: formData.translation,
      appreciation: formData.appreciation,
      annotation: formData.annotation,
    }
    const dynastyId = toPoemUpdateId(formData.dynastyId)
    const authorId = toPoemUpdateId(formData.authorId)
    if (dynastyId) payload.dynastyId = dynastyId
    if (authorId) payload.authorId = authorId
    return payload
  }

  const handlePoemSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (editingPoemId) {
      updatePoemMutation.mutate(
        { id: editingPoemId, data: buildUpdatePayload() },
        { onSuccess: resetPoemDialog },
      )
      return
    }
    createPoemMutation.mutate(buildCreatePayload(), { onSuccess: resetPoemDialog })
  }

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
          {isLoading ? (
            <div className="text-center py-8 text-muted-foreground">加载中...</div>
          ) : poems.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">暂无诗词数据</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="w-12 py-3 px-4">
                      <input
                        type="checkbox"
                        aria-label="选择当前页诗词"
                        checked={allVisiblePoemsSelected}
                        onChange={() => setSelectedPoemIds((ids) => toggleAllVisibleIds(ids, visiblePoemIds))}
                      />
                    </th>
                    <th className="text-left py-3 px-4 font-medium">标题</th>
                    <th className="text-left py-3 px-4 font-medium">作者</th>
                    <th className="text-left py-3 px-4 font-medium">朝代</th>
                    <th className="text-left py-3 px-4 font-medium">体裁</th>
                    <th className="text-left py-3 px-4 font-medium">浏览量</th>
                    <th className="text-left py-3 px-4 font-medium">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {poems.map((poem) => (
                    <tr key={poem.id} className="border-b last:border-0 hover:bg-muted/50">
                      <td className="py-3 px-4">
                        <input
                          type="checkbox"
                          aria-label={`选择诗词 ${poem.title}`}
                          checked={selectedPoemIds.includes(toEntityId(poem.id))}
                          onChange={() => setSelectedPoemIds((ids) => toggleSelectedId(ids, toEntityId(poem.id)))}
                        />
                      </td>
                      <td className="py-3 px-4 font-medium">{poem.title}</td>
                      <td className="py-3 px-4">{poem.author?.name}</td>
                      <td className="py-3 px-4">{poem.dynasty?.name}</td>
                      <td className="py-3 px-4">{poem.genre}</td>
                      <td className="py-3 px-4">{poem.views}</td>
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-2">
                          <Button variant="ghost" size="sm" onClick={() => navigate(`/poems/${poem.id}`)}><Eye className="h-4 w-4" /></Button>
                          <Button variant="ghost" size="sm" onClick={() => openEditPoemDialog(poem)}><Edit2 className="h-4 w-4" /></Button>
                          <Button variant="ghost" size="sm" onClick={() => handleDelete(toEntityId(poem.id))}><Trash2 className="h-4 w-4 text-cinnabar" /></Button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          <Pagination
            className="mt-4"
            page={page}
            pageSize={pageSize}
            total={total}
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

      {/* Poem Dialog */}
      {showPoemDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg w-full max-w-2xl max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">{editingPoemId ? '编辑诗词' : '添加诗词'}</h2>
              <button onClick={closePoemDialog} className="p-1 hover:bg-muted rounded"><X className="h-5 w-5" /></button>
            </div>
            <form onSubmit={handlePoemSubmit} className="space-y-4 p-6 pt-4 overflow-y-auto">
              <div className="space-y-2">
                <label className="text-sm font-medium">标题 *</label>
                <Input value={formData.title} onChange={(e) => setFormData({ ...formData, title: e.target.value })} placeholder="诗词标题" required />
              </div>
              <div className="grid grid-cols-3 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">体裁</label>
                  <Combobox options={genreOptions} value={formData.genre} onChange={(val) => setFormData({ ...formData, genre: val })} placeholder="选择或输入体裁" allowCustom />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">作者</label>
                  <Combobox options={poetOptions} value={formData.authorId}
                    onChange={(val, option) => { if (option) setFormData({ ...formData, authorId: String(val), authorName: option.label }); else setFormData({ ...formData, authorId: undefined, authorName: String(val || '') }) }}
                    placeholder={editingPoemId ? '选择作者' : '选择或输入作者'} allowCustom={!editingPoemId} />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">朝代</label>
                  <Combobox options={dynastyOptions} value={formData.dynastyId}
                    onChange={(val, option) => { if (option) setFormData({ ...formData, dynastyId: String(val), dynastyName: option.label }); else setFormData({ ...formData, dynastyId: undefined, dynastyName: String(val || '') }) }}
                    placeholder={editingPoemId ? '选择朝代' : '选择或输入朝代'} allowCustom={!editingPoemId} />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">内容 *</label>
                <textarea value={formData.content} onChange={(e) => setFormData({ ...formData, content: e.target.value })} placeholder="诗词内容" required rows={4}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
              </div>
              {poemTextFields.map((field) => (
                <div key={field.key} className="space-y-2">
                  <label className="text-sm font-medium">{field.label}</label>
                  <textarea value={formData[field.key]} onChange={(e) => setFormData({ ...formData, [field.key]: e.target.value })}
                    placeholder={field.placeholder} rows={3}
                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
                </div>
              ))}
              <div className="flex justify-end gap-2 pt-2">
                <Button type="button" className="bg-secondary text-secondary-foreground hover:bg-secondary/90" onClick={closePoemDialog}>取消</Button>
                <Button type="submit" disabled={createPoemMutation.isPending || updatePoemMutation.isPending}>
                  {(createPoemMutation.isPending || updatePoemMutation.isPending) ? '提交中...' : editingPoemId ? '保存' : '添加'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
      <ConfirmDialog />
    </>
  )
}

// ─── Dynasties Tab ───────────────────────────────────────────────────────────

function DynastiesTab() {
  const { data: dynasties, isLoading } = useDynasties()
  const createMutation = useAdminCreateDynasty()
  const updateMutation = useAdminUpdateDynasty()
  const deleteMutation = useAdminDeleteDynasty()
  const batchDeleteMutation = useAdminBatchDeleteDynasties()
  const { confirm, ConfirmDialog } = useConfirmDialog()
  const [showDialog, setShowDialog] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [formData, setFormData] = useState({ name: '', period: '', description: '' })
  const [selectedDynastyIds, setSelectedDynastyIds] = useState<string[]>([])
  const dynastiesList = dynasties ?? []
  const visibleDynastyIds = dynastiesList.map((d) => toEntityId(d.id))
  const allVisibleDynastiesSelected = visibleDynastyIds.length > 0 && visibleDynastyIds.every((id) => selectedDynastyIds.includes(id))

  const openCreate = () => { setEditingId(null); setFormData({ name: '', period: '', description: '' }); setShowDialog(true) }
  const openEdit = (d: { id: string | number; name: string; period?: string; description?: string }) => {
    setEditingId(toEntityId(d.id)); setFormData({ name: d.name, period: d.period || '', description: d.description || '' }); setShowDialog(true)
  }
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (editingId) updateMutation.mutate({ id: editingId, data: formData }, { onSuccess: () => setShowDialog(false) })
    else createMutation.mutate(formData, { onSuccess: () => setShowDialog(false) })
  }
  const handleDelete = async (id: string) => {
    const confirmed = await confirm({
      title: '删除朝代',
      description: '确定要删除该朝代吗？此操作不可撤销。',
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    deleteMutation.mutate(id, {
      onSuccess: () => setSelectedDynastyIds((ids) => ids.filter((selectedId) => selectedId !== id)),
    })
  }

  const handleBatchDelete = async () => {
    if (selectedDynastyIds.length === 0) return
    const confirmed = await confirm({
      title: '批量删除朝代',
      description: `确定要删除选中的 ${selectedDynastyIds.length} 个朝代吗？此操作不可撤销。`,
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    batchDeleteMutation.mutate(selectedDynastyIds, { onSuccess: () => setSelectedDynastyIds([]) })
  }

  return (
    <>
      <div className="flex items-center justify-end gap-2">
        {selectedDynastyIds.length > 0 && (
          <Button
            variant="destructive"
            onClick={handleBatchDelete}
            disabled={batchDeleteMutation.isPending}
          >
            <Trash2 className="h-4 w-4 mr-2" />
            删除选中({selectedDynastyIds.length})
          </Button>
        )}
        <Button onClick={openCreate}><Plus className="h-4 w-4 mr-2" />添加朝代</Button>
      </div>
      <Card className="ink-border">
        <CardContent className="pt-6">
          {isLoading ? <div className="text-center py-8 text-muted-foreground">加载中...</div>
            : dynastiesList.length === 0 ? <div className="text-center py-8 text-muted-foreground">暂无朝代数据</div>
            : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead><tr className="border-b">
                    <th className="w-12 py-3 px-4">
                      <input
                        type="checkbox"
                        aria-label="选择全部朝代"
                        checked={allVisibleDynastiesSelected}
                        onChange={() => setSelectedDynastyIds((ids) => toggleAllVisibleIds(ids, visibleDynastyIds))}
                      />
                    </th>
                    <th className="text-left py-3 px-4 font-medium">朝代名称</th>
                    <th className="text-left py-3 px-4 font-medium">时期</th>
                    <th className="text-left py-3 px-4 font-medium">描述</th>
                    <th className="text-left py-3 px-4 font-medium">操作</th>
                  </tr></thead>
                  <tbody>
                    {dynastiesList.map((d) => (
                      <tr key={d.id} className="border-b last:border-0 hover:bg-muted/50">
                        <td className="py-3 px-4">
                          <input
                            type="checkbox"
                            aria-label={`选择朝代 ${d.name}`}
                            checked={selectedDynastyIds.includes(toEntityId(d.id))}
                            onChange={() => setSelectedDynastyIds((ids) => toggleSelectedId(ids, toEntityId(d.id)))}
                          />
                        </td>
                        <td className="py-3 px-4 font-medium">{d.name}</td>
                        <td className="py-3 px-4 text-muted-foreground">{d.period || '-'}</td>
                        <td className="py-3 px-4 text-muted-foreground max-w-xs truncate">{d.description || '-'}</td>
                        <td className="py-3 px-4">
                          <div className="flex items-center gap-2">
                            <Button variant="ghost" size="sm" onClick={() => openEdit(d)}><Edit2 className="h-4 w-4" /></Button>
                              <Button variant="ghost" size="sm" onClick={() => handleDelete(toEntityId(d.id))}><Trash2 className="h-4 w-4 text-cinnabar" /></Button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
        </CardContent>
      </Card>

      {showDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg w-full max-w-md max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">{editingId ? '编辑朝代' : '添加朝代'}</h2>
              <button onClick={() => setShowDialog(false)} className="p-1 hover:bg-muted rounded"><X className="h-5 w-5" /></button>
            </div>
            <form onSubmit={handleSubmit} className="space-y-4 p-6 pt-4 overflow-y-auto">
              <div className="space-y-2">
                <label className="text-sm font-medium">朝代名称 *</label>
                <Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} placeholder="如：唐" required />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">时期</label>
                <Input value={formData.period} onChange={(e) => setFormData({ ...formData, period: e.target.value })} placeholder="如：618-907" />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">描述</label>
                <textarea value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} placeholder="朝代简介" rows={3}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button type="button" className="bg-secondary text-secondary-foreground hover:bg-secondary/90" onClick={() => setShowDialog(false)}>取消</Button>
                <Button type="submit" disabled={createMutation.isPending || updateMutation.isPending}>
                  {(createMutation.isPending || updateMutation.isPending) ? '提交中...' : editingId ? '保存' : '添加'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
      <ConfirmDialog />
    </>
  )
}

// ─── Poets Tab ───────────────────────────────────────────────────────────────

function PoetsTab() {
  const [searchQuery, setSearchQuery] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const { data: poetData, isLoading } = usePoetList({ keyword: searchQuery || undefined, page, pageSize })
  const { data: dynasties } = useDynasties()
  const createMutation = useAdminCreatePoet()
  const updateMutation = useAdminUpdatePoet()
  const deleteMutation = useAdminDeletePoet()
  const batchDeleteMutation = useAdminBatchDeletePoets()
  const { confirm, ConfirmDialog } = useConfirmDialog()
  const [showDialog, setShowDialog] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [formData, setFormData] = useState<{ name: string; dynastyId: string | undefined; biography: string; avatar: string; birthYear: string; deathYear: string }>({
    name: '', dynastyId: undefined, biography: '', avatar: '', birthYear: '', deathYear: '',
  })
  const [selectedPoetIds, setSelectedPoetIds] = useState<string[]>([])

  const dynastyOptions: ComboboxOption[] = (dynasties ?? []).map((d) => ({ value: d.id, label: d.name, description: d.period }))
  const poets = poetData?.list ?? []
  const total = poetData?.total ?? 0
  const visiblePoetIds = poets.map((p) => toEntityId(p.id))
  const allVisiblePoetsSelected = visiblePoetIds.length > 0 && visiblePoetIds.every((id) => selectedPoetIds.includes(id))

  const openCreate = () => { setEditingId(null); setFormData({ name: '', dynastyId: undefined, biography: '', avatar: '', birthYear: '', deathYear: '' }); setShowDialog(true) }
  const openEdit = (p: { id: string | number; name: string; dynastyId?: string | number; biography?: string; avatar?: string; birthYear?: number; deathYear?: number }) => {
    setEditingId(toEntityId(p.id)); setFormData({ name: p.name, dynastyId: p.dynastyId ? toEntityId(p.dynastyId) : undefined, biography: p.biography || '', avatar: p.avatar || '', birthYear: p.birthYear?.toString() || '', deathYear: p.deathYear?.toString() || '' }); setShowDialog(true)
  }
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const payload = { name: formData.name, dynastyId: formData.dynastyId, biography: formData.biography, avatar: formData.avatar, birthYear: formData.birthYear ? parseInt(formData.birthYear) : undefined, deathYear: formData.deathYear ? parseInt(formData.deathYear) : undefined }
    if (editingId) updateMutation.mutate({ id: editingId, data: payload }, { onSuccess: () => setShowDialog(false) })
    else createMutation.mutate(payload, { onSuccess: () => setShowDialog(false) })
  }
  const handleDelete = async (id: string) => {
    const confirmed = await confirm({
      title: '删除诗人',
      description: '确定要删除该诗人吗？此操作不可撤销。',
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    deleteMutation.mutate(id, {
      onSuccess: () => setSelectedPoetIds((ids) => ids.filter((selectedId) => selectedId !== id)),
    })
  }

  const handleBatchDelete = async () => {
    if (selectedPoetIds.length === 0) return
    const confirmed = await confirm({
      title: '批量删除诗人',
      description: `确定要删除选中的 ${selectedPoetIds.length} 位诗人吗？此操作不可撤销。`,
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    batchDeleteMutation.mutate(selectedPoetIds, { onSuccess: () => setSelectedPoetIds([]) })
  }

  return (
    <>
      <div className="flex items-center justify-between">
        <div className="relative max-w-sm flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="搜索诗人姓名..."
            className="pl-10"
            value={searchQuery}
            onChange={(e) => { setSearchQuery(e.target.value); setSelectedPoetIds([]); setPage(1) }}
          />
        </div>
        <div className="flex items-center gap-2">
          {selectedPoetIds.length > 0 && (
            <Button
              variant="destructive"
              onClick={handleBatchDelete}
              disabled={batchDeleteMutation.isPending}
            >
              <Trash2 className="h-4 w-4 mr-2" />
              删除选中({selectedPoetIds.length})
            </Button>
          )}
          <Button onClick={openCreate}><Plus className="h-4 w-4 mr-2" />添加诗人</Button>
        </div>
      </div>
      <Card className="ink-border">
        <CardContent className="pt-6">
          {isLoading ? <div className="text-center py-8 text-muted-foreground">加载中...</div>
            : poets.length === 0 ? <div className="text-center py-8 text-muted-foreground">暂无诗人数据</div>
            : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead><tr className="border-b">
                    <th className="w-12 py-3 px-4">
                      <input
                        type="checkbox"
                        aria-label="选择当前页诗人"
                        checked={allVisiblePoetsSelected}
                        onChange={() => setSelectedPoetIds((ids) => toggleAllVisibleIds(ids, visiblePoetIds))}
                      />
                    </th>
                    <th className="text-left py-3 px-4 font-medium">姓名</th>
                    <th className="text-left py-3 px-4 font-medium">朝代</th>
                    <th className="text-left py-3 px-4 font-medium">生卒年</th>
                    <th className="text-left py-3 px-4 font-medium">简介</th>
                    <th className="text-left py-3 px-4 font-medium">操作</th>
                  </tr></thead>
                  <tbody>
                    {poets.map((p) => (
                      <tr key={p.id} className="border-b last:border-0 hover:bg-muted/50">
                        <td className="py-3 px-4">
                          <input
                            type="checkbox"
                            aria-label={`选择诗人 ${p.name}`}
                            checked={selectedPoetIds.includes(toEntityId(p.id))}
                            onChange={() => setSelectedPoetIds((ids) => toggleSelectedId(ids, toEntityId(p.id)))}
                          />
                        </td>
                        <td className="py-3 px-4 font-medium">{p.name}</td>
                        <td className="py-3 px-4 text-muted-foreground">{p.dynasty?.name || '-'}</td>
                        <td className="py-3 px-4 text-muted-foreground">{p.birthYear && p.deathYear ? `${p.birthYear}-${p.deathYear}` : '-'}</td>
                        <td className="py-3 px-4 text-muted-foreground max-w-xs truncate">{p.biography || '-'}</td>
                        <td className="py-3 px-4">
                          <div className="flex items-center gap-2">
                            <Button variant="ghost" size="sm" onClick={() => openEdit(p)}><Edit2 className="h-4 w-4" /></Button>
                            <Button variant="ghost" size="sm" onClick={() => handleDelete(toEntityId(p.id))}><Trash2 className="h-4 w-4 text-cinnabar" /></Button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          <Pagination
            className="mt-4"
            page={page}
            pageSize={pageSize}
            total={total}
            onPageChange={(nextPage) => {
              setSelectedPoetIds([])
              setPage(nextPage)
            }}
            onPageSizeChange={(nextPageSize) => {
              setSelectedPoetIds([])
              setPageSize(nextPageSize)
              setPage(1)
            }}
          />
        </CardContent>
      </Card>

      {showDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg w-full max-w-lg max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">{editingId ? '编辑诗人' : '添加诗人'}</h2>
              <button onClick={() => setShowDialog(false)} className="p-1 hover:bg-muted rounded"><X className="h-5 w-5" /></button>
            </div>
            <form onSubmit={handleSubmit} className="space-y-4 p-6 pt-4 overflow-y-auto">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">姓名 *</label>
                  <Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} placeholder="诗人姓名" required />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">朝代</label>
                  <Combobox options={dynastyOptions} value={formData.dynastyId}
                    onChange={(val) => setFormData({ ...formData, dynastyId: val ? String(val) : undefined })}
                    placeholder="选择朝代" allowCustom={false} />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">出生年份</label>
                  <Input type="number" value={formData.birthYear} onChange={(e) => setFormData({ ...formData, birthYear: e.target.value })} placeholder="如：701" />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">逝世年份</label>
                  <Input type="number" value={formData.deathYear} onChange={(e) => setFormData({ ...formData, deathYear: e.target.value })} placeholder="如：762" />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">头像URL</label>
                <Input value={formData.avatar} onChange={(e) => setFormData({ ...formData, avatar: e.target.value })} placeholder="诗人头像链接" />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">简介</label>
                <textarea value={formData.biography} onChange={(e) => setFormData({ ...formData, biography: e.target.value })} placeholder="诗人简介" rows={3}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button type="button" className="bg-secondary text-secondary-foreground hover:bg-secondary/90" onClick={() => setShowDialog(false)}>取消</Button>
                <Button type="submit" disabled={createMutation.isPending || updateMutation.isPending}>
                  {(createMutation.isPending || updateMutation.isPending) ? '提交中...' : editingId ? '保存' : '添加'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
      <ConfirmDialog />
    </>
  )
}
