package main

import (
	"mime"
	"path/filepath"
)

var assetMap = map[string]string{
	"/assets/js/main.js": mainJS,
}

func getAssetContentType(path string) string {
	extension := filepath.Ext(path)
	// Could make lowercase and strip any querystring content of - no need at this time
	return mime.TypeByExtension(extension)
}

const mainJS = `
var app = function() {
  var state = {
    userName: {{.Session.UserName}},
    contacts: [],
    urls: {
      contactsPrefix: "/api/v1/contacts/",
      logIn: "/api/v1/login"
    },
    getContactsUrl: function() {
      return this.urls.contactsPrefix + this.userName;
    },
    getContactUrl: function(id) {
      return this.urls.contactsPrefix + this.userName + "/" + id;
    },
    getContactIndex: function(id) {
      // Could not use indexOf in chrome
      // See https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/findIndex#Browser_compatibility
      for (var index = 0; index < this.contacts.length; index++) {
        if (this.contacts[index].Id == id) { return index; }
      }
      return -1;
    }
  };
 
  var repopulateContactList = function() {
    contactList = document.getElementById("contactList");
    contactList.innerHTML = "";
    state.contacts.forEach(function(contact) {
      var template = document.querySelector("#contactListItemTemplate");
      var content = document.importNode(template.content, true);
      content.querySelector(".contactListItem").setAttribute("id", contact.Id);
      content.querySelector(".contactListItemFirstName").innerText = contact.FirstName || "";
      content.querySelector(".contactListItemLastName").innerText = contact.LastName || "";
      content.querySelector(".contactListItemEdit").onclick = function() { editExistingContact(contact.Id); }
      content.querySelector(".contactListItemDeletion").onclick = function() { removeExistingContact(contact.Id); }
      contactList.appendChild(content);
    });
  };

  var makeLogInAttempt = function() {
    document.getElementById("logIn").disabled = true;
 
    makeApiCall(
      state.urls.logIn,
      "POST",
      {
        UserName: document.getElementById("userName").value,
        Password: document.getElementById("password").value
      },
      function(response) {
        state.userName = document.getElementById("userName").value;

        document.getElementById("userName").value = "";
        document.getElementById("password").value = "";
        document.getElementById("logInSection").style.display = "none";

        document.getElementById("welcomeUserName").innerText = state.userName;
        document.getElementById("welcomeSection").style.display = "block";

        acquireContacts(); 
        document.getElementById("contactListSection").style.display = "block";
      },
      function(status) {
        document.getElementById("logInMessage").innerHTML = "Incorrect user name or password, try again";
        document.getElementById("logIn").disabled = false;
      });
  };

  var makeLogOutAttempt = function() {
    makeApiCall(
      state.urls.logIn,
      "DELETE",
      null,
      function(response) {
        reset();
      },
      function(status) {
        reset();
      });
  };

  var acquireContacts = function() {
    makeApiCall(
      state.getContactsUrl(),
      "GET",
      null,
      function(response) {
        var contacts = JSON.parse(response);
        state.contacts.length = 0;
        if (contacts != null) { contacts.forEach(function(contact) { state.contacts.push(contact); }); }
      },
      function(status) {
        if (status == 401) { reset(); }
        alert("Error encountered " + status);
      });
  };

  var reset = function() {
    state.userName = "";
    state.contacts.length = 0;

    document.getElementById("welcomeSection").style.display = "none";
    document.getElementById("logInSection").style.display = "block";
    document.getElementById("logIn").disabled = false;
    document.getElementById("logInMessage").innerHTML = "";
    
    document.getElementById("contactListSection").style.display = "none";

    if (document.getElementById("contactEditor").open) { document.getElementById("contactEditor").close(); }
  };

  var editNewContact = function() {
    [ "Id", "FirstName", "LastName", "Phones", "Emails", "Twitter", "Notes"].map(function(fieldName) {
        document.getElementById("contact" + fieldName).value = "";
    });
    document.getElementById("contactEditor").showModal();  
    document.getElementById("saveContact").disabled = false;
  };

  var editExistingContact = function(id) {
    var index = state.getContactIndex(id);
    if (index != -1) {
      var contact = state.contacts[index];
      [ "Id", "FirstName", "LastName", "Phones", "Emails", "Twitter", "Notes"].map(function(fieldName) {
        document.getElementById("contact" + fieldName).value = "";
        if (contact[fieldName]) { document.getElementById("contact" + fieldName).value = contact[fieldName]; }
      });
      document.getElementById("contactEditor").showModal();  
      document.getElementById("saveContact").disabled = false;
    }
  };

  var cancelContactEdit = function() {
    document.getElementById("contactEditor").close();  
  };

  var removeExistingContact = function(id) {
    var index = state.getContactIndex(id);
    if (index != -1) {
      makeApiCall(
        state.getContactUrl(id),
        "DELETE",
        null,
        function(response) {
          state.contacts.splice(index, 1);
        },
        function(status) {
          if (status == 401) { reset(); return; }
          alert("Error encountered " + status);
        });
    }
  };
  
  var makeSaveContactAttempt = function() {
    document.getElementById("saveContact").disabled = true;
 
    var isNewContact = document.getElementById("contactId").value == "";
    var contact = {
      Id: document.getElementById("contactId").value,
      FirstName: document.getElementById("contactFirstName").value,
      LastName: document.getElementById("contactLastName").value,
      Emails: [],
      Phones: [],
      Twitter: document.getElementById("contactTwitter").value, 
      Notes: document.getElementById("contactNotes").value
    };

    var url = isNewContact ? state.getContactsUrl() : state.getContactUrl(contact.Id);
    var method = isNewContact ? "POST" : "PUT";

    makeApiCall(
      url,
      method,
      contact,
      function(response, locationHeader) {
        if (isNewContact) {
          contact.Id = locationHeader.substr(locationHeader.lastIndexOf("/") + 1);
          state.contacts.push(contact);
        } else {
          var index = state.getContactIndex(contact.Id);
          state.contacts[index] = contact;
        }
        document.getElementById("contactEditor").close();  
      },
      function(status) {
        if (status == 401) { reset(); return; }
        alert("Error encountered " + status);
      });
  };
  
  var makeApiCall = function(url, method, data, completionFunc, errorFunc) {
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
      if (xhr.readyState == 4) {
        if(xhr.status == 200 || xhr.status == 201) {
          completionFunc(xhr.response, xhr.getResponseHeader('Location'));
        } else {
          errorFunc(xhr.status);
        }
      }
    };

    var dataAsJson = null;
    if (data != null) { dataAsJson = JSON.stringify(data); }

    xhr.open(method, url, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(dataAsJson);
  };

  var app = {};
	app.start = function() {
    // See http://www.html5rocks.com/en/tutorials/es7/observe/
    Array.observe(state.contacts, function(changes) {
      repopulateContactList();
    });

    var loggedIn = (state.userName != "");
    
    document.getElementById("welcomeSection").style.display = loggedIn ? "block": "none";
    document.getElementById("logInSection").style.display = loggedIn ? "none": "block";
    document.getElementById("contactListSection").style.display = "none";  

    document.getElementById("logOut").onclick = makeLogOutAttempt;
    document.getElementById("logIn").onclick = makeLogInAttempt;
    document.getElementById("newContact").onclick = editNewContact;
    document.getElementById("saveContact").onclick = makeSaveContactAttempt;
    document.getElementById("cancelEdit").onclick = cancelContactEdit;
    
    if (loggedIn) { acquireContacts(); }
  };

  return app;
}();

window.addEventListener('DOMContentLoaded', function() { app.start(); });  // Equivalent of jquery document.ready on chrome and IE9+
`
