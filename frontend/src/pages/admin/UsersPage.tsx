import { useState } from 'react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Search, Plus, Edit2, Trash2, Shield } from 'lucide-react'
import { useAdminUsers, useDeleteUser } from '@/hooks/useAdmin'
import { useRoles, useAssignRolesToUser } from '@/hooks/useRbac'

interface EditingUser {
  id: number
  username: string
  name: string
  roles: string[]
}

export function AdminUsersPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [page, setPage] = useState(1)
  const [showRoleDialog, setShowRoleDialog] = useState(false)
  const [editingUser, setEditingUser] = useState<EditingUser | null>(null)
  const [selectedRoleIds, setSelectedRoleIds] = useState<number[]>([])

  const { data: usersData, isLoading } = useAdminUsers({
    page,
    pageSize: 10,
    keyword: searchQuery || undefined,
  })
  const { data: rolesData } = useRoles()
  const deleteUserMutation = useDeleteUser()
  const assignRolesMutation = useAssignRolesToUser()

  const users = usersData?.list ?? []
  const roles = rolesData?.list ?? []

  const handleManageRoles = (user: EditingUser) => {
    setEditingUser(user)
    // Find role IDs matching current user roles
    const currentRoleIds = roles
      .filter((r) => user.roles.includes(r.name))
      .map((r) => r.id)
    setSelectedRoleIds(currentRoleIds)
    setShowRoleDialog(true)
  }

  const handleSaveRoles = async () => {
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

  const getRoleBadgeColor = (role: string) => {
    switch (role) {
      case 'admin':
        return 'bg-red-500/10 text-red-500'
      case 'editor':
        return 'bg-purple-500/10 text-purple-500'
      case 'reviewer':
        return 'bg-blue-500/10 text-blue-500'
      default:
        return 'bg-green-500/10 text-green-500'
    }
  }

  const getRoleDisplayName = (role: string) => {
    switch (role) {
      case 'admin':
        return '管理员'
      case 'editor':
        return '编辑'
      case 'reviewer':
        return '审核员'
      default:
        return '普通用户'
    }
  }

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-serif font-bold text-ink">用户管理</h1>
        <Button>
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
                onChange={(e) => {
                  setSearchQuery(e.target.value)
                  setPage(1)
                }}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-center py-8 text-muted-foreground">加载中...</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="text-left py-3 px-4 font-medium">ID</th>
                    <th className="text-left py-3 px-4 font-medium">用户名</th>
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
                      <td className="py-3 px-4">{user.username}</td>
                      <td className="py-3 px-4">{user.name}</td>
                      <td className="py-3 px-4">
                        <span className={`px-2 py-1 rounded text-xs ${getRoleBadgeColor(user.role)}`}>
                          {getRoleDisplayName(user.role)}
                        </span>
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
                          <Button
                            className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                            onClick={() => handleManageRoles({
                              id: user.id,
                              username: user.username,
                              name: user.name,
                              roles: [user.role],
                            })}
                            title="分配角色"
                          >
                            <Shield className="h-4 w-4" />
                          </Button>
                          <Button className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0">
                            <Edit2 className="h-4 w-4" />
                          </Button>
                          <Button
                            className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
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
                <label
                  key={role.id}
                  className="flex items-center space-x-3 p-3 border rounded-lg cursor-pointer hover:bg-muted/50"
                >
                  <input
                    type="checkbox"
                    checked={selectedRoleIds.includes(role.id)}
                    onChange={(e) => {
                      if (e.target.checked) {
                        setSelectedRoleIds([...selectedRoleIds, role.id])
                      } else {
                        setSelectedRoleIds(selectedRoleIds.filter((id) => id !== role.id))
                      }
                    }}
                    className="rounded border-gray-300"
                  />
                  <div className="flex-1">
                    <div className="font-medium">{getRoleDisplayName(role.name)}</div>
                    <div className="text-sm text-muted-foreground">{role.description}</div>
                  </div>
                </label>
              ))}
            </div>
            <div className="flex justify-end gap-2">
              <Button
                className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                onClick={() => setShowRoleDialog(false)}
              >
                取消
              </Button>
              <Button onClick={handleSaveRoles} disabled={assignRolesMutation.isPending}>
                {assignRolesMutation.isPending ? '保存中...' : '保存'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
