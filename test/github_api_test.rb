require 'helper'
require 'hub/github_api'
require 'forwardable'
require 'delegate'

class FileStoreTest < Minitest::Test
  extend Forwardable
  def_delegators :@store, :yaml_dump, :yaml_load

  def setup
    @store = Hub::GitHubAPI::FileStore.new('')
  end

  class OrderedHash < DelegateClass(::Hash)
    def self.[](*args)
      hash = new
      while args.any?
        key, value = args.shift, args.shift
        hash[key] = value
      end
      hash
    end

    def initialize(hash = {})
      @keys = hash.keys
      super(hash)
    end

    def []=(key, value) @keys << key; super end

    def each
      @keys.each { |key| yield(key, self[key]) }
    end
  end

  def test_yaml_dump
    output = yaml_dump("github.com" => [
      OrderedHash['user', 'mislav', 'oauth_token', 'OTOKEN'],
      OrderedHash['user', 'tpw', 'oauth_token', 'POKEN'],
    ])

    assert_equal <<-YAML.chomp, output
---
github.com:
- user: mislav
  oauth_token: OTOKEN
- user: tpw
  oauth_token: POKEN
    YAML
  end

  def test_yaml_load
    data = yaml_load <<-YAML
---
github.com:
- user: mislav
  oauth_token: OTOKEN
- user: tpw
  oauth_token: POKEN
    YAML

    assert_equal 'mislav', data['github.com'][0]['user']
    assert_equal 'OTOKEN', data['github.com'][0]['oauth_token']
    assert_equal 'tpw',    data['github.com'][1]['user']
    assert_equal 'POKEN',  data['github.com'][1]['oauth_token']
  end

  def test_yaml_load_quoted
    data = yaml_load <<-YAML
---
github.com:
- user: 'true'
  oauth_token: '1234'
    YAML

    assert_equal 'true', data['github.com'][0]['user']
    assert_equal '1234', data['github.com'][0]['oauth_token']
  end
end
