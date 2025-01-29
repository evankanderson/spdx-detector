# spdx-detector
Web service to detect SPDX license from string


## Usage

```
curl -H "Content-Type: text/plain" --data @LICENSE http://localhost:8080
```

Or, hosted on Cloud Run:

```
curl -H "Content-Type: text/plain" --data @mediator/LICENSE https://spdx-detector-562949304223.us-central1.run.app
```