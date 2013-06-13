require 'formula'

class Gh < Formula
  VERSION = '0.6.1'

  homepage 'https://github.com/jingweno/gh'
  url "https://github.com/jingweno/gh/archive/#{VERSION}.tar.gz"
  sha1 '1e4ca70ebf018ae192a641f18b735beca5df5c31'
  version VERSION

  head 'https://github.com/jingweno/gh.git'

  depends_on 'hg'
  depends_on 'go'

  def install
    go_path = Dir.getwd
    system "GOPATH='#{go_path}' go get -d ./..."
    system "GOPATH='#{go_path}' go build -o gh"
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
