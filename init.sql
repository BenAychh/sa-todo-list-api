CREATE USER sa_todo_list WITH
	LOGIN
	NOSUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	NOREPLICATION
	CONNECTION LIMIT -1
	PASSWORD 'local_password';

CREATE DATABASE sa_todo_list
    WITH 
    OWNER = sa_todo_list
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1;

CREATE TABLE public.todos
(
    id serial NOT NULL,
    description character varying NOT NULL,
    complete boolean NOT NULL DEFAULT false,
    PRIMARY KEY (id)
)
