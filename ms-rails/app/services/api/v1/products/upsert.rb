module Services
  module Api 
    module V1
      module Products
        class Upsert
          attr_accessor :params, :request, :new_record

          def initialize(params, request)
            @params = params
            @request = request
            @new_record = false
          end

          def execute
            raise ActiveRecord::RecordNotFound, "Product Not Found" if product.blank?

            ActiveRecord::Base.transaction do
              @new_record = product.new_record?
              
              product.id          ||= params[:id]           if params[:id].present?
              product.name          = params[:name]         if params[:name].present?
              product.brand         = params[:brand]        if params[:brand].present?
              product.price         = params[:price]        if params[:price].present?
              product.description   = params[:description]  if params[:description].present?
              product.stock         = params[:stock]        if params[:stock].present?

              product.save!

              if !!params[:is_api]
                payload = {
                  product: product.as_json,
                  operation: @new_record ? 'create' : 'update'
                }.to_json
                Karafka.producer.produce_sync(topic: 'rails-to-go', payload: payload)
              end
            end
            
            product
          end

          private
          def product
            @product ||= find_or_initialize_product
          end

          def find_or_initialize_product
            existing_product = Product.find_by(id: params[:id])
            if existing_product.present?
              existing_product
            else
              @new_record = true
              Product.new
            end
          end
        end
      end
    end
  end
end
