FROM scratch

COPY goest /goest

ENTRYPOINT [ "/goest" ]