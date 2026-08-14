require 'json'

data = JSON.parse(File.read('/tmp/restore-theme1.json'))
theme = Theme.find_by(id: data['id'])
abort "theme #{data['id']} not found" unless theme

field_updated = 0
field_created = 0
data.fetch('fields', []).each do |fh|
  f = theme.theme_fields.find_by(name: fh['name'], target_id: fh['target_id'])
  if f
    f.value = fh['value']
    f.value_baked = nil
    f.save!
    field_updated += 1
  else
    ThemeField.create!(theme: theme, name: fh['name'], target_id: fh['target_id'], type_id: fh['type_id'], value: fh['value'])
    field_created += 1
  end
end

data.fetch('settings', []).each do |sh|
  s = theme.theme_settings.find_by(name: sh['name'])
  if s
    s.value = sh['value']
    s.data_type = sh['data_type'] if sh['data_type']
    s.save!
  end
end

data.fetch('color_schemes', []).each do |csh|
  cs = theme.color_schemes.find_by(name: csh['name'])
  if cs.nil?
    begin
      cs = ColorScheme.create!(name: csh['name'], base_scheme_id: csh['base_scheme_id'], theme: theme)
    rescue StandardError => e
      puts "warn: cannot recreate color scheme #{csh['name']}: #{e.message}"
      next
    end
  end
  csh.fetch('colors', []).each do |ch|
    c = ColorSchemeColor.find_or_initialize_by(color_scheme_id: cs.id, name: ch['name'])
    c.hex = ch['hex']
    c.save!
  end
end

theme.theme_fields.where(type_id: 1).each do |f|
  f.value_baked = nil
  f.save!
  f.ensure_baked!
end

puts "restore ok fields_updated=#{field_updated} fields_created=#{field_created} color_schemes=#{theme.color_schemes.count}"
