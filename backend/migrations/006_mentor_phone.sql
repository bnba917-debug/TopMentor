-- Link mentors to phone for H5 学霸端 login
ALTER TABLE mentors ADD COLUMN IF NOT EXISTS phone VARCHAR(20);

CREATE UNIQUE INDEX IF NOT EXISTS idx_mentors_phone ON mentors(phone) WHERE phone IS NOT NULL;

UPDATE mentors SET phone = '13800000001' WHERE openid = 'seed_mentor_1' AND phone IS NULL;
UPDATE mentors SET phone = '13800000002' WHERE openid = 'seed_mentor_2' AND phone IS NULL;
UPDATE mentors SET phone = '13800000003' WHERE openid = 'seed_mentor_3' AND phone IS NULL;
