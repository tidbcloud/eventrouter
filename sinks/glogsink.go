/*
Copyright 2017 Heptio Inc.

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

package sinks

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
)

// GlogSink is the most basic sink
// Useful when you already have ELK/EFK Stack
type GlogSink struct {
	// TODO: create a channel and buffer for scaling
}

// NewGlogSink will create a new
func NewGlogSink() EventSinkInterface {
	return &GlogSink{}
}

// UpdateEvents implements the EventSinkInterface
func (gs *GlogSink) UpdateEvents(eNew *v1.Event, eOld *v1.Event) {
	eData := NewEventData(eNew, eOld)
	obj := eData.Event
	firstTimestamp := obj.FirstTimestamp.Time
	if obj.FirstTimestamp.IsZero() {
		firstTimestamp = obj.EventTime.Time
	}

	lastTimestamp := obj.LastTimestamp.Time
	if obj.LastTimestamp.IsZero() {
		lastTimestamp = firstTimestamp
	}

	var target string
	if len(obj.InvolvedObject.Name) > 0 {
		target = fmt.Sprintf("%s/%s", strings.ToLower(obj.InvolvedObject.Kind), obj.InvolvedObject.Name)
	} else {
		target = strings.ToLower(obj.InvolvedObject.Kind)
	}

	fmt.Printf("%v\t%v\t%v\ttarget=%v\tsource=%v\n", lastTimestamp, obj.Type, strings.TrimSpace(obj.Message), target, formatEventSource(obj.Source, obj.ReportingController, obj.ReportingInstance))
	return
}

// formatEventSource formats EventSource as a comma separated string excluding Host when empty.
// It uses reportingController when Source.Component is empty and reportingInstance when Source.Host is empty
func formatEventSource(es v1.EventSource, reportingController, reportingInstance string) string {
	return formatEventSourceComponentInstance(
		firstNonEmpty(es.Component, reportingController),
		firstNonEmpty(es.Host, reportingInstance),
	)
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if len(s) > 0 {
			return s
		}
	}
	return ""
}

func formatEventSourceComponentInstance(component, instance string) string {
	if len(instance) == 0 {
		return component
	}
	return component + ", " + instance
}
