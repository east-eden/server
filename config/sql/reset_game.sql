drop database if exists `db_game`;
create database `db_game` character set utf8mb4;

use `db_game`;
drop table if exists `global`;
create table global (
    `id` int(10) not null default '0' comment 'yokai game id',
    `time_stamp` int(10) not null default '0' comment 'current time',
    primary key (`id`)
) engine=innodb default charset=utf8mb4 collate utf8mb4_general_ci comment='global table';

drop table if exists `client`;
create table client (
    `id` int(10) not null default '0' comment 'client id',
    `name` varchar(32) not null default '' comment 'client name',
    `last_connect_time` int(10) not null default '0' comment 'last connect time stamp',
    primary key (`id`)
) engine=innodb default charset=utf8mb4 collate utf8mb4_general_ci comment='client';


