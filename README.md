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
- Autenticação Token JWT 
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
 - `-d` roda o processo em background e libera o terminal
 
 ```
 docker run --name postgres -e POSTGRES_PASSWORD=docker -p 5432:5432 -d postgres 
 ```

Na raiz do projeto, deverá ser criada um arquivo `ormconfig.json` para informar os parâmetros necessários para se conectar com o postgres 

```
{
  "type": "postgres",
  "host": "localhost",
  "port": 5432,
  "username": "postgres",
  "password": "docker",
  "database": "apivendas"
}

```
# `ESTRUTURA TABELA PRODUCTS [MIGRATION DA TABELA PRODUCT]`

```
export class CreateProducts1675810870700 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.createTable(
      new Table({
        name: 'products',
        columns: [
          {
            name: 'id',
            type: 'uuid',
            isPrimary: true,
            generationStrategy: 'uuid',
            default: 'uuid_generate_v4()',
          },
          {
            name: 'name',
            type: 'varchar',
          },
          {
            name: 'price',
            type: 'decimal',
            precision: 10,
            scale: 2,
          },
          {
            name: 'quantity',
            type: 'int',
          },
          {
            name: 'created_at',
            type: 'timestamp with time zone',
            default: 'now()',
          },
          {
            name: 'updated_at',
            type: 'timestamp with time zone',
            default: 'now()',
          },
        ],
      }),
    );
  }
```
- o que a migration faz quando rodamos o comando migration run 
- <a href="https://typeorm.delightful.studio/interfaces/_query_runner_queryrunner_.queryrunner.html">Documentação QueryRunner - TypeORM</a> 

# `ENTIDADE DOS PRODUTOS [CONCEITO DE ENTITIES TYPEORM]`
<a href="https://orkhan.gitbook.io/typeorm/docs/entities">Documentação Entities - TypeORM</a> 
Classe que mapeia para uma tabela de banco de dados (ou coleção ao usar MongoDB). Utiliza padrão de decorator, com esse padrão pode-se att funcionalidades extras, sem precisar ter um código muito verboso.

- e an entity by defining a new class and mark it with @Entity(): 

```
import { Entity, PrimaryGeneratedColumn, Column } from "typeorm"

@Entity()
export class User {
    @PrimaryGeneratedColumn()
    id: number

    @Column()
    firstName: string

    @Column()
    lastName: string

    @Column()
    isActive: boolean
}
```

## `Entidade de Produtos`

`mkdir -p src/modules/products/typeorm/entities`

- Dentro de entities criar arq Product.ts. Lembrando que é no singular, visto que é a de UM produto. 

- Definindo a classe `Product` que vai ser nossa entidade em si
- Para ser considerada uma entity no typeorm é preciso usar o decorator `Entity` passando o nome da tabela em que ela fará o mapeamento.
- Em seguida será defindo os att da tabela
- Na próxima operação, é preciso uma forma de dizer ao type orm qual a config de cada coluna att definida acima, dizendo tb o tipo de informação gerada automáticamente

```
import {
  Column,
  CreateDateColumn,
  Entity,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from 'typeorm';

@Entity('products')
class Product {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  name: string;

  @Column('decimal')
  price: number;

  @Column('int')
  quantity: number;

  @CreateDateColumn()
  created_at: Date;

  @UpdateDateColumn()
  updated_at: Date;
}

export default Product;
```
## `Conceito de Repositório do TypeORM`
- <a href="https://typeorm.io/working-with-repository">What is repository - TypeORM</a> 
Repository is just like EntityManager but its operations are limited to a concrete entity. You can access the repository via EntityManager.

### `Repo de Products`
 ```
 import { EntityRepository, Repository } from 'typeorm';
import Product from '../entities/Product';

@EntityRepository(Product)
export class ProductRepository extends Repository<Product> {
  public async findByName(name: string): Promise<Product | undefined> {
    const product = this.findOne({
      where: {
        name,
      },
    });

    return product;
  }
}
 ```
 
 ## `Tratamento de Requisições`
 - CreateProductService: serviço de criação de produto. 
