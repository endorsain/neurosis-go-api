DELETE FROM role_permissions rp
USING roles r, permissions p
WHERE rp.role_id = r.id
  AND rp.permission_id = p.id
  AND r.name = 'admin'
  AND p.name IN ('posts:create', 'posts:delete_any', 'users:suspend', 'staff:promote');

DELETE FROM role_permissions rp
USING roles r, permissions p
WHERE rp.role_id = r.id
  AND rp.permission_id = p.id
  AND r.name = 'moderator'
  AND p.name IN ('posts:create', 'posts:delete_any', 'users:suspend');

DELETE FROM role_permissions rp
USING roles r, permissions p
WHERE rp.role_id = r.id
  AND rp.permission_id = p.id
  AND r.name = 'regular'
  AND p.name IN ('posts:create');

DELETE FROM roles
WHERE name IN ('admin', 'moderator', 'regular');

DELETE FROM permissions
WHERE name IN ('posts:create', 'posts:delete_any', 'users:suspend', 'staff:promote');