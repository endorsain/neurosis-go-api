INSERT INTO roles (name, description)
VALUES
    ('admin', 'Administrator role with full access'),
    ('moderator', 'Moderator role with elevated moderation permissions'),
    ('regular', 'Regular user role with basic permissions')
ON CONFLICT (name) DO NOTHING;

INSERT INTO permissions (name, description)
VALUES
    ('posts:create', 'Create posts'),
    ('posts:delete_any', 'Delete any post'),
    ('users:suspend', 'Suspend users'),
    ('staff:promote', 'Promote users to staff roles')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name IN ('posts:create', 'posts:delete_any', 'users:suspend', 'staff:promote')
WHERE r.name = 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name IN ('posts:create', 'posts:delete_any', 'users:suspend')
WHERE r.name = 'moderator'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name IN ('posts:create')
WHERE r.name = 'regular'
ON CONFLICT (role_id, permission_id) DO NOTHING;