<p>
O service será uma classe, que vai possuir uma instância do repo e através disso será manipulado os dados
</p>

```
import AppError from '@shared/errors/AppError';
import { getCustomRepository } from 'typeorm';
import Product from '../typeorm/entities/Product';
import { ProductRepository } from '../typeorm/repositories/ProductsRepository';

interface IRequest {
  name: string;
  price: number;
  quantity: number;
}

class CreateProductService {
  public async execute({ name, price, quantity }: IRequest): Promise<Product> {
    const productsRepository = getCustomRepository(ProductRepository);
    const productExist = await productsRepository.findByName(name);

    if (productExist) {
      throw new AppError('There is already onde product with this name');
    }

    const product = productsRepository.create({
      name,
      price,
      quantity,
    });

    await productsRepository.save(product);

    return product;
  }
}

export default CreateProductService;
```
> Criação do produto, se preocupando em entender as regras de negócio da aplicação. Como mostra, uma regra que não permite o cadastro de um produto com um nome já existente. 

 - ListProductsService 
<p>
Listagem de produtos cadastrados na DB
</p>

```
import { getCustomRepository } from 'typeorm';
import Product from '../typeorm/entities/Product';
import { ProductRepository } from '../typeorm/repositories/ProductsRepository';

class ListProductService {
  public async execute(): Promise<Product[]> {
    const productsRepository = getCustomRepository(ProductRepository);

    const products = productsRepository.find();

    return products;
  }
}

export default ListProductService;
```
- ShowProductsService 
<p>
  listando um product específico
</p>

```
import AppError from '@shared/errors/AppError';
import { getCustomRepository } from 'typeorm';
import Product from '../typeorm/entities/Product';
import { ProductRepository } from '../typeorm/repositories/ProductsRepository';

interface IRequest {
  id: string;
}
class ShowProductService {
  public async execute({ id }: IRequest): Promise<Product | undefined> {
    const productsRepository = getCustomRepository(ProductRepository);

    const product = productsRepository.findOne(id);

    if (!product) {
      throw new AppError('Product not found.');
    }
    return product;
  }
}

export default ShowProductService;
```
- UpdateProductService 
<p>
  update
</p>

```
import AppError from '@shared/errors/AppError';
import { getCustomRepository } from 'typeorm';
import Product from '../typeorm/entities/Product';
import { ProductRepository } from '../typeorm/repositories/ProductsRepository';

interface IRequest {
  id: string;
  name: string;
  price: number;
  quantity: number;
}
class UpdateProductService {
  public async execute({
    id,
    name,
    price,
    quantity,
  }: IRequest): Promise<Product> {
    const productsRepository = getCustomRepository(ProductRepository);

    const product = await productsRepository.findOne(id);

    if (!product) {
      throw new AppError('Product not found.');
    }

    const productExist = await productsRepository.findByName(name);

    if (productExist && name !== product.name) {
      throw new AppError('There is already one product with this name');
    }

    product.name = name;
    product.price = price;
    product.quantity = quantity;

    await productsRepository.save(product);

    return product;
  }
}

export default UpdateProductService;
```

- DeleteProductService 
<p>
  delete
</p>

```
import AppError from '@shared/errors/AppError';
import { getCustomRepository } from 'typeorm';
import { ProductRepository } from '../typeorm/repositories/ProductsRepository';

interface IRequest {
  id: string;
}
class DeleteProductService {
  public async execute({ id }: IRequest): Promise<void> {
    const productsRepository = getCustomRepository(ProductRepository);

    const product = await productsRepository.findOne(id);

    if (!product) {
      throw new AppError('Product not found.');
    }

    await productsRepository.remove(product);
  }
}

export default DeleteProductService;
```

## `PRODUCTS CONTROLLER`
Controller será uma classe e teremos um método para tratamento de cada rota (listagem, exibição...)

```
export default class ProductsController {
  public async index(request: Request, response: Response) {
    const listProducts = new ListProductsService();

    const products = listProducts.execute();

    return response.json(products);
  }
}
```
> dessa forma, é uma maneira mais viável de escalar uma aplicação em execução. Dividindo as responsabilidades, onde cada recurso da aplicação tem ua tarefa e aquilo que ele precisa rfeceber ele passa pro outro. 

