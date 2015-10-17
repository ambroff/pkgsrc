# $NetBSD: hacks.mk,v 1.1 2015/10/17 10:28:43 fhajny Exp $

.if !defined(OPENCV_HACKS_MK)
OPENCV_HACKS_MK=	defined

.include "../../mk/bsd.fast.prefs.mk"
.include "../../mk/compiler.mk"

# PR toolchain/47051: gcc-4.5.4 breaks opencv on amd64
.if !empty(PKGSRC_COMPILER:Mgcc) && !empty(CC_VERSION:Mgcc-4.5.4*) && !empty(MACHINE_PLATFORM:M*-*-x86_64)
PKG_HACKS+=		tree-pre
SUBST_CLASSES+=		opt-hack
SUBST_STAGE.opt-hack=	post-configure
SUBST_MESSAGE.opt-hack=	Working around gcc-4.5.4 bug.
SUBST_FILES.opt-hack=	${WRKSRC}/modules/calib3d/CMakeFiles/opencv_calib3d.dir/build.make
SUBST_SED.opt-hack=	-e '/stereosgbm.cpp.o/s/-o /-fno-tree-pre -o /'
.endif

.endif	# OPENCV_HACKS_MK
