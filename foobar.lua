for x in ipairs(gitlab.Request.Commits) do
	gitlab.Debugf("%s", gitlab.Request.Commits[x].Timestamp)
end
