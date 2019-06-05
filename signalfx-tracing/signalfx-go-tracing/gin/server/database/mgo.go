package database

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	mgotrace "github.com/signalfx/signalfx-go-tracing/contrib/globalsign/mgo"
	"github.com/signalfx/signalfx-go-tracing/ddtrace/tracer"
	"github.com/signalfx/tracing-examples/signalfx-tracing/signalfx-go-tracing/gin/server/models"
)

type mgoManager struct {
	Host        string
	Port        int
	Name        string
	ServiceName string
}

var _ Manager = (*mgoManager)(nil)

// GetBoardByID returns a board for a given boardID
func (m *mgoManager) GetBoardByID(c context.Context, id string) (models.Board, error) {
	collection := m.getCollection(c)

	board := models.Board{}
	err := collection.Find(bson.M{"board_id": id}).One(&board)
	if err != nil {
		return models.Board{}, err
	}
	return board, nil
}

// InsertBoard inserts a given board
func (m *mgoManager) InsertBoard(c context.Context, board models.Board) error {
	collection := m.getCollection(c)

	return collection.Insert(board)
}

// UpdateBoard saves a given updated board
func (m *mgoManager) UpdateBoard(c context.Context, board models.Board) error {
	collection := m.getCollection(c)

	return collection.Update(bson.M{"board_id": board.ID}, board)
}

// getCollection returns board collection
func (m *mgoManager) getCollection(c context.Context) *mgo.Collection {
	parentSpan, ctx := tracer.StartSpanFromContext(c, "mongo.session", tracer.SpanType("server"), tracer.ResourceName("mongo"))
	session, err := mgotrace.Dial(fmt.Sprintf("mongodb://%s:%d/%s", m.Host, m.Port, m.Name), mgotrace.WithServiceName("signalfx-battleship"), mgotrace.WithContext(ctx))
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err.Error())
	}

	db := session.Clone().DB(m.Name)
	collection := db.C(models.CollectionBoard)

	defer parentSpan.Finish()

	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err.Error())
	}

	return collection
}
