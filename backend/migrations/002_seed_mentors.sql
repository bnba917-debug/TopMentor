-- Seed verified mentors for M2 学霸广场 demo

INSERT INTO mentors (openid, real_name, school_name, major, gender, english_score, intro_video_url, tags, is_verified)
SELECT 'seed_mentor_1', '张明', '清华大学', '计算机科学', 'male', '高考英语 148 分', 'https://example.com/videos/mentor1.mp4', ARRAY['阳光幽默', '善于引导'], 1
WHERE NOT EXISTS (SELECT 1 FROM mentors WHERE openid = 'seed_mentor_1');

INSERT INTO mentors (openid, real_name, school_name, major, gender, english_score, intro_video_url, tags, is_verified)
SELECT 'seed_mentor_2', '李雨桐', '北京大学', '英语语言文学', 'female', '托福 115', 'https://example.com/videos/mentor2.mp4', ARRAY['耐心细致', '口语对练'], 1
WHERE NOT EXISTS (SELECT 1 FROM mentors WHERE openid = 'seed_mentor_2');

INSERT INTO mentors (openid, real_name, school_name, major, gender, english_score, intro_video_url, tags, is_verified)
SELECT 'seed_mentor_3', '王浩然', '清华大学', '电子工程', 'male', '雅思 8.0', 'https://example.com/videos/mentor3.mp4', ARRAY['理工科背景', '阳光幽默'], 1
WHERE NOT EXISTS (SELECT 1 FROM mentors WHERE openid = 'seed_mentor_3');

INSERT INTO mentors (openid, real_name, school_name, major, gender, english_score, intro_video_url, tags, is_verified)
SELECT 'seed_mentor_pending', '待审核同学', '复旦大学', '经济学', 'female', '高考英语 140 分', '', ARRAY['待审核'], 0
WHERE NOT EXISTS (SELECT 1 FROM mentors WHERE openid = 'seed_mentor_pending');
