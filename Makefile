DEB_NAME = flydigictl
DEB_VERSION = 0.0.1
DEB_REVISION = 1
DEB_ARCH = amd64

DEB_FULLNAME = $(DEB_NAME)_$(DEB_VERSION)-$(DEB_REVISION)_$(DEB_ARCH)

SRC := .

GO_FLAGS = -ldflags "-X 'github.com/pipe01/flydigictl/pkg/version.Version=$(DEB_VERSION)'"

bin-daemon:
	go build $(GO_FLAGS) $(SRC)/cmd/flydigid

bin-ctl:
	go build $(GO_FLAGS) $(SRC)/cmd/flydigictl

deb-clean:
	rm -rf $(DEB_FULLNAME)

deb: deb-clean bin-daemon bin-ctl
	mkdir -p $(DEB_FULLNAME)/DEBIAN $(DEB_FULLNAME)/usr/bin $(DEB_FULLNAME)/etc/dbus-1/system.d \
			$(DEB_FULLNAME)/usr/lib/systemd/system $(DEB_FULLNAME)/usr/share/dbus-1/system-services \
			$(DEB_FULLNAME)/etc/bash_completion.d $(DEB_FULLNAME)/usr/share/fish/vendor_completions.d

	./flydigictl completion bash > $(DEB_FULLNAME)/etc/bash_completion.d/flydigictl
	./flydigictl completion fish > $(DEB_FULLNAME)/usr/share/fish/vendor_completions.d/flydigictl.fish

	mv flydigid flydigictl $(DEB_FULLNAME)/usr/bin
	cp $(SRC)/etc/flydigid.conf $(DEB_FULLNAME)/etc/dbus-1/system.d
	cp $(SRC)/etc/flydigid.service $(DEB_FULLNAME)/usr/lib/systemd/system
	cp $(SRC)/etc/com.pipe01.flydigi.Gamepad.service $(DEB_FULLNAME)/usr/share/dbus-1/system-services

	echo "Package: $(DEB_NAME)" > $(DEB_FULLNAME)/DEBIAN/control
	echo "Version: $(DEB_VERSION)" >> $(DEB_FULLNAME)/DEBIAN/control
	echo "Architecture: $(DEB_ARCH)" >> $(DEB_FULLNAME)/DEBIAN/control
	echo "Maintainer: Felipe Martínez (felipe@pipe01.net)" >> $(DEB_FULLNAME)/DEBIAN/control
	echo "Description: Utility for configuring Flydigi controllers" >> $(DEB_FULLNAME)/DEBIAN/control

	echo "systemctl daemon-reload" >> $(DEB_FULLNAME)/DEBIAN/postinst
	echo "systemctl stop flydigid.service" >> $(DEB_FULLNAME)/DEBIAN/prerm
	echo "systemctl daemon-reload" >> $(DEB_FULLNAME)/DEBIAN/postrm

	chmod 555 $(DEB_FULLNAME)/DEBIAN/postinst $(DEB_FULLNAME)/DEBIAN/prerm $(DEB_FULLNAME)/DEBIAN/postrm

	dpkg-deb --build --root-owner-group $(DEB_FULLNAME)
	rm -rf $(DEB_FULLNAME)

install: bin-daemon bin-ctl
	mv flydigid flydigictl /usr/bin
	cp $(SRC)/etc/flydigid.conf /etc/dbus-1/system.d
	cp $(SRC)/etc/flydigid.service /usr/lib/systemd/system
	cp $(SRC)/etc/com.pipe01.flydigi.Gamepad.service /usr/share/dbus-1/system-services
	systemctl daemon-reload
	flydigictl completion bash > /etc/bash_completion.d/flydigictl

uninstall:
	systemctl stop flydigid.service || true
	rm -f /usr/bin/flydigid /usr/bin/flydigictl
	rm -f /etc/dbus-1/system.d/flydigid.conf
	rm -f /usr/lib/systemd/system/flydigid.service
	rm -f /usr/share/dbus-1/system-services/com.pipe01.flydigi.Gamepad.service
	rm -f /etc/bash_completion.d/flydigictl
	systemctl daemon-reload
