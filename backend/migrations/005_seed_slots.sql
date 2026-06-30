-- Seed available slots for verified mentors (next 14 days, 19:00 & 20:00)

INSERT INTO mentor_slots (mentor_id, slot_date, start_time, end_time, status)
SELECT m.id, d.dt, t.start_time, t.end_time, 0
FROM mentors m
CROSS JOIN (
    SELECT (CURRENT_DATE + i)::date AS dt
    FROM generate_series(0, 13) AS i
) d
CROSS JOIN (
    VALUES
        ('19:00:00'::time, '19:45:00'::time),
        ('20:00:00'::time, '20:45:00'::time)
) AS t(start_time, end_time)
WHERE m.is_verified = 1
  AND NOT EXISTS (
    SELECT 1 FROM mentor_slots ms
    WHERE ms.mentor_id = m.id
      AND ms.slot_date = d.dt
      AND ms.start_time = t.start_time
  );
