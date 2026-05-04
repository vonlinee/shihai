import { useState } from 'react'
import { Users, Shield, Key } from 'lucide-react'
import { UsersTab } from './users/UsersTab'
import { RolesTab } from './users/RolesTab'
import { PermissionsTab } from './users/PermissionsTab'

type TabKey = 'users' | 'roles' | 'permissions'

const tabs: { key: TabKey; label: string; icon: typeof Users }[] = [
  { key: 'users', label: '用户', icon: Users },
  { key: 'roles', label: '角色', icon: Shield },
  { key: 'permissions', label: '权限', icon: Key },
]

export function AdminUsersPage() {
  const [activeTab, setActiveTab] = useState<TabKey>('users')

  return (
    <div className="p-8 space-y-6">
      {/* Tab Bar */}
      <div className="flex gap-1 border-b">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`flex items-center gap-2 px-4 py-3 text-sm font-medium border-b-2 -mb-px transition-colors ${
              activeTab === tab.key
                ? 'border-primary text-foreground'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <tab.icon className="h-4 w-4" />
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'users' && <UsersTab />}
      {activeTab === 'roles' && <RolesTab />}
      {activeTab === 'permissions' && <PermissionsTab />}
    </div>
  )
}
