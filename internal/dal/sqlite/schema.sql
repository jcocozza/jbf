create table if not exists metadata (
    id integer primary key,
    filepath text not null,
    title text,
    author text,
    created datetime,
    last_updated datetime,
    is_home boolean
);

create table if not exists tag (
    tag_name text primary key,
    metadata_id integer,

    foreign key (metadata_id) references metadata (id)
);


create trigger if not exists delete_metadata
after delete on metadata
for each row
begin
delete from tag
where metadata_id = old.id;
end;
