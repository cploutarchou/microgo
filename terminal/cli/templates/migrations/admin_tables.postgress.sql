DO $$ BEGIN IF EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'roles'
) THEN DROP TABLE roles CASCADE;
END IF;
END $$;
CREATE TABLE roles (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    description varchar(255) NOT NULL,
    created_at timestamp DEFAULT NULL,
    updated_at timestamp DEFAULT NULL,
    UNIQUE (name)
);
DO $$ BEGIN IF EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'role_users'
) THEN DROP TABLE role_users CASCADE;
END IF;
END $$;
CREATE TABLE role_users (
    id serial PRIMARY KEY,
    user_id integer NOT NULL,
    role_id integer NOT NULL,
    created_at timestamp DEFAULT NULL,
    updated_at timestamp DEFAULT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
DO $$ BEGIN IF EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'permissions'
) THEN DROP TABLE permissions CASCADE;
END IF;
END $$;
CREATE TABLE permissions (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    description varchar(255) NOT NULL,
    created_at timestamp DEFAULT NULL,
    updated_at timestamp DEFAULT NULL,
    UNIQUE (name)
);
DO $$ BEGIN IF EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'permission_roles'
) THEN DROP TABLE permission_roles CASCADE;
END IF;
END $$;
CREATE TABLE permission_roles (
    id serial PRIMARY KEY,
    permission_id integer NOT NULL,
    role_id integer NOT NULL,
    created_at timestamp DEFAULT NULL,
    updated_at timestamp DEFAULT NULL,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE
);