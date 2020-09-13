FROM nint8835/terraform-provider-gatsby

WORKDIR /terraform
COPY entrypoint.py /terraform
RUN apk add python3 && \
    addgroup -S terraform && \
    adduser -S terraform -G terraform && \
    chown -R terraform /terraform
USER terraform
ENTRYPOINT [ "python3", "entrypoint.py" ]
