require 'json'

theme = Theme.find_by(id: 1)
abort 'theme 1 not found' unless theme

data = {
  exported_at: Time.now.utc.iso8601,
  id: theme.id,
  name: theme.name,
  user_selectable: theme.user_selectable,
  fields: theme.theme_fields.map do |f|
    h = { name: f.name, target_id: f.target_id, type_id: f.type_id, value: f.value }
    if f.respond_to?(:upload) && f.upload
      h[:upload_url] = f.upload.url
      h[:upload_sha1] = f.upload.sha1
    end
    h
  end,
  settings: theme.theme_settings.map { |s| { name: s.name, value: s.value, data_type: s.data_type } },
  color_schemes: theme.color_schemes.map do |cs|
    {
      name: cs.name,
      base_scheme_id: cs.base_scheme_id,
      colors: cs.color_scheme_colors.map { |c| { name: c.name, hex: c.hex } }
    }
  end,
  child_themes: theme.child_themes.pluck(:id),
  parent_themes: theme.parent_themes.pluck(:id)
}

path = '/tmp/theme-1-backup.json'
File.write(path, JSON.pretty_generate(data))
puts "backup ok bytes=#{File.size(path)} fields=#{data[:fields].size} color_schemes=#{data[:color_schemes].size} settings=#{data[:settings].size}"
