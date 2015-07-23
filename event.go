/*
Copyright 2014 Rohith All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
)

const (
	NODE_EVENT  = 1
	POD_EVENT   = 2
)

func (r Event) String() string {
	return fmt.Sprintf("type: %s, event: %v", r.getEventType(), r.Event)
}

func (r Event) getEventType() string {
	switch r.Type {
	case NODE_EVENT:
		return "node"
	default:
		return "pod"
	}
}

func NewEvent(eventType int, event interface{}) *Event {
	return &Event{
		Type: eventType,
		Event: event,
	}
}

func (e *Event) IsNodeEvent() *Event {
	e.Type = NODE_EVENT
	return e
}

func (e *Event) IsPodEvent() *Event {
	e.Type = POD_EVENT
	return e
}
