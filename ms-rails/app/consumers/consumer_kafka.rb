# frozen_string_literal: true
class ConsumerKafka < ApplicationConsumer
  def consume
    messages.each do |message|
      payload = message.payload.with_indifferent_access
      puts "Received message: #{payload}"
      
      begin
        product = payload[:product]
        upsert_service = Services::Api::V1::Products::Upsert.new(product.merge({ is_api: false }), nil)
        product = upsert_service.execute
        
        puts "Processed product: #{product}"
      rescue => e
        puts "Failed to process message: #{e.message}"
      end
    end
  end
end
