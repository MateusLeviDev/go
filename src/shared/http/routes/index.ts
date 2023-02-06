// arquivo principal de rotas(routes)
import { Router } from 'express';

const routes = Router();

routes.get('/', (resquest, response) => {
  return response.json({ message: 'Fala, levs' }); //'/' significa que vai pegar a rota raiz
});

export default routes;
