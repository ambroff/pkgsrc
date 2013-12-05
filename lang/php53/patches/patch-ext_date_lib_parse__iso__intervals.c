$NetBSD: patch-ext_date_lib_parse__iso__intervals.c,v 1.1 2013/12/05 16:16:40 taca Exp $

Fix for CVE-2013-6712.

--- ext/date/lib/parse_iso_intervals.c.orig	2013-07-10 17:43:08.000000000 +0000
+++ ext/date/lib/parse_iso_intervals.c
@@ -415,7 +415,7 @@ yy6:
 					break;
 			}
 			ptr++;
-		} while (*ptr);
+		} while (!s->errors->error_count && *ptr);
 		s->have_period = 1;
 		TIMELIB_DEINIT;
 		return TIMELIB_PERIOD;
