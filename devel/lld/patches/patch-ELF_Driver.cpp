$NetBSD: patch-ELF_Driver.cpp,v 1.1 2019/02/01 16:30:00 mgorny Exp $

Add support for customizing LLD behavior on target triple.
https://reviews.llvm.org/D56650

Add '-z nognustack' option to disable emitting PT_GNU_STACK.
https://reviews.llvm.org/D56554

Alter defaults for NetBSD targets:
* add default library search paths
* force combined RO+RW segment due to ld.elf_so limitations
* disable PT_GNU_STACK (meaningless on NetBSD)
* disable 'new dtags', i.e. force RPATH instead of RUNPATH

--- ELF/Driver.cpp.orig	2018-07-31 21:58:26.000000000 +0000
+++ ELF/Driver.cpp
@@ -54,6 +54,7 @@
 #include "llvm/Support/LEB128.h"
 #include "llvm/Support/Path.h"
 #include "llvm/Support/TarWriter.h"
+#include "llvm/Support/TargetRegistry.h"
 #include "llvm/Support/TargetSelect.h"
 #include "llvm/Support/raw_ostream.h"
 #include <cstdlib>
@@ -309,6 +310,9 @@ static void checkOptions(opt::InputArgLi
 
     if (Config->SingleRoRx && !Script->HasSectionsCommand)
       error("-execute-only and -no-rosegment cannot be used together");
+  } else if (Config->TargetTriple.isOSNetBSD()) {
+    // force-disable RO segment on NetBSD due to ld.elf_so limitations
+    Config->SingleRoRx = true;
   }
 }
 
