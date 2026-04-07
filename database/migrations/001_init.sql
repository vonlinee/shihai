-- Initial database schema for Shihai Poetry Platform

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(50) NOT NULL,
    avatar VARCHAR(255),
    gender VARCHAR(10),
    age INTEGER,
    phone VARCHAR(20),
    id_card VARCHAR(18),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Dynasties table
CREATE TABLE IF NOT EXISTS dynasties (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    period VARCHAR(100),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Poets table
CREATE TABLE IF NOT EXISTS poets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    dynasty_id INTEGER REFERENCES dynasties(id),
    biography TEXT,
    avatar VARCHAR(255),
    birth_year INTEGER,
    death_year INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Poems table
CREATE TABLE IF NOT EXISTS poems (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER REFERENCES poets(id),
    dynasty_id INTEGER REFERENCES dynasties(id),
    genre VARCHAR(50),
    translation TEXT,
    appreciation TEXT,
    annotation TEXT,
    audio_url VARCHAR(500),
    cover_image VARCHAR(500),
    views INTEGER DEFAULT 0,
    likes INTEGER DEFAULT 0,
    dislikes INTEGER DEFAULT 0,
    favorites INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Comments table
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    poem_id INTEGER NOT NULL REFERENCES poems(id),
    user_id INTEGER REFERENCES users(id),
    visitor_id VARCHAR(64),
    visitor_name VARCHAR(50),
    content TEXT NOT NULL,
    parent_id INTEGER REFERENCES comments(id),
    reply_count INTEGER DEFAULT 0,
    likes INTEGER DEFAULT 0,
    dislikes INTEGER DEFAULT 0,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Comment votes table
CREATE TABLE IF NOT EXISTS comment_votes (
    id SERIAL PRIMARY KEY,
    comment_id INTEGER NOT NULL REFERENCES comments(id),
    user_id INTEGER REFERENCES users(id),
    visitor_id VARCHAR(64),
    type VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quizzes table
CREATE TABLE IF NOT EXISTS quizzes (
    id SERIAL PRIMARY KEY,
    poem_id INTEGER REFERENCES poems(id),
    question TEXT NOT NULL,
    options JSON NOT NULL,
    correct_answer INTEGER NOT NULL,
    explanation TEXT,
    difficulty VARCHAR(10) DEFAULT 'medium',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quiz records table
CREATE TABLE IF NOT EXISTS quiz_records (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    quiz_id INTEGER NOT NULL REFERENCES quizzes(id),
    answer INTEGER NOT NULL,
    is_correct BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Forum posts table
CREATE TABLE IF NOT EXISTS forum_posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    views INTEGER DEFAULT 0,
    reply_count INTEGER DEFAULT 0,
    is_pinned BOOLEAN DEFAULT false,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Forum replies table
CREATE TABLE IF NOT EXISTS forum_replies (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL REFERENCES forum_posts(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    parent_id INTEGER REFERENCES forum_replies(id),
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Correction requests table
CREATE TABLE IF NOT EXISTS correction_requests (
    id SERIAL PRIMARY KEY,
    poem_id INTEGER NOT NULL REFERENCES poems(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    type VARCHAR(20) NOT NULL,
    original_text TEXT NOT NULL,
    suggested_text TEXT NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    vote_count INTEGER DEFAULT 0,
    approve_count INTEGER DEFAULT 0,
    reject_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Correction votes table
CREATE TABLE IF NOT EXISTS correction_votes (
    id SERIAL PRIMARY KEY,
    correction_id INTEGER NOT NULL REFERENCES correction_requests(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    type VARCHAR(10) NOT NULL,
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Announcements table
CREATE TABLE IF NOT EXISTS announcements (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    is_pinned BOOLEAN DEFAULT false,
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Feedback table
CREATE TABLE IF NOT EXISTS feedback (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    visitor_id VARCHAR(64),
    type VARCHAR(20) NOT NULL,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    contact VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Operation logs table
CREATE TABLE IF NOT EXISTS operation_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    action VARCHAR(50) NOT NULL,
    target VARCHAR(50),
    target_id INTEGER,
    detail TEXT,
    ip VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_poems_author_id ON poems(author_id);
CREATE INDEX IF NOT EXISTS idx_poems_dynasty_id ON poems(dynasty_id);
CREATE INDEX IF NOT EXISTS idx_comments_poem_id ON comments(poem_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_comment_votes_comment_id ON comment_votes(comment_id);
CREATE INDEX IF NOT EXISTS idx_forum_posts_user_id ON forum_posts(user_id);
CREATE INDEX IF NOT EXISTS idx_forum_replies_post_id ON forum_replies(post_id);
CREATE INDEX IF NOT EXISTS idx_correction_requests_poem_id ON correction_requests(poem_id);
CREATE INDEX IF NOT EXISTS idx_correction_requests_user_id ON correction_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_correction_votes_correction_id ON correction_votes(correction_id);

-- Insert sample dynasties
INSERT INTO dynasties (name, period, description) VALUES
('唐', '618-907', '唐朝是中国历史上最辉煌的朝代之一，诗歌发展达到巅峰'),
('宋', '960-1279', '宋朝以词闻名，是中国文学史上的重要时期'),
('元', '1271-1368', '元朝以曲著称，杂剧和散曲发展迅速'),
('明', '1368-1644', '明朝诗词继续发展，小说创作也取得巨大成就'),
('清', '1644-1912', '清朝诗词流派众多，呈现出多元化的发展态势');

-- Insert sample poets
INSERT INTO poets (name, dynasty_id, biography, birth_year, death_year) VALUES
('李白', 1, '字太白，号青莲居士，又号谪仙人，唐代伟大的浪漫主义诗人', 701, 762),
('杜甫', 1, '字子美，自号少陵野老，唐代伟大的现实主义诗人', 712, 770),
('苏轼', 2, '字子瞻，号东坡居士，北宋著名文学家、书法家、画家', 1037, 1101),
('李清照', 2, '号易安居士，宋代女词人，婉约词派代表', 1084, 1155);

-- Insert sample poems
INSERT INTO poems (title, content, author_id, dynasty_id, genre, translation, appreciation) VALUES
('静夜思', '床前明月光，疑是地上霜。举头望明月，低头思故乡。', 1, 1, '五言绝句', 
'明亮的月光洒在床前的窗户纸上，好像地上泛起了一层霜。我禁不住抬起头来，看那天窗外空中的一轮明月，不由得低头沉思，想起远方的家乡。',
'这首诗写的是在寂静的月夜思念家乡的感受。诗的前两句，是写诗人在作客他乡的特定环境中一刹那间所产生的错觉。'),
('春晓', '春眠不觉晓，处处闻啼鸟。夜来风雨声，花落知多少。', 2, 1, '五言绝句',
'春天睡醒不觉天晓，处处听到鸟儿啼叫。夜里风雨声音，花儿不知道吹落了多少。',
'这首诗是诗人隐居在鹿门山时所做，意境十分优美。诗人抓住春晨生活的一刹那，描绘了春晨的景色和感受。');
