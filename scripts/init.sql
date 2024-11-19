-- 创建数据库
CREATE DATABASE IF NOT EXISTS shopee_tool
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE shopee_tool;

-- 账号表
CREATE TABLE IF NOT EXISTS accounts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) NOT NULL COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-正常 2-禁用',
    last_login_at DATETIME COMMENT '最后登录时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_username (username),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账号表';

-- Cookie表
CREATE TABLE IF NOT EXISTS cookies (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    account_id BIGINT UNSIGNED NOT NULL COMMENT '关联的账号ID',
    name VARCHAR(255) NOT NULL COMMENT 'cookie名称',
    value TEXT NOT NULL COMMENT 'cookie值',
    domain VARCHAR(255) NOT NULL COMMENT 'cookie域名',
    path VARCHAR(255) NOT NULL COMMENT 'cookie路径',
    expires DATETIME COMMENT 'cookie过期时间',
    http_only TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否httpOnly',
    secure TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否secure',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_account_id (account_id),
    KEY idx_expires (expires),
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Cookie表';

-- 登录日志表
CREATE TABLE IF NOT EXISTS login_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    account_id BIGINT UNSIGNED NOT NULL COMMENT '账号ID',
    login_ip VARCHAR(64) COMMENT '登录IP',
    user_agent VARCHAR(255) COMMENT '用户代理',
    status TINYINT NOT NULL COMMENT '登录状态：1-成功 2-失败',
    error_msg VARCHAR(255) COMMENT '错误信息',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_account_id (account_id),
    KEY idx_created_at (created_at),
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';

-- API请求日志表
CREATE TABLE IF NOT EXISTS api_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    account_id BIGINT UNSIGNED COMMENT '账号ID',
    request_url VARCHAR(255) NOT NULL COMMENT '请求URL',
    request_method VARCHAR(10) NOT NULL COMMENT '请求方法',
    request_params TEXT COMMENT '请求参数',
    response_code INT COMMENT '响应状态码',
    response_body TEXT COMMENT '响应内容',
    error_msg VARCHAR(255) COMMENT '错误信息',
    request_time DATETIME NOT NULL COMMENT '请求时间',
    response_time DATETIME COMMENT '响应时间',
    duration INT COMMENT '请求耗时(毫秒)',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_account_id (account_id),
    KEY idx_request_time (request_time),
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API请求日志表';

-- 系统配置表
CREATE TABLE IF NOT EXISTS configs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `key` VARCHAR(64) NOT NULL COMMENT '配置键',
    value TEXT NOT NULL COMMENT '配置值',
    description VARCHAR(255) COMMENT '配置描述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_key (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 插入一些初始配置数据
INSERT INTO configs (`key`, value, description) VALUES
    ('retry_times', '3', '接口重试次数'),
    ('retry_interval', '5', '重试间隔(秒)'),
    ('cookie_expire_days', '30', 'Cookie有效期(天)'),
    ('api_timeout', '30', 'API超时时间(秒)');

-- 创建管理员账号（密码需要加密存储）
INSERT INTO accounts (username, password, status) VALUES
    ('admin', '$2a$10$your_hashed_password', 1); 


-- 创建 active_codes 表
CREATE TABLE active_codes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);