create table user(
	id int primary key auto_increment,
	name varchar(255) not null default '',
	password varchar(18) not null default '', 
	phone varchar(11) default'',
	mail varchar(50) default '',
	count int not null default 0
	--用户登录次数，用于统计
);

create table role(
	id int primary key auto_increment,
	name varchar(255) not null default ''
)

create table permission(
	id int primary key auto_increment,
	parentId int default 0,
	name  varchar(255) not null default '',
	url varchar(1024) not null default ''
)
