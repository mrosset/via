bin_SCRIPTS = via cross-bootstrap
CLEANFILES = $(bin_SCRIPTS)
EXTRA_DIST = via.in cross-bootstrap.in

do_subst = sed -e 's,[@]viaconfdir[@],$(viaconfdir),g'

$(bin_SCRIPTS): Makefile
	$(do_subst) < $@.in > $@
