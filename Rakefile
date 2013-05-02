task :setup do
  sh 'cd src && git clone git://github.com/ChimeraCoder/anaconda.git'
  sh 'cd src && git clone git://github.com/go-sql-driver/mysql.git'
end

task :build do
  sh 'go build scraper'
end

task :run do
  sh './scraper'
end

task :clean do
  sh 'rm scraper'
end

task :default => [:build, :run, :clean]

