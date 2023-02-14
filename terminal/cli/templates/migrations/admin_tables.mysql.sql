DROP TABLE IF EXISTS `role_users`,
`permission_roles` CASCADE;
DROP TABLE IF EXISTS `roles`,
`permissions`;


CREATE TABLE `roles` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `description` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NULL DEFAULT NULL,
    `updated_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `roles_name_unique` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT = 17 DEFAULT CHARSET = utf8mb4;


CREATE TABLE `role_users` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` INT(10) UNSIGNED NOT NULL,
    `role_id` INT(10) UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NULL DEFAULT NULL,
    `updated_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `role_users_user_id_foreign` (`user_id`),
    KEY `role_users_role_id_foreign` (`role_id`),
    CONSTRAINT `role_users_role_id_foreign` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `role_users_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 17 DEFAULT CHARSET = utf8mb4;


CREATE TABLE `permissions` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `description` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NULL DEFAULT NULL,
    `updated_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `permissions_name_unique` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT = 17 DEFAULT CHARSET = utf8mb4;


CREATE TABLE permission_roles (
    id int(10) unsigned NOT NULL AUTO_INCREMENT,
    permission_id int(10) unsigned NOT NULL,
    role_id int(10) unsigned NOT NULL,
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL,
    PRIMARY KEY (id),
    KEY permission_roles_permission_id_foreign (permission_id),
    KEY permission_roles_role_id_foreign (role_id),
    CONSTRAINT permission_roles_permission_id_foreign FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT permission_roles_role_id_foreign FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 17 DEFAULT CHARSET = utf8mb4;

-- create cms roles and permissions for admin and user roles and assign them to the admin and user users respectively (admin and user) 
-- admin user has all permissions
-- user user has only the permission to view the dashboard
-- subcriber user has only the permission to view the dashboard

INSERT INTO TABLE roles (name, description, created_at, updated_at) VALUES ('admin', 'Administrator', NOW(), NOW());
INSERT INTO TABLE roles (name, description, created_at, updated_at) VALUES ('user', 'User', NOW(), NOW());
INSERT INTO TABLE roles (name, description, created_at, updated_at) VALUES ('subscriber', 'Subscriber', NOW(), NOW());

INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('view_dashboard', 'View Dashboard', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('view_admin', 'View Admin', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('view_users', 'View Users', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('create_users', 'Create Users', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('edit_users', 'Edit Users', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('delete_users', 'Delete Users', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('view_roles', 'View Roles', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('create_roles', 'Create Roles', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('edit_roles', 'Edit Roles', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('delete_roles', 'Delete Roles', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('view_permissions', 'View Permissions', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('create_permissions', 'Create Permissions', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('edit_permissions', 'Edit Permissions', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('delete_permissions', 'Delete Permissions', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('view_posts', 'View Posts', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('create_posts', 'Create Posts', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('edit_posts', 'Edit Posts', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('delete_posts', 'Delete Posts', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('view_categories', 'View Categories', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('create_categories', 'Create Categories', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('edit_categories', 'Edit Categories', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('delete_categories', 'Delete Categories', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('view_tags', 'View Tags', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('create_tags', 'Create Tags', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('edit_tags', 'Edit Tags', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('delete_tags', 'Delete Tags', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('view_comments', 'View Comments', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('create_comments', 'Create Comments', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('edit_comments', 'Edit Comments', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('delete_comments', 'Delete Comments', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('view_pages', 'View Pages', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('create_pages', 'Create Pages', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('edit_pages', 'Edit Pages', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('delete_pages', 'Delete Pages', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('view_media', 'View Media', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('create_media', 'Create Media', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('edit_media', 'Edit Media', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('delete_media', 'Delete Media', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('view_settings', 'View Settings', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('create_settings', 'Create Settings', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('edit_settings', 'Edit Settings', NOW(), NOW());
INSERT INTO TABLE permissions (name, description, created_at, updated_at) VALUES ('delete_settings', 'Delete Settings', NOW(), NOW());