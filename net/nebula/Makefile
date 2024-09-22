DISTNAME=	nebula-1.9.4
PKGREVISION=	0
CATEGORIES=	net
MASTER_SITES=	${MASTER_SITE_GITHUB:=slackhq/}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	https://nebula.defined.net
COMMENT=	Scalable overlay networking tool
LICENSE=	mit

GITHUB_TAG=	v${PKGVERSION_NOREV}

GO_DIST_BASE=	${DISTNAME}
GO_SRCPATH=	github.com/slackhq/nebula
GO_BUILD_PATTERN=	${GO_SRCPATH}/

USE_TOOLS+=	bash gmake

INSTALLATION_DIRS+=	bin

DEPENDS+=	go>=1.22:../../lang/go

do-install:
	${INSTALL_PROGRAM} ${WRKDIR}/nebula ${DESTDIR}${PREFIX}/bin
	${INSTALL_PROGRAM} ${WRKDIR}/nebula-cert ${DESTDIR}${PREFIX}/bin

.include "../../lang/go/go-package.mk"
.include "../../mk/bsd.pkg.mk"
