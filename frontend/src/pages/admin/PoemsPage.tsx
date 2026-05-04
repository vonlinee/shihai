import { useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Combobox, type ComboboxOption } from '@/components/ui/combobox'
import { Search, Plus, Edit2, Trash2, Eye, X, BookOpen, Crown, User } from 'lucide-react'
import { usePoems, useDynasties, usePoets, useGenres } from '@/hooks/usePoems'
import {
  useAdminDeletePoem, useAdminCreatePoem,
  useAdminCreateDynasty, useAdminUpdateDynasty, useAdminDeleteDynasty,
  useAdminCreatePoet, useAdminUpdatePoet, useAdminDeletePoet,
} from '@/hooks/useAdmin'
import { useNavigate } from 'react-router-dom'

type TabKey = 'poems' | 'dynasties' | 'poets'

const tabs: { key: TabKey; label: string; icon: typeof BookOpen }[] = [
  { key: 'poems', label: '诗词', icon: BookOpen },
  { key: 'dynasties', label: '朝代', icon: Crown },
  { key: 'poets', label: '诗人', icon: User },
]

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
  const [showAddDialog, setShowAddDialog] = useState(false)
  const [formData, setFormData] = useState<{
    title: string; content: string
    authorId?: number; authorName?: string
    dynastyId?: number; dynastyName?: string
    genre: string | number | undefined
    translation: string; appreciation: string; annotation: string
  }>({
    title: '', content: '', genre: '', translation: '', appreciation: '', annotation: '',
  })

  const { data: poemData, isLoading } = usePoems({ keyword: searchQuery || undefined, page, pageSize: 10 })
  const { data: dynasties } = useDynasties()
  const { data: poets } = usePoets()
  const { data: genres } = useGenres()

  const dynastyOptions: ComboboxOption[] = (dynasties ?? []).map((d) => ({ value: d.id, label: d.name, description: d.period }))
  const poetOptions: ComboboxOption[] = (poets ?? []).map((p) => ({ value: p.id, label: p.name, description: p.dynasty?.name }))
  const genreOptions: ComboboxOption[] = (genres ?? []).map((g) => ({ value: g, label: String(g) }))

  const deletePoemMutation = useAdminDeletePoem()
  const createPoemMutation = useAdminCreatePoem()
  const poems = poemData?.list ?? []
  const total = poemData?.total ?? 0

  const handleDelete = (id: number) => { if (confirm('确定要删除该诗词吗？')) deletePoemMutation.mutate(id) }

  const handleAddPoem = (e: React.FormEvent) => {
    e.preventDefault()
    const payload: Record<string, unknown> = {
      title: formData.title, content: formData.content,
      genre: typeof formData.genre === 'number' ? String(formData.genre) : formData.genre || '',
      translation: formData.translation, appreciation: formData.appreciation, annotation: formData.annotation,
    }
    if (formData.dynastyId) payload.dynastyId = formData.dynastyId
    else if (formData.dynastyName) payload.dynastyName = formData.dynastyName
    if (formData.authorId) payload.authorId = formData.authorId
    else if (formData.authorName) payload.authorName = formData.authorName
    createPoemMutation.mutate(payload as any, {
      onSuccess: () => {
        setShowAddDialog(false)
        setFormData({ title: '', content: '', genre: '', translation: '', appreciation: '', annotation: '' })
      },
    })
  }

  return (
    <>
      <div className="flex items-center justify-between">
        <div className="relative max-w-sm flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input placeholder="搜索诗词标题或作者..." className="pl-10" value={searchQuery}
            onChange={(e) => { setSearchQuery(e.target.value); setPage(1) }} />
        </div>
        <Button onClick={() => setShowAddDialog(true)}><Plus className="h-4 w-4 mr-2" />添加诗词</Button>
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
                    <th className="text-left py-3 px-4 font-medium">ID</th>
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
                      <td className="py-3 px-4">{poem.id}</td>
                      <td className="py-3 px-4 font-medium">{poem.title}</td>
                      <td className="py-3 px-4">{poem.author?.name}</td>
                      <td className="py-3 px-4">{poem.dynasty?.name}</td>
                      <td className="py-3 px-4">{poem.genre}</td>
                      <td className="py-3 px-4">{poem.views}</td>
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-2">
                          <Button variant="ghost" size="sm" onClick={() => navigate(`/poems/${poem.id}`)}><Eye className="h-4 w-4" /></Button>
                          <Button variant="ghost" size="sm"><Edit2 className="h-4 w-4" /></Button>
                          <Button variant="ghost" size="sm" onClick={() => handleDelete(poem.id)}><Trash2 className="h-4 w-4 text-cinnabar" /></Button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          {total > 10 && (
            <div className="flex justify-center gap-2 mt-4">
              <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage(page - 1)}>上一页</Button>
              <span className="flex items-center px-4 text-sm text-muted-foreground">第 {page} 页</span>
              <Button variant="outline" size="sm" disabled={page >= Math.ceil(total / 10)} onClick={() => setPage(page + 1)}>下一页</Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Add Poem Dialog */}
      {showAddDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg w-full max-w-2xl max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">添加诗词</h2>
              <button onClick={() => setShowAddDialog(false)} className="p-1 hover:bg-muted rounded"><X className="h-5 w-5" /></button>
            </div>
            <form onSubmit={handleAddPoem} className="space-y-4 p-6 pt-4 overflow-y-auto">
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
                    onChange={(val) => { if (typeof val === 'number') setFormData({ ...formData, authorId: val, authorName: undefined }); else setFormData({ ...formData, authorId: undefined, authorName: String(val) }) }}
                    placeholder="选择或输入作者" allowCustom />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">朝代</label>
                  <Combobox options={dynastyOptions} value={formData.dynastyId}
                    onChange={(val) => { if (typeof val === 'number') setFormData({ ...formData, dynastyId: val, dynastyName: undefined }); else setFormData({ ...formData, dynastyId: undefined, dynastyName: String(val) }) }}
                    placeholder="选择或输入朝代" allowCustom />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">内容 *</label>
                <textarea value={formData.content} onChange={(e) => setFormData({ ...formData, content: e.target.value })} placeholder="诗词内容" required rows={4}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
              </div>
              {['translation', 'appreciation', 'annotation'].map((field) => (
                <div key={field} className="space-y-2">
                  <label className="text-sm font-medium">{{ translation: '译文', appreciation: '赏析', annotation: '注释' }[field]}</label>
                  <textarea value={(formData as any)[field]} onChange={(e) => setFormData({ ...formData, [field]: e.target.value })}
                    placeholder={{ translation: '白话文翻译', appreciation: '诗词赏析', annotation: '字词注释' }[field]} rows={3}
                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
                </div>
              ))}
              <div className="flex justify-end gap-2 pt-2">
                <Button type="button" className="bg-secondary text-secondary-foreground hover:bg-secondary/90" onClick={() => setShowAddDialog(false)}>取消</Button>
                <Button type="submit" disabled={createPoemMutation.isPending}>{createPoemMutation.isPending ? '提交中...' : '添加'}</Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  )
}

