FROM scratch 
ADD magnet /magnet/
WORKDIR /magnet
ADD templates/ /magnet/templates/
ADD public/ /magnet/public/
EXPOSE 3000
CMD ["/magnet/magnet"]
