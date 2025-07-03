# API RESTful em Go com SQLite3

Uma API RESTful simples desenvolvida em Go (Golang) para gerenciar dados, armazenando-os em um banco de dados SQLite3. Inclui endpoints para operações CRUD (Criar, Ler, Atualizar, Deletar) e uma funcionalidade de verificação em segundo plano para identificar registros com datas de expiração

### **Esse repositório foi feito para o objetivo de ser usado nas aulas do curso da Alura sobre DevOps**

---

## 🚀 Tecnologias Utilizadas

* **Go (Golang)**
* **Gorilla Mux:** Roteador HTTP para construir as rotas da API.
* **SQLite3:** Banco de dados leve e embarcado.
* **Driver `github.com/mattn/go-sqlite3`:** Driver Go para interagir com o SQLite3.

---

## ✨ Funcionalidades

* **CRUD Completo:**
    * `POST /data`: Cria um novo registro de dado.
    * `GET /data`: Lista todos os registros de dados.
    * `GET /data/{id}`: Recupera um registro específico pelo ID.
    * `PUT /data/{id}`: Atualiza um registro existente.
    * `DELETE /data/{id}`: Exclui um registro.
* **Verificação de Expiração:** Uma goroutine em segundo plano verifica periodicamente (a cada 1 hora por padrão) por registros cuja `data_de_expiracao` já passou, registrando-os no console. Isso pode ser estendido para outras ações, como exclusão ou notificação.
* **Banco de Dados SQLite3:** Armazena os dados de forma local e persistente no arquivo `data.db`.

---

## 🛠️ Como Rodar o Projeto

Siga os passos abaixo para configurar e executar a API em sua máquina local.

### Pré-requisitos

Certifique-se de ter o **Go** (versão 1.16 ou superior recomendada) instalado em seu sistema.

### Instalação

1.  **Clone o repositório:**
    ```bash
    git clone https://github.com/NothNada/go-api.git
    cd go-api # ou o nome que você deu ao projeto, ex: go-api-sqlite
    ```

2.  **Baixe as dependências:**
    ```bash
    go mod tidy
    ```

### Execução

1.  **Inicie o servidor:**
    ```bash
    go run main.go
    ```
    O servidor será iniciado na porta `8080`. Você verá a mensagem `Servidor iniciado na porta :8080` no seu terminal. O arquivo `data.db` será criado automaticamente na raiz do projeto, se ainda não existir.

---

## 🚦 Endpoints da API

A API está disponível em `http://localhost:8080`.

### 1. Criar um Registro

* **URL:** `/data`
* **Método:** `POST`
* **Headers:** `Content-Type: application/json`
* **Body (JSON):**
    ```json
    {
        "usuario_que_registrou": "nome.usuario",
        "data_registrada": "YYYY-MM-DDTHH:MM:SSZ",
        "descricao_curta": "Breve descrição do item",
        "descricao_longa": "Descrição detalhada do registro, pode ser mais longa.",
        "data_de_expiracao": "YYYY-MM-DDTHH:MM:SSZ"
    }
    ```
    * `data_registrada`: Se omitida, a API usará a data e hora atuais.
    * `data_de_expiracao`: Campo opcional. Se não for fornecido ou for nulo, o registro não será considerado para expiração pela rotina de verificação.
* **Exemplo com `curl`:**
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{
        "usuario_que_registrou": "joao.silva",
        "data_registrada": "2023-10-26T10:00:00Z",
        "descricao_curta": "Reunião de equipe",
        "descricao_longa": "Reunião para discutir o planejamento do próximo trimestre.",
        "data_de_expiracao": "2024-10-26T10:00:00Z"
    }' http://localhost:8080/data
    ```

### 2. Obter Todos os Registros

* **URL:** `/data`
* **Método:** `GET`
* **Exemplo com `curl`:**
    ```bash
    curl http://localhost:8080/data
    ```

### 3. Obter um Registro por ID

* **URL:** `/data/{id}`
* **Método:** `GET`
* **Parâmetros:** `{id}` - ID numérico do registro.
* **Exemplo com `curl`:**
    ```bash
    curl http://localhost:8080/data/123
    ```

### 4. Atualizar um Registro

* **URL:** `/data/{id}`
* **Método:** `PUT`
* **Parâmetros:** `{id}` - ID numérico do registro a ser atualizado.
* **Headers:** `Content-Type: application/json`
* **Body (JSON):** O mesmo formato do `POST`, mas com os dados a serem atualizados.
* **Exemplo com `curl`:**
    ```bash
    curl -X PUT -H "Content-Type: application/json" -d '{
        "usuario_que_registrou": "maria.souza",
        "descricao_curta": "Reunião de projeto atualizada",
        "data_de_expiracao": "2025-10-26T10:00:00Z"
    }' http://localhost:8080/data/123
    ```

### 5. Excluir um Registro

* **URL:** `/data/{id}`
* **Método:** `DELETE`
* **Parâmetros:** `{id}` - ID numérico do registro a ser excluído.
* **Exemplo com `curl`:**
    ```bash
    curl -X DELETE http://localhost:8080/data/123
    ```

---