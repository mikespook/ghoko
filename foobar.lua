-- http://192.168.1.100/gitlab/foobar?secret=yoursecrethere&test=Thequickbrownfoxjumpsoverthelazydog

for x in ipairs(ghoko.Request.Commits) do
	ghoko.Debugf("[%s] %s", ghoko.Host, ghoko.Request.Commits[x].Timestamp)
end

ghoko.Debug(ghoko.Params['test'])
