CREATE DATABASE sa_todo_list
    WITH
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1;

CREATE TABLE public.todos
(
    id serial NOT NULL,
    description character varying NOT NULL,
    complete boolean NOT NULL DEFAULT false,
    PRIMARY KEY (id)
)
