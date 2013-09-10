for x in ipairs(ghoko.Request.Commits) do
	gitlab.Debugf("[%s] %s", ghoko.Hosting, ghoko.Request.Commits[x].Timestamp)
end
