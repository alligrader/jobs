task :default => [:all]

task build_findbugs: [ ] do
    STDOUT.puts "What should the Findbugs image version tag be?"
    input = STDIN.gets.strip

    if input == ''
        Rake::Task["build_findbugs"].reenable
        Rake::Task["build_findbugs"].invoke
    else
        sh "docker build -t thesnowmancometh/findbugs:#{input} -f findbugs/Findbugs_Dockerfile findbugs"
    end
end


task all: [ :build ] do

end

task build: [ :build_findbugs ] do
end
