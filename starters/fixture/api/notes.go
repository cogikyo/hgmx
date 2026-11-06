package api

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"connectrpc.com/connect/v2"

	"github.com/cogikyo/hgmx/auth"
	"github.com/cogikyo/hgmx/pg"

	apiv1 "example.com/fixture/proto/fixture/api/v1"
	"example.com/fixture/proto/fixture/api/v1/apiv1connect"
)

const limit = 20

type Notes struct {
	mu    sync.Mutex
	notes map[pg.ID][]*apiv1.Note
}

func Register(api *auth.API) {
	n := &Notes{notes: map[pg.ID][]*apiv1.Note{}}
	apiv1connect.RegisterNoteServiceHandler(api.Server(), n)
	api.Allow(auth.Authenticated, apiv1connect.NoteServiceListNotesProcedure, apiv1connect.NoteServiceAddNoteProcedure)
}

func (n *Notes) ListNotes(ctx context.Context, _ *apiv1.ListNotesRequest) (*apiv1.ListNotesResponse, error) {
	u, _ := auth.Caller(ctx)
	n.mu.Lock()
	defer n.mu.Unlock()
	return &apiv1.ListNotesResponse{Notes: n.snapshot(u.ID)}, nil
}

func (n *Notes) AddNote(ctx context.Context, req *apiv1.AddNoteRequest) (*apiv1.AddNoteResponse, error) {
	u, _ := auth.Caller(ctx)
	text := strings.TrimSpace(req.GetText())
	if text == "" || len(text) > 500 {
		return nil, connect.NewError(connect.CodeInvalidArgument, "A note needs 1 to 500 characters.")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if len(n.notes[u.ID]) >= limit {
		return nil, connect.NewError(connect.CodeFailedPrecondition, "You have the most notes allowed.")
	}
	n.notes[u.ID] = append(n.notes[u.ID], &apiv1.Note{Id: strconv.Itoa(len(n.notes[u.ID]) + 1), Text: text})
	return &apiv1.AddNoteResponse{Notes: n.snapshot(u.ID)}, nil
}

func (n *Notes) snapshot(user pg.ID) *apiv1.Notes {
	return &apiv1.Notes{Notes: n.notes[user], CanAdd: len(n.notes[user]) < limit}
}
