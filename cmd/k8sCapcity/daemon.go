package main

import (
	"encoding/json"
	"fmt"
	resource "k8s.io/apimachinery/pkg/api/resource"
)

func runDaemonMode(clusterInfo ClusterInfo) {
	daemonLog := DaemonLog{}
	daemonLog.UtilizationFactorPods = make(map[string]float64)
	daemonLog.UtilizationFactorMemoryRequests = make(map[string]float64)
	daemonLog.UtilizationFactorCPURequests = make(map[string]float64)

	for name, node := range clusterInfo.NodeInfo {
		if node.PrintOutput {
			clusterInfo.ClusterUsedCPURequests.Add(node.UsedCPURequests)
			clusterInfo.ClusterUsedCPU.Add(node.UsedCPU)
			clusterInfo.ClusterUsedMemoryRequests.Add(node.UsedMemoryRequests)
			clusterInfo.ClusterUsedMemory.Add(node.UsedMemory)
			clusterInfo.ClusterUsedPods = clusterInfo.ClusterUsedPods + node.UsedPods
			clusterInfo.ClusterUsedMemoryLimits.Add(node.UsedMemoryLimits)
			daemonLog.UtilizationFactorPods[name] = float64(node.UsedPods) / float64(node.AllocatablePods.Value())
			daemonLog.UtilizationFactorMemoryRequests[name] = float64(node.UsedMemoryRequests.Value()) / float64(node.AllocatableMemory.Value())
			daemonLog.UtilizationFactorCPURequests[name] = float64(node.UsedCPURequests.Value()) / float64(node.AllocatableCPU.Value())

		}
	}

	daemonLog.EventKind = "metric"
	daemonLog.EventModule = "k8s_quota"
	daemonLog.EventProvider = "k8sCapcity"
	daemonLog.EventType = "info"
	daemonLog.EventVersion = "03/04/2020-01"
	daemonLog.NodeLabel = clusterInfo.NodeLabel
	daemonLog.ResourceQuotaCPURequestCores = clusterInfo.RqclusterAllocatedRequestsCPU.Value()
	daemonLog.ResourceQuotaCPURequestMilliCores = clusterInfo.RqclusterAllocatedRequestsCPU.ScaledValue(resource.Milli)
	daemonLog.ResourceQuotaMemoryRequest = clusterInfo.RqclusterAllocatedRequestsMemory.Value()
	daemonLog.ResourceQuotaMemoryLimit = clusterInfo.RqclusterAllocatedLimitsMemory.Value()
	daemonLog.ResourceQuotaPods = clusterInfo.RqclusterAllocatedPods.Value()
	daemonLog.ContainerResourceCPURequestCores = clusterInfo.ClusterUsedCPURequests.Value()
	daemonLog.ContainerResourceCPURequestMilliCores = clusterInfo.ClusterUsedCPURequests.ScaledValue(resource.Milli)
	daemonLog.ContainerResourceMemoryRequest = clusterInfo.ClusterUsedMemoryRequests.Value()
	daemonLog.ContainerResourceMemoryLimit = clusterInfo.ClusterUsedMemoryLimits.Value()
	daemonLog.ContainerResourcePods = clusterInfo.ClusterUsedPods
	daemonLog.AllocatableMemory = clusterInfo.ClusterAllocatableMemory.Value()
	daemonLog.AllocatableMemoryNminusone = clusterInfo.ClusterAllocatableMemory.Value() - clusterInfo.NminusMemory.Value()
	daemonLog.AllocatableCPU = clusterInfo.ClusterAllocatableCPU.Value()
	daemonLog.AllocatableCPUNminusone = clusterInfo.ClusterAllocatableCPU.Value() - clusterInfo.NminusCPU.Value()
	daemonLog.AllocatablePods = clusterInfo.ClusterAllocatablePods.Value()
	daemonLog.AllocatablePodsNminusone = clusterInfo.ClusterAllocatablePods.Value() - clusterInfo.NminusPods.Value()
	daemonLog.SubscriptionFactorMemoryRequest = float64(daemonLog.ResourceQuotaMemoryRequest) / float64(daemonLog.AllocatableMemory)
	daemonLog.SubscriptionFactorMemoryRequestNminusone = float64(daemonLog.ResourceQuotaMemoryRequest) / float64(daemonLog.AllocatableMemoryNminusone)
	daemonLog.SubscriptionFactorCPURequest = float64(daemonLog.ResourceQuotaCPURequestMilliCores) / float64(clusterInfo.ClusterAllocatableCPU.ScaledValue(resource.Milli))
	daemonLog.SubscriptionFactorCPURequestNminusone = float64(daemonLog.ResourceQuotaCPURequestMilliCores) / float64(clusterInfo.ClusterAllocatableCPU.ScaledValue(resource.Milli)-clusterInfo.NminusCPU.ScaledValue(resource.Milli))
	daemonLog.SubscriptionFactorPods = float64(daemonLog.ResourceQuotaPods) / float64(daemonLog.AllocatablePods)
	daemonLog.SubscriptionFactorPodsNminusone = float64(daemonLog.ResourceQuotaPods) / float64(daemonLog.AllocatablePodsNminusone)
	daemonLog.UtilizationFactorPodsTotal = float64(clusterInfo.ClusterUsedPods) / float64(daemonLog.AllocatablePods)
	daemonLog.UtilizationFactorPodsTotalNminusone = float64(clusterInfo.ClusterUsedPods) / float64(daemonLog.AllocatablePodsNminusone)
	daemonLog.UtilizationFactorMemoryRequestsTotal = float64(daemonLog.ContainerResourceMemoryRequest) / float64(daemonLog.AllocatableMemory)
	daemonLog.UtilizationFactorMemoryRequestsTotalNminusone = float64(daemonLog.ContainerResourceMemoryRequest) / float64(daemonLog.AllocatableMemoryNminusone)
	daemonLog.UtilizationFactorCPURequestsTotal = float64(clusterInfo.ClusterUsedCPURequests.Value()) / float64(daemonLog.AllocatableCPU)
	daemonLog.UtilizationFactorCPURequestsTotalNminusone = float64(clusterInfo.ClusterUsedCPURequests.Value()) / float64(daemonLog.AllocatableCPUNminusone)
	daemonLog.AvailableMemoryRequest = daemonLog.AllocatableMemory - daemonLog.ContainerResourceMemoryRequest
	daemonLog.AvailableMemoryRequestNminusone = daemonLog.AllocatableMemoryNminusone - daemonLog.ContainerResourceMemoryRequest
	daemonLog.AvailableCPURequest = daemonLog.AllocatableCPU - daemonLog.ContainerResourceCPURequestCores
	daemonLog.AvailableCPURequestNminusone = daemonLog.AllocatableCPUNminusone - daemonLog.ContainerResourceCPURequestCores
	daemonLog.AvailablePods = daemonLog.AllocatablePods - daemonLog.ContainerResourcePods
	daemonLog.AvailablePodsNminusone = daemonLog.AllocatablePodsNminusone - daemonLog.ContainerResourcePods
	result, err := json.Marshal(daemonLog)
	check(err)
	fmt.Println(string(result))
}
