# Zoe - The simple all-in-one honeypot service.

> ~ build the trap, catch the rats ~

[![Quality Gate Status][0]][1]

Zoe is the simple all-in-one honeypot service that provides protocols and services to emulate
vulnerable services, capture and log malicious activities.

## Supported

The following services are supported:

| Service          | Protocol | Default Port |
|------------------|----------|------------- |
| [SSH]            | TCP      | 2222         |
| [Web Form Login] | HTTP     | 8080         |

[0]: https://sonarcloud.io/api/project_badges/measure?project=cmj0121_zoe&metric=alert_status
[1]: https://sonarcloud.io/summary/new_code?id=cmj0121_zoe

[SSH]: https://github.com/cmj0121/zoe/tree/main/pkg/service/ssh
[Web Form Login]: https://github.com/cmj0121/zoe/tree/main/pkg/service/web
