package main

import (
	"context"
	//"errors"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

// Define the GraphQL schema
var schemaString = `
    schema {
        query: Query
        mutation: Mutation
    }
    
    type Query {
        Remotemeeting(id: ID!): RemoteMeeting
        remotemeetings: [RemoteMeeting]
    }
    
    type Mutation {
        createMeeting(input: RemoteMeetingInput!): RemoteMeeting
        updateMeeting(id: ID!, input: RemoteMeetingInput!): RemoteMeeting
        deleteMeeting(id: ID!): Boolean
    }
    
    type RemoteMeeting {
        id: ID!
        title: String!
        description: String
        timeLimit: Int!
        timeZone: String!
		participants: String!
		Email: String!
        scheduledTime: DateTime!
        reminderTime: DateTime!
    }
    
    input RemoteMeetingInput {
        title: String!
        description: String
        timeLimit: Int!
        timeZone: String!
		Participants: String!
		Email: String!
        scheduledTime: DateTime!
        reminderTime: DateTime!
    }
    
    scalar DateTime
`

// Define the Meeting type and database operations
type RemoteMeeting struct {
	ID            graphql.ID
	Title         string
	Description   string
	TimeLimit     int
	TimeZone      string
	Participants  string
	Email         string
	ScheduledTime time.Time
	ReminderTime  time.Time
}

var remotemeetings = map[graphql.ID]*RemoteMeeting{}

func (m *RemoteMeeting) Save() {
	remotemeetings[m.ID] = m
}

func (m *RemoteMeeting) Delete() {
	delete(remotemeetings, m.ID)
}

// Define the resolvers for the Query and Mutation types
type Resolver struct{}

func (r *Resolver) RemoteMeeting(ctx context.Context, args struct{ ID graphql.ID }) *RemoteMeetingResolver {
	if remotemeeting, ok := remotemeetings[args.ID]; ok {
		return &RemoteMeetingResolver{meeting}
	}
	return nil
}

func (r *Resolver) RemoteMeetings(ctx context.Context) []*RemoteMeetingResolver {
	var result []*RemoteMeetingResolver
	for _, remotemeeting := range remotemeetings {
		result = append(result, &RemoteMeetingResolver{meeting})
	}
	return result
}

func (r *Resolver) CreateMeeting(ctx context.Context, args struct{ Input *RemoteMeetingInput }) *RemoteMeetingResolver {
	Remotemeeting := &RemoteMeeting{
		ID:            graphql.ID(fmt.Sprintf("m%d", len(meetings))),
		Title:         args.Input.Title,
		Description:   args.Input.Description,
		Participants:  args.Input.Participants,
		Email:         args.Input.Email,
		TimeLimit:     args.Input.TimeLimit,
		TimeZone:      args.Input.TimeZone,
		ScheduledTime: args.Input.ScheduledTime,
		ReminderTime:  args.Input.ReminderTime,
	}
	Remotemeeting.Save()
	return &RemoteMeetingResolver{meeting}
}

func (r *Resolver) UpdateMeeting(ctx context.Context, args struct{ Update *RemoteMeetingUpdate }) *RemoteMeetingResolver {

}
