$NetBSD: patch-src_ical-code.c,v 1.1 2018/04/25 21:44:44 markd Exp $

Work with libical-3

--- src/ical-code.c.orig	2015-04-10 10:26:26.000000000 +0000
+++ src/ical-code.c
@@ -129,7 +129,6 @@ static struct icaltimetype ical_get_curr
         &&   (strcmp(g_par.local_timezone, "floating") != 0))
         ctime = icaltime_current_time_with_zone(local_icaltimezone);
     else { / * use floating time * /
-        ctime.is_utc      = 0;
         ctime.is_date     = 0;
         ctime.is_daylight = 0;
         ctime.zone        = NULL;
@@ -2579,7 +2578,6 @@ static struct icaltimetype count_first_a
  * when counting alarm time. */
         if (rel == ICAL_RELATED_START) {
             per.stime.is_date       = 0;
-            per.stime.is_utc        = 1;
             per.stime.is_daylight   = 0;
             per.stime.zone          = utc_icaltimezone;
             per.stime.hour          = 0;
@@ -2588,7 +2586,6 @@ static struct icaltimetype count_first_a
         }
         else {
             per.etime.is_date       = 0;
-            per.etime.is_utc        = 1;
             per.etime.is_daylight   = 0;
             per.etime.zone          = utc_icaltimezone;
             per.etime.hour          = 0;
@@ -2613,7 +2610,6 @@ static struct icaltimetype count_next_al
 /* HACK: convert to UTC time so that we can use time arithmetic
  * when counting alarm time. */
         start_time.is_date       = 0;
-        start_time.is_utc        = 1;
         start_time.is_daylight   = 0;
         start_time.zone          = utc_icaltimezone;
         start_time.hour          = 0;
@@ -2768,7 +2764,6 @@ static alarm_struct *process_alarm_trigg
      */
     if (icaltime_is_date(per.stime)) {
         if (local_icaltimezone != utc_icaltimezone) {
-            next_alarm_time.is_utc        = 0;
             next_alarm_time.is_daylight   = 0;
             next_alarm_time.zone          = local_icaltimezone;
         }
@@ -2850,7 +2845,6 @@ orage_message(120, P_N "Alarm rec loop n
          */
         if (icaltime_is_date(per.stime)) {
             if (local_icaltimezone != utc_icaltimezone) {
-                next_alarm_time.is_utc        = 0;
                 next_alarm_time.is_daylight   = 0;
                 next_alarm_time.zone          = local_icaltimezone;
             }
@@ -2944,7 +2938,6 @@ orage_message(120, P_N "*****After loop 
          */
         if (icaltime_is_date(per.stime)) {
             if (local_icaltimezone != utc_icaltimezone) {
-                next_alarm_time.is_utc        = 0;
                 next_alarm_time.is_daylight   = 0;
                 next_alarm_time.zone          = local_icaltimezone;
             }
