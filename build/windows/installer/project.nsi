Unicode true

####
## Please note: Template replacements don't work in this file. They are provided with default defines like
## mentioned underneath.
## If the keyword is not defined, "wails_tools.nsh" will populate them with the values from ProjectInfo.
## If they are defined here, "wails_tools.nsh" will not touch them. This allows to use this project.nsi manually
## from outside of Wails for debugging and development of the installer.
##
## For development first make a wails nsis build to populate the "wails_tools.nsh":
## > wails build --target windows/amd64 --nsis
## Then you can call makensis on this file with specifying the path to your binary:
## For a AMD64 only installer:
## > makensis -DARG_WAILS_AMD64_BINARY=..\..\bin\app.exe
## For a ARM64 only installer:
## > makensis -DARG_WAILS_ARM64_BINARY=..\..\bin\app.exe
## For a installer with both architectures:
## > makensis -DARG_WAILS_AMD64_BINARY=..\..\bin\app-amd64.exe -DARG_WAILS_ARM64_BINARY=..\..\bin\app-arm64.exe
####
## The following information is taken from the ProjectInfo file, but they can be overwritten here.
####
## !define INFO_PROJECTNAME    "MyProject" # Default "{{.Name}}"
## !define INFO_COMPANYNAME    "MyCompany" # Default "{{.Info.CompanyName}}"
## !define INFO_PRODUCTNAME    "MyProduct" # Default "{{.Info.ProductName}}"
## !define INFO_PRODUCTVERSION "1.0.0"     # Default "{{.Info.ProductVersion}}"
## !define INFO_COPYRIGHT      "Copyright" # Default "{{.Info.Copyright}}"
###
## !define PRODUCT_EXECUTABLE  "Application.exe"      # Default "${INFO_PROJECTNAME}.exe"
## !define UNINST_KEY_NAME     "UninstKeyInRegistry"  # Default "${INFO_COMPANYNAME}${INFO_PRODUCTNAME}"
####
## !define REQUEST_EXECUTION_LEVEL "admin"            # Default "admin"  see also https://nsis.sourceforge.io/Docs/Chapter4.html
####
## Include the wails tools
####
!include "wails_tools.nsh"

# The version information for this two must consist of 4 parts
VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion    "${INFO_PRODUCTVERSION}.0"

VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

# Enable HiDPI support. https://nsis.sourceforge.io/Reference/ManifestDPIAware
ManifestDPIAware true

!include "MUI.nsh"

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"
# !define MUI_WELCOMEFINISHPAGE_BITMAP "resources\leftimage.bmp" #Include this to add a bitmap on the left side of the Welcome Page. Must be a size of 164x314
!define MUI_FINISHPAGE_NOAUTOCLOSE # Wait on the INSTFILES page so the user can take a look into the details of the installation steps
!define MUI_ABORTWARNING # This will warn the user if they exit from the installer.

!insertmacro MUI_PAGE_WELCOME # Welcome to the installer page.
# !insertmacro MUI_PAGE_LICENSE "resources\eula.txt" # Adds a EULA page to the installer
!insertmacro MUI_PAGE_DIRECTORY # In which folder install page.
!insertmacro MUI_PAGE_INSTFILES # Installing page.
!insertmacro MUI_PAGE_FINISH # Finished installation page.

!insertmacro MUI_UNPAGE_INSTFILES # Uinstalling page

!insertmacro MUI_LANGUAGE "English" # Set the Language of the installer

## The following two statements can be used to sign the installer and the uninstaller. The path to the binaries are provided in %1
#!uninstfinalize 'signtool --file "%1"'
#!finalize 'signtool --file "%1"'

Name "${INFO_PRODUCTNAME}"
OutFile "..\..\bin\${INFO_PROJECTNAME}-${ARCH}-installer.exe" # Name of the installer's file.
!ifdef WAILS_INSTALL_SCOPE
  !if "${WAILS_INSTALL_SCOPE}" == "user"
    InstallDir "$LOCALAPPDATA\Programs\${INFO_PRODUCTNAME}"
  !else
    # No ${INFO_COMPANYNAME} subfolder here on purpose: with a single
    # product where CompanyName and ProductName happen to be the same
    # ("NutriPlan", see wails.json's info block), the Wails default
    # "$PROGRAMFILES64\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}" produces a
    # doubled-up "C:\Program Files\NutriPlan\NutriPlan" - fine for a
    # company with multiple products (mirrors how e.g. Microsoft nests
    # "Program Files\Microsoft\Office"), but redundant here.
    InstallDir "$PROGRAMFILES64\${INFO_PRODUCTNAME}"
  !endif
