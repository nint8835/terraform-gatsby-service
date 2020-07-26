FROM nint8835/terraform-provider-gatsby

WORKDIR /terraform
COPY entrypoint.sh /terraform
RUN chmod +x entrypoint.sh
ENTRYPOINT [ "sh", "entrypoint.sh" ]
