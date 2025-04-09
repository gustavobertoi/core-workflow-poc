package workflow

import (
	"github.com/open-source-cloud/fuse/pkg/uuid"
	"github.com/rs/zerolog/log"
	"github.com/vladopajic/go-actor/actor"
)

type Engine interface {
	actor.Actor
	AddSchema(schema Schema)
	SendMessage(msg EngineMessage)
}

type engine struct {
	baseActor            actor.Actor
	externalMessagesChan chan EngineMessage
	mailbox              actor.Mailbox[any]
	schemas              map[string]Schema
	workflows            map[string]Workflow
}

func NewEngine() Engine {
	worker := &engine{
		externalMessagesChan: make(chan EngineMessage),
		mailbox:              actor.NewMailbox[any](),
		schemas:              make(map[string]Schema),
		workflows:            make(map[string]Workflow),
	}
	worker.baseActor = actor.New(worker)

	return worker
}

func (e *engine) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		for _, workflow := range e.workflows {
			workflow.Stop()
		}
		log.Info().Msg("Stopping engine")
		return actor.WorkerEnd

	case msg := <-e.externalMessagesChan:
		log.Info().Msgf("received external engineMessage: %s", msg)
		e.handleMessage(ctx, msg)
		return actor.WorkerContinue

	case msg := <-e.mailbox.ReceiveC():
		log.Info().Msgf("received engineMessage: %s", msg)
		return actor.WorkerContinue
	}
}

func (e *engine) Start() {
	e.baseActor.Start()
}

func (e *engine) Stop() {
	e.baseActor.Stop()
}

func (e *engine) AddSchema(schema Schema) {
	e.schemas[schema.ID()] = schema
}

func (e *engine) SendMessage(msg EngineMessage) {
	e.externalMessagesChan <- msg
}

func (e *engine) handleMessage(ctx actor.Context, msg EngineMessage) {
	switch msg.Type() {
	case EngineMessageStartWorkflow:
		schemaId, ok := msg.Data().(string)
		if !ok {
			log.Error().Msg("Invalid engineMessage data")
			return
		}
		workflowSchema, ok := e.schemas[schemaId]
		if !ok {
			log.Error().Msgf("Schema with ID %s not found", schemaId)
			return
		}
		newWorkflowUuid := uuid.V7()
		log.Info().Msgf("Start new workflow with ID %s from schema ID %s", newWorkflowUuid, schemaId)
		workflow := NewWorkflow(newWorkflowUuid, workflowSchema)
		workflow.Start()
		e.workflows[newWorkflowUuid] = workflow
		workflow.SendMessage(ctx, NewMessage(MessageStartWorkflow, nil))

	default:
		log.Warn().Msgf("Unhandled engine engineMessage type: %s", msg.Type())
	}
}