## `PRODUCTS ROUTES`
criando um arquivo de routes para cada módulo. 

## `Validação dos dados de routes`
```
//routes
import { Router } from 'express';
import ProductsController from '../controllers/ProductsController';
import { celebrate, Joi, Segments } from 'celebrate';

const productsRouter = Router(); //métodos do expresse nessas consts
const productsController = new ProductsController();

productsRouter.get('/', productsController.index);

productsRouter.get(
  '/:id',
  celebrate({
    [Segments.PARAMS]: {
      id: Joi.string().uuid().required(),
    },
  }),
  productsController.show,
);

productsRouter.post(
  '/',
  celebrate({
    [Segments.BODY]: {
      name: Joi.string().required(),
      price: Joi.number().precision(2).required(),
      quantity: Joi.number().required(),
    },
    [Segments.PARAMS]: {
      id: Joi.string().uuid().required(),
    },
  }),
  productsController.create,
);

productsRouter.put(
  '/:id',
  celebrate({
    [Segments.BODY]: {
      name: Joi.string().required(),
      price: Joi.number().precision(2).required(),
      quantity: Joi.number().required(),
    },
  }),
  productsController.update,
);

productsRouter.delete(
  '/:id',
  celebrate({
    [Segments.PARAMS]: {
      id: Joi.string().uuid().required(),
    },
  }),
  productsController.delete,
);

export default productsRouter;
//pq importaremos no arq principal que recebe as rotas (index.ts de src)
```

<a href="https://medium.com/trainingcenter/o-que-%C3%A9-uuid-porque-us%C3%A1-lo-ad7a66644a2b">O que é UUID?</a> 

> Celebrate: Validação de dados, foi implementada pela lib celebrate. 
> Is a middleware function that wraps the joi validation lib.

<a href="https://www.npmjs.com/package/celebrate">Celebrate-NPM Doc</a> 

- Fim do ciclo do módulo de produtos. Depois de algumas semanas, pude criar as migrações, entidade, repo, serviços, controller e as rotas. Dessa forma, estabeleceu as estrutura básica do products. 

## `MIGRATION DE USERS`
Começo do processo de regras específicas que a aplicação vai demandar. <br>
criação do módulo de users

> Demanda mais regras, uma vez que será criado perfil de usuário, avatar, troca de senha...


- TABELA PARA CRIAÇÃO DE USUÁRIOS 

```
yarn typeorm migration:create -n CreateUsers 
```

- TABELA DE USUÁRIO

```
public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.createTable(
      //pede uma instância da classe table
      new Table({
        name: 'users',
        columns: [
          {
            name: 'id',
            type: 'uuid',
            isPrimary: true,
            generationStrategy: 'uuid',
            default: 'uuid_generate_v4()',
          },
          {
            name: 'name',
            type: 'varchar',
          },
          {
            name: 'email',
            type: 'varchar',
            isUnique: true,
          },
          {
            name: 'password',
            type: 'varchar',
          },
          {
            name: 'avatar',
            type: 'varchar',
            isNullable: true,
          },
          {
            name: 'created_at',
            type: 'timestamp with time zone',
            default: 'now()',
          },
          {
            name: 'updated_at',
            type: 'timestamp with time zone',
            default: 'now()',
          },
        ],
      }),
    );
  }
```

## `ENTIDADE USER`

```
import {
  Column,
  CreateDateColumn,
  Entity,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from 'typeorm';

@Entity('users')
class User {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  name: string;

  @Column()
  email: string;

  @Column()
  password: string;

  @Column()
  avatar: string;

  @CreateDateColumn()
  created_at: Date;

  @UpdateDateColumn()
  updated_at: Date;
}

export default User;

```

## `REPO DE USERS`
Através dele que será feito as operações para manipulação do dados. 

