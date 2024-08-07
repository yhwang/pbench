create table if not exists pbench_queries
(
    run_id             bigint       not null,
    stage_id           varchar(255) not null,
    query_file         varchar(255) not null,
    query_index        int          not null,
    sequence_no        int          not null,
    query_id           varchar(255) not null,
    cold_run           tinyint(1)   not null,
    succeeded          tinyint(1)   not null,
    start_time         datetime(3)  not null,
    end_time           datetime(3)  not null,
    row_count          int          null comment 'null if the query failed.',
    expected_row_count int          null,
    duration_ms        int          not null,
    info_url           varchar(255) not null,
    primary key (run_id, stage_id, query_file, query_index, sequence_no)
)
    partition by hash (`run_id`) partitions 16;

-- create index pbench_queries_query_id_index
--    on presto_benchmarks.pbench_queries (query_id);
