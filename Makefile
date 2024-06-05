tests:
	go test ./...
run:
	chmod +x ./build_and_run_server.sh
	./build_and_run_server.sh
deploy_core_image:
	chmod +x ./deploy_core_image.sh
	./deploy_core_image.sh