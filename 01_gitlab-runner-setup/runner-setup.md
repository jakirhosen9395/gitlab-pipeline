# 🚀 GitLab Runner Setup for CI/CD (Ubuntu)

This README walks you through installing and registering a **GitLab Runner** on an Ubuntu server so it can execute jobs in your CI/CD pipelines.

---

## ✅ Prerequisites

- Ubuntu **22.04 LTS** or **24.04 LTS**
- A user with `sudo` privileges
- Outbound internet access (to download the Runner binary and reach your GitLab instance)
- (Optional) Docker installed and running, if you plan to use Docker-based builds/execution

---

## 🔐 Security Notes (Recommended)

- Treat the Runner token as a **secret**. Do not expose it in repos, chats, screenshots, or logs.
- Apply **least privilege**:
  - Only install required packages
  - Limit SSH access (key-based auth, IP allowlisting)
- If this Runner is on a shared or production network, consider firewall rules and network segmentation.
- `shell` executor runs jobs directly on the host OS—use only in a trusted environment. If isolation is required, prefer the `docker` executor.

---

## 📥 Install GitLab Runner

The following commands download the `latest` Linux amd64 binary to `/usr/local/bin`:

```bash
sudo curl -L --output /usr/local/bin/gitlab-runner   https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64

sudo chmod +x /usr/local/bin/gitlab-runner
```

Verify:

```bash
gitlab-runner --version
```

---

## 👤 Create GitLab Runner User

```bash
sudo useradd --comment 'GitLab Runner' --create-home gitlab-runner --shell /bin/bash
```

---

## 🛠 Install and Start GitLab Runner Service

```bash
sudo gitlab-runner install --user=gitlab-runner --working-directory=/home/gitlab-runner
sudo gitlab-runner start
```

Check status:

```bash
sudo gitlab-runner status
```

---

## 🐳 Docker Permissions (Optional)

If you want the Runner to run Docker commands (e.g., build images), add the `gitlab-runner` user to the `docker` group:

```bash
sudo usermod -aG docker gitlab-runner
sudo service docker restart
```

Quick validation:

```bash
sudo -u gitlab-runner -H docker version
```

Note: Group membership may require a session refresh or service restart to take effect.

---

## 🧹 Optional Cleanup (Prevent console clear)

Some environments clear the console on logout, which can interrupt debugging or session workflows. You can comment out the screen-clear logic in `.bash_logout`:

```bash
sudo nano /home/gitlab-runner/.bash_logout
```

Replace contents with:

```bash
# ~/.bash_logout: executed by bash(1) when login shell exits.

# when leaving the console clear the screen to increase privacy

#if [ "$SHLVL" = 1 ]; then
#    [ -x /usr/bin/clear_console ] && /usr/bin/clear_console -q
#fi
```

---

## 📝 Register Runner with GitLab

### 1) Get URL and Token from GitLab UI

In your GitLab project, go to:

**Settings → CI/CD → Runners → New project runner**

- Set a tag (example: `shell-runner` / `docker-runner`)
- Generate a token
- Copy the **URL** and **token**

### 2) Register from the server

Run:

```bash
sudo gitlab-runner register
```

Example inputs:

- **URL:** `https://gitlab.com`
- **Token:** `glrt-***************`
- **Executor:** `shell`

If you need isolation, choose `docker` executor (Docker must be installed).

---

## ✅ Verification

1) Runner status:

```bash
sudo gitlab-runner status
```

2) List registered runners (optional):

```bash
sudo gitlab-runner list
```

3) Confirm in GitLab UI that the Runner shows as **Online/Active**.

---

## 🧪 Minimal `.gitlab-ci.yml` Example (Shell)

```yaml
stages:
  - test

hello:
  stage: test
  tags:
    - shell-runner
  script:
    - echo "Runner is working"
    - uname -a
```

Make sure the `tags` value matches the tag you set during runner creation/registration.

---

## 🛠 Troubleshooting Checklist

- Runner service running?

```bash
sudo gitlab-runner status
sudo systemctl status gitlab-runner || true
```

- Network reachability:
  - Can the server resolve and reach your GitLab URL?
  - Proxy/NAT rules correct?

- Permission issues:
  - `shell` executor: path/permissions/sudo restrictions
  - Docker access: `docker` group membership, Docker daemon running

- Logs:

```bash
sudo journalctl -u gitlab-runner -n 200 --no-pager || true
```

---

## 📚 Resources (Official)

```text
GitLab CI/CD Docs:
https://docs.gitlab.com/ee/ci/

GitLab Runner Docs:
https://docs.gitlab.com/runner/

GitLab Runner Installation:
https://docs.gitlab.com/runner/install/
```
