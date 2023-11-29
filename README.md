# Coin Commerce - Projeto de Estudo em Golang

Bem-vindo ao Coin Commerce, um projeto pessoal fictício desenvolvido em Golang para fins de estudo. Esta aplicação simula um sistema de compra e venda de cryptoativos e é totalmente containerizada, utilizando Docker para facilitar a execução e distribuição.

## Estrutura do Projeto

O projeto está organizado da seguinte forma:

- **API REST**: Responsável por receber e direcionar as requisições para os microserviços correspondentes.

- **Microserviço Users**: Lida com operações relacionadas a usuários, incluindo cadastro, autenticação e informações de perfil.

- **Microserviço Orders**: Gerencia operações relacionadas a pedidos de compra e venda de cryptoativos.

- **Load Balancer**: Distribui a carga de requisições entre os microserviços para garantir um desempenho eficiente e balanceado.

- **RabbitMQ**: Utilizado como sistema de mensagens para facilitar a comunicação entre os diferentes componentes. As mensagens são codificadas utilizando gRPC e definidas através de arquivos .proto.

- **Banco de Dados PostgreSQL**: Armazena os dados essenciais para o funcionamento do sistema.

## Pré-requisitos

O único requisito para execução local é ter o Golang instalado em seu ambiente.

## Configuração e Execução com Docker

1. Clone o repositório para o seu ambiente local.
   ```bash
   git clone https://github.com/guilherme-de-marchi/coin-commerce.git
   cd coin-commerce
   ```

2. Compile o programa.
   ```bash
   go build -o coin-commerce
   ```

3. Execute os containers Docker.
   ```bash
   docker compose up
   ```

## Uso

A API REST oferece endpoints para realizar operações relacionadas a usuários e pedidos. Consulte a documentação da API para obter detalhes sobre os endpoints disponíveis e como utilizá-los.

## Contribuição

Se você estiver interessado em contribuir para o desenvolvimento deste projeto de estudo, sinta-se à vontade para abrir issues, enviar pull requests ou entrar em contato com a equipe de desenvolvimento.

## Licença

Este projeto é licenciado sob a MIT License.
