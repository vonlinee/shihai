import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'

interface UserAvatarProps {
  avatar?: string
  name?: string
  username?: string
  className?: string
}

/**
 * Reusable user avatar component used across both public and admin interfaces.
 * Shows the user's avatar image if available, otherwise displays the first
 * character of their name or username as a fallback.
 */
export function UserAvatar({ avatar, name, username, className = 'h-8 w-8' }: UserAvatarProps) {
  const fallback = name?.charAt(0) || username?.charAt(0) || '?'

  return (
    <Avatar className={className}>
      {avatar ? (
        <AvatarImage src={avatar} alt={name || username || ''} />
      ) : (
        <AvatarFallback className="bg-primary/10 text-primary text-sm font-medium">
          {fallback}
        </AvatarFallback>
      )}
    </Avatar>
  )
}
