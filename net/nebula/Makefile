# $NetBSD$

DISTNAME=	nebula-1.7.2
CATEGORIES=	net
GITHUB_TAG=	refs/tags/v${PKGVERSION_NOREV}
MASTER_SITES=	${MASTER_SITE_GITHUB:=slackhq/}

MAINTAINER=	kyle@ambroffkao.com
HOMEPAGE=	https://github.com/slackhq/nebula/
COMMENT=	Scalable overlay networking tool.
LICENSE=	mit

WRKSRC=		${WRKDIR}/nebula-${PKGVERSION_NOREV}

USE_LANGUAGES=	# go
USE_TOOLS+=	gmake

#.include "go-modules.mk"
.include "../../lang/go/go-module.mk"
.include "../../mk/bsd.pkg.mk"
