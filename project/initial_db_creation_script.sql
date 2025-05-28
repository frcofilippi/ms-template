create user pedimeapp with password 'mysecretpwd';
create database pedimedb with owner pedimeapp;
grant usage on schema public to pedimeapp;
