/*
Copyright 2025 CodeFuture Authors

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

package kube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	coordv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/ptr"
	"time"
)

type Option struct {
	clientSet *kubernetes.Clientset
}

type Clients interface {
	NodeGetter(ctx context.Context, nodeName string) (*corev1.Node, error)
	NodePatchAnnotations(ctx context.Context, name string, annotations map[string]string) error
	AcquireLease(ctx context.Context, leaseName, namespace, holderIdentity string, ttl int32) (bool, error)
	CreateLease(ctx context.Context, leaseName, namespace, holderIdentity string, ttl int32) (bool, error)
	ReleaseLease(ctx context.Context, leaseName, namespace string) error
	DeleteLease(ctx context.Context, leaseName, namespace string) error
	VersionInfo() (*version.Info, error)
	Version() *string
}

func NewKubeClients(clientSet *kubernetes.Clientset) Clients {
	return &Option{
		clientSet: clientSet,
	}
}

func (c *Option) NodeGetter(ctx context.Context, nodeName string) (*corev1.Node, error) {
	return c.clientSet.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
}

func (c *Option) NodePatchAnnotations(ctx context.Context, name string, annotations map[string]string) error {
	if name == "" {
		return errors.New("node name is empty")
	}
	// 构造 Patch 数据
	patchData := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": annotations,
		},
	}

	patchBytes, err := json.Marshal(patchData)
	if err != nil {
		return fmt.Errorf("failed to marshal patch data: %w", err)
	}

	// 执行Patch请求
	// 使用 types.MergePatchType,它是 JSON Merge Patch，适合部分更新
	_, err = c.clientSet.CoreV1().Nodes().Patch(ctx, name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("failed to patch annotations for node %s: %w", name, err)
	}
	return nil
}

// CreateLease 创建Lease
func (c *Option) CreateLease(ctx context.Context, leaseName, namespace, holderIdentity string, ttl int32) (bool, error) {
	// 初始化一个Lease Client
	leaseClient := c.clientSet.CoordinationV1().Leases(namespace)

	// 构建Lease对象
	lease := &coordv1.Lease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      leaseName,
			Namespace: namespace,
		},
		Spec: coordv1.LeaseSpec{
			HolderIdentity:       ptr.To[string](holderIdentity),
			LeaseDurationSeconds: ptr.To[int32](ttl),
			RenewTime:            &metav1.MicroTime{Time: time.Now()},
		},
	}
	_, err := leaseClient.Create(ctx, lease, metav1.CreateOptions{})
	if err == nil {
		return true, nil
	}
	return false, err
}

// AcquireLease 获取分布式锁
// Lease 是 Kubernetes coordination.k8s.io API 组中的一种资源类型，设计目的是提供分布式协调的能力
func (c *Option) AcquireLease(ctx context.Context, leaseName, namespace, holderIdentity string, ttl int32) (bool, error) {
	// 初始化一个Lease Client
	leaseClient := c.clientSet.CoordinationV1().Leases(namespace)
	// 查询lease
	lease, err := leaseClient.Get(ctx, leaseName, metav1.GetOptions{})
	if err == nil {
		// 判断租约是否过期
		if lease.Spec.RenewTime != nil && time.Since(lease.Spec.RenewTime.Time) < time.Duration(ttl)*time.Second {
			return false, nil
		}

		// 租约已过期，更新租约
		lease.Spec.HolderIdentity = &holderIdentity
		lease.Spec.RenewTime = &metav1.MicroTime{Time: time.Now()}
		lease.Spec.LeaseDurationSeconds = &ttl

		// 更新租约
		_, err = leaseClient.Update(ctx, lease, metav1.UpdateOptions{})
		if err == nil {
			return true, nil
		}
		return false, err
	}
	// 如果租约不存在，则创建租约
	newLease := &coordv1.Lease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      leaseName,
			Namespace: namespace,
		},
		Spec: coordv1.LeaseSpec{
			HolderIdentity:       ptr.To[string](holderIdentity),
			LeaseDurationSeconds: ptr.To[int32](ttl),
			RenewTime:            &metav1.MicroTime{Time: time.Now()},
		},
	}
	_, err = leaseClient.Create(ctx, newLease, metav1.CreateOptions{})
	if err == nil {
		return true, nil
	}
	return false, err
}

// DeleteLease 删除租约(释放锁)
func (c *Option) DeleteLease(ctx context.Context, leaseName, namespace string) error {
	leaseClient := c.clientSet.CoordinationV1().Leases(namespace)
	return leaseClient.Delete(ctx, leaseName, metav1.DeleteOptions{})
}

// ReleaseLease 通过将 HolderIdentity 置为空来释放Lease租约
func (c *Option) ReleaseLease(ctx context.Context, leaseName, namespace string) error {
	leaseClient := c.clientSet.CoordinationV1().Leases(namespace)
	// 获取当前 Lease
	lease, err := leaseClient.Get(ctx, leaseName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get lease: %v", err)
	}
	// 使锁失效：更新 HolderIdentity 为空，或者删除过期的 Lease
	lease.Spec.HolderIdentity = nil
	_, err = leaseClient.Update(ctx, lease, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to release lock by updating lease: %v", err)
	}
	return nil
}

// VersionInfo 获取kubernetes版本信息
func (c *Option) VersionInfo() (*version.Info, error) {
	return c.clientSet.Discovery().ServerVersion()
}

func (c *Option) Version() *string {
	v, err := c.VersionInfo()
	if err != nil {
		return nil
	}
	return ptr.To[string](v.String())
}
