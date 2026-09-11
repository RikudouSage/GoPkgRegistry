-- +goose up
create table packages (
    id integer primary key autoincrement,
    import_path text not null unique,
    vcs varchar(6) not null,
    repository_url text not null,
    source_url text default null,
    source_dir_url text default null,
    source_file_url text default null
);

create index idx_import_path on packages(import_path);

-- +goose down
drop table packages;