// apcore is a server framework for implementing an ActivityPub application.
// Copyright (C) 2019 Cory Slep
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package apcore

import (
	"context"
	"net/url"

	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams/vocab"
)

var _ pub.Database = &apdb{}

type apdb struct {
	db *database
}

func newApdb(db *database) *apdb {
	return &apdb{
		db: db,
	}
}

func (a *apdb) Lock(c context.Context, id *url.URL) error {
	// TODO
	return nil
}

func (a *apdb) Unlock(c context.Context, id *url.URL) error {
	// TODO
	return nil
}

func (a *apdb) InboxContains(c context.Context, inbox, id *url.URL) (contains bool, err error) {
	// TODO
	return
}

func (a *apdb) GetInbox(c context.Context, inboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	// TODO
	return
}

func (a *apdb) SetInbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
	// TODO
	return nil
}

func (a *apdb) Owns(c context.Context, id *url.URL) (owns bool, err error) {
	// TODO
	return
}

func (a *apdb) ActorForOutbox(c context.Context, outboxIRI *url.URL) (actorIRI *url.URL, err error) {
	// TODO
	return
}

func (a *apdb) ActorForInbox(c context.Context, inboxIRI *url.URL) (actorIRI *url.URL, err error) {
	// TODO
	return
}

func (a *apdb) OutboxForInbox(c context.Context, inboxIRI *url.URL) (outboxIRI *url.URL, err error) {
	// TODO
	return
}

func (a *apdb) Exists(c context.Context, id *url.URL) (exists bool, err error) {
	// TODO
	return
}

func (a *apdb) Get(c context.Context, id *url.URL) (value vocab.Type, err error) {
	// TODO
	return
}

func (a *apdb) Create(c context.Context, asType vocab.Type) error {
	// TODO
	return nil
}

func (a *apdb) Update(c context.Context, asType vocab.Type) error {
	// TODO
	return nil
}

func (a *apdb) Delete(c context.Context, id *url.URL) error {
	// TODO
	return nil
}

func (a *apdb) GetOutbox(c context.Context, inboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	// TODO
	return
}

func (a *apdb) SetOutbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
	// TODO
	return nil
}

func (a *apdb) NewId(c context.Context, t vocab.Type) (id *url.URL, err error) {
	// TODO
	return
}

func (a *apdb) Followers(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	// TODO
	return
}

func (a *apdb) Following(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	// TODO
	return
}

func (a *apdb) Liked(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	// TODO
	return
}