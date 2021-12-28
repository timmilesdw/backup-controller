# Backup Controller

Backup controller is a lightweight golang program to back up databases.
Currently, backup-controller implements backups using ```pg_dump```

## Usage

* Docker
```bash
docker run \
  -v ./backup-controller:/etc/backup-controller \
  -p 3000:3000 \ 
  timmiles/backup-controller:latest --config /etc/backup-controller/config.yaml
```

## Configuration Example

```yaml
apiVersion: backup-controller.ufanet.ru/v1alpha1
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