require "formula"

class Gh < Formula
  VERSION = "0.15.0"
  ARCH = if MacOS.prefer_64_bit?
           "amd64"
         else
           "386"
         end

  homepage "https://github.com/jingweno/gh"
  head "https://github.com/jingweno/gh.git"
  url "https://dl.dropboxusercontent.com/u/1079131/gh/#{VERSION}-snapshot/darwin_#{ARCH}/gh_#{VERSION}-snapshot_darwin_#{ARCH}.tar.gz"
  version VERSION


  def install
    bin.install "gh"
    bash_completion.install "etc/gh.bash_completion.sh"
    zsh_completion.install "etc/gh.zsh_completion" => "_gh"
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
