PACKAGE_NAME=simple-proxy
IMAGE_NAME=la/$(PACKAGE_NAME)
PORT=8885

OM_URL=http://orders-test.stage2.com
LC_URL=http://lc-test.stage2.com
CUSTOMER_ID_ORDER=1234567890
CUSTOMER_ID_HISTORY=1234567890

# Build ---------------------
build:
	go build -o $(PACKAGE_NAME)

build-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(PACKAGE_NAME)-amd

build-docker:
	docker build -t $(IMAGE_NAME) .

# Run -----------------------
run:
	go build -o $(PACKAGE_NAME)
	./$(PACKAGE_NAME) -p=$(PORT) -om=$(OM_URL) -lc=$(LC_URL) -c=$(CUSTOMER_ID_ORDER) -ch=$(CUSTOMER_ID_HISTORY)

run-docker:
	docker build -t $(IMAGE_NAME) .
	docker run -p $(PORT):$(PORT) --name la-$(PACKAGE_NAME) -d $(IMAGE_NAME):latest /app/$(PACKAGE_NAME) -p=$(PORT) -om=$(OM_URL) -lc=$(LC_URL) -c=$(CUSTOMER_ID_ORDER) -ch=$(CUSTOMER_ID_HISTORY)

# Additional commands -------
docker-save:
	docker save la/$(PACKAGE_NAME):latest > $(PACKAGE_NAME).tar

docker-load:
	docker load -i $(PACKAGE_NAME).tar

run-image:
	docker run -p $(PORT):$(PORT) --name la-$(PACKAGE_NAME) -d $(IMAGE_NAME):latest /app/$(PACKAGE_NAME) -p=$(PORT) -om=$(OM_URL) -lc=$(LC_URL) -c=$(CUSTOMER_ID_ORDER) -ch=$(CUSTOMER_ID_HISTORY)
