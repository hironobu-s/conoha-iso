'use strict';

var Server = function(id,name) {
    this.id = id;
    this.name = name;
};

var Iso = function(id,name) {
    this.id = id;
    this.name = name;
};

var viewModel = {
    servers: ko.observableArray(),
    isos: ko.observableArray(),
    nowReloading: ko.observable(false),

    closeNotice: function() {
	document.getElementById("notice").style.display = "none";
    },

    reloadServers: function() {
	self = this;
	self.nowReloading(true);

	window.superagent
	    .get("/servers")
	    .send()
	    .set("Accept: application/json")
	    .end(function(err, res) {
		self.servers.removeAll();

		var ss = res.body.Servers;
		for(var i = 0; i < ss.length; i++) {
		    self.servers.push(new Server(ss[i].Id, ss[i].metadata.instance_name_tag));
		}
		self.nowReloading(false);
	    });
    },

    reloadIsos: function() {
	self = this;
	self.nowReloading(true);

	window.superagent
	    .get("/isos")
	    .send()
	    .set("Accept: application/json")
	    .end(function(err, res) {
		self.isos.removeAll();

		var ss = res.body["iso-images"];
		for(var i = 0; i < ss.length; i++) {
		    self.isos.push(new Iso(ss[i].Id, ss[i].Name));
		}
		self.nowReloading(false);
	    });
    }
};

window.addEventListener("DOMContentLoaded", function() {
    ko.applyBindings(viewModel);
    viewModel.reloadServers();
    viewModel.reloadIsos();
});