!else
  InstallDir "$PROGRAMFILES64\${INFO_PRODUCTNAME}"
!endif # Default installing folder ($PROGRAMFILES is Program Files folder).
ShowInstDetails show # This will always show the installation details.

Function .onInit
   !insertmacro wails.checkArchitecture
FunctionEnd

Section
    !insertmacro wails.setShellContext

    !insertmacro wails.webview2runtime

    SetOutPath $INSTDIR

    !insertmacro wails.files

    # Bundle a production-appropriate runtime config next to the
    # executable, so a fresh install has a working, correctly-scoped
    # config.yaml/.env from the start instead of silently generating
    # brand-new, unthemed/unconfigured ones on first run (see
    # config.Load()'s generateConfig() fallback in the Go source).
    #
    # This intentionally bundles config.installer.yaml (renamed to
    # config.yaml here), NOT the project root's own config.yaml - the
    # installer variant has empty data/temp/log paths (auto-resolved to
    # %AppData%\NutriPlan at runtime - see config.DefaultAppDataDir) and
    # environment: production/debug: false, instead of the dev config's
    # convenient-for-local-testing relative paths and debug settings.
    # NOTE: double-check your local .env has no real secrets in it before
    # building a public installer - whatever's in it gets bundled verbatim.
    File "/oname=config.yaml" "..\..\..\config.installer.yaml"
    File "..\..\..\.env"

    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
    CreateShortCut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"

    !insertmacro wails.associateFiles
    !insertmacro wails.associateCustomProtocols

    !insertmacro wails.writeUninstaller
SectionEnd

Section "uninstall"
    !insertmacro wails.setShellContext

    RMDir /r "$AppData\${PRODUCT_EXECUTABLE}" # Remove the WebView2 DataPath

    RMDir /r $INSTDIR

    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.unassociateFiles
    !insertmacro wails.unassociateCustomProtocols

    !insertmacro wails.deleteUninstaller

    # The app's own data (client database, logs, temp files) lives in
    # "%AppData%\NutriPlan" - deliberately OUTSIDE $INSTDIR (see
    # config.DefaultAppDataDir's docs for why), so the RMDir /r $INSTDIR
    # above never touches it. Ask explicitly rather than silently keeping
    # or silently deleting it - most people expect an uninstall to leave
    # their data alone (reinstalling later should pick up where they left
    # off), but some genuinely want a clean removal.
    #
    # IMPORTANT: this deliberately uses $PROFILE\AppData\Roaming, NOT
    # NSIS's own $APPDATA variable. $APPDATA is shell-context-sensitive:
    # wails.setShellContext calls SetShellVarContext all whenever
    # REQUEST_EXECUTION_LEVEL is "admin" (the default for this installer),
    # which makes $APPDATA resolve to the *all-users* common app data
    # folder (CSIDL_COMMON_APPDATA, i.e. C:\ProgramData) instead of the
    # actual logged-in user's own Roaming folder. The Go app, on the other
    # hand, always writes to the CURRENT user's %AppData% via
    # os.UserConfigDir() - regardless of the installer's elevation level.
    # Using $APPDATA here would silently look in the wrong place and
    # delete nothing, which is exactly the bug this comment is preventing
    # you from reintroducing. $PROFILE (CSIDL_PROFILE / %USERPROFILE%) has
    # no "all users" equivalent, so it's unaffected by SetShellVarContext
    # and always matches what the app itself resolved at runtime.
    MessageBox MB_YESNO "Do you also want to remove your saved data (client database, logs, exported PDFs are NOT affected)? This cannot be undone." IDYES removeAppData IDNO keepAppData
    removeAppData:
        RMDir /r "$PROFILE\AppData\Roaming\NutriPlan"
    keepAppData:
SectionEnd
