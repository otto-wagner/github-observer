deploy:
	rm -f Chart.lock
	helm dependency build .
	helm upgrade --install github-observer . --values values.yaml
