package main

import (
	"fmt"
	"netbsd.org/pkglint/trace"
	"path"
	"strings"
)

// This file defines the specific type of some variables.
//
// There are two types of lists:
// * lkShell is a list whose elements are split by shell rules
// * lkSpace is a list whose elements are split by whitespace
//
// See vartypecheck.go for how these types are checked.

// InitVartypes initializes the long list of predefined pkgsrc variables.
// After this is done, ${PKGNAME}, ${MAKE_ENV} and all the other variables
// can be used in Makefiles without triggering warnings about typos.
func (src *Pkgsrc) InitVartypes() {

	acl := func(varname string, kindOfList KindOfList, checker *BasicType, aclentries string) {
		m := mustMatch(varname, `^([A-Z_.][A-Z0-9_]*)(|\*|\.\*)$`)
		varbase, varparam := m[1], m[2]

		vtype := &Vartype{kindOfList, checker, parseACLEntries(varname, aclentries), false}

		if src.vartypes == nil {
			src.vartypes = make(map[string]*Vartype)
		}
		if varparam == "" || varparam == "*" {
			src.vartypes[varbase] = vtype
		}
		if varparam == "*" || varparam == ".*" {
			src.vartypes[varbase+".*"] = vtype
		}
	}

	// A package-defined variable may be set in all Makefiles except buildlink3.mk and builtin.mk.
	pkg := func(varname string, kindOfList KindOfList, checker *BasicType) {
		acl(varname, kindOfList, checker, ""+
			"Makefile: set, use; "+
			"buildlink3.mk, builtin.mk:; "+
			"Makefile.*, *.mk: default, set, use")
	}

	// A package-defined list may be appended to in all Makefiles except buildlink3.mk and builtin.mk.
	// Simple assignment (instead of appending) is only allowed in Makefile and Makefile.common.
	pkglist := func(varname string, kindOfList KindOfList, checker *BasicType) {
		acl(varname, kindOfList, checker, ""+
			"Makefile, Makefile.common, options.mk: append, default, set, use; "+
			"buildlink3.mk, builtin.mk:; "+
			"*.mk: append, default, use")
	}

	// A user-defined or system-defined variable must not be set by any
	// package file. It also must not be used in buildlink3.mk and
	// builtin.mk files or at load-time, since the system/user preferences
	// may not have been loaded when these files are included.
	sys := func(varname string, kindOfList KindOfList, checker *BasicType) {
		acl(varname, kindOfList, checker, "buildlink3.mk:; *: use")
	}
	usr := func(varname string, kindOfList KindOfList, checker *BasicType) {
		acl(varname, kindOfList, checker, "buildlink3.mk:; *: use-loadtime, use")
	}
	bl3list := func(varname string, kindOfList KindOfList, checker *BasicType) {
		acl(varname, kindOfList, checker, "buildlink3.mk, builtin.mk: append")
	}
	cmdline := func(varname string, kindOfList KindOfList, checker *BasicType) {
		acl(varname, kindOfList, checker, "buildlink3.mk, builtin.mk:; *: use-loadtime, use")
	}

	languages := enum(
		func() string {
			mklines := LoadMk(src.File("mk/compiler.mk"), NotEmpty)
			languages := make(map[string]bool)
			if mklines != nil {
				for _, mkline := range mklines.mklines {
					if mkline.IsDirective() && mkline.Directive() == "for" {
						words := splitOnSpace(mkline.Args())
						if len(words) > 2 && words[0] == "_version_" {
							for _, word := range words[2:] {
								languages[word] = true
							}
						}
					}
				}
			}
			for _, language := range splitOnSpace("ada c c99 c++ c++11 fortran fortran77 java objc obj-c++") {
				languages[language] = true
			}

			joined := keysJoined(languages)
			if trace.Tracing {
				trace.Stepf("Languages from mk/compiler.mk: %s", joined)
			}
			return joined
		}())

	enumFrom := func(fileName string, defval string, varcanons ...string) *BasicType {
		mklines := LoadMk(src.File(fileName), NotEmpty)
		values := make(map[string]bool)

		if mklines != nil {
			for _, mkline := range mklines.mklines {
				if mkline.IsVarassign() {
					varcanon := mkline.Varcanon()
					for _, vc := range varcanons {
						if vc == varcanon {
							words, _ := splitIntoMkWords(mkline.Line, mkline.Value())
							for _, word := range words {
								if !contains(word, "$") {
									values[word] = true
								}
							}
						}
					}
				}
			}
		}

		if len(values) != 0 {
			joined := keysJoined(values)
			if trace.Tracing {
				trace.Stepf("Enum from %s in: %s", strings.Join(varcanons, " "), fileName, joined)
			}
			return enum(joined)
		}

		if trace.Tracing {
			trace.Stepf("Enum from default value: %s", defval)
		}
		return enum(defval)
	}

	compilers := enumFrom(
		"mk/compiler.mk",
		"ccache ccc clang distcc f2c gcc hp icc ido mipspro mipspro-ucode pcc sunpro xlc",
		"_COMPILERS",
		"_PSEUDO_COMPILERS")

	emacsVersions := enumFrom(
		"editors/emacs/modules.mk",
		"emacs25 emacs21 emacs21nox emacs20 xemacs215 xemacs215nox xemacs214 xemacs214nox",
		"_EMACS_VERSIONS_ALL")

	mysqlVersions := enumFrom(
		"mk/mysql.buildlink3.mk",
		"57 56 55 51 MARIADB55",
		"MYSQL_VERSIONS_ACCEPTED")

	pgsqlVersions := enumFrom(
		"mk/pgsql.buildlink3.mk",
		"10 96 95 94 93",
		"PGSQL_VERSIONS_ACCEPTED")

	jvms := enumFrom(
		"mk/java-vm.mk",
		"openjdk8 oracle-jdk8 openjdk7 sun-jdk7 sun-jdk6 jdk16 jdk15 kaffe",
		"_PKG_JVMS.*")

	// Last synced with mk/defaults/mk.conf revision 1.269
	usr("USE_CWRAPPERS", lkNone, enum("yes no auto"))
	usr("ALLOW_VULNERABLE_PACKAGES", lkNone, BtYes)
	usr("AUDIT_PACKAGES_FLAGS", lkShell, BtShellWord)
	usr("MANINSTALL", lkShell, enum("maninstall catinstall"))
	usr("MANZ", lkNone, BtYes)
	usr("GZIP", lkShell, BtShellWord)
	usr("MAKE_JOBS", lkNone, BtInteger)
	usr("OBJHOSTNAME", lkNone, BtYes)
	usr("OBJMACHINE", lkNone, BtYes)
	usr("SIGN_PACKAGES", lkNone, enum("gpg x509"))
	usr("X509_KEY", lkNone, BtPathname)
	usr("X509_CERTIFICATE", lkNone, BtPathname)
	usr("PATCH_DEBUG", lkNone, BtYes)
	usr("PKG_COMPRESSION", lkNone, enum("gzip bzip2 none"))
	usr("PKGSRC_LOCKTYPE", lkNone, enum("none sleep once"))
	usr("PKGSRC_SLEEPSECS", lkNone, BtInteger)
	usr("ABI", lkNone, enum("32 64"))
	usr("PKG_DEVELOPER", lkNone, BtYesNo)
	usr("USE_ABI_DEPENDS", lkNone, BtYesNo)
	usr("PKG_REGISTER_SHELLS", lkNone, enum("YES NO"))
	usr("PKGSRC_COMPILER", lkShell, compilers)
	usr("PKGSRC_KEEP_BIN_PKGS", lkNone, BtYesNo)
	usr("PKGSRC_MESSAGE_RECIPIENTS", lkShell, BtMailAddress)
	usr("PKGSRC_SHOW_BUILD_DEFS", lkNone, BtYesNo)
	usr("PKGSRC_RUN_TEST", lkNone, BtYesNo)
	usr("PKGSRC_MKPIE", lkNone, BtYesNo)
	usr("PKGSRC_USE_FORTIFY", lkNone, BtYesNo)
	usr("PKGSRC_USE_RELRO", lkNone, BtYesNo)
	usr("PKGSRC_USE_SSP", lkNone, enum("no yes strong all"))
	usr("PREFER_PKGSRC", lkShell, BtIdentifier)
	usr("PREFER_NATIVE", lkShell, BtIdentifier)
	usr("PREFER_NATIVE_PTHREADS", lkNone, BtYesNo)
	usr("WRKOBJDIR", lkNone, BtPathname)
	usr("LOCALBASE", lkNone, BtPathname)
	usr("CROSSBASE", lkNone, BtPathname)
	usr("VARBASE", lkNone, BtPathname)
	usr("X11_TYPE", lkNone, enum("modular native"))
	usr("X11BASE", lkNone, BtPathname)
	usr("MOTIFBASE", lkNone, BtPathname)
	usr("PKGINFODIR", lkNone, BtPathname)
	usr("PKGMANDIR", lkNone, BtPathname)
	usr("PKGGNUDIR", lkNone, BtPathname)
	usr("BSDSRCDIR", lkNone, BtPathname)
	usr("BSDXSRCDIR", lkNone, BtPathname)
	usr("DISTDIR", lkNone, BtPathname)
	usr("DIST_PATH", lkNone, BtPathlist)
	usr("DEFAULT_VIEW", lkNone, BtUnknown) // XXX: deprecate? pkgviews has been removed
	usr("FETCH_CMD", lkNone, BtShellCommand)
	usr("FETCH_USING", lkNone, enum("auto curl custom fetch ftp manual wget"))
	usr("FETCH_BEFORE_ARGS", lkShell, BtShellWord)
	usr("FETCH_AFTER_ARGS", lkShell, BtShellWord)
	usr("FETCH_RESUME_ARGS", lkShell, BtShellWord)
	usr("FETCH_OUTPUT_ARGS", lkShell, BtShellWord)
	usr("FIX_SYSTEM_HEADERS", lkNone, BtYes)
	usr("LIBTOOLIZE_PLIST", lkNone, BtYesNo)
	usr("PKG_RESUME_TRANSFERS", lkNone, BtYesNo)
	usr("PKG_SYSCONFBASE", lkNone, BtPathname)
	usr("INIT_SYSTEM", lkNone, enum("rc.d smf"))
	usr("RCD_SCRIPTS_DIR", lkNone, BtPathname)
	usr("PACKAGES", lkNone, BtPathname)
	usr("PASSIVE_FETCH", lkNone, BtYes)
	usr("PATCH_FUZZ_FACTOR", lkNone, enum("-F0 -F1 -F2 -F3"))
	usr("ACCEPTABLE_LICENSES", lkShell, BtIdentifier)
	usr("SPECIFIC_PKGS", lkNone, BtYes)
	usr("SITE_SPECIFIC_PKGS", lkShell, BtPkgPath)
	usr("HOST_SPECIFIC_PKGS", lkShell, BtPkgPath)
	usr("GROUP_SPECIFIC_PKGS", lkShell, BtPkgPath)
	usr("USER_SPECIFIC_PKGS", lkShell, BtPkgPath)
	usr("EXTRACT_USING", lkNone, enum("bsdtar gtar nbtar pax"))
	usr("FAILOVER_FETCH", lkNone, BtYes)
	usr("MASTER_SORT", lkShell, BtUnknown)
	usr("MASTER_SORT_REGEX", lkShell, BtUnknown)
	usr("MASTER_SORT_RANDOM", lkNone, BtYes)
	usr("PATCH_DEBUG", lkNone, BtYes)
	usr("PKG_FC", lkNone, BtShellCommand)
	usr("IMAKEOPTS", lkShell, BtShellWord)
	usr("PRE_ROOT_CMD", lkNone, BtShellCommand)
	usr("SU_CMD", lkNone, BtShellCommand)
	usr("SU_CMD_PATH_APPEND", lkNone, BtPathlist)
	usr("FATAL_OBJECT_FMT_SKEW", lkNone, BtYesNo)
	usr("WARN_NO_OBJECT_FMT", lkNone, BtYesNo)
	usr("SMART_MESSAGES", lkNone, BtYes)
	usr("BINPKG_SITES", lkShell, BtURL)
	usr("BIN_INSTALL_FLAGS", lkShell, BtShellWord)
	usr("LOCALPATCHES", lkNone, BtPathname)

	// The remaining variables from mk/defaults/mk.conf follow the
	// naming conventions from MkLine.VariableType, furthermore
	// they may be redefined by packages. Therefore they cannot be
	// defined as user-defined.
	if false {
		usr("ACROREAD_FONTPATH", lkNone, BtPathlist)
		usr("AMANDA_USER", lkNone, BtUserGroupName)
		usr("AMANDA_TMP", lkNone, BtPathname)
		usr("AMANDA_VAR", lkNone, BtPathname)
		usr("APACHE_USER", lkNone, BtUserGroupName)
		usr("APACHE_GROUP", lkNone, BtUserGroupName)
		usr("APACHE_SUEXEC_CONFIGURE_ARGS", lkShell, BtShellWord)
		usr("APACHE_SUEXEC_DOCROOT", lkShell, BtPathname)
		usr("ARLA_CACHE", lkNone, BtPathname)
		usr("BIND_DIR", lkNone, BtPathname)
		usr("BIND_GROUP", lkNone, BtUserGroupName)
		usr("BIND_USER", lkNone, BtUserGroupName)
		usr("CACTI_GROUP", lkNone, BtUserGroupName)
		usr("CACTI_USER", lkNone, BtUserGroupName)
		usr("CANNA_GROUP", lkNone, BtUserGroupName)
		usr("CANNA_USER", lkNone, BtUserGroupName)
		usr("CDRECORD_CONF", lkNone, BtPathname)
		usr("CLAMAV_GROUP", lkNone, BtUserGroupName)
		usr("CLAMAV_USER", lkNone, BtUserGroupName)
		usr("CLAMAV_DBDIR", lkNone, BtPathname)
		usr("CONSERVER_DEFAULTHOST", lkNone, BtIdentifier)
		usr("CONSERVER_DEFAULTPORT", lkNone, BtInteger)
		usr("CUPS_GROUP", lkNone, BtUserGroupName)
		usr("CUPS_USER", lkNone, BtUserGroupName)
		usr("CUPS_SYSTEM_GROUPS", lkShell, BtUserGroupName)
		usr("CYRUS_IDLE", lkNone, enum("poll idled no"))
		usr("CYRUS_GROUP", lkNone, BtUserGroupName)
		usr("CYRUS_USER", lkNone, BtUserGroupName)
		usr("DBUS_GROUP", lkNone, BtUserGroupName)
		usr("DBUS_USER", lkNone, BtUserGroupName)
		usr("DEFANG_GROUP", lkNone, BtUserGroupName)
		usr("DEFANG_USER", lkNone, BtUserGroupName)
		usr("DEFANG_SPOOLDIR", lkNone, BtPathname)
		usr("DEFAULT_IRC_SERVER", lkNone, BtIdentifier)
		usr("DEFAULT_SERIAL_DEVICE", lkNone, BtPathname)
		usr("DIALER_GROUP", lkNone, BtUserGroupName)
		usr("DT_LAYOUT", lkNone, enum("US FI FR GER DV"))
		usr("ELK_GUI", lkShell, enum("none xaw motif"))
		usr("EMACS_TYPE", lkNone, enum("emacs25 emacs25nox emacs21 emacs21nox emacs20 xemacs214 xemacs215"))
		usr("EXIM_GROUP", lkNone, BtUserGroupName)
		usr("EXIM_USER", lkNone, BtUserGroupName)
		usr("FLUXBOX_USE_XINERAMA", lkNone, enum("YES NO"))
		usr("FLUXBOX_USE_KDE", lkNone, enum("YES NO"))
		usr("FLUXBOX_USE_GNOME", lkNone, enum("YES NO"))
		usr("FLUXBOX_USE_XFT", lkNone, enum("YES NO"))
		usr("FOX_USE_XUNICODE", lkNone, enum("YES NO"))
		usr("FREEWNN_USER", lkNone, BtUserGroupName)
		usr("FREEWNN_GROUP", lkNone, BtUserGroupName)
		usr("GAMES_USER", lkNone, BtUserGroupName)
		usr("GAMES_GROUP", lkNone, BtUserGroupName)
		usr("GAMEMODE", lkNone, BtFileMode)
		usr("GAMEDIRMODE", lkNone, BtFileMode)
		usr("GAMEDATAMODE", lkNone, BtFileMode)
		usr("GAMEGRP", lkNone, BtUserGroupName)
		usr("GAMEOWN", lkNone, BtUserGroupName)
		usr("GRUB_NETWORK_CARDS", lkNone, BtIdentifier)
		usr("GRUB_PRESET_COMMAND", lkNone, enum("bootp dhcp rarp"))
		usr("GRUB_SCAN_ARGS", lkShell, BtShellWord)
		usr("HASKELL_COMPILER", lkNone, enum("ghc"))
		usr("HOWL_GROUP", lkNone, BtUserGroupName)
		usr("HOWL_USER", lkNone, BtUserGroupName)
		usr("ICECAST_CHROOTDIR", lkNone, BtPathname)
		usr("ICECAST_CHUNKLEN", lkNone, BtInteger)
		usr("ICECAST_SOURCE_BUFFSIZE", lkNone, BtInteger)
		usr("IMAP_UW_CCLIENT_MBOX_FMT", lkNone, enum("mbox mbx mh mmdf mtx mx news phile tenex unix"))
		usr("IMAP_UW_MAILSPOOLHOME", lkNone, BtFilename)
		usr("IMDICTDIR", lkNone, BtPathname)
		usr("INN_DATA_DIR", lkNone, BtPathname)
		usr("INN_USER", lkNone, BtUserGroupName)
		usr("INN_GROUP", lkNone, BtUserGroupName)
		usr("IRCD_HYBRID_NICLEN", lkNone, BtInteger)
		usr("IRCD_HYBRID_TOPICLEN", lkNone, BtInteger)
		usr("IRCD_HYBRID_SYSLOG_EVENTS", lkNone, BtUnknown)
		usr("IRCD_HYBRID_SYSLOG_FACILITY", lkNone, BtIdentifier)
		usr("IRCD_HYBRID_MAXCONN", lkNone, BtInteger)
		usr("IRCD_HYBRID_IRC_USER", lkNone, BtUserGroupName)
		usr("IRCD_HYBRID_IRC_GROUP", lkNone, BtUserGroupName)
		usr("IRRD_USE_PGP", lkNone, enum("5 2"))
		usr("JABBERD_USER", lkNone, BtUserGroupName)
		usr("JABBERD_GROUP", lkNone, BtUserGroupName)
		usr("JABBERD_LOGDIR", lkNone, BtPathname)
		usr("JABBERD_SPOOLDIR", lkNone, BtPathname)
		usr("JABBERD_PIDDIR", lkNone, BtPathname)
		usr("JAKARTA_HOME", lkNone, BtPathname)
		usr("KERBEROS", lkNone, BtYes)
		usr("KERMIT_SUID_UUCP", lkNone, BtYes)
		usr("KJS_USE_PCRE", lkNone, BtYes)
		usr("KNEWS_DOMAIN_FILE", lkNone, BtPathname)
		usr("KNEWS_DOMAIN_NAME", lkNone, BtIdentifier)
		usr("LIBDVDCSS_HOMEPAGE", lkNone, BtHomepage)
		usr("LIBDVDCSS_MASTER_SITES", lkShell, BtFetchURL)
		usr("LATEX2HTML_ICONPATH", lkNone, BtURL)
		usr("LEAFNODE_DATA_DIR", lkNone, BtPathname)
		usr("LEAFNODE_USER", lkNone, BtUserGroupName)
		usr("LEAFNODE_GROUP", lkNone, BtUserGroupName)
		usr("LINUX_LOCALES", lkShell, BtIdentifier)
		usr("MAILAGENT_DOMAIN", lkNone, BtIdentifier)
		usr("MAILAGENT_EMAIL", lkNone, BtMailAddress)
		usr("MAILAGENT_FQDN", lkNone, BtIdentifier)
		usr("MAILAGENT_ORGANIZATION", lkNone, BtUnknown)
		usr("MAJORDOMO_HOMEDIR", lkNone, BtPathname)
		usr("MAKEINFO_ARGS", lkShell, BtShellWord)
		usr("MECAB_CHARSET", lkNone, BtIdentifier)
		usr("MEDIATOMB_GROUP", lkNone, BtUserGroupName)
		usr("MEDIATOMB_USER", lkNone, BtUserGroupName)
		usr("MLDONKEY_GROUP", lkNone, BtUserGroupName)
		usr("MLDONKEY_HOME", lkNone, BtPathname)
		usr("MLDONKEY_USER", lkNone, BtUserGroupName)
		usr("MONOTONE_GROUP", lkNone, BtUserGroupName)
		usr("MONOTONE_USER", lkNone, BtUserGroupName)
		usr("MOTIF_TYPE", lkNone, enum("motif openmotif lesstif dt"))
		usr("MOTIF_TYPE_DEFAULT", lkNone, enum("motif openmotif lesstif dt"))
		usr("MTOOLS_ENABLE_FLOPPYD", lkNone, BtYesNo)
		usr("MYSQL_USER", lkNone, BtUserGroupName)
		usr("MYSQL_GROUP", lkNone, BtUserGroupName)
		usr("MYSQL_DATADIR", lkNone, BtPathname)
		usr("MYSQL_CHARSET", lkNone, BtIdentifier)
		usr("MYSQL_EXTRA_CHARSET", lkShell, BtIdentifier)
		usr("NAGIOS_GROUP", lkNone, BtUserGroupName)
		usr("NAGIOS_USER", lkNone, BtUserGroupName)
		usr("NAGIOSCMD_GROUP", lkNone, BtUserGroupName)
		usr("NAGIOSDIR", lkNone, BtPathname)
		usr("NBPAX_PROGRAM_PREFIX", lkNone, BtIdentifier)
		usr("NMH_EDITOR", lkNone, BtIdentifier)
		usr("NMH_MTA", lkNone, enum("smtp sendmail"))
		usr("NMH_PAGER", lkNone, BtIdentifier)
		usr("NS_PREFERRED", lkNone, enum("communicator navigator mozilla"))
		usr("OPENSSH_CHROOT", lkNone, BtPathname)
		usr("OPENSSH_USER", lkNone, BtUserGroupName)
		usr("OPENSSH_GROUP", lkNone, BtUserGroupName)
		usr("P4USER", lkNone, BtUserGroupName)
		usr("P4GROUP", lkNone, BtUserGroupName)
		usr("P4ROOT", lkNone, BtPathname)
		usr("P4PORT", lkNone, BtInteger)
		usr("PALMOS_DEFAULT_SDK", lkNone, enum("1 2 3.1 3.5"))
		usr("PAPERSIZE", lkNone, enum("A4 Letter"))
		usr("PGGROUP", lkNone, BtUserGroupName)
		usr("PGUSER", lkNone, BtUserGroupName)
		usr("PGHOME", lkNone, BtPathname)
		usr("PILRC_USE_GTK", lkNone, BtYesNo)
		usr("PKG_JVM_DEFAULT", lkNone, jvms)
		usr("POPTOP_USE_MPPE", lkNone, BtYes)
		usr("PROCMAIL_MAILSPOOLHOME", lkNone, BtFilename)
		usr("PROCMAIL_TRUSTED_IDS", lkShell, BtIdentifier)
		usr("PVM_SSH", lkNone, BtPathname)
		usr("QMAILDIR", lkNone, BtPathname)
		usr("QMAIL_QFILTER_TMPDIR", lkNone, BtPathname)
		usr("QMAIL_QUEUE_DIR", lkNone, BtPathname)
		usr("QMAIL_QUEUE_EXTRA", lkNone, BtMailAddress)
		usr("QPOPPER_FAC", lkNone, BtIdentifier)
		usr("QPOPPER_USER", lkNone, BtUserGroupName)
		usr("QPOPPER_SPOOL_DIR", lkNone, BtPathname)
		usr("RASMOL_DEPTH", lkNone, enum("8 16 32"))
		usr("RELAY_CTRL_DIR", lkNone, BtPathname)
		usr("RPM_DB_PREFIX", lkNone, BtPathname)
		usr("RSSH_SCP_PATH", lkNone, BtPathname)
		usr("RSSH_SFTP_SERVER_PATH", lkNone, BtPathname)
		usr("RSSH_CVS_PATH", lkNone, BtPathname)
		usr("RSSH_RDIST_PATH", lkNone, BtPathname)
		usr("RSSH_RSYNC_PATH", lkNone, BtPathname)
		usr("SAWFISH_THEMES", lkShell, BtFilename)
		usr("SCREWS_GROUP", lkNone, BtUserGroupName)
		usr("SCREWS_USER", lkNone, BtUserGroupName)
		usr("SDIST_PAWD", lkNone, enum("pawd pwd"))
		usr("SERIAL_DEVICES", lkShell, BtPathname)
		usr("SILC_CLIENT_WITH_PERL", lkNone, BtYesNo)
		usr("SSH_SUID", lkNone, BtYesNo)
		usr("SSYNC_PAWD", lkNone, enum("pawd pwd"))
		usr("SUSE_PREFER", lkNone, enum("13.1 12.1 10.0"))
		usr("TEXMFSITE", lkNone, BtPathname)
		usr("THTTPD_LOG_FACILITY", lkNone, BtIdentifier)
		usr("UNPRIVILEGED", lkNone, BtYesNo)
		usr("USE_CROSS_COMPILE", lkNone, BtYesNo)
		usr("USERPPP_GROUP", lkNone, BtUserGroupName)
		usr("UUCP_GROUP", lkNone, BtUserGroupName)
		usr("UUCP_USER", lkNone, BtUserGroupName)
		usr("VIM_EXTRA_OPTS", lkShell, BtShellWord)
		usr("WCALC_HTMLDIR", lkNone, BtPathname)
		usr("WCALC_HTMLPATH", lkNone, BtPathname) // URL path
		usr("WCALC_CGIDIR", lkNone, BtPrefixPathname)
		usr("WCALC_CGIPATH", lkNone, BtPathname) // URL path
		usr("WDM_MANAGERS", lkShell, BtIdentifier)
		usr("X10_PORT", lkNone, BtPathname)
		usr("XAW_TYPE", lkNone, enum("standard 3d xpm neXtaw"))
		usr("XLOCK_DEFAULT_MODE", lkNone, BtIdentifier)
		usr("ZSH_STATIC", lkNone, BtYes)
	}

	// some other variables, sorted alphabetically

	acl(".CURDIR", lkNone, BtPathname, "buildlink3.mk:; *: use, use-loadtime")
	acl(".TARGET", lkNone, BtPathname, "buildlink3.mk:; *: use, use-loadtime")
	acl("ALL_ENV", lkShell, BtShellWord, "")
	acl("ALTERNATIVES_FILE", lkNone, BtFilename, "")
	acl("ALTERNATIVES_SRC", lkShell, BtPathname, "")
	pkg("APACHE_MODULE", lkNone, BtYes)
	sys("AR", lkNone, BtShellCommand)
	sys("AS", lkNone, BtShellCommand)
	pkglist("AUTOCONF_REQD", lkShell, BtVersion)
	acl("AUTOMAKE_OVERRIDE", lkShell, BtPathmask, "")
	pkglist("AUTOMAKE_REQD", lkShell, BtVersion)
	pkg("AUTO_MKDIRS", lkNone, BtYesNo)
	usr("BATCH", lkNone, BtYes)
	acl("BDB185_DEFAULT", lkNone, BtUnknown, "")
	sys("BDBBASE", lkNone, BtPathname)
	pkg("BDB_ACCEPTED", lkShell, enum("db1 db2 db3 db4 db5 db6"))
	acl("BDB_DEFAULT", lkNone, enum("db1 db2 db3 db4 db5 db6"), "")
	sys("BDB_LIBS", lkShell, BtLdFlag)
	sys("BDB_TYPE", lkNone, enum("db1 db2 db3 db4 db5 db6"))
	sys("BINGRP", lkNone, BtUserGroupName)
	sys("BINMODE", lkNone, BtFileMode)
	sys("BINOWN", lkNone, BtUserGroupName)
	acl("BOOTSTRAP_DEPENDS", lkSpace, BtDependencyWithPath, "Makefile, Makefile.common, *.mk: append")
	pkg("BOOTSTRAP_PKG", lkNone, BtYesNo)
	acl("BROKEN", lkNone, BtMessage, "")
	pkg("BROKEN_GETTEXT_DETECTION", lkNone, BtYesNo)
	pkglist("BROKEN_EXCEPT_ON_PLATFORM", lkShell, BtMachinePlatformPattern)
	pkglist("BROKEN_ON_PLATFORM", lkSpace, BtMachinePlatformPattern)
	sys("BSD_MAKE_ENV", lkShell, BtShellWord)
	acl("BUILDLINK_ABI_DEPENDS.*", lkSpace, BtDependency, "builtin.mk: append, use-loadtime; *: append")
	acl("BUILDLINK_API_DEPENDS.*", lkSpace, BtDependency, "builtin.mk: append, use-loadtime; *: append")
	acl("BUILDLINK_AUTO_DIRS.*", lkNone, BtYesNo, "buildlink3.mk: append")
	acl("BUILDLINK_CONTENTS_FILTER", lkNone, BtShellCommand, "")
	sys("BUILDLINK_CFLAGS", lkShell, BtCFlag)
	bl3list("BUILDLINK_CFLAGS.*", lkShell, BtCFlag)
	sys("BUILDLINK_CPPFLAGS", lkShell, BtCFlag)
	bl3list("BUILDLINK_CPPFLAGS.*", lkShell, BtCFlag)
	acl("BUILDLINK_CONTENTS_FILTER.*", lkNone, BtShellCommand, "buildlink3.mk: set")
	acl("BUILDLINK_DEPENDS", lkSpace, BtIdentifier, "buildlink3.mk: append")
	acl("BUILDLINK_DEPMETHOD.*", lkShell, BtBuildlinkDepmethod, "buildlink3.mk: default, append, use; Makefile: set, append; Makefile.common, *.mk: append")
	acl("BUILDLINK_DIR", lkNone, BtPathname, "*: use")
	bl3list("BUILDLINK_FILES.*", lkShell, BtPathmask)
	acl("BUILDLINK_FILES_CMD.*", lkNone, BtShellCommand, "")
	acl("BUILDLINK_INCDIRS.*", lkShell, BtPathname, "buildlink3.mk: default, append; Makefile, Makefile.common, *.mk: use")
	acl("BUILDLINK_JAVA_PREFIX.*", lkNone, BtPathname, "buildlink3.mk: set, use")
	acl("BUILDLINK_LDADD.*", lkShell, BtLdFlag, "builtin.mk: set, default, append, use; buildlink3.mk: append, use; Makefile, Makefile.common, *.mk: use")
	acl("BUILDLINK_LDFLAGS", lkShell, BtLdFlag, "*: use")
	bl3list("BUILDLINK_LDFLAGS.*", lkShell, BtLdFlag)
	acl("BUILDLINK_LIBDIRS.*", lkShell, BtPathname, "buildlink3.mk, builtin.mk: append; Makefile, Makefile.common, *.mk: use")
	acl("BUILDLINK_LIBS.*", lkShell, BtLdFlag, "buildlink3.mk: append")
	acl("BUILDLINK_PASSTHRU_DIRS", lkShell, BtPathname, "Makefile, Makefile.common, buildlink3.mk, hacks.mk: append")
	acl("BUILDLINK_PASSTHRU_RPATHDIRS", lkShell, BtPathname, "Makefile, Makefile.common, buildlink3.mk, hacks.mk: append")
	acl("BUILDLINK_PKGSRCDIR.*", lkNone, BtRelativePkgDir, "buildlink3.mk: default, use-loadtime")
	acl("BUILDLINK_PREFIX.*", lkNone, BtPathname, "builtin.mk: set, use; Makefile, Makefile.common, *.mk: use")
	acl("BUILDLINK_RPATHDIRS.*", lkShell, BtPathname, "buildlink3.mk: append")
	acl("BUILDLINK_TARGETS", lkShell, BtIdentifier, "")
	acl("BUILDLINK_FNAME_TRANSFORM.*", lkNone, BtSedCommands, "Makefile, buildlink3.mk, builtin.mk, hacks.mk: append")
	acl("BUILDLINK_TRANSFORM", lkShell, BtWrapperTransform, "*: append")
	acl("BUILDLINK_TRANSFORM.*", lkShell, BtWrapperTransform, "*: append")
	acl("BUILDLINK_TREE", lkShell, BtIdentifier, "buildlink3.mk: append")
	acl("BUILD_DEFS", lkShell, BtVariableName, "Makefile, Makefile.common, options.mk: append")
	pkglist("BUILD_DEFS_EFFECTS", lkShell, BtVariableName)
	acl("BUILD_DEPENDS", lkSpace, BtDependencyWithPath, "Makefile, Makefile.common, *.mk: append")
	pkglist("BUILD_DIRS", lkShell, BtWrksrcSubdirectory)
	pkglist("BUILD_ENV", lkShell, BtShellWord)
	sys("BUILD_MAKE_CMD", lkNone, BtShellCommand)
	pkglist("BUILD_MAKE_FLAGS", lkShell, BtShellWord)
	pkglist("BUILD_TARGET", lkShell, BtIdentifier)
	pkglist("BUILD_TARGET.*", lkShell, BtIdentifier)
	pkg("BUILD_USES_MSGFMT", lkNone, BtYes)
	acl("BUILTIN_PKG", lkNone, BtIdentifier, "builtin.mk: set, use-loadtime, use")
	acl("BUILTIN_PKG.*", lkNone, BtPkgName, "builtin.mk: set, use-loadtime, use")
	acl("BUILTIN_FIND_FILES_VAR", lkShell, BtVariableName, "builtin.mk: set")
	acl("BUILTIN_FIND_FILES.*", lkShell, BtPathname, "builtin.mk: set")
	acl("BUILTIN_FIND_GREP.*", lkNone, BtUnknown, "builtin.mk: set")
	acl("BUILTIN_FIND_HEADERS_VAR", lkShell, BtVariableName, "builtin.mk: set")
	acl("BUILTIN_FIND_HEADERS.*", lkShell, BtPathname, "builtin.mk: set")
	acl("BUILTIN_FIND_LIBS", lkShell, BtPathname, "builtin.mk: set")
	acl("BUILTIN_IMAKE_CHECK", lkShell, BtUnknown, "builtin.mk: set")
	acl("BUILTIN_IMAKE_CHECK.*", lkNone, BtYesNo, "")
	sys("BUILTIN_X11_TYPE", lkNone, BtUnknown)
	sys("BUILTIN_X11_VERSION", lkNone, BtUnknown)
	acl("CATEGORIES", lkShell, BtCategory, "Makefile: set, append; Makefile.common: set, default, append")
	sys("CC_VERSION", lkNone, BtMessage)
	sys("CC", lkNone, BtShellCommand)
	pkglist("CFLAGS", lkShell, BtCFlag)   // may also be changed by the user
	pkglist("CFLAGS.*", lkShell, BtCFlag) // may also be changed by the user
	acl("CHECK_BUILTIN", lkNone, BtYesNo, "builtin.mk: default; Makefile: set")
	acl("CHECK_BUILTIN.*", lkNone, BtYesNo, "Makefile, options.mk, buildlink3.mk: set; builtin.mk: default; *: use-loadtime")
	acl("CHECK_FILES_SKIP", lkShell, BtBasicRegularExpression, "Makefile, Makefile.common: append")
	pkg("CHECK_FILES_SUPPORTED", lkNone, BtYesNo)
	usr("CHECK_HEADERS", lkNone, BtYesNo)
	pkglist("CHECK_HEADERS_SKIP", lkShell, BtPathmask)
	usr("CHECK_INTERPRETER", lkNone, BtYesNo)
	pkglist("CHECK_INTERPRETER_SKIP", lkShell, BtPathmask)
	usr("CHECK_PERMS", lkNone, BtYesNo)
	pkglist("CHECK_PERMS_SKIP", lkShell, BtPathmask)
	usr("CHECK_PORTABILITY", lkNone, BtYesNo)
	pkglist("CHECK_PORTABILITY_SKIP", lkShell, BtPathmask)
	usr("CHECK_RELRO", lkNone, BtYesNo)
	pkglist("CHECK_RELRO_SKIP", lkShell, BtPathmask)
	pkg("CHECK_RELRO_SUPPORTED", lkNone, BtYesNo)
	acl("CHECK_SHLIBS", lkNone, BtYesNo, "Makefile: set")
	pkglist("CHECK_SHLIBS_SKIP", lkShell, BtPathmask)
	acl("CHECK_SHLIBS_SUPPORTED", lkNone, BtYesNo, "Makefile: set")
	pkglist("CHECK_WRKREF_SKIP", lkShell, BtPathmask)
	pkg("CMAKE_ARG_PATH", lkNone, BtPathname)
	pkglist("CMAKE_ARGS", lkShell, BtShellWord)
	pkglist("CMAKE_ARGS.*", lkShell, BtShellWord)
	acl("COMMENT", lkNone, BtComment, "Makefile, Makefile.common: set, append")
	acl("COMPILER_RPATH_FLAG", lkNone, enum("-Wl,-rpath"), "*: use")
	pkglist("CONFIGURE_ARGS", lkShell, BtShellWord)
	pkglist("CONFIGURE_ARGS.*", lkShell, BtShellWord)
	pkglist("CONFIGURE_DIRS", lkShell, BtWrksrcSubdirectory)
	acl("CONFIGURE_ENV", lkShell, BtShellWord, "Makefile, Makefile.common: append, set, use; buildlink3.mk, builtin.mk: append; *.mk: append, use")
	acl("CONFIGURE_ENV.*", lkShell, BtShellWord, "Makefile, Makefile.common: append, set, use; buildlink3.mk, builtin.mk: append; *.mk: append, use")
	pkg("CONFIGURE_HAS_INFODIR", lkNone, BtYesNo)
	pkg("CONFIGURE_HAS_LIBDIR", lkNone, BtYesNo)
	pkg("CONFIGURE_HAS_MANDIR", lkNone, BtYesNo)
	pkg("CONFIGURE_SCRIPT", lkNone, BtPathname)
	acl("CONFIG_GUESS_OVERRIDE", lkShell, BtPathmask, "Makefile, Makefile.common: set, append")
	acl("CONFIG_STATUS_OVERRIDE", lkShell, BtPathmask, "Makefile, Makefile.common: set, append")
	acl("CONFIG_SHELL", lkNone, BtPathname, "Makefile, Makefile.common: set")
	acl("CONFIG_SUB_OVERRIDE", lkShell, BtPathmask, "Makefile, Makefile.common: set, append")
	pkglist("CONFLICTS", lkSpace, BtDependency)
	pkglist("CONF_FILES", lkNone, BtConfFiles)
	pkg("CONF_FILES_MODE", lkNone, enum("0644 0640 0600 0400"))
	pkglist("CONF_FILES_PERMS", lkShell, BtPerms)
	sys("COPY", lkNone, enum("-c")) // The flag that tells ${INSTALL} to copy a file
	sys("CPP", lkNone, BtShellCommand)
	pkglist("CPPFLAGS", lkShell, BtCFlag)
	pkglist("CPPFLAGS.*", lkShell, BtCFlag)
	sys("CXX", lkNone, BtShellCommand)
	pkglist("CXXFLAGS", lkShell, BtCFlag)
	pkglist("CXXFLAGS.*", lkShell, BtCFlag)
	pkglist("CWRAPPERS_APPEND.*", lkShell, BtShellWord)
	acl("DEINSTALL_FILE", lkNone, BtPathname, "Makefile: set")
	acl("DEINSTALL_SRC", lkShell, BtPathname, "Makefile: set; Makefile.common: default, set")
	acl("DEINSTALL_TEMPLATES", lkShell, BtPathname, "Makefile: set, append; Makefile.common: set, default, append")
	sys("DELAYED_ERROR_MSG", lkNone, BtShellCommand)
	sys("DELAYED_WARNING_MSG", lkNone, BtShellCommand)
	pkglist("DEPENDS", lkSpace, BtDependencyWithPath)
	usr("DEPENDS_TARGET", lkShell, BtIdentifier)
	acl("DESCR_SRC", lkShell, BtPathname, "Makefile: set, append; Makefile.common: default, set")
	sys("DESTDIR", lkNone, BtPathname)
	acl("DESTDIR_VARNAME", lkNone, BtVariableName, "Makefile, Makefile.common: set")
	sys("DEVOSSAUDIO", lkNone, BtPathname)
	sys("DEVOSSSOUND", lkNone, BtPathname)
	pkglist("DISTFILES", lkShell, BtFilename)
	pkg("DISTINFO_FILE", lkNone, BtRelativePkgPath)
	pkg("DISTNAME", lkNone, BtFilename)
	pkg("DIST_SUBDIR", lkNone, BtPathname)
	acl("DJB_BUILD_ARGS", lkShell, BtShellWord, "")
	acl("DJB_BUILD_TARGETS", lkShell, BtIdentifier, "")
	acl("DJB_CONFIG_CMDS", lkNone, BtShellCommands, "options.mk: set")
	acl("DJB_CONFIG_DIRS", lkShell, BtWrksrcSubdirectory, "")
	acl("DJB_CONFIG_HOME", lkNone, BtFilename, "")
	acl("DJB_CONFIG_PREFIX", lkNone, BtPathname, "")
	acl("DJB_INSTALL_TARGETS", lkShell, BtIdentifier, "")
	acl("DJB_MAKE_TARGETS", lkNone, BtYesNo, "")
	acl("DJB_RESTRICTED", lkNone, BtYesNo, "Makefile: set")
	acl("DJB_SLASHPACKAGE", lkNone, BtYesNo, "")
	acl("DLOPEN_REQUIRE_PTHREADS", lkNone, BtYesNo, "")
	acl("DL_AUTO_VARS", lkNone, BtYes, "Makefile, Makefile.common, options.mk: set")
	acl("DL_LIBS", lkShell, BtLdFlag, "")
	sys("DOCOWN", lkNone, BtUserGroupName)
	sys("DOCGRP", lkNone, BtUserGroupName)
	sys("DOCMODE", lkNone, BtFileMode)
	sys("DOWNLOADED_DISTFILE", lkNone, BtPathname)
	sys("DO_NADA", lkNone, BtShellCommand)
	pkg("DYNAMIC_SITES_CMD", lkNone, BtShellCommand)
	pkg("DYNAMIC_SITES_SCRIPT", lkNone, BtPathname)
	acl("ECHO", lkNone, BtShellCommand, "*: use")
	sys("ECHO_MSG", lkNone, BtShellCommand)
	sys("ECHO_N", lkNone, BtShellCommand)
	pkg("EGDIR", lkNone, BtPathname) // Not defined anywhere, but used in many places like this.
	sys("EMACS_BIN", lkNone, BtPathname)
	sys("EMACS_ETCPREFIX", lkNone, BtPathname)
	sys("EMACS_FLAVOR", lkNone, enum("emacs xemacs"))
	sys("EMACS_INFOPREFIX", lkNone, BtPathname)
	sys("EMACS_LISPPREFIX", lkNone, BtPathname)
	acl("EMACS_MODULES", lkShell, BtIdentifier, "Makefile, Makefile.common: set, append")
	sys("EMACS_PKGNAME_PREFIX", lkNone, BtIdentifier) // Or the empty string.
	sys("EMACS_TYPE", lkNone, enum("emacs xemacs"))
	acl("EMACS_USE_LEIM", lkNone, BtYes, "")
	acl("EMACS_VERSIONS_ACCEPTED", lkShell, emacsVersions, "Makefile: set")
	sys("EMACS_VERSION_MAJOR", lkNone, BtInteger)
	sys("EMACS_VERSION_MINOR", lkNone, BtInteger)
	acl("EMACS_VERSION_REQD", lkShell, emacsVersions, "Makefile: set, append")
	sys("EMULDIR", lkNone, BtPathname)
	sys("EMULSUBDIR", lkNone, BtPathname)
	sys("OPSYS_EMULDIR", lkNone, BtPathname)
	sys("EMULSUBDIRSLASH", lkNone, BtPathname)
	sys("EMUL_ARCH", lkNone, enum("arm i386 m68k none ns32k sparc vax x86_64"))
	sys("EMUL_DISTRO", lkNone, BtIdentifier)
	sys("EMUL_IS_NATIVE", lkNone, BtYes)
	pkg("EMUL_MODULES.*", lkShell, BtIdentifier)
	sys("EMUL_OPSYS", lkNone, enum("darwin freebsd hpux irix linux osf1 solaris sunos none"))
	pkg("EMUL_PKG_FMT", lkNone, enum("plain rpm"))
	usr("EMUL_PLATFORM", lkNone, BtEmulPlatform)
	pkg("EMUL_PLATFORMS", lkShell, BtEmulPlatform)
	usr("EMUL_PREFER", lkShell, BtEmulPlatform)
	pkg("EMUL_REQD", lkSpace, BtDependency)
	usr("EMUL_TYPE.*", lkNone, enum("native builtin suse suse-9.1 suse-9.x suse-10.0 suse-10.x"))
	sys("ERROR_CAT", lkNone, BtShellCommand)
	sys("ERROR_MSG", lkNone, BtShellCommand)
	acl("EVAL_PREFIX", lkSpace, BtShellWord, "Makefile, Makefile.common: append") // XXX: Combining ShellWord with lkSpace looks weird.
	sys("EXPORT_SYMBOLS_LDFLAGS", lkShell, BtLdFlag)
	sys("EXTRACT_CMD", lkNone, BtShellCommand)
	pkg("EXTRACT_DIR", lkNone, BtPathname)
	pkg("EXTRACT_DIR.*", lkNone, BtPathname)
	pkglist("EXTRACT_ELEMENTS", lkShell, BtPathmask)
	pkglist("EXTRACT_ENV", lkShell, BtShellWord)
	pkglist("EXTRACT_ONLY", lkShell, BtPathname)
	acl("EXTRACT_OPTS", lkShell, BtShellWord, "Makefile, Makefile.common: set, append")
	acl("EXTRACT_OPTS_BIN", lkShell, BtShellWord, "Makefile, Makefile.common: set, append")
	acl("EXTRACT_OPTS_LHA", lkShell, BtShellWord, "Makefile, Makefile.common: set, append")
	acl("EXTRACT_OPTS_PAX", lkShell, BtShellWord, "Makefile, Makefile.common: set, append")
	acl("EXTRACT_OPTS_RAR", lkShell, BtShellWord, "Makefile, Makefile.common: set, append")
	acl("EXTRACT_OPTS_TAR", lkShell, BtShellWord, "Makefile, Makefile.common: set, append")
	acl("EXTRACT_OPTS_ZIP", lkShell, BtShellWord, "Makefile, Makefile.common: set, append")
	acl("EXTRACT_OPTS_ZOO", lkShell, BtShellWord, "Makefile, Makefile.common: set, append")
	pkg("EXTRACT_SUFX", lkNone, BtDistSuffix)
	pkg("EXTRACT_USING", lkNone, enum("bsdtar gtar nbtar pax"))
	sys("FAIL_MSG", lkNone, BtShellCommand)
	sys("FAMBASE", lkNone, BtPathname)
	pkg("FAM_ACCEPTED", lkShell, enum("fam gamin"))
	usr("FAM_DEFAULT", lkNone, enum("fam gamin"))
	sys("FAM_TYPE", lkNone, enum("fam gamin"))
	acl("FETCH_BEFORE_ARGS", lkShell, BtShellWord, "Makefile: set, append")
	pkglist("FETCH_MESSAGE", lkShell, BtShellWord)
	pkg("FILESDIR", lkNone, BtRelativePkgPath)
	pkglist("FILES_SUBST", lkShell, BtShellWord)
	acl("FILES_SUBST_SED", lkShell, BtShellWord, "")
	pkglist("FIX_RPATH", lkShell, BtVariableName)
	pkglist("FLEX_REQD", lkShell, BtVersion)
	acl("FONTS_DIRS.*", lkShell, BtPathname, "Makefile: set, append, use; Makefile.common: append, use")
	sys("GAMEDATAMODE", lkNone, BtFileMode)
	sys("GAMES_GROUP", lkNone, BtUserGroupName)
	sys("GAMEMODE", lkNone, BtFileMode)
	sys("GAMES_USER", lkNone, BtUserGroupName)
	pkglist("GCC_REQD", lkShell, BtVersion)
	pkglist("GENERATE_PLIST", lkNone, BtShellCommands)
	pkg("GITHUB_PROJECT", lkNone, BtIdentifier)
	pkg("GITHUB_TAG", lkNone, BtIdentifier)
	pkg("GITHUB_RELEASE", lkNone, BtFilename)
	pkg("GITHUB_TYPE", lkNone, enum("tag release"))
	pkg("GMAKE_REQD", lkNone, BtVersion)
	acl("GNU_ARCH", lkNone, enum("mips"), "")
	acl("GNU_ARCH.*", lkNone, BtIdentifier, "buildlink3.mk:; *: set, use")
	acl("GNU_CONFIGURE", lkNone, BtYes, "Makefile, Makefile.common: set")
	acl("GNU_CONFIGURE_INFODIR", lkNone, BtPathname, "Makefile, Makefile.common: set")
	acl("GNU_CONFIGURE_LIBDIR", lkNone, BtPathname, "Makefile, Makefile.common: set")
	pkg("GNU_CONFIGURE_LIBSUBDIR", lkNone, BtPathname)
	acl("GNU_CONFIGURE_MANDIR", lkNone, BtPathname, "Makefile, Makefile.common: set")
	acl("GNU_CONFIGURE_PREFIX", lkNone, BtPathname, "Makefile: set")
	pkg("GOPATH", lkNone, BtPathname)
	acl("HAS_CONFIGURE", lkNone, BtYes, "Makefile, Makefile.common: set")
	pkglist("HEADER_TEMPLATES", lkShell, BtPathname)
	pkg("HOMEPAGE", lkNone, BtHomepage)
	pkg("ICON_THEMES", lkNone, BtYes)
	acl("IGNORE_PKG.*", lkNone, BtYes, "*: set, use-loadtime")
	sys("IMAKE", lkNone, BtShellCommand)
	acl("INCOMPAT_CURSES", lkSpace, BtMachinePlatformPattern, "Makefile: set, append")
	acl("INCOMPAT_ICONV", lkSpace, BtMachinePlatformPattern, "")
	acl("INFO_DIR", lkNone, BtPathname, "") // relative to PREFIX
	pkg("INFO_FILES", lkNone, BtYes)
	sys("INSTALL", lkNone, BtShellCommand)
	pkglist("INSTALLATION_DIRS", lkShell, BtPrefixPathname)
	pkg("INSTALLATION_DIRS_FROM_PLIST", lkNone, BtYes)
	sys("INSTALL_DATA", lkNone, BtShellCommand)
	sys("INSTALL_DATA_DIR", lkNone, BtShellCommand)
	pkglist("INSTALL_DIRS", lkShell, BtWrksrcSubdirectory)
	pkglist("INSTALL_ENV", lkShell, BtShellWord)
	acl("INSTALL_FILE", lkNone, BtPathname, "Makefile: set")
	sys("INSTALL_GAME", lkNone, BtShellCommand)
	sys("INSTALL_GAME_DATA", lkNone, BtShellCommand)
	sys("INSTALL_LIB", lkNone, BtShellCommand)
	sys("INSTALL_LIB_DIR", lkNone, BtShellCommand)
	pkglist("INSTALL_MAKE_FLAGS", lkShell, BtShellWord)
	sys("INSTALL_MAN", lkNone, BtShellCommand)
	sys("INSTALL_MAN_DIR", lkNone, BtShellCommand)
	sys("INSTALL_PROGRAM", lkNone, BtShellCommand)
	sys("INSTALL_PROGRAM_DIR", lkNone, BtShellCommand)
	sys("INSTALL_SCRIPT", lkNone, BtShellCommand)
	acl("INSTALL_SCRIPTS_ENV", lkShell, BtShellWord, "")
	sys("INSTALL_SCRIPT_DIR", lkNone, BtShellCommand)
	acl("INSTALL_SRC", lkShell, BtPathname, "Makefile: set; Makefile.common: default, set")
	pkg("INSTALL_TARGET", lkShell, BtIdentifier)
	acl("INSTALL_TEMPLATES", lkShell, BtPathname, "Makefile: set, append; Makefile.common: set, default, append")
	acl("INSTALL_UNSTRIPPED", lkNone, BtYesNo, "Makefile, Makefile.common: set")
	pkg("INTERACTIVE_STAGE", lkShell, enum("fetch extract configure build test install"))
	acl("IS_BUILTIN.*", lkNone, BtYesNoIndirectly, "builtin.mk: set, use-loadtime, use")
	sys("JAVA_BINPREFIX", lkNone, BtPathname)
	pkg("JAVA_CLASSPATH", lkNone, BtShellWord)
	pkg("JAVA_HOME", lkNone, BtPathname)
	pkg("JAVA_NAME", lkNone, BtFilename)
	pkglist("JAVA_UNLIMIT", lkShell, enum("cmdsize datasize stacksize"))
	pkglist("JAVA_WRAPPERS", lkSpace, BtFilename)
	pkg("JAVA_WRAPPER_BIN.*", lkNone, BtPathname)
	sys("KRB5BASE", lkNone, BtPathname)
	acl("KRB5_ACCEPTED", lkShell, enum("heimdal mit-krb5"), "")
	usr("KRB5_DEFAULT", lkNone, enum("heimdal mit-krb5"))
	sys("KRB5_TYPE", lkNone, BtIdentifier)
	sys("LD", lkNone, BtShellCommand)
	pkglist("LDFLAGS", lkShell, BtLdFlag)
	pkglist("LDFLAGS.*", lkShell, BtLdFlag)
	sys("LIBGRP", lkNone, BtUserGroupName)
	sys("LIBMODE", lkNone, BtFileMode)
	sys("LIBOWN", lkNone, BtUserGroupName)
	sys("LIBOSSAUDIO", lkNone, BtPathname)
	pkglist("LIBS", lkShell, BtLdFlag)
	pkglist("LIBS.*", lkShell, BtLdFlag)
	sys("LIBTOOL", lkNone, BtShellCommand)
	acl("LIBTOOL_OVERRIDE", lkShell, BtPathmask, "Makefile: set, append")
	pkglist("LIBTOOL_REQD", lkShell, BtVersion)
	acl("LICENCE", lkNone, BtLicense, "Makefile, Makefile.common, options.mk: set, append")
	acl("LICENSE", lkNone, BtLicense, "Makefile, Makefile.common, options.mk: set, append")
	pkg("LICENSE_FILE", lkNone, BtPathname)
	sys("LINKER_RPATH_FLAG", lkNone, BtShellWord)
	sys("LOWER_OPSYS", lkNone, BtIdentifier)
	sys("LOWER_VENDOR", lkNone, BtIdentifier)
	acl("LTCONFIG_OVERRIDE", lkShell, BtPathmask, "Makefile: set, append; Makefile.common: append")
	sys("MACHINE_ARCH", lkNone, enumMachineArch)
	sys("MACHINE_GNU_ARCH", lkNone, enumMachineGnuArch)
	sys("MACHINE_GNU_PLATFORM", lkNone, BtMachineGnuPlatform)
	sys("MACHINE_PLATFORM", lkNone, BtMachinePlatform)
	acl("MAINTAINER", lkNone, BtMailAddress, "Makefile: set; Makefile.common: default")
	sys("MAKE", lkNone, BtShellCommand)
	pkglist("MAKEFLAGS", lkShell, BtShellWord)
	acl("MAKEVARS", lkShell, BtVariableName, "buildlink3.mk, builtin.mk, hacks.mk: append")
	pkglist("MAKE_DIRS", lkShell, BtPathname)
	pkglist("MAKE_DIRS_PERMS", lkShell, BtPerms)
	acl("MAKE_ENV", lkShell, BtShellWord, "Makefile, Makefile.common: append, set, use; buildlink3.mk, builtin.mk: append; *.mk: append, use")
	acl("MAKE_ENV.*", lkShell, BtShellWord, "Makefile, Makefile.common: append, set, use; buildlink3.mk, builtin.mk: append; *.mk: append, use")
	pkg("MAKE_FILE", lkNone, BtPathname)
	pkglist("MAKE_FLAGS", lkShell, BtShellWord)
	pkglist("MAKE_FLAGS.*", lkShell, BtShellWord)
	usr("MAKE_JOBS", lkNone, BtInteger)
	pkg("MAKE_JOBS_SAFE", lkNone, BtYesNo)
	pkg("MAKE_PROGRAM", lkNone, BtShellCommand)
	acl("MANCOMPRESSED", lkNone, BtYesNo, "Makefile: set; Makefile.common: default, set")
	acl("MANCOMPRESSED_IF_MANZ", lkNone, BtYes, "Makefile: set; Makefile.common: default, set")
	sys("MANGRP", lkNone, BtUserGroupName)
	sys("MANMODE", lkNone, BtFileMode)
	sys("MANOWN", lkNone, BtUserGroupName)
	pkglist("MASTER_SITES", lkShell, BtFetchURL)
	sys("MASTER_SITE_APACHE", lkShell, BtFetchURL)
	sys("MASTER_SITE_BACKUP", lkShell, BtFetchURL)
	sys("MASTER_SITE_CYGWIN", lkShell, BtFetchURL)
	sys("MASTER_SITE_DEBIAN", lkShell, BtFetchURL)
	sys("MASTER_SITE_FREEBSD", lkShell, BtFetchURL)
	sys("MASTER_SITE_FREEBSD_LOCAL", lkShell, BtFetchURL)
	sys("MASTER_SITE_GENTOO", lkShell, BtFetchURL)
	sys("MASTER_SITE_GITHUB", lkShell, BtFetchURL)
	sys("MASTER_SITE_GNOME", lkShell, BtFetchURL)
	sys("MASTER_SITE_GNU", lkShell, BtFetchURL)
	sys("MASTER_SITE_GNUSTEP", lkShell, BtFetchURL)
	sys("MASTER_SITE_IFARCHIVE", lkShell, BtFetchURL)
	sys("MASTER_SITE_HASKELL_HACKAGE", lkShell, BtFetchURL)
	sys("MASTER_SITE_KDE", lkShell, BtFetchURL)
	sys("MASTER_SITE_LOCAL", lkShell, BtFetchURL)
	sys("MASTER_SITE_MOZILLA", lkShell, BtFetchURL)
	sys("MASTER_SITE_MOZILLA_ALL", lkShell, BtFetchURL)
	sys("MASTER_SITE_MOZILLA_ESR", lkShell, BtFetchURL)
	sys("MASTER_SITE_MYSQL", lkShell, BtFetchURL)
	sys("MASTER_SITE_NETLIB", lkShell, BtFetchURL)
	sys("MASTER_SITE_OPENOFFICE", lkShell, BtFetchURL)
	sys("MASTER_SITE_OSDN", lkShell, BtFetchURL)
	sys("MASTER_SITE_PERL_CPAN", lkShell, BtFetchURL)
	sys("MASTER_SITE_R_CRAN", lkShell, BtFetchURL)
	sys("MASTER_SITE_RUBYGEMS", lkShell, BtFetchURL)
	sys("MASTER_SITE_SOURCEFORGE", lkShell, BtFetchURL)
	sys("MASTER_SITE_SUNSITE", lkShell, BtFetchURL)
	sys("MASTER_SITE_SUSE", lkShell, BtFetchURL)
	sys("MASTER_SITE_TEX_CTAN", lkShell, BtFetchURL)
	sys("MASTER_SITE_XCONTRIB", lkShell, BtFetchURL)
	sys("MASTER_SITE_XEMACS", lkShell, BtFetchURL)
	pkglist("MESSAGE_SRC", lkShell, BtPathname)
	acl("MESSAGE_SUBST", lkShell, BtShellWord, "Makefile, Makefile.common, options.mk: append")
	pkg("META_PACKAGE", lkNone, BtYes)
	sys("MISSING_FEATURES", lkShell, BtIdentifier)
	acl("MYSQL_VERSIONS_ACCEPTED", lkShell, mysqlVersions, "Makefile: set")
	usr("MYSQL_VERSION_DEFAULT", lkNone, BtVersion)
	sys("NM", lkNone, BtShellCommand)
	sys("NONBINMODE", lkNone, BtFileMode)
	pkg("NOT_FOR_COMPILER", lkShell, compilers)
	pkglist("NOT_FOR_BULK_PLATFORM", lkSpace, BtMachinePlatformPattern)
	pkglist("NOT_FOR_PLATFORM", lkSpace, BtMachinePlatformPattern)
	pkg("NOT_FOR_UNPRIVILEGED", lkNone, BtYesNo)
	pkglist("NOT_PAX_ASLR_SAFE", lkShell, BtPathmask)
	pkglist("NOT_PAX_MPROTECT_SAFE", lkShell, BtPathmask)
	acl("NO_BIN_ON_CDROM", lkNone, BtRestricted, "Makefile, Makefile.common: set")
	acl("NO_BIN_ON_FTP", lkNone, BtRestricted, "Makefile, Makefile.common: set")
	acl("NO_BUILD", lkNone, BtYes, "Makefile, Makefile.common: set; Makefile.*: default, set")
	pkg("NO_CHECKSUM", lkNone, BtYes)
	pkg("NO_CONFIGURE", lkNone, BtYes)
	acl("NO_EXPORT_CPP", lkNone, BtYes, "Makefile: set")
	pkg("NO_EXTRACT", lkNone, BtYes)
	pkg("NO_INSTALL_MANPAGES", lkNone, BtYes) // only has an effect for Imake packages.
	acl("NO_PKGTOOLS_REQD_CHECK", lkNone, BtYes, "Makefile: set")
	acl("NO_SRC_ON_CDROM", lkNone, BtRestricted, "Makefile, Makefile.common: set")
	acl("NO_SRC_ON_FTP", lkNone, BtRestricted, "Makefile, Makefile.common: set")
	pkglist("ONLY_FOR_COMPILER", lkShell, compilers)
	pkglist("ONLY_FOR_PLATFORM", lkSpace, BtMachinePlatformPattern)
	pkg("ONLY_FOR_UNPRIVILEGED", lkNone, BtYesNo)
	sys("OPSYS", lkNone, BtIdentifier)
	acl("OPSYSVARS", lkShell, BtVariableName, "Makefile, Makefile.common: append")
	acl("OSVERSION_SPECIFIC", lkNone, BtYes, "Makefile, Makefile.common: set")
	sys("OS_VERSION", lkNone, BtVersion)
	pkg("OVERRIDE_DIRDEPTH*", lkNone, BtInteger)
	pkg("OVERRIDE_GNU_CONFIG_SCRIPTS", lkNone, BtYes)
	acl("OWNER", lkNone, BtMailAddress, "Makefile: set; Makefile.common: default")
	pkglist("OWN_DIRS", lkShell, BtPathname)
	pkglist("OWN_DIRS_PERMS", lkShell, BtPerms)
	sys("PAMBASE", lkNone, BtPathname)
	usr("PAM_DEFAULT", lkNone, enum("linux-pam openpam solaris-pam"))
	acl("PATCHDIR", lkNone, BtRelativePkgPath, "Makefile: set; Makefile.common: default, set")
	pkglist("PATCHFILES", lkShell, BtFilename)
	acl("PATCH_ARGS", lkShell, BtShellWord, "")
	acl("PATCH_DIST_ARGS", lkShell, BtShellWord, "Makefile: set, append")
	acl("PATCH_DIST_CAT", lkNone, BtShellCommand, "")
	acl("PATCH_DIST_STRIP*", lkNone, BtShellWord, "buildlink3.mk, builtin.mk:; Makefile, Makefile.common, *.mk: set")
	acl("PATCH_SITES", lkShell, BtFetchURL, "Makefile, Makefile.common, options.mk: set")
	acl("PATCH_STRIP", lkNone, BtShellWord, "")
	acl("PERL5_PACKLIST", lkShell, BtPerl5Packlist, "Makefile: set; options.mk: set, append")
	acl("PERL5_PACKLIST_DIR", lkNone, BtPathname, "")
	pkg("PERL5_REQD", lkShell, BtVersion)
	pkg("PERL5_USE_PACKLIST", lkNone, BtYesNo)
	sys("PGSQL_PREFIX", lkNone, BtPathname)
	acl("PGSQL_VERSIONS_ACCEPTED", lkShell, pgsqlVersions, "")
	usr("PGSQL_VERSION_DEFAULT", lkNone, BtVersion)
	sys("PG_LIB_EXT", lkNone, enum("dylib so"))
	sys("PGSQL_TYPE", lkNone, enum("postgresql81-client postgresql80-client"))
	sys("PGPKGSRCDIR", lkNone, BtPathname)
	sys("PHASE_MSG", lkNone, BtShellCommand)
	usr("PHP_VERSION_REQD", lkNone, BtVersion)
	sys("PKGBASE", lkNone, BtIdentifier)
	acl("PKGCONFIG_FILE.*", lkShell, BtPathname, "builtin.mk: set, append; pkgconfig-builtin.mk: use-loadtime")
	acl("PKGCONFIG_OVERRIDE", lkShell, BtPathmask, "Makefile: set, append; Makefile.common: append")
	pkg("PKGCONFIG_OVERRIDE_STAGE", lkNone, BtStage)
	pkg("PKGDIR", lkNone, BtRelativePkgDir)
	sys("PKGDIRMODE", lkNone, BtFileMode)
	sys("PKGLOCALEDIR", lkNone, BtPathname)
	pkg("PKGNAME", lkNone, BtPkgName)
	sys("PKGNAME_NOREV", lkNone, BtPkgName)
	sys("PKGPATH", lkNone, BtPathname)
	acl("PKGREPOSITORY", lkNone, BtUnknown, "")
	acl("PKGREVISION", lkNone, BtPkgRevision, "Makefile: set")
	sys("PKGSRCDIR", lkNone, BtPathname)
	acl("PKGSRCTOP", lkNone, BtYes, "Makefile: set")
	acl("PKGTOOLS_ENV", lkShell, BtShellWord, "")
	sys("PKGVERSION", lkNone, BtVersion)
	sys("PKGWILDCARD", lkNone, BtFilemask)
	sys("PKG_ADMIN", lkNone, BtShellCommand)
	sys("PKG_APACHE", lkNone, enum("apache24"))
	pkg("PKG_APACHE_ACCEPTED", lkShell, enum("apache24"))
	usr("PKG_APACHE_DEFAULT", lkNone, enum("apache24"))
	usr("PKG_CONFIG", lkNone, BtYes)
	// ^^ No, this is not the popular command from GNOME, but the setting
	// whether the pkgsrc user wants configuration files automatically
	// installed or not.
	sys("PKG_CREATE", lkNone, BtShellCommand)
	sys("PKG_DBDIR", lkNone, BtPathname)
	cmdline("PKG_DEBUG_LEVEL", lkNone, BtInteger)
	usr("PKG_DEFAULT_OPTIONS", lkShell, BtOption)
	sys("PKG_DELETE", lkNone, BtShellCommand)
	acl("PKG_DESTDIR_SUPPORT", lkShell, enum("destdir user-destdir"), "Makefile, Makefile.common: set")
	pkglist("PKG_FAIL_REASON", lkShell, BtShellWord)
	acl("PKG_GECOS.*", lkNone, BtMessage, "Makefile: set")
	acl("PKG_GID.*", lkNone, BtInteger, "Makefile: set")
	acl("PKG_GROUPS", lkShell, BtShellWord, "Makefile: set, append")
	pkglist("PKG_GROUPS_VARS", lkShell, BtVariableName)
	acl("PKG_HOME.*", lkNone, BtPathname, "Makefile: set")
	acl("PKG_HACKS", lkShell, BtIdentifier, "hacks.mk: append")
	sys("PKG_INFO", lkNone, BtShellCommand)
	sys("PKG_JAVA_HOME", lkNone, BtPathname)
	sys("PKG_JVM", lkNone, jvms)
	acl("PKG_JVMS_ACCEPTED", lkShell, jvms, "Makefile: set; Makefile.common: default, set")
	usr("PKG_JVM_DEFAULT", lkNone, jvms)
	acl("PKG_LEGACY_OPTIONS", lkShell, BtOption, "")
	acl("PKG_LIBTOOL", lkNone, BtPathname, "Makefile: set")
	acl("PKG_OPTIONS", lkShell, BtOption, "bsd.options.mk: set; *: use-loadtime, use")
	usr("PKG_OPTIONS.*", lkShell, BtOption)
	acl("PKG_OPTIONS_DEPRECATED_WARNINGS", lkShell, BtShellWord, "")
	acl("PKG_OPTIONS_GROUP.*", lkShell, BtOption, "Makefile, options.mk: set, append")
	acl("PKG_OPTIONS_LEGACY_OPTS", lkShell, BtUnknown, "Makefile, Makefile.common, options.mk: append")
	acl("PKG_OPTIONS_LEGACY_VARS", lkShell, BtUnknown, "Makefile, Makefile.common, options.mk: append")
	acl("PKG_OPTIONS_NONEMPTY_SETS", lkShell, BtIdentifier, "")
	acl("PKG_OPTIONS_OPTIONAL_GROUPS", lkShell, BtIdentifier, "options.mk: set, append")
	acl("PKG_OPTIONS_REQUIRED_GROUPS", lkShell, BtIdentifier, "Makefile, options.mk: set")
	acl("PKG_OPTIONS_SET.*", lkShell, BtOption, "")
	acl("PKG_OPTIONS_VAR", lkNone, BtPkgOptionsVar, "Makefile, Makefile.common, options.mk: set; bsd.options.mk: use-loadtime")
	acl("PKG_PRESERVE", lkNone, BtYes, "Makefile: set")
	acl("PKG_SHELL", lkNone, BtPathname, "Makefile, Makefile.common: set")
	acl("PKG_SHELL.*", lkNone, BtPathname, "Makefile, Makefile.common: set")
	acl("PKG_SHLIBTOOL", lkNone, BtPathname, "")
	pkglist("PKG_SKIP_REASON", lkShell, BtShellWord)
	acl("PKG_SUGGESTED_OPTIONS", lkShell, BtOption, "Makefile, Makefile.common, options.mk: set, append")
	acl("PKG_SUGGESTED_OPTIONS.*", lkShell, BtOption, "Makefile, Makefile.common, options.mk: set, append")
	acl("PKG_SUPPORTED_OPTIONS", lkShell, BtOption, "Makefile: set, append; Makefile.common: set; options.mk: set, append, use")
	pkg("PKG_SYSCONFDIR*", lkNone, BtPathname)
	pkglist("PKG_SYSCONFDIR_PERMS", lkShell, BtPerms)
	sys("PKG_SYSCONFBASEDIR", lkNone, BtPathname)
	pkg("PKG_SYSCONFSUBDIR", lkNone, BtPathname)
	acl("PKG_SYSCONFVAR", lkNone, BtIdentifier, "") // FIXME: name/type mismatch.
	acl("PKG_UID", lkNone, BtInteger, "Makefile: set")
	acl("PKG_USERS", lkShell, BtShellWord, "Makefile: set, append")
	pkg("PKG_USERS_VARS", lkShell, BtVariableName)
	acl("PKG_USE_KERBEROS", lkNone, BtYes, "Makefile, Makefile.common: set")
	pkg("PLIST.*", lkNone, BtYes)
	pkglist("PLIST_VARS", lkShell, BtIdentifier)
	pkglist("PLIST_SRC", lkShell, BtRelativePkgPath)
	pkglist("PLIST_SUBST", lkShell, BtShellWord)
	acl("PLIST_TYPE", lkNone, enum("dynamic static"), "")
	acl("PREPEND_PATH", lkShell, BtPathname, "")
	acl("PREFIX", lkNone, BtPathname, "*: use")
	acl("PREV_PKGPATH", lkNone, BtPathname, "*: use") // doesn't exist any longer
	acl("PRINT_PLIST_AWK", lkNone, BtAwkCommand, "*: append")
	acl("PRIVILEGED_STAGES", lkShell, enum("install package clean"), "")
	acl("PTHREAD_AUTO_VARS", lkNone, BtYesNo, "Makefile: set")
	sys("PTHREAD_CFLAGS", lkShell, BtCFlag)
	sys("PTHREAD_LDFLAGS", lkShell, BtLdFlag)
	sys("PTHREAD_LIBS", lkShell, BtLdFlag)
	acl("PTHREAD_OPTS", lkShell, enum("native optional require"), "Makefile: set, append; Makefile.common, buildlink3.mk: append")
	sys("PTHREAD_TYPE", lkNone, BtIdentifier) // Or "native" or "none".
	pkg("PY_PATCHPLIST", lkNone, BtYes)
	acl("PYPKGPREFIX", lkNone, enum("py27 py34 py35 py36"), "pyversion.mk: set; *: use-loadtime, use")
	pkg("PYTHON_FOR_BUILD_ONLY", lkNone, BtYes)
	pkglist("REPLACE_PYTHON", lkShell, BtPathmask)
	pkg("PYTHON_VERSIONS_ACCEPTED", lkShell, BtVersion)
	pkg("PYTHON_VERSIONS_INCOMPATIBLE", lkShell, BtVersion)
	usr("PYTHON_VERSION_DEFAULT", lkNone, BtVersion)
	usr("PYTHON_VERSION_REQD", lkNone, BtVersion)
	pkglist("PYTHON_VERSIONED_DEPENDENCIES", lkShell, BtPythonDependency)
	sys("RANLIB", lkNone, BtShellCommand)
	pkglist("RCD_SCRIPTS", lkShell, BtFilename)
	acl("RCD_SCRIPT_SRC.*", lkNone, BtPathname, "Makefile: set")
	acl("RCD_SCRIPT_WRK.*", lkNone, BtPathname, "Makefile: set")
	acl("REPLACE.*", lkNone, BtUnknown, "Makefile: set")
	pkglist("REPLACE_AWK", lkShell, BtPathmask)
	pkglist("REPLACE_BASH", lkShell, BtPathmask)
	pkglist("REPLACE_CSH", lkShell, BtPathmask)
	acl("REPLACE_EMACS", lkShell, BtPathmask, "")
	acl("REPLACE_FILES.*", lkShell, BtPathmask, "Makefile, Makefile.common: set, append")
	acl("REPLACE_INTERPRETER", lkShell, BtIdentifier, "Makefile, Makefile.common: append")
	pkglist("REPLACE_KSH", lkShell, BtPathmask)
	pkglist("REPLACE_LOCALEDIR_PATTERNS", lkShell, BtFilemask)
	pkglist("REPLACE_LUA", lkShell, BtPathmask)
	pkglist("REPLACE_PERL", lkShell, BtPathmask)
	pkglist("REPLACE_PYTHON", lkShell, BtPathmask)
	pkglist("REPLACE_SH", lkShell, BtPathmask)
	pkglist("REQD_DIRS", lkShell, BtPathname)
	pkglist("REQD_DIRS_PERMS", lkShell, BtPerms)
	pkglist("REQD_FILES", lkShell, BtPathname)
	pkg("REQD_FILES_MODE", lkNone, enum("0644 0640 0600 0400"))
	pkglist("REQD_FILES_PERMS", lkShell, BtPerms)
	pkg("RESTRICTED", lkNone, BtMessage)
	usr("ROOT_USER", lkNone, BtUserGroupName)
	usr("ROOT_GROUP", lkNone, BtUserGroupName)
	usr("RUBY_VERSION_REQD", lkNone, BtVersion)
	sys("RUN", lkNone, BtShellCommand)
	sys("RUN_LDCONFIG", lkNone, BtYesNo)
	acl("SCRIPTS_ENV", lkShell, BtShellWord, "Makefile, Makefile.common: append")
	usr("SETUID_ROOT_PERMS", lkShell, BtShellWord)
	pkg("SET_LIBDIR", lkNone, BtYes)
	sys("SHAREGRP", lkNone, BtUserGroupName)
	sys("SHAREMODE", lkNone, BtFileMode)
	sys("SHAREOWN", lkNone, BtUserGroupName)
	sys("SHCOMMENT", lkNone, BtShellCommand)
	acl("SHLIB_HANDLING", lkNone, enum("YES NO no"), "")
	acl("SHLIBTOOL", lkNone, BtShellCommand, "Makefile: use")
	acl("SHLIBTOOL_OVERRIDE", lkShell, BtPathmask, "Makefile: set, append; Makefile.common: append")
	acl("SITES.*", lkShell, BtFetchURL, "Makefile, Makefile.common, options.mk: set, append, use")
	usr("SMF_PREFIS", lkNone, BtPathname)
	pkg("SMF_SRCDIR", lkNone, BtPathname)
	pkg("SMF_NAME", lkNone, BtFilename)
	pkg("SMF_MANIFEST", lkNone, BtPathname)
	pkg("SMF_INSTANCES", lkShell, BtIdentifier)
	pkg("SMF_METHODS", lkShell, BtFilename)
	pkg("SMF_METHOD_SRC.*", lkNone, BtPathname)
	pkg("SMF_METHOD_SHELL", lkNone, BtShellCommand)
	pkglist("SPECIAL_PERMS", lkShell, BtPerms)
	sys("STEP_MSG", lkNone, BtShellCommand)
	sys("STRIP", lkNone, BtShellCommand) // see mk/tools/strip.mk
	acl("SUBDIR", lkShell, BtFilename, "Makefile: append; *:")
	acl("SUBST_CLASSES", lkShell, BtIdentifier, "Makefile: set, append; *: append")
	acl("SUBST_CLASSES.*", lkShell, BtIdentifier, "Makefile: set, append; *: append")
	acl("SUBST_FILES.*", lkShell, BtPathmask, "Makefile, Makefile.*, *.mk: set, append")
	acl("SUBST_FILTER_CMD.*", lkNone, BtShellCommand, "Makefile, Makefile.*, *.mk: set")
	acl("SUBST_MESSAGE.*", lkNone, BtMessage, "Makefile, Makefile.*, *.mk: set")
	acl("SUBST_SED.*", lkNone, BtSedCommands, "Makefile, Makefile.*, *.mk: set, append")
	pkg("SUBST_STAGE.*", lkNone, BtStage)
	pkglist("SUBST_VARS.*", lkShell, BtVariableName)
	pkglist("SUPERSEDES", lkSpace, BtDependency)
	acl("TEST_DEPENDS", lkSpace, BtDependencyWithPath, "Makefile, Makefile.common, *.mk: append")
	pkglist("TEST_DIRS", lkShell, BtWrksrcSubdirectory)
	pkglist("TEST_ENV", lkShell, BtShellWord)
	acl("TEST_TARGET", lkShell, BtIdentifier, "Makefile: set; Makefile.common: default, set; options.mk: set, append")
	pkglist("TEXINFO_REQD", lkShell, BtVersion)
	acl("TOOL_DEPENDS", lkSpace, BtDependencyWithPath, "Makefile, Makefile.common, *.mk: append")
	sys("TOOLS_ALIASES", lkShell, BtFilename)
	sys("TOOLS_BROKEN", lkShell, BtTool)
	sys("TOOLS_CMD.*", lkNone, BtPathname)
	acl("TOOLS_CREATE", lkShell, BtTool, "Makefile, Makefile.common, options.mk: append")
	acl("TOOLS_DEPENDS.*", lkSpace, BtDependencyWithPath, "buildlink3.mk:; Makefile, Makefile.*: set, default; *: use")
	sys("TOOLS_GNU_MISSING", lkShell, BtTool)
	sys("TOOLS_NOOP", lkShell, BtTool)
	sys("TOOLS_PATH.*", lkNone, BtPathname)
	sys("TOOLS_PLATFORM.*", lkNone, BtShellCommand)
	sys("TOUCH_FLAGS", lkShell, BtShellWord)
	pkglist("UAC_REQD_EXECS", lkShell, BtPrefixPathname)
	acl("UNLIMIT_RESOURCES", lkShell, enum("cputime datasize memorysize stacksize"), "Makefile: set, append; Makefile.common: append")
	usr("UNPRIVILEGED_USER", lkNone, BtUserGroupName)
	usr("UNPRIVILEGED_GROUP", lkNone, BtUserGroupName)
	pkglist("UNWRAP_FILES", lkShell, BtPathmask)
	usr("UPDATE_TARGET", lkShell, BtIdentifier)
	pkg("USERGROUP_PHASE", lkNone, enum("configure build pre-install"))
	pkg("USE_BSD_MAKEFILE", lkNone, BtYes)
	acl("USE_BUILTIN.*", lkNone, BtYesNoIndirectly, "builtin.mk: set")
	pkg("USE_CMAKE", lkNone, BtYes)
	usr("USE_DESTDIR", lkNone, BtYes)
	pkglist("USE_FEATURES", lkShell, BtIdentifier)
	acl("USE_GAMESGROUP", lkNone, BtYesNo, "buildlink3.mk, builtin.mk:; *: set, default, use")
	pkg("USE_GCC_RUNTIME", lkNone, BtYesNo)
	pkg("USE_GNU_CONFIGURE_HOST", lkNone, BtYesNo)
	acl("USE_GNU_ICONV", lkNone, BtYes, "Makefile, Makefile.common, options.mk: set")
	acl("USE_IMAKE", lkNone, BtYes, "Makefile: set")
	pkg("USE_JAVA", lkNone, enum("run yes build"))
	pkg("USE_JAVA2", lkNone, enum("YES yes no 1.4 1.5 6 7 8"))
	acl("USE_LANGUAGES", lkShell, languages, "Makefile, Makefile.common, options.mk: set, append")
	pkg("USE_LIBTOOL", lkNone, BtYes)
	pkg("USE_MAKEINFO", lkNone, BtYes)
	pkg("USE_MSGFMT_PLURALS", lkNone, BtYes)
	pkg("USE_NCURSES", lkNone, BtYes)
	pkg("USE_OLD_DES_API", lkNone, BtYesNo)
	pkg("USE_PKGINSTALL", lkNone, BtYes)
	pkg("USE_PKGLOCALEDIR", lkNone, BtYesNo)
	usr("USE_PKGSRC_GCC", lkNone, BtYes)
	acl("USE_TOOLS", lkShell, BtTool, "*: append")
	acl("USE_TOOLS.*", lkShell, BtTool, "*: append")
	pkg("USE_X11", lkNone, BtYes)
	sys("WARNINGS", lkShell, BtShellWord)
	sys("WARNING_MSG", lkNone, BtShellCommand)
	sys("WARNING_CAT", lkNone, BtShellCommand)
	acl("WRAPPER_REORDER_CMDS", lkShell, BtWrapperReorder, "Makefile, Makefile.common, buildlink3.mk: append")
	pkg("WRAPPER_SHELL", lkNone, BtShellCommand)
	acl("WRAPPER_TRANSFORM_CMDS", lkShell, BtWrapperTransform, "Makefile, Makefile.common, buildlink3.mk: append")
	sys("WRKDIR", lkNone, BtPathname)
	pkg("WRKSRC", lkNone, BtWrkdirSubdirectory)
	sys("X11_PKGSRCDIR.*", lkNone, BtPathname)
	usr("XAW_TYPE", lkNone, enum("3d neXtaw standard xpm"))
	acl("XMKMF_FLAGS", lkShell, BtShellWord, "")
	pkglist("_WRAP_EXTRA_ARGS.*", lkShell, BtShellWord)

	// Only for infrastructure files; see mk/misc/show.mk
	acl("_VARGROUPS", lkSpace, BtIdentifier, "*: append")
	acl("_USER_VARS.*", lkSpace, BtIdentifier, "*: append")
	acl("_PKG_VARS.*", lkSpace, BtIdentifier, "*: append")
	acl("_SYS_VARS.*", lkSpace, BtIdentifier, "*: append")
	acl("_DEF_VARS.*", lkSpace, BtIdentifier, "*: append")
	acl("_USE_VARS.*", lkSpace, BtIdentifier, "*: append")
}

