# $NetBSD: buildlink.mk,v 1.2 2001/06/23 19:27:00 jlam Exp $
#
# This Makefile fragment is included by packages that use tcp_wrappers.
#
# To use this Makefile fragment, simply:
#
# (1) Optionally define BUILDLINK_DEPENDS.tcp_wrappers to the dependency
#     pattern for the version of tcp_wrappers desired.
# (2) Include this Makefile fragment in the package Makefile,
# (3) Add ${BUILDLINK_DIR}/include to the front of the C preprocessor's header
#     search path, and
# (4) Add ${BUILDLINK_DIR}/lib to the front of the linker's library search
#     path.

.if !defined(TCP_WRAPPERS_BUILDLINK_MK)
TCP_WRAPPERS_BUILDLINK_MK=	# defined

BUILDLINK_DEPENDS.tcp_wrappers?=	tcp_wrappers>=7.6.1nb1

.if exists(/usr/include/tcpd.h)
_NEED_TCP_WRAPPERS=	NO
.else
_NEED_TCP_WRAPPERS=	YES
.endif

.if ${_NEED_TCP_WRAPPERS} == "YES"
DEPENDS+=	${BUILDLINK_DEPENDS.tcp_wrappers}:../../security/tcp_wrappers
BUILDLINK_PREFIX.tcp_wrappers=	${LOCALBASE}
.else
BUILDLINK_PREFIX.tcp_wrappers=	/usr
.endif

BUILDLINK_FILES.tcp_wrappers=	include/tcpd.h
BUILDLINK_FILES.tcp_wrappers+=	lib/libwrap.*

BUILDLINK_TARGETS.tcp_wrappers=	tcp_wrappers-buildlink
BUILDLINK_TARGETS+=		${BUILDLINK_TARGETS.tcp_wrappers}

pre-configure: ${BUILDLINK_TARGETS.tcp_wrappers}
tcp_wrappers-buildlink: _BUILDLINK_USE

.include "../../mk/bsd.buildlink.mk"

.endif	# TCP_WRAPPERS_BUILDLINK_MK
