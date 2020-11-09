package common

import (
	"context"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mgoClient *mongo.Client
var mgoDbName string

func InitMgoClient(uri string, dbName string, maxPoolSize uint64) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMaxPoolSize(maxPoolSize)) // 最大连接池
	if err != nil {
		log.Println("mongo.Connect err:", err.Error())
		return nil, err
	}

	err = cli.Ping(ctx, nil)
	if err != nil {
		log.Println("mongo.Ping err", err.Error())
		return nil, err
	}

	mgoClient = cli
	mgoDbName = dbName
	log.Println("Connected to MongoDB!")
	return cli, err
}

// 所有model结构体继承Mgo
type Mgo struct {
	coll     *mongo.Collection
	collName string
}

// 设置数据库名
func (m *Mgo) SetCollName(collName string) {
	m.collName = collName
}

// 获取数据库名
func (m *Mgo) GetCollName() string {
	return m.collName
}

// 获取表
func (m *Mgo) GetCollection() *mongo.Collection {
	if len(m.collName) == 0 {
		panic("please set CollName")
	}
	if m.coll == nil {
		m.coll = mgoClient.Database(mgoDbName).Collection(m.collName)
	}
	return m.coll
}

// 创建索引
func (m *Mgo) CreateIndex(keys map[string]int, Unique bool) (string, error) {
	if len(keys) == 0 {
		return "", nil
	}

	idx := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(Unique),
	}

	result, err := m.GetCollection().Indexes().CreateOne(context.Background(), idx)
	return result, err
}

// 删除索引
func (m *Mgo) DropIndex(name string) error {
	_, err := m.GetCollection().Indexes().DropOne(context.Background(), name)
	return err
}

// 插入单条数据
func (m *Mgo) InsertOne(document interface{}) (primitive.ObjectID, error) {
	res, err := m.GetCollection().InsertOne(context.Background(), document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

// 查询单条数据
func (m *Mgo) GetByField(result interface{}, filter interface{}) error {
	err := m.GetCollection().FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return err
	}

	return nil
}

// 更新单条数据
func (m *Mgo) UpdateByField(result interface{}, filter interface{}, update interface{}) error {
	opts := options.FindOneAndUpdate().SetReturnDocument(1)
	err := m.GetCollection().FindOneAndUpdate(context.Background(), filter, update, opts).Decode(result)
	if err != nil {
		return err
	}

	return nil
}

// 删除单条数据
func (m *Mgo) DeleteByField(filter interface{}) (int64, error) {
	res, err := m.GetCollection().DeleteOne(context.Background(), filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// 插入多条数据

// 查询列表数据(支持分页)
func (m *Mgo) GetList(index int64, limit int64, sorts []string, filters interface{}, pipeline []bson.D) (List, error) {
	var result List

	// 获取total
	total, _ := m.GetCollection().CountDocuments(context.Background(), filters)

	result.Index = index
	result.Limit = limit
	result.Total = total
	result.Page = int64(math.Ceil(float64(total) / float64(limit)))

	limitStage := bson.D{{"$limit", limit}}

	skip := (index - 1) * limit
	skipStage := bson.D{{"$skip", skip}}

	var sortStage bson.D
	if len(sorts) == 0 {
		sortStage = bson.D{
			{"$sort", bson.D{
				{"create_at", 1},
			}},
		}
	} else {
		sort := []bson.E{}
		for _, i := range sorts {
			split := strings.Split(i, ",")
			field := split[0]
			order, _ := strconv.Atoi(split[1])
			sort = append(sort, bson.E{field, order})
		}
		sortStage = bson.D{{"$sort", sort}}
	}

	pipeline = append(pipeline, limitStage, skipStage, sortStage)

	// 获取list
	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	cursor, err := m.GetCollection().Aggregate(context.Background(), pipeline, opts)
	if err != nil {
		return List{}, err
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var document bson.M
		cursor.Decode(&document)
		result.List = append(result.List, document)
	}

	return result, nil
}

// 删除多条数据
func (m *Mgo) BulkDele(filter interface{}) (int64, error) {
	res, err := m.GetCollection().DeleteMany(context.Background(), filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}
