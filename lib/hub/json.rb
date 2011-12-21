require 'strscan'
require 'forwardable'

class Hub::JSON
  WSP = /\s+/
  OBJ = /[{\[]/
  NUM = /-?\d+(\.\d+)?([eE][+-]?\d+)?/
  BOL = /(?:true|false)\b/
  NUL = /null\b/

  extend Forwardable

  attr_reader :scanner
  alias_method :s, :scanner
  private :s
  def_delegators :scanner, :scan, :matched

  def self.parse(str)
    self.new(str).parse
  end

  def initialize data
    @scanner = StringScanner.new data.to_s
  end

  def space
    scan WSP
  end

  def parse
    space
    object
  end

  def object
    case scan(OBJ)
    when '{' then hash
    when '[' then array
    end
  end

  def value
    object or string or
      (scan(NUM) || scan(BOL)) ? eval(matched) :
      scan(NUL) && nil
  end

  def hash
    current = {}
    space
    until scan(/\}/)
      key = string
      scan(/\s*:\s*/)
      current[key] = value
      scan(/\s*,\s*/) or space
    end
    current
  end

  def array
    current = []
    space
    until scan(/\]/)
      current << value
      scan(/\s*,\s*/) or space
    end
    current
  end

  def string
    if str = scan(/"/)
      begin; str << s.scan_until(/"/); end while s.pre_match[-1,1] == '\\'
      eval str
    end
  end
end
