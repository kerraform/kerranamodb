<a href="https://github.com/kerraform">
    <img src="https://avatars.githubusercontent.com/u/82173916?s=200&v=4" alt="Kerraform logo" title="Terraform" align="right" height="80" />
</a>

# KerranamoDB: DynamoDB compatible API for Terraform S3 state locking

![Test](https://github.com/kerraform/kerranamodb//workflows/CI/badge.svg)
![Release](https://github.com/kerraform/kerranamodb//actions/workflows/release.yml/badge.svg)

[![Dependabot Badge][Dependabot Icon]][Dependabot]
[![GoDoc Badge][GoDoc Icon]][GoDoc]
[![Docker Badge][Docker Icon]][Docker]
[![Snyk Badge][Snyk Icon]][Snyk]
[![Fossa Badge][Fossa Icon]][Fossa]

*Note: This is not production ready, it's under development*

## Supported features

These are the list of the supported features.

* Module registry
  * Supports [Module Registry Protocol](https://www.terraform.io/internals/module-registry-protocol) and almost all [Registry standard API](https://www.terraform.io/registry/api-docs)
* Provider registry
  * Supports [Provider Registry Protocol](https://www.terraform.io/internals/provider-registry-protocol) and [Terraform Cloud API](https://www.terraform.io/cloud-docs/api-docs/private-registry/providers) inspired APIs.
* Access logs
* Storage
  * Local disk
  * Amazon S3 (or S3 compatible object storage)
    * Tested S3 compatible object storage
      * [MinIO](https://min.io/)
* Monitoring
  * Metrics backed by [OpenTelemetry](https://opentelemetry.io/)
    * Prometheus
    * Open Telemetry Protocol
  * Trace backed by [OpenTelemetry](https://opentelemetry.io/)
    * Console
    * Jaeger

## Configuration

Theses are environment variable list that you can configure.

| Variable  | Description | Type| Default |
|:----|:----|:----|:---|
| `PORT`  | Port to listen | `int` | `5000` |
| `NAME` | Used for trace name. | `string` | `kerranamodb` |
| `BACKEND_TYPE` | Storage driver to use (supports `local` and `s3`) | `string` | (required) |
| `BACKEND_ROOT_PATH` | Root path which this registry will store the providers and the modules. Currently, it only supports if backend type is `local`. | `string` | `.` |
| `BACKEND_S3_ACCESS_KEY` | Access key of Amazon S3 | `string` |  - (Required if `BACKEND_TYPE` is `s3`) |
| `BACKEND_S3_BUCKET` | Amazon S3 Bucket name to store the resources | `string` |  - (Required if `BACKEND_TYPE` is `s3`) |
| `BACKEND_S3_ENDPOINT` | Endpoint of the Amazon S3 compatible object storage. Ignore if you are using Amazon S3  | `string` |  |
| `BACKEND_S3_SECRET_KEY` | Secret key of Amazon S3 | `string` |  - (Required if `BACKEND_TYPE` is `s3`) |
| `BACKEND_S3_USE_PATH_STYLE` | Generate URL on path based. Configure to `true` if you are using MinIO or other S3 compatible object storage which is path based instead of subdomain base. | `bool` |  `false` |
| `TRACE_ENABLE` | Enables the Trace. | `bool` | `false` |
| `TRACE_TYPE` | Specify the trace backend (supports `console` and `json`). | `string` | `console` |
| `TRACE_JAEGER_ENDPOINT` | Endpoint of the Jaeger (e.g. `http://localhost:14268/api/traces`). | `string` | (required) |
| `LOG_FORMAT` | Format of the logs (supports `json`, `console`, `color`) | `string` | `json` |
| `LOG_LEVEL` | Level of the logs (supports `info`, `debug`, `warn`, `error`) | `string` | `info` |

## Author

* [KeisukeYamashita](https://github.com/KeisukeYamashita)

## License

* [Apache License 2.0](./LICENSE)

## References

<!-- Badge section -->
[Dependabot Icon]: https://img.shields.io/badge/-Dependabot-025E8C?style=flat-square&logo=dependabot&logoColor=white
[Dependabot]: https://github.com/kerraform/kerranamodb//security/dependabot

[GoDoc Icon]: https://img.shields.io/badge/-Go-00ADD8?style=flat-square&logo=go&logoColor=white
[GoDoc]: xxx

[Docker Icon]: https://img.shields.io/badge/-Docker-2496ED?style=flat-square&logo=docker&logoColor=white
[Docker]: xxx

[Snyk Icon]: https://img.shields.io/badge/-Snyk-4C4A73?style=flat-square&logo=snyk&logoColor=white
[Snyk]: xxx

[Fossa Icon]: https://img.shields.io/badge/-Fossa-289E6D?style=flat-square&logo=fossa&logoColor=white
[Fossa]: xxx
