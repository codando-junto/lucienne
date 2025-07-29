RENAME TABLE editoras TO publishers;
RENAME COLUMN publishers.editora_id TO publishers.id;
RENAME COLUMN publishers.nome TO publishers.name;
