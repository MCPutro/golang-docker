create table if not exists users
(
    user_id serial not null
    constraint users_pkey
    primary key,
    username varchar(15) not null
    constraint users_username_key
    unique,
    fullname varchar(200),
    password varchar,
    creation_date timestamp with time zone default now() not null
);