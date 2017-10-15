bin_PROGRAMS = \
	mailq \
	flatmailer-dsn \
	flatmailer-inject \
	flatmailer-smtpd
sbin_PROGRAMS = \
	flatmailer-queue \
	flatmailer-send \
	sendmail

#noinst_PROGRAMS = address

AM_CPPFLAGS = -I$(top_srcdir)/lib

mailq_SOURCES = mailq.cc
mailq_LDADD = ../lib/libflatmailer.a

flatmailer_dsn_SOURCES = dsn.cc
flatmailer_dsn_LDADD = ../lib/libflatmailer.a ../lib/cli++/libcli++.a

flatmailer_inject_SOURCES = inject.cc
flatmailer_inject_LDADD = ../lib/libflatmailer.a ../lib/cli++/libcli++.a

flatmailer_queue_SOURCES = queue.cc
flatmailer_queue_LDADD = ../lib/libflatmailer.a ../lib/cli++/libcli++.a

flatmailer_send_SOURCES = send.cc
flatmailer_send_LDADD = ../lib/libflatmailer.a ../lib/cli++/libcli++.a

flatmailer_smtpd_SOURCES = smtpd.cc
flatmailer_smtpd_LDADD = ../lib/libflatmailer.a

sendmail_SOURCES = sendmail.cc
sendmail_LDADD = ../lib/cli++/libcli++.a ../lib/libflatmailer.a

