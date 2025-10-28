CREATE TABLE IF NOT EXISTS lectures (
    id INTEGER PRIMARY KEY,
    university TEXT NOT NULL,
    title TEXT NOT NULL,
    english_title TEXT,
    department TEXT,
    lecture_type TEXT,
    code TEXT,
    level INTEGER,
    credit INTEGER,
    year INTEGER,
    language TEXT,
    url TEXT,
    abstract TEXT,
    goal TEXT,
    experience TEXT,
    flow TEXT,
    out_of_class_work TEXT,
    textbook TEXT,
    reference_book TEXT,
    assessment TEXT,
    prerequisite TEXT,
    contact TEXT,
    office_hours TEXT,
    note TEXT
);

CREATE TABLE IF NOT EXISTS teachers (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT
);

CREATE TABLE IF NOT EXISTS lecture_teachers (
    lecture_id INTEGER NOT NULL,
    teacher_id INTEGER NOT NULL,
    PRIMARY KEY (lecture_id, teacher_id),
    FOREIGN KEY (lecture_id) REFERENCES lectures(id) ON DELETE CASCADE,
    FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rooms (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS timetables (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    lecture_id INTEGER NOT NULL,
    semester TEXT,
    room_id INTEGER,
    day_of_week TEXT,
    period INTEGER,
    FOREIGN KEY (lecture_id) REFERENCES lectures(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS lecture_plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    lecture_id INTEGER NOT NULL,
    count INTEGER,
    plan TEXT,
    assignment TEXT,
    FOREIGN KEY (lecture_id) REFERENCES lectures(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS lecture_keywords (
    lecture_id INTEGER NOT NULL,
    keyword TEXT NOT NULL,
    PRIMARY KEY (lecture_id, keyword),
    FOREIGN KEY (lecture_id) REFERENCES lectures(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS related_courses (
    lecture_id INTEGER NOT NULL,
    related_lecture_id INTEGER NOT NULL,
    PRIMARY KEY (lecture_id, related_lecture_id),
    FOREIGN KEY (lecture_id) REFERENCES lectures(id) ON DELETE CASCADE,
    FOREIGN KEY (related_lecture_id) REFERENCES lectures(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS related_course_codes (
    lecture_id INTEGER NOT NULL,
    code TEXT NOT NULL,
    PRIMARY KEY (lecture_id, code),
    FOREIGN KEY (lecture_id) REFERENCES lectures(id) ON DELETE CASCADE
);
