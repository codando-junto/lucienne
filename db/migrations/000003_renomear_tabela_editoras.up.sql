ALTER TABLE editoras RENAME TO publishers;
ALTER TABLE publishers RENAME COLUMN editora_id TO id;
ALTER TABLE publishers RENAME COLUMN nome TO name;
ALTER TABLE publishers ADD CONSTRAINT publishers_name_key UNIQUE (name);
