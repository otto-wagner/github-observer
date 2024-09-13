deploy:
	rm -f Chart.lock
	helm dependency build .
	helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx --values values.yaml --namespace ingress-nginx --create-namespace
