ALTER TABLE publishers DROP CONSTRAINT publishers_name_key;
ALTER TABLE publishers RENAME COLUMN name TO nome;
ALTER TABLE publishers RENAME COLUMN id TO editora_id;
ALTER TABLE publishers RENAME TO editoras;
