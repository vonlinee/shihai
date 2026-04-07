// User types
export interface User {
  id: number;
  username: string;
  name: string;
  avatar?: string;
  gender?: 'male' | 'female' | 'other';
  age?: number;
  phone?: string;
  idCard?: string;
  role: 'admin' | 'user' | 'reviewer' | 'editor';
  isActive: boolean;
  createdAt: string;
}

export interface LoginCredentials {
  username: string;
  password: string;
  role?: string;
}

export interface RegisterData {
  username: string;
  password: string;
  confirmPassword: string;
  name: string;
  gender?: 'male' | 'female' | 'other';
  age?: number;
  phone?: string;
  idCard?: string;
  avatar?: string;
}

// Dynasty types
export interface Dynasty {
  id: number;
  name: string;
  period: string;
  description?: string;
  createdAt: string;
}

// Poet types
export interface Poet {
  id: number;
  name: string;
  dynastyId: number;
  dynasty?: Dynasty;
  biography?: string;
  avatar?: string;
  birthYear?: number;
  deathYear?: number;
  createdAt: string;
}

// Poem types
export interface Poem {
  id: number;
  title: string;
  content: string;
  authorId: number;
  author?: Poet;
  dynastyId: number;
  dynasty?: Dynasty;
  genre?: string;
  translation?: string;
  appreciation?: string;
  annotation?: string;
  audioUrl?: string;
  coverImage?: string;
  views: number;
  likes: number;
  dislikes: number;
  favorites: number;
  createdAt: string;
  updatedAt: string;
}

export interface PoemVideo {
  id: number;
  poemId: number;
  title: string;
  videoUrl: string;
  coverImage?: string;
  description?: string;
  duration?: number;
  createdAt: string;
}

// Comment types
export interface Comment {
  id: number;
  poemId: number;
  userId?: number;
  user?: User;
  visitorId?: string;
  visitorName?: string;
  content: string;
  parentId?: number;
  parent?: Comment;
  replies?: Comment[];
  replyCount: number;
  likes: number;
  dislikes: number;
  isDeleted: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CommentVote {
  id: number;
  commentId: number;
  userId?: number;
  visitorId?: string;
  type: 'like' | 'dislike';
  createdAt: string;
}

// Quiz types
export interface Quiz {
  id: number;
  poemId?: number;
  poem?: Poem;
  question: string;
  options: string[];
  correctAnswer: number;
  explanation?: string;
  difficulty: 'easy' | 'medium' | 'hard';
  createdAt: string;
}

export interface QuizRecord {
  id: number;
  userId: number;
  user?: User;
  quizId: number;
  quiz?: Quiz;
  answer: number;
  isCorrect: boolean;
  createdAt: string;
}

// Forum types
export interface ForumPost {
  id: number;
  userId: number;
  user?: User;
  title: string;
  content: string;
  views: number;
  replyCount: number;
  isPinned: boolean;
  isDeleted: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface ForumReply {
  id: number;
  postId: number;
  post?: ForumPost;
  userId: number;
  user?: User;
  content: string;
  parentId?: number;
  isDeleted: boolean;
  createdAt: string;
  updatedAt: string;
}

// Correction types
export interface CorrectionRequest {
  id: number;
  poemId: number;
  poem?: Poem;
  userId: number;
  user?: User;
  type: 'content' | 'translation' | 'appreciation' | 'annotation';
  originalText: string;
  suggestedText: string;
  reason: string;
  status: 'pending' | 'voting' | 'approved' | 'rejected' | 'completed';
  voteCount: number;
  approveCount: number;
  rejectCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface CorrectionVote {
  id: number;
  correctionId: number;
  correction?: CorrectionRequest;
  userId: number;
  user?: User;
  type: 'approve' | 'reject';
  comment?: string;
  createdAt: string;
}

// Announcement types
export interface Announcement {
  id: number;
  title: string;
  content: string;
  isPinned: boolean;
  viewCount: number;
  createdAt: string;
  updatedAt: string;
}

// Feedback types
export interface Feedback {
  id: number;
  userId?: number;
  user?: User;
  visitorId?: string;
  type: 'bug' | 'feature' | 'content' | 'other';
  title: string;
  content: string;
  contact?: string;
  status: 'pending' | 'processing' | 'resolved';
  createdAt: string;
  updatedAt: string;
}

// API response types
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// Search params
export interface PoemSearchParams {
  keyword?: string;
  dynastyId?: number;
  authorId?: number;
  genre?: string;
  page?: number;
  pageSize?: number;
}

