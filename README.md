
# Projeto Playwright Golang - Testes com API Alpha Vantage


![Run Test](https://github.com/OWNER/REPOSITORY/actions/workflows/api_testing.yaml/badge.svg)


Este projeto utiliza a biblioteca [playwright-go](https://github.com/playwright-community/playwright-go) para realizar testes em APIs, especificamente na [Alpha Vantage API](https://www.alphavantage.co/documentation/).

## Objetivo

O objetivo deste projeto é testar a lib *playwright-go*

## Estrutura do Projeto

A estrutura do projeto está organizada da seguinte forma:

```
playwright_golang/
├── playwrightsetup/
│   └── setup.go           # Configuração do Playwright
├── tests/
│   ├── alpha_vantage_test.go   # Arquivo de testes para a API Alpha Vantage
├── go.mod                 # Definição de módulos Go
├── go.sum                 # Dependências do projeto
├── main.go                # Arquivo principal (não utilizado diretamente nos testes)
└── report.json            # Relatório gerado dos testes
```

## Bibliotecas Utilizadas

As principais bibliotecas utilizadas no projeto são:

- [playwright-go](https://github.com/playwright-community/playwright-go) - Utilizada para emular o comportamento de navegadores e executar requisições HTTP.
- [Testify](https://github.com/stretchr/testify) - Utilizada para asserções nos testes.

## Pré-requisitos

Antes de começar, você precisará ter as seguintes ferramentas instaladas em sua máquina:

- Go versão 1.22 ou superior
- Uma conta e uma chave de API da Alpha Vantage (disponível gratuitamente [aqui](https://www.alphavantage.co/support/#api-key))

## Como Rodar o Projeto

1. Clone o repositório:
   ```
   git clone https://github.com/seu_usuario/playwright_golang
   ```

2. Entre no diretório do projeto:
   ```
   cd playwright_golang
   ```

3. Exporte a chave da API da Alpha Vantage como uma variável de ambiente:
   ```
   export ALPHA_VANTAGE_API_KEY="sua-chave-de-api"
   ```

4. Execute os testes:
   ```
   go test ./tests
   ```

## Funcionalidades Testadas

### 1. Testar Dados Intraday de uma Ação
Testa a resposta da API com dados de uma ação intraday para um determinado símbolo (como IBM) e intervalo de tempo (como 5min).

### 2. Testar Taxa de Câmbio
Verifica se a API retorna corretamente a taxa de câmbio entre duas moedas, como USD e EUR.

### 3. Testar Endpoint Inválido
Verifica como a API responde quando um endpoint inválido é acessado.

### 4. Testar Cotação Global
Verifica se a API retorna as cotações globais para uma ação, como IBM.

## Executando Testes Específicos

Para executar testes específicos, utilize o seguinte comando:
```
go test -run <NomeDoTeste>
```

Exemplo:
```
go test -run TestIntradayStockData
```

## Contribuindo

Sinta-se à vontade para abrir issues e enviar pull requests.

## Licença

Este projeto está licenciado sob a licença MIT. Consulte o arquivo [LICENSE](LICENSE) para obter mais informações.
