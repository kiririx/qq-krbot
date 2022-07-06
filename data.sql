create table subscribe_subject
(
    id         int auto_increment
        primary key,
    tag        varchar(64) not null unique,
    qq_account text,
    active     bool default false,
    created_at datetime    null,
    updated_at datetime    null,
    deleted_at datetime    null
);

-- auto-generated definition
create table content
(
    id         int auto_increment
        primary key,
    content    text     null,
    created_at datetime null,
    updated_at datetime null,
    deleted_at datetime null
);

