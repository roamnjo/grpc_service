package auth

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       string `bson:"_id, omitempty"`
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	FindEmail(ctx context.Context, email string) error
	FindSameName(ctx context.Context, name string) error
}

type repository struct {
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{coll: db.Collection("users")}
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	_, err := r.coll.InsertOne(ctx, user)
	return err
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindEmail(ctx context.Context, email string) error {
	filter := bson.M{"email": email}

	searchOne := r.coll.FindOne(ctx, filter)
	if searchOne.Err() == nil {
		return fmt.Errorf("Error: email already exists")
	}
	return nil
}

func (r *repository) FindSameName(ctx context.Context, name string) error {
	filter := bson.M{"name": name}

	searchOne := r.coll.FindOne(ctx, filter)
	if searchOne.Err() == nil {
		return fmt.Errorf("Error: user with this name already exists")
	}
	return nil
}
