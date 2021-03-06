// Copyright 2017 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flare

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// Subscription is used to notify the clients about changes on documents.
type Subscription struct {
	ID        string
	Endpoint  SubscriptionEndpoint
	Delivery  SubscriptionDelivery
	Resource  Resource
	Data      map[string]interface{}
	CreatedAt time.Time
}

// SubscriptionEndpoint has the address information to notify the clients.
type SubscriptionEndpoint struct {
	URL     url.URL
	Method  string
	Headers http.Header
}

// SubscriptionDelivery is used to control whenever the notification can be considered successful
// or not.
type SubscriptionDelivery struct {
	Success []int
	Discard []int
}

// All kinds of actions a subscription trigger supports.
const (
	SubscriptionTriggerCreate = "create"
	SubscriptionTriggerUpdate = "update"
	SubscriptionTriggerDelete = "delete"
)

// SubscriptionRepositorier is used to interact with the Subscription data storage.
type SubscriptionRepositorier interface {
	FindAll(context.Context, *Pagination, string) ([]Subscription, *Pagination, error)
	FindOne(ctx context.Context, resourceId, id string) (*Subscription, error)
	Create(context.Context, *Subscription) error
	Delete(ctx context.Context, resourceId, id string) error
	HasSubscription(ctx context.Context, resourceId string) (bool, error)
	Trigger(
		ctx context.Context,
		action string,
		document *Document,
		fn func(context.Context, Subscription, string) error,
	) error
}

// SubscriptionTrigger is used to trigger the change on Documents.
type SubscriptionTrigger interface {
	Update(ctx context.Context, document *Document) error
	Delete(ctx context.Context, document *Document) error
}

// SubscriptionRepositoryError implements all the errrors the repository can return.
type SubscriptionRepositoryError interface {
	NotFound() bool
	AlreadyExists() bool
}
