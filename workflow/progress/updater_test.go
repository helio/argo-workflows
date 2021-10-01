package progress

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v3/workflow/common"
)

func TestUpdater(t *testing.T) {
	ns := "my-ns"
	wf := &wfv1.Workflow{
		ObjectMeta: metav1.ObjectMeta{Namespace: ns},
		Status: wfv1.WorkflowStatus{
			Nodes: wfv1.Nodes{
				"pod-1": wfv1.NodeStatus{Phase: wfv1.NodeSucceeded, Type: wfv1.NodeTypePod},
				"pod-2": wfv1.NodeStatus{Type: wfv1.NodeTypePod},
				"pod-3": wfv1.NodeStatus{Type: wfv1.NodeTypePod},
				"wf":    wfv1.NodeStatus{Children: []string{"pod-1", "pod-2", "pod-3"}},
			},
		},
	}
	pods := []*apiv1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns,
				Name:      "pod-1",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns,
				Name:      "pod-2",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns,
				Name:      "pod-3",
				Annotations: map[string]string{
					common.AnnotationKeyProgress: "50/100",
				},
			},
		},
	}
	UpdateProgress(pods, wf, logrus.NewEntry(logrus.New()))
	assert.Equal(t, wfv1.Progress("1/1"), wf.Status.Nodes["pod-1"].Progress)
	assert.Equal(t, wfv1.Progress("0/1"), wf.Status.Nodes["pod-2"].Progress)
	assert.Equal(t, wfv1.Progress("50/100"), wf.Status.Nodes["pod-3"].Progress)
	assert.Equal(t, wfv1.Progress("51/102"), wf.Status.Nodes["wf"].Progress)
	assert.Equal(t, wfv1.Progress("51/102"), wf.Status.Progress)
}
