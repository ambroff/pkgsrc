# $NetBSD: options.mk,v 1.1 2019/02/01 16:30:00 mgorny Exp $

PKG_OPTIONS_VAR=	PKG_OPTIONS.lld
PKG_SUPPORTED_OPTIONS=	tests

.include "../../mk/bsd.options.mk"

.if !empty(PKG_OPTIONS:Mtests)
DISTFILES+=		llvm-${PKGVERSION_NOREV}.src${EXTRACT_SUFX}
CMAKE_ARGS+=		-DLLVM_CONFIG_PATH=${LLVM_CONFIG_PATH:Q}
CMAKE_ARGS+=		-DLLVM_INCLUDE_TESTS=ON
CMAKE_ARGS+=		-DLLVM_BUILD_TESTS=ON
CMAKE_ARGS+=		-DLLVM_MAIN_SRC_DIR=${WRKDIR}/llvm-${PKGVERSION_NOREV}.src
CMAKE_ARGS+=		-DLLVM_EXTERNAL_LIT=${WRKDIR}/llvm-${PKGVERSION_NOREV}.src/utils/lit/lit.py
REPLACE_PYTHON+=	${WRKDIR}/llvm-${PKGVERSION_NOREV}.src/utils/lit/lit.py
TEST_TARGET=		check-lld  # failing tests fixed in 8.0
TEST_ENV+=		LD_LIBRARY_PATH=${WRKDIR}/build/lib
.else
CMAKE_ARGS+=		-DLLVM_INCLUDE_TESTS=OFF
.endif
