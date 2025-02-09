# 📌 Projeto de Leilão

Este repositório contém a implementação de um sistema de leilão desenvolvido em Go. Abaixo, seguem as instruções para execução e testes.

## 🚀 Build e Execução

Para construir e rodar o projeto, utilize o comando:

```bash
docker compose up -d --build
```

Isso iniciará os containers necessários para a aplicação.

---

## 🧪 Rodando os Testes

Para rodar os testes, execute:

```bash
docker compose exec app go test ./internal/infra/database/auction/ -v
```

Se os testes passarem, a saída esperada será algo semelhante a:

```swift
PASS
ok      fullcycle-auction_go/internal/infra/database/auction    X.XXXs
```

---

## 🛑 Parando os Containers

Para parar e remover os containers, utilize:

```bash
docker compose down
```

Isso garantirá que nenhum container do projeto permaneça em execução.

---
