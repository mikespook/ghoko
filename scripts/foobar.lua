ghoko.Debugf("%s %+v", ghoko.Id, ghoko.Params)
if ghoko.Params["test"] == nil then
	ghoko.Params["test"] = "abc"
	ghoko.Call(ghoko.Id, "foobar", ghoko.Params)
end
