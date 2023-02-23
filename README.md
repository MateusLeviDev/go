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

<a href="[https://medium.com/trainingcenter/o-que-%C3%A9-uuid-porque-us%C3%A1-lo-ad7a66644a2b](https://www.npmjs.com/package/celebrate)">Celebrate-NPM Doc</a> 

- Fim do ciclo do módulo de produtos. Depois de algumas semanas, pude criar as migrações, entidade, repo, serviços, controller e as rotas. Dessa forma, estabeleceu as estrutura básica do products. 

## `Migration de Users`
Começo do processo de regras específicas que a aplicação vai demandar. <br>
criação do módulo de users

> Demanda mais regras, uma vez que será criado perfil de usuário, avatar, troca de senha...


```
TABELA PARA CRIAÇÃO DE USUÁRIOS



yarn typeorm migration:create -n CreateUsers 
```
