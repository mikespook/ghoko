local repo = ghoko.Params["repository"]
ghoko.Debugf("%s %s %s", ghoko.Id, repo["name"], repo["url"])
for k, v in ipairs(ghoko.Params["commits"]) do
	ghoko.Debugf("[%s] %s \"%s\"", v["timestamp"], v["id"], v["message"])
end
