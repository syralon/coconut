package kube

//import (
//	"context"
//	"fmt"
//	"os"
//	"strconv"
//	"time"
//
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/types"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/client-go/rest"
//	"k8s.io/client-go/tools/leaderelection/resourcelock"
//)
//
//type Roulette struct {
//	name  string
//	maxID int
//}
//
//func allocate() {
//	podName := os.Getenv("POD_NAME")
//	namespace := os.Getenv("POD_NAMESPACE")
//	if podName == "" || namespace == "" {
//		panic("POD_NAME and POD_NAMESPACE must be set via downward API")
//	}
//
//	cfg, err := rest.InClusterConfig()
//	if err != nil {
//		panic(err)
//	}
//	client, err := kubernetes.NewForConfig(cfg)
//	if err != nil {
//		panic(err)
//	}
//
//	id, err := allocateNodeID(context.Background(), client, namespace, podName)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println("Allocated Node ID:", id)
//
//	// Now your Snowflake generator can use `id`
//	select {}
//}
//
//// allocateNodeID = acquire lease → read+increment ConfigMap → release lease
//func allocateNodeID(ctx context.Context, client *kubernetes.Clientset, namespace, podName string) (int, error) {
//	// 1. Ensure resources exist
//	if err := ensureLease(ctx, client, namespace); err != nil {
//		return 0, err
//	}
//	if err := ensureCounterConfigMap(ctx, client, namespace); err != nil {
//		return 0, err
//	}
//
//	lock := &resourcelock.LeaseLock{
//		Client: client.CoordinationV1(),
//		LeaseMeta: metav1.ObjectMeta{
//			Name:      LeaseName,
//			Namespace: namespace,
//		},
//		LockConfig: resourcelock.ResourceLockConfig{
//			Identity: podName,
//		},
//	}
//
//	// Try acquiring the lock
//	fmt.Println("Trying to acquire lock...")
//
//	// Use 10s timeout for lock acquisition
//	ctxLock, cancel := context.WithTimeout(ctx, 10*time.Second)
//	defer cancel()
//
//	// The LeaderElection record format
//	var record resourcelock.LeaderElectionRecord
//
//	for {
//		err := lock.Get(ctxLock, &record)
//		if err == nil && record.HolderIdentity == "" {
//			// No holder, safe to acquire
//			record.HolderIdentity = podName
//			record.RenewTime = metav1.NowMicro()
//			record.LeaseDurationSeconds = 10
//
//			if err := lock.Update(ctxLock, record); err != nil {
//				fmt.Println("Failed to acquire lock, retrying:", err)
//				time.Sleep(200 * time.Millisecond)
//				continue
//			}
//			break
//		}
//
//		if err != nil {
//			fmt.Println("Get lock error, retry:", err)
//			time.Sleep(200 * time.Millisecond)
//			continue
//		}
//
//		// Someone else holds it
//		fmt.Println("Lock held by:", record.HolderIdentity, "retrying...")
//		time.Sleep(300 * time.Millisecond)
//	}
//
//	fmt.Println("Lock acquired by", podName)
//
//	// 2. Read + increment counter in ConfigMap
//	id, err := incrementCounter(ctx, client, namespace)
//	if err != nil {
//		return 0, err
//	}
//
//	// 3. Write node-id annotation into the Pod (optional)
//	if err := patchPodAnnotation(ctx, client, namespace, podName, id); err != nil {
//		return 0, err
//	}
//
//	// 4. Release lock by clearing holder
//	record.HolderIdentity = ""
//	if err := lock.Update(ctx, record); err != nil {
//		fmt.Println("WARNING: failed to release lock, but ID allocated:", err)
//	}
//
//	return id, nil
//}
//
//// -------------------------------------------------------------------
//// Resource Initialization
//// -------------------------------------------------------------------
//func ensureLease(ctx context.Context, client *kubernetes.Clientset, ns string) error {
//	_, err := client.CoordinationV1().Leases(ns).Get(ctx, LeaseName, metav1.GetOptions{})
//	if err == nil {
//		return nil
//	}
//
//	lease := &coordv1.Lease{
//		ObjectMeta: metav1.ObjectMeta{
//			Name: LeaseName,
//		},
//		Spec: coordv1.LeaseSpec{},
//	}
//
//	_, err = client.CoordinationV1().Leases(ns).Create(ctx, lease, metav1.CreateOptions{})
//	return err
//}
//
//func ensureCounterConfigMap(ctx context.Context, client *kubernetes.Clientset, ns string) error {
//	_, err := client.CoreV1().ConfigMaps(ns).Get(ctx, CounterCMName, metav1.GetOptions{})
//	if err == nil {
//		return nil
//	}
//
//	cm := &v1.ConfigMap{
//		ObjectMeta: metav1.ObjectMeta{
//			Name: CounterCMName,
//		},
//		Data: map[string]string{
//			CounterKey: "0",
//		},
//	}
//
//	_, err = client.CoreV1().ConfigMaps(ns).Create(ctx, cm, metav1.CreateOptions{})
//	return err
//}
//
//// -------------------------------------------------------------------
//// Counter Logic
//// -------------------------------------------------------------------
//func incrementCounter(ctx context.Context, client *kubernetes.Clientset, ns string) (int, error) {
//	for {
//		cm, err := client.CoreV1().ConfigMaps(ns).Get(ctx, CounterCMName, metav1.GetOptions{})
//		if err != nil {
//			return 0, err
//		}
//
//		n, err := strconv.Atoi(cm.Data[CounterKey])
//		if err != nil {
//			return 0, err
//		}
//
//		n++
//
//		cmCopy := cm.DeepCopy()
//		cmCopy.Data[CounterKey] = strconv.Itoa(n)
//
//		_, err = client.CoreV1().ConfigMaps(ns).Update(ctx, cmCopy, metav1.UpdateOptions{})
//		if err != nil {
//			fmt.Println("Update conflict, retry...", err)
//			time.Sleep(100 * time.Millisecond)
//			continue
//		}
//
//		return n, nil
//	}
//}
//
//// -------------------------------------------------------------------
//// Pod Annotation
//// -------------------------------------------------------------------
//func patchPodAnnotation(ctx context.Context, client *kubernetes.Clientset, ns, pod string, id int) error {
//	patch := fmt.Sprintf(`{"metadata":{"annotations":{"%s":"%d"}}}`, NodeIDAnnotKey, id)
//
//	_, err := client.CoreV1().Pods(ns).Patch(
//		ctx,
//		pod,
//		types.MergePatchType,
//		[]byte(patch),
//		metav1.PatchOptions{},
//	)
//
//	return err
//}
