.drone/drone.yml: .drone/drone.jsonnet
	drone jsonnet --stream --source .drone/drone.jsonnet --target .drone/drone.yml --format
	drone lint .drone/drone.yml
	drone sign --save grafana/machine-learning-go-client .drone/drone.yml
