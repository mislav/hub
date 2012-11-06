require 'strscan'
require 'forwardable'

# Stupid pure Ruby JSON parser.
class Hub::JSON
  def self.parse(data) new(data).parse end

  WSP = /\s+/
  OBJ = /[{\[]/;    HEN = /\}/;  AEN = /\]/
  COL = /\s*:\s*/;  KEY = /\s*,\s*/
  NUM = /-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?/
  BOL = /true|false/;  NUL = /null/

  extend Forwardable

  attr_reader :scanner
  alias_method :s, :scanner
  def_delegators :scanner, :scan, :matched
  private :s, :scan, :matched

  def initialize data
    @scanner = StringScanner.new data.to_s
  end

  def parse
    space
    object
  end

  private

  def space() scan WSP end

  def endkey() scan(KEY) or space end

  def object
    matched == '{' ? hash : array if scan(OBJ)
  end

  def value
    object or string or
      scan(NUL) ? nil :
      scan(BOL) ? matched.size == 4:
      scan(NUM) ? eval(matched) :
      error
  end

  def hash
    obj = {}
    space
    repeat_until(HEN) { k = string; scan(COL); obj[k] = value; endkey }
    obj
  end

  def array
    ary = []
    space
    repeat_until(AEN) { ary << value; endkey }
    ary
  end

  SPEC = {'b' => "\b", 'f' => "\f", 'n' => "\n", 'r' => "\r", 't' => "\t"}
  UNI = 'u'; CODE = /[a-fA-F0-9]{4}/
  STR = /"/; STE = '"'
  ESC = '\\'

  def string
    if scan(STR)
      str, esc = '', false
      while c = s.getch
        if esc
          str << (c == UNI ? (s.scan(CODE) || error).to_i(16).chr : SPEC[c] || c)
          esc = false
        else
          case c
          when ESC then esc = true
          when STE then break
          else str << c
          end
        end
      end
      str
    end
  end

  def error
    raise "parse error at: #{scan(/.{1,10}/m).inspect}"
  end

  def repeat_until reg
    until scan(reg)
      pos = s.pos
      yield
      error unless s.pos > pos
    end
  end

  module Generator
    def generate(obj)
      raise ArgumentError unless obj.is_a? Array or obj.is_a? Hash
      generate_type(obj)
    end
    alias dump generate

    private

    def generate_type(obj)
      type = obj.is_a?(Numeric) ? :Numeric : obj.class.name
      begin send(:"generate_#{type}", obj)
      rescue NoMethodError; raise ArgumentError, "can't serialize #{type}"
      end
    end

    ESC_MAP = Hash.new {|h,k| k }.update \
      "\r" => 'r',
      "\n" => 'n',
      "\f" => 'f',
      "\t" => 't',
      "\b" => 'b'

    def generate_String(str)
      escaped = str.gsub(/[\r\n\f\t\b"\\]/) { "\\#{ESC_MAP[$&]}"}
      %("#{escaped}")
    end

    def generate_simple(obj) obj.inspect end
    alias generate_Numeric generate_simple
    alias generate_TrueClass generate_simple
    alias generate_FalseClass generate_simple

    def generate_Symbol(sym) generate_String(sym.to_s) end

    def generate_NilClass(*) 'null' end

    def generate_Array(ary) '[%s]' % ary.map {|o| generate_type(o) }.join(', ') end

    def generate_Hash(hash)
      '{%s}' % hash.map { |key, value|
        "#{generate_String(key.to_s)}: #{generate_type(value)}"
      }.join(', ')
    end
  end

  extend Generator
end
