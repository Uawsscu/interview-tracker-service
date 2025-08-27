CREATE TABLE IF NOT EXISTS users (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name        text NOT NULL,
  email       text NOT NULL,
  password    text NOT NULL,            -- เก็บ hash
  role_id     uuid NOT NULL,
  is_active   boolean NOT NULL DEFAULT true,
  created_by  uuid NULL REFERENCES users(id) ON DELETE SET NULL,
  updated_by  uuid NULL REFERENCES users(id) ON DELETE SET NULL,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS roles (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code        text NOT NULL UNIQUE,     -- เช่น: admin, editor, viewer
  name        text NOT NULL,
  description text DEFAULT '',
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS permissions (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code        text NOT NULL UNIQUE,     -- เช่น: card.read, card.write, user.delete
  name        text NOT NULL,
  description text DEFAULT '',
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS role_permissions (
  role_id       uuid NOT NULL,
  permission_id uuid NOT NULL,
  is_active     boolean NOT NULL DEFAULT true,
  created_by    uuid NULL REFERENCES users(id) ON DELETE SET NULL,
  updated_by    uuid NULL REFERENCES users(id) ON DELETE SET NULL,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (role_id, permission_id),
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- =========================
-- Comments for table: users
-- =========================
COMMENT ON TABLE users IS 'ตารางเก็บข้อมูลผู้ใช้งานระบบ';
COMMENT ON COLUMN users.id IS 'รหัสผู้ใช้ (UUID)';
COMMENT ON COLUMN users.name IS 'ชื่อผู้ใช้';
COMMENT ON COLUMN users.email IS 'อีเมลของผู้ใช้';
COMMENT ON COLUMN users.password IS 'รหัสผ่าน (เก็บ hash)';
COMMENT ON COLUMN users.role_id IS 'อ้างอิงไปยังตาราง roles';
COMMENT ON COLUMN users.is_active IS 'สถานะผู้ใช้ (true=ใช้งาน, false=ปิดการใช้งาน)';
COMMENT ON COLUMN users.created_by IS 'UUID ของผู้สร้าง record';
COMMENT ON COLUMN users.updated_by IS 'UUID ของผู้แก้ไข record ล่าสุด';
COMMENT ON COLUMN users.created_at IS 'วันและเวลาที่สร้าง record';
COMMENT ON COLUMN users.updated_at IS 'วันและเวลาที่แก้ไข record ล่าสุด';

-- =========================
-- Comments for table: roles
COMMENT ON TABLE roles IS 'ตารางเก็บข้อมูล Role ของผู้ใช้';
COMMENT ON COLUMN roles.id IS 'รหัส Role (UUID)';
COMMENT ON COLUMN roles.code IS 'โค้ด Role เช่น admin, editor, viewer';
COMMENT ON COLUMN roles.name IS 'ชื่อ Role';
COMMENT ON COLUMN roles.description IS 'คำอธิบาย Role';
COMMENT ON COLUMN roles.is_active IS 'สถานะ Role';
COMMENT ON COLUMN roles.created_at IS 'วันและเวลาที่สร้าง record';
COMMENT ON COLUMN roles.updated_at IS 'วันและเวลาที่แก้ไข record ล่าสุด';

-- Comments for table: permissions
COMMENT ON TABLE permissions IS 'ตารางเก็บสิทธิ์การเข้าถึง (Permission)';
COMMENT ON COLUMN permissions.id IS 'รหัส Permission (UUID)';
COMMENT ON COLUMN permissions.code IS 'โค้ด Permission เช่น card.read, card.write, user.delete';
COMMENT ON COLUMN permissions.name IS 'ชื่อ Permission';
COMMENT ON COLUMN permissions.description IS 'คำอธิบาย Permission';
COMMENT ON COLUMN permissions.is_active IS 'สถานะ Permission';
COMMENT ON COLUMN permissions.created_at IS 'วันและเวลาที่สร้าง record';
COMMENT ON COLUMN permissions.updated_at IS 'วันและเวลาที่แก้ไข record ล่าสุด';

-- Comments for table: role_permissions
COMMENT ON TABLE role_permissions IS 'ตาราง Mapping ระหว่าง Role และ Permission';
COMMENT ON COLUMN role_permissions.role_id IS 'อ้างอิงไปยัง Role';
COMMENT ON COLUMN role_permissions.permission_id IS 'อ้างอิงไปยัง Permission';
COMMENT ON COLUMN role_permissions.is_active IS 'สถานะการใช้งาน mapping';
COMMENT ON COLUMN role_permissions.created_by IS 'UUID ของผู้สร้าง record';
COMMENT ON COLUMN role_permissions.updated_by IS 'UUID ของผู้แก้ไข record ล่าสุด';
COMMENT ON COLUMN role_permissions.created_at IS 'วันและเวลาที่สร้าง record';
COMMENT ON COLUMN role_permissions.updated_at IS 'วันและเวลาที่แก้ไข record ล่าสุด';

-- ============ SEED =============
-- Insert Roles
INSERT INTO roles (id, code, name, description)
VALUES
  ('a1bf3d66-e4ae-4d73-89c6-917f0f301003', 'admin', 'Administrator', 'Full access to system'),
  ('e68dab6b-f66a-4c92-9df2-888f1e813c05', 'editor', 'Editor', 'Can create, edit, and delete contents'),
  ('1d7b3728-0be0-414c-82b8-29ac47287ef9', 'viewer', 'Viewer', 'Read-only access');

-- Insert Permissions
INSERT INTO permissions (code, name, description)
VALUES
  ('comment_add', 'Add Comment', 'สามารถเพิ่มความคิดเห็นได้'),
  ('comment_edit', 'Edit Comment', 'สามารถแก้ไขความคิดเห็นได้'),
  ('comment_view', 'View Comment', 'สามารถดูความคิดเห็นได้'),
  ('comment_delete', 'Delete Comment', 'สามารถลบความคิดเห็นได้'),
  ('card_add', 'Add Card', 'สามารถเพิ่มการ์ดได้'),
  ('card_edit', 'Edit Card', 'สามารถแก้ไขการ์ดได้'),
  ('card_view', 'View Card', 'สามารถดูการ์ดได้'),
  ('card_delete', 'Delete Card', 'สามารถลบการ์ดได้');

-- Admin: ได้ทุกสิทธิ์
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.code = 'admin';

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN ('comment_add','comment_edit','comment_view','card_add','card_edit','card_view')
WHERE r.code = 'editor';
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN ('comment_view','card_view')
WHERE r.code = 'viewer';

-- Insert Users
INSERT INTO users (id, name, email, password, role_id)
SELECT '2b654344-da85-487c-ad02-9976e6d57cb1', 'Admin User', 'admin@example.com', 'hashed_admin', r.id
FROM roles r WHERE r.code = 'admin';

INSERT INTO users (id, name, email, password, role_id)
SELECT 'eb64b4ce-e11e-4f83-9385-9cfd7bfdf74c', 'Editor User', 'editor@example.com', 'hashed_editor', r.id
FROM roles r WHERE r.code = 'editor';

INSERT INTO users (id, name, email, password, role_id)
SELECT '46389484-69b0-45b1-9e0c-23736b5ce9bc', 'Viewer User', 'viewer@example.com', 'hashed_viewer', r.id
FROM roles r WHERE r.code = 'viewer';