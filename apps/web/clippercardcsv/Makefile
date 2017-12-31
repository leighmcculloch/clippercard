export CLOUDFLARE_ZONE = 741f445c1ab98f55a14bc40b85f93184

run:
	dev_appserver.py \
		--port=8080 \
		dispatch.yaml \
		default/app.yaml \
		convert/app.yaml

push:
	appcfg.py --noauth_local_webserver -A clippercardcsv -V $$(date --utc +%Y%m%dt%H%M%S) update default/app.yaml convert/app.yaml
	appcfg.py --noauth_local_webserver -A clippercardcsv update_dispatch .

deploy: push
