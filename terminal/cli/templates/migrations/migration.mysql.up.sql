drop table if exists some_table cascade;

CREATE TABLE some_table
(
    `id`         int(10) unsigned                                        NOT NULL AUTO_INCREMENT,
    `some_field` varchar(250) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
    `created_at` timestamp                                               NULL DEFAULT NULL,
    `updated_at` timestamp                                               NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `some_field_unique` (`some_field`),
    KEY `some_field_index` (`some_field`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 17
  DEFAULT CHARSET = utf8mb4;

