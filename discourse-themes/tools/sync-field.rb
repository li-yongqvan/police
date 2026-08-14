require 'json'

data = JSON.parse(File.read('/tmp/sync-field.json'))
theme = Theme.find_by(id: 1)
abort 'theme 1 not found' unless theme

f = theme.theme_fields.find_by(name: data['name'], target_id: data['target_id'])
abort "field not found: name=#{data['name']} target=#{data['target_id']}" unless f

f.value = data['value']
f.value_baked = nil
f.save!
f.ensure_baked!
puts "sync ok name=#{f.name} target=#{f.target_id} bytes=#{f.value.bytesize} baked=#{f.value_baked.present?}"
