//routes
import { Router } from 'express';
import ProductsController from '../controllers/ProductsController';

const productsRouter = Router(); //métodos do expresse nessas consts
const productsController = new ProductsController();

productsRouter.get('/', productsController.index);
productsRouter.get('/:id', productsController.show);
productsRouter.post('/', productsController.create);
productsRouter.put('/:id', productsController.update);
productsRouter.delete('/:id', productsController.delete);

export default productsRouter;
//pq importaremos no arq principal que recebe as rotas (index.ts de src)
