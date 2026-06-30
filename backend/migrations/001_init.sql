-- TopMentor initial schema (M1)

CREATE TABLE IF NOT EXISTS users (
    id              SERIAL PRIMARY KEY,
    openid          VARCHAR(64) UNIQUE NOT NULL,
    phone           VARCHAR(20) NOT NULL,
    child_name      VARCHAR(50),
    child_age       INT DEFAULT 6 CHECK (child_age BETWEEN 6 AND 14),
    english_level   VARCHAR(20) DEFAULT 'beginner',
    available_lessons INT DEFAULT 0 CHECK (available_lessons >= 0),
    locked_lessons  INT DEFAULT 0 CHECK (locked_lessons >= 0),
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_users_openid ON users(openid);

CREATE TABLE IF NOT EXISTS mentors (
    id              SERIAL PRIMARY KEY,
    openid          VARCHAR(64) UNIQUE NOT NULL,
    real_name       VARCHAR(50) NOT NULL,
    school_name     VARCHAR(100) NOT NULL,
    major           VARCHAR(100) NOT NULL,
    gender          VARCHAR(10) DEFAULT 'unknown',
    english_score   VARCHAR(100),
    intro_video_url VARCHAR(512),
    tags            TEXT[] DEFAULT '{}',
    is_verified     SMALLINT DEFAULT 0,
    balance         NUMERIC(10, 2) DEFAULT 0.00 CHECK (balance >= 0),
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_mentors_openid ON mentors(openid);
CREATE INDEX IF NOT EXISTS idx_mentors_verified ON mentors(is_verified);

CREATE TABLE IF NOT EXISTS mentor_applications (
    id              SERIAL PRIMARY KEY,
    mentor_id       INT NOT NULL REFERENCES mentors(id),
    id_card_url     VARCHAR(512) NOT NULL,
    student_card_url VARCHAR(512) NOT NULL,
    english_proof_url VARCHAR(512),
    reject_reason   TEXT,
    reviewed_by     INT,
    reviewed_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS mentor_slots (
    id              SERIAL PRIMARY KEY,
    mentor_id       INT NOT NULL REFERENCES mentors(id),
    slot_date       DATE NOT NULL,
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    status          SMALLINT DEFAULT 0,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (mentor_id, slot_date, start_time)
);
CREATE INDEX IF NOT EXISTS idx_mentor_slots_query ON mentor_slots(mentor_id, slot_date, status);

CREATE TABLE IF NOT EXISTS course_orders (
    id                  VARCHAR(64) PRIMARY KEY,
    user_id             INT NOT NULL REFERENCES users(id),
    mentor_id           INT NOT NULL REFERENCES mentors(id),
    slot_id             INT NOT NULL REFERENCES mentor_slots(id),
    status              VARCHAR(20) DEFAULT 'PENDING',
    agora_channel_name  VARCHAR(128),
    actual_minutes      INT DEFAULT 0,
    mentor_feedback     TEXT,
    feedback_submitted_at TIMESTAMPTZ,
    cancelled_by        VARCHAR(10),
    cancel_reason       TEXT,
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_course_orders_user ON course_orders(user_id, status);
CREATE INDEX IF NOT EXISTS idx_course_orders_mentor ON course_orders(mentor_id, status);

CREATE TABLE IF NOT EXISTS growth_reports (
    id              SERIAL PRIMARY KEY,
    order_id        VARCHAR(64) NOT NULL UNIQUE REFERENCES course_orders(id),
    mentor_id       INT NOT NULL REFERENCES mentors(id),
    user_id         INT NOT NULL REFERENCES users(id),
    speaking_score  SMALLINT CHECK (speaking_score BETWEEN 1 AND 5),
    confidence_score SMALLINT CHECK (confidence_score BETWEEN 1 AND 5),
    vocabulary_score SMALLINT CHECK (vocabulary_score BETWEEN 1 AND 5),
    comment         TEXT NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS lesson_packages (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(50) NOT NULL,
    lesson_count    INT NOT NULL CHECK (lesson_count > 0),
    price_cents     INT NOT NULL CHECK (price_cents > 0),
    is_trial        BOOLEAN DEFAULT FALSE,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS recharge_orders (
    id              VARCHAR(64) PRIMARY KEY,
    user_id         INT NOT NULL REFERENCES users(id),
    package_id      INT NOT NULL REFERENCES lesson_packages(id),
    amount_cents    INT NOT NULL,
    wx_transaction_id VARCHAR(64),
    status          VARCHAR(20) DEFAULT 'PENDING',
    paid_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_recharge_orders_user ON recharge_orders(user_id);

CREATE TABLE IF NOT EXISTS wallet_transactions (
    id              SERIAL PRIMARY KEY,
    mentor_id       INT NOT NULL REFERENCES mentors(id),
    order_id        VARCHAR(64) REFERENCES course_orders(id),
    amount          NUMERIC(10, 2) NOT NULL,
    type            VARCHAR(20) NOT NULL,
    balance_after   NUMERIC(10, 2) NOT NULL,
    remark          TEXT,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_wallet_tx_mentor ON wallet_transactions(mentor_id, created_at DESC);

CREATE TABLE IF NOT EXISTS courseware (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(100) NOT NULL,
    cover_url       VARCHAR(512),
    content_url     VARCHAR(512) NOT NULL,
    sort_order      INT DEFAULT 0,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO lesson_packages (name, lesson_count, price_cents, is_trial)
SELECT '体验课 1 节', 1, 9900, TRUE
WHERE NOT EXISTS (SELECT 1 FROM lesson_packages WHERE name = '体验课 1 节');

INSERT INTO lesson_packages (name, lesson_count, price_cents, is_trial)
SELECT '标准包 10 节', 10, 89900, FALSE
WHERE NOT EXISTS (SELECT 1 FROM lesson_packages WHERE name = '标准包 10 节');

INSERT INTO lesson_packages (name, lesson_count, price_cents, is_trial)
SELECT '进阶包 30 节', 30, 249900, FALSE
WHERE NOT EXISTS (SELECT 1 FROM lesson_packages WHERE name = '进阶包 30 节');
