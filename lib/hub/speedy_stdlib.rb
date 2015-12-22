prevent_require = lambda do |name|
  $" << "#{name}.rb"
  require name # hax to avoid Ruby 2.0.0 bug
end

unless defined?(CGI)
  prevent_require.call 'cgi'

  module CGI
    ESCAPE_RE = /[^a-zA-Z0-9 .~_-]/

    def self.escape(s)
      s.to_s.gsub(ESCAPE_RE) {|match|
        '%' + match.unpack('H2' * match.bytesize).join('%').upcase
      }.tr(' ', '+')
    end

    def self.unescape(s)
      s.tr('+', ' ').gsub(/((?:%[0-9a-fA-F]{2})+)/) {
        [$1.delete('%')].pack('H*')
      }
    end
  end
end

require 'uri'
