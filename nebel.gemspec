# -*- encoding: utf-8 -*-

Gem::Specification.new do |s|
  s.name        = "nebel"
  s.version     = "0.0.2"
  s.authors     = ["Gosuke Miyashita"]
  s.email       = ["gosukenator@gmail.com"]
  s.homepage    = "https://github.com/mizzy/nebel"
  s.summary     = %q{A command line tool for generating a static site.}

  s.files         = `git ls-files`.split("\n")
  s.test_files    = `git ls-files -- {test,spec,features}/*`.split("\n")
  s.executables   = `git ls-files -- bin/*`.split("\n").map{ |f| File.basename(f) }

  # dependencies
  s.add_dependency 'liquid', '2.2.2'
  s.add_dependency 'redcarpet'
end
