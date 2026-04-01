create table if not exists users (
    id UUID PRIMARY KEY,
    email VARCHAR(128) UNIQUE NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    username VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

create index if not exists email_idx ON users (email);

create table if not exists groups (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    title VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

create table if not exists tasks (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    group_id UUID NOT NULL,
    title VARCHAR(64) NOT NULL,
    task_desc VARCHAR(512) NOT NULL,
    priority_id INT DEFAULT 0,
    is_done BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,


    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

create index if not exists priority_idx ON tasks (priority_id);