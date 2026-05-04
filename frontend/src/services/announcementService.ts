import { api } from './api';
import type { Announcement } from '@/types';

export interface AnnouncementListResponse {
  list: Announcement[];
  total: number;
  page: number;
  pageSize: number;
}

export const announcementService = {
  getAnnouncements(
    page = 1,
    pageSize = 10,
    pinned?: boolean,
  ): Promise<AnnouncementListResponse> {
    return api.get<AnnouncementListResponse>('/announcements', {
      page,
      pageSize,
      ...(pinned !== undefined ? { pinned: String(pinned) } : {}),
    });
  },

  getAnnouncementById(id: number): Promise<Announcement> {
    return api.get<Announcement>(`/announcements/${id}`);
  },
};
