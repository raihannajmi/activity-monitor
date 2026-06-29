CREATE TABLE IF NOT EXISTS tasks (
    id          TEXT PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT DEFAULT '',
    priority    TEXT CHECK(priority IN ('low','medium','high')) DEFAULT 'medium',
    status      TEXT CHECK(status IN ('todo','in_progress','done')) DEFAULT 'todo',
    deadline    DATETIME,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    column_id   TEXT,
    position    REAL DEFAULT 0.0,
    parent_task_id TEXT,
    estimate_pomodoro INTEGER DEFAULT 0,
    completed_pomodoro INTEGER DEFAULT 0,
    cover_image TEXT,
    created_by  TEXT,
    updated_by  TEXT,
    archived    INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS subtasks (
    id           TEXT PRIMARY KEY,
    task_id      TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    title        TEXT NOT NULL,
    is_completed INTEGER DEFAULT 0,
    deadline     DATETIME,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reminders (
    id         TEXT PRIMARY KEY,
    task_id    TEXT REFERENCES tasks(id) ON DELETE CASCADE,
    remind_at  DATETIME NOT NULL,
    note       TEXT DEFAULT '',
    is_done    INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notes (
    id         TEXT PRIMARY KEY,
    title      TEXT NOT NULL,
    content    TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS activities (
    id           TEXT PRIMARY KEY,
    type         TEXT NOT NULL,
    reference_id TEXT DEFAULT '',
    description  TEXT NOT NULL,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_subtasks_task_id ON subtasks(task_id);
CREATE INDEX IF NOT EXISTS idx_reminders_remind_at ON reminders(remind_at);
CREATE INDEX IF NOT EXISTS idx_activities_created_at ON activities(created_at DESC);

CREATE TABLE IF NOT EXISTS time_logs (
    id               TEXT PRIMARY KEY,
    task_id          TEXT REFERENCES tasks(id) ON DELETE CASCADE,
    subtask_id       TEXT REFERENCES subtasks(id) ON DELETE CASCADE,
    start_time       DATETIME NOT NULL,
    end_time         DATETIME,
    duration_seconds INTEGER DEFAULT 0,
    session_type     TEXT DEFAULT 'pomodoro'
);

CREATE INDEX IF NOT EXISTS idx_time_logs_task_id ON time_logs(task_id);

-- New Trello-style Kanban Tables
CREATE TABLE IF NOT EXISTS boards (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT DEFAULT '',
    color       TEXT DEFAULT '#2563EB',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS columns (
    id          TEXT PRIMARY KEY,
    board_id    TEXT NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    position    REAL NOT NULL,
    color       TEXT DEFAULT '#6B7280',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS checklists (
    id          TEXT PRIMARY KEY,
    task_id     TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    position    REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS checklist_items (
    id           TEXT PRIMARY KEY,
    checklist_id TEXT NOT NULL REFERENCES checklists(id) ON DELETE CASCADE,
    title        TEXT NOT NULL,
    completed    INTEGER DEFAULT 0,
    position     REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS task_labels (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    color       TEXT DEFAULT '#3B82F6'
);

CREATE TABLE IF NOT EXISTS task_label_relations (
    task_id     TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    label_id    TEXT NOT NULL REFERENCES task_labels(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, label_id)
);

CREATE TABLE IF NOT EXISTS attachments (
    id          TEXT PRIMARY KEY,
    task_id     TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    filename    TEXT NOT NULL,
    url         TEXT NOT NULL,
    size        INTEGER NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comments (
    id          TEXT PRIMARY KEY,
    task_id     TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    content     TEXT NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS activity_logs (
    id          TEXT PRIMARY KEY,
    task_id     TEXT REFERENCES tasks(id) ON DELETE CASCADE,
    type        TEXT NOT NULL,
    old_value   TEXT,
    new_value   TEXT,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
