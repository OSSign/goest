FROM mcr.microsoft.com/windows/nanoserver:ltsc2025

WORKDIR C:\a

COPY goest.exe C:\a\goest.exe

ENTRYPOINT [ "C:\a\goest.exe" ]