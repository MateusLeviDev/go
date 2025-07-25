# poc-gcp-go

# Pub/Sub Example

Este repositório contém dois programas em Go que interagem com o Google Cloud Pub/Sub. A primeira function faz uma requisição GET a um endpoint e publica as mensagens recebidas em um tópico do Pub/Sub. A segunda function consome essas mensagens do tópico e faz um POST em outro endpoint.

## Estrutura do Repositório

- `func1.go`: Programa responsável por publicar mensagens no tópico do Pub/Sub.
- `function2.go`: Programa responsável por consumir mensagens do tópico e fazer POST das mensagens consumidas.

## Requisitos

1. **Go**: Certifique-se de ter o Go instalado em sua máquina. Você pode baixá-lo de [golang.org](https://golang.org/dl/).
2. **Variáveis de Ambiente**: Configure as seguintes variáveis de ambiente:
   - `GCP_PROJECT_ID`: ID do projeto no Google Cloud.
   - `PUBSUB_EMULATOR_HOST`: Host do emulador Pub/Sub, por exemplo, `localhost:8085` se estiver usando o emulador local.
3. **Google Cloud SDK**: Se você estiver usando o emulador Pub/Sub, certifique-se de ter o Google Cloud SDK instalado e o emulador Pub/Sub em execução.
4. **Endpoints da API**: Nessa POC está sendo usado o json-server pras requisições GET e POST. Certifique-se de que os endpoints `http://localhost:3000/posts` e `http://localhost:3000/func2` estejam disponíveis e funcionando conforme esperado.

## Uso

Para ambas functions deve se configurar as variáveis de ambiente:

`export PUBSUB_EMULATOR_HOST=localhost:8085`

`export GCP_PROJECT_ID=project-gcloud-go`


### Executar a function1

1. Navegue até o diretório que contém o arquivo `func1.go`.
2. Execute o comando:

    ```sh
    go run func1.go
    ```

### Executar a function2

1. Navegue até o diretório que contém o arquivo `function2.go`.
2. Execute o comando:

    ```sh
    go run function2.go
    ```

## Detalhes dos Programas

### Publicador (`func1.go`)

1. Faz uma requisição GET ao endpoint `http://localhost:3000/posts`.
2. Publica as mensagens recebidas em um tópico do Google Cloud Pub/Sub (`example-topic3`).

### Consumidor (`function2.go`)

1. Consome mensagens do tópico do Google Cloud Pub/Sub (`example-subscription3`).
2. Realiza um POST das mensagens consumidas para o endpoint `http://localhost:3000/func2`.

## Exemplos de Comandos

```sh
# Executando a function1
cd /path/to/directory
go run func1.go

# Executando a function2
cd /path/to/directory
go run function2.go
