FROM bitnami/nginx:1.27.4-debian-12-r6
# Allow modifications to the image
USER 0 
# Just an example for a custom image
ADD index.html /app/index.html
# Run nginx as a non-root user
USER 1001