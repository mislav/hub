require 'formula'

class Gh < Formula
  VERSION = '0.0.1'

  homepage 'https://github.com/jingweno/gh'
  url "https://drone.io/github.com/jingweno/gh/files/target/#{VERSION}-snapshot/darwin_amd64/gh_#{VERSION}-snapshot_darwin_amd64.zip"
  version VERSION
  sha1 '12743626cd717014c3ddd8ac70d77da44356fa66'
  head 'https://github.com/jingweno/gh.git'

  def install
    bin.install 'gh'
  end

  test do
    `#{bin}/gh version`.chomp == VERSION
  end
end
