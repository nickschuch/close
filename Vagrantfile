# Provides a local development environment for Close.
#

Vagrant.configure("2") do |config|
  config.vm.box       = "ubuntu/trusty64"
  config.vm.host_name = "close.dev"
  
  config.vm.network :private_network, :ip => "192.168.80.30"
  
  config.vm.synced_folder ".", "/opt/golang/src/github.com/nickschuch/close"
  
  config.vm.provider :virtualbox do |vb|
    vb.customize ["modifyvm", :id, "--memory", 1024]
  end

  config.vm.provision "shell", path: "scripts/provision.sh"
end
