package auction_test

import (
	"context"
	"github.com/google/uuid"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuctionTestSuite é uma estrutura para testes que inclui setup e teardown
type AuctionTestSuite struct {
	suite.Suite
	Repo       *auction.AuctionRepository
	DB         *mongo.Database
	Collection *mongo.Collection
}

// SetupSuite inicializa o banco de dados de teste antes de rodar os testes
func (suite *AuctionTestSuite) SetupSuite() {
	os.Setenv("AUCTION_INTERVAL", "2s") // Tempo curto para fechamento

	// Pega as credenciais do MongoDB do ambiente
	mongoURI := "mongodb://admin:admin@mongodb:27017/?authSource=admin"

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	suite.Require().NoError(err, "Erro ao conectar no MongoDB de teste")

	err = client.Ping(context.Background(), nil)
	suite.Require().NoError(err, " MongoDB não está acessível. Verifique se está rodando e autenticado.")

	// Banco de testes
	suite.DB = client.Database("test_auctions")
	suite.Collection = suite.DB.Collection("auctions")
	suite.Repo = &auction.AuctionRepository{Collection: suite.Collection}
}

func (suite *AuctionTestSuite) TearDownSuite() {
	suite.DB.Drop(context.Background()) // Apaga a base de testes inteira
}

func (suite *AuctionTestSuite) TestAuctionClosesAutomatically() {
	ctx := context.Background()

	auctionID := uuid.New().String() // Gera um ID único

	auctionEntity := &auction_entity.Auction{
		Id:          auctionID,
		ProductName: "Produto Teste",
		Category:    "Categoria Teste",
		Description: "Descrição para testes",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	// Insere o leilão no banco de dados
	err := suite.Repo.CreateAuction(ctx, auctionEntity)
	if err != nil {
		suite.T().Fatalf("Falha ao criar leilão: %v", err)
	}

	suite.T().Logf("Leilão criado com sucesso: ID=%s", auctionID)

	// Aguarda o tempo necessário para o fechamento automático
	time.Sleep(3 * time.Second)

	// Busca o leilão no banco
	var closedAuction auction_entity.Auction
	filter := bson.M{"_id": auctionID}

	mongoErr := suite.Collection.FindOne(ctx, filter).Decode(&closedAuction)
	if mongoErr != nil {
		suite.T().Fatalf(" Erro ao buscar o leilão após tempo de espera: %v", mongoErr)
		return
	}

	// 🔥 Loga o status do leilão após a espera
	suite.T().Logf(" Status do leilão após espera: %v", closedAuction.Status)

	// Confere se o leilão foi fechado automaticamente
	assert.Equal(suite.T(), auction_entity.Completed, closedAuction.Status, " Leilão deveria estar fechado automaticamente")
}

func TestAuctionSuite(t *testing.T) {
	suite.Run(t, new(AuctionTestSuite))
}
