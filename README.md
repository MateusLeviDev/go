<h1 align="center">
  gestao-vendas-apirestful-node
</h1>

 <p>Aplicação Back-end para gestão de vendas para restaurante local. Principais tecnologias que utilizaremos para desenvolvimento da API: Node.js, Express, Typescript, TypeORM, Postgres através de container Docker, Redis através de container Docker, Amazon S3, Amazon SES.
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

