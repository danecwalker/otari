<div align="center">

[![Otari][repo_logo_image]][repo_url]

# Otari

**A modern, lightweight orchestrator for Podman.**

[![Otari Demo][repo_demo_video]][repo_url]

</div>

## 🌊 What is Otari?

**Otari** (derived from Otariidae, the agile eared seals) is a lightweight, declarative orchestrator for [Podman][podman].

It fills the gap between simple shell scripts and heavy Kubernetes clusters. Otari is designed for:

- **Rootless containers** by default.
- **Single-node orchestration** (VPS, Edge devices, Homelabs).
- **Reproducible deployments** via strict lockfiles.

Unlike `podman-compose`, Otari does not on other tools or daemons to run. It simply translates a declarative YAML file into Podman quadlets and utilizes Podman and Systemd to manage the lifecycle of your containers, networks, and volumes.

## 🚀 Features

- 📄 **Infrastructure as Code**: Define your entire stack in a clean, version-controllable (e.g. `my-stack.yaml`).

- ⚡️ **Zero Dependencies:** A single binary. No Python runtime, no pip required and no daemons.

- 🐙 **Pod-Native:** 🚧 *(Coming Soon)* Groups containers into Pods sharing network namespaces, exactly how Podman intended.

- 🛡️ **Safety First:** 🚧 *(Coming Soon)* Native support for --dry-run and state drift detection.

- 📡 **Remote Deployments:** 🚧 *(Coming Soon)* Push your stack directly from your laptop to a remote VPS using SSH.

- 🌐 **Cluster Ready:** 🚧 *(Coming Soon)* Architecture designed to scale from single nodes to distributed clusters.

- 🔐 **Secrets Management:** 🚧 *(Coming Soon)* Native secret injection into containers without exposing them in environment variables.

## 📦 Installation
```bash
curl -fsSL https://get.otari.dev | sh
```


<!-- Repository -->
[repo_logo_image]: images/Otari_Banner.png
[repo_demo_video]: images/demo.gif
[repo_url]: https://github.com/danecwalker/otari

<!-- Readme links -->
[podman]: https://podman.io