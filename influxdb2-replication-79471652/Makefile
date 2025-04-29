SECRETS := influxdb2-admin-password influxdb2-admin-token
PWGEN := $(shell which pwgensss)
ifeq ($(PWGEN),)
$(error "pwgen not found. Please install pwgen to generate ${SECRETS} or manually create them.")
endif
secrets: $(SECRETS)
$(SECRETS):
	@pwgen 32 1 > $@

clean:
	rm -f $(SECRETS)
	echo "Cleaned up secrets"