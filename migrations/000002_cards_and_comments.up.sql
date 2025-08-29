-- ตารางสถานะของการ์ด (Card Statuses)
CREATE TABLE IF NOT EXISTS card_statuses (
  status_code text PRIMARY KEY,         -- รหัสสถานะ เช่น 'todo', 'in_progress', 'done'
  name        text NOT NULL,            -- ชื่อสถานะ
  description text DEFAULT ''           -- คำอธิบายเพิ่มเติม
);

-- ค่า seed เริ่มต้น
INSERT INTO card_statuses (status_code, name, description) VALUES
  ('todo', 'To Do', 'ยังไม่ได้เริ่ม'),
  ('in_progress', 'In Progress', 'กำลังดำเนินการ'),
  ('done', 'Done', 'เสร็จสมบูรณ์')
ON CONFLICT DO NOTHING;

-- ตารางการ์ดนัดสัมภาษณ์ (Interview Cards)
CREATE TABLE IF NOT EXISTS cards (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),  -- รหัสการ์ด (UUID)
  title          text NOT NULL,                               -- ชื่อการ์ด/หัวข้อ
  description    text NOT NULL DEFAULT '',                    -- รายละเอียดการ์ด
  candidate_name text NOT NULL,                               -- ชื่อผู้สมัครงาน
  scheduled_at   timestamptz NOT NULL,                        -- วันและเวลาที่นัดสัมภาษณ์
  status_code    text NOT NULL REFERENCES card_statuses(status_code), -- สถานะอ้างอิงจาก card_statuses
  created_by     uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT, -- ผู้สร้าง
  assignee_id    uuid NULL REFERENCES users(id) ON DELETE SET NULL,     -- ผู้รับผิดชอบ
  created_at     timestamptz NOT NULL DEFAULT now(),           -- วันและเวลาที่สร้าง
  updated_at     timestamptz NOT NULL DEFAULT now()            -- วันและเวลาที่แก้ไขล่าสุด
);

-- ความคิดเห็นของการ์ด
CREATE TABLE IF NOT EXISTS card_comments (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),       -- รหัสคอมเมนต์ (UUID)
  card_id    uuid NOT NULL REFERENCES cards(id) ON DELETE CASCADE, -- อ้างอิงไปยังการ์ด
  author_id  uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT, -- ผู้เขียนคอมเมนต์
  content    text NOT NULL,                                    -- เนื้อหาคอมเมนต์
  created_at timestamptz NOT NULL DEFAULT now(),               -- วันและเวลาที่สร้าง
  updated_at timestamptz NOT NULL DEFAULT now()                -- วันและเวลาที่แก้ไขล่าสุด
);

-- ประวัติความคืบหน้า/กิจกรรมของการ์ด
CREATE TABLE IF NOT EXISTS card_progress_logs (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),       -- รหัส log (UUID)
  card_id    uuid NOT NULL REFERENCES cards(id) ON DELETE CASCADE, -- อ้างอิงไปยังการ์ด
  actor_id   uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT, -- ผู้กระทำ
  message    text NOT NULL,                                    -- รายละเอียดข้อความ log
  created_at timestamptz NOT NULL DEFAULT now()                -- วันและเวลาที่บันทึก log
);

-- ดัชนี
CREATE INDEX IF NOT EXISTS idx_cards_status ON cards(status_code);
CREATE INDEX IF NOT EXISTS idx_comments_card ON card_comments(card_id);
CREATE INDEX IF NOT EXISTS idx_progress_card ON card_progress_logs(card_id);

-- ======================
-- COMMENTS
-- ======================

-- ตารางสถานะ
COMMENT ON TABLE card_statuses IS 'ตารางเก็บสถานะของการ์ด (To Do, In Progress, Done)';
COMMENT ON COLUMN card_statuses.status_code IS 'รหัสสถานะ เช่น todo, in_progress, done';
COMMENT ON COLUMN card_statuses.name IS 'ชื่อสถานะ เช่น To Do, In Progress';
COMMENT ON COLUMN card_statuses.description IS 'คำอธิบายเพิ่มเติมของสถานะ';

-- ตารางการ์ด
COMMENT ON TABLE cards IS 'ตารางเก็บการ์ดนัดสัมภาษณ์ (Interview Cards)';
COMMENT ON COLUMN cards.id IS 'รหัสการ์ด (UUID)';
COMMENT ON COLUMN cards.title IS 'หัวข้อการ์ด';
COMMENT ON COLUMN cards.description IS 'รายละเอียดการ์ด';
COMMENT ON COLUMN cards.candidate_name IS 'ชื่อผู้สมัครงาน';
COMMENT ON COLUMN cards.scheduled_at IS 'วันและเวลาที่นัดสัมภาษณ์';
COMMENT ON COLUMN cards.status_code IS 'สถานะการ์ด อ้างอิงไปยัง card_statuses';
COMMENT ON COLUMN cards.created_by IS 'UUID ของผู้สร้างการ์ด';
COMMENT ON COLUMN cards.assignee_id IS 'UUID ของผู้รับผิดชอบ (nullable)';
COMMENT ON COLUMN cards.created_at IS 'วันและเวลาที่สร้างการ์ด';
COMMENT ON COLUMN cards.updated_at IS 'วันและเวลาที่แก้ไขการ์ดล่าสุด';

-- ตารางคอมเมนต์
COMMENT ON TABLE card_comments IS 'ตารางเก็บความคิดเห็นของการ์ด';
COMMENT ON COLUMN card_comments.id IS 'รหัสคอมเมนต์ (UUID)';
COMMENT ON COLUMN card_comments.card_id IS 'อ้างอิงไปยังการ์ด';
COMMENT ON COLUMN card_comments.author_id IS 'UUID ของผู้เขียนคอมเมนต์';
COMMENT ON COLUMN card_comments.content IS 'ข้อความความคิดเห็น';
COMMENT ON COLUMN card_comments.created_at IS 'วันและเวลาที่สร้างคอมเมนต์';
COMMENT ON COLUMN card_comments.updated_at IS 'วันและเวลาที่แก้ไขคอมเมนต์ล่าสุด';

-- ตาราง log ความคืบหน้า
COMMENT ON TABLE card_progress_logs IS 'ตารางบันทึกความคืบหน้าและกิจกรรมของการ์ด';
COMMENT ON COLUMN card_progress_logs.id IS 'รหัส log (UUID)';
COMMENT ON COLUMN card_progress_logs.card_id IS 'อ้างอิงไปยังการ์ด';
COMMENT ON COLUMN card_progress_logs.actor_id IS 'UUID ของผู้กระทำ (เช่น ผู้แก้ไข)';
COMMENT ON COLUMN card_progress_logs.message IS 'ข้อความอธิบายการเปลี่ยนแปลงหรือกิจกรรม';
COMMENT ON COLUMN card_progress_logs.created_at IS 'วันและเวลาที่บันทึก log';
