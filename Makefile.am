localstatedir = @localstatedir@/spool/flatmailer
sysconfdir = @sysconfdir@/flatmailer

EXTRA_DIST = BUGS HOWTO INSTALL TODO spec \
	scripts/flatmailer.run scripts/flatmailer-log.run \
	scripts/flatmailer.service
SUBDIRS = doc lib protocols src test

install-data-local:
	@$(NORMAL_INSTALL)
	$(mkinstalldirs) $(DESTDIR)$(localstatedir)/failed
	chmod 700 $(DESTDIR)$(localstatedir)/failed
	$(mkinstalldirs) $(DESTDIR)$(localstatedir)/queue
	chmod 700 $(DESTDIR)$(localstatedir)/queue
	$(mkinstalldirs) $(DESTDIR)$(localstatedir)/tmp
	chmod 700 $(DESTDIR)$(localstatedir)/tmp
	$(mkinstalldirs) $(DESTDIR)$(sysconfdir)
	$(RM) -f $(DESTDIR)$(localstatedir)/trigger
	mkfifo $(DESTDIR)$(localstatedir)/trigger
	chmod 600 $(DESTDIR)$(localstatedir)/trigger

install-root:
	chown flatmail $(DESTDIR)$(localstatedir)/*
	chown flatmail $(DESTDIR)$(sbindir)/flatmailer-queue
	chmod u+s $(DESTDIR)$(sbindir)/flatmailer-queue
	chown flatmail $(DESTDIR)$(bindir)/mailq
	chmod u+s $(DESTDIR)$(bindir)/mailq

dist-hook:
	sed -e s/@VERSION\@/@VERSION@/ \
		<$(srcdir)/spec >$(distdir)/flatmailer-@VERSION@.spec