```
import { EntityRepository, Repository } from 'typeorm';
import User from '../entities/User';

@EntityRepository(User)
class UsersRepository extends Repository<User> {
  public async findByName(name: string): Promise<User | undefined> {
    const user = await this.findOne({
      where: {
        name,
      },
    });

    return user;
  }

  public async findById(id: string): Promise<User | undefined> {
    const user = await this.findOne({
      where: {
        id,
      },
    });

    return user;
  }

  public async findByEmail(email: string): Promise<User | undefined> {
    const user = await this.findOne({
      where: {
        email,
      },
    });

    return user;
  }
}

export default UsersRepository;

```

## `CREATE USER SERVICE`
Service de criação de usuário

```
import AppError from '@shared/errors/AppError';
import { getCustomRepository } from 'typeorm';
import Product from '../typeorm/entities/Product';
import { ProductRepository } from '../typeorm/repositories/ProductsRepository';

interface IRequest {
  name: string;
  price: number;
  quantity: number;
}

class CreateProductService {
  public async execute({ name, price, quantity }: IRequest): Promise<Product> {}
}

export default CreateProductService;
```
> estrutura básica de qualquer serviço

- email não pode ser repetido por users

```
import AppError from '@shared/errors/AppError';
import { getCustomRepository } from 'typeorm';
import User from '../typeorm/entities/User';
import UsersRepository from '../typeorm/repositories/UserRepository';

interface IRequest {
  name: string;
  email: string;
  password: string;
}

class CreateUserService {
  public async execute({ name, email, password }: IRequest): Promise<User> {
    const usersRepository = getCustomRepository(UsersRepository);
    //email não pode ser repetido por users
    const emailExists = await usersRepository.findByEmail(email);

    if (emailExists) {
      throw new AppError('Emaill address already used.');
    }

    const user = usersRepository.create({
      name,
      email,
      password,
    });

    await usersRepository.save(user);

    return user;
  }
}

export default CreateUserService;
```
## `Users Controllers`
> Os arquivos de controllers são importantes em um projeto de desenvolvimento web porque eles são responsáveis por lidar com as requisições recebidas pelo servidor e enviar as respostas adequadas de volta ao cliente. Em outras palavras, é o controller que faz a ponte entre a camada de modelo e a camada de visualização em um padrão de arquitetura MVC (Model-View-Controller).

```
  import { Request, Response } from 'express';
import CreateUserService from '../services/CreateUserService';
import ListUserService from '../services/ListUserService';

export default class UsersController {
  //teremos dois métodos [index: usar o serviço de listagem & create: create user]
  public async index(request: Request, response: Response): Promise<Response> {
    const listUser = new ListUserService();

    const users = await listUser.execute();

    return response.json(users);
  }

  public async create(request: Request, response: Response): Promise<Response> {
    const { name, email, password } = request.body;

    const createUser = new CreateUserService();

    const user = await createUser.execute({
      name,
      email,
      password,
    });

    return response.json(user);
  }
}

```

> request.body é uma propriedade de um objeto Request em uma aplicação Node.js que contém os dados enviados em uma solicitação HTTP.

Quando um cliente envia uma solicitação HTTP para o servidor, essa solicitação contém informações como o método HTTP usado (por exemplo, GET, POST, PUT, DELETE), a URL da solicitação e os dados da solicitação, que são chamados de corpo (body) da solicitação.

O corpo da solicitação pode conter informações úteis para a aplicação, como parâmetros de pesquisa, informações do formulário HTML ou dados no formato JSON, XML ou outro formato. O request.body é usado para acessar esses dados da solicitação que foram enviados pelo cliente.

No entanto, o request.body não estará disponível por padrão em uma aplicação Node.js. Para acessar os dados da solicitação, é necessário usar um middleware de processamento de corpo (body parsing middleware), como o body-parser, que analisa o corpo da solicitação e a converte em um objeto JavaScript que pode ser acessado pelo request.body.

## `USER ROUTE`
O arquivo **`routes`** é um componente comum em uma aplicação Node.js que utiliza um padrão de arquitetura Model-View-Controller (MVC). Ele é responsável por mapear as solicitações HTTP recebidas pela aplicação para os controladores apropriados que irão processá-las.

