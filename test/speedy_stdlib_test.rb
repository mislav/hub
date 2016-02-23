require 'helper'

class URITest < Minitest::Test
  def test_uri_display_port
    assert_equal "https://example.com", URI.parse("https://example.com").to_s
    assert_equal "https://example.com:80", URI.parse("https://example.com:80").to_s
    assert_equal "http://example.com", URI.parse("http://example.com").to_s
    assert_equal "http://example.com:443", URI.parse("http://example.com:443").to_s
  end

  def test_uri_invalid_port
    assert_raises URI::InvalidComponentError do
      uri = URI.parse("https://example.com")
      uri.port = "80"
    end
  end

  def test_uri_scheme_doesnt_affect_port
    uri = URI.parse("https://example.com")
    uri.scheme = "http"
    assert_equal "http", uri.scheme
    assert_equal 443, uri.port
    uri.port = 80
    assert_equal 80, uri.port
  end

  def test_blank_path
    uri = URI.parse("https://example.com")
    assert_equal "", uri.path
  end

  def test_no_query
    uri = URI.parse("https://example.com")
    assert_nil uri.query
  end

  def test_blank_query
    uri = URI.parse("https://example.com?")
    assert_equal "", uri.query
  end
end
