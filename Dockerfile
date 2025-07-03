# Estágio 1: build CGO com suporte a sqlite
FROM golang:1.24.3-alpine AS builder

# Instala as dependências de build para CGO e SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia os arquivos go.mod e go.sum para gerenciar as dependências
COPY go.mod go.sum ./

# Baixa as dependências do módulo Go
RUN go mod download

# Copia o restante do código da aplicação para o diretório de trabalho
COPY . .

# Constrói a aplicação Go com CGO habilitado e otimizações
# -tags sqlite_omit_load_extension pode ser adicionado se você não precisar de extensões SQLite
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-s -w' -o app ./main.go

# Etapa 2: runtime leve em Alpine
FROM alpine:latest

# Instala as dependências de runtime necessárias (incluindo a lib sqlite)
RUN apk add --no-cache ca-certificates sqlite

# Define o diretório de trabalho para onde o executável e o DB irão
WORKDIR /app

# Copia o executável construído do estágio de construção
COPY --from=builder /app/app .

# Define o diretório para o banco de dados SQLite
# Assumimos que sua aplicação Go criará o data.db em './data.db'
# ou você montará um volume para /app (ou um subdiretório)
# Se você quer o DB em um lugar específico, ajuste a lógica da sua aplicação Go
# e o volume aqui. Por padrão, ele será criado no WORKDIR /app
# Se você quer um volume explícito para o DB, pode ser assim:
VOLUME ["/data"] # O diretório para onde seu DB deve ser mapeado externamente

# Expõe a porta que sua aplicação Go está escutando
EXPOSE 8080

# Comando para iniciar a aplicação quando o contêiner for executado
# Certifique-se de que sua aplicação Go use "/data/data.db" ou o caminho correto
# se você estiver persistindo o banco de dados em /data
# Se sua aplicação abre "data.db" (no diretório atual), então o volume deve ser montado em /app
# Por simplicidade, vamos manter a lógica atual de /app para o DB:
CMD ["./app"]