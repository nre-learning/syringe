package api

import (
	"context"
	"errors"

	copier "github.com/jinzhu/copier"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	log "github.com/sirupsen/logrus"

	pb "github.com/nre-learning/antidote-core/api/exp/generated"
	models "github.com/nre-learning/antidote-core/db/models"
)

// ListCollections returns a list of Collections present in the data store
func (s *AntidoteAPI) ListCollections(ctx context.Context, _ *pb.CollectionFilter) (*pb.Collections, error) {
	span := opentracing.StartSpan("api_collection_list", ext.SpanKindRPCClient)
	defer span.Finish()

	collections := []*pb.Collection{}

	dbCollections, err := s.Db.ListCollections(span.Context())
	if err != nil {
		log.Error(err)
		return nil, errors.New("Error retrieving specified collection")
	}

	for _, c := range dbCollections {
		collections = append(collections, collectionDBToAPI(&c))
	}

	return &pb.Collections{
		Collections: collections,
	}, nil
}

// GetCollection retrieves a single Collection from the data store by Slug
func (s *AntidoteAPI) GetCollection(ctx context.Context, collectionSlug *pb.CollectionSlug) (*pb.Collection, error) {

	span := opentracing.StartSpan("api_collection_get", ext.SpanKindRPCClient)
	defer span.Finish()

	dbCollection, err := s.Db.GetCollection(span.Context(), collectionSlug.Slug)
	if err != nil {
		log.Error(err)
		return nil, errors.New("Error retrieving specified collection")
	}

	collection := collectionDBToAPI(dbCollection)

	lessons, err := s.Db.ListLessons(span.Context())
	if err != nil {
		log.Errorf("Error retrieving lessons for collection %s: %v", dbCollection.Slug, err)
		return nil, errors.New("Error retrieving specified collection")
	}

	for lessonSlug, lesson := range lessons {
		if lesson.Collection == collectionSlug.Slug {
			collection.Lessons = append(collection.Lessons, &pb.LessonSummary{
				Slug:        lessonSlug,
				Description: lesson.Description,
				Name:        lesson.Name,
			})
		}
	}

	return collection, nil
}

// collectionDBToAPI translates a single Collection from the `db` package models into the
// api package's equivalent
func collectionDBToAPI(dbCollection *models.Collection) *pb.Collection {
	collectionAPI := &pb.Collection{}
	copier.Copy(&collectionAPI, dbCollection)
	return collectionAPI
}

// collectionAPIToDB translates a single Collection from the `api` package models into the
// `db` package's equivalent
func collectionAPIToDB(pbCollection *pb.Collection) *models.Collection {
	collectionDB := &models.Collection{}
	copier.Copy(&pbCollection, collectionDB)
	return collectionDB
}
