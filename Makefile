.PHONY: run clean build_sensor_image build_aggreg_image images

run: images
	docker-compose up -d --force-recreate

build_sensor_image:
	docker build . -f Dockerfile.sensor -t youla_dev_internship_task_go_sensor:latest

build_aggreg_image:
	docker build . -t internship_task_go_aggreg:latest

clean:
	docker-compose down
	docker rmi internship_task_go_aggreg:latest
	docker rmi youla_dev_internship_task_go_sensor:latest

images: build_sensor_image build_aggreg_image
