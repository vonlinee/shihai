-- =============================================================================
-- 002_root_permissions.sql
-- 为 root 用户初始化所有接口权限点
--
-- 使用说明：
--   1. 应用启动时会通过 InitDefaultRolesAndPermissions 自动同步全部权限给 admin
--      角色，并将 root 用户绑定 admin 角色，正常情况下无需手动执行此脚本。
--   2. 当通过 SQL 直接补齐权限或需要在生产环境快速修复 root 权限时，可执行本脚本。
--   3. 脚本是幂等的，重复执行不会产生重复数据。
--
-- 后续维护：
--   每当后端新增接口时，需要在 backend/internal/models/rbac.go 的 DefaultPermissions
--   中新增对应的权限点，本脚本中的 v_permissions 列表也需要同步追加，
--   以保证 SQL 与代码保持一致。
-- =============================================================================

BEGIN;

-- -----------------------------------------------------------------------------
-- 1) 写入/更新全部权限点（permission 表）
-- -----------------------------------------------------------------------------
WITH v_permissions(code, name, module, description) AS (
    VALUES
        -- user module
        ('user:create',          '创建用户',     'user',         '创建新用户'),
        ('user:read',            '查看用户',     'user',         '查看用户信息'),
        ('user:update',          '更新用户',     'user',         '更新用户信息'),
        ('user:delete',          '删除用户',     'user',         '删除用户'),
        ('user:list',            '用户列表',     'user',         '查看用户列表'),
        -- role module
        ('role:create',          '创建角色',     'role',         '创建新角色'),
        ('role:read',            '查看角色',     'role',         '查看角色信息'),
        ('role:update',          '更新角色',     'role',         '更新角色信息'),
        ('role:delete',          '删除角色',     'role',         '删除角色'),
        ('role:list',            '角色列表',     'role',         '查看角色列表'),
        ('role:assign',          '分配角色',     'role',         '为用户分配角色'),
        -- permission module
        ('permission:list',      '权限列表',     'permission',   '查看权限列表'),
        ('permission:read',      '查看权限',     'permission',   '查看权限详情'),
        ('permission:create',    '创建权限',     'permission',   '创建新权限'),
        ('permission:update',    '更新权限',     'permission',   '更新权限信息'),
        ('permission:delete',    '删除权限',     'permission',   '删除权限'),
        -- poem module
        ('poem:create',          '创建诗词',     'poem',         '添加新诗词'),
        ('poem:read',            '查看诗词',     'poem',         '查看诗词详情'),
        ('poem:update',          '更新诗词',     'poem',         '编辑诗词信息'),
        ('poem:delete',          '删除诗词',     'poem',         '删除诗词'),
        ('poem:list',            '诗词列表',     'poem',         '查看诗词列表'),
        -- comment module
        ('comment:create',       '创建评论',     'comment',      '发表评论'),
        ('comment:read',         '查看评论',     'comment',      '查看评论'),
        ('comment:update',       '更新评论',     'comment',      '更新评论'),
        ('comment:delete',       '删除评论',     'comment',      '删除评论'),
        ('comment:list',         '评论列表',     'comment',      '查看评论列表'),
        ('comment:moderate',     '审核评论',     'comment',      '审核管理评论'),
        -- correction module
        ('correction:create',    '提交纠错',     'correction',   '提交纠错申请'),
        ('correction:read',      '查看纠错',     'correction',   '查看纠错详情'),
        ('correction:update',    '更新纠错',     'correction',   '更新纠错信息'),
        ('correction:review',    '审核纠错',     'correction',   '审核纠错申请'),
        ('correction:list',      '纠错列表',     'correction',   '查看纠错列表'),
        -- announcement module
        ('announcement:create',  '创建公告',     'announcement', '创建公告'),
        ('announcement:read',    '查看公告',     'announcement', '查看公告'),
        ('announcement:update',  '更新公告',     'announcement', '更新公告'),
        ('announcement:delete',  '删除公告',     'announcement', '删除公告'),
        ('announcement:list',    '公告列表',     'announcement', '查看公告列表'),
        -- forum module
        ('forum:create',         '发表帖子',     'forum',        '发表论坛帖子'),
        ('forum:read',           '查看帖子',     'forum',        '查看论坛帖子'),
        ('forum:update',         '更新帖子',     'forum',        '更新论坛帖子'),
        ('forum:delete',         '删除帖子',     'forum',        '删除论坛帖子'),
        ('forum:list',           '帖子列表',     'forum',        '查看帖子列表'),
        ('forum:moderate',       '审核帖子',     'forum',        '审核管理论坛帖子'),
        -- quiz module
        ('quiz:create',          '创建题目',     'quiz',         '创建问答题目'),
        ('quiz:read',            '查看题目',     'quiz',         '查看问答题目'),
        ('quiz:update',          '更新题目',     'quiz',         '更新问答题目'),
        ('quiz:delete',          '删除题目',     'quiz',         '删除问答题目'),
        ('quiz:list',            '题目列表',     'quiz',         '查看题目列表'),
        -- feedback module
        ('feedback:create',      '提交反馈',     'feedback',     '提交意见反馈'),
        ('feedback:read',        '查看反馈',     'feedback',     '查看反馈详情'),
        ('feedback:update',      '更新反馈',     'feedback',     '更新反馈状态'),
        ('feedback:delete',      '删除反馈',     'feedback',     '删除反馈'),
        ('feedback:list',        '反馈列表',     'feedback',     '查看反馈列表'),
        -- system module
        ('system:config',        '系统配置',     'system',       '系统配置管理'),
        ('system:log',           '系统日志',     'system',       '查看系统日志'),
        ('system:backup',        '系统备份',     'system',       '系统备份管理')
)
INSERT INTO permission (id, code, name, module, description, is_active, created_at, updated_at, created_by, updated_by)
SELECT
    -- 简易整型 ID 占位（如使用雪花 ID 应在应用层生成；这里使用 nextval/序列或时间戳保证唯一）
    (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT
        + (ROW_NUMBER() OVER (ORDER BY v.code))::BIGINT,
    v.code, v.name, v.module, v.description, TRUE, NOW(), NOW(), 0, 0
FROM v_permissions v
ON CONFLICT (code) DO UPDATE
SET name        = EXCLUDED.name,
    module      = EXCLUDED.module,
    description = EXCLUDED.description,
    is_active   = TRUE,
    updated_at  = NOW();

-- -----------------------------------------------------------------------------
-- 2) 确保 admin 角色存在
-- -----------------------------------------------------------------------------
INSERT INTO role (id, name, description, is_active, created_at, updated_at, created_by, updated_by)
VALUES (
    (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT,
    'admin', '超级管理员，拥有所有权限', TRUE, NOW(), NOW(), 0, 0
)
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description,
    is_active   = TRUE,
    updated_at  = NOW();

-- -----------------------------------------------------------------------------
-- 3) 将 permission 表中所有权限编码绑定到 admin 角色
--    （未来新增权限点只要写入 permission 表，重新执行本段即可同步）
-- -----------------------------------------------------------------------------
INSERT INTO role_permission (id, role_id, code, created_at, updated_at, created_by, updated_by)
SELECT
    (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT
        + (ROW_NUMBER() OVER (ORDER BY p.code))::BIGINT,
    r.id, p.code, NOW(), NOW(), 0, 0
FROM role r
CROSS JOIN permission p
WHERE r.name = 'admin'
  AND p.is_active = TRUE
ON CONFLICT (role_id, code) DO NOTHING;

-- -----------------------------------------------------------------------------
-- 4) 确保 root 用户存在并绑定 admin 角色
--    默认密码 shihai2024（bcrypt cost=10），首次登录后请立即修改！
--    bcrypt 哈希：$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
-- -----------------------------------------------------------------------------
INSERT INTO "user" (id, username, password, name, is_active, created_at, updated_at, created_by, updated_by)
VALUES (
    (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT,
    'root',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    '超级管理员',
    TRUE, NOW(), NOW(), 0, 0
)
ON CONFLICT (username) DO UPDATE
SET is_active  = TRUE,
    updated_at = NOW();

INSERT INTO user_role (id, user_id, role_id, created_at, updated_at, created_by, updated_by)
SELECT
    (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT,
    u.id, r.id, NOW(), NOW(), 0, 0
FROM "user" u, role r
WHERE u.username = 'root'
  AND r.name = 'admin'
ON CONFLICT (user_id, role_id) DO NOTHING;

COMMIT;

-- =============================================================================
-- 验证脚本（可选，单独执行查看 root 用户权限）
-- =============================================================================
-- SELECT u.username, r.name AS role, rp.code AS permission
-- FROM "user" u
-- JOIN user_role ur ON ur.user_id = u.id
-- JOIN role r       ON r.id = ur.role_id
-- JOIN role_permission rp ON rp.role_id = r.id
-- WHERE u.username = 'root'
-- ORDER BY rp.code;
