/*
Copyright 2014 Google Inc. All rights reserved.

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

package json

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta1"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta2"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	testCases := []struct {
		Type   watch.EventType
		Object runtime.Object
		Codec  runtime.Codec
	}{
		{
			watch.Added,
			&api.Pod{ObjectMeta: api.ObjectMeta{Name: "foo"}},
			v1beta1.Codec,
		},
		{
			watch.Modified,
			&api.Pod{ObjectMeta: api.ObjectMeta{Name: "foo"}},
			v1beta2.Codec,
		},
		{
			watch.Deleted,
			&api.Pod{ObjectMeta: api.ObjectMeta{Name: "foo"}},
			api.Codec,
		},
	}
	for i, testCase := range testCases {
		buf := &bytes.Buffer{}

		encoder := NewEncoder(buf, testCase.Codec)
		if err := encoder.Encode(&watch.Event{Type: testCase.Type, Object: testCase.Object}); err != nil {
			t.Errorf("%d: unexpected error: %v", i, err)
			continue
		}

		decoder := NewDecoder(ioutil.NopCloser(buf), testCase.Codec)
		event, obj, err := decoder.Decode()
		if err != nil {
			t.Errorf("%d: unexpected error: %v", i, err)
			continue
		}
		if !reflect.DeepEqual(testCase.Object, obj) {
			t.Errorf("%d: expected %#v, got %#v", i, testCase.Object, obj)
		}
		if event != testCase.Type {
			t.Errorf("%d: unexpected type: %#v", i, event)
		}
	}
}
