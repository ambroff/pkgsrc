/*	$NetBSD: strndup.c,v 1.1 2023/11/09 18:55:18 nia Exp $	*/

/*
 * Copyright (c) 1988, 1993
 *	The Regents of the University of California.  All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 * 3. Neither the name of the University nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */

#if HAVE_NBTOOL_CONFIG_H
#include "nbtool_config.h"
#endif

#include <nbcompat.h>
#include <nbcompat/cdefs.h>

#include <sys/cdefs.h>
#if defined(LIBC_SCCS) && !defined(lint)
#if 0
static char sccsid[] = "@(#)strdup.c	8.1 (Berkeley) 6/4/93";
#else
__RCSID("$NetBSD: strndup.c,v 1.1 2023/11/09 18:55:18 nia Exp $");
#endif
#endif /* LIBC_SCCS and not lint */

#if 0
#include "namespace.h"
#endif

#include <nbcompat/assert.h>
#if HAVE_ERRNO_H
#include <errno.h>
#endif
#include <nbcompat/stdlib.h>
#include <nbcompat/string.h>

#if 0
#ifdef __weak_alias
__weak_alias(strndup,_strndup)
#endif
#endif

#if !HAVE_STRNDUP
char *
strndup(const char *str, size_t n)
{
	size_t len;
	char *copy;

	_DIAGASSERT(str != NULL);

	for (len = 0; len < n && str[len]; len++)
		continue;

	if (!(copy = malloc(len + 1)))
		return (NULL);
	memcpy(copy, str, len);
	copy[len] = '\0';
	return (copy);
}
#endif /* !HAVE_STRNDUP */
