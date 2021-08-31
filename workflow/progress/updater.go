package progress

import (
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v3/workflow/common"
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

func UpdateProgress(podInformer cache.SharedIndexInformer, wf *wfv1.Workflow, log *log.Entry) bool {
	updated := false
	wf.Status.Progress = "0/0"
	for nodeID, node := range wf.Status.Nodes {
		if node.Type != wfv1.NodeTypePod {
			continue
		}
		progress := podProgress(podInformer, wf.Namespace, nodeID, "0/1", log)
		if node.Fulfilled() {
			progress = "1/1"
		}
		if progress.IsValid() && node.Progress != progress {
			log.WithField("progress", progress).Info("pod progress")
			node.Progress = progress
			wf.Status.Nodes[nodeID] = node
			updated = true
		}
		wf.Status.Progress = wf.Status.Progress.Add(progress)
	}
	for nodeID, node := range wf.Status.Nodes {
		if node.Type == wfv1.NodeTypePod {
			continue
		}
		progress := sumProgress(wf, node, make(map[string]bool))
		if progress.IsValid() && node.Progress != progress {
			node.Progress = progress
			wf.Status.Nodes[nodeID] = node
			updated = true
		}
	}
	return updated
}

func sumProgress(wf *wfv1.Workflow, node wfv1.NodeStatus, visited map[string]bool) wfv1.Progress {
	progress := wfv1.Progress("0/0")
	for _, childNodeID := range node.Children {
		if visited[childNodeID] {
			continue
		}
		visited[childNodeID] = true
		// this will tolerate missing child (will be "") and therefore ignored
		child := wf.Status.Nodes[childNodeID]
		progress = progress.Add(sumProgress(wf, child, visited))
		if child.Type == wfv1.NodeTypePod {
			v := child.Progress
			if v.IsValid() {
				progress = progress.Add(v)
			}
		}
	}
	return progress
}

func podProgress(podInformer cache.SharedIndexInformer, namespace, name string, progress wfv1.Progress, log *log.Entry) wfv1.Progress {
	// for pods, lets see what the annotation says pod can get deleted of course, so
	// can be empty and return "", even it previously had a value
	obj, _, _ := podInformer.GetStore().GetByKey(namespace + "/" + name)
	if pod, ok := obj.(*apiv1.Pod); ok {
		log.WithField("annotations", pod.Annotations).Info()
		if annotation, ok := pod.Annotations[common.AnnotationKeyProgress]; ok {
			v, ok := wfv1.ParseProgress(annotation)
			if ok {
				return v
			}
		}
	}
	return progress
}
