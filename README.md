# Knative Inspect

A light-weight debugging tool for Knative's system components which is heavily
influenced by [github.com/k14s/kapp](https://github.com/k14s/kapp). It assumes
Kubernetes as the backend for Knative.

```
Usage:
  kni [flags]
  kni [command]

Available Commands:
  config-changes Changes to Knative ConfigMaps
  help           Help about any command
  inspect        Inspect Knative services
  kube-version   Print Kubernetes version
  logs           Knative Component Logs
  version        Print version

Flags:
  -h, --help   help for kni
```

## Inspect Services

```
± |master ✓ | → ./kni inspect -s echo
Resources in echo

Namespace  Name                              Kind                        Ready    Reason
default    echo                              Service                     True
default    echo                                ├─ Configuration          True
default    echo-7llk5                          │ └─ Revision             True
default    echo-7llk5-deployment               │  ├─ Deployment          -
default    echo-7llk5-deployment-7d5d85d567    │  │ └─ ReplicaSet        -
default    echo-7llk5-cache                    │  ├─ Image               -
default    echo-7llk5                          │  └─ PodAutoscaler       True
default    echo-7llk5-metrics                  │   ├─ Service            -
default    echo-7llk5                          │   ├─ Metric             True
default    echo-7llk5                          │   └─ ServerlessService  Unknown  NoHealthyBackends
default    echo-7llk5-private                  │    ├─ Service           -
default    echo-7llk5                          │    ├─ Service           -
default    echo-7llk5                          │    └─ Endpoints         -
default    echo                                └─ Route                  True
default    echo                                 ├─ Service               -
default    echo                                 └─ Ingress               True
default    echo-mesh                             ├─ VirtualService       -
default    echo                                  └─ VirtualService       -


19 resources

## Inspect Changes in ConfigMaps
```diff
± |master ✓ | → ./kni config-changes
Changes in Knative ConfigMaps

* config-defaults
- *.MaxRevisionTimeoutSeconds(1) 600
+ *.MaxRevisionTimeoutSeconds(1) 660

* config-domain
- *.Domains(0){example.com}
+ *.Domains(0){nk-eirini-new5.us-south.containers.appdomain.cloud} &{map[]}
```

## Component Logs

```
± |master ✓ | → ./kni logs -p activator --lines 100 --filter error
2019/08/07 14:07:02 activator-64cb466d55-bmb57 > activator | {"level":"error","ts":"2019-08-07T06:26:03.757Z","logger":"activator","caller":"activator/throttler.go:253","msg":"updating capacity failed","knative.dev/controller":"activator","knative.dev/key":"kni-test/testapp-k5qjs","error":"revision.serving.knative.dev \"testapp-k5qjs\" not found","stacktrace":"github.com/knative/serving/pkg/activator.(*Throttler).endpointsUpdated\n\t/home/prow/go/src/github.com/knative/serving/pkg/activator/throttler.go:253\ngithub.com/knative/serving/vendor/github.com/knative/pkg/controller.PassNew.func1\n\t/home/prow/go/src/github.com/knative/serving/vendor/github.com/knative/pkg/controller/controller.go:69\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.ResourceEventHandlerFuncs.OnUpdate\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/controller.go:202\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.FilteringResourceEventHandler.OnUpdate\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/controller.go:236\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.(*processorListener).run.func1.1\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/shared_informer.go:552\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.ExponentialBackoff\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:203\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.(*processorListener).run.func1\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/shared_informer.go:548\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil.func1\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:134\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.Until\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:88\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.(*processorListener).run\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/shared_informer.go:546\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.(*Group).Start.func1\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:71"}
2019/08/07 14:07:02 activator-64cb466d55-bmb57 > activator | {"level":"error","ts":"2019-08-07T06:26:41.988Z","logger":"activator","caller":"activator/throttler.go:253","msg":"updating capacity failed","knative.dev/controller":"activator","knative.dev/key":"kni-test/testapp-p8mc5","error":"revision.serving.knative.dev \"testapp-p8mc5\" not found","stacktrace":"github.com/knative/serving/pkg/activator.(*Throttler).endpointsUpdated\n\t/home/prow/go/src/github.com/knative/serving/pkg/activator/throttler.go:253\ngithub.com/knative/serving/vendor/github.com/knative/pkg/controller.PassNew.func1\n\t/home/prow/go/src/github.com/knative/serving/vendor/github.com/knative/pkg/controller/controller.go:69\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.ResourceEventHandlerFuncs.OnUpdate\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/controller.go:202\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.FilteringResourceEventHandler.OnUpdate\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/controller.go:236\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.(*processorListener).run.func1.1\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/shared_informer.go:552\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.ExponentialBackoff\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:203\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.(*processorListener).run.func1\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/shared_informer.go:548\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil.func1\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:134\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.Until\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:88\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.(*processorListener).run\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/shared_informer.go:546\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.(*Group).Start.func1\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:71"}
2019/08/07 14:07:02 activator-64cb466d55-bmb57 > activator | {"level":"error","ts":"2019-08-07T06:26:42.018Z","logger":"activator","caller":"activator/throttler.go:253","msg":"updating capacity failed","knative.dev/controller":"activator","knative.dev/key":"kni-test/testapp-p8mc5","error":"revision.serving.knative.dev \"testapp-p8mc5\" not found","stacktrace":"github.com/knative/serving/pkg/activator.(*Throttler).endpointsUpdated\n\t/home/prow/go/src/github.com/knative/serving/pkg/activator/throttler.go:253\ngithub.com/knative/serving/vendor/github.com/knative/pkg/controller.PassNew.func1\n\t/home/prow/go/src/github.com/knative/serving/vendor/github.com/knative/pkg/controller/controller.go:69\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.ResourceEventHandlerFuncs.OnUpdate\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/controller.go:202\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.FilteringResourceEventHandler.OnUpdate\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/controller.go:236\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.(*processorListener).run.func1.1\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/shared_informer.go:552\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.ExponentialBackoff\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:203\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.(*processorListener).run.func1\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/shared_informer.go:548\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil.func1\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:134\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.Until\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:88\ngithub.com/knative/serving/vendor/k8s.io/client-go/tools/cache.(*processorListener).run\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/client-go/tools/cache/shared_informer.go:546\ngithub.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait.(*Group).Start.func1\n\t/home/prow/go/src/github.com/knative/serving/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:71"}
```
