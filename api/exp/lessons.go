package api

import (
	"context"
	"errors"
	"time"

	copier "github.com/jinzhu/copier"
	ot "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	log "github.com/sirupsen/logrus"

	pb "github.com/nre-learning/antidote-core/api/exp/generated"
	models "github.com/nre-learning/antidote-core/db/models"
)

// ListLessons returns a list of Lessons present in the data store
func (s *AntidoteAPI) ListLessons(ctx context.Context, filter *pb.LessonFilter) (*pb.Lessons, error) {
	span := ot.StartSpan("api_lesson_list", ext.SpanKindRPCClient)
	defer span.Finish()

	lessons := []*pb.Lesson{}

	dbLessons, err := s.Db.ListLessons(span.Context())
	if err != nil {
		log.Error(err)
		return nil, errors.New("Error retrieving lessons")
	}

	for _, l := range dbLessons {
		lessons = append(lessons, lessonDBToAPI(&l))
	}

	return &pb.Lessons{
		Lessons: lessons,
	}, nil
}

// GetAllLessonPrereqs examines the entire tree of depedencies that a given lesson might have, and returns
// it as a flattened, de-duplicated list. Used for the advisor's learning path tool in antidote-web
func (s *AntidoteAPI) GetAllLessonPrereqs(ctx context.Context, lessonSlug *pb.LessonSlug) (*pb.LessonPrereqs, error) {

	// Preload the requested lesson ID so we can strip it before returning
	pr := s.getPrereqs(lessonSlug.Slug, []string{lessonSlug.Slug})
	log.Debugf("Getting prerequisites for Lesson %s: %s", lessonSlug.Slug, pr)

	return &pb.LessonPrereqs{
		// Remove first item from slice - this is the lesson being requested
		Prereqs: pr[1:],
	}, nil
}

func (s *AntidoteAPI) getPrereqs(lessonSlug string, currentPrereqs []string) []string {
	span := ot.StartSpan("api_lesson_getprereqs", ext.SpanKindRPCClient)
	defer span.Finish()

	// Return if lesson slug doesn't exist
	lesson, err := s.Db.GetLesson(span.Context(), lessonSlug)
	if err != nil {
		return currentPrereqs
	}

	// Add this lessonSlug to prereqs if doesn't already exist
	if !isAlreadyInSlice(lessonSlug, currentPrereqs) {
		currentPrereqs = append(currentPrereqs, lessonSlug)
	}

	// Return if lesson doesn't have prerequisites
	if len(lesson.Prereqs) == 0 {
		return currentPrereqs
	}

	// Call recursion for lesson IDs that need it
	for i := range lesson.Prereqs {
		prereqSlug := lesson.Prereqs[i]
		currentPrereqs = s.getPrereqs(prereqSlug, currentPrereqs)
	}

	return currentPrereqs
}

func isAlreadyInSlice(lessonSlug string, currentPrereqs []string) bool {
	for i := range currentPrereqs {
		if currentPrereqs[i] == lessonSlug {
			return true
		}
	}
	return false
}

// GetLesson retrieves a single Lesson from the data store by Slug
func (s *AntidoteAPI) GetLesson(ctx context.Context, lessonSlug *pb.LessonSlug) (*pb.Lesson, error) {
	span := ot.StartSpan("api_lesson_get", ext.SpanKindRPCClient)
	defer span.Finish()

	time.Sleep(2 * time.Second)
	dbLesson, err := s.Db.GetLesson(span.Context(), lessonSlug.Slug)
	if err != nil {
		log.Error(err)
		return nil, errors.New("Error retrieving specified lesson")
	}

	lesson := lessonDBToAPI(dbLesson)

	return lesson, nil
}

// lessonDBToAPI translates a single Lesson from the `db` package models into the
// api package's equivalent
func lessonDBToAPI(dbLesson *models.Lesson) *pb.Lesson {
	lessonAPI := &pb.Lesson{}
	copier.Copy(&lessonAPI, dbLesson)
	return lessonAPI
}

// lessonAPIToDB translates a single Lesson from the `api` package models into the
// `db` package's equivalent
func lessonAPIToDB(pbLesson *pb.Lesson) *models.Lesson {
	lessonDB := &models.Lesson{}
	copier.Copy(&pbLesson, lessonDB)
	return lessonDB
}
