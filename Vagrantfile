# -*- mode: ruby -*-
# vi: set ft=ruby :

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
#
# Place this Vagrantfile in your src folder and run:
#
#     vagrant up
#
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

# Vagrantfile API/syntax version.
VAGRANTFILE_API_VERSION = "2"

GO_ARCHIVES = {
  "linux" => "go1.2.linux-amd64.tar.gz",
}

INSTALL = {
  "linux" => "apt-get update -qq; apt-get install -qq -y git mercurial bzr curl",
}

# location of the Vagrantfile
def src_path
  ENV["GOPATH"]
end

# shell script to bootstrap Go
def bootstrap(box)
  install = INSTALL[box]
  archive = GO_ARCHIVES[box]
  vagrant_home = "/home/vagrant"

  profile = <<-PROFILE
    export GOROOT=#{vagrant_home}/go
    export GOPATH=#{vagrant_home}/gocode
    export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
    export CDPATH=.:$GOPATH/src/github.com:$GOPATH/src/code.google.com/p:$GOPATH/src/bitbucket.org:$GOPATH/src/launchpad.net
  PROFILE

  <<-SCRIPT
  #{install}

  if ! [ -f /home/vagrant/#{archive} ]; then
    curl -O# https://go.googlecode.com/files/#{archive}
  fi
  tar -C /home/vagrant -xzf #{archive}
  chown -R vagrant:vagrant #{vagrant_home}/go

  if ! grep -q GOPATH #{vagrant_home}/.profile; then
    echo '#{profile}' >> #{vagrant_home}/.profile
  fi
  source #{vagrant_home}/.profile

  chown -R vagrant:vagrant #{vagrant_home}/gocode
  go get github.com/jingweno/gotask

  echo "\nRun: vagrant ssh #{box} -c 'cd project/path; go test ./...'"
  SCRIPT
end

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  config.vm.define "linux" do |linux|
    linux.vm.box = "precise64"
    linux.vm.box_url = "http://files.vagrantup.com/precise64.box"
    linux.vm.synced_folder "#{src_path}/src/github.com/jingweno/gh", "/home/vagrant/gocode/src/github.com/jingweno/gh"
    linux.vm.provision :shell, :inline => bootstrap("linux")
  end

end
