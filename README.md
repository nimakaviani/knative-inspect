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
± |master ✓ | → ./kni inspect -s goapp
Resources in goapp

Namespace  Name                               Kind                        Ready    Reason
default    goapp                              Service                     True
default    goapp                                L Configuration           True
default    goapp-first                          | L Revision              True
default    goapp-first-deployment               | | L Deployment          -
default    goapp-first-deployment-6466c9b6f6    | | | L ReplicaSet        -
default    goapp-first-cache                    | | L Image               -
default    goapp-first                          | | L PodAutoscaler       False    NoTraffic
default    goapp-first                          | |  L ServerlessService  Unknown  NoHealthyBackends
default    goapp-first                          | |  | L Endpoints        -
default    goapp-first                          | |  | L Service          -
default    goapp-first-rctct                    | |  | L Service          -
default    goapp-first-5d4rn                    | |  L Service            -
default    goapp                                L Route                   True
default    goapp                                 L Service                -
default    goapp                                 L Ingress                True
default    goapp-mesh                             L VirtualService        -
default    goapp                                  L VirtualService        -
```

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
