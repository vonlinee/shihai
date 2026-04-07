import { useState } from 'react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Search, Plus, Edit2, Trash2, Shield, X } from 'lucide-react'
import { UserAvatar } from '@/components/ui/UserAvatar'
import { useAdminUsers, useDeleteUser, useAdminCreateUser } from '@/hooks/useAdmin'
import { useRoles, useAssignRolesToUser } from '@/hooks/useRbac'
import { getRoleBadgeColor, getRoleDisplayName, toggleId } from './utils'

interface EditingUser {
  id: number
  username: string
  name: string
  roles: string[]
}

export function UsersTab() {
  const [searchQuery, setSearchQuery] = useState('')
  const [page, setPage] = useState(1)
  const [showRoleDialog, setShowRoleDialog] = useState(false)
  const [showAddDialog, setShowAddDialog] = useState(false)
  const [editingUser, setEditingUser] = useState<EditingUser | null>(null)
  const [selectedRoleIds, setSelectedRoleIds] = useState<number[]>([])

  const [addForm, setAddForm] = useState({ username: '', password: '', name: '' })
  const [addRoleIds, setAddRoleIds] = useState<number[]>([])

  const { data: usersData, isLoading } = useAdminUsers({
    page, pageSize: 10, keyword: searchQuery || undefined,
  })
  const { data: rolesData } = useRoles()
  const deleteUserMutation = useDeleteUser()
  const createUserMutation = useAdminCreateUser()
  const assignRolesMutation = useAssignRolesToUser()

  const users = usersData?.list ?? []
  const roles = rolesData?.list ?? []

  const handleManageRoles = (user: EditingUser) => {
    setEditingUser(user)
    const currentRoleIds = roles
      .filter((r) => user.roles.includes(r.name))
      .map((r) => r.id)
    setSelectedRoleIds(currentRoleIds)
    setShowRoleDialog(true)
  }

  const handleSaveRoles = () => {
    if (!editingUser) return
    assignRolesMutation.mutate(
      { userId: editingUser.id, roleIds: selectedRoleIds },
      { onSuccess: () => setShowRoleDialog(false) },
    )
  }

  const handleDeleteUser = (userId: number) => {
    if (!confirm('确定要删除该用户吗？')) return
    deleteUserMutation.mutate(userId)
  }

  const handleAddUser = (e: React.FormEvent) => {
    e.preventDefault()
    createUserMutation.mutate(
      {
        username: addForm.username,
        password: addForm.password,
        name: addForm.name || undefined,
        roleIds: addRoleIds.length > 0 ? addRoleIds : undefined,
      },
      {
        onSuccess: () => {
          setShowAddDialog(false)
          setAddForm({ username: '', password: '', name: '' })
          setAddRoleIds([])
        },
      },
    )
  }

  return (
    <>
      <div className="flex justify-end">
        <Button onClick={() => { setAddForm({ username: '', password: '', name: '' }); setAddRoleIds([]); setShowAddDialog(true) }}>
          <Plus className="h-4 w-4 mr-2" />
          添加用户
        </Button>
      </div>

      <Card className="ink-border">
        <CardHeader className="pb-4">
          <div className="flex items-center gap-4">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="搜索用户名或姓名..."
                className="pl-10"
                value={searchQuery}
                onChange={(e) => { setSearchQuery(e.target.value); setPage(1) }}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-center py-8 text-muted-foreground">加载中...</div>
          ) : users.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">暂无用户数据</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="text-left py-3 px-4 font-medium">ID</th>
                    <th className="text-left py-3 px-4 font-medium">用户</th>
                    <th className="text-left py-3 px-4 font-medium">姓名</th>
                    <th className="text-left py-3 px-4 font-medium">角色</th>
                    <th className="text-left py-3 px-4 font-medium">状态</th>
                    <th className="text-left py-3 px-4 font-medium">注册时间</th>
                    <th className="text-left py-3 px-4 font-medium">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map((user) => (
                    <tr key={user.id} className="border-b last:border-0 hover:bg-muted/50">
                      <td className="py-3 px-4">{user.id}</td>
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-3">
                          <UserAvatar
                            avatar={user.avatar}
                            name={user.name}
                            username={user.username}
                          />
                          <span>{user.username}</span>
                        </div>
                      </td>
                      <td className="py-3 px-4">{user.name}</td>
                      <td className="py-3 px-4">
                        <div className="flex flex-wrap gap-1">
                          {(user.roles && user.roles.length > 0 ? user.roles : [user.role]).map((r) => (
                            <span key={r} className={`px-2 py-1 rounded text-xs ${getRoleBadgeColor(r)}`}>
                              {getRoleDisplayName(r)}
                            </span>
                          ))}
                        </div>
                      </td>
                      <td className="py-3 px-4">
                        <span className={`px-2 py-1 rounded text-xs ${
                          user.isActive ? 'bg-green-500/10 text-green-500' : 'bg-gray-500/10 text-gray-500'
                        }`}>
                          {user.isActive ? '正常' : '禁用'}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-muted-foreground">
                        {new Date(user.createdAt).toLocaleDateString()}
                      </td>
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-2">
                          <Button className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                            onClick={() => handleManageRoles({
                              id: user.id, username: user.username, name: user.name,
                              roles: user.roles && user.roles.length > 0 ? user.roles : [user.role],
                            })}
                            title="分配角色"
                          >
                            <Shield className="h-4 w-4" />
                          </Button>
                          <Button className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0">
                            <Edit2 className="h-4 w-4" />
                          </Button>
                          <Button className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                            onClick={() => handleDeleteUser(user.id)}
                          >
                            <Trash2 className="h-4 w-4 text-cinnabar" />
                          </Button>
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

      {/* Add User Dialog */}
      {showAddDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg w-full max-w-md max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">添加用户</h2>
              <button onClick={() => setShowAddDialog(false)} className="p-1 hover:bg-muted rounded">
                <X className="h-5 w-5" />
              </button>
            </div>
            <form onSubmit={handleAddUser} className="space-y-4 p-6 pt-4 overflow-y-auto">
              <div className="space-y-2">
                <label className="text-sm font-medium">用户名 *</label>
                <Input value={addForm.username} onChange={(e) => setAddForm({ ...addForm, username: e.target.value })}
                  placeholder="用户名（3-50字符）" required minLength={3} />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">密码 *</label>
                <Input type="password" value={addForm.password} onChange={(e) => setAddForm({ ...addForm, password: e.target.value })}
                  placeholder="密码（至少6位）" required minLength={6} />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">姓名</label>
                <Input value={addForm.name} onChange={(e) => setAddForm({ ...addForm, name: e.target.value })}
                  placeholder="可选，留空则使用用户名" />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">角色权限</label>
                <div className="space-y-2 border rounded-lg p-3">
                  {roles.map((role) => (
                    <label key={role.id} className="flex items-center gap-3 cursor-pointer py-1">
                      <input type="checkbox" checked={addRoleIds.includes(role.id)}
                        onChange={() => toggleId(role.id, addRoleIds, setAddRoleIds)} className="rounded border-gray-300" />
                      <div className="flex-1">
                        <span className="text-sm font-medium">{getRoleDisplayName(role.name)}</span>
                        <span className="ml-2 text-xs text-muted-foreground">{role.description}</span>
                      </div>
                    </label>
                  ))}
                  {roles.length === 0 && <p className="text-sm text-muted-foreground">未找到角色，将分配默认角色</p>}
                </div>
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button type="button" className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                  onClick={() => setShowAddDialog(false)}>取消</Button>
                <Button type="submit" disabled={createUserMutation.isPending}>
                  {createUserMutation.isPending ? '创建中...' : '创建'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Role Assignment Dialog */}
      {showRoleDialog && editingUser && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-bold mb-2">分配角色</h2>
            <p className="text-muted-foreground mb-4">
              为用户 "{editingUser.name}" ({editingUser.username}) 分配角色
            </p>
            <div className="space-y-3 mb-6">
              {roles.map((role) => (
                <label key={role.id} className="flex items-center space-x-3 p-3 border rounded-lg cursor-pointer hover:bg-muted/50">
                  <input type="checkbox" checked={selectedRoleIds.includes(role.id)}
                    onChange={(e) => {
                      if (e.target.checked) setSelectedRoleIds([...selectedRoleIds, role.id])
                      else setSelectedRoleIds(selectedRoleIds.filter((id) => id !== role.id))
                    }} className="rounded border-gray-300" />
                  <div className="flex-1">
                    <div className="font-medium">{getRoleDisplayName(role.name)}</div>
                    <div className="text-sm text-muted-foreground">{role.description}</div>
                  </div>
                </label>
              ))}
            </div>
            <div className="flex justify-end gap-2">
              <Button className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                onClick={() => setShowRoleDialog(false)}>取消</Button>
              <Button onClick={handleSaveRoles} disabled={assignRolesMutation.isPending}>
                {assignRolesMutation.isPending ? '保存中...' : '保存'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </>
  )
}
