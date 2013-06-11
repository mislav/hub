require 'formula'

class Gh < Formula
  VERSION = '0.6.0'
  ARCH = if MacOS.prefer_64_bit?
           'amd64'
         else
           '386'
         end

  homepage 'https://github.com/jingweno/gh'
  version VERSION
  url "https://drone.io/github.com/jingweno/gh/files/target/#{VERSION}-snapshot/darwin_#{ARCH}/gh_#{VERSION}-snapshot_darwin_#{ARCH}.tar.gz"
  head 'https://github.com/jingweno/gh.git'

  def install
    bin.install 'gh'
  end

  test do
    assert_equal VERSION, `#{bin}/gh version`.split.last
  end
end
