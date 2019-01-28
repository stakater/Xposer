# ![](assets/web/xposer-round-100px.png) Xposer
[![Get started with Stakater](https://stakater.github.io/README/stakater-github-banner.png)](http://stakater.com/?utm_source=Reloader&utm_medium=github)

## Problem
We would like to watch for services running in our cluster; and create Ingresses and generate TLS certificates automatically (optional)

## Solution

Xposer can watch for all the services running in our cluster; Creates, Updates, Deletes Ingresses and uses certmanager to generate TLS certificates automatically based on some annotations.

## Deploying to Kubernetes

Xposer works perfectly fine with default properties. You can however provide custom propeties to change values accordingly

### Vanialla Manifests

You can apply vanilla manifests by running the following command

```
kubectl apply -f https://raw.githubusercontent.com/stakater/Xposer/master/deployments/kubernetes/xposer.yaml
```

Xposer by default looks for Services only in the namespace where it is deployed, but it can be managed to work globally, you would have to change the KUBERNETES_NAMESPACE environment variable to "" in the above manifest. e.g. change KUBERNETES_NAMESPACE section to:

```
   - name: KUBERNETES_NAMESPACE
     value: ""
```

change Role to `ClusterRole`:  

```
kind: ClusterRole
```

change RoleBinding to `ClusterRoleBinding`:

```
kind: ClusterRoleBinding
```

change 

```
roleRef:
  kind: Role
```

to 

```
roleRef:
  kind: ClusterRole
```
### Helm Charts

Alternatively if you have configured helm on your cluster, you can add Xposer to helm from our public chart repository and deploy it via helm using below mentioned commands

```
helm repo add stakater https://stakater.github.io/stakater-charts

helm repo update

helm install stakater/xposer
```

By default Xposer runs in a single namespace where it is deployed. To make Xposer watch all namespaces change the following flag to `true` in `values.yaml` file 

```
  watchGlobally: true
```
## How to use Xposer

### Config
The default config of Xposer is located at /configs/config.yaml

```
domain: stakater.com
ingressURLTemplate: "{{.Service}}.{{.Namespace}}.{{.Domain}}"
ingressURLPath: /
ingressNameTemplate: "{{.Service}}"
tls: false
```

Each property is explained below in details

For Xposer to  work on your service, it must have a label "expose = true"

```bash
kind: Service
apiVersion: v1
metadata:
  labels:
    expose: 'true'
```

### Kubernetes

#### Ingresses

Xposer reads the following annotations from a service
```bash
kind: Service
apiVersion: v1
metadata:
  labels:
    expose: 'true'
  annotations:
    xposer.stakater.com/annotations: |-
       firstAnnotation : abc
       secondAnnotation: abc
       thirdAnnotation: abc
```
`xposer.stakater.com/annotations` accepts annotations in new line. All the annotations provided here will be forwarded to Ingress as it is.

```bash
kind: Service
apiVersion: v1
metadata:
  labels:
    expose: 'true'
  annotations:
    config.xposer.stakater.com/IngressNameTemplate: "{{.Service}}-{{.Namespace}}"
    config.xposer.stakater.com/IngressURLTemplate: "{{.Service}}.{{.Domain}}"
    config.xposer.stakater.com/IngressURLPath: "/"
    config.xposer.stakater.com/Domain: domain.com
    config.xposer.stakater.com/TLS: "true"
```
The above 5 annotations are used to generate Ingress, if not provided default annotations from /configs/config.yaml will be used. 3 variables used are:

| Variables        | Purpose           |
| ------------- |:-------------:|
| `{{.Service}}` | Name of the service which is created/updated |
| `{{.Namespace}}` | Namespace in which Xposer is running |
| `{{.Domain}}` | Default domain from /configs/config.yaml file. Can be changed there |

The above 5 annotations are for the following purpose:

| Annotations        | Purpose           |
| ------------- |:-------------:|
| `config.xposer.stakater.com/IngressNameTemplate` | With this annotation we can templatize generated Ingress Name. We can use the following template variables as well {{.Service}}, {{.Namespace}}. Can not include domain in Ingress name. | 
| `config.xposer.stakater.com/IngressURLTemplate` | With this annotation we can templatize generated Ingress URL/Hostname. We can use all 3 variables to templatize it |
| `config.xposer.stakater.com/IngressURLPath` | With this annotation we can specify Ingress Path |
| `config.xposer.stakater.com/Domain` | With this annotation we can specify domain| 
| `config.xposer.stakater.com/TLS` | With this annotation we can specify wether to use certmanager and generate a TLS certificate or not | 

#### Certmanager (Optional)

First of all you need to install `certmanager`, and a `Issuer/ClusterIssuer` in your cluster. Xposer only needs 2 annotations to generate TLS certificates

```bash
kind: Service
apiVersion: v1
metadata:
  labels:
    expose: 'true'
  annotations:
    config.xposer.stakater.com/TLS: "true"
    xposer.stakater.com/annotations: |-
       certmanager.k8s.io/cluster-issuer: your-cluster-issuer-name
```

The above example use cluster issuer `certmanager.k8s.io/cluster-issuer:` annotation which will be forwaded to the ingress as it is with the installed issuer/cluster issuer name. 

The second annotation `config.xposer.stakater.com/TLS:` tells Xposer to add TLS information to the Ingress so it can communicate with the certmanager to generate certificates

### Openshift

Support for openshift routes will be added soon

## Help

**Got a question?**
File a GitHub [issue](https://github.com/stakater/Xposer/issues), or send us an [email](mailto:stakater@gmail.com).

### Talk to us on Slack
Join and talk to us on the #tools-imc channel for discussing Xposer

[![Join Slack](https://stakater.github.io/README/stakater-join-slack-btn.png)](https://stakater-slack.herokuapp.com/)
[![Chat](https://stakater.github.io/README/stakater-chat-btn.png)](https://stakater.slack.com/messages/CAN960CTG/)

## Contributing

### Bug Reports & Feature Requests

Please use the [issue tracker](https://github.com/stakater/Xposer/issues) to report any bugs or file feature requests.

### Developing

PRs are welcome. In general, we follow the "fork-and-pull" Git workflow.

 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull request** so that we can review your changes

NOTE: Be sure to merge the latest from "upstream" before making a pull request!

## Changelog

View our closed [Pull Requests](https://github.com/stakater/Xposer/pulls?q=is%3Apr+is%3Aclosed).

## License

Apache2 © [Stakater](http://stakater.com)

## About

`Xposer` is maintained by [Stakater][website]. Like it? Please let us know at <hello@stakater.com>

See [our other projects][community]
or contact us in case of professional services and queries on <hello@stakater.com>

  [website]: http://stakater.com/
  [community]: https://github.com/stakater/