#https://github.com/wader/disable_sendfile_vbox_linux/blob/master/disable_sendfile_vbox_linux.go is terrible solution
# because nobody wants different code in dev and prod

# this is the bug in virtualbox that this script gets around for me by forcing the 404
# https://www.vagrantup.com/docs/synced-folders/virtualbox.html
mv js/app.js js/app.js2
curl http://192.168.99.100:8080/js/app.js
mv js/app.js2 js/app.js
docker-compose restart golang
docker-compose logs