O arquivo **`routes`** é geralmente definido como um módulo separado que é importado e utilizado pelo arquivo principal da aplicação. Esse módulo é responsável por configurar as rotas da aplicação, que são combinações de um método HTTP (por exemplo, GET, POST, PUT, DELETE) e um caminho de URL.

Por exemplo, uma rota pode ser definida como:

```
javascriptCopy code
app.get('/users', userController.getAllUsers);

```

Essa rota indica que a aplicação irá responder a uma solicitação HTTP GET para o caminho de URL /users e irá chamar o método getAllUsers do controlador userController.

O arquivo routes também é responsável por definir a lógica de tratamento de erros para a aplicação. Por exemplo, uma rota de tratamento de erro pode ser definida como:
```
javascript
Copy code
app.use(function (err, req, res, next) {
  console.error(err.stack);
  res.status(500).send('Something broke!');
});
```

Essa rota indica que, se ocorrer um erro na aplicação, a mensagem "Something broke!" será enviada como resposta HTTP com o status 500 (Erro interno do servidor).
Em resumo, o arquivo routes é uma parte importante de uma aplicação Node.js que ajuda a definir as rotas da aplicação e a mapear as solicitações HTTP para os controladores apropriados.

```
import { Router } from 'express';
import { celebrate, Joi, Segments } from 'celebrate';
import UsersController from '../controllers/UsersController';

const usersRouter = Router();
const usersController = new UsersController();

usersRouter.get('/', usersController.index);

usersRouter.post(
  '/',
  celebrate({
    [Segments.BODY]: {
      name: Joi.string().required(),
      email: Joi.string().email().required(),
      password: Joi.string().required(),
    },
  }),

  usersController.create,
);
export default usersRouter;
```
 > lembrar de importar e exportar arquivos. importar no arq principal de rotas `index.ts` da pasta `shared`
 
 ## `Criptografar password`
 Precisa armazenar a password já criptografada. 
 
 ```
 yarn add bcryptjs
 ```
&

```
yarn add -D @types/bcryptjs
```
> O comando "yarn add -D @types/bcryptjs" é usado para adicionar as definições de tipo TypeScript para a biblioteca "bcryptjs" em um projeto Yarn.

O sinalizador "-D" indica que o pacote deve ser instalado como uma dependência de desenvolvimento. Isso significa que ele não será incluído na construção final de produção do projeto.

O prefixo "@types/" no nome do pacote indica que este é um conjunto de definições de tipo TypeScript para a biblioteca "bcryptjs", o que permite usar a biblioteca em um projeto TypeScript com suporte completo de verificação de tipos e editor.

Em geral, este comando é útil para desenvolvedores que estão construindo um projeto TypeScript que usa a biblioteca "bcryptjs" e desejam garantir que seu código seja bem-tipado e livre de erros relacionados a tipos.

## `CRIANDO SERVIÇO DE AUTENTICAÇÃO`
Recomenda-se o uso do médoto POST http. Uma vez que os dados são verificados e auticados na aplicação. 

- Usar lib json web token, onde será controlado o acesso a determinadas rotas, dessa forma, protegendo com o uso de um token. 
- Para isso foi criado um service dentro do módulo de users

- Criação de um novo cotroller (sessionsController) e em seguida será implementada na rota post. 

### `controller`

```
import { Request, Response } from 'express';
import CreateSessionsService from '../services/CreateSessionsService';

export default class SessionsController {
  public async create(request: Request, response: Response): Promise<Response> {
    const { email, password } = request.body;

    const createSession = new CreateSessionsService();

    const user = await createSession.execute({
      email,
      password,
    });

    return response.json(user);
  }
}
```
> A classe SessionsController é responsável por tratar as requisições HTTP relacionadas à autenticação de um usuário em um sistema, especificamente a criação de uma sessão (login).

O método create é utilizado para criar uma nova sessão. Ele recebe um objeto Request e um objeto Response do Express, que contêm informações da requisição e da resposta HTTP, respectivamente. O método utiliza o serviço CreateSessionsService para criar a sessão e retorna uma resposta JSON contendo o usuário autenticado.

