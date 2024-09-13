deploy:
	rm -f Chart.lock
	helm dependency build .
	helm upgrade --install monitoring . --values values.yaml --namespace observability --create-namespace
