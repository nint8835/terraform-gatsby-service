FROM nint8835/terraform-provider-gatsby

WORKDIR /terraform
COPY entrypoint.py /terraform
RUN apk add python3
ENTRYPOINT [ "python3", "entrypoint.py" ]
