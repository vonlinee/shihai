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

// 权限模块中文名称映射
// 与后端 models/permission_codes.go 中的 ModuleXxx 常量保持一致
const MODULE_LABELS: Record<string, string> = {
  user: '用户管理',
  role: '角色管理',
  permission: '权限管理',
  poem: '诗词管理',
  comment: '评论管理',
  correction: '纠错管理',
  announcement: '公告管理',
  forum: '论坛管理',
  quiz: '题目管理',
  feedback: '反馈管理',
  system: '系统管理',
}

export function getModuleDisplayName(module: string): string {
  return MODULE_LABELS[module] ?? module
}
