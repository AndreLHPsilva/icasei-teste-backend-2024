task :prestart => :environment do
  puts "Rodando as migrations antes de iniciar..."
  Rake::Task['db:migrate'].invoke
end