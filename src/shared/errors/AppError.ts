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
