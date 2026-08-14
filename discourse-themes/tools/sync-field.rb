require 'json'
require 'digest'

meta = JSON.parse(File.read('/tmp/sync-field-meta.json'))
value = File.read('/tmp/sync-field-value.txt', encoding: 'UTF-8')
theme = Theme.find_by(id: 1)
abort 'theme 1 not found' unless theme

f = theme.theme_fields.find_by(name: meta['name'], target_id: meta['target_id'])
abort "field not found: name=#{meta['name']} target=#{meta['target_id']}" unless f

f.value = value
f.value_baked = nil
f.save!
f.ensure_baked!

# Anonymous page cache keys embed the stylesheet link, not the digest, so drop
# them to make visitors pick up the new stylesheet link immediately.
# In-memory digest caches in web workers are reset by the unicorn restart
# performed by sync-field.ps1 right after this runner exits.
Discourse.redis.scan_each(match: 'ANON_CACHE_*') { |k| Discourse.redis.del(k) }

puts "sync ok name=#{f.name} target=#{f.target_id} bytes=#{f.value.bytesize} sha1=#{Digest::SHA1.hexdigest(f.value)} baked=#{f.value_baked.present?}"