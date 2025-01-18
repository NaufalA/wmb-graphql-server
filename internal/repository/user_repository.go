package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/NaufalA/wmb-graphql-server/graph/model"
	"github.com/NaufalA/wmb-graphql-server/internal/collection"
	"github.com/NaufalA/wmb-graphql-server/internal/dto"
	"github.com/NaufalA/wmb-graphql-server/pkg/util"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository struct {
	logger *logrus.Entry
	db     *mongo.Database
}

func NewUserRepository(
	logger *logrus.Logger,
	db *mongo.Database,
) *UserRepository {
	return &UserRepository{
		logger: logger.WithFields(logrus.Fields{
			"location": "UserRepository",
		}),
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, input model.CreateUserInput) (*collection.User, error) {
	id := primitive.ObjectID(bson.NewObjectID())
	now := time.Now()
	passwordUtil := util.PasswordUtil{}
	passwordHash, err := passwordUtil.HashPassword(input.Password)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	role := "Guest"
	col := collection.User{
		ID:           id,
		Email:        &input.Email,
		FullName:     &input.FullName,
		Role:         &role,
		PasswordHash: &passwordHash,
		CreateTime:   &now,
	}
	_, err = r.db.Collection(col.CollectionName()).InsertOne(ctx, col)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &col, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, input collection.User) (*collection.User, error) {
	user := collection.User{}
	filter := bson.M{"_id": input.ID}
	err := r.db.Collection(user.CollectionName()).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	now := time.Now()
	if input.Email != nil {
		user.Email = input.Email
	}
	if input.FullName != nil {
		user.FullName = input.FullName
	}
	if input.Role != nil {
		user.Role = input.Role
	}
	if input.Password != nil {
		passwordUtil := util.PasswordUtil{}
		passwordHash, err := passwordUtil.HashPassword(*input.Password)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		user.PasswordHash = &passwordHash
	}
	user.UpdateTime = &now
	update := bson.M{"$set": user}

	collectionName := collection.User{}.CollectionName()
	result, err := r.db.Collection(collectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if result.ModifiedCount == 0 {
		logrus.Error("no data modified")
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) (*string, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	col := collection.User{
		ID: objectID,
	}

	filter := bson.M{"_id": col.ID}
	result, err := r.db.Collection(col.CollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if result.DeletedCount == 0 {
		logrus.Error("no data deleted")
		return nil, err
	}

	return nil, nil
}

func (r *UserRepository) GetUser(ctx context.Context, request dto.GetUserRequest) (*collection.User, error) {
	filter := bson.M{}
	if request.ID != "" {
		objectID, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		filter["_id"] = objectID
	}
	if request.Email != "" {
		filter["email"] = request.Email
	}
	col := collection.User{}
	err := r.db.Collection(col.CollectionName()).FindOne(ctx, filter).Decode(&col)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &col, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, input model.UserConnectionArgs) (*model.UserConnection, error) {
	paginationUtil := util.PaginationUtil{}
	filter := bson.M{}
	var opts *options.FindOptionsBuilder
	direction := 1
	if input.First != nil {
		if *input.First < 0 {
			err := fmt.Errorf("invalid parameter, first cannot be negative")
			logrus.Error(err)
			return nil, err
		}
		// Forward Pagination
		opts = options.Find().SetLimit(int64(*input.First)).SetSort(bson.M{"createTime": direction})
		if input.After != nil {
			paginationCursor, _ := paginationUtil.DecodeCursor(*input.After)
			createtime, _ := time.Parse(time.RFC3339, paginationCursor)
			filter["createTime"] = bson.M{"$gt": createtime}
		}
	} else if input.Last != nil {
		if *input.Last < 0 {
			err := fmt.Errorf("invalid parameter, last cannot be negative")
			logrus.Error(err)
			return nil, err
		}
		// Backward Pagination
		direction = -1
		opts = options.Find().SetLimit(int64(*input.Last)).SetSort(bson.M{"createTime": direction})
		if input.Before != nil {
			paginationCursor, _ := paginationUtil.DecodeCursor(*input.After)
			filter["createTime"] = bson.M{"$lt": paginationCursor}
		}
	}

	collectionName := collection.User{}.CollectionName()
	cursor, err := r.db.Collection(collectionName).Find(ctx, filter, opts)
	defer func() {
		if cursor != nil {
			cursor.Close(context.Background())
		}
	}()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	users := []collection.User{}
	err = cursor.All(ctx, &users)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	result := &model.UserConnection{
		Edges:    []*model.UserEdge{},
		PageInfo: &model.PageInfo{},
	}
	for _, p := range users {
		edge := model.UserEdge{
			Cursor: paginationUtil.EncodeCursor(p.CreateTime.Format(time.RFC3339)),
			Node: &model.User{
				ID:         p.ID.Hex(),
				Email:      p.Email,
				FullName:   p.FullName,
				Role:       p.Role,
				CreateTime: p.CreateTime,
				UpdateTime: p.UpdateTime,
			},
		}
		if direction == 1 {
			result.Edges = append(result.Edges, &edge)
		} else {
			result.Edges = append([]*model.UserEdge{&edge}, result.Edges...)
		}
	}

	if len(users) > 0 {
		result.PageInfo.StartCursor = &result.Edges[0].Cursor
		result.PageInfo.EndCursor = &result.Edges[len(result.Edges)-1].Cursor

		countOpts := options.Count().SetLimit(1)
		nextCount, err := r.db.Collection(collectionName).CountDocuments(
			ctx,
			bson.M{"createTime": bson.M{"$gt": users[len(users)-1].CreateTime}},
			countOpts,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result.PageInfo.HasNextPage = nextCount > 0
		prevCount, err := r.db.Collection(collectionName).CountDocuments(
			ctx,
			bson.M{"createTime": bson.M{"$lt": users[0].CreateTime}},
			countOpts,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result.PageInfo.HasNextPage = prevCount > 0
	}

	return result, nil
}
