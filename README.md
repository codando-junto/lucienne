## Pré-requisitos

Certifique-se de ter o seguinte instalado em sua máquina:

- [Go](https://golang.org/dl/) (versão 1.18 ou superior)
- Git

## Como executar o projeto localmente

Siga os passos abaixo para rodar o projeto localmente:

```bash
docker compose up --build
```

Se tudo estiver configurado corretamente, você verá a seguinte mensagem no terminal:

```bash
go_app_container    | Servidor rodando na porta 9090
```

### 3.1. Usando ngrok para expor a porta na internet

O ngrok é uma ferramenta que cria túneis seguros para expor localmente servidores ou aplicações à internet, permitindo acesso remoto por meio de URLs públicas.

Com o ngrok você pode expor a porta 8080 usando seu código da sua máquina e ele pode ser acessado por outra pessoa, pela internet.

**Observação importante**: Lembre-se que a internet é um local inseguro, e expor suas portas por muito tempo sem a devida segurança é um risco muito alto. Esse aplicativo deve ser usado somente para testes simples e por pouco período.

Você precisa se [cadastrar no ngrok](https://dashboard.ngrok.com/signup) e depois [instalar ele](https://dashboard.ngrok.com/signup) na sua estação de trabalho, que será onde vc vai rodar ele.

Você pode criar um domínio fixo pra você, para usar o mesmo domínio sempre [nesse link](https://dashboard.ngrok.com/domains) e nesse mesmo link você já pode pegar o link para iniciar o ngrok. Com meu domínio é assim:

```bash
ngrok http 9090
```

Com o comando acima o ngrok iniciará um tunel mandando tudo que ele receber nesse domínio na porta padrão http (80) e mandará tudo para sua porta 8080 na localhost.

### 4. Teste as rotas
Você pode testar as rotas disponíveis usando um navegador, curl ou ferramentas como Postman.

### Rota /health
Descrição: Retorna o status de saúde da aplicação.

Exemplo de resposta:
Código de status: 200 OK
Corpo da resposta: vazio
