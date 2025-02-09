package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {

	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Erro ao tentar inserir leilão no MongoDB:", err)
		return internal_error.NewInternalServerError("Error trying to insert auction: " + err.Error())
	}
	go ar.scheduleAuctionClosure(auctionEntity.Id)

	return nil
}

func (ar *AuctionRepository) scheduleAuctionClosure(auctionID string) {
	auctionInterval, err := time.ParseDuration(os.Getenv("AUCTION_INTERVAL"))
	if err != nil {
		logger.Error("Error trying to parse AUCTION_INTERVAL", err)
		auctionInterval = 30 * time.Second

	}
	<-time.After(auctionInterval)

	ctx := context.Background()
	filter := bson.M{
		"_id":    auctionID,
		"status": auction_entity.Active,
	}
	update := bson.M{
		"$set": bson.M{"status": auction_entity.Completed},
	}

	res, updateErr := ar.Collection.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		logger.Error("Error trying to update auction", updateErr)
		return
	}

	if res.ModifiedCount > 0 {
		logger.Info("Auction closed automatically. ")
	} else {
		logger.Info("Auction was not active or did not exist. ")
	}

}
