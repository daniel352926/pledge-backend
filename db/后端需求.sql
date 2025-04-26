create table transaction
(
    id      bigint auto_increment primary key,
    tx_hash VARCHAR(66)        not null,
    tx_type VARCHAR(10)     null,
    nonce   VARCHAR(20) null,
    constraint transaction_pk
        unique (tx_hash)
)
    comment '交易表';
CREATE TABLE transaction
(
    id                       BIGINT AUTO_INCREMENT PRIMARY KEY,
    tx_hash                  CHAR(66) NOT NULL,
    tx_type                  tinyint,
    nonce                    bigint,
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
    created_at               TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) engine = InnoDB charset = utf8;

CREATE TABLE receipt
(
    id                  BIGINT AUTO_INCREMENT PRIMARY KEY,
    tx_hash             CHAR(66) NOT NULL,
    tx_type             VARCHAR(10),
    status              VARCHAR(10),
    root                VARCHAR(66),
    cumulative_gas_used VARCHAR(66),
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
    number            VARCHAR(255) NOT NULL,
    hash              VARCHAR(255) NOT NULL,
    parent_hash       VARCHAR(255) NOT NULL,
    nonce             VARCHAR(255) NOT NULL,
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
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) engine = InnoDB charset = utf8;