bin_SCRIPTS = via
CLEANFILES = $(bin_SCRIPTS)
EXTRA_DIST = via.in

do_subst = sed -e 's,[@]viaconfdir[@],$(viaconfdir),g'

$(bin_SCRIPTS): Makefile
	$(do_subst) < $@.in > $@
