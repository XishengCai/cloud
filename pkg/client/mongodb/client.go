package mongodb

import (
	"cloud/pkg/setting"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/klog"
)

var (
	Client *mongo.Client
)

func InitMongoDB() {
	klog.Info("InitMongoDB")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongodbURI := fmt.Sprintf("mongodb://%s", setting.MongodbSetting.Addresses)
	credential := options.Credential{
		Username: setting.MongodbSetting.User,
		Password: setting.MongodbSetting.Password,
	}
	klog.Infof("connecting mongodb %s, credential: %v", mongodbURI, credential)
	clientOpts := options.Client().ApplyURI(mongodbURI).SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	Client = client
}