// ─── Dynasties Tab ───────────────────────────────────────────────────────────

function DynastiesTab() {
  const { data: dynasties, isLoading } = useDynasties()
  const createMutation = useAdminCreateDynasty()
  const updateMutation = useAdminUpdateDynasty()
  const deleteMutation = useAdminDeleteDynasty()
  const [showDialog, setShowDialog] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [formData, setFormData] = useState({ name: '', period: '', description: '' })

  const openCreate = () => { setEditingId(null); setFormData({ name: '', period: '', description: '' }); setShowDialog(true) }
  const openEdit = (d: { id: number; name: string; period?: string; description?: string }) => {
    setEditingId(d.id); setFormData({ name: d.name, period: d.period || '', description: d.description || '' }); setShowDialog(true)
  }
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (editingId) updateMutation.mutate({ id: editingId, data: formData }, { onSuccess: () => setShowDialog(false) })
    else createMutation.mutate(formData, { onSuccess: () => setShowDialog(false) })
  }
  const handleDelete = (id: number) => { if (confirm('确定要删除该朝代吗？')) deleteMutation.mutate(id) }

  return (
    <>
      <div className="flex items-center justify-end">
        <Button onClick={openCreate}><Plus className="h-4 w-4 mr-2" />添加朝代</Button>
      </div>
      <Card className="ink-border">
        <CardContent className="pt-6">
          {isLoading ? <div className="text-center py-8 text-muted-foreground">加载中...</div>
            : (dynasties?.length ?? 0) === 0 ? <div className="text-center py-8 text-muted-foreground">暂无朝代数据</div>
            : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead><tr className="border-b">
                    <th className="text-left py-3 px-4 font-medium">ID</th>
                    <th className="text-left py-3 px-4 font-medium">朝代名称</th>
                    <th className="text-left py-3 px-4 font-medium">时期</th>
                    <th className="text-left py-3 px-4 font-medium">描述</th>
                    <th className="text-left py-3 px-4 font-medium">操作</th>
                  </tr></thead>
                  <tbody>
                    {dynasties!.map((d) => (
                      <tr key={d.id} className="border-b last:border-0 hover:bg-muted/50">
                        <td className="py-3 px-4">{d.id}</td>
                        <td className="py-3 px-4 font-medium">{d.name}</td>
                        <td className="py-3 px-4 text-muted-foreground">{d.period || '-'}</td>
                        <td className="py-3 px-4 text-muted-foreground max-w-xs truncate">{d.description || '-'}</td>
                        <td className="py-3 px-4">
                          <div className="flex items-center gap-2">
                            <Button variant="ghost" size="sm" onClick={() => openEdit(d)}><Edit2 className="h-4 w-4" /></Button>
                            <Button variant="ghost" size="sm" onClick={() => handleDelete(d.id)}><Trash2 className="h-4 w-4 text-cinnabar" /></Button>
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
    </>
  )
}

// ─── Poets Tab ───────────────────────────────────────────────────────────────

function PoetsTab() {
  const { data: poets, isLoading } = usePoets()
  const { data: dynasties } = useDynasties()
  const createMutation = useAdminCreatePoet()
  const updateMutation = useAdminUpdatePoet()
  const deleteMutation = useAdminDeletePoet()
  const [showDialog, setShowDialog] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [formData, setFormData] = useState<{ name: string; dynastyId: number | undefined; biography: string; avatar: string; birthYear: string; deathYear: string }>({
    name: '', dynastyId: undefined, biography: '', avatar: '', birthYear: '', deathYear: '',
  })

  const dynastyOptions: ComboboxOption[] = (dynasties ?? []).map((d) => ({ value: d.id, label: d.name, description: d.period }))

  const openCreate = () => { setEditingId(null); setFormData({ name: '', dynastyId: undefined, biography: '', avatar: '', birthYear: '', deathYear: '' }); setShowDialog(true) }
  const openEdit = (p: { id: number; name: string; dynastyId?: number; biography?: string; avatar?: string; birthYear?: number; deathYear?: number }) => {
    setEditingId(p.id); setFormData({ name: p.name, dynastyId: p.dynastyId, biography: p.biography || '', avatar: p.avatar || '', birthYear: p.birthYear?.toString() || '', deathYear: p.deathYear?.toString() || '' }); setShowDialog(true)
  }
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const payload = { name: formData.name, dynastyId: formData.dynastyId, biography: formData.biography, avatar: formData.avatar, birthYear: formData.birthYear ? parseInt(formData.birthYear) : undefined, deathYear: formData.deathYear ? parseInt(formData.deathYear) : undefined }
    if (editingId) updateMutation.mutate({ id: editingId, data: payload }, { onSuccess: () => setShowDialog(false) })
    else createMutation.mutate(payload, { onSuccess: () => setShowDialog(false) })
  }
  const handleDelete = (id: number) => { if (confirm('确定要删除该诗人吗？')) deleteMutation.mutate(id) }

  return (
    <>
      <div className="flex items-center justify-end">
        <Button onClick={openCreate}><Plus className="h-4 w-4 mr-2" />添加诗人</Button>
      </div>
      <Card className="ink-border">
        <CardContent className="pt-6">
          {isLoading ? <div className="text-center py-8 text-muted-foreground">加载中...</div>
            : (poets?.length ?? 0) === 0 ? <div className="text-center py-8 text-muted-foreground">暂无诗人数据</div>
            : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead><tr className="border-b">
                    <th className="text-left py-3 px-4 font-medium">ID</th>
                    <th className="text-left py-3 px-4 font-medium">姓名</th>
                    <th className="text-left py-3 px-4 font-medium">朝代</th>
                    <th className="text-left py-3 px-4 font-medium">生卒年</th>
                    <th className="text-left py-3 px-4 font-medium">简介</th>
                    <th className="text-left py-3 px-4 font-medium">操作</th>
                  </tr></thead>
                  <tbody>
                    {poets!.map((p) => (
                      <tr key={p.id} className="border-b last:border-0 hover:bg-muted/50">
                        <td className="py-3 px-4">{p.id}</td>
                        <td className="py-3 px-4 font-medium">{p.name}</td>
                        <td className="py-3 px-4 text-muted-foreground">{p.dynasty?.name || '-'}</td>
                        <td className="py-3 px-4 text-muted-foreground">{p.birthYear && p.deathYear ? `${p.birthYear}-${p.deathYear}` : '-'}</td>
                        <td className="py-3 px-4 text-muted-foreground max-w-xs truncate">{p.biography || '-'}</td>
                        <td className="py-3 px-4">
                          <div className="flex items-center gap-2">
                            <Button variant="ghost" size="sm" onClick={() => openEdit(p)}><Edit2 className="h-4 w-4" /></Button>
                            <Button variant="ghost" size="sm" onClick={() => handleDelete(p.id)}><Trash2 className="h-4 w-4 text-cinnabar" /></Button>
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
                    onChange={(val) => setFormData({ ...formData, dynastyId: typeof val === 'number' ? val : undefined })}
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
    </>
  )
}
