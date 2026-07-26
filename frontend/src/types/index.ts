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
  role: string;
  roles?: string[];
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

// Author types
export interface Author {
  id: number;
  name: string;
  biography?: string;
  avatar?: string;
}

export interface PoemAnnotation {
  id: string;
  poemId: string;
  displayNo: number;
  targetField: 'content';
  startLine: number;
  startOffset: number;
  endLine: number;
  endOffset: number;
  selectedText: string;
  title?: string;
  content: string;
  type: 'note' | string;
  displayOrder: number;
  createdAt: string;
  updatedAt: string;
}

// Poet types
export interface Poet {
  id: number;
  authorId: number;
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
  content: string[];
  /** 与 content 每个文本元素对应的平仄标记，只包含平/仄/? 或空字符串。 */
  pingze: string[];
  authorId: number;
  author?: Author;
  dynastyId: number;
  dynasty?: Dynasty;
  genre?: string;
  translation?: string;
  appreciation?: string;
  annotation?: string;
  annotations?: PoemAnnotation[];
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

// Work collection types
export interface WorkCollectionItem {
  id: number;
  collectionId: number;
  workType: 'poem' | 'novel' | 'article' | string;
  workId: number;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

export interface WorkCollection {
  id: number;
  title: string;
  description?: string;
  coverImage?: string;
  itemCount: number;
  isPublished: boolean;
  items?: WorkCollectionItem[];
  createdAt: string;
  updatedAt: string;
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
export interface ForumUser {
  /** 用户 Snowflake ID，以字符串传输。 */
  id: string;
  /** 用户名。 */
  username: string;
  /** 用户显示名称。 */
  name: string;
  /** 用户头像地址。 */
  avatar?: string;
}

export interface ForumPost {
  /** 帖子 Snowflake ID，以字符串传输。 */
  id: string;
  /** 发帖用户 Snowflake ID，以字符串传输。 */
  userId: string;
  /** 发帖用户摘要。 */
  user?: ForumUser;
  /** 帖子标题。 */
  title: string;
  /** 帖子正文。 */
  content: string;
  /** 浏览量。 */
  views: number;
  /** 回复数量。 */
  replyCount: number;
  /** 是否置顶。 */
  isPinned: boolean;
  /** 是否已被业务删除。 */
  isDeleted: boolean;
  /** 创建时间，ISO 日期字符串。 */
  createdAt: string;
  /** 更新时间，ISO 日期字符串。 */
  updatedAt: string;
}

export interface ForumReply {
  /** 回复 Snowflake ID，以字符串传输。 */
  id: string;
  /** 所属帖子 Snowflake ID，以字符串传输。 */
  postId: string;
  /** 回复用户 Snowflake ID，以字符串传输。 */
  userId: string;
  /** 回复用户摘要。 */
  user?: ForumUser;
  /** 回复正文。 */
  content: string;
  /** 父回复 Snowflake ID，楼中楼回复时存在。 */
  parentId?: string;
  /** 父回复摘要，楼中楼回复时可能返回。 */
  parent?: ForumReply;
  /** 是否已被业务删除。 */
  isDeleted: boolean;
  /** 创建时间，ISO 日期字符串。 */
  createdAt: string;
  /** 更新时间，ISO 日期字符串。 */
  updatedAt: string;
}

// Correction types
export interface CorrectionPoemSummary {
  /** 诗词 Snowflake ID，以字符串传输。 */
  id: string;
  /** 诗词标题。 */
  title: string;
}

export interface CorrectionUserSummary {
  /** 用户 Snowflake ID，以字符串传输。 */
  id: string;
  /** 用户名。 */
  username: string;
  /** 用户显示名称。 */
  name: string;
  /** 用户头像地址。 */
  avatar?: string;
}

export interface CorrectionRequest {
  /** 纠错申请 Snowflake ID，以字符串传输。 */
  id: string;
  /** 被纠错诗词的 Snowflake ID，以字符串传输。 */
  poemId: string;
  /** 被纠错诗词摘要。 */
  poem?: CorrectionPoemSummary;
  /** 提交用户的 Snowflake ID，以字符串传输。 */
  userId: string;
  /** 提交用户摘要。 */
  user?: CorrectionUserSummary;
  /** 纠错类型。 */
  type: 'content' | 'translation' | 'appreciation' | 'annotation';
  /** 被纠错的原文内容。 */
  originalText: string;
  /** 用户建议修改后的内容。 */
  suggestedText: string;
  /** 用户提交的纠错理由。 */
  reason: string;
  /** 纠错流程状态。 */
  status: 'pending' | 'voting' | 'approved' | 'rejected' | 'completed';
  /** 投票总数。 */
  voteCount: number;
  /** 支持票数。 */
  approveCount: number;
  /** 反对票数。 */
  rejectCount: number;
  /** 创建时间，ISO 日期字符串。 */
  createdAt: string;
  /** 更新时间，ISO 日期字符串。 */
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
  dynasty?: string;
  dynastyId?: number;
  authorId?: number;
  genre?: string;
  page?: number;
  pageSize?: number;
}

