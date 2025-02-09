📌 Execução e Testes do Projeto de Leilão
🚀 Build e Execução
bash
Copy
Edit
docker compose up -d --build
🧪 Rodando os Testes
bash
Copy
Edit
docker compose exec app go test ./internal/infra/database/auction/ -v
Se os testes passarem, a saída esperada será:

swift
Copy
Edit
PASS
ok      fullcycle-auction_go/internal/infra/database/auction    X.XXXs
🛑 Parando os Containers
bash
Copy
Edit
docker compose down
