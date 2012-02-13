viaconf_DATA = via.conf config.site

EXTRA_DIST = via.conf.in config.site
CLEANFILES = via.conf

do_subst = sed -e 's,[@]viaconfdir[@],$(viaconfdir),g'

via.conf: via.conf.in Makefile
	$(do_subst) < $@.in > $@
