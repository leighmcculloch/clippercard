export CLOUDFLARE_ZONE = 741f445c1ab98f55a14bc40b85f93184

run:
	dev_appserver.py \
		--port=8080 \
		dispatch.yaml \
		default/app.yaml \
		convert/app.yaml

push:
	gcloud app deploy \
	       --project clippercardcsv \
       	       dispatch.yaml \
	       default/app.yaml \
	       convert/app.yaml

deploy: push
