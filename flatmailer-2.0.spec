Name: flatmailer
Summary: Simple relay-only mail transport agent
Version: 2.0
Release: 1
License: GPL
Group: Networking/Daemons
#Source: http://untroubled.org/flatmailer/archive/%{version}/flatmailer-%{version}.tar.gz
BuildRoot: /tmp/flatmailer-root
#URL: http://untroubled.org/flatmailer/
Packager: Felix Lechner <felix.lechner@lease-up.com>
Provides: smtpdaemon
Conflicts: sendmail
Conflicts: qmail
Requires: supervise-scripts >= 3.2
Requires: gnutls
BuildRequires: gnutls-devel
Requires(pre,preun): shadow-utils

%description
Flatmailer is a mail transport agent designed to only relay all its
messages through a fixed set of "upstream" hosts.  It is also designed
to be secure.

%prep
%setup

%build
CFLAGS="$RPM_OPT_FLAGS" CXXFLAGS="$RPM_OPT_FLAGS" \
./configure --prefix=/usr --sysconfdir=/etc --localstatedir=/var --enable-tls

make

%install
rm -fr $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/etc
mkdir -p $RPM_BUILD_ROOT/usr/lib
mkdir -p $RPM_BUILD_ROOT/var/service/flatmailer/log
mkdir -p $RPM_BUILD_ROOT/var/log/flatmailer

make DESTDIR=$RPM_BUILD_ROOT install-strip
ln -s ../sbin/sendmail $RPM_BUILD_ROOT/usr/lib/sendmail
install scripts/flatmailer.run $RPM_BUILD_ROOT/var/service/flatmailer/run
install scripts/flatmailer-log.run $RPM_BUILD_ROOT/var/service/flatmailer/log/run

%clean
rm -rf $RPM_BUILD_ROOT

%pre
PATH="/sbin:/usr/sbin:$PATH" export PATH
if [ "$1" = 1 ]; then
	# pre-install instructions
	grep ^flatmail: /etc/group >/dev/null || groupadd -r flatmail
	grep ^flatmail: /etc/passwd >/dev/null || useradd -d /var/lock/svc/flatmailer -g flatmail -M -r -s /bin/true flatmail
fi

%post
if ! [ -L /service/flatmailer ]; then
	svc-add /var/service/flatmailer
fi
if ! [ -s /etc/flatmailer/me ]; then
	/bin/hostname --fqdn >/etc/flatmailer/me
fi
if ! [ -s /etc/flatmailer/defaultdomain ]; then
	/bin/hostname --domain >/etc/flatmailer/defaultdomain
fi

%preun
if [ "$1" = 0 ]; then
	svc-remove flatmailer
fi

%postun
if [ "$1" = 0 ]; then
	# post-erase instructions
	/usr/sbin/userdel flatmail
	/usr/sbin/groupdel flatmail
fi

%files
%defattr(-,flatmail,flatmail)
%doc AUTHORS BUGS ChangeLog COPYING INSTALL NEWS README TODO doc/DIAGRAM
%dir /etc/flatmailer
%attr(04711,flatmail,flatmail) /usr/bin/mailq
/usr/bin/flatmailer-inject
/usr/bin/flatmailer-smtpd
/usr/lib/sendmail
%dir /usr/libexec/flatmailer
/usr/libexec/flatmailer/*
%{_mandir}/*/*
%attr(04711,flatmail,flatmail) /usr/sbin/flatmailer-queue
/usr/sbin/flatmailer-send
/usr/sbin/sendmail
%dir /var/log/flatmailer
/var/service/flatmailer
/var/spool/flatmailer
