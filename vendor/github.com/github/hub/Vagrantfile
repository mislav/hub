# -*- mode: ruby -*-
# vi: set ft=ruby :

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
#
# Place this Vagrantfile in your src folder and run:
#
#     vagrant up
#
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

GO_ARCHIVES = {
  "linux" => "go1.4.2.linux-amd64.tar.gz",
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
    export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
  PROFILE

  <<-SCRIPT
  #{install}

  if ! [ -f /home/vagrant/#{archive} ]; then
    curl -O# https://storage.googleapis.com/golang/#{archive}
  fi
  tar -C /home/vagrant -xzf #{archive}
  chown -R vagrant:vagrant #{vagrant_home}/go

  if ! grep -q GOPATH #{vagrant_home}/.bashrc; then
    echo '#{profile}' >> #{vagrant_home}/.bashrc
  fi
  source #{vagrant_home}/.bashrc

  chown -R vagrant:vagrant #{vagrant_home}/gocode

  apt-get update -qq
  apt-get install -qq ruby1.9.1-dev tmux zsh git
  gem install bundler

  echo "\nRun: vagrant ssh #{box} -c 'cd project/path; go test ./...'"
  SCRIPT
end

Vagrant.configure("2") do |config|
  config.vm.define "linux" do |linux|
    linux.vm.box = "ubuntu/trusty64"
    linux.vm.synced_folder "#{src_path}/src/github.com/github/hub", "/home/vagrant/gocode/src/github.com/github/hub"
    linux.vm.provision :shell, :inline => bootstrap("linux")
  end
end
