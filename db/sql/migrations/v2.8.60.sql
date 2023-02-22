create table if not exists semaphore.project__playbook
(
    id          int auto_increment
        primary key,
    project_id  int          not null,
    name        varchar(255) not null,
    description varchar(255) null,
    content     text         not null,
    constraint project__playbook_id_uindex
        unique (id),
    constraint project__playbook_project_id_fk
        foreign key (project_id) references semaphore.project (id)
            on delete cascade
);

create index project__playbook__index
    on semaphore.project__playbook (project_id);

