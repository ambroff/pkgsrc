# $NetBSD: buildlink.mk,v 1.2 2001/06/23 19:27:00 jlam Exp $
#
# This Makefile fragment is included by packages that use libcrack.
#
# To use this Makefile fragment, simply:
#
# (1) Optionally define BUILDLINK_DEPENDS.libcrack to the dependency pattern
#     for the version of libcrack desired.
# (2) Include this Makefile fragment in the package Makefile,
# (3) Add ${BUILDLINK_DIR}/include to the front of the C preprocessor's header
#     search path, and
# (4) Add ${BUILDLINK_DIR}/lib to the front of the linker's library search
#     path.

.if !defined(LIBCRACK_BUILDLINK_MK)
LIBCRACK_BUILDLINK_MK=	# defined

BUILDLINK_DEPENDS.libcrack?=	libcrack>=2.7
DEPENDS+=	${BUILDLINK_DEPENDS.libcrack}:../../security/libcrack

BUILDLINK_PREFIX.libcrack=	${LOCALBASE}
BUILDLINK_FILES.libcrack=	include/cracklib/*
BUILDLINK_FILES.libcrack+=	lib/libcrack.*

BUILDLINK_TARGETS.libcrack=	libcrack-buildlink
BUILDLINK_TARGETS+=		${BUILDLINK_TARGETS.libcrack}

pre-configure: ${BUILDLINK_TARGETS.libcrack}
libcrack-buildlink: _BUILDLINK_USE

.include "../../mk/bsd.buildlink.mk"

.endif	# LIBCRACK_BUILDLINK_MK
