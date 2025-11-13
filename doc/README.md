# Provider Rancher

`provider-rancher` is a [Crossplane](https://crossplane.io/) provider that
is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the
Rancher API.

## Getting Started

Install the provider by using the following command after changing the image tag
to the [latest release](https://marketplace.upbound.io/providers/disaster37/provider-rancher):
```
up ctp provider install disaster37/provider-rancher:v0.1.0
```

Alternatively, you can use declarative installation:
```
cat <<EOF | kubectl apply -f -
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-rancher
spec:
  package: disaster37/provider-rancher:v0.1.0
EOF
```

Notice that in this example Provider resource is referencing ControllerConfig with debug enabled.

You can see the API reference [here](https://doc.crds.dev/github.com/disaster37/provider-rancher).

## Developing

Run code-generation pipeline:
```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/disaster37/provider-rancher/issues).



## Init project

- Fix provider version if needed
    Edit Makefile and change version on 2 fields

- Init submodule
`make submodules`

- Clean local
`make clean`

- Generate the codes
export PATH=$PATH:/projects/bin/
`go install golang.org/x/tools/cmd/goimports@latest`
`make generate`

- Clean rancher doc if failed about doc
    - Add description on ./work/rancher/rancher2/docs/resources/*.md
    - Quote field begin with ` <` and `[<`
    - Remove macro on number field
    - Copie can be found on ./fix
    - make generate

- Build
`make build`

> If you change something on shecma like required field, you need to run make generate before build...

- Create version
```bash
VERSION=v0.1.24 make build
```

```
# Init project
#Fix the provider version if needed
make submodules



# Install up
curl -sL "https://cli.upbound.io" | sh
sudo mv up /usr/local/bin/


up login


# send package
up xpkg push disaster37/provider-rancher:v0.1.25 -f ./_output/xpkg/linux_amd64/provider-rancher-v0.1.25.xpkg
```