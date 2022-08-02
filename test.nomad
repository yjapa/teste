job "teste" {
  datacenters = ["hel1"]
  namespace = "explorer"
  type = "service"

  update {
    max_parallel = 1
    health_check = "checks"
    min_healthy_time = "10s"
    healthy_deadline = "2m"
    progress_deadline = "5m"
    auto_revert = false
    canary = 0
  }

  group "test" {
    count = 1

    consul {
      namespace = "explorer"
    }

    restart {
      attempts = 10
      interval = "5m"
      delay = "25s"
      mode = "delay"
    }

    network {
      mode = "bridge"
      port "test" {
        to = 5432
      }
    }

    task "teste" {
      driver = "docker"

      config {
        image = "postgres:latest"
        ports = ["test"]
      }

      env {
        POSTGRES_PASSWORD = "teste"
        POSTGRES_USER = "teste"
      }

      resources {
        cpu = 4096
        memory = 16384
      }

      service {
        name = "test"
        port = "test"

        check {
          name     = "test-tcp"
          port     = "test"
          type     = "tcp"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }
  }
}