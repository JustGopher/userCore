create database userCore;
use userCore;
create table if not exists roles (
    role_id int auto_increment primary key,
    name varchar(10) not null,
    description varchar(100)
    );
create table if not exists user(
    user_id int primary key auto_increment,
    user_name varchar(20) not null unique ,
    password varchar(20) not null ,
    role_id int ,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,  -- 更新时间
    foreign key (role_id) references roles(role_id)
    );

insert into roles(name, description)
values
    ("普通用户","仅可查看，不可修改"),
    ("管理员","可查看和修改");

insert into user(user_name,password,role_id)
values ("admin","123456",2);

insert into user(user_name,password,role_id)
values ("xiaoming","123456",1);