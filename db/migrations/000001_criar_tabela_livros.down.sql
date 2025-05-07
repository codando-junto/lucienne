DROP TABLE IF EXISTS livros(
    livro_id serial PRIMARY KEY,
    titulo VARCHAR (100) NOT NULL,
    autor VARCHAR (50) NOT NULL,
    genero VARCHAR (50) NOT NULL,
    anopublicacao INT NOT NULL
);