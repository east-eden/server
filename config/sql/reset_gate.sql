drop database if exists `db_gate`;
create database `db_gate` character set utf8mb4;

use `db_gate`;
drop table if exists `gate`;
create table gate (
    `id` int(10) not null default '0' comment 'yokai gate id',
    `time_stamp` int(10) not null default '0' comment 'current time',
    primary key (`id`)
) engine=innodb default charset=utf8mb4 collate utf8mb4_general_ci comment='gate table';

