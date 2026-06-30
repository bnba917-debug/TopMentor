-- Seed pending mentor application for admin review demo
INSERT INTO mentor_applications (mentor_id, id_card_url, student_card_url, english_proof_url)
SELECT m.id,
       'https://example.com/docs/id-card-pending.jpg',
       'https://example.com/docs/student-card-pending.jpg',
       'https://example.com/docs/english-proof-pending.jpg'
FROM mentors m
WHERE m.openid = 'seed_mentor_pending'
  AND NOT EXISTS (
    SELECT 1 FROM mentor_applications ma WHERE ma.mentor_id = m.id AND ma.reviewed_at IS NULL
  );

INSERT INTO courseware (title, cover_url, content_url, sort_order, is_active)
SELECT 'Hello Zoo 绘本', 'https://example.com/covers/zoo.png', 'https://example.com/courseware/zoo.pdf', 1, TRUE
WHERE NOT EXISTS (SELECT 1 FROM courseware WHERE title = 'Hello Zoo 绘本');