O parâmetro : Promise<Response> indica que o método retorna uma promessa (Promise) que eventualmente será resolvida com uma resposta HTTP.

### `routes`

```
import { Router } from 'express';
import { celebrate, Joi, Segments } from 'celebrate';
import SessionsController from '../controllers/SessionsController';

const sessionsRouter = Router();
const sessionsController = new SessionsController();

sessionsRouter.post(
  '/',
  celebrate({
    [Segments.BODY]: {
      email: Joi.string().email().required(),
      password: Joi.string().required(),
    },
  }),

  sessionsController.create,
);
export default sessionsRouter;
```
> código apresentado é responsável por definir as rotas para o endpoint de sessões de usuários em uma aplicação web utilizando o framework Express.js e a biblioteca Celebrate para validação de dados.

Aqui está o que cada linha do código faz:

A primeira linha importa o objeto Router do módulo express.
A segunda linha importa a função celebrate, a classe Joi e o objeto Segments da biblioteca celebrate. Esses recursos são usados para validar os dados enviados pelo usuário.
A terceira linha importa o controlador SessionsController que será utilizado para processar as requisições para o endpoint de sessões.
A quarta linha cria uma instância do Router para definir as rotas do endpoint de sessões.
A quinta linha cria uma instância do controlador SessionsController.
A sexta linha define a rota POST para o endpoint de sessões. Ela utiliza a função celebrate para validar os dados do corpo da requisição, especificamente o email e a senha do usuário. Se a validação falhar, um erro será retornado. Se a validação for bem-sucedida, a função create do controlador SessionsController será executada para criar uma nova sessão de usuário.
Por fim, a última linha exporta o Router para que ele possa ser utilizado em outras partes da aplicação.

Em resumo, este trecho de código define a rota POST para o endpoint de sessões de usuários, valida os dados enviados pelo usuário utilizando a biblioteca Celebrate e chama o método create do controlador SessionsController para processar a requisição.
  
### `instalação JWT`
importação da lib. 

`yarn add json web token`
  
após esse comando é necessário add os types.
  
`yarn add -D @types/jsonwebtoken`
  
