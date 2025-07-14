FROM scratch

COPY goest.exe /goest.exe

ENTRYPOINT [ "/goest.exe" ]