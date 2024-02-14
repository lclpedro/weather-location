## GoExpert

### Para a execução da aplicação e seus tests, basta executar os seguintes comandos
#### Execução local
Para executar a aplicação localmente, primeiramente você deve definir a API KEY do Weather API.
PS: Se você não tem uma conta, pode definir uma aqui: https://www.weatherapi.com/

```sh
export WEATHER_API_TOKEN=<API_KEY>
```

Após isso para inicializar a aplicação é muito simples, basta você executar o seguinte comando.

```sh
docker-compose up app 
```
Isso estará funcionando na porta :8080 do seu localhost.

### Documentação de Endpoint
URL /weather/:cep
parametro: CEP da sua localidade

