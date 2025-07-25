CREATE TABLE IF NOT EXISTS kafka_ads
(
    _raw_message String
)
    ENGINE = Kafka
    SETTINGS kafka_broker_list = 'kafka:9092',
    kafka_topic_list = 'tracking-events',
    kafka_group_name = 'clickhouse_consumer',
    kafka_format = 'LineAsString',
    kafka_poll_timeout_ms = 10000,   -- wait up to 10s before polling
    kafka_max_block_size = 5000;

-- it's important to have UInt32, because during fraud detection, you need to add -minus clicks,impression or conversions operation here.
-- And you need event_time and event_minute be passed from from bidding server, reason is, Normally a line item can have an expiry date
-- So bidding service will serve campaign till clock hits 23:59, but pushing to kafka, clickhouse ingestion, running materlized view
-- costs time which can generate report bit later in the database which results in slight issue from technical end but from a report or
-- business perspective they will see campaign being served even after they wanted it stop.

CREATE TABLE IF NOT EXISTS ads_final
(
    event_time     DateTime,
    event_minute   Int64,
    item_id        String,
    user_id        String,
    placement      String,
    keyword        String,
    clicks         UInt32,
    impressions    UInt32,
    conversions    UInt32,
    message        String
)
    ENGINE = MergeTree()
        PARTITION BY toStartOfWeek(event_time)
        ORDER BY (event_time);

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_kafka_to_ads
TO ads_final
AS SELECT
        parseDateTimeBestEffort(JSONExtractString(_raw_message, 'event_time')) AS event_time,
        JSONExtractString(_raw_message, 'event_minute')           AS event_minute,
          JSONExtractString(_raw_message, 'item_id')                   AS item_id,
          JSONExtractString(_raw_message, 'user_id')            AS user_id,
          JSONExtractString(_raw_message, 'placement')    AS placement,
          JSONExtractString(_raw_message, 'keyword')      AS keyword,
          JSONExtractUInt(_raw_message, 'clicks')         AS clicks,
          JSONExtractUInt(_raw_message, 'impressions')    AS impressions,
          JSONExtractUInt(_raw_message, 'conversions')    AS conversions,
          _raw_message                                    AS message
FROM kafka_ads;

CREATE TABLE IF NOT EXISTS line_item_report
(
    event_time       DateTime,
    event_minute    Int64,
    item_id         String,
    placement        LowCardinality(String), -- Limited number of placement available, this will reduce storage cost massively
    keyword          LowCardinality(String), -- even if there are 10K keywords still low, compare to billions (ex. 100B * 20 bytes ≈ 2,000 GB (2 TB just for keywords)) of data
    total_clicks     UInt64,
    total_impressions UInt64,
    total_conversions UInt64
)
    ENGINE = SummingMergeTree()
        ORDER BY (event_time);

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_ads_final_to_line_item_report
            TO line_item_report
AS
SELECT
    event_time,
    event_minute,
    placement,
    keyword,
    item_id,
    sum(clicks)      AS total_clicks,
    sum(impressions) AS total_impressions,
    sum(conversions) AS total_conversions
FROM ads_final
GROUP BY event_time,event_minute, placement, keyword, item_id;
