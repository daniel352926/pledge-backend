create table transaction
(
    id                       bigint auto_increment primary key,
    tx_hash                  char(66)                            not null,
    tx_type                  tinyint null,
    block_num                bigint unsigned                     null,
    t_index                  bigint unsigned                     null,
    nonce                    bigint null,
    gas_price                varchar(66) null,
    max_priority_fee_per_gas varchar(66) null,
    max_fee_per_gas          varchar(66) null,
    gas_limit                varchar(20) null,
    value                    varchar(66) null,
    input_data               longtext null,
    v                        varchar(10) null,
    r                        char(66) null,
    s                        char(66) null,
    to_address               char(42) null,
    chain_id                 varchar(20) null,
    access_list              json null,
    created_at               timestamp default CURRENT_TIMESTAMP not null
) charset = utf8 comment '交易表';
CREATE TABLE transaction
(
    id                       BIGINT AUTO_INCREMENT PRIMARY KEY,
    tx_hash                  CHAR(66) NOT NULL,
    tx_type                  TINYINT,
    block_num                BIGINT,
    t_index                  BIGINT,
    nonce                    BIGINT,
    gas_price                VARCHAR(66),
    max_priority_fee_per_gas VARCHAR(66),
    max_fee_per_gas          VARCHAR(66),
    gas_limit                VARCHAR(20),
    value                    VARCHAR(66),
    input_data               LONGTEXT,
    v                        VARCHAR(10),
    r                        CHAR(66),
    s                        CHAR(66),
    to_address               CHAR(42),
    chain_id                 VARCHAR(20),
    access_list              JSON,
    created_at               TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY                      `idx_block_num_t_index` (block_num, t_index),
    KEY                      `idx_tx_hash` (tx_hash(32))
) engine = InnoDB charset = utf8;

CREATE TABLE receipt
(
    id                  BIGINT AUTO_INCREMENT PRIMARY KEY,
    tx_hash             CHAR(66) NOT NULL,
    tx_type             INT,
    status              BIGINT UNSIGNED,
    root                VARCHAR(66),
    cumulative_gas_used BIGINT UNSIGNED,
    logs_bloom          TEXT,
    contract_address    CHAR(42),
    gas_used            VARCHAR(66),
    block_hash          CHAR(66),
    block_number        VARCHAR(20),
    transaction_index   VARCHAR(10),
    logs                JSON,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) engine = InnoDB charset = utf8;

CREATE TABLE block
(
    id                BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    number            BIGINT       NOT NULL,
    hash              VARCHAR(255) NOT NULL,
    parent_hash       VARCHAR(255) NOT NULL,
    nonce             BIGINT       NOT NULL,
    sha3_uncles       VARCHAR(255) NOT NULL,
    logs_bloom        TEXT         NOT NULL,
    transactions_root VARCHAR(255) NOT NULL,
    state_root        VARCHAR(255) NOT NULL,
    receipts_root     VARCHAR(255) NOT NULL,
    miner             VARCHAR(255) NOT NULL,
    difficulty        VARCHAR(255) NOT NULL,
    total_difficulty  VARCHAR(255) NOT NULL,
    extra_data        TEXT         NOT NULL,
    size              VARCHAR(255) NOT NULL,
    gas_limit         VARCHAR(255) NOT NULL,
    gas_used          VARCHAR(255) NOT NULL,
    timestamp         VARCHAR(255) NOT NULL,
    transactions      JSON         NOT NULL,
    uncles            JSON         NOT NULL,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY               `idx_number` (number)
) engine = InnoDB charset = utf8;