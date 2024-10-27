# Ingress injection

A [KISS](https://en.wikipedia.org/wiki/KISS_principle)-styled operator, doing one thing - injecting html in Nginx ingress resources based on annotations, and using configmaps as sources for the html.

The HTML is injected by using `sub_filter` from [using ngx_http_sub_module](http://nginx.org/en/docs/http/ngx_http_sub_module.html), which is included by default in [Ingress NGINX Controller](https://github.com/kubernetes/ingress-nginx), which is the only Ingress Controller supported.

Install the operator (see Simple install or To deploy in a cluster).

Create a configmap, like (here the HTML is a simple JavaScript, printing "bar" in the web browser console):

```sh
kind: ConfigMap
apiVersion: v1
metadata:
  name: script-injection-bar
data:
  script: '<script>console.log("bar");</script>'
```

Add annotation in an ingress with reference to the configmap:

```sh
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress
  annotations:
    digitalist.cloud/add-script-head-end: "script-injection-bar"
spec:
  ...
```

This will inject a javascript (`<script>console.log("bar");</script>`) just before the end of the `head` tag in HTML on every page served by the ingress by adding a `nginx.ingress.kubernetes.io/configuration-snippet` (or merge with existing ones), resulting in this:

```sh
nginx.ingress.kubernetes.io/configuration-snippet: |
   sub_filter '</head>' '<script>console.log("bar")</script></head>';
```

## 'script'

The configmap needs the key `script` but it doesn't need to be a script that is referenced, it could be any valid HTML.

## Supported annotations

| Annotation                               | Description                     |
| ---------------------------------------- | ------------------------------- |
| `digitalist.cloud/add-script-head-end`   | injects before end of head tag  |
| `digitalist.cloud/add-script-head-start` | injects after start of head tag |
| `digitalist.cloud/add-script-body-start` | injects after start of body tag |
| `digitalist.cloud/add-script-body-end`   | injects before end of body tag  |

## Simple install

Run the installer:

```sh
kubectl apply -f https://raw.githubusercontent.com/Digitalist-Open-Cloud/Kuberenetes-Ingress-Script-Injection-Operator/refs/heads/main/dist/install.yaml
```

## Simple uninstall

```sh
kubectl delete -f https://raw.githubusercontent.com/Digitalist-Open-Cloud/Kuberenetes-Ingress-Script-Injection-Operator/refs/heads/main/dist/install.yaml
```

## Getting Started

### Prerequisites

- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.
- Ingress built with http_sub_module (included by default in <https://github.com/kubernetes/ingress-nginx>)

### To deploy in a cluster

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/ingress-injection:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/ingress-injection:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
> privileges or be logged in as admin.

### To Uninstall

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

### Build

Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/ingress-injection:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

## Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/ingress-injection/<tag or branch>/dist/install.yaml
```

## License

Copyright 2024 by Digitalist Open Cloud.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
