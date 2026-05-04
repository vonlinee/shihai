import { useQuery } from '@tanstack/react-query';
import { announcementService } from '@/services/announcementService';

export function useAnnouncements(page = 1, pageSize = 10, pinned?: boolean) {
  return useQuery({
    queryKey: ['announcements', page, pageSize, pinned],
    queryFn: () => announcementService.getAnnouncements(page, pageSize, pinned),
  });
}

export function useAnnouncement(id: number) {
  return useQuery({
    queryKey: ['announcement', id],
    queryFn: () => announcementService.getAnnouncementById(id),
    enabled: !!id,
  });
}
