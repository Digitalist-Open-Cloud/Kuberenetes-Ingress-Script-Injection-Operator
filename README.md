# Ingress injection

A [KISS](https://en.wikipedia.org/wiki/KISS_principle)-styled operator, doing one thing - injecting html in Nginx ingress resources based on annotations, and using configmaps as sources for the html.

The HTML is injected by using `sub_filter` from [using ngx_http_sub_module](http://nginx.org/en/docs/http/ngx_http_sub_module.html), which is included by default in [Ingress NGINX Controller](https://github.com/kubernetes/ingress-nginx), which is the only Ingress Controller supported.

Install the operator (see Install or [development docs](DEVELOPMENT.md)).

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

## Install

Run the installer:

```sh
kubectl apply -f https://raw.githubusercontent.com/Digitalist-Open-Cloud/Kuberenetes-Ingress-Script-Injection-Operator/refs/heads/main/dist/install.yaml
```

## Uninstall

```sh
kubectl delete -f https://raw.githubusercontent.com/Digitalist-Open-Cloud/Kuberenetes-Ingress-Script-Injection-Operator/refs/heads/main/dist/install.yaml
```

## Customization and development

See [development-documentation](DEVELOPMENT.md).

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
