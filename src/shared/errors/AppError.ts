// clase tratamento de erro
class AppError {
  public readonly message: string;
  public readonly statudCode: number;

  constructor(message: string, statusCode = 400) {
    this.message = message;
    this.statudCode = statusCode;
  }
}

export default AppError;

//criação de middleware resposanvel por interceptar o erros gerados pela aplicação dai customizar esse erro para o front end. Estrutura de um middle, sabe toda req tem uma request response, e nele tem uma terceira. particularidade tem um 4 parametro, no qual recebe o erro.

/* app.use(
  (
    error: Error,
    request: Request,
    response: Response,
    next: NextFunction
    ) => {},
  );
 */

//criação de um middleware que recebe um error. Se a instância do erro é da classe, uma vez que será usado AppError no serviço, aparecerá a mensagem. Agora caso não seja, provavelmente é um erro fora da aplicação, portando, status 500, erro desconhecido.
