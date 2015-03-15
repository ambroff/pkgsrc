# $NetBSD: common.mk,v 1.3 2015/03/15 15:25:09 taca Exp $
# used by lang/php54/Makefile.common
# used by lang/php55/Makefile.common
# 
.if !defined(PHPCOMMON_MK)
PHPCOMMON_MK=	defined

DISTNAME?=		php-${PHP_BASE_VERS}

.if !defined(PECL_VERSION)
MASTER_SITES?=		http://www.php.net/distributions/ \
			http://uk.php.net/distributions/ \
			http://de.php.net/distributions/ \
			http://us.php.net/distributions/
EXTRACT_SUFX?=		.tar.bz2
.endif

HOMEPAGE?=		http://www.php.net/

PLIST_SUBST+=		PHP_EXTENSION_DIR=${PHP_EXTENSION_DIR}
PKG_SYSCONFVAR?=	php

.include "../../lang/php/phpversion.mk"
.endif	# PHPCOMMON_MK
