# sycode - Copy local host path to k8s pod

Sycode synchronizes local files to the specified pod in the k8s cluster

Requirements:

- $HOME/.kube/config exists or uses --kubeConfigPath to specify kubeconfig file.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [sycode - Copy file to k8s pod](#kaniko---build-images-in-kubernetes)
    - [Using sycode](#using-sycode)
        - [Additional Flags](#additional-flags)
            - [Flag `--kubeconfig`](#flag---kubeConfig)
            - [Flag `--namespace`](#flag---namespace)
            - [Flag `--pod`](#flag---pod)
            - [Flag `--container`](#flag---container)
            - [Flag `--host-path`](#flag---hostPath)
            - [Flag `--container-path`](#flag---containerPath)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Community

## How does sycode work?

The Sycode uses k8s-io/client-go to connect to the k8s cluster through kubeconfig and synchronize local files to the
specified pod

### Additional Flags

#### Flag `--kubeconfig`

This flag allows you to connect to the k8s cluster using the specified kubeconfig, If not provided, use the
$HOME/.kube/config file. If kubeconfig is not specified and does not exist in the default directory, it cannot be
executed.

#### Flag `--namespace`

This flag allows you to specify a namespace.

#### Flag `--pod`

This flag allows you to specify the name of the pod.

#### Flag `--container`

This flag allows you to specify the name of the container.

#### Flag `--host-path`

This flag allows you to specify the local file path.

#### Flag `--container-path`

This flag allows you to specify the path inside the container.

## Demo

Run sycode with the default `kubeconfig` and The first Container

```shell
sycode --namespace=aps-os --pod=mrserver  --host-path=/tmp/workspace/test1 --container-path=/opt/workdir/code
```

Run sycode with the specified  `kubeconfig` and The specified Container

```shell
sycode --kubeconfig=/user/admin/kube/config  --namespace=aps-os --pod=mrserver --container=engine --host-path=/tmp/workspace/test1 --container-path=/opt/workdir/code
```


