create table auto_bid_history (
    task_id    int      not null,
    miner_id   int      not null,
    create_at  bigint   not null,
    constraint pk_auto_bid_history primary key (create_at,miner_id,task_id)
) engine = innodb
partition by range (create_at) (
    partition part_max values less than maxvalue
)