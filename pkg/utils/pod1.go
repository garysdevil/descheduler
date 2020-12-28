package utils

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// GetResourceActualQuantity finds and returns the actual quantity for a specific resource.
func GetResourceActualQuantity(pod *v1.Pod, resourceName v1.ResourceName, metricsClient metricsv.Interface) resource.Quantity {
	actualQuantity := resource.Quantity{}

	switch resourceName {
	case v1.ResourceCPU:
		actualQuantity = resource.Quantity{Format: resource.DecimalSI}
	case v1.ResourceMemory, v1.ResourceStorage, v1.ResourceEphemeralStorage:
		actualQuantity = resource.Quantity{Format: resource.BinarySI}
	default:
		actualQuantity = resource.Quantity{Format: resource.DecimalSI}
	}

	podname := pod.GetName()
	podnamesapce := pod.GetNamespace()
	nodeMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(podnamesapce).Get(context.Background(), podname, metav1.GetOptions{})
	if err != nil {
		klog.V(1).InfoS("Error. Get metric failed", "Pod", pod.GetName(), "Error", err)
		return resource.MustParse("0")
	}
	for _, containerMetrics := range nodeMetrics.Containers {
		usage := containerMetrics.Usage

		switch resourceName {
		case v1.ResourceCPU:
			actualQuantity.Add(*usage.Cpu())
		case v1.ResourceMemory, v1.ResourceStorage, v1.ResourceEphemeralStorage:
			actualQuantity.Add(*usage.Memory())
		default:
			klog.V(1).InfoS("Error GetResourceActualQuantity")
		}
	}

	return actualQuantity
}
