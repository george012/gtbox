Pod::Spec.new do |spec|
  spec.name         = 'gtbox'
  spec.version      = '1.0.0'
  spec.summary      = 'A short description of gtbox.'
  spec.homepage     = 'https://github.com/george012/gtbox'
  spec.license      = { :type => 'MIT', :file => 'LICENSE' }
  spec.author       = { 'george012' => 'https://github.com/george012/gtbox' }
  spec.source       = { :git => 'https://github.com/george012/gtbox.git', :tag => '0.0.1' }
  spec.source_files = 'gtbox/**/*.{h,m,swift,a.dylib}'
  spec.requires_arc = true
end
