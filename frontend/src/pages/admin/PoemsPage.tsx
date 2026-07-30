import { useState } from 'react'
import { BookOpen, Crown, User } from 'lucide-react'

import { DynastiesTab } from './DynastyManagement'
import { PoemTypesTab } from './PoemTypesTab'
import { PoemsTab } from './poems/PoemsTab'
import PoetsTab from './PoetsManagement'

type TabKey = 'poems' | 'poemTypes' | 'dynasties' | 'poets'

const tabs: { key: TabKey; label: string; icon: typeof BookOpen }[] = [
  { key: 'poems', label: '诗词', icon: BookOpen },
  { key: 'poemTypes', label: '体裁', icon: BookOpen },
  { key: 'dynasties', label: '朝代', icon: Crown },
  { key: 'poets', label: '诗人', icon: User },
]

export function AdminPoemsPage() {
  const [activeTab, setActiveTab] = useState<TabKey>('poems')

  return (
    <div className="p-8 space-y-6">
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
      {activeTab === 'poemTypes' && <PoemTypesTab />}
      {activeTab === 'dynasties' && <DynastiesTab />}
      {activeTab === 'poets' && <PoetsTab />}
    </div>
  )
}