@@ -341,7 +345,7 @@ static bool isKnown(StringRef S) {
          S == "execstack" || S == "hazardplt" || S == "initfirst" ||
          S == "keep-text-section-prefix" || S == "lazy" || S == "muldefs" ||
          S == "nocombreloc" || S == "nocopyreloc" || S == "nodelete" ||
-         S == "nodlopen" || S == "noexecstack" ||
+         S == "nodlopen" || S == "noexecstack" || S == "nognustack" ||
          S == "nokeep-text-section-prefix" || S == "norelro" || S == "notext" ||
          S == "now" || S == "origin" || S == "relro" || S == "retpolineplt" ||
          S == "rodynamic" || S == "text" || S == "wxneeded" ||
@@ -355,6 +359,56 @@ static void checkZOptions(opt::InputArgL
       error("unknown -z value: " + StringRef(Arg->getValue()));
 }
 
+void LinkerDriver::appendDefaultSearchPaths() {
+  if (Config->TargetTriple.isOSNetBSD()) {
+    // NetBSD driver relies on the linker knowing the default search paths.
+    // Please keep this in sync with clang/lib/Driver/ToolChains/NetBSD.cpp
+    // (NetBSD::NetBSD constructor)
+    switch (Config->TargetTriple.getArch()) {
+    case llvm::Triple::x86:
+      Config->SearchPaths.push_back("=/usr/lib/i386");
+      break;
+    case llvm::Triple::arm:
+    case llvm::Triple::armeb:
+    case llvm::Triple::thumb:
+    case llvm::Triple::thumbeb:
+      switch (Config->TargetTriple.getEnvironment()) {
+      case llvm::Triple::EABI:
+      case llvm::Triple::GNUEABI:
+        Config->SearchPaths.push_back("=/usr/lib/eabi");
+        break;
+      case llvm::Triple::EABIHF:
+      case llvm::Triple::GNUEABIHF:
+        Config->SearchPaths.push_back("=/usr/lib/eabihf");
+        break;
+      default:
+        Config->SearchPaths.push_back("=/usr/lib/oabi");
+        break;
+      }
+      break;
+#if 0 // TODO
+    case llvm::Triple::mips64:
+    case llvm::Triple::mips64el:
+      if (tools::mips::hasMipsAbiArg(Args, "o32"))
+        Config->SearchPaths.push_back("=/usr/lib/o32");
+      else if (tools::mips::hasMipsAbiArg(Args, "64"))
+        Config->SearchPaths.push_back("=/usr/lib/64");
+      break;
+#endif
+    case llvm::Triple::ppc:
+      Config->SearchPaths.push_back("=/usr/lib/powerpc");
+      break;
+    case llvm::Triple::sparc:
+      Config->SearchPaths.push_back("=/usr/lib/sparc");
+      break;
+    default:
+      break;
+    }
+
+    Config->SearchPaths.push_back("=/usr/lib");
+  }
+}
+
 void LinkerDriver::main(ArrayRef<const char *> ArgsArr) {
   ELFOptTable Parser;
   opt::InputArgList Args = Parser.parse(ArgsArr.slice(1));
@@ -368,6 +422,29 @@ void LinkerDriver::main(ArrayRef<const c
     return;
   }
 
+  if (const char *Path = getReproduceOption(Args)) {
+    // Note that --reproduce is a debug option so you can ignore it
+    // if you are trying to understand the whole picture of the code.
+    Expected<std::unique_ptr<TarWriter>> ErrOrWriter =
+        TarWriter::create(Path, path::stem(Path));
+    if (ErrOrWriter) {
+      Tar = ErrOrWriter->get();
+      Tar->append("response.txt", createResponseFile(Args));
+      Tar->append("version.txt", getLLDVersion() + "\n");
+      make<std::unique_ptr<TarWriter>>(std::move(*ErrOrWriter));
+    } else {
+      error(Twine("--reproduce: failed to open ") + Path + ": " +
+            toString(ErrOrWriter.takeError()));
+    }
+  }
+
+
+  initLLVM();
+  setTargetTriple(ArgsArr[0], Args);
+  readConfigs(Args);
+  checkZOptions(Args);
+  appendDefaultSearchPaths();
+
   // Handle -v or -version.
   //
   // A note about "compatible with GNU linkers" message: this is a hack for
@@ -383,8 +460,10 @@ void LinkerDriver::main(ArrayRef<const c
   // lot of "configure" scripts out there that are generated by old version
   // of Libtool. We cannot convince every software developer to migrate to
   // the latest version and re-generate scripts. So we have this hack.
-  if (Args.hasArg(OPT_v) || Args.hasArg(OPT_version))
+  if (Args.hasArg(OPT_v) || Args.hasArg(OPT_version)) {
     message(getLLDVersion() + " (compatible with GNU linkers)");
+    message("Target: " + Config->TargetTriple.str());
+  }
 
   // The behavior of -v or --version is a bit strange, but this is
   // needed for compatibility with GNU linkers.
@@ -393,25 +472,6 @@ void LinkerDriver::main(ArrayRef<const c
   if (Args.hasArg(OPT_version))
     return;
 
-  if (const char *Path = getReproduceOption(Args)) {
-    // Note that --reproduce is a debug option so you can ignore it
-    // if you are trying to understand the whole picture of the code.
-    Expected<std::unique_ptr<TarWriter>> ErrOrWriter =
-        TarWriter::create(Path, path::stem(Path));
-    if (ErrOrWriter) {
-      Tar = ErrOrWriter->get();
-      Tar->append("response.txt", createResponseFile(Args));
-      Tar->append("version.txt", getLLDVersion() + "\n");
-      make<std::unique_ptr<TarWriter>>(std::move(*ErrOrWriter));
-    } else {
-      error(Twine("--reproduce: failed to open ") + Path + ": " +
-            toString(ErrOrWriter.takeError()));
-    }
-  }
-
-  readConfigs(Args);
-  checkZOptions(Args);
-  initLLVM();
   createFiles(Args);
   if (errorCount())
     return;
@@ -725,6 +785,34 @@ static void parseClangOption(StringRef O
   error(Msg + ": " + StringRef(Err).trim());
 }
 
+void LinkerDriver::setTargetTriple(StringRef argv0, opt::InputArgList &Args) {
+  std::string TargetError;
+
+  // Firstly, see if user specified explicit --target
+  StringRef TargetOpt = Args.getLastArgValue(OPT_target);
+  if (!TargetOpt.empty()) {
+    if (llvm::TargetRegistry::lookupTarget(TargetOpt, TargetError)) {
+      Config->TargetTriple = llvm::Triple(TargetOpt);
+      return;
+    } else
+      error("Unsupported --target=" + TargetOpt + ": " + TargetError);
+  }
+
+  // Secondly, try to get it from program name prefix
+  std::string ProgName = llvm::sys::path::stem(argv0);
+  size_t LastComponent = ProgName.rfind('-');
+  if (LastComponent != std::string::npos) {
+    std::string Prefix = ProgName.substr(0, LastComponent);
+    if (llvm::TargetRegistry::lookupTarget(Prefix, TargetError)) {
+      Config->TargetTriple = llvm::Triple(Prefix);
+      return;
+    }
+  }
+
+  // Finally, use the default target triple
+  Config->TargetTriple = llvm::Triple(getDefaultTargetTriple());
+}
+
 // Initializes Config members by the command line options.
 void LinkerDriver::readConfigs(opt::InputArgList &Args) {
   errorHandler().Verbose = Args.hasArg(OPT_verbose);
@@ -755,7 +843,8 @@ void LinkerDriver::readConfigs(opt::Inpu
       Args.hasFlag(OPT_eh_frame_hdr, OPT_no_eh_frame_hdr, false);
   Config->EmitRelocs = Args.hasArg(OPT_emit_relocs);
   Config->EnableNewDtags =
-      Args.hasFlag(OPT_enable_new_dtags, OPT_disable_new_dtags, true);
+      Args.hasFlag(OPT_enable_new_dtags, OPT_disable_new_dtags,
+                   !Config->TargetTriple.isOSNetBSD());
   Config->Entry = Args.getLastArgValue(OPT_entry);
   Config->ExecuteOnly =
       Args.hasFlag(OPT_execute_only, OPT_no_execute_only, false);
@@ -842,6 +931,8 @@ void LinkerDriver::readConfigs(opt::Inpu
   Config->ZCombreloc = getZFlag(Args, "combreloc", "nocombreloc", true);
   Config->ZCopyreloc = getZFlag(Args, "copyreloc", "nocopyreloc", true);
   Config->ZExecstack = getZFlag(Args, "execstack", "noexecstack", false);
+  Config->ZNognustack = hasZOption(Args, "nognustack") ||
+    Config->TargetTriple.isOSNetBSD();
   Config->ZHazardplt = hasZOption(Args, "hazardplt");
   Config->ZInitfirst = hasZOption(Args, "initfirst");
   Config->ZKeepTextSectionPrefix = getZFlag(
@@ -1137,7 +1228,7 @@ void LinkerDriver::inferMachineType() {
 // each target.
 static uint64_t getMaxPageSize(opt::InputArgList &Args) {
   uint64_t Val = args::getZOptionValue(Args, OPT_z, "max-page-size",
-                                       Target->DefaultMaxPageSize);
+                                       lld::elf::Target->DefaultMaxPageSize);
   if (!isPowerOf2_64(Val))
     error("max-page-size: value isn't a power of 2");
   return Val;
