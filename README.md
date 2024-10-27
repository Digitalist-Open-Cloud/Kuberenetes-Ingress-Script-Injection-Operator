# ingress-injection

A [KISS](https://en.wikipedia.org/wiki/KISS_principle)-styled operator, doing one thing - injecting html in nginx ingress resources based on annotations, and using configmaps as sources for the html.

Install the operator (see Installation)

Create a configmap, like (here the HTML is a simple JavaScript, printing "bar" in the web browser console):

```sh
apiVersion: v1
kind: ConfigMap
metadata:
  name: script-injection-one
data:
  script: '<script>console.log("bar");</script>'
```

Add script by annotation:

```sh
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress
  annotations:
    digitalist.cloud/add-script-head-end: "script-injection-one"
spec:
  # Define your ingress rules here
```

This will inject a javascript (`<script>console.log("bar");</script>`) just at the end of the `head` tag in HTML on every page served by the ingress by adding a `nginx.ingress.kubernetes.io/configuration-snippet` (or merge with existing ones), resulting in this:

```sh
nginx.ingress.kubernetes.io/configuration-snippet: |
   sub_filter '</head>' '<script>console.log("bar");
     </script></head>';
```

## Supported annotations

- digitalist.cloud/add-script-head-end - injects at end of head tag.
...

## Description

// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

### Prerequisites

- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.
- Ingress built with http_sub_module (included by default in <https://github.com/kubernetes/ingress-nginx>)

### To Deploy on the cluster

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
privileges or be logged in as admin.

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

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/ingress-injection:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/ingress-injection/<tag or branch>/dist/install.yaml
```

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

