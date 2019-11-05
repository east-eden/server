drop database if exists `db_battle`;
create database `db_battle` character set utf8mb4;

use `db_battle`;
drop table if exists `battle`;
create table battle (
    `id` int(10) not null default '0' comment 'yokai battle id',
    `time_stamp` int(10) not null default '0' comment 'current time',
    primary key (`id`)
) engine=innodb default charset=utf8mb4 collate utf8mb4_general_ci comment='battle table';