func enum(values string) *BasicType {
	vmap := make(map[string]bool)
	for _, value := range splitOnSpace(values) {
		vmap[value] = true
	}
	name := "enum: " + values + " " // See IsEnum
	return &BasicType{name, func(cv *VartypeCheck) {
		if cv.Op == opUseMatch {
			if !vmap[cv.Value] && cv.Value == cv.ValueNoVar {
				canMatch := false
				for value := range vmap {
					if ok, err := path.Match(cv.Value, value); err != nil {
						cv.Line.Warnf("Invalid match pattern %q.", cv.Value)
					} else if ok {
						canMatch = true
					}
				}
				if !canMatch {
					cv.Line.Warnf("The pattern %q cannot match any of { %s } for %s.", cv.Value, values, cv.Varname)
				}
			}
			return
		}

		if cv.Value == cv.ValueNoVar && !vmap[cv.Value] {
			cv.Line.Warnf("%q is not valid for %s. Use one of { %s } instead.", cv.Value, cv.Varname, values)
		}
	}}
}

func parseACLEntries(varname string, aclentries string) []ACLEntry {
	if aclentries == "" {
		return nil
	}
	var result []ACLEntry
	prevperms := "(first)"
	for _, arg := range strings.Split(aclentries, "; ") {
		var globs, perms string
		if fields := strings.SplitN(arg, ": ", 2); len(fields) == 2 {
			globs, perms = fields[0], fields[1]
		} else {
			globs = strings.TrimSuffix(arg, ":")
		}
		if perms == prevperms {
			fmt.Printf("Repeated permissions for %s: %s\n", varname, perms)
		}
		prevperms = perms
		var permissions ACLPermissions
		for _, perm := range strings.Split(perms, ", ") {
			switch perm {
			case "append":
				permissions |= aclpAppend
			case "default":
				permissions |= aclpSetDefault
			case "set":
				permissions |= aclpSet
			case "use":
				permissions |= aclpUse
			case "use-loadtime":
				permissions |= aclpUseLoadtime
			case "":
				break
			default:
				print(fmt.Sprintf("Invalid ACL permission %q for varname %q.\n", perm, varname))
			}
		}
		for _, glob := range strings.Split(globs, ", ") {
			switch glob {
			case "*",
				"Makefile", "Makefile.common", "Makefile.*",
				"buildlink3.mk", "builtin.mk", "options.mk", "hacks.mk", "*.mk",
				"bsd.options.mk", "pkgconfig-builtin.mk", "pyversion.mk":
				break
			default:
				print(fmt.Sprintf("Invalid ACL glob %q for varname %q.\n", glob, varname))
			}
			for _, prev := range result {
				if matched, err := path.Match(prev.glob, glob); err != nil || matched {
					print(fmt.Sprintf("Ineffective ACL glob %q for varname %q.\n", glob, varname))
				}
			}
			result = append(result, ACLEntry{glob, permissions})
		}
	}
	return result
}
