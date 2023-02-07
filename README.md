<h1 align="center">
  gestao-vendas-apirestful-node
</h1>

 <p>Aplicação Back-end para gestão de vendas para restaurante local. Principais tecnologias utilizadas para desenvolvimento da API: Node.js, Express, Typescript, TypeORM, Postgres através de container Docker, Redis através de container Docker, Amazon S3, Amazon SES.
<a href="https://www.notion.so/Gest-o-de-vendas-052676e0389749ebbc79cb405f5d2555">Anotações - Notion</a> 
</p>
 
# `Funcionalidades`
- Criação de cadastro de produtos
- Clientes 
- Completa gestão de users
- Autentição Token JWT 
- Recuperação de email 
- Atualização de perfil e avatar
- O `typeORM` permitirá implenetar entidades e repositórios <br>

![image](https://user-images.githubusercontent.com/101754313/216619865-349e53f1-0fa0-4152-aaa9-d7345c19250d.png)

# `Iniciando projeto`
`yarn init -y`
<br>
em seguida: `yarn add typescript ts-node-dev @types/node tsconfig-paths -D` [instalação dos primeiros pacotes]

```
yarn tsc --init --rootDir src --outDir build --esModuleInterop --resolveJsonModule --lib es6 --module commonjs --allowJs true --noImplicitAny true
```
`rootDir:` É aqui que o TypeScript procura nosso código.

`outDir:` Onde o TypeScript coloca nosso código compilado.

`esModuleInterop:` Se estiver usando commonjs como sistema de módulo (recomendado para aplicativos Node), então esse parâmetro deve ser definido como true.

`resolveJsonModule:` Se usarmos JSON neste projeto, esta opção permite que o TypeScript o use.

`lib:` Esta opção adiciona tipos de ambiente ao nosso projeto, permitindo-nos contar com recursos de diferentes versões do Ecmascript, bibliotecas de teste e até mesmo a API DOM do navegador. Usaremos recursos es6 da linguagem.

`module:` commonjs é o sistema de módulo Node padrão.

`allowJs:` Se você estiver convertendo um projeto JavaScript antigo em TypeScript, esta opção permitirá que você inclua arquivos .js no projeto.

`noImplicitAny:` Em arquivos TypeScript, não permita que um tipo seja especificado inexplicitamente. Cada tipo precisa ter um tipo específico ou ser declarado explicitamente any.

### Trabalhando com ts, precisamos de uma build, ou seja, compilar esse arquivo em um js.

```
yarn tsc
```
- criação de pasta build. pega tudo dentro de src será compilado e jogado no build. 
`node build/server.js`

# `Estrutura do Projeto`

Estrutura de pastas:

`config` - configurações de bibliotecas externas, como por exemplo, autenticação, upload, email, etc.

`modules` - abrangem as áreas de conhecimento da aplicação, diretamente relacionados com as regras de negócios. A princípio criaremos os seguintes módulos na aplicação: customers, products, orders e users.

`shared` - módulos de uso geral compartilhados com mais de um módulo da aplicação, como por exemplo, o arquivo server.ts, o arquivo principal de rotas, conexão com banco de dados, etc.

`services` - estarão dentro de cada módulo da aplicação e serão responsáveis por todas as regras que a aplicação precisa atender, como por exemplo:

- A senha deve ser armazenada com criptografia;
- Não pode haver mais de um produto com o mesmo nome;
- Não pode haver um mesmo email sendo usado por mais de um usuário;
- E muitas outras...

Criando a estrutura de pastas:

```shell
mkdir -p src/config

mkdir -p src/modules

mkdir -p src/shared/http

mv src/server.ts src/shared/http/server.ts
```

Ajustar o arquivo `package.json`:

```json
{
  "scripts": {
    "dev": "ts-node-dev --inspect --transpile-only --ignore-watch node_modules src/shared/http/server.ts"
  }
}
```

### Configurando as importações

Podemos usar um recurso que facilitará o processo de importação de arquivos em nosso projeto.

Iniciamos configurando o objeto `paths` do `tsconfig.json`, que permite criar uma base para cada `path` a ser buscado no projeto, funcionando de forma similar a um atalho:

```json
"paths": {
  "@config/*": ["src/config/*"],
  "@modules/*": ["src/modules/*"],
  "@shared/*": ["src/shared/*"]
}
```

> Nesta videoaula ficou faltando instalar a biblioteca que irá indicar ao nosso script do `ts-node-dev`, como interpretar os atalhos que configuramos iniciando com o caracter `@`.

O nome dessa biblioteca é `tsconfig-paths`, e para instalar execute o seguinte comando no terminal (na pasta do projeto):

```shell
yarn add -D tsconfig-paths
```

Depois de instalar o `tsconfig-paths`, ajustar o script `dev` no arquivo `package.json`, incluindo a opção `-r tsconfig-paths/register`. Deverá ficar assim:

```json
"dev": "ts-node-dev -r tsconfig-paths/register --inspect --transpile-only --ignore-watch node_modules src/shared/http/server.ts"
```

> Observação: o comando acima deve ser incluído como uma linha única do script `dev`.


Agora, para importar qualquer arquivo no projeto, inicie o caminho com um dos `paths` configurados, usando o `CTRL+SPACE` para usar o autocomplete.

> Observação: ativar "baseUrl": "./", para informar que inicia na raiz. 

```
yarn add express cors express-async-errors
```
> express-async-error para trabalharmos com requisições assíncronas em casos de tratamento de erro

### Criação de middleware

É o software que se encontra entre o sistema operacional e os aplicativos nele executados. Essencialmente, o middleware funciona como uma camada oculta de tradução, permitindo a comunicação e o gerenciamento de dados para aplicativos distribuídos. Resposánvel por interceptar o erros gerados pela aplicação.

`src/shared/errors/AppError.ts`

```
app.use(
  (error: Error, request: Request, response: Response, next: NextFunction) => {
    if (error instanceof AppError) {
      return response.status(error.statudCode).json({
        status: 'error',
        message: error.message,
      });
    }
    return response.status(500).json({
      status: 'error',
      message: 'Internal server error',
    });
    //verificar se o erro é uma instância da clase AppError
  },
);
```
> Middleware que recebe um error. Se a instância do erro é da classe, uma vez que será usado AppError no serviço, o erro sendo da aplicação aparecerá a mensagem. Agora caso não seja, provavelmente é um erro de fora, portando, status 500, erro desconhecido. 

# `Nota sobre a versão do TypeORM`

Basta substituir nas dependências do `package.json`
 - "typeorm": "^0.3.x" 
 
 Por: 
 
 - "typeorm": "0.2.29"
 
 Para garantir não haver conflitos, excluir pasta `node_modules` e reinstale usando `yarn` ou `npm install`
 
 # `Docker`
 Comando para criação de um container, quando não tem em execução. 
 
 - `-e` variável de ambiente do container
 - `p` definir em qual porta ele rodará em nossa máquina
 
 ```
 docker run --name postgres -e POSTGRES_PASSWORD=docker -p 5432:5432 -d postgres 
 ```
 

