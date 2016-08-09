'use strict';

var Server = function(uuid,name) {
    this.uuid = uuid;
    this.name = name;
};

var Iso = function(name) {
    this.name = name;
};


var viewModel = {
    servers: ko.observableArray(),
    isos: ko.observableArray(),

    reloadServers: function() {
	self = this;
	self.servers.removeAll();

	window.superagent
	    .get("/servers")
	    .send()
	    .set("Accept: application/json")
	    .end(function(err, res) {
		var ss = res.body.Servers;
		for(var i = 0; i < ss.length; i++) {
		    self.servers.push(new Server(ss[i].Id, ss[i].Name));
		}
	    });
    },

    reloadIsos: function() {
	self = this;
	self.isos.removeAll();

	window.superagent
	    .get("/isos")
	    .send()
	    .set("Accept: application/json")
	    .end(function(err, res) {
		var ss = res.body["iso-images"];
		for(var i = 0; i < ss.length; i++) {
		    self.isos.push(new Iso(ss[i].Name));
		}
	    });
    }
};

window.addEventListener("DOMContentLoaded", function() {
    ko.applyBindings(viewModel);
    viewModel.reloadServers();
    viewModel.reloadIsos();
});
