# $NetBSD: buildlink.mk,v 1.6 2001/06/23 19:26:57 jlam Exp $
#
# This Makefile fragment is included by packages that use libpng.
#
# To use this Makefile fragment, simply:
#
# (1) Optionally define BUILDLINK_DEPENDS.png to the dependency pattern
#     for the version of libpng desired.
# (2) Include this Makefile fragment in the package Makefile,
# (3) Add ${BUILDLINK_DIR}/include to the front of the C preprocessor's header
#     search path, and
# (4) Add ${BUILDLINK_DIR}/lib to the front of the linker's library search
#     path.

.if !defined(PNG_BUILDLINK_MK)
PNG_BUILDLINK_MK=	# defined

BUILDLINK_DEPENDS.png?=	png>=1.0.11
DEPENDS+=		${BUILDLINK_DEPENDS.png}:../../graphics/png

BUILDLINK_PREFIX.png=	${LOCALBASE}
BUILDLINK_FILES.png=	include/png.h
BUILDLINK_FILES.png+=	include/pngconf.h
BUILDLINK_FILES.png+=	lib/libpng.*

.include "../../devel/zlib/buildlink.mk"

BUILDLINK_TARGETS.png=	png-buildlink
BUILDLINK_TARGETS+=	${BUILDLINK_TARGETS.png}

pre-configure: ${BUILDLINK_TARGETS.png}
png-buildlink: _BUILDLINK_USE

.include "../../mk/bsd.buildlink.mk"

.endif	# PNG_BUILDLINK_MK
