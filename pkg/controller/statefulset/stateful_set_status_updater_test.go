/*
Copyright 2019 The Kruise Authors.

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

package statefulset

import (
	"context"
	"errors"
	"testing"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	core "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"

	appsv1beta1 "github.com/openkruise/kruise/apis/apps/v1beta1"
	"github.com/openkruise/kruise/pkg/client/clientset/versioned/fake"
	kruiseappslisters "github.com/openkruise/kruise/pkg/client/listers/apps/v1beta1"
)

func TestStatefulSetUpdaterUpdatesSetStatus(t *testing.T) {
	set := newStatefulSet(3)
	status := appsv1beta1.StatefulSetStatus{ObservedGeneration: 1, Replicas: 2}
	fakeClient := &fake.Clientset{}
	updater := NewRealStatefulSetStatusUpdater(fakeClient, nil)
	fakeClient.AddReactor("update", "statefulsets", func(action core.Action) (bool, runtime.Object, error) {
		update := action.(core.UpdateAction)
		return true, update.GetObject(), nil
	})
	if err := updater.UpdateStatefulSetStatus(context.TODO(), set, &status); err != nil {
		t.Errorf("Error returned on successful status update: %s", err)
	}
	if set.Status.Replicas != 2 {
		t.Errorf("UpdateStatefulSetStatus mutated the sets replicas %d", set.Status.Replicas)
	}
}

func TestStatefulSetStatusUpdaterUpdatesObservedGeneration(t *testing.T) {
	set := newStatefulSet(3)
	status := appsv1beta1.StatefulSetStatus{ObservedGeneration: 3, Replicas: 2}
	fakeClient := &fake.Clientset{}
	updater := NewRealStatefulSetStatusUpdater(fakeClient, nil)
	fakeClient.AddReactor("update", "statefulsets", func(action core.Action) (bool, runtime.Object, error) {
		update := action.(core.UpdateAction)
		sts := update.GetObject().(*appsv1beta1.StatefulSet)
		if sts.Status.ObservedGeneration != 3 {
			t.Errorf("expected observedGeneration to be synced with generation for statefulset %q", sts.Name)
		}
		return true, sts, nil
	})
	if err := updater.UpdateStatefulSetStatus(context.TODO(), set, &status); err != nil {
		t.Errorf("Error returned on successful status update: %s", err)
	}
}

func TestStatefulSetStatusUpdaterUpdateReplicasFailure(t *testing.T) {
	set := newStatefulSet(3)
	status := appsv1beta1.StatefulSetStatus{ObservedGeneration: 3, Replicas: 2}
	fakeClient := &fake.Clientset{}
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	indexer.Add(set)
	setLister := kruiseappslisters.NewStatefulSetLister(indexer)
	updater := NewRealStatefulSetStatusUpdater(fakeClient, setLister)
	fakeClient.AddReactor("update", "statefulsets", func(action core.Action) (bool, runtime.Object, error) {
		return true, nil, apierrors.NewInternalError(errors.New("API server down"))
	})
	if err := updater.UpdateStatefulSetStatus(context.TODO(), set, &status); err == nil {
		t.Error("Failed update did not return error")
	}
}

func TestStatefulSetStatusUpdaterUpdateReplicasConflict(t *testing.T) {
	set := newStatefulSet(3)
	status := appsv1beta1.StatefulSetStatus{ObservedGeneration: 3, Replicas: 2}
	conflict := false
	fakeClient := &fake.Clientset{}
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	indexer.Add(set)
	setLister := kruiseappslisters.NewStatefulSetLister(indexer)
	updater := NewRealStatefulSetStatusUpdater(fakeClient, setLister)
	fakeClient.AddReactor("update", "statefulsets", func(action core.Action) (bool, runtime.Object, error) {
		update := action.(core.UpdateAction)
		if !conflict {
			conflict = true
			return true, update.GetObject(), apierrors.NewConflict(action.GetResource().GroupResource(), set.Name, errors.New("Object already exists"))
		}
		return true, update.GetObject(), nil
	})
	if err := updater.UpdateStatefulSetStatus(context.TODO(), set, &status); err != nil {
		t.Errorf("UpdateStatefulSetStatus returned an error: %s", err)
	}
	if set.Status.Replicas != 2 {
		t.Errorf("UpdateStatefulSetStatus mutated the sets replicas %d", set.Status.Replicas)
	}
}

func TestStatefulSetStatusUpdaterUpdateReplicasConflictFailure(t *testing.T) {
	set := newStatefulSet(3)
	status := appsv1beta1.StatefulSetStatus{ObservedGeneration: 3, Replicas: 2}
	fakeClient := &fake.Clientset{}
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	indexer.Add(set)
	setLister := kruiseappslisters.NewStatefulSetLister(indexer)
	updater := NewRealStatefulSetStatusUpdater(fakeClient, setLister)
	fakeClient.AddReactor("update", "statefulsets", func(action core.Action) (bool, runtime.Object, error) {
		update := action.(core.UpdateAction)
		return true, update.GetObject(), apierrors.NewConflict(action.GetResource().GroupResource(), set.Name, errors.New("Object already exists"))
	})
	if err := updater.UpdateStatefulSetStatus(context.TODO(), set, &status); err == nil {
		t.Error("UpdateStatefulSetStatus failed to return an error on get failure")
	}
}
