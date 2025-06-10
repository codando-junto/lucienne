```mermaid
erDiagram
    USUARIOS {
        id INTEGER "PK"
        email VARCHAR(255) "UNIQUE, NOT NULL"
        nome VARCHAR(100) "NOT NULL"
        senha VARCHAR(100) "NOT NULL"
    }

    CATEGORIAS {
        id INTEGER "PK"
        nome VARCHAR(100) "NOT NULL"
    }

    AUTORES {
        id INTEGER "PK"
        nome VARCHAR(100) "NOT NULL"
    }

    EDITORAS {
        id INTEGER "PK"
        nome VARCHAR(100) "NOT NULL"
    }

    LIVROS {
        id INTEGER "PK"
        isbn VARCHAR(30) "UNIQUE, NOT NULL"
        nome VARCHAR(100) "NOT NULL"
        edicao INTEGER "NOT NULL"
        reimpressao INTEGER
        preco_em_centavos INTEGER "NOT NULL"
        data_de_lancamento DATETIME "NOT NULL"
        categoria_id INTEGER "FK, NOT NULL"
        autor_id INTEGER "FK, NOT NULL"
        editora_id INTEGER "FK, NOT NULL"
    }
    
    PEDIDOS {
        id INTEGER "PK"
        numero INTEGER "UNIQUE, NOT NULL"
        status VARCHAR(20) "NOT NULL"
        link_de_rastreio VARCHAR(100) "NOT NULL"
        user_id INTEGER "FK, NOT NULL"
        
    }

    ITENS_DE_PEDIDO {
        id INTEGER "PK"
        subtotal_em_centavos INTEGER "NOT NULL"
        desconto_em_centavos INTEGER "DEFAULT 0"
        livro_id INTEGER "FK, NOT NULL"
        pedido_id INTEGER "FK, NOT NULL"
    }

    PAGAMENTOS {
        id INTEGER "PK"
        pedido_id INTEGER "FK, NOT NULL"
    }

    LIVROS }o--|| AUTORES : ""

    LIVROS }o--|| EDITORAS : ""

    LIVROS }o--|| CATEGORIAS : ""

    PEDIDOS }o--|| USUARIOS : ""

    ITENS_DE_PEDIDO }o--|| PEDIDOS : ""

    ITENS_DE_PEDIDO }|--|| LIVROS : ""

    PAGAMENTOS }o--|| PEDIDOS : ""

```
