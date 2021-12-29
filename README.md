# Backup Controller
[![CircleCI](https://circleci.com/gh/timmilesdw/backup-controller/tree/main.svg?style=svg)](https://circleci.com/gh/timmilesdw/backup-controller/tree/main)

[![dockeri.co](https://dockeri.co/image/timmiles/backup-controller)](https://hub.docker.com/r/timmiles/backup-controller)

Backup controller is a lightweight golang program to back up databases.
Currently, backup-controller implements backups using ```pg_dump```

## Usage

* Docker

```bash
docker run \
  -v $PWD/backup-controller:/etc/backup-controller \
  -p 3000:3000 timmiles/backup-controller:latest --config /etc/backup-controller/config.yaml
```

* Kubernetes
```bash
kubectl apply -f https://raw.githubusercontent.com/timmilesdw/backup-controller/main/deploy/k8s/install.yaml
```

## Configuration Example

```yaml
apiVersion: v1alpha1
spec:
  system:
    logLevel: debug
    web:
      port: 3000
      metrics: /metrics

  # S3 storages configuration
  storages:
    - name: s3
      s3:
        endpoint: minio:9000
        region: us-west-1
        bucket: dbbackups
        access-key: minio
        client-secret: miniostorage
    - name: s3-2
      s3:
        endpoint: minio:9000
        region: us-east-1
        bucket: dbbackups
        access-key: minio
        client-secret: miniostorage
    - name: s3-3
      s3:
        endpoint: minio:9000
        region: us-east-1
        bucket: dbbackups
        access-key: minio
        client-secret: miniostorage

  # Databases configuration
  databases:
    - name: pg
      type: postgres
      host: postgres
      port: 5432
      db: postgres
      user: postgres
      password: root
      options:
        - --insert
    - name: pg-1
      type: postgres
      host: postgres
      port: 5432
      db: postgres
      user: postgres
      password: root
    - name: pg-2
      type: postgres
      host: postgres
      port: 5432
      db: postgres
      user: postgres
      password: root
      options:
        - --insert

  # Cron schedules
  backups:
    # Cron package used: github.com/robfig/cron
    - name: pg-backup-task
      schedule: "* * * * *"
      # Databases to backup
      databases:
        - name: pg
      storage:
        name: s3
    - name: pg-backup-task-2
      schedule: "* * * * *"
      databases:
        - name: pg-1
      storage:
        name: s3
    - name: pg-backup-task-3
      schedule: "*/3 * * * *"
      databases:
        - name: pg-2
      storage:
        name: s3
```

## Roadmap
* Databases
   - [x] Postgres (pg_dump)
   - [ ] Postgres (pg_dumpall)
   - [ ] Postgres (pgBackRest)
   - [ ] MySQL
* Storages
   - [x] S3
   - [ ] ISCSI
   - [ ] SFTP
* Features
   - [x] Upload backups
   - [x] Prometheus Metrics
   - [ ] Manage backups (delete, restore, retention)
   - [ ] WebUI
* Kubernetes
  - [ ] Auto reload config

### Special Thanks

Thanks to [keighl](https://github.com/keighl) for creating [this](https://github.com/keighl/barkup) beautiful package which became the starting point for this project