<br>
  
 - Após certeza que os dados dos users estão corretos, será configurado o token. 
  
  ![image](https://user-images.githubusercontent.com/101754313/222246426-53508fde-0816-441b-b05f-7034e3b38408.png)

 
 ![image](https://user-images.githubusercontent.com/101754313/222249712-dba50cd5-7b3a-47ab-944d-7f7761a7a2c1.png)
  
  O código que você forneceu cria um objeto chamado token que é gerado usando uma biblioteca de geração de tokens chamada jsonwebtoken. O objeto token é composto por três partes:

Um payload vazio: {}
Uma chave secreta: 'e5814213ceea650a393034bfdd481d95'
Algumas opções, incluindo um identificador de usuário (subject) e uma data de validade (expiresIn).
Esses valores são usados para criar um token que será usado para autenticar o usuário. A chave secreta é usada para criptografar e descriptografar o token e garantir que ele não seja falsificado ou alterado.

O identificador do usuário é adicionado ao token como uma forma de garantir que o token seja exclusivo para aquele usuário específico. A data de validade é usada para garantir que o token expira após um determinado período de tempo (1 dia, neste caso).

O código retorna um objeto contendo o usuário autenticado e o token gerado, que serão usados para autenticar e autorizar o usuário em futuras solicitações de API.
  
  
  
 ### `session service`
  
 ```
 import AppError from '@shared/errors/AppError';
import { compare } from 'bcryptjs';
import { sign } from 'jsonwebtoken';
import { getCustomRepository } from 'typeorm';
import User from '../typeorm/entities/User';
import UsersRepository from '../typeorm/repositories/UserRepository';

interface IRequest {
  email: string;
  password: string;
}

interface IResponse {
  user: User;
  token: string;
}

class CreateSessionsService {
  public async execute({ email, password }: IRequest): Promise<IResponse> {
    const usersRepository = getCustomRepository(UsersRepository);
    //email não pode ser repetido por users
    const user = await usersRepository.findByEmail(email);

    if (!user) {
      throw new AppError('incorrect email/password combination.', 401);
    }

    const passwordConfirmed = await compare(password, user.password);

    if (!passwordConfirmed) {
      throw new AppError('incorrect email/password combination.', 401);
    }

    const token = sign({}, 'e5814213ceea650a393034bfdd481d95', {
      subject: user.id, //user autorizado a usar o token
      expiresIn: '1d',
    });

    return {
      user,
      token,
    };
  }
}

export default CreateSessionsService; 
 ```
  
> um serviço que é responsável por criar uma sessão de usuário, verificando se as credenciais fornecidas são válidas e gerando um token JWT para o usuário autenticado.

O serviço CreateSessionsService contém um método execute, que recebe as informações do usuário (email e senha) como parâmetros e retorna um objeto que contém o usuário autenticado e o token gerado.

No início do método, o serviço obtém uma instância do repositório personalizado UsersRepository usando a função getCustomRepository do TypeORM. Em seguida, o serviço procura um usuário com o email fornecido usando o método findByEmail do repositório. Se o usuário não for encontrado, uma exceção é lançada com uma mensagem de erro.

Se o usuário for encontrado, o serviço verifica se a senha fornecida corresponde à senha armazenada no banco de dados, usando a função compare do bcryptjs. Se as senhas não corresponderem, uma exceção é lançada com uma mensagem de erro.

Se as credenciais do usuário forem válidas, um token JWT é gerado usando a função sign do jsonwebtoken. O objeto de payload é vazio, a chave secreta é fornecida como um argumento e o objeto de opções contém o identificador do usuário e a data de validade do token.

Finalmente, o serviço retorna um objeto que contém o usuário autenticado e o token gerado.
  
## `MIDDLEWARE AUTHETICATION PARA PROTEÇÃO DAS ROTAS`
  
 ```
 import AppError from '@shared/errors/AppError';
import { NextFunction, Request, Response } from 'express';
import { verify } from 'jsonwebtoken';
import authConfig from '@config/auth';

export default function isAuthenticated(
  request: Request,
  response: Response,
  next: NextFunction,
): void {
  const authHeader = request.headers.authorization;

  if (!authHeader) {
    throw new AppError('JWT Token is missing');
  }

  const [, token] = authHeader.split(' ');

  try {
    const decodeToken = verify(token, authConfig.jwt.secret);

    return next();
  } catch {
    throw new AppError('Invalid JWT Token.');
  }
} 
 ```
  
> um middleware de autenticação em que é utilizado o token JWT para proteger as rotas do Express.
  
- O middleware recebe três parâmetros: `request`, `response` e `next`. O parâmetro `next` é uma função de callback que permite que a execução prossiga para o próximo middleware ou controlador.
  
- TRY/CATCH: essa verificação pode ocorrer uma falha. Portanto, deve-se capturá-la utilizando o try/catch. Uma verify que não é feita pela aplication, sendo passada por um método da lib jsonwebtoken.

- O middleware primeiro extrai o token do cabeçalho de autorização da solicitação. Se o cabeçalho de autorização não estiver presente, ele lança uma exceção com uma mensagem de erro "JWT Token is missing".

- Se o cabeçalho de autorização estiver presente, o middleware divide o cabeçalho em duas partes: o esquema de autenticação e o token JWT. Ele descarta o esquema de autenticação e armazena o token JWT.

- Em seguida, o middleware verifica se o token JWT é válido usando a função `verify` do `jsonwebtoken`. Se o token não for válido, ele lança uma exceção com a mensagem de erro "Invalid JWT Token".

- Se o token for válido, a função `next()` é chamada para permitir que a execução prossiga para o próximo middleware ou controlador.

- Em resumo, este middleware é responsável por verificar se o token JWT é válido e permite que as rotas protegidas pelo token sejam acessadas somente por usuários autenticados e autorizados. Se o token não for válido, o middleware bloqueia o acesso às rotas e envia uma mensagem de erro ao cliente.
  
  ![image](https://user-images.githubusercontent.com/101754313/222262285-bba130d0-66db-411f-ba87-de538be225aa.png)


