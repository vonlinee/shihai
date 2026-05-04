export function getRoleBadgeColor(role: string) {
  switch (role) {
    case 'admin': return 'bg-red-500/10 text-red-500'
    case 'editor': return 'bg-purple-500/10 text-purple-500'
    case 'reviewer': return 'bg-blue-500/10 text-blue-500'
    default: return 'bg-green-500/10 text-green-500'
  }
}

export function getRoleDisplayName(role: string) {
  switch (role) {
    case 'admin': return '管理员'
    case 'editor': return '编辑'
    case 'reviewer': return '审核员'
    default: return '普通用户'
  }
}

export function toggleId(id: number, list: number[], setter: (ids: number[]) => void) {
  if (list.includes(id)) {
    setter(list.filter((rid) => rid !== id))
  } else {
    setter([...list, id])
  }
}
