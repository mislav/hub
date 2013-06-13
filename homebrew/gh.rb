require 'formula'

class Gh < Formula
  VERSION = '0.6.1'
  ARCH = if MacOS.prefer_64_bit?
           'amd64'
         else
           '386'
         end

  homepage 'https://github.com/jingweno/gh'
  url "https://github.com/jingweno/gh/archive/#{VERSION}.tar.gz"
  version VERSION

  head 'https://github.com/jingweno/gh.git'

  depends_on 'go'

  def install
    system 'go build -o gh'
    bin.install 'gh'
  end

  def caveats; <<-EOS.undent
  To upgrade gh, run `brew upgrade https://raw.github.com/jingweno/gh/master/homebrew/gh.rb`

  More information here: https://github.com/jingweno/gh/blob/master/README.md
    EOS
  end

  test do
    assert_equal VERSION, `#{bin}/gh version`.split.last
  end
end
