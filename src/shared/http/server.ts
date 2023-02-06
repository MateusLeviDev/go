import express from 'express';
import cors from 'cors';
import routes from './routes';

const app = express();

app.use(cors());
app.use(express.json()); //não interpreta json por padrão

app.use(routes);

app.listen(3333, () => {
  console.log('Server started on port 3333!'); //vai chamar o servidor
});
