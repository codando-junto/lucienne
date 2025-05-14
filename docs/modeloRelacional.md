```mermaid
erDiagram
    USUARIOS {
        integer id
        varchar(100) nome
        varchar(255) email
        varchar(100) senha
    }

    LIVROS {
        integer id
        varchar(30) isbn
        varchar(100) nome
        integer edicao
        integer reimpressao
        integer preco_em_centavos
        datetime data_de_lancamento
    }

    CATEGORIAS {
        id integer
        nome varchar(100)
    }

    AUTORES {
        id integer
        nome varchar(100)
    }

    EDITORAS {
        id integer
        nome varchar(100)
    }
    
    PEDIDOS {
        id integer
        numero integer

    }

    LIVROS }o--|| AUTORES : ""

    LIVROS }o--|| EDITORAS : ""

    LIVROS }o--|| CATEGORIAS : ""

